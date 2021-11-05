package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dohr-michael/roll-and-paper-bot/config"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/models"
	"log"
)

func (s *Services) Register(dis *discordgo.Session, state *models.ServerState) {
	cmds := []*discordgo.ApplicationCommand{
		s.roll(state.Language),
	}

	registeredCommands, _ := dis.ApplicationCommands(config.DiscordAppId(), state.Id)
	if registeredCommands != nil {
		for _, current := range registeredCommands {
			log.Printf("unregister '%s' into guild %s", current.Name, state.Id)
			if err := dis.ApplicationCommandDelete(config.DiscordAppId(), state.Id, current.ID); err != nil {
				log.Printf("cannot delete command %s, %s", current.Name, err.Error())
			}
		}
	}
	for _, cmd := range cmds {
		log.Printf("register command '%s' into guild '%s'", cmd.Name, state.Id)
		if _, err := dis.ApplicationCommandCreate(config.DiscordAppId(), state.Id, cmd); err != nil {
			log.Printf("cannot create command '%s' :  %s", cmd.Name, err.Error())
		}
	}
}
