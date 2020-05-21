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

	selected := make(map[int]bool)
	max := rw.maxWeight

	selection := make([]int, size)

	for i := 0; i < size; i++ {
		// select a next one
		for {
			curSelection := rw.selectOne(max)
			if _, ok := selected[curSelection]; !ok {
				selection[i] = curSelection
				selected[curSelection] = true
				break
			}
		}
	}
	return selection
}

// selectOne item randomly from the weighted items.
// Param f_max is the maximum weight of the population.
// Returns index of the selected item.
func (rw *rouletteSA) selectOne(fMax pos.Stake) int {
	for {
		// select randomly one of the individuals
		i := int(uint(rw.rand64()) % uint(len(rw.weights)))
		// the selection is accepted with probability fitness(i) / fMax
		randWeight := pos.Stake(rw.rand64()) % fMax
		if randWeight < rw.weights[i] {
			return i
		}
	}
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
