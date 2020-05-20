package utils

import (
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/Fantom-foundation/go-lachesis/inter/pos"
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

func Test_StochasticPermutation_concurency(t *testing.T) {
	assertar := assert.New(t)

	weights := getTestWeightsIncreasing(10)
	expect := StochasticPermutation(len(weights), weights, hashOf(common.Hash{}, 0))

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {

		tmpWeights := make([]pos.Stake, len(weights))
		copy(tmpWeights, weights)

		wg.Add(1)
		go func(w []pos.Stake) {
			defer wg.Done()

			got := StochasticPermutation(len(tmpWeights), tmpWeights, hashOf(common.Hash{}, 0))
			assertar.Equal(got, expect)
		}(tmpWeights)
	}

	wg.Wait()
}
