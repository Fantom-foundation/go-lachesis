package utils

import (
	"github.com/Fantom-foundation/go-lachesis/inter/pos"
)

// rouletteSA implements the following paper:
// "Roulette-wheel selection via stochastic acceptance (Adam Liposki, Dorota Lipowska - 2011)"
type rouletteSA struct {
	weights   []pos.Stake
	maxWeight pos.Stake
	deterministicRand
}

// NSelection randomly chooses a sample from the array of weights.
// Returns first {size} entries of {weights} permutation.
func (rw *rouletteSA) NSelection(size int) []int {
	if len(rw.weights) < size {
		panic("the permutation size must be less or equal to weights size")
	}
	if len(rw.weights) == 0 {
		return make([]int, 0)
	}

	selection := make([]int, size)
	keep := make([]int, len(rw.weights))
	keepSize := len(rw.weights)
	for i := 0; i < len(rw.weights); i++ {
		keep[i] = i
	}

	for i := 0; i < size; i++ {
		// select a next one
		curSelection := int(uint(rw.rand64()) % uint(keepSize))
		selection[i] = keep[curSelection]

		last := keepSize - 1
		if curSelection < last {
			keep[curSelection], keep[last] = keep[last], keep[curSelection]
		}
		keepSize--
	}
	return selection
}

func maxOf(w []pos.Stake) pos.Stake {
	if len(w) == 1 {
		return w[0]
	}
	if len(w) < 2 {
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
