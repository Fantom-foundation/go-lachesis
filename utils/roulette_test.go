package utils

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func Test_StochasticPermutation_correctness(t *testing.T) {
	testCorrectPermutation(StochasticPermutation, t, getTestWeightsIncreasing(1))
	testCorrectPermutation(StochasticPermutation, t, getTestWeightsIncreasing(30))
	testCorrectPermutation(StochasticPermutation, t, getTestWeightsEqual(1000))
}

func Test_StochasticPermutation_determinism(t *testing.T) {
	weightsArr := getTestWeightsIncreasing(5)

	assertar := assert.New(t)

	assertar.Equal([]int{2, 3, 4, 1, 0}, StochasticPermutation(len(weightsArr), weightsArr, hashOf(common.Hash{}, 0)))
	assertar.Equal([]int{2, 3, 4, 1, 0}, StochasticPermutation(len(weightsArr), weightsArr, hashOf(common.Hash{}, 1)))
	assertar.Equal([]int{2, 3, 4, 1, 0}, StochasticPermutation(len(weightsArr), weightsArr, hashOf(common.Hash{}, 3)))
	assertar.Equal([]int{2, 3}, StochasticPermutation(len(weightsArr)/2, weightsArr, hashOf(common.Hash{}, 4)))
}

//func Test_StochasticPermutation_determinism_concurency(t *testing.T) {
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
