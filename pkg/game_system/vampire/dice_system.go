package vampire

import (
	"fmt"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/dices"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/dices/roller"
)

func withDifficulty(difficulty int) roller.ResultOption {
	return roller.Count(func(result roller.DiceResult) bool { return result.Value >= difficulty })
}

var (
	is1  = roller.Count(func(result roller.DiceResult) bool { return result.Value == 1 })
	is10 = roller.Count(func(result roller.DiceResult) bool { return result.Value == 10 })
)

func Roll(count, difficulty, nbSpecs int) *roller.Result {
	r := roller.FromDice(dices.D10, count, []roller.Option{withDifficulty(difficulty)})
	roll := r.Roll(fmt.Sprintf("%dd : count(>= %d) + count(= 10) x %d - count(= 0)", count, difficulty, nbSpecs))
	return roll.
		Minus(roller.FromResultReference(roll, []roller.Option{is1}).Roll("")).
		Add(roller.FromResultReference(roll, []roller.Option{is10}).Roll("").Mul(roller.NewResult(float64(nbSpecs), nil, "")))
}
