package roller

import (
	"github.com/dohr-michael/roll-and-paper-bot/pkg/dices"
	"github.com/thoas/go-funk"
	"sort"
)

type Option interface{}

type DiceResultOption func(DiceResult, dices.Dice) DiceResult

func Explode(validator Validator) DiceResultOption {
	return func(result DiceResult, dice dices.Dice) DiceResult {
		if validator(dice, result) {
			results := []int{result.Value}
			for validator(dice, DiceResult{results[len(results)-1], nil}) {
				results = append(results, dice.Roll())
			}
			return DiceResult{funk.SumInt(results), results}
		}
		return result
	}
}

type ResultOption func(*Result) *Result

type SortDirection string

const (
	Asc  SortDirection = "asc"
	Desc SortDirection = "desc"
)

func Sort(direction SortDirection) ResultOption {
	return func(r *Result) *Result {
		if r == nil {
			return nil
		}
		if len(r.Dices) == 0 {
			return r
		}
		d := make([]DiceResult, len(r.Dices))
		for idx, c := range r.Dices {
			d[idx] = c
		}
		sort.Slice(d, func(i, j int) bool {
			if direction == Asc {
				return d[i].Value < d[j].Value
			}
			return d[i].Value > d[j].Value
		})
		return &Result{
			Value:      r.Value,
			Dices:      d,
			Expression: r.Expression,
		}
	}
}

func Count(predicate func(result DiceResult) bool) ResultOption {
	return func(r *Result) *Result {
		if r == nil {
			return nil
		}
		if len(r.Dices) == 0 {
			return r
		}
		result := 0
		for _, c := range r.Dices {
			if predicate(c) {
				result++
			}
		}
		return &Result{
			Value:          float64(result),
			Dices:          r.Dices,
			Expression:     r.Expression,
		}
	}
}
