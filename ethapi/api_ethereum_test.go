package ethapi

import (
	"context"
	"github.com/Fantom-foundation/go-lachesis/hash"
	"testing"

	"github.com/stretchr/testify/assert"
)

// PublicEthereumAPI

func TestPublicEthereumAPI_GasPrice(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicEthereumAPI(b)
	assert.NotPanics(t, func() {
		res, err := api.GasPrice(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}

func TestPublicEthereumAPI_ProtocolVersion(t *testing.T) {
	b := NewTestBackend()

	api := NewPublicEthereumAPI(b)
	assert.NotPanics(t, func() {
		res := api.ProtocolVersion()
		assert.NotEmpty(t, res)
	})
}

func TestPublicEthereumAPI_Syncing(t *testing.T) {
	b := NewTestBackend()

	api := NewPublicEthereumAPI(b)
	assert.NotPanics(t, func() {
		res, err := api.Syncing()
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
	assert.NotPanics(t, func() {
		b.Returned("Progress", PeerProgress{
			CurrentEpoch:     1,
			CurrentBlock:     2,
			CurrentBlockHash: hash.Event{3},
			CurrentBlockTime: 1,
			HighestBlock:     5,
			HighestEpoch:     6,
		})
		_, err := api.Syncing()
		assert.NoError(t, err)
	})
}
