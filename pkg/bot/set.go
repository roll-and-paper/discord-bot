package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	gp "github.com/dohr-michael/roll-and-paper-bot/pkg/models"
	"github.com/dohr-michael/roll-and-paper-bot/tools/discord"
	"log"
	"strings"
)

var setCommands map[string]func(*Services, *discordgo.Message, *discordgo.Session, *gp.ServerState, []string) error
var allSetCommands string

func init() {
	setCommands = map[string]func(*Services, *discordgo.Message, *discordgo.Session, *gp.ServerState, []string) error{
		"master":      setMaster,
		"player":      setPlayer,
		"game-system": setGameSystem,
	}
	tmp := make([]string, 0)
	for key := range setCommands {
		tmp = append(tmp, key)
	}
	allSetCommands = strings.Join(tmp, ", ")
}

func (s *Services) Set(msg *discordgo.Message, sess *discordgo.Session, state *gp.ServerState, args []string) error {
	if len(args) == 0 {
		_, err := discord.SendMessage(msg.ChannelID, sess, state.Language, "errors.commands.set.bad-command", map[string]interface{}{"Cmd": "", "AllCmd": allSetCommands})
		return err
	}
	subCommand := args[0]
	params := args[1:]
	fn, ok := setCommands[subCommand]
	if !ok {
		_, err := discord.SendMessage(msg.ChannelID, sess, state.Language, "errors.commands.set.bad-command", map[string]interface{}{"Cmd": subCommand, "AllCmd": allSetCommands})
		return err
	}
	if msg.ChannelID != state.Channels.WithBot {
		_, err := discord.SendMessage(msg.ChannelID, sess, state.Language, "errors.commands.set.unauthorized", nil)
		return err
	}
	log.Printf(" - start 'set %s'", subCommand)
	if params[0] == "--help" {
		_, err := discord.SendMessage(msg.ChannelID, sess, state.Language, fmt.Sprintf("messages.set.%s.help", subCommand), state)
		return err
	}

	return fn(s, msg, sess, state, params)
}
