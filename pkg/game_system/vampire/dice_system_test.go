package vampire

import (
	"github.com/dohr-michael/roll-and-paper-bot/pkg/dices"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/dices/roller"
	"github.com/magiconair/properties/assert"
	"strconv"
	"testing"
)

func TestRoll(t *testing.T) {
	// var twentyFirsts = []int{7, 9, 10, 10, 7, 9, 1, 5, 9, 4, 1, 8, 9, 5, 1, 6, 8, 3, 5, 4}
	dices.SetRandomSeed(56)
	tests := []struct {
		count      int
		difficulty int
		nbSpecs    int
		want       *roller.Result
	}{
		{7, 8, 1, roller.NewResult(5,
			[]roller.DiceResult{
				{7, nil},
				{9, nil},
				{10, nil},
				{10, nil},
				{7, nil},
				{9, nil},
				{1, nil},
			}, "7d : count(>= 8) + count(= 10) x 1 - count(= 1)"),
		},
		{1, 0, 1, roller.NewResult(1,
			[]roller.DiceResult{
				{5, nil},
			}, "1d : count(= 10) x 1 - count(= 1)"),
		},
		{1, 5, 0, roller.NewResult(1,
			[]roller.DiceResult{
				{9, nil},
			}, "1d : count(>= 5) - count(= 1)"),
		},
		{1, 0, 0, roller.NewResult(1,
			[]roller.DiceResult{
				{4, nil},
			}, "1d : - count(= 1)"),
		},
	}
	for idx, tt := range tests {
		t.Run(strconv.FormatInt(int64(idx), 10), func(t *testing.T) {
			assert.Equal(t, Roll(tt.count, tt.difficulty, tt.nbSpecs), tt.want)
			//if got := Roll(tt.count, tt.difficulty, tt.nbSpecs); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("Roll() = %v, want %v", got, tt.want)
			//}
		})
	}
}
