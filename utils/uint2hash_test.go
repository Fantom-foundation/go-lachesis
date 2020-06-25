package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_H256toU64(t *testing.T) {
	for _, exp := range []uint64{
		0, 1, 250, 0xffffffffffffffff,
	} {
		h := U64to256(exp)
		got := H256toU64(h)
		require.Equal(t, exp, got)
	}
}
