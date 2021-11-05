package roll

type Params interface {
	DiceSystemName() string
}

type GenericParams struct{ Expression string }

func (p *GenericParams) DiceSystemName() string { return "generic" }

type VampireDarkAgesParams struct {
	Dices          int
	Difficulty     int
	Specialisation int
}

func (p *VampireDarkAgesParams) DiceSystemName() string { return "vampire-dark-ages" }
