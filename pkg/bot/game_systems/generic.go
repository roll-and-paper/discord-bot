package game_systems

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/dices/parser"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/dices/roller"
	"strings"
)

type Generic struct{}

func (g *Generic) Roll(msg *discordgo.Message, sess *discordgo.Session, args []string, fallback GameSystem) (*roller.Result, error) {
	cmd := strings.ToLower(strings.Join(args, " "))
	r, err := parser.Parse(cmd)
	if err != nil {
		return nil, err
	}
	return r.Roll(cmd), nil
}
