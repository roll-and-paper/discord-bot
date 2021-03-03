package parser

import "github.com/dohr-michael/roll-and-paper-bot/pkg/dices/roller"

func Parse(value string) (*roller.Roller, error) {
	all := NewScanner(value).ScanAll()
	a := &analyser{
		tokens:  all,
		size:    len(all),
		current: 0,
	}
	return a.Parse()
}

type analyser struct {
	tokens  []Token
	size    int
	current int
}

func (a *analyser) Parse() (*roller.Roller, error) {

	return nil, nil
}

func (a *analyser) scan() (tok Token) {
	if a.current >= a.size {
		return Token{EOF, ""}
	}
	tok = a.tokens[a.current]
	a.current++
	return
}

func (a *analyser) scanIgnoreWhitespace() (tok Token) {
	tok = a.scan()
	if tok.Type == WS {
		tok = a.scan()
	}
	return
}

func (a *analyser) unscan() { a.current-- }
