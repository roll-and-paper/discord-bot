package parser

import (
	"fmt"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/dices"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/dices/roller"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func toIterator(value string) *iterator {
	if len(value) == 0 {
		return newIterator([]Token{})
	}
	iter, _ := NewLexer(value).ScanAll()
	return newIterator(iter)
}

func Test_Parse(t *testing.T) {
	tests := []struct {
		args  string
		want  roller.Roller
		want1 error
	}{
		{"", nil, fmt.Errorf("expression not valid")},
		{"!", nil, fmt.Errorf("illegal token ! ")},
		{"3d6s3", nil, fmt.Errorf("expression not finished")},
		{"3d", nil, fmt.Errorf("bad dice format")},
		{"3d6", roller.FromDice(dices.D6, 3, nil), nil},
		{"3d6s", roller.FromDice(dices.D6, 3, []roller.Option{roller.Sort(roller.Desc)}), nil},
	}
	for idx, tt := range tests {
		t.Run(strconv.FormatInt(int64(idx), 10), func(t *testing.T) {
			got, got1 := Parse(tt.args)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func Test_options(t *testing.T) {
	tests := []struct {
		args        *iterator
		want        []roller.Option
		want1       bool
		want2       int
		idx         int
		shouldPanic bool
	}{
		{toIterator(""), []roller.Option{}, false, 0, 0, false},
		{toIterator("55"), []roller.Option{}, false, 0, 0, false},
		{toIterator("k5"), []roller.Option{roller.Keep(5, false)}, true, 2, 2, false},
		{toIterator("kl5"), []roller.Option{roller.Keep(5, true)}, true, 2, 2, false},
		{toIterator("s"), []roller.Option{roller.Sort(roller.Desc)}, true, 1, 1, false},
		{toIterator("sl"), []roller.Option{roller.Sort(roller.Asc)}, true, 1, 1, false},
	}
	for idx, tt := range tests {
		t.Run(strconv.FormatInt(int64(idx), 10), func(t *testing.T) {
			if tt.shouldPanic {
				assert.Panics(t, func() { options(tt.args) })
			} else {
				got, got1, got2 := options(tt.args)
				assert.Equal(t, tt.want, got)
				assert.Equal(t, tt.want1, got1)
				assert.Equal(t, tt.want2, got2)
				assert.Equal(t, tt.idx, tt.args.idx)
			}
		})
	}
}

func Test_arithmetic(t *testing.T) {
	tests := []struct {
		args        *iterator
		want        roller.Roller
		want1       bool
		want2       int
		idx         int
		shouldPanic bool
	}{
		{toIterator(""), nil, false, 0, 0, false},
		{toIterator("5d6+"), nil, false, 0, 0, true},
		{toIterator("5d6-"), nil, false, 0, 0, true},
		{toIterator("5d6*"), nil, false, 0, 0, true},
		{toIterator("5d6/"), nil, false, 0, 0, true},
		{toIterator("5d6|"), nil, false, 0, 0, true},
		{toIterator("5d6**"), nil, false, 0, 0, true},
		{toIterator("5d6%"), nil, false, 0, 0, true},
		{toIterator("5d6+5d6"), roller.FromDice(dices.NewFromMax(6), 5, nil).Add(roller.FromDice(dices.NewFromMax(6), 5, nil)), true, 7, 7, false},
		{toIterator("5d6-5d6"), roller.FromDice(dices.NewFromMax(6), 5, nil).Minus(roller.FromDice(dices.NewFromMax(6), 5, nil)), true, 7, 7, false},
		{toIterator("5d6*5d6"), roller.FromDice(dices.NewFromMax(6), 5, nil).Mul(roller.FromDice(dices.NewFromMax(6), 5, nil)), true, 7, 7, false},
		{toIterator("5d6/5d6"), roller.FromDice(dices.NewFromMax(6), 5, nil).Div(roller.FromDice(dices.NewFromMax(6), 5, nil)), true, 7, 7, false},
		{toIterator("5d6|5d6"), roller.FromDice(dices.NewFromMax(6), 5, nil).IntDiv(roller.FromDice(dices.NewFromMax(6), 5, nil)), true, 7, 7, false},
		{toIterator("5d6**5d6"), roller.FromDice(dices.NewFromMax(6), 5, nil).Pow(roller.FromDice(dices.NewFromMax(6), 5, nil)), true, 7, 7, false},
		{toIterator("5d6%5d6"), roller.FromDice(dices.NewFromMax(6), 5, nil).Mod(roller.FromDice(dices.NewFromMax(6), 5, nil)), true, 7, 7, false},
		{toIterator("5d6"), roller.FromDice(dices.NewFromMax(6), 5, nil), true, 3, 3, false},
		//{toIterator("5d[1..3]"), roller.FromDice(dices.NewFromFaces([]int{1, 2, 3}), 5, nil), true, 7, 7, false},
		//{toIterator("5d[1,3,5]"), roller.FromDice(dices.NewFromFaces([]int{1, 3, 5}), 5, nil), true, 9, 9, false},
		//{toIterator("5"), roller.FromValue(5, nil), true, 1, 1, false},
		//{toIterator("[5,6]"), roller.FromValues([]float64{5, 6}, []roller.DiceResult{{Value: 5}, {Value: 6}}, nil), true, 5, 5, false},
		//{newIterator([]Token{{Dice, "d"}, {BraceOut, "]"}}), nil, false, 0, 0},
		//{newIterator([]Token{{Dice, "d"}, {Number, "6"}}), dices.NewFromMax(6), true, 2, 2},
		//{newIterator([]Token{{Dice, "d"}, {BraceIn, "["}, {Number, "1"}, {DoubleDot, ".."}, {Number, "4"}, {BraceOut, "]"}}), dices.NewFromFaces([]int{1, 2, 3, 4}), true, 6, 6},
		//{newIterator([]Token{{Dice, "d"}, {BraceIn, "["}, {Number, "2"}, {Comma, ","}, {Number, "4"}, {BraceOut, "]"}}), dices.NewFromFaces([]int{2, 4}), true, 6, 6},
	}
	for idx, tt := range tests {
		t.Run(strconv.FormatInt(int64(idx), 10), func(t *testing.T) {
			if tt.shouldPanic {
				assert.Panics(t, func() { arithmetic(tt.args) })
			} else {
				got, got1, got2 := arithmetic(tt.args)
				assert.Equal(t, tt.want, got)
				assert.Equal(t, tt.want1, got1)
				assert.Equal(t, tt.want2, got2)
				assert.Equal(t, tt.idx, tt.args.idx)
			}
		})
	}
}

func Test_roll(t *testing.T) {
	tests := []struct {
		args        *iterator
		want        roller.Roller
		want1       bool
		want2       int
		idx         int
		shouldPanic bool
	}{
		{toIterator(""), nil, false, 0, 0, false},
		{toIterator("5d"), nil, false, 0, 0, true},
		{toIterator("(((5d6))"), nil, false, 9, 9, true},
		{toIterator("(((5d6)"), nil, false, 9, 9, true},
		{toIterator("(((5d6"), nil, false, 9, 9, true},
		{toIterator("(5d6"), nil, false, 9, 9, true},
		{toIterator("(d6)"), nil, false, 9, 9, true},
		{toIterator("5d6"), roller.FromDice(dices.NewFromMax(6), 5, nil), true, 3, 3, false},
		{toIterator("(5d6)"), roller.FromDice(dices.NewFromMax(6), 5, nil), true, 5, 5, false},
		{toIterator("(((5d6)))"), roller.FromDice(dices.NewFromMax(6), 5, nil), true, 9, 9, false},
		{toIterator("(5d6)d6"), roller.FromRoller(dices.NewFromMax(6), roller.FromDice(dices.NewFromMax(6), 5, nil), nil), true, 7, 7, false},
		{toIterator("(5d6+3)d6"), roller.FromRoller(dices.NewFromMax(6), roller.FromDice(dices.NewFromMax(6), 5, nil).Add(roller.FromValue(3, nil)), nil), true, 9, 9, false},
		{toIterator("5d[1..3]"), roller.FromDice(dices.NewFromFaces([]int{1, 2, 3}), 5, nil), true, 7, 7, false},
		{toIterator("5d[1,3,5]"), roller.FromDice(dices.NewFromFaces([]int{1, 3, 5}), 5, nil), true, 9, 9, false},
		{toIterator("5"), roller.FromValue(5, nil), true, 1, 1, false},
		{toIterator("[5,6]"), roller.FromValues([]float64{5, 6}, []roller.DiceResult{{Value: 5}, {Value: 6}}, nil), true, 5, 5, false},
		//{newIterator([]Token{{Dice, "d"}, {BraceOut, "]"}}), nil, false, 0, 0},
		//{newIterator([]Token{{Dice, "d"}, {Number, "6"}}), dices.NewFromMax(6), true, 2, 2},
		//{newIterator([]Token{{Dice, "d"}, {BraceIn, "["}, {Number, "1"}, {DoubleDot, ".."}, {Number, "4"}, {BraceOut, "]"}}), dices.NewFromFaces([]int{1, 2, 3, 4}), true, 6, 6},
		//{newIterator([]Token{{Dice, "d"}, {BraceIn, "["}, {Number, "2"}, {Comma, ","}, {Number, "4"}, {BraceOut, "]"}}), dices.NewFromFaces([]int{2, 4}), true, 6, 6},
	}
	for idx, tt := range tests {
		t.Run(strconv.FormatInt(int64(idx), 10), func(t *testing.T) {
			if tt.shouldPanic {
				assert.Panics(t, func() { roll(tt.args) })
			} else {
				got, got1, got2 := roll(tt.args)
				assert.Equal(t, tt.want, got)
				assert.Equal(t, tt.want1, got1)
				assert.Equal(t, tt.want2, got2)
				assert.Equal(t, tt.idx, tt.args.idx)
			}
		})
	}
}

func Test_dice(t *testing.T) {
	tests := []struct {
		args        *iterator
		want        dices.Dice
		want1       bool
		want2       int
		idx         int
		shouldPanic bool
	}{
		{toIterator(""), nil, false, 0, 0, false},
		{toIterator("d]"), nil, false, 0, 0, true},
		{toIterator("d"), nil, false, 0, 0, true},
		{toIterator("d6"), dices.NewFromMax(6), true, 2, 2, false},
		{toIterator("d[1..4]"), dices.NewFromFaces([]int{1, 2, 3, 4}), true, 6, 6, false},
		{toIterator("d[2,4]"), dices.NewFromFaces([]int{2, 4}), true, 6, 6, false},
	}
	for idx, tt := range tests {
		t.Run(strconv.FormatInt(int64(idx), 10), func(t *testing.T) {
			if tt.shouldPanic {
				assert.Panics(t, func() { dice(tt.args) })
			} else {
				got, got1, got2 := dice(tt.args)
				assert.Equal(t, tt.want, got)
				assert.Equal(t, tt.want1, got1)
				assert.Equal(t, tt.want2, got2)
				assert.Equal(t, tt.idx, tt.args.idx)
			}
		})
	}
}

func Test_numbers(t *testing.T) {
	tests := []struct {
		args  *iterator
		want  []int
		want1 bool
		want2 int
		idx   int
	}{
		{toIterator(""), []int{}, false, 0, 0},
		{toIterator("[]"), []int{}, true, 2, 2},
		{toIterator("[123]"), []int{123}, true, 3, 3},
		{toIterator("[-123]"), []int{-123}, true, 4, 4},
		{newIterator([]Token{{BraceIn, "["}, {Minus, "-"}, {Number, "123"}, {Number, "123"}, {BraceOut, "]"}}), []int{}, false, 0, 0},
		{toIterator("[-123,123]"), []int{-123, 123}, true, 6, 6},
	}
	for idx, tt := range tests {
		t.Run(strconv.FormatInt(int64(idx), 10), func(t *testing.T) {
			got, got1, got2 := numbers(tt.args)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, tt.want2, got2)
			assert.Equal(t, tt.idx, tt.args.idx)
		})
	}
}

func Test_number(t *testing.T) {
	tests := []struct {
		args  *iterator
		want  int
		want1 bool
		want2 int
		idx   int
	}{
		{toIterator(""), 0, false, 0, 0},
		{newIterator([]Token{{Minus, "-"}, {Number, "12a3"}}), 0, false, 0, 0},
		{newIterator([]Token{{Number, "12a3"}}), 0, false, 0, 0},
		{toIterator("123"), 123, true, 1, 1},
		{toIterator("-123"), -123, true, 2, 2},
		{toIterator("-d"), 0, false, 0, 0},
		{toIterator("d"), 0, false, 0, 0},
	}
	for idx, tt := range tests {
		t.Run(strconv.FormatInt(int64(idx), 10), func(t *testing.T) {
			got, got1, got2 := number(tt.args)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, tt.want2, got2)
			assert.Equal(t, tt.idx, tt.args.idx)
		})
	}
}

func Test_rrange(t *testing.T) {
	tests := []struct {
		args  *iterator
		want  []int
		want1 bool
		want2 int
		idx   int
	}{
		{toIterator(""), []int{}, false, 0, 0},
		{toIterator("[]"), []int{}, false, 0, 0},
		{toIterator("["), []int{}, false, 0, 0},
		{toIterator("[1"), []int{}, false, 0, 0},
		{toIterator("[-1..1"), []int{}, false, 0, 0},
		{toIterator("[1]"), []int{}, false, 0, 0},
		{toIterator("[-1]"), []int{}, false, 0, 0},
		{toIterator("[-1,1]"), []int{}, false, 0, 0},
		{toIterator("[-1..1]"), []int{-1, 0, 1}, true, 6, 6},
		{toIterator("[-1..-1]"), []int{-1}, true, 7, 7},
		{toIterator("[1..-1]"), []int{1, 0, -1}, true, 6, 6},
	}
	for idx, tt := range tests {
		t.Run(strconv.FormatInt(int64(idx), 10), func(t *testing.T) {
			got, got1, got2 := rrange(tt.args)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, tt.want2, got2)
			assert.Equal(t, tt.idx, tt.args.idx)
		})
	}
}
