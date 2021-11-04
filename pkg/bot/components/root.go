package components

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/models"
)

var Application = map[string]func(sess *discordgo.Session, evt *discordgo.InteractionCreate, state *models.ServerState){
	"roll": Roll,
}
