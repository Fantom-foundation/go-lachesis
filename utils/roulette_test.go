package utils

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// Test for the Roulette wheel selection implementation

func TestRouletteSA_GetMax(t *testing.T) {
	var w = []float64{0.1, 0.2, 0.15, 0.3, 0.25}
	var RW = NewRouletteSA(w)
	assert.Equal(t, RW.MaxWeight, 0.3)
}

func TestRouletteSA_Counter(t *testing.T) {
	var w = []float64{0.1, 0.2, 0.15, 0.3, 0.25}
	var RW = NewRouletteSA(w)
	counter := RW.Counter(1000, RW.MaxWeight)

	// assert
	assert.True(t, counter[0] < counter[2])
	assert.True(t, counter[2] < counter[1])
	assert.True(t, counter[1] < counter[4])
	assert.True(t, counter[4] < counter[3])
}

func TestRouletteSA_Selection(t *testing.T) {
	var w = []float64{0.1, 0.2, 0.15, 0.3, 0.25}
	var RW = NewRouletteSA(w)
	i := RW.Selection(RW.MaxWeight)

	assert.True(t, int(i) < len(w))
}


func TestRouletteSA_NSelection(t *testing.T) {
	var w = []float64{0.1, 0.2, 0.15, 0.3, 0.25}
	var RW = NewRouletteSA(w)
	selection := RW.NSelection(2)

	assert.Equal(t, 2, len(selection))
	assert.True(t, selection[0] != selection[1])
}

func TestRouletteSA_NSelectionThree(t *testing.T) {
	var w = []float64{0.1, 0.2, 0.15, 0.3, 0.25}
	var RW = NewRouletteSA(w)
	selection := RW.NSelection(3)

	assert.Equal(t, 3, len(selection))
	assert.True(t, selection[0] != selection[1])
	assert.True(t, selection[0] != selection[2])
	assert.True(t, selection[1] != selection[2])
}


func TestRouletteSA_Large_NSelectionThree(t *testing.T) {
	var w = []float64{
		0.1, 0.2, 0.15, 0.3, 0.25,
		0.6, 0.9, 1, 1.2, 3.4, 5.6}
	var RW = NewRouletteSA(w)
	selection := RW.NSelection(3)

	assert.Equal(t, 3, len(selection))
	assert.True(t, selection[0] != selection[1])
	assert.True(t, selection[0] != selection[2])
	assert.True(t, selection[1] != selection[2])
}