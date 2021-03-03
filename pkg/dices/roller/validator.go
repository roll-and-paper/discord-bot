package roller

import "github.com/dohr-michael/roll-and-paper-bot/pkg/dices"

type Validator func(dices.Dice, DiceResult) bool

func MaxValue(dice dices.Dice, result DiceResult) bool {
	return len(result.Exploded) == 0 && result.Value == dice.MaxValue()
}

func HasValue(value int) Validator {
	return func(dice dices.Dice, result DiceResult) bool {
		return len(result.Exploded) == 0 && result.Value == value
	}
}
