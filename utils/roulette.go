package utils

import (
	"math/rand"
	"time"
)

// This file contains an implementation of the following paper
// Roulette-wheel selection via stochastic acceptance (Adam Liposki, Dorota Lipowska - 2011)

// RouletteSA implementation.
type RouletteSA struct {
	Weights   []float64
	MaxWeight float64
}

// NewRouletteSA constructor.
func NewRouletteSA(w []float64) *RouletteSA {
	if len(w) <= 0 {
		panic("the size must be positive")
	}

	max := maxOf(w)
	return &RouletteSA{
		Weights:   w,
		MaxWeight: max,
	}
}

// NSelection randomly chooses a sample from the array of weights
// Returns first {size} entries of {weights} permutation.
func (rw *RouletteSA) NSelection(size int) []uint {
	if len(rw.Weights) < size {
		panic("the permutation size must be less or equal to weights size")
	}
	if len(rw.Weights) == 0 {
		return make([]uint, 0)
	}

	selected := make(map[uint]bool)
	max := rw.MaxWeight

	selection := make([]uint, size)
	for i := 0; i < size; i++ {
		// select a next one
		for {
			curSelection := rw.Selection(max)
			if _, ok := selected[curSelection]; !ok {
				selection[i] = curSelection
				selected[curSelection] = true
				//fmt.Printf("selection[%d]=%d\n", i, selection[i]);
				break
			}
		}
	}

	return selection
}

// maxOf computes the largest value in an array.
func maxOf(w []float64) float64 {
	if len(w) == 0 {
		return -1
	}
	max := w[0]
	for _, v := range w {
		if max < v {
			max = v
		}
	}
	return max
}

// Selection selects a single item randomly from the weighted items.
// param f_max is the maximum weight of the population
// returns index of the selected item
func (rw *RouletteSA) Selection(f_max float64) uint {
	n := len(rw.Weights)
	rand.Seed(time.Now().UnixNano())
	for {
		// Select randomly one of the individuals
		i := uint(rand.Intn(n))
		// The selection is accepted with probability fitness(i) / f_max
		if rand.Float64() < rw.Weights[i]/f_max {
			// fmt.Printf("i=%d\n", i)
			return i
		}
	}
}

// Counter
func (rw *RouletteSA) Counter(n_select int, f_max float64) []uint {
	n := len(rw.Weights)
	var counter = make([]uint, n)
	for i := 0; i < n_select; i++ {
		index := rw.Selection(f_max)
		counter[index]++
	}
	//for i := 0; i < n; i++ {
	//	fmt.Printf("counter[%d]=%d\n", i, counter[i]);
	//}
	return counter
}
