package roller

type operation int

var (
	add    operation = 0
	minus  operation = 1
	mul    operation = 2
	div    operation = 3
	intDiv operation = 4
	pow    operation = 5
	mod    operation = 6
)

func applyOps(op operation, r1 Roller, r2 Roller, s string) *Result {
	switch op {
	case minus:
		return r1.Roll(s).Minus(r2.Roll(s))
	case mul:
		return r1.Roll(s).Mul(r2.Roll(s))
	case div:
		return r1.Roll(s).Div(r2.Roll(s))
	case intDiv:
		return r1.Roll(s).IntDiv(r2.Roll(s))
	case pow:
		return r1.Roll(s).Pow(r2.Roll(s))
	case mod:
		return r1.Roll(s).Mod(r2.Roll(s))
	default: // Add by default
		return r1.Roll(s).Add(r2.Roll(s))
	}
}
