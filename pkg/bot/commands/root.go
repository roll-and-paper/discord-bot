package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dohr-michael/roll-and-paper-bot/config"
	"github.com/dohr-michael/roll-and-paper-bot/i18n"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/models"
	"log"
)

func commands(lang string) []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		//{
		//	Name:        "init",
		//	Description: i18n.Must(lang, "commands.description.init", nil),
		//},
		{
			Name:        "roll",
			Description: i18n.Must(lang, "commands.description.roll", nil),
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "generic",
					Description: i18n.Must(lang, "commands.roll.generic.description", nil),
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
							Name:        "query",
							Description: i18n.Must(lang, "commands.roll.generic.options.query.description", nil),
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "vampire-dark-ages",
					Description: i18n.Must(lang, "commands.roll.vampire-dark-ages.description", nil),
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    true,
							Name:        "dices",
							Description: i18n.Must(lang, "commands.roll.vampire-dark-ages.options.dices.description", nil),
						},
						{
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    false,
							Name:        "difficulty",
							Description: i18n.Must(lang, "commands.roll.vampire-dark-ages.options.difficulty.description", nil),
						},
						{
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    false,
							Name:        "specialisation",
							Description: i18n.Must(lang, "commands.roll.vampire-dark-ages.options.specialisation.description", nil),
						},
					},
				},
			},
		},
		//{
		//	Name:        "buttons",
		//	Description: "Test the buttons if you got courage",
		//},
		//{
		//	Name: "selects",
		//	Options: []*discordgo.ApplicationCommandOption{
		//		{
		//			Type:        discordgo.ApplicationCommandOptionSubCommand,
		//			Name:        "multi",
		//			Description: "Multi-item select menu",
		//		},
		//		{
		//			Type:        discordgo.ApplicationCommandOptionSubCommand,
		//			Name:        "single",
		//			Description: "Single-item select menu",
		//		},
		//	},
		//	Description: "Lo and behold: dropdowns are coming",
		//},
	}
}

func Register(dis *discordgo.Session, state *models.ServerState) {
	registeredCommands, _ := dis.ApplicationCommands(config.DiscordAppId(), state.Id)
	if registeredCommands != nil {
		for _, current := range registeredCommands {
			log.Printf("unregister '%s' into guild %s", current.Name, state.Id)
			if err := dis.ApplicationCommandDelete(config.DiscordAppId(), state.Id, current.ID); err != nil {
				log.Printf("cannot delete command %s, %s", current.Name, err.Error())
			}
		}
	}
	for _, cmd := range commands(state.Language) {
		log.Printf("register command '%s' into guild '%s'", cmd.Name, state.Id)
		if _, err := dis.ApplicationCommandCreate(config.DiscordAppId(), state.Id, cmd); err != nil {
			log.Printf("cannot create command '%s' :  %s", cmd.Name, err.Error())
		}
	}
}
