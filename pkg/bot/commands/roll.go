package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dohr-michael/roll-and-paper-bot/i18n"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/models"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/services/roll"
)

const RollName = "roll"

func (*Services) roll(lang string) *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        RollName,
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
						Name:        "expression",
						Description: i18n.Must(lang, "commands.roll.generic.options.expression.description", nil),
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
	}
}

func (s *Services) handleRoll(sess *discordgo.Session, evt *discordgo.InteractionCreate, state *models.ServerState) {
	data := evt.ApplicationCommandData().Options[0]
	var params roll.Params
	switch data.Name {
	case "generic":
		params = &roll.GenericParams{Expression: data.Options[0].StringValue()}
	case "vampire-dark-ages":
		var dices, diff, spe int64
		for _, c := range data.Options {
			switch c.Name {
			case "dices":
				dices = c.IntValue()
			case "difficulty":
				diff = c.IntValue()
			case "specialisation":
				spe = c.IntValue()
			}
		}
		params = &roll.VampireDarkAgesParams{
			Dices:          int(dices),
			Difficulty:     int(diff),
			Specialisation: int(spe),
		}
	}
	if params == nil {
		s.printResult(sess, evt, state, "", fmt.Errorf(i18n.Must(state.Language, "errors.commands.roll.unknown-game-system", nil)))
		return
	}
	res, err := s.underlying.Roll(state, params)
	s.printResult(sess, evt, state, res, err)
}
