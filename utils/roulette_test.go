package utils

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func Benchmark_StochasticPermutation(b *testing.B) {
	benchmarkPermutation(StochasticPermutation, b)
}

func Test_StochasticPermutation_correctness(t *testing.T) {
	testCorrectPermutation(StochasticPermutation, t, getTestWeightsIncreasing(1))
	testCorrectPermutation(StochasticPermutation, t, getTestWeightsIncreasing(1))
	testCorrectPermutation(StochasticPermutation, t, getTestWeightsEqual(1))
}

func Test_StochasticPermutation_determinism(t *testing.T) {
	require := require.New(t)

	weightsArr := getTestWeightsIncreasing(5)

	require.Equal([]int{1}, StochasticPermutation(len(weightsArr), weightsArr, hashOf(common.Hash{}, 0)))
	require.Equal([]int{0}, StochasticPermutation(len(weightsArr), weightsArr, hashOf(common.Hash{}, 1)))
	require.Equal([]int{4}, StochasticPermutation(len(weightsArr), weightsArr, hashOf(common.Hash{}, 3)))
	require.Equal([]int{4}, StochasticPermutation(len(weightsArr)/2, weightsArr, hashOf(common.Hash{}, 4)))
}

func Test_StochasticPermutation_concurency(t *testing.T) {
	testPermutationConcurency(StochasticPermutation, t)
}
