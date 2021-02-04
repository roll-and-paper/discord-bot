package pkg

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/dohr-michael/roll-and-paper-bot/config"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/models"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/services"
	"github.com/dohr-michael/roll-and-paper-bot/tools/discord"
	"github.com/dohr-michael/roll-and-paper-bot/tools/storage"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"log"
	"strings"
	"time"
)

func onMessageHandler(col storage.Storage) func(*discordgo.Session, *discordgo.MessageCreate) {
	serv := services.New(col)
	return func(sess *discordgo.Session, msg *discordgo.MessageCreate) {
		log.Printf("on message")

		if msg.Author.ID == sess.State.User.ID {
			return
		}
		ctx := context.Background()
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

	dis, err := discordgo.New("Bot " + config.DiscordToken())
	if err != nil {
		return err
	}

	dis.AddHandler(onMessageHandler(s))

	dis.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions

	err = dis.Open()
	if err != nil {
		return err
	}
	defer func() { _ = dis.Close() }()

	r := gin.Default()

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
