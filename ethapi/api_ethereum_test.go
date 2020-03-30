package ethapi

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
)

// PublicEthereumAPI

func TestPublicEthereumAPI_GasPrice(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()
	b.Returned("SuggestPrice", big.NewInt(1))

	api := NewPublicEthereumAPI(b)
	assert.NotPanics(t, func() {
		res, err := api.GasPrice(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}

func TestPublicEthereumAPI_ProtocolVersion(t *testing.T) {
	b := NewTestBackend()
	b.Returned("ProtocolVersion", 1)

	api := NewPublicEthereumAPI(b)
	assert.NotPanics(t, func() {
		res := api.ProtocolVersion()
		assert.NotEmpty(t, res)
	})
}

func TestPublicEthereumAPI_Syncing(t *testing.T) {
	b := NewTestBackend()
	ts := inter.Timestamp(time.Now().Add(-91 * time.Minute).UnixNano())
	b.Returned("Progress", PeerProgress{
		CurrentEpoch:     1,
		CurrentBlock:     2,
		CurrentBlockHash: hash.Event{3},
		CurrentBlockTime: ts,
		HighestBlock:     5,
		HighestEpoch:     6,
	})

	api := NewPublicEthereumAPI(b)
	assert.NotPanics(t, func() {
		res, err := api.Syncing()
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
