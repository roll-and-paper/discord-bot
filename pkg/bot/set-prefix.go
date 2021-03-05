package bot

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/models"
	"github.com/dohr-michael/roll-and-paper-bot/tools/discord"
)

func setPrefix(s *Services, msg *discordgo.Message, sess *discordgo.Session, state *models.ServerState, args []string) error {
	if len(args) == 0 {
		_, err := discord.SendMessage(msg.ChannelID, sess, state.Language, "errors.commands.set.prefix.missing", state)
		return err
	}

	if err := s.Save(state, Fields{"config.prefix": args[0]}, context.Background()); err != nil {
		return err
	}

	return sess.MessageReactionAdd(msg.ChannelID, msg.ID, "👌")
}
