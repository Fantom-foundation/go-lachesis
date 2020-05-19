package utils

import (
	"math/rand"
	"sync"
)

// This file contains an implementation of the following paper
// Roulette-wheel selection via stochastic acceptance (Adam Liposki, Dorota Lipowska - 2011)
type RouletteSA struct {
	Weights   []uint64
	MaxWeight uint64
	seed      int64
	mu        sync.Mutex
}

func NewRouletteSA(weigths []uint64) *RouletteSA {
	if len(weigths) <= 0 {
		panic("the size must be positive")
	}

	max := GetMax(weigths)
	var average uint64 = weigths[0]
	for _, v := range weigths[1:] {
		average = +v
		average = average >> 1
	}
	return &RouletteSA{
		Weights:   weigths,
		MaxWeight: max,
		seed:      int64(average),
		mu:        sync.Mutex{},
	}
}

// NSelection randomly chooses a sample from the array of weights
// Returns first {size} entries of {weights} permutation.
func (rw *RouletteSA) NSelection(size int) []uint {
	rw.mu.Lock()
	defer rw.mu.Unlock()

	if len(rw.Weights) < size {
		panic("the permutation size must be less or equal to weights size")
	}
	if len(rw.Weights) == 0 {
		return make([]uint, 0)
	}
	selected := make(map[uint]bool)
	max := rw.MaxWeight

	selection := make([]uint, size)

	rand.Seed(rw.seed)

	for i := 0; i < size; i++ {
		// select a next one
		for {
			curSelection := rw.Selection(max)
			if _, ok := selected[curSelection]; !ok {
				selection[i] = curSelection
				selected[curSelection] = true
				break
			}
		}
	}
	return selection
}

// GetMax computes the largest value in an array
func GetMax(w []uint64) uint64 {
	if len(w) < 2 {
		if len(w) == 1 {
			return w[0]
		}
		return 0
	}
	max := w[0]
	for _, v := range w[1:] {
		if max < v {
			max = v
		}
	}
	return max
}

// Selection selects a single item randomly from the weighted items.
// param f_max is the maximum weight of the population
// returns index of the selected item
func (rw *RouletteSA) Selection(fMax uint64) uint {
	n := len(rw.Weights)

	for {
		// Select randomly one of the individuals
		i := rand.Intn(n)
		// The selection is accepted with probability fitness(i) / fMax
		if rand.Float64() < float64(rw.Weights[i])/float64(fMax) {
			return uint(i)
		}
	}
}

// Counter
func (rw *RouletteSA) Counter(nSelect int, fMax uint64) []uint {
	n := len(rw.Weights)
	var counter = make([]uint, n)
	for i := 0; i < nSelect; i++ {
		index := rw.Selection(fMax)
		counter[index]++
	}
	return counter
}
