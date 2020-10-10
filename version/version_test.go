package version

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

type testcase struct {
	vMajor, vMinor, vPatch uint16
	result                 uint64
}

func TestAsBigInt(t *testing.T) {
	require := require.New(t)

	prev := testcase{0, 0, 0, 0}
	for _, next := range []testcase{
		{0, 0, 1, 1},
		{0, 0, 2, 2},
		{0, 1, 0, 1000000},
		{0, 1, math.MaxUint16, 1065535},
		{1, 0, 0, 1000000000000},
		{1, 0, math.MaxUint16, 1000000065535},
		{1, 1, 0, 1000001000000},
		{2, 9, 9, 2000009000009},
		{3, 1, 0, 3000001000000},
		{math.MaxUint16, math.MaxUint16, math.MaxUint16, 65535065535065535},
	} {
		a := asU64(prev.vMajor, prev.vMinor, prev.vPatch)
		b := asU64(next.vMajor, next.vMinor, next.vPatch)
		require.Greater(b, a)
		require.Equal(b, next.result)
		prev = next
	}
}
