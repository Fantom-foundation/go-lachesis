package ethapi

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
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
}
