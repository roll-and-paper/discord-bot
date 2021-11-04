package bot

import (
	"github.com/bwmarrin/discordgo"
	gp "github.com/dohr-michael/roll-and-paper-bot/pkg/models"
	"github.com/dohr-michael/roll-and-paper-bot/tools/discord"
	"log"
)

func (s *Services) Apply(msg *discordgo.Message, sess *discordgo.Session, state *gp.ServerState, cmd string, args ...string) error {
	log.Printf("key: %s, args : %v", cmd, args)
	switch cmd {
	case "ping":
		_, err := discord.ReplyTo(msg, sess, state.Language, "pong", state)
		return err
	case "init":
		return s.Init(msg, sess, state)
	case "set":
		return s.Set(msg, sess, state, args)
	case "help":
		_, err := discord.SendMessage(msg.ChannelID, sess, state.Language, "messages.help", state)
		return err
	default: // Default action is trying to roll
		return s.Roll(msg, sess, state, append([]string{cmd}, args...))
	}
}
