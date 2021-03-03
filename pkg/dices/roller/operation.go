package roller

type operation func(Roller, Roller) func(string) *Result

func add(r1 Roller, r2 Roller) func(string) *Result {
	return func(s string) *Result { return r1.Roll(s).Add(r2.Roll(s)) }
}
func minus(r1 Roller, r2 Roller) func(string) *Result {
	return func(s string) *Result { return r1.Roll(s).Minus(r2.Roll(s)) }
}
func mul(r1 Roller, r2 Roller) func(string) *Result {
	return func(s string) *Result { return r1.Roll(s).Mul(r2.Roll(s)) }
}
func div(r1 Roller, r2 Roller) func(string) *Result {
	return func(s string) *Result { return r1.Roll(s).Div(r2.Roll(s)) }
}
func intDiv(r1 Roller, r2 Roller) func(string) *Result {
	return func(s string) *Result { return r1.Roll(s).IntDiv(r2.Roll(s)) }
}
func pow(r1 Roller, r2 Roller) func(string) *Result {
	return func(s string) *Result { return r1.Roll(s).Pow(r2.Roll(s)) }
}
func mod(r1 Roller, r2 Roller) func(string) *Result {
	return func(s string) *Result { return r1.Roll(s).Mod(r2.Roll(s)) }
}
