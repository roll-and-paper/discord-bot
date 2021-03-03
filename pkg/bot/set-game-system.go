package bot

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/models"
	"github.com/dohr-michael/roll-and-paper-bot/tools/discord"
	"github.com/thoas/go-funk"
)

func setGameSystem(s *Services, msg *discordgo.Message, sess *discordgo.Session, state *models.ServerState, args []string) error {
	if len(args) == 0 {
		_, err := discord.SendMessage(msg.ChannelID, sess, state.Language, "errors.commands.set.game-system.missing", state)
		return err
	}
	wantedGameSystem := args[0]
	if !funk.Contains(funk.Keys(rollSystems), wantedGameSystem) {
		_, err := discord.SendMessage(msg.ChannelID, sess, state.Language, "errors.commands.set.game-system.unknown-game-system", map[string]string{"Wanted": wantedGameSystem, "AllSystems": allRollSystems})
		return err
	}

	if err := s.Save(state, Fields{"config.gameSystem": wantedGameSystem}, context.Background()); err != nil {
		return err
	}

	return sess.MessageReactionAdd(msg.ChannelID, msg.ID, "👌")
}
