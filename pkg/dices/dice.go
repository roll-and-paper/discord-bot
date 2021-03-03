package dices

import "github.com/thoas/go-funk"

var (
	D4   = NewFromMax(4)
	D6   = NewFromMax(6)
	D8   = NewFromMax(8)
	D10  = NewFromMax(10)
	D12  = NewFromMax(12)
	D20  = NewFromMax(20)
	D100 = NewFromMax(100)
)

type Dice interface {
	MaxValue() int
	Faces() []int
	Roll() int
}

type dice struct {
	max      int
	maxValue int
	faces    []int
}

func (d *dice) MaxValue() int {
	return d.maxValue
}

func (d *dice) Faces() []int {
	return d.faces
}

func (d *dice) Roll() int {
	result := random.Intn(d.max)
	if len(d.faces) == 0 {
		return result + 1
	} else {
		return d.faces[result]
	}
}

func NewFromMax(max int) Dice {
	if max == 0 {
		panic("cannot create dice with zero face")
	}
	return &dice{maxValue: max, max: max, faces: []int{}}
}

func NewFromFaces(faces []int) Dice {
	if len(faces) == 0 {
		panic("cannot create dice with zero face")
	}
	return &dice{
		faces:    faces,
		max:      len(faces),
		maxValue: funk.MaxInt(faces).(int),
	}
}
