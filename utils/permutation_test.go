package utils

import (
	"crypto/sha256"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/Fantom-foundation/go-lachesis/common/littleendian"
	"github.com/Fantom-foundation/go-lachesis/inter/pos"
)

type permutator func(size int, weights []pos.Stake, seed common.Hash) []int

// test that WeightedPermutation provides a correct permutaition
func testCorrectPermutation(f permutator, t *testing.T, weightsArr []pos.Stake) {
	require := require.New(t)

	perm := f(len(weightsArr), weightsArr, common.Hash{})
	require.Equal(len(weightsArr), len(perm))

	met := make(map[int]bool)
	for _, p := range perm {
		require.True(p >= 0)
		require.True(p < len(weightsArr))
		require.False(met[p])
		met[p] = true
	}
}

func testPermutationConcurency(f permutator, t *testing.T) {
	require := require.New(t)

	weights := getTestWeightsIncreasing(10)
	expect := f(len(weights), weights, hashOf(common.Hash{}, 0))

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			got := f(len(weights), weights, hashOf(common.Hash{}, 0))
			require.Equal(got, expect)
		}()
	}

	wg.Wait()
}

func getTestWeightsIncreasing(num int) []pos.Stake {
	weights := make([]pos.Stake, num)
	for i := 0; i < num; i++ {
		weights[i] = pos.Stake(i+1) * 1000
	}
	return weights
}

func getTestWeightsEqual(num int) []pos.Stake {
	weights := make([]pos.Stake, num)
	for i := 0; i < num; i++ {
		weights[i] = 1000
	}
	return weights
}

func hashOf(a common.Hash, b uint32) common.Hash {
	hasher := sha256.New()
	hasher.Write(a.Bytes())
	hasher.Write(littleendian.Int32ToBytes(uint32(b)))
	return common.BytesToHash(hasher.Sum(nil))
}
