package utils

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/Fantom-foundation/go-lachesis/inter/pos"
)

// WeightedPermutation builds weighted random permutation
// Returns first {size} entries of {weights} permutation.
// Call with {size} == len(weights) to get the whole permutation.
func WeightedPermutation(size int, weights []pos.Stake, seed common.Hash) []int {
	if len(weights) < size {
		panic("the permutation size must be less or equal to weights size")
	}
	if len(weights) == 0 {
		return make([]int, 0)
	}

	tree := weightedShuffleTree{
		weights: weights,
		nodes:   make([]weightedShuffleNode, len(weights)),
		seed:    seed,
	}
	tree.build(0)

	permutation := make([]int, size)
	for i := 0; i < size; i++ {
		permutation[i] = tree.retrieve(0)
	}
	return permutation
}

// StochasticPermutation builds weighted random permutation
// Returns first {size} entries of {weights} permutation.
// Call with {size} == len(weights) to get the whole permutation.
func StochasticPermutation(size int, weights []pos.Stake, seed common.Hash) []int {
	if len(weights) < size {
		panic("the permutation size must be less or equal to weights size")
	}
	if len(weights) == 0 {
		return make([]int, 0)
	}

	wt := make([]uint64, len(weights))
	for i, v := range weights {
		wt[i] = uint64(v)
	}

	roulette := NewRouletteSA(wt)
	permutation := roulette.NSelection(size)

	result := make([]int, len(permutation))
	for i, v := range permutation {
		result[i] = int(v)
	}
	return result
}
