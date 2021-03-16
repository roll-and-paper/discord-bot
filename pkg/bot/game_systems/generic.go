package game_systems

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/dices/roller"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/game_system/generic"
	"strings"
)

type Generic struct{}

func (g *Generic) Roll(msg *discordgo.Message, sess *discordgo.Session, args []string, fallback GameSystem) (*roller.Result, error) {
	cmd := strings.ToLower(strings.Join(args, " "))
	return generic.Roll(cmd)
}
