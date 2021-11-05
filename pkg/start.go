package pkg

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dohr-michael/roll-and-paper-bot/config"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/bot/commands"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/bot/components"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/bot/legacy"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/models"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/services"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/shared"
	"github.com/dohr-michael/roll-and-paper-bot/tools/discord"
	"github.com/dohr-michael/roll-and-paper-bot/tools/storage"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"log"
	"strings"
	"time"
)

var (
	discordPromCollector = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "discord_durations_seconds",
			Help:       "Discord duration distributions.",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}, []string{"service"},
	)
)

func init() {
	prometheus.MustRegister(discordPromCollector)
	prometheus.MustRegister(prometheus.NewBuildInfoCollector())
}

func onEvent(id, guildId, command string, col storage.Storage, lock *discord.Lock, fn func(state *models.ServerState)) {
	log.Printf("process %s", command)
	ctx := context.Background()
	lockId := fmt.Sprintf("%s-%s", id, command)
	if !lock.Received(lockId, ctx) {
		log.Printf("%s already processed %s", command, id)
		return
	}
	defer func() { _ = lock.End(lockId, ctx) }()
	start := time.Now()
	defer func() {
		discordPromCollector.WithLabelValues(command).Observe(float64(time.Now().Sub(start).Milliseconds()))
	}()
	state := models.NewServerState(guildId)

	if err := col.FindOneOrCreate(guildId, state, ctx); err != nil {
		log.Printf("error when fetching guild state %v", err)
		return
	}
	fn(state)
}

func onReadyHandler(cmd *commands.Services, col storage.Storage, lock *discord.Lock) func(*discordgo.Session, *discordgo.Ready) {
	return func(session *discordgo.Session, evt *discordgo.Ready) {
		for _, g := range evt.Guilds {
			onEvent(fmt.Sprintf("%s-%s", shared.Revision, g.ID), g.ID, fmt.Sprintf("init_commands_%s", g.ID), col, lock, func(state *models.ServerState) {
				cmd.Register(session, state)
			})
		}
	}
}

func onInteractionHandler(cmd *commands.Services, col storage.Storage, lock *discord.Lock) func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(session *discordgo.Session, evt *discordgo.InteractionCreate) {
		onEvent(evt.ID, evt.GuildID, "on_interaction", col, lock, func(state *models.ServerState) {
			cmd.Handle(session, evt, state)
			switch evt.Type {
			case discordgo.InteractionApplicationCommand:
				if fn, ok := components.Application[evt.ApplicationCommandData().Name]; ok {
					fn(session, evt, state)
				}
			case discordgo.InteractionMessageComponent:
				//if fn, ok := components.Message[evt.MessageComponentData().CustomID]; ok {
				//	fn(session, evt, state, serv)
				//}
			}
		})
	}
}

func onMessageHandler(serv *legacy.Services, col storage.Storage, lock *discord.Lock) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(sess *discordgo.Session, msg *discordgo.MessageCreate) {
		log.Printf("on message")

		if msg.Author.ID == sess.State.User.ID {
			return
		}

		onEvent(msg.ID, msg.GuildID, "on_message", col, lock, func(state *models.ServerState) {
			if !strings.HasPrefix(msg.Content, state.Config.Prefix) {
				return
			}
			cmd := strings.Split(strings.TrimPrefix(msg.Content, state.Config.Prefix), " ")
			if err := serv.Apply(msg.Message, sess, state, cmd[0], cmd[1:]...); err != nil {
				log.Printf("error : %v", err)
				if _, err := discord.SendMessage(msg.ChannelID, sess, state.Language, "messages.oops", map[string]string{"Details": err.Error()}); err != nil {
					log.Printf("failed to send message 'messages.oops' to %s", msg.GuildID)
				}

			}
		})
	}
}

func initMongodb() (*mongo.Database, func(), error) {
	cstring, err := connstring.Parse(config.MongoUri())
	if err != nil {
		return nil, nil, err
	}
	muri := options.Client().ApplyURI(cstring.String())
	mongoClient, err := mongo.NewClient(muri)
	if err != nil {
		return nil, nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = mongoClient.Connect(ctx)
	if err != nil {
		return nil, nil, err
	}
	database := mongoClient.Database(cstring.Database)
	return database, func() {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		_ = mongoClient.Disconnect(ctx)
	}, nil
}

func Start() error {
	database, closeDatabase, err := initMongodb()
	if err != nil {
		return err
	}
	defer closeDatabase()
	s := storage.NewMongoDBStorage(database, "discord_guilds")
	lock := discord.NewLock(database)

	dis, err := discordgo.New("Bot " + config.DiscordToken())
	if err != nil {
		return err
	}

	srv := services.NewServices(s)
	cmd := commands.NewServices(srv)
	leg := legacy.New(srv, s)

	dis.AddHandler(func(s *discordgo.Session, r *discordgo.GuildCreate) {
		log.Printf("guild join %s", r.ID)
	})
	dis.AddHandler(func(s *discordgo.Session, r *discordgo.GuildDelete) {
		log.Printf("guild leave %s", r.ID)
	})

	dis.AddHandler(onReadyHandler(cmd, s, lock))
	dis.AddHandler(onInteractionHandler(cmd, s, lock))
	dis.AddHandler(onMessageHandler(leg, s, lock))

	dis.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions

	err = dis.Open()
	if err != nil {
		return err
	}
	defer func() { _ = dis.Close() }()

	r := gin.Default()
	r.GET("/@/metrics", gin.WrapH(promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	})))

	r.GET("/@/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	r.GET("/@/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	return r.Run()
}
