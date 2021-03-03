package roller

import (
	"github.com/dohr-michael/roll-and-paper-bot/pkg/dices"
	"github.com/thoas/go-funk"
	"math"
)

type Roller interface {
	options() []Option
	Add(other Roller) Roller
	Minus(other Roller) Roller
	Mul(other Roller) Roller
	Div(other Roller) Roller
	IntDiv(other Roller) Roller
	Pow(other Roller) Roller
	Mod(other Roller) Roller
	WithOptions(options []Option) Roller
	Roll(expression string) *Result
}

func applyResultOptions(options []Option, result *Result) *Result {
	res := result
	for _, opt := range options {
		switch v := opt.(type) {
		case ResultOption:
			res = v(res)
		}
	}
	return res
}
func applyDiceResultOptions(options []Option, result DiceResult, dice dices.Dice) DiceResult {
	res := result
	for _, opt := range options {
		switch v := opt.(type) {
		case DiceResultOption:
			res = v(res, dice)
		}
	}
	return res
}

type base struct {
	underlying Roller
	opts       []Option
}

func (t *base) options() []Option          { return t.opts }
func (t *base) Add(other Roller) Roller    { return newWithOperation(t.underlying, other, add) }
func (t *base) Minus(other Roller) Roller  { return newWithOperation(t.underlying, other, minus) }
func (t *base) Mul(other Roller) Roller    { return newWithOperation(t.underlying, other, mul) }
func (t *base) Div(other Roller) Roller    { return newWithOperation(t.underlying, other, div) }
func (t *base) IntDiv(other Roller) Roller { return newWithOperation(t.underlying, other, intDiv) }
func (t *base) Pow(other Roller) Roller    { return newWithOperation(t.underlying, other, pow) }
func (t *base) Mod(other Roller) Roller    { return newWithOperation(t.underlying, other, mod) }

// Roller from static value
func FromValue(value float64, options []Option) Roller {
	r := &fromValue{Value: value}
	r.base = &base{opts: options, underlying: r}
	return r
}

type fromValue struct {
	*base
	Value float64
}

func (t *fromValue) WithOptions(options []Option) Roller { return FromValue(t.Value, options) }
func (t *fromValue) Roll(expression string) *Result {
	return applyResultOptions(t.options(), NewResult(t.Value, nil, expression))
}

// Roller from list of values with details of results
func FromValues(values []float64, results []DiceResult, options []Option) Roller {
	r := &fromValues{
		Values:      values,
		DiceResults: results,
	}
	r.base = &base{opts: options, underlying: r}
	return r
}

type fromValues struct {
	*base
	Values      []float64
	DiceResults []DiceResult
}

func (t *fromValues) WithOptions(options []Option) Roller {
	return FromValues(t.Values, t.DiceResults, options)
}
func (t *fromValues) Roll(expression string) *Result {
	return applyResultOptions(t.options(), NewResult(funk.Sum(t.Values), t.DiceResults, expression))
}

// Roller from dice and count of rolls
func FromDice(dice dices.Dice, count int, options []Option) Roller {
	r := &fromDice{
		Dice:  dice,
		Count: count,
	}
	r.base = &base{opts: options, underlying: r}
	return r
}

type fromDice struct {
	*base
	Dice  dices.Dice
	Count int
}

func (t *fromDice) WithOptions(options []Option) Roller {
	return FromDice(t.Dice, t.Count, options)
}
func (t *fromDice) Roll(expression string) *Result {
	results := make([]DiceResult, 0)
	var sum = 0.0
	for i := 0; i < t.Count; i++ {
		current := DiceResult{t.Dice.Roll(), nil}
		diceResult := applyDiceResultOptions(t.options(), current, t.Dice)
		results = append(results, diceResult)
		sum += float64(diceResult.Value)
	}
	return applyResultOptions(t.options(), NewResult(sum, results, expression))
}

// Roller from another Roller
func FromRoller(dice dices.Dice, roller Roller, options []Option) Roller {
	r := &fromRoller{
		Dice:   dice,
		Roller: roller,
	}
	r.base = &base{opts: options, underlying: r}
	return r
}

type fromRoller struct {
	*base
	Dice   dices.Dice
	Roller Roller
}

func (t *fromRoller) WithOptions(options []Option) Roller {
	return FromRoller(t.Dice, t.Roller, options)
}

func (t *fromRoller) Roll(expression string) *Result {
	count := int(math.Floor(t.Roller.Roll(expression).Value))
	return FromDice(t.Dice, count, t.options()).Roll(expression)
}

// Roller from result reference of another Roller
func FromResultReference(result *Result, options []Option) Roller {
	r := &fromResultReference{Result: result}
	r.base = &base{opts: options, underlying: r}
	return r
}

type fromResultReference struct {
	*base
	Result *Result
}

func (t *fromResultReference) WithOptions(options []Option) Roller {
	return FromResultReference(t.Result, options)
}

func (t *fromResultReference) Roll(expression string) *Result {
	res := applyResultOptions(t.options(), t.Result)
	// Remove dices reference, only value is important here.
	res.Dices = nil
	return res
}

// Roller with operation
func newWithOperation(left Roller, right Roller, ops operation) *withOperation {
	return &withOperation{
		base: &base{
			opts: append(left.options(), right.options()...),
		},
		left:  left,
		right: right,
		ops:   ops,
	}
}

type withOperation struct {
	*base
	left  Roller
	right Roller
	ops   operation
}

func (t *withOperation) WithOptions(options []Option) Roller {
	return newWithOperation(t.left.WithOptions(options), t.right.WithOptions(nil), t.ops)
}

func (t *withOperation) Roll(expression string) *Result {
	return applyResultOptions(t.options(), t.ops(t.left, t.right)(expression))
}
