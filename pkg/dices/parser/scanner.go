package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

const eof = rune(0)

var simpleBinding = map[rune]TokenType{
	'[': BraceIn,
	']': BraceOut,
	'(': BracketIn,
	')': BracketOut,
	'|': Pipe,
	'd': Dice,
	'D': Dice,
	'=': Equal,
	'&': And,
	'^': XOr,
	'%': Modulo,
	'+': Plus,
	'-': Minus,
	'x': Multiplication,
	'/': Divide,
	'÷': Divide,
	'k': Keep,
	'K': KeepAndExplode,
	's': Sort,
	'e': Explode,
	'c': Count,
	'r': Reroll,
	'R': RerollUntil,
	'a': RerollAndAdd,
	'm': Merge,
	'i': IfOperator,
	'p': Painter,
	'f': Filter,
	'y': Split,
	'u': Unique,
	't': AllSameExplode,
	'g': Group,
	'b': Bind,
	'o': Occurrences,
}

type Scanner struct{ r *bufio.Reader }

func NewLexer(value string) *Scanner { return &Scanner{r: bufio.NewReader(strings.NewReader(value))} }

func (s *Scanner) ScanAll() ([]Token, error) {
	result := make([]Token, 0)
	for {
		t, v := s.scan()
		if t == EOF {
			break
		}
		result = append(result, Token{t, v})
		if t == Illegal {
			return nil, fmt.Errorf("illegal token %s ", v)
		}
	}
	return result, nil
}

func (s *Scanner) scan() (TokenType, string) {
	ch := s.read()
	// read
	if isWhitespace(ch) {
		s.unread()
		return s.readWhitespace()
	} else if isDigit(ch) {
		s.unread()
		return s.readNumber()
	} else if t, ok := simpleBinding[ch]; ok {
		return t, string(ch)
	} else if ch == '>' {
		if next := s.read(); next == '=' {
			return GreaterOrEqual, ">="
		}
		s.unread()
		return GreaterThen, string(ch)
	} else if ch == '<' {
		if next := s.read(); next == '=' {
			return LesserOrEqual, "<="
		}
		s.unread()
		return LesserThan, string(ch)
	} else if ch == '!' {
		if next := s.read(); next == '=' {
			return Different, "!="
		}
		s.unread()
	} else if ch == '*' {
		if next := s.read(); next == '*' {
			return Pow, "**"
		}
		s.unread()
		return Multiplication, string(ch)
	} else if ch == '.' {
		if next := s.read(); next == '.' {
			return DoubleDot, ".."
		}
		s.unread()
		return Dot, "."
	} else if ch == eof {
		return EOF, ""
	}
	return Illegal, string(ch)
}

func (s *Scanner) readNumber() (TokenType, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())
	for {
		if ch := s.read(); ch == eof {
			break
		} else if isDigit(ch) {
			buf.WriteRune(ch)
		} else {
			s.unread()
			break
		}
	}
	return Number, buf.String()
}

func (s *Scanner) readWhitespace() (TokenType, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}
	return WS, buf.String()
}

func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }
func isDigit(ch rune) bool      { return ch >= '0' && ch <= '9' }

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) unread() { _ = s.r.UnreadRune() }
