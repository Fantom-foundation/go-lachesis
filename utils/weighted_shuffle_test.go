package utils

import (
	"crypto/sha256"
	"github.com/Fantom-foundation/go-lachesis/inter/pos"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/Fantom-foundation/go-lachesis/common/littleendian"
)

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

// test that WeightedPermutation provides a correct permaition
func testCorrectPermutation(t *testing.T, weightsArr []pos.Stake) {
	assertar := assert.New(t)

	perm := WeightedPermutation(len(weightsArr), weightsArr, common.Hash{})
	assertar.Equal(len(weightsArr), len(perm))

	met := make(map[int]bool)
	for _, p := range perm {
		assertar.True(p >= 0)
		assertar.True(p < len(weightsArr))
		assertar.False(met[p])
		met[p] = true
	}
}

func Test_Permutation_correctness(t *testing.T) {
	testCorrectPermutation(t, getTestWeightsIncreasing(1))
	testCorrectPermutation(t, getTestWeightsIncreasing(30))
	testCorrectPermutation(t, getTestWeightsEqual(1000))
}

func hashOf(a common.Hash, b uint32) common.Hash {
	hasher := sha256.New()
	hasher.Write(a.Bytes())
	hasher.Write(littleendian.Int32ToBytes(uint32(b)))
	return common.BytesToHash(hasher.Sum(nil))
}

func Test_Permutation_determinism(t *testing.T) {
	weightsArr := getTestWeightsIncreasing(5)

	assertar := assert.New(t)

	assertar.Equal([]int{2, 3, 4, 1, 0}, WeightedPermutation(len(weightsArr), weightsArr, hashOf(common.Hash{}, 0)))
	assertar.Equal([]int{2, 3, 4, 1, 0}, WeightedPermutation(len(weightsArr), weightsArr, hashOf(common.Hash{}, 1)))
	assertar.Equal([]int{2, 3, 4, 1, 0}, WeightedPermutation(len(weightsArr), weightsArr, hashOf(common.Hash{}, 3)))
	assertar.Equal([]int{2, 3}, WeightedPermutation(len(weightsArr)/2, weightsArr, hashOf(common.Hash{}, 4)))
}

//func Test_Permutation_determinism_concurency(t *testing.T) {
//
//	assertar := assert.New(t)
//
//	weights := getTestWeightsIncreasing(10)
//	permutation := WeightedPermutation(len(weights), weights, hashOf(common.Hash{}, 0))
//	wg := sync.WaitGroup{}
//
//	for i := 0; i < 100; i++ {
//		wg.Add(1)
//		var tmpWeights = make([]pos.Stake, len(weights))
//		copy(tmpWeights, weights)
//
//		var tmpPermutation = make([]int, len(permutation))
//		copy(tmpPermutation, permutation)
//
//		go func(w []pos.Stake, perm []int) {
//			defer wg.Done()
//
//			p := WeightedPermutation(len(tmpWeights), tmpWeights, hashOf(common.Hash{}, 0))
//			assertar.Equal(p, tmpPermutation)
//		}(tmpWeights, tmpPermutation)
//	}
//	wg.Wait()
//}
