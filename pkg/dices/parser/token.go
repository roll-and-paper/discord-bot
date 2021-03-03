package parser

type TokenType int

const (
	// Special
	Illegal TokenType = iota
	EOF
	WS

	// General
	Number     //  [0-9]+(.[0-9]+)?
	BraceIn    // [
	BraceOut   // ]
	BracketIn  // (
	BracketOut // )
	DoubleDot  // ..
	Dot        // .
	Pipe       // |
	Dice       // d, D

	// Logic Operator
	GreaterOrEqual // >=
	LesserOrEqual  // <=
	GreaterThen    // >
	LesserThan     // <
	Equal          // =
	Different      // !=

	// Logic Operation
	And // &
	XOr // ^

	// Conditional Operation
	Modulo // %

	// Arithmetic Operation
	Plus           // +
	Minus          // -
	Multiplication // x, *
	Divide         // /, ÷
	Pow            // **

	// Options
	Keep           // k
	KeepAndExplode // K
	Sort           // s
	Explode        // e

	Count        // c
	Reroll       // r
	RerollUntil  // R
	RerollAndAdd // a

	Merge      // m
	IfOperator // i
	Painter    // p
	Filter     // f

	Split          // y
	Unique         // u
	AllSameExplode // t
	Group          // g

	Bind        // b
	Occurrences // o
)

type Token struct {
	Type  TokenType
	Value string
}
