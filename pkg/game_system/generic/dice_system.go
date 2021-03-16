package generic

import (
	"github.com/dohr-michael/roll-and-paper-bot/pkg/dices/parser"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/dices/roller"
)

func Roll(expression string) (*roller.Result, error) {
	r, err := parser.Parse(expression)
	if err != nil {
		return nil, err
	}
	return r.Roll(expression), nil
}
