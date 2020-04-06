package inter

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"

	"github.com/Fantom-foundation/go-lachesis/hash"
)

func TestEventSerialization(t *testing.T) {
	assertar := assert.New(t)
	require := require.New(t)

	events := FakeFuzzingEvents()
	for i, e0 := range events {
		dsc := fmt.Sprintf("iter#%d", i)

		buf, err := rlp.EncodeToBytes(e0)
		require.NoError(err, dsc)

		assertar.Equal(len(buf), e0.Size())
		assertar.Equal(len(buf), e0.CalcSize())

		e1 := &Event{}
		err = rlp.DecodeBytes(buf, e1)
		require.NoError(err, dsc)

		if e1.Sig == nil {
			e1.Sig = []uint8{}
		}

		assertar.Equal(len(buf), e1.CalcSize())
		assertar.Equal(len(buf), e1.Size())

		require.Equal(e0.EventHeader, e1.EventHeader, dsc)
	}
}

func TestEventHash(t *testing.T) {
	require := require.New(t)
	var (
		events = FakeFuzzingEvents()
		hashes = make([]hash.Event, len(events))
	)

	t.Run("Calculation", func(t *testing.T) {
		for i, e := range events {
			hashes[i] = e.Hash()
		}
	})

	t.Run("Comparison", func(t *testing.T) {
		for i, e := range events {
			h := e.Hash()
			require.Equal(hashes[i], h, "Non-deterministic event hash detected")
			for _, other := range hashes[i+1:] {
				require.NotEqual(other, h, "Event hash collision detected")
			}
		}
	})
}
