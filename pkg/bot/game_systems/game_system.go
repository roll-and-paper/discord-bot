package game_systems

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/dices/roller"
)

type GameSystem interface {
	Roll(msg *discordgo.Message, sess *discordgo.Session, args []string, fallback GameSystem) (*roller.Result, error)
}
