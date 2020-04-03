package ethapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Fantom-foundation/go-lachesis/hash"
)

// PublicEthereumAPI

func TestPublicEthereumAPI_GasPrice(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicEthereumAPI(b)
	require.NotPanics(t, func() {
		res, err := api.GasPrice(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
}

func TestPublicEthereumAPI_ProtocolVersion(t *testing.T) {
	b := NewTestBackend()

	api := NewPublicEthereumAPI(b)
	require.NotPanics(t, func() {
		res := api.ProtocolVersion()
		require.NotEmpty(t, res)
	})
}

func TestPublicEthereumAPI_Syncing(t *testing.T) {
	b := NewTestBackend()

	api := NewPublicEthereumAPI(b)
	require.NotPanics(t, func() {
		res, err := api.Syncing()
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
	require.NotPanics(t, func() {
		b.Returned("Progress", PeerProgress{
			CurrentEpoch:     1,
			CurrentBlock:     2,
			CurrentBlockHash: hash.Event{3},
			CurrentBlockTime: 1,
			HighestBlock:     5,
			HighestEpoch:     6,
		})
		_, err := api.Syncing()
		require.NoError(t, err)
	})
}
