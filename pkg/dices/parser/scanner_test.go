package parser

import (
	"reflect"
	"strconv"
	"testing"
)

func TestScanner_ScanAll(t *testing.T) {
	tests := []struct {
		input   string
		want    []Token
		onError bool
	}{
		{"   ", []Token{{WS, "   "}}, false},
		{"123", []Token{{Number, "123"}}, false},
		{"123.321", []Token{{Number, "123"}, {Dot, "."}, {Number, "321"}}, false},
		{"123.", []Token{{Number, "123"}, {Dot, "."}}, false},
		{"[]()", []Token{{BraceIn, "["}, {BraceOut, "]"}, {BracketIn, "("}, {BracketOut, ")"}}, false},
		{"..", []Token{{DoubleDot, ".."}}, false},
		{"...", []Token{{DoubleDot, ".."}, {Dot, "."}}, false},
		{"....", []Token{{DoubleDot, ".."}, {DoubleDot, ".."}}, false},
		{"=>=<=><!=", []Token{{Equal, "="}, {GreaterOrEqual, ">="}, {LesserOrEqual, "<="}, {GreaterThen, ">"}, {LesserThan, "<"}, {Different, "!="}}, false},
		{"&^", []Token{{And, "&"}, {XOr, "^"}}, false},
		{"%", []Token{{Modulo, "%"}}, false},
		{"+-x/÷***", []Token{{Plus, "+"}, {Minus, "-"}, {Multiplication, "x"}, {Divide, "/"}, {Divide, "÷"}, {Pow, "**"}, {Multiplication, "*"}}, false},
		{"kKse", []Token{{Keep, "k"}, {KeepAndExplode, "K"}, {Sort, "s"}, {Explode, "e"}}, false},
		{"crRa", []Token{{Count, "c"}, {Reroll, "r"}, {RerollUntil, "R"}, {RerollAndAdd, "a"}}, false},
		{"mipf", []Token{{Merge, "m"}, {IfOperator, "i"}, {Painter, "p"}, {Filter, "f"}}, false},
		{"yutgbo", []Token{{Split, "y"}, {Unique, "u"}, {AllSameExplode, "t"}, {Group, "g"}, {Bind, "b"}, {Occurrences, "o"}}, false},
	}
	for idx, tt := range tests {
		t.Run(strconv.FormatInt(int64(idx), 10), func(t *testing.T) {
			s := NewLexer(tt.input)
			if got, err := s.ScanAll(); tt.onError && err != nil || !tt.onError && err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ScanAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
