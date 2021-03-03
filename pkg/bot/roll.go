package bot

import (
	"github.com/bwmarrin/discordgo"
	game_systems "github.com/dohr-michael/roll-and-paper-bot/pkg/bot/game_systems"
	gp "github.com/dohr-michael/roll-and-paper-bot/pkg/models"
	"github.com/dohr-michael/roll-and-paper-bot/tools/discord"
	"github.com/thoas/go-funk"
	"strings"
)

var rollSystems map[string]game_systems.GameSystem
var allRollSystems string
var genericRollSystem = &game_systems.Generic{}

func init() {
	rollSystems = map[string]game_systems.GameSystem{
		"vampire-dark-ages": &game_systems.VampireSystem{},
	}
	allRollSystems = strings.Join(funk.Keys(rollSystems).([]string), ", ")
}

func (s *Services) Roll(msg *discordgo.Message, sess *discordgo.Session, state *gp.ServerState, args []string) error {
	gs, ok := rollSystems[state.Config.GameSystem]
	if !ok {
		gs = genericRollSystem
	}
	result, err := gs.Roll(msg, sess, args, genericRollSystem)
	if err != nil {
		return err
	}
	if result != nil {
		_, err = discord.ReplyToWithContent(msg, sess, result.String())
		return err
	}
	return nil
}
