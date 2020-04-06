package topicsdb

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPosToBytes(t *testing.T) {
	require := require.New(t)

	for expect := uint8(0); ; /*see break*/ expect += 0x0f {
		bb := posToBytes(expect)
		got := bytesToPos(bb)

		require.Equal(expect, got)

		if expect == 0xff {
			break
		}
	}
}
