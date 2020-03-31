package ethapi

import (
	"context"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

// PublicTransactionPoolAPI

func TestPublicTransactionPoolAPI_FillTransaction(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	nonceLock := new(AddrLocker)
	api := NewPublicTransactionPoolAPI(b, nonceLock)
	assert.NotPanics(t, func() {
		gas := hexutil.Uint64(1)
		res, err := api.FillTransaction(ctx, SendTxArgs{
			From:     common.Address{1},
			To:       &common.Address{2},
			Gas:      &gas,
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
func TestPublicTransactionPoolAPI_GetBlockTransactionCountByHash(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	nonceLock := new(AddrLocker)
	api := NewPublicTransactionPoolAPI(b, nonceLock)
	assert.NotPanics(t, func() {
		_ = api.GetBlockTransactionCountByHash(ctx, common.Hash{1})
	})
}
func TestPublicTransactionPoolAPI_GetBlockTransactionCountByNumber(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	nonceLock := new(AddrLocker)
	api := NewPublicTransactionPoolAPI(b, nonceLock)
	assert.NotPanics(t, func() {
		_ = api.GetBlockTransactionCountByNumber(ctx, rpc.BlockNumber(1))
	})
}
func TestPublicTransactionPoolAPI_GetRawTransactionByBlockHashAndIndex(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	nonceLock := new(AddrLocker)
	api := NewPublicTransactionPoolAPI(b, nonceLock)
	assert.NotPanics(t, func() {
		_ = api.GetRawTransactionByBlockHashAndIndex(ctx, common.Hash{1}, hexutil.Uint(0))
	})
}
func TestPublicTransactionPoolAPI_GetRawTransactionByBlockNumberAndIndex(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	nonceLock := new(AddrLocker)
	api := NewPublicTransactionPoolAPI(b, nonceLock)
	assert.NotPanics(t, func() {
		_ = api.GetRawTransactionByBlockNumberAndIndex(ctx, rpc.BlockNumber(1), hexutil.Uint(0))
	})
}
func TestPublicTransactionPoolAPI_GetRawTransactionByHash(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	nonceLock := new(AddrLocker)
	api := NewPublicTransactionPoolAPI(b, nonceLock)
	assert.NotPanics(t, func() {
		res, err := api.GetRawTransactionByHash(ctx, common.Hash{1})
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
func TestPublicTransactionPoolAPI_GetTransactionByBlockHashAndIndex(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	nonceLock := new(AddrLocker)
	api := NewPublicTransactionPoolAPI(b, nonceLock)
	assert.NotPanics(t, func() {
		res := api.GetTransactionByBlockHashAndIndex(ctx, common.Hash{1}, hexutil.Uint(0))
		assert.NotEmpty(t, res)
	})
}
func TestPublicTransactionPoolAPI_GetTransactionByBlockNumberAndIndex(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	nonceLock := new(AddrLocker)
	api := NewPublicTransactionPoolAPI(b, nonceLock)
	assert.NotPanics(t, func() {
		res := api.GetTransactionByBlockNumberAndIndex(ctx, rpc.BlockNumber(1), hexutil.Uint(0))
		assert.NotEmpty(t, res)
	})
}
func TestPublicTransactionPoolAPI_GetTransactionByHash(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	nonceLock := new(AddrLocker)
	api := NewPublicTransactionPoolAPI(b, nonceLock)
	assert.NotPanics(t, func() {
		res, err := api.GetTransactionByHash(ctx, common.Hash{1})
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
func TestPublicTransactionPoolAPI_GetTransactionCount(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	nonceLock := new(AddrLocker)
	api := NewPublicTransactionPoolAPI(b, nonceLock)
	assert.NotPanics(t, func() {
		res, err := api.GetTransactionCount(ctx, common.Address{1}, rpc.BlockNumber(1))
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
func TestPublicTransactionPoolAPI_GetTransactionReceipt(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	nonceLock := new(AddrLocker)
	api := NewPublicTransactionPoolAPI(b, nonceLock)
	assert.NotPanics(t, func() {
		res, err := api.GetTransactionReceipt(ctx, common.Hash{1})
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
func TestPublicTransactionPoolAPI_PendingTransactions(t *testing.T) {
	b := NewTestBackend()

	nonceLock := new(AddrLocker)
	api := NewPublicTransactionPoolAPI(b, nonceLock)
	assert.NotPanics(t, func() {
		_, err := api.PendingTransactions()
		assert.NoError(t, err)
	})
}
func TestPublicTransactionPoolAPI_Resend(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	nonceLock := new(AddrLocker)
	api := NewPublicTransactionPoolAPI(b, nonceLock)
	assert.NotPanics(t, func() {
		gasPrice := hexutil.Big(*big.NewInt(1))
		gasLimit := hexutil.Uint64(1)
		nonce := hexutil.Uint64(1)
		gas := hexutil.Uint64(0)
		_, _ = api.Resend(ctx, SendTxArgs{
			From:     common.Address{1},
			To:       &common.Address{2},
			Nonce:	  &nonce,
			Gas:	  &gas,
		}, &gasPrice, &gasLimit)
	})
}
