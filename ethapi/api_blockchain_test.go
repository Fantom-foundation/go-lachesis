package ethapi

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
)

// PublicBlockChainAPI

func TestPublicBlockChainAPI_Call(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	gas := hexutil.Uint64(0)
	gasPrice := hexutil.Big(*big.NewInt(0))
	value := hexutil.Big(*big.NewInt(0))
	data := hexutil.Bytes([]byte{1, 2, 3})
	code := hexutil.Bytes([]byte{1, 2, 3})
	balance := &hexutil.Big{}
	nonce := hexutil.Uint64(1)

	api := NewPublicBlockChainAPI(b)
	assert.NotPanics(t, func() {
		api.Call(ctx, CallArgs{
			Gas:      &gas,
			GasPrice: &gasPrice,
			Value:    &value,
			Data:     &data,
		}, rpc.BlockNumber(1), &map[common.Address]account{
			common.HexToAddress("0x0"): account{
				Nonce:   &nonce,
				Code:    &code,
				Balance: &balance,
				StateDiff: &map[common.Hash]common.Hash{
					common.Hash{1}: {1},
				},
			},
		})
		// assert.NoError(t, err)
	})
}
func TestPublicBlockChainAPI_BlockNumber(t *testing.T) {
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	assert.NotPanics(t, func() {
		res := api.BlockNumber()
		assert.NotEmpty(t, res)
	})
}
func TestPublicBlockChainAPI_ChainID(t *testing.T) {
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	assert.NotPanics(t, func() {
		res := api.ChainID()
		assert.NotEmpty(t, res)
	})
}
func TestPublicBlockChainAPI_EstimateGas(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	assert.NotPanics(t, func() {
		_, _ = api.EstimateGas(ctx, CallArgs{})
	})
}
func TestPublicBlockChainAPI_GetBalance(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	assert.NotPanics(t, func() {
		balance, err := api.GetBalance(ctx, common.Address{1}, rpc.BlockNumber(1))
		assert.NoError(t, err)
		assert.Equal(t, big.NewInt(10), balance.ToInt())
	})
}
func TestPublicBlockChainAPI_GetBlockByHash(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	assert.NotPanics(t, func() {
		res, err := api.GetBlockByHash(ctx, common.Hash{1}, true)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
func TestPublicBlockChainAPI_GetBlockByNumber(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	assert.NotPanics(t, func() {
		res, err := api.GetBlockByNumber(ctx, rpc.BlockNumber(rpc.PendingBlockNumber), true)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
func TestPublicBlockChainAPI_GetCode(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	assert.NotPanics(t, func() {
		res, err := api.GetCode(ctx, common.Address{1}, rpc.BlockNumber(1))
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
func TestPublicBlockChainAPI_GetHeaderByHash(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	assert.NotPanics(t, func() {
		res := api.GetHeaderByHash(ctx, common.HexToHash("0x1"))
		assert.NotEmpty(t, res)
	})
}
func TestPublicBlockChainAPI_GetHeaderByNumber(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	assert.NotPanics(t, func() {
		res, err := api.GetHeaderByNumber(ctx, rpc.BlockNumber(rpc.PendingBlockNumber))
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
func TestPublicBlockChainAPI_GetProof(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	assert.NotPanics(t, func() {
		res, err := api.GetProof(ctx, common.Address{1}, []string{"1"}, rpc.BlockNumber(1))
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
func TestPublicBlockChainAPI_GetStorageAt(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	assert.NotPanics(t, func() {
		res, err := api.GetStorageAt(ctx, common.Address{1}, "1", rpc.BlockNumber(1))
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
func TestPublicBlockChainAPI_GetUncleByBlockHashAndIndex(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	assert.NotPanics(t, func() {
		_, err := api.GetUncleByBlockHashAndIndex(ctx, common.Hash{1}, hexutil.Uint(1))
		assert.NoError(t, err)
	})
}
func TestPublicBlockChainAPI_GetUncleByBlockNumberAndIndex(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	assert.NotPanics(t, func() {
		_, err := api.GetUncleByBlockNumberAndIndex(ctx, rpc.BlockNumber(1), hexutil.Uint(1))
		assert.NoError(t, err)
	})
}
func TestPublicBlockChainAPI_GetUncleCountByBlockHash(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	assert.NotPanics(t, func() {
		api.GetUncleCountByBlockHash(ctx, common.Hash{1})
	})
}
func TestPublicBlockChainAPI_GetUncleCountByBlockNumber(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	assert.NotPanics(t, func() {
		api.GetUncleCountByBlockNumber(ctx, rpc.BlockNumber(1))
	})
}
