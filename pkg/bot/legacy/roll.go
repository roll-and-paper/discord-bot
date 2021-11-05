package legacy

import (
	"github.com/bwmarrin/discordgo"
	gp "github.com/dohr-michael/roll-and-paper-bot/pkg/models"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/services/roll"
	"github.com/dohr-michael/roll-and-paper-bot/tools/discord"
	"strconv"
	"strings"
)

func parseVampireDarkAges(cmd string) roll.Params {
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
		return &roll.GenericParams{Expression: cmd}
	} else {
		return &roll.VampireDarkAgesParams{
			Dices:          dices,
			Difficulty:     diff,
			Specialisation: spe,
		}
	}
}

func (s *Services) Roll(msg *discordgo.Message, sess *discordgo.Session, state *gp.ServerState, args []string) error {
	cmd := strings.ToLower(strings.Join(args, " "))
	var params roll.Params
	switch state.Config.GameSystem {
	case "vampire-dark-ages":
		params = parseVampireDarkAges(cmd)
	default:
		params = &roll.GenericParams{Expression: cmd}
	}
	result, err := s.underlying.Roll(state, params)

	if err != nil {
		return err
	}

	_, err = discord.ReplyToWithContent(msg, sess, result)
	return err
}
