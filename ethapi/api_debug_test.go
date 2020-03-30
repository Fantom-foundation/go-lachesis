package ethapi

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/Fantom-foundation/go-lachesis/evmcore"
)

// PublicDebugAPI

func TestPublicDebugAPI_GetBlockRlp(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()
	b.Returned("BlockByNumber", &evmcore.EvmBlock{
		EvmHeader:    evmcore.EvmHeader{
			Number:     big.NewInt(1),
		},
		Transactions: nil,
	})

	api := NewPublicDebugAPI(b)
	assert.NotPanics(t, func() {
		res, err := api.GetBlockRlp(ctx, 1)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
func TestPublicDebugAPI_PrintBlock(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()
	b.Returned("BlockByNumber", &evmcore.EvmBlock{
		EvmHeader:    evmcore.EvmHeader{
			Number:     big.NewInt(1),
		},
		Transactions: nil,
	})

	api := NewPublicDebugAPI(b)
	assert.NotPanics(t, func() {
		res, err := api.PrintBlock(ctx, 1)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
func TestPublicDebugAPI_SeedHash(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()
	b.Returned("BlockByNumber", &evmcore.EvmBlock{
		EvmHeader:    evmcore.EvmHeader{
			Number:     big.NewInt(1),
		},
		Transactions: nil,
	})

	api := NewPublicDebugAPI(b)
	assert.NotPanics(t, func() {
		res, err := api.SeedHash(ctx, 1)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
func TestPublicDebugAPI_TestSignCliqueBlock(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicDebugAPI(b)
	assert.NotPanics(t, func() {
		res, err := api.TestSignCliqueBlock(ctx, common.Address{}, 1)
		assert.Error(t, err)
		assert.NotEmpty(t, res)
	})
}

// PrivateDebugAPI

func TestPrivateDebugAPI_SetHead(t *testing.T) {
	b := NewTestBackend()

	api := NewPrivateDebugAPI(b)
	assert.NotPanics(t, func() {
		err := api.SetHead(1)
		assert.Error(t, err)
	})
}
