package version

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAsBigInt(t *testing.T) {
	require := require.New(t)

	prev := [3]uint64{0, 0, 0}
	for _, next := range [][3]uint64{
		{0, 0, 1},
		{0, 0, 2},
		{0, 1, 0},
		{0, 1, math.MaxUint64},
		{1, 0, 0},
		{1, 0, math.MaxUint64},
		{1, 1, 0},
		{2, 9, 9},
		{3, 1, 0},
		{math.MaxUint64, 0, 0},
	} {
		a := asBigInt(prev[0], prev[1], prev[2])
		b := asBigInt(next[0], next[1], next[2])
		require.Greater(b.Cmp(a), 0)
		require.LessOrEqual(len(b.Bytes()), 32)
		prev = next
	}
}
