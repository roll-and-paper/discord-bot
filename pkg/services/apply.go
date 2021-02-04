package services

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/models"
	"github.com/dohr-michael/roll-and-paper-bot/tools/discord"
	"log"
)

func (s *Services) Apply(msg *discordgo.Message, sess *discordgo.Session, state *models.ServerState, cmd string, args ...string) error {
	log.Printf("key: %s, args : %v", cmd, args)
	switch cmd {
	case "ping":
		_, err := discord.ReplyTo(msg, sess, state.Language, "pong", state)
		return err
	case "init":
		return s.Init(msg, sess, state)
	case "set":
		return s.Set(msg, sess, state, args)
	}

	return nil
}
