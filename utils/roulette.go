package utils

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/Fantom-foundation/go-lachesis/inter/pos"
)

// rouletteSA implements the following paper:
// "Roulette-wheel selection via stochastic acceptance (Adam Liposki, Dorota Lipowska - 2011)"
type rouletteSA struct {
	weights   []pos.Stake
	maxWeight pos.Stake
	deterministicRand
}

func newRouletteSA(weights []pos.Stake, seed common.Hash) *rouletteSA {
	rw := &rouletteSA{
		weights: weights,
	}
	rw.seed = seed

	for _, w := range rw.weights {
		if rw.maxWeight < w {
			rw.maxWeight = w
		}
	}

	return rw
}

// NSelection randomly chooses a sample from the array of weights.
// Returns first {size} entries of {weights} permutation.
func (rw *rouletteSA) NSelection(size int) []int {
	permutation := make([]int, size)

	restLen := len(rw.weights)
	rest := make([]int, restLen)
	for i := 0; i < restLen; i++ {
		rest[i] = i
	}

	for i := 0; i < size; {
		nextOne := int(uint(rw.rand64()) % uint(restLen))
		selected := rest[nextOne]
		if pos.Stake(rw.rand64())%rw.maxWeight >= rw.weights[selected] {
			continue
		}
		permutation[i] = selected
		i++
		restLen--
		rest[nextOne] = rest[restLen]
	}

	return permutation
}
