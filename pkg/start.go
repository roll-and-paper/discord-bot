package pkg

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dohr-michael/roll-and-paper-bot/config"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/bot"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/models"
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

func onMessageHandler(col storage.Storage, lock *discord.Lock) func(*discordgo.Session, *discordgo.MessageCreate) {
	serv := bot.New(col)
	return func(sess *discordgo.Session, msg *discordgo.MessageCreate) {
		log.Printf("on message")

		if msg.Author.ID == sess.State.User.ID {
			return
		}
		ctx := context.Background()
		lockId := fmt.Sprintf("%s-%s", msg.ID, "on_message")
		if !lock.Received(lockId, ctx) {
			log.Printf("message already processed %s", msg.ID)
			return
		}
		defer func() { _ = lock.End(lockId, ctx) }()
		start := time.Now()
		defer func() {
			discordPromCollector.WithLabelValues("on_message").Observe(float64(time.Now().Sub(start).Milliseconds()))
		}()
		state := models.NewServerState(msg.GuildID)

		if err := col.FindOneOrCreate(msg.GuildID, state, ctx); err != nil {
			log.Printf("error when fetching guild state %v", err)
			return
		}

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

	dis.AddHandler(onMessageHandler(s, lock))

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
