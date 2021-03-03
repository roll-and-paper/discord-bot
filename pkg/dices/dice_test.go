package dices

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//var twentyFirsts = []int{6, 6, 3, 1, 2, 2, 4, 3, 3, 2, 2, 6, 4, 5, 3, 4, 3, 6, 6, 3}
//var tenFirsts = []int{6, 6, 3, 1, 2, 2, 4, 3, 3, 2}
//var fiveFirsts = []int{6, 6, 3, 1, 2}

func TestNewFromFaces_Roll(t *testing.T) {
	SetRandomSeed(42)
	dice := NewFromFaces([]int{2, 4, 6, 8, 10, 12})
	var results = make([]int, 0)
	for i := 0; i < 10; i++ {
		results = append(results, dice.Roll())
	}
	assert.Equal(t, dice.MaxValue(), 12)
	assert.Equal(t, dice.Faces(), []int{2, 4, 6, 8, 10, 12})
	assert.Equal(t, results, []int{12, 12, 6, 2, 4, 4, 8, 6, 6, 4})
}

func TestNewFromMax_Roll(t *testing.T) {
	SetRandomSeed(42)
	dice := NewFromMax(6)
	var results = make([]int, 0)
	for i := 0; i < 10; i++ {
		results = append(results, dice.Roll())
	}
	assert.Equal(t, dice.MaxValue(), 6)
	assert.Equal(t, dice.Faces(), []int{})
	assert.Equal(t, results, []int{6, 6, 3, 1, 2, 2, 4, 3, 3, 2})
}

func TestError(t *testing.T) {
	assert.Panics(t, func() { NewFromFaces(nil) }, "")
	assert.Panics(t, func() { NewFromMax(0) }, "")
}
