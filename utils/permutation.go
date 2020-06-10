package utils

import (
	"crypto/sha256"

	"github.com/ethereum/go-ethereum/common"

	"github.com/Fantom-foundation/go-lachesis/common/littleendian"
	"github.com/Fantom-foundation/go-lachesis/inter/pos"
)

type deterministicRand struct {
	seed      common.Hash
	seedIndex int
}

func (r *deterministicRand) rand64() uint64 {
	if r.seedIndex == 32 {
		hasher := sha256.New() // use sha2 instead of sha3 for speed
		hasher.Write(r.seed.Bytes())
		r.seed = common.BytesToHash(hasher.Sum(nil))
		r.seedIndex = 0
	}
	// use not used parts of old seed, instead of calculating new one
	res := littleendian.BytesToInt64(r.seed[r.seedIndex : r.seedIndex+8])
	r.seedIndex += 8
	return res
}

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
	}
	tree.seed = seed
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

	roulette := &rouletteSA{
		weights: weights,
	}
	roulette.seed = seed

	result := roulette.NSelection(size)

	return result
}
