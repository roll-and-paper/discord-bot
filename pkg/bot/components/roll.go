package components

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/dices/roller"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/game_system/generic"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/game_system/vampire"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/models"
)

var rollSystems map[string]func(sess *discordgo.Session, evt *discordgo.InteractionCreate, data *discordgo.ApplicationCommandInteractionDataOption, state *models.ServerState)

func init() {
	rollSystems = map[string]func(sess *discordgo.Session, evt *discordgo.InteractionCreate, data *discordgo.ApplicationCommandInteractionDataOption, state *models.ServerState){
		"vampire-dark-ages": rollVampireDarkAges,
	}
}

func Roll(sess *discordgo.Session, evt *discordgo.InteractionCreate, state *models.ServerState) {
	data := evt.ApplicationCommandData().Options[0]
	fn, ok := rollSystems[data.Name]
	if !ok {
		fn = rollGeneric
	}
	fn(sess, evt, data, state)
}

func printResult(sess *discordgo.Session, evt *discordgo.InteractionCreate, state *models.ServerState, result *roller.Result, err error) {
	if err != nil {
		// TODO Print errors
		return
	}
	_ = sess.InteractionRespond(evt.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: result.String(),
		},
	})
}

func rollGeneric(sess *discordgo.Session, evt *discordgo.InteractionCreate, data *discordgo.ApplicationCommandInteractionDataOption, state *models.ServerState) {
	cmd, ok := data.Options[0].Value.(string)
	if !ok {
		printResult(sess, evt, state, nil, fmt.Errorf("command must be string"))
		return
	}
	res, err := generic.Roll(cmd)
	printResult(sess, evt, state, res, err)
}

func rollVampireDarkAges(sess *discordgo.Session, evt *discordgo.InteractionCreate, data *discordgo.ApplicationCommandInteractionDataOption, state *models.ServerState) {
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
	res := vampire.Roll(int(dices), int(diff), int(spe))
	printResult(sess, evt, state, res, nil)
}
