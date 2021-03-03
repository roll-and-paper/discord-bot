package roller

import (
	"reflect"
	"testing"
)

func TestDiceResult_String(t1 *testing.T) {
	tests := []struct {
		name   string
		fields DiceResult
		want   string
	}{
		{"1", DiceResult{6, []int{}}, "6"},
		{"2", DiceResult{6, []int{1, 2, 3}}, "6 [1 2 3]"},
		{"3", DiceResult{666, []int{6, 6, 6}}, "666 [6 6 6]"},
		{"4", DiceResult{666, nil}, "666"},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := tt.fields
			if got := t.String(); got != tt.want {
				t1.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResult_Add(t1 *testing.T) {
	tests := []struct {
		name   string
		fields *Result
		args   *Result
		want   *Result
	}{
		{"1", NewResult(1.0, []DiceResult{{1, nil}}, "v1"), NewResult(1.0, []DiceResult{{1, nil}}, "v2"), NewResult(2.0, []DiceResult{{1, nil}, {1, nil}}, "v1")},
		{"2", NewResult(2.0, []DiceResult{{1, nil}, {1, nil}}, "v1"), NewResult(3.0, []DiceResult{{3, nil}}, "v2"), NewResult(5.0, []DiceResult{{1, nil}, {1, nil}, {3, nil}}, "v1")},
		{"3", NewResult(2.0, nil, "v1"), NewResult(3.0, []DiceResult{{3, nil}}, "v2"), NewResult(5.0, []DiceResult{{3, nil}}, "v1")},
		{"4", NewResult(2.0, nil, "v1"), NewResult(3.0, nil, "v2"), NewResult(5.0, nil, "v1")},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := tt.fields
			if got := t.Add(tt.args); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResult_Div(t1 *testing.T) {
	tests := []struct {
		name   string
		fields *Result
		args   *Result
		want   *Result
	}{
		{"1", NewResult(1.0, []DiceResult{{1, nil}}, "v1"), NewResult(1.0, []DiceResult{{1, nil}}, "v2"), NewResult(1.0, []DiceResult{{1, nil}, {1, nil}}, "v1")},
		{"2", NewResult(2.0, []DiceResult{{1, nil}, {1, nil}}, "v1"), NewResult(3.0, []DiceResult{{3, nil}}, "v2"), NewResult(2.0/3.0, []DiceResult{{1, nil}, {1, nil}, {3, nil}}, "v1")},
		{"3", NewResult(2.0, nil, "v1"), NewResult(3.0, []DiceResult{{3, nil}}, "v2"), NewResult(2.0/3.0, []DiceResult{{3, nil}}, "v1")},
		{"4", NewResult(2.0, nil, "v1"), NewResult(3.0, nil, "v2"), NewResult(2.0/3.0, nil, "v1")},
		{"5", NewResult(2.0, nil, "v1"), NewResult(2.0, nil, "v2"), NewResult(1.0, nil, "v1")},
		{"6", NewResult(4.0, nil, "v1"), NewResult(2.0, nil, "v2"), NewResult(2.0, nil, "v1")},
		{"7", NewResult(5.0, nil, "v1"), NewResult(2.0, nil, "v2"), NewResult(2.5, nil, "v1")},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := tt.fields
			if got := t.Div(tt.args); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Div() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResult_IntDiv(t1 *testing.T) {
	tests := []struct {
		name   string
		fields *Result
		args   *Result
		want   *Result
	}{
		{"1", NewResult(1.0, []DiceResult{{1, nil}}, "v1"), NewResult(1.0, []DiceResult{{1, nil}}, "v2"), NewResult(1.0, []DiceResult{{1, nil}, {1, nil}}, "v1")},
		{"2", NewResult(2.0, []DiceResult{{1, nil}, {1, nil}}, "v1"), NewResult(3.0, []DiceResult{{3, nil}}, "v2"), NewResult(0.0, []DiceResult{{1, nil}, {1, nil}, {3, nil}}, "v1")},
		{"3", NewResult(2.0, nil, "v1"), NewResult(3.0, []DiceResult{{3, nil}}, "v2"), NewResult(0.0, []DiceResult{{3, nil}}, "v1")},
		{"4", NewResult(2.0, nil, "v1"), NewResult(3.0, nil, "v2"), NewResult(0.0, nil, "v1")},
		{"5", NewResult(2.0, nil, "v1"), NewResult(2.0, nil, "v2"), NewResult(1.0, nil, "v1")},
		{"6", NewResult(4.0, nil, "v1"), NewResult(2.0, nil, "v2"), NewResult(2.0, nil, "v1")},
		{"7", NewResult(5.0, nil, "v1"), NewResult(2.0, nil, "v2"), NewResult(2.0, nil, "v1")},
		{"8", NewResult(5.0, nil, "v1"), NewResult(3.0, nil, "v2"), NewResult(1.0, nil, "v1")},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := tt.fields
			if got := t.IntDiv(tt.args); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("IntDiv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResult_Minus(t1 *testing.T) {
	tests := []struct {
		name   string
		fields *Result
		args   *Result
		want   *Result
	}{
		{"1", NewResult(1.0, []DiceResult{{1, nil}}, "v1"), NewResult(1.0, []DiceResult{{1, nil}}, "v2"), NewResult(0.0, []DiceResult{{1, nil}, {1, nil}}, "v1")},
		{"2", NewResult(2.0, []DiceResult{{1, nil}, {1, nil}}, "v1"), NewResult(3.0, []DiceResult{{3, nil}}, "v2"), NewResult(-1.0, []DiceResult{{1, nil}, {1, nil}, {3, nil}}, "v1")},
		{"3", NewResult(2.0, nil, "v1"), NewResult(3.0, []DiceResult{{3, nil}}, "v2"), NewResult(-1.0, []DiceResult{{3, nil}}, "v1")},
		{"4", NewResult(2.0, nil, "v1"), NewResult(3.0, nil, "v2"), NewResult(-1.0, nil, "v1")},
		{"5", NewResult(2.0, nil, "v1"), NewResult(2.0, nil, "v2"), NewResult(0.0, nil, "v1")},
		{"6", NewResult(4.0, nil, "v1"), NewResult(2.0, nil, "v2"), NewResult(2.0, nil, "v1")},
		{"7", NewResult(5.0, nil, "v1"), NewResult(2.0, nil, "v2"), NewResult(3.0, nil, "v1")},
		{"8", NewResult(5.0, nil, "v1"), NewResult(3.0, nil, "v2"), NewResult(2.0, nil, "v1")},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := tt.fields
			if got := t.Minus(tt.args); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Minus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResult_Mod(t1 *testing.T) {
	tests := []struct {
		name   string
		fields *Result
		args   *Result
		want   *Result
	}{
		{"1", NewResult(1.0, []DiceResult{{1, nil}}, "v1"), NewResult(1.0, []DiceResult{{1, nil}}, "v2"), NewResult(0.0, []DiceResult{{1.0, nil}, {1, nil}}, "v1")},
		{"2", NewResult(2.0, []DiceResult{{1, nil}, {1, nil}}, "v1"), NewResult(3.0, []DiceResult{{3, nil}}, "v2"), NewResult(2.0, []DiceResult{{1, nil}, {1, nil}, {3, nil}}, "v1")},
		{"3", NewResult(2.0, nil, "v1"), NewResult(3.0, []DiceResult{{3, nil}}, "v2"), NewResult(2.0, []DiceResult{{3, nil}}, "v1")},
		{"4", NewResult(2.0, nil, "v1"), NewResult(3.0, nil, "v2"), NewResult(2.0, nil, "v1")},
		{"5", NewResult(2.0, nil, "v1"), NewResult(2.0, nil, "v2"), NewResult(0.0, nil, "v1")},
		{"6", NewResult(4.0, nil, "v1"), NewResult(2.0, nil, "v2"), NewResult(0.0, nil, "v1")},
		{"7", NewResult(5.0, nil, "v1"), NewResult(2.0, nil, "v2"), NewResult(1.0, nil, "v1")},
		{"8", NewResult(5.0, nil, "v1"), NewResult(3.0, nil, "v2"), NewResult(2.0, nil, "v1")},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := tt.fields
			if got := t.Mod(tt.args); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Mod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResult_Mul(t1 *testing.T) {
	tests := []struct {
		name   string
		fields *Result
		args   *Result
		want   *Result
	}{
		{"1", NewResult(1.0, []DiceResult{{1, nil}}, "v1"), NewResult(1.0, []DiceResult{{1, nil}}, "v2"), NewResult(1.0, []DiceResult{{1.0, nil}, {1, nil}}, "v1")},
		{"2", NewResult(2.0, []DiceResult{{1, nil}, {1, nil}}, "v1"), NewResult(3.0, []DiceResult{{3, nil}}, "v2"), NewResult(6.0, []DiceResult{{1, nil}, {1, nil}, {3, nil}}, "v1")},
		{"3", NewResult(2.0, nil, "v1"), NewResult(3.0, []DiceResult{{3, nil}}, "v2"), NewResult(6.0, []DiceResult{{3, nil}}, "v1")},
		{"4", NewResult(2.0, nil, "v1"), NewResult(3.0, nil, "v2"), NewResult(6.0, nil, "v1")},
		{"5", NewResult(2.0, nil, "v1"), NewResult(2.0, nil, "v2"), NewResult(4.0, nil, "v1")},
		{"6", NewResult(4.0, nil, "v1"), NewResult(2.0, nil, "v2"), NewResult(8.0, nil, "v1")},
		{"7", NewResult(5.0, nil, "v1"), NewResult(2.0, nil, "v2"), NewResult(10.0, nil, "v1")},
		{"8", NewResult(5.0, nil, "v1"), NewResult(3.0, nil, "v2"), NewResult(15.0, nil, "v1")},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := tt.fields
			if got := t.Mul(tt.args); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Mul() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResult_Pow(t1 *testing.T) {
	tests := []struct {
		name   string
		fields *Result
		args   *Result
		want   *Result
	}{
		{"1", NewResult(1.0, []DiceResult{{1, nil}}, "v1"), NewResult(1.0, []DiceResult{{1, nil}}, "v2"), NewResult(1.0, []DiceResult{{1.0, nil}, {1, nil}}, "v1")},
		{"2", NewResult(2.0, []DiceResult{{1, nil}, {1, nil}}, "v1"), NewResult(3.0, []DiceResult{{3, nil}}, "v2"), NewResult(8.0, []DiceResult{{1, nil}, {1, nil}, {3, nil}}, "v1")},
		{"3", NewResult(2.0, nil, "v1"), NewResult(3.0, []DiceResult{{3, nil}}, "v2"), NewResult(8.0, []DiceResult{{3, nil}}, "v1")},
		{"4", NewResult(2.0, nil, "v1"), NewResult(3.0, nil, "v2"), NewResult(8.0, nil, "v1")},
		{"5", NewResult(2.0, nil, "v1"), NewResult(2.0, nil, "v2"), NewResult(4.0, nil, "v1")},
		{"6", NewResult(4.0, nil, "v1"), NewResult(2.0, nil, "v2"), NewResult(16.0, nil, "v1")},
		{"7", NewResult(5.0, nil, "v1"), NewResult(2.0, nil, "v2"), NewResult(25.0, nil, "v1")},
		{"8", NewResult(5.0, nil, "v1"), NewResult(3.0, nil, "v2"), NewResult(125.0, nil, "v1")},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := tt.fields
			if got := t.Pow(tt.args); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Pow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResult_String(t1 *testing.T) {
	tests := []struct {
		name   string
		fields *Result
		want   string
	}{
		{"1", NewResult(1.0, []DiceResult{{1, nil}}, "v1"), "# 1\nDetails:[v1 (1)]"},
		{"2", NewResult(4.0, []DiceResult{{1, nil}, {2, nil}, {1, nil}}, "v1"), "# 4\nDetails:[v1 (1 2 1)]"},
		{"3", NewResult(12.0, []DiceResult{{1, nil}, {7, []int{6, 1}}, {4, nil}}, "v1"), "# 12\nDetails:[v1 (1 7 [6 1] 4)]"},
		{"4", NewResult(666.6, nil, "v1"), "# 666.6\nDetails:[v1 ()]"},
		{"5", NewResult(1.0/3.0, nil, "v1"), "# 0.3333\nDetails:[v1 ()]"},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := tt.fields
			if got := t.String(); got != tt.want {
				t1.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
