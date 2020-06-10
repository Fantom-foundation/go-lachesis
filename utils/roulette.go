package utils

import (
	"github.com/Fantom-foundation/go-lachesis/inter/pos"
)

// rouletteSA implements the following paper:
// "Roulette-wheel selection via stochastic acceptance (Adam Liposki, Dorota Lipowska - 2011)"
type rouletteSA struct {
	weights []pos.Stake
	deterministicRand
}

// NSelection randomly chooses a sample from the array of weights.
// Returns first {size} entries of {weights} permutation.
func (rw *rouletteSA) NSelection(size int) []int {
	// weights size is checked before
	permutation := make([]int, size)

	restLen := len(rw.weights)
	rest := make([]int, restLen)
	for i := 0; i < restLen; i++ {
		rest[i] = i
	}

	for i := 0; i < size; i++ {
		nextOne := int(uint(rw.rand64()) % uint(restLen))
		permutation[i] = rest[nextOne]
		restLen--
		rest[nextOne] = rest[restLen]
	}

	return permutation
}
