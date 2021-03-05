package parser

import "github.com/dohr-michael/roll-and-paper-bot/pkg/dices/roller"

func newIterator(tokens []Token) *iterator {
	return &iterator{
		tokens: tokens,
		length: len(tokens),
		idx:    0,
	}
}

func Parse(value string) (*roller.Roller, error) {
	//all, err := NewLexer(value).ScanAll()
	//if err != nil {
	//	return nil, err
	//}
	//a := &iterator{
	//	tokens:  all,
	//	size:    len(all),
	//	current: 0,
	//}
	return nil, nil
}

//func

type iterator struct {
	tokens []Token
	length int
	idx    int
}

func (t *iterator) currentAndMove(potentialCount ...int) (tok Token) {
	count := 1
	if len(potentialCount) > 0 {
		count = potentialCount[0]
	}
	result := t.current()
	if result.Type != EOF {
		t.idx += count
	}
	return result
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
