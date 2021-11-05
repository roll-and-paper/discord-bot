package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/models"
)

func (*Services) printResult(sess *discordgo.Session, evt *discordgo.InteractionCreate, state *models.ServerState, result string, err error) {
	if err != nil {
		return
	}
	_ = sess.InteractionRespond(evt.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: result,
		},
	})
}
