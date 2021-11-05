package roll

import (
	"github.com/dohr-michael/roll-and-paper-bot/pkg/game_system/generic"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/game_system/vampire"
)

func Generic(params *GenericParams) (string, error) {
	res, err := generic.Roll(params.Expression)
	if err != nil {
		return "", err
	}
	return res.String(), nil
}

func VampireDarkAges(params *VampireDarkAgesParams) (string, error) {
	return vampire.Roll(params.Dices, params.Difficulty, params.Specialisation).String(), nil
}
