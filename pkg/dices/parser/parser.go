package parser

import (
	"fmt"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/dices"
	"github.com/dohr-michael/roll-and-paper-bot/pkg/dices/roller"
	"github.com/thoas/go-funk"
	"strconv"
)

func newIterator(tokens []Token) *iterator {
	return &iterator{
		tokens: tokens,
		length: len(tokens),
		idx:    0,
	}
}

func Parse(value string) (roller.Roller, error) {
	all, err := NewLexer(value).ScanAll()
	if err != nil {
		return nil, err
	}
	return r(newIterator(all))
}

func r(tokens *iterator) (res roller.Roller, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch c := r.(type) {
			case error:
				err = c
			}
		}
	}()
	tmp, ok, _ := arithmetic(tokens)
	if !ok {
		err = fmt.Errorf("expression not valid")
		return
	}
	opts, ok, _ := options(tokens)
	if ok {
		tmp = tmp.WithOptions(opts)
	}
	if tokens.current().Type != EOF {
		err = fmt.Errorf("expression not finished")
		return
	}
	res = tmp
	return
}

func options(tokens *iterator) ([]roller.Option, bool, int) {
	res := make([]roller.Option, 0)
	count := 0
	for {
		opt, ok, c := option(tokens)
		if !ok {
			break
		}
		res = append(res, opt...)
		count += c
	}

	return res, len(res) > 0, count
}

func option(tokens *iterator) ([]roller.Option, bool, int) {
	tok, c := tokens.currentAndMove()
	switch tok.Type {
	case Keep:
		if n, ok, c1 := number(tokens); ok {
			return []roller.Option{roller.Keep(n, false)}, true, c + c1
		}
	case KeepLower:
		if n, ok, c1 := number(tokens); ok {
			return []roller.Option{roller.Keep(n, true)}, true, c + c1
		}
	case Sort:
		return []roller.Option{roller.Sort(roller.Desc)}, true, c
	case SortAsc:
		return []roller.Option{roller.Sort(roller.Asc)}, true, c
		//case KeepLower:
		//case Explode:
		//case KeepAndExplode:

	}
	tokens.unmove(c)
	return []roller.Option{}, false, 0
}

func arithmetic(tokens *iterator) (roller.Roller, bool, int) {
	if r, ok1, c1 := roll(tokens); ok1 {
		takeNext := func() (roller.Roller, int) {
			v, ok, c := roll(tokens)
			if !ok {
				panic(fmt.Errorf("invalid roller next to arithmetic operator"))
			}
			return v, c
		}
		opp, c2 := tokens.currentAndMove()
		switch opp.Type {
		case Plus:
			n, c3 := takeNext()
			return r.Add(n), true, c1 + c2 + c3
		case Minus:
			n, c3 := takeNext()
			return r.Minus(n), true, c1 + c2 + c3
		case Pow:
			n, c3 := takeNext()
			return r.Pow(n), true, c1 + c2 + c3
		case Multiplication:
			n, c3 := takeNext()
			return r.Mul(n), true, c1 + c2 + c3
		case Divide:
			n, c3 := takeNext()
			return r.Div(n), true, c1 + c2 + c3
		case Pipe:
			n, c3 := takeNext()
			return r.IntDiv(n), true, c1 + c2 + c3
		case Modulo:
			n, c3 := takeNext()
			return r.Mod(n), true, c1 + c2 + c3
		default:
			tokens.unmove(c2)
		}
		return r, true, c1
	}
	return nil, false, 0
}

func roll(tokens *iterator) (roller.Roller, bool, int) {
	if tokens.current().Type == BracketIn {
		_, c1 := tokens.currentAndMove()
		if n, ok, c2 := arithmetic(tokens); ok {
			if tokens.current().Type == BracketOut {
				_, c3 := tokens.currentAndMove()
				d, ok, c4 := dice(tokens)
				if ok {
					return roller.FromRoller(d, n, nil), true, c1 + c2 + c3 + c4
				}
				return n, true, c1 + c2 + c3
			} else {
				panic(fmt.Errorf("no `)` found after `(`"))
			}

		} else {
			panic(fmt.Errorf("invalid tokens between bracket"))
		}
	}
	if n, ok1, c1 := number(tokens); ok1 {
		d, ok2, c2 := dice(tokens)
		if ok2 {
			return roller.FromDice(d, n, nil), true, c1 + c2
		}
		return roller.FromValue(float64(n), nil), true, c1
	}
	if n, ok, c := numbers(tokens); ok {
		return roller.FromValues(
			funk.Map(n, func(v int) float64 { return float64(v) }).([]float64),
			funk.Map(n, func(v int) roller.DiceResult { return roller.DiceResult{Value: v} }).([]roller.DiceResult),
			nil,
		), true, c
	}
	return nil, false, 0
}

