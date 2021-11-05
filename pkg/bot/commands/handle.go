package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/models"
)

func (s *Services) Handle(sess *discordgo.Session, evt *discordgo.InteractionCreate, state *models.ServerState) {
	switch evt.Type {
	case discordgo.InteractionApplicationCommand:
		if fn, ok := s.handlers[evt.ApplicationCommandData().Name]; ok {
			fn(sess, evt, state)
		}
	case discordgo.InteractionMessageComponent:
		//if fn, ok := components.Message[evt.MessageComponentData().CustomID]; ok {
		//	fn(session, evt, state, serv)
		//}
	}
}
