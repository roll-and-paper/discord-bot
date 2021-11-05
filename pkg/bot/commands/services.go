package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/models"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/services"
)

type Services struct {
	underlying *services.Services
	handlers   map[string]func(*discordgo.Session, *discordgo.InteractionCreate, *models.ServerState)
}

func NewServices(underlying *services.Services) *Services {
	result := &Services{underlying: underlying}
	result.handlers = map[string]func(*discordgo.Session, *discordgo.InteractionCreate, *models.ServerState){
		RollName: result.handleRoll,
	}
	return result
}