func dice(tokens *iterator) (dices.Dice, bool, int) {
	if tokens.current().Type == Dice {
		_, c := tokens.currentAndMove()
		n, ok1, c1 := number(tokens)
		if ok1 {
			return dices.NewFromMax(n), true, c + c1
		}
		r1, ok2, c2 := rrange(tokens)
		if ok2 {
			return dices.NewFromFaces(r1), true, c + c1 + c2
		}
		r2, ok3, c3 := numbers(tokens)
		if ok3 {
			return dices.NewFromFaces(r2), true, c + c1 + c2 + c3
		}
		tokens.unmove(c + c1 + c2 + c3)
		panic(fmt.Errorf("bad dice format"))
	}
	return nil, false, 0
}

func number(tokens *iterator) (int, bool, int) {
	if tokens.current().Type == Minus && tokens.get(tokens.idx+1).Type == Number {
		_, count1 := tokens.currentAndMove()
		v, count2 := tokens.currentAndMove()
		result, err := strconv.Atoi(fmt.Sprintf("-%v", v.Value))
		if err != nil {
			tokens.unmove(count1 + count2)
			return 0, false, 0
		}
		return result, true, count1 + count2
	} else if tokens.current().Type == Number {
		v, count := tokens.currentAndMove()
		result, err := strconv.Atoi(fmt.Sprintf("%v", v.Value))
		if err != nil {
			tokens.unmove(count)
			return 0, false, 0
		}
		return result, true, count
	}
	return 0, false, 0
}

func numbers(tokens *iterator) ([]int, bool, int) {
	count := 0
	result := make([]int, 0)
	if tokens.current().Type == BraceIn {
		_, c := tokens.currentAndMove()
		count += c
		for {
			v, ok, cc := number(tokens)
			count += cc
			if !ok {
				break
			}
			result = append(result, v)
			if tokens.current().Type == Comma {
				_, ccc := tokens.currentAndMove()
				count += ccc
			} else {
				break
			}
		}
		if tokens.current().Type == BraceOut {
			_, ccc := tokens.currentAndMove()
			count += ccc
			return result, true, count
		}
		tokens.unmove(count)
	}
	return []int{}, false, 0
}

func rrange(tokens *iterator) ([]int, bool, int) {
	if tokens.current().Type == BraceIn {
		_, c1 := tokens.currentAndMove()
		min, ok, c2 := number(tokens)
		if !ok {
			tokens.unmove(c1 + c2)
			return []int{}, false, 0
		}
		if tokens.current().Type != DoubleDot {
			tokens.unmove(c1 + c2)
			return []int{}, false, 0
		}
		_, c3 := tokens.currentAndMove()
		max, ok, c4 := number(tokens)
		if !ok || tokens.current().Type != BraceOut {
			tokens.unmove(c1 + c2 + c3 + c4)
			return []int{}, false, 0
		}
		_, c5 := tokens.currentAndMove()
		result := make([]int, 0)
		if max > min {
			for i := min; i <= max; i++ {
				result = append(result, i)
			}
		} else if min > max {
			for i := min; i >= max; i-- {
				result = append(result, i)
			}
		} else {
			result = append(result, min)
		}
		return result, true, c1 + c2 + c3 + c4 + c5
	}
	return []int{}, false, 0
}

type iterator struct {
	tokens []Token
	length int
	idx    int
}

func (t *iterator) currentAndMove() (tok Token, count int) {
	count = 1
	result := t.current()
	if result.Type != EOF {
		t.idx += count
	} else {
		count = 0
	}
	return result, count
}

func (t *iterator) get(idx int) Token {
	if idx >= t.length {
		return Token{
			Type:  EOF,
			Value: "",
		}
	}
	return t.tokens[idx]
}

func (t *iterator) current() Token { return t.get(t.idx) }

func (t *iterator) unmove(count int) { t.idx -= count }
