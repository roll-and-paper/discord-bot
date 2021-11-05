package services

import (
	"fmt"
	"github.com/dohr-michael/roll-and-paper-bot/i18n"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/models"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/services/roll"
)

func (s *Services) Roll(state *models.ServerState, params roll.Params) (string, error) {
	switch d := params.(type) {
	case *roll.VampireDarkAgesParams:
		return roll.VampireDarkAges(d)
	case *roll.GenericParams:
		return roll.Generic(d)
	default:
		return "", fmt.Errorf(i18n.Must(state.Language, "errors.commands.roll.unknown-game-system", nil))
	}
}
