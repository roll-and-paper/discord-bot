package roller

import (
	"github.com/dohr-michael/roll-and-paper-bot/pkg/dices"
	"strconv"
	"testing"
)

func TestHasValue(t *testing.T) {
	tests := []struct {
		args   int
		result DiceResult
		want   bool
	}{
		{6, DiceResult{6, nil}, true},
		{6, DiceResult{6, []int{6}}, false},
		{5, DiceResult{5, nil}, true},
		{5, DiceResult{5, []int{5}}, false},
		{6, DiceResult{5, nil}, false},
		{6, DiceResult{5, []int{6}}, false},
	}
	for idx, tt := range tests {
		t.Run(strconv.FormatInt(int64(idx), 10), func(t *testing.T) {
			fn := HasValue(tt.args)
			if got := fn(dices.NewFromMax(6), tt.result); got != tt.want {
				t.Errorf("HasValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxValue(t *testing.T) {
	tests := []struct {
		result DiceResult
		want   bool
	}{
		{DiceResult{6, nil}, true},
		{DiceResult{6, []int{6}}, false},
		{DiceResult{5, nil}, false},
		{DiceResult{5, []int{5}}, false},
		{DiceResult{5, nil}, false},
		{DiceResult{5, []int{6}}, false},
	}
	for idx, tt := range tests {
		t.Run(strconv.FormatInt(int64(idx), 10), func(t *testing.T) {
			if got := MaxValue(dices.NewFromMax(6), tt.result); got != tt.want {
				t.Errorf("MaxValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
