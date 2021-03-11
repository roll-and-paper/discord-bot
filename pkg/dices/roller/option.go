package roller

import (
	"github.com/dohr-michael/roll-and-paper-bot/pkg/dices"
	"github.com/thoas/go-funk"
	"sort"
)

type Option interface{}

type DiceResultOption interface {
	Option
	Apply(DiceResult, dices.Dice) DiceResult
}

type explodeOpt struct {
	validator Validator
}

func (s *explodeOpt) Apply(result DiceResult, dice dices.Dice) DiceResult {
	if s.validator(dice, result) {
		results := []int{result.Value}
		for s.validator(dice, DiceResult{results[len(results)-1], nil}) {
			results = append(results, dice.Roll())
		}
		return DiceResult{funk.SumInt(results), results}
	}
	return result
}

func Explode(validator Validator) DiceResultOption {
	return &explodeOpt{validator}
}

type ResultOption interface {
	Option
	Apply(*Result) *Result
}

// Sort

type SortDirection string

const (
	Asc  SortDirection = "asc"
	Desc SortDirection = "desc"
)

type sortOpt struct {
	direction SortDirection
}

var (
	sortOptAsc  = &sortOpt{Asc}
	sortOptDesc = &sortOpt{Desc}
)

func (s *sortOpt) Apply(r *Result) *Result {
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
		if s.direction == Asc {
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

func Sort(direction SortDirection) ResultOption {
	if direction == Asc {
		return sortOptAsc
	}
	return sortOptDesc
}

// Count

type countOpt struct {
	predicate func(DiceResult) bool
}

func (s *countOpt) Apply(r *Result) *Result {
	if r == nil {
		return nil
	}
	if len(r.Dices) == 0 {
		return r
	}
	result := 0
	for _, c := range r.Dices {
		if s.predicate(c) {
			result++
		}
	}
	return &Result{
		Value:      float64(result),
		Dices:      r.Dices,
		Expression: r.Expression,
	}
}

func Count(predicate func(result DiceResult) bool) ResultOption {
	return &countOpt{predicate}
}

// Keep

type keepOpt struct {
	count int
	lower bool
}

func (s *keepOpt) Apply(r *Result) *Result {
	if r == nil {
		return nil
	}
	if len(r.Dices) == 0 {
		return r
	}
	cd := make([]DiceResult, len(r.Dices))
	for idx, c := range r.Dices {
		cd[idx] = c
	}
	sort.Slice(cd, func(i, j int) bool { return cd[i].Value < cd[j].Value })
	asFloat := funk.Map(cd, func(c DiceResult) float64 { return float64(c.Value) }).([]float64)
	var res float64
	if s.count >= len(asFloat) {
		res = funk.Sum(asFloat)
	} else if s.lower {
		res = funk.Sum(asFloat[:s.count])
	} else {
		res = funk.Sum(asFloat[len(asFloat)-s.count:])
	}
	return &Result{
		Value:      res,
		Dices:      r.Dices,
		Expression: r.Expression,
	}
}

func Keep(count int, lower bool) ResultOption {
	return &keepOpt{count, lower}
}
