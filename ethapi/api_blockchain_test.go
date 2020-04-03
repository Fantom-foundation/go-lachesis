package ethapi

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
)

var (
	ErrBackendTest = errors.New("backend test error")
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
	require.NotPanics(t, func() {
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
		// require.NoError(t, err)
	})
}
func TestPublicBlockChainAPI_BlockNumber(t *testing.T) {
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	require.NotPanics(t, func() {
		res := api.BlockNumber()
		require.NotEmpty(t, res)
	})
}
func TestPublicBlockChainAPI_ChainID(t *testing.T) {
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	require.NotPanics(t, func() {
		res := api.ChainID()
		require.NotEmpty(t, res)
	})
}
func TestPublicBlockChainAPI_EstimateGas(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	require.NotPanics(t, func() {
		_, _ = api.EstimateGas(ctx, CallArgs{})
	})
}
func TestPublicBlockChainAPI_GetBalance(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	require.NotPanics(t, func() {
		balance, err := api.GetBalance(ctx, common.Address{1}, rpc.BlockNumber(1))
		require.NoError(t, err)
		require.Equal(t, big.NewInt(10), balance.ToInt())
	})
	require.NotPanics(t, func() {
		b.Error("StateAndHeaderByNumber", ErrBackendTest)
		_, err := api.GetBalance(ctx, common.Address{1}, rpc.BlockNumber(1))
		require.Error(t, err)
	})
}
func TestPublicBlockChainAPI_GetBlockByHash(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	require.NotPanics(t, func() {
		res, err := api.GetBlockByHash(ctx, common.Hash{1}, true)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
	require.NotPanics(t, func() {
		b.Returned("GetBlock", nil)
		b.Error("GetBlock", ErrBackendTest)
		res, err := api.GetBlockByHash(ctx, common.Hash{1}, true)
		require.Error(t, err)
		require.Empty(t, res)
	})
}
func TestPublicBlockChainAPI_GetBlockByNumber(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	require.NotPanics(t, func() {
		res, err := api.GetBlockByNumber(ctx, rpc.BlockNumber(rpc.PendingBlockNumber), true)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
	require.NotPanics(t, func() {
		b.Error("BlockByNumber", ErrBackendTest)
		_, err := api.GetBlockByNumber(ctx, rpc.BlockNumber(rpc.PendingBlockNumber), true)
		require.Error(t, err)
	})
}
func TestPublicBlockChainAPI_GetCode(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	require.NotPanics(t, func() {
		res, err := api.GetCode(ctx, common.Address{1}, rpc.BlockNumber(1))
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
	require.NotPanics(t, func() {
		b.Error("StateAndHeaderByNumber", ErrBackendTest)
		res, err := api.GetCode(ctx, common.Address{1}, rpc.BlockNumber(1))
		require.Error(t, err)
		require.Empty(t, res)
	})
}
func TestPublicBlockChainAPI_GetHeaderByHash(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	require.NotPanics(t, func() {
		res := api.GetHeaderByHash(ctx, common.HexToHash("0x1"))
		require.NotEmpty(t, res)
	})
	require.NotPanics(t, func() {
		b.Returned("HeaderByHash", nil)
		res := api.GetHeaderByHash(ctx, common.HexToHash("0x1"))
		require.Empty(t, res)
	})
}
func TestPublicBlockChainAPI_GetHeaderByNumber(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	require.NotPanics(t, func() {
		res, err := api.GetHeaderByNumber(ctx, rpc.BlockNumber(rpc.PendingBlockNumber))
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
	require.NotPanics(t, func() {
		b.Error("HeaderByNumber", ErrBackendTest)
		_, err := api.GetHeaderByNumber(ctx, rpc.BlockNumber(rpc.PendingBlockNumber))
		require.Error(t, err)
	})
}
func TestPublicBlockChainAPI_GetProof(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	require.NotPanics(t, func() {
		res, err := api.GetProof(ctx, common.Address{1}, []string{"1"}, rpc.BlockNumber(1))
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
	require.NotPanics(t, func() {
		b.Error("StateAndHeaderByNumber", ErrBackendTest)
		_, err := api.GetProof(ctx, common.Address{1}, []string{"1"}, rpc.BlockNumber(1))
		require.Error(t, err)
	})
}
func TestPublicBlockChainAPI_GetStorageAt(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	require.NotPanics(t, func() {
		res, err := api.GetStorageAt(ctx, common.Address{1}, "1", rpc.BlockNumber(1))
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
	require.NotPanics(t, func() {
		b.Error("StateAndHeaderByNumber", ErrBackendTest)
		res, err := api.GetStorageAt(ctx, common.Address{1}, "1", rpc.BlockNumber(1))
		require.Error(t, err)
		require.Empty(t, res)
	})
}
func TestPublicBlockChainAPI_GetUncleByBlockHashAndIndex(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	require.NotPanics(t, func() {
		_, err := api.GetUncleByBlockHashAndIndex(ctx, common.Hash{1}, hexutil.Uint(1))
		require.NoError(t, err)
	})
	require.NotPanics(t, func() {
		b.Returned("GetBlock", nil)
		b.Error("GetBlock", ErrBackendTest)
		_, err := api.GetUncleByBlockHashAndIndex(ctx, common.Hash{1}, hexutil.Uint(1))
		require.Error(t, err)
	})
}
func TestPublicBlockChainAPI_GetUncleByBlockNumberAndIndex(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	require.NotPanics(t, func() {
		_, err := api.GetUncleByBlockNumberAndIndex(ctx, rpc.BlockNumber(1), hexutil.Uint(1))
		require.NoError(t, err)
	})
	require.NotPanics(t, func() {
		b.Returned("BlockByNumber", nil)
		b.Error("BlockByNumber", ErrBackendTest)
		_, err := api.GetUncleByBlockNumberAndIndex(ctx, rpc.BlockNumber(1), hexutil.Uint(1))
		require.Error(t, err)
	})
}
func TestPublicBlockChainAPI_GetUncleCountByBlockHash(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	require.NotPanics(t, func() {
		api.GetUncleCountByBlockHash(ctx, common.Hash{1})
	})
	require.NotPanics(t, func() {
		b.Returned("GetBlock", nil)
		b.Error("GetBlock", ErrBackendTest)
		api.GetUncleCountByBlockHash(ctx, common.Hash{1})
	})
}
func TestPublicBlockChainAPI_GetUncleCountByBlockNumber(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicBlockChainAPI(b)
	require.NotPanics(t, func() {
		api.GetUncleCountByBlockNumber(ctx, rpc.BlockNumber(1))
	})
	require.NotPanics(t, func() {
		b.Returned("BlockByNumber", nil)
		b.Error("BlockByNumber", ErrBackendTest)
		api.GetUncleCountByBlockNumber(ctx, rpc.BlockNumber(1))
	})
}
