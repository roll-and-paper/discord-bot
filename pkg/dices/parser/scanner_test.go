package parser

import (
	"reflect"
	"strconv"
	"testing"
)

func TestScanner_ScanAll(t *testing.T) {
	tests := []struct {
		input string
		want  []Token
	}{
		{"   ", []Token{{WS, "   "}}},
		{"123", []Token{{Number, "123"}}},
		{"123.321", []Token{{Number, "123"}, {Dot, "."}, {Number, "321"}}},
		{"123.", []Token{{Number, "123"}, {Dot, "."}}},
		{"[]()", []Token{{BraceIn, "["}, {BraceOut, "]"}, {BracketIn, "("}, {BracketOut, ")"}}},
		{"..", []Token{{DoubleDot, ".."}}},
		{"...", []Token{{DoubleDot, ".."}, {Dot, "."}}},
		{"....", []Token{{DoubleDot, ".."}, {DoubleDot, ".."}}},
		{"=>=<=><!=", []Token{{Equal, "="}, {GreaterOrEqual, ">="}, {LesserOrEqual, "<="}, {GreaterThen, ">"}, {LesserThan, "<"}, {Different, "!="}}},
		{"&^", []Token{{And, "&"}, {XOr, "^"}}},
		{"%", []Token{{Modulo, "%"}}},
		{"+-x/÷***", []Token{{Plus, "+"}, {Minus, "-"}, {Multiplication, "x"}, {Divide, "/"}, {Divide, "÷"}, {Pow, "**"}, {Multiplication, "*"}}},
		{"kKse", []Token{{Keep, "k"}, {KeepAndExplode, "K"}, {Sort, "s"}, {Explode, "e"}}},
		{"crRa", []Token{{Count, "c"}, {Reroll, "r"}, {RerollUntil, "R"}, {RerollAndAdd, "a"}}},
		{"mipf", []Token{{Merge, "m"}, {IfOperator, "i"}, {Painter, "p"}, {Filter, "f"}}},
		{"yutgbo", []Token{{Split, "y"}, {Unique, "u"}, {AllSameExplode, "t"}, {Group, "g"}, {Bind, "b"}, {Occurrences, "o"}}},
	}
	for idx, tt := range tests {
		t.Run(strconv.FormatInt(int64(idx), 10), func(t *testing.T) {
			s := NewScanner(tt.input)
			if got := s.ScanAll(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ScanAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
