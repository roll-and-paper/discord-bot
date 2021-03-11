package roller

import (
	"github.com/dohr-michael/roll-and-paper-bot/pkg/dices"
	"github.com/thoas/go-funk"
	"reflect"
	"strconv"
	"testing"
)

//var twentyFirsts = []int{6, 6, 3, 1, 2, 2, 4, 3, 3, 2, 2, 6, 4, 5, 3, 4, 3, 6, 6, 3}
var fiveFirsts = []int{6, 6, 3, 1, 2}
var fiveFirstsSum = funk.SumInt(fiveFirsts)
var fiveFirstsRollResults = funk.Map(fiveFirsts, func(v int) DiceResult { return DiceResult{v, nil} }).([]DiceResult)
var fiveFirstsSorted = []int{1, 2, 3, 6, 6}
var fiveFirstsSortedRollResults = funk.Map(fiveFirstsSorted, func(v int) DiceResult { return DiceResult{v, nil} }).([]DiceResult)

func TestRoller_Roll(t *testing.T) {

	dices.SetRandomSeed(42)
	tests := []struct {
		args Roller
		want *Result
	}{
		{FromValue(5.0, nil), &Result{5.0, nil, ""}},
		{FromValue(5.0, nil).Add(FromValue(6.0, nil)), &Result{11.0, nil, ""}},
		{FromValue(5.0, nil).Minus(FromValue(4.0, nil)), &Result{1.0, nil, ""}},
		{FromValue(5.0, nil).Mul(FromValue(2.0, nil)), &Result{10.0, nil, ""}},
		{FromValue(5.0, nil).Div(FromValue(2.0, nil)), &Result{2.5, nil, ""}},
		{FromValue(5.0, nil).IntDiv(FromValue(2.0, nil)), &Result{2.0, nil, ""}},
		{FromValue(5.0, nil).Pow(FromValue(2.0, nil)), &Result{25.0, nil, ""}},
		{FromValue(5.0, nil).Mod(FromValue(2.0, nil)), &Result{1.0, nil, ""}},
		{FromValues([]float64{1, 2, 3}, []DiceResult{{1, nil}, {2, nil}, {3, nil}}, nil), &Result{6, []DiceResult{{1, nil}, {2, nil}, {3, nil}}, ""}},
		{FromDice(dices.NewFromMax(6), 5, nil), &Result{float64(fiveFirstsSum), fiveFirstsRollResults, ""}},
		{FromRoller(dices.NewFromMax(6), FromValue(2.0, nil), nil), &Result{6.0, []DiceResult{{2, nil}, {4, nil}}, ""}},
	}
	for idx, tt := range tests {
		t.Run(strconv.FormatInt(int64(idx), 10), func(t *testing.T) {
			if got := tt.args.Roll(""); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Roller_Roll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoller_WithOptions(t *testing.T) {
	dices.SetRandomSeed(42)
	opts := []Option{Sort(Desc)}
	tests := []struct {
		args Roller
		opts []Option
		want *Result
	}{
		{FromValue(5.0, nil), opts, &Result{5.0, nil, ""}},
		{FromValues([]float64{1, 2, 3}, []DiceResult{{1, nil}, {2, nil}, {3, nil}}, nil), opts, &Result{6, []DiceResult{{3, nil}, {2, nil}, {1, nil}}, ""}},
		{FromValues([]float64{1}, []DiceResult{{1, nil}}, nil).Add(FromValues([]float64{3, 2}, []DiceResult{{3, nil}, {2, nil}}, nil)), opts, &Result{6, []DiceResult{{3, nil}, {2, nil}, {1, nil}}, ""}},
		{FromDice(dices.NewFromMax(6), 5, nil), opts, &Result{float64(fiveFirstsSum), funk.Reverse(fiveFirstsSortedRollResults).([]DiceResult), ""}},
		{FromRoller(dices.NewFromMax(6), FromValue(2.0, nil), nil), opts, &Result{6.0, []DiceResult{{4, nil}, {2, nil}}, ""}},
		{FromDice(dices.NewFromMax(6), 1, nil), []Option{Explode(HasValue(3))}, &Result{8.0, []DiceResult{{8, []int{3, 3, 2}}}, ""}},
		{FromResultReference(&Result{12, []DiceResult{{3, nil}, {5, nil}}, ""}, nil), opts, &Result{12.0, nil, ""}},
	}
	for idx, tt := range tests {
		t.Run(strconv.FormatInt(int64(idx), 10), func(t *testing.T) {
			if got := tt.args.WithOptions(tt.opts).Roll(""); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Roller_WithOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
