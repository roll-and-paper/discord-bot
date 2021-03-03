package roller

import (
	"github.com/dohr-michael/roll-and-paper-bot/pkg/dices"
	"reflect"
	"strconv"
	"testing"
)

func Test_Sort(t *testing.T) {
	tests := []struct {
		name   string
		fields SortDirection
		args   *Result
		want   *Result
	}{
		{"1", Asc, &Result{12.0, []DiceResult{{4, nil}, {4, []int{3, 2}}, {3, nil}, {5, nil}}, ""}, &Result{12.0, []DiceResult{{3, nil}, {4, nil}, {4, []int{3, 2}}, {5, nil}}, ""}},
		{"2", Desc, &Result{12.0, []DiceResult{{4, nil}, {4, []int{3, 2}}, {3, nil}, {5, nil}}, ""}, &Result{12.0, []DiceResult{{5, nil}, {4, nil}, {4, []int{3, 2}}, {3, nil}}, ""}},
		{"2", Desc, &Result{12.0, nil, ""}, &Result{12, nil, ""}},
		{"2", Desc, nil, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			so := Sort(tt.fields)
			if got := so(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Count(t *testing.T) {
	dices.SetRandomSeed(42)
	tests := []struct {
		name      string
		predicate func(DiceResult) bool
		args      *Result
		want      *Result
	}{
		{"1", func(r DiceResult) bool { return r.Value >= 4 }, &Result{12.0, []DiceResult{{4, nil}, {4, []int{3, 2}}, {3, nil}, {5, nil}}, ""}, &Result{3.0, []DiceResult{{4, nil}, {4, []int{3, 2}}, {3, nil}, {5, nil}}, ""}},
	}
	for idx, tt := range tests {
		t.Run(strconv.FormatInt(int64(idx), 10), func(t *testing.T) {
			so := Count(tt.predicate)
			if got := so(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Explode(t *testing.T) {
	//var twentyFirsts = []int{6, 6, 3, 1, 2, 2, 4, 3, 3, 2, 2, 6, 4, 5, 3, 4, 3, 6, 6, 3}

	dices.SetRandomSeed(42)
	tests := []struct {
		args DiceResult
		want DiceResult
	}{
		{DiceResult{6, nil}, DiceResult{21, []int{6, 6, 6, 3}}},
		{DiceResult{2, nil}, DiceResult{2, nil}},
	}

	for idx, tt := range tests {
		t.Run(strconv.FormatInt(int64(idx), 10), func(t *testing.T) {
			so := Explode(MaxValue)
			if got := so(tt.args, dices.NewFromMax(6)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}
