package game_systems

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/dices/roller"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/game_system/vampire"
	"strconv"
	"strings"
)

type VampireSystem struct{}

func (v *VampireSystem) Roll(msg *discordgo.Message, sess *discordgo.Session, args []string, fallback GameSystem) (*roller.Result, error) {
	cmd := strings.ToLower(strings.Join(args, " "))
	options := strings.Split(cmd, ",")
	found := false
	dices, diff, spe := 0, 0, 0
	var err error
	for _, item := range options {
		trim := strings.TrimSpace(item)
		if strings.HasPrefix(trim, "diff:") || strings.HasPrefix(trim, "diff=") {
			diff, err = strconv.Atoi(strings.TrimPrefix(strings.TrimPrefix(trim, "diff:"), "diff="))
			if err != nil {
				diff = 0
			}
		}
		if strings.HasPrefix(trim, "spe:") || strings.HasPrefix(trim, "spe=") {
			spe, err = strconv.Atoi(strings.TrimPrefix(strings.TrimPrefix(trim, "spe:"), "spe="))
			if err != nil {
				spe = 0
			}
		}
		if strings.HasSuffix(trim, "d") {
			found = true
			dices, err = strconv.Atoi(strings.TrimSuffix(trim, "d"))
			if err != nil {
				dices = 0
			}
		}
	}
	if !found {
		return fallback.Roll(msg, sess, args, nil)
	}
	return vampire.Roll(dices, diff, spe), nil
}
