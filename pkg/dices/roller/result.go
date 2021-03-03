package roller

import (
	"fmt"
	"github.com/thoas/go-funk"
	"math"
	"strconv"
	"strings"
)

func NewResult(value float64, dices []DiceResult, expression string) *Result {
	return &Result{
		Value:      value,
		Dices:      dices,
		Expression: expression,
	}
}

type Result struct {
	Value          float64
	Dices          []DiceResult
	Expression     string
}

func (t *Result) appendDices(other *Result) []DiceResult {
	return append(t.Dices, other.Dices...)
}

func (t *Result) Add(other *Result) *Result {
	return NewResult(t.Value+other.Value, t.appendDices(other), t.Expression)
}
func (t *Result) Minus(other *Result) *Result {
	return NewResult(t.Value-other.Value, t.appendDices(other), t.Expression)
}
func (t *Result) Mul(other *Result) *Result {
	return NewResult(t.Value*other.Value, t.appendDices(other), t.Expression)
}
func (t *Result) Div(other *Result) *Result {
	return NewResult(t.Value/other.Value, t.appendDices(other), t.Expression)
}
func (t *Result) IntDiv(other *Result) *Result {
	return NewResult(math.Floor(t.Value/other.Value), t.appendDices(other), t.Expression)
}
func (t *Result) Pow(other *Result) *Result {
	return NewResult(math.Pow(t.Value, other.Value), t.appendDices(other), t.Expression)
}
func (t *Result) Mod(other *Result) *Result {
	return NewResult(math.Mod(t.Value, other.Value), t.appendDices(other), t.Expression)
}

func (t *Result) String() string {
	v := ""
	if math.Floor(t.Value) == t.Value {
		v = strconv.FormatInt(int64(math.Floor(t.Value)), 10)
	} else {
		v = fmt.Sprintf("%.4g", t.Value)
	}
	dices := funk.Map(t.Dices, func(d DiceResult) string { return d.String() }).([]string)
	return fmt.Sprintf("# %s\nDetails:[%s (%s)]", v, t.Expression, strings.Join(dices, " "))
}

type DiceResult struct {
	Value    int
	Exploded []int
}

func (t DiceResult) String() string {
	valueStr := strconv.FormatInt(int64(t.Value), 10)
	if len(t.Exploded) == 0 {
		return valueStr
	}
	v := funk.Map(t.Exploded, func(c int) string {
		return strconv.FormatInt(int64(c), 10)
	}).([]string)
	return fmt.Sprintf("%s [%s]", valueStr, strings.Join(v, " "))
}
