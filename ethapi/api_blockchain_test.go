package ethapi

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"

	"github.com/Fantom-foundation/go-lachesis/evmcore"
)

// PublicBlockChainAPI

func TestPublicBlockChainAPI_BlockNumber(t *testing.T) {
	b := NewTestBackend()
	b.Returned("HeaderByNumber", &evmcore.EvmHeader{
		Number:     big.NewInt(1),
	})

	api := NewPublicBlockChainAPI(b)
	assert.NotPanics(t, func() {
		res := api.BlockNumber()
		assert.NotEmpty(t, res)
	})
}
func TestPublicBlockChainAPI_ChainID(t *testing.T) {
	b := NewTestBackend()
	b.Returned("ChainConfig", &params.ChainConfig{
		ChainID: big.NewInt(1),
	})

	api := NewPublicBlockChainAPI(b)
	assert.NotPanics(t, func() {
		res := api.ChainID()
		assert.NotEmpty(t, res)
	})
}
func TestPublicBlockChainAPI_EstimateGas(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()
	SetBackendStateDB(b)
	b.Returned("RPCGasCap", big.NewInt(1))

	api := NewPublicBlockChainAPI(b)
	assert.NotPanics(t, func() {
		_, _ = api.EstimateGas(ctx, CallArgs{})
	})
}
func TestPublicBlockChainAPI_GetBalance(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()
	stateDB := SetBackendStateDB(b)
	stateDB.AddBalance(common.Address{1}, big.NewInt(10))

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
	b.Returned("GetBlock", &evmcore.EvmBlock{
		EvmHeader:    evmcore.EvmHeader{
			Number:     big.NewInt(1),
			Hash:       common.Hash{2},
			ParentHash: common.Hash{3},
			Root:       common.Hash{4},
			TxHash:     common.Hash{5},
			Time:       6,
			Coinbase:   common.Address{7},
			GasLimit:   8,
			GasUsed:    9,
		},
		Transactions: nil,
	})
	b.Returned("GetTd", big.NewInt(1))

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
	b.Returned("BlockByNumber", &evmcore.EvmBlock{
		EvmHeader:    evmcore.EvmHeader{
			Number:     big.NewInt(1),
			Hash:       common.Hash{2},
			ParentHash: common.Hash{3},
			Root:       common.Hash{4},
			TxHash:     common.Hash{5},
			Time:       6,
			Coinbase:   common.Address{7},
			GasLimit:   8,
			GasUsed:    9,
		},
		Transactions: nil,
	})
	b.Returned("GetTd", big.NewInt(1))

	api := NewPublicBlockChainAPI(b)
	assert.NotPanics(t, func() {
		res, err := api.GetBlockByNumber(ctx, rpc.BlockNumber(1), true)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
func TestPublicBlockChainAPI_GetCode(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()
	stateDB := SetBackendStateDB(b)
	stateDB.AddBalance(common.Address{1}, big.NewInt(10))
	stateDB.SetCode(common.Address{1}, []byte{1, 2, 3})

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
	b.Returned("HeaderByHash", &evmcore.EvmHeader{
		Number:     big.NewInt(1),
	})
	b.Returned("GetTd", big.NewInt(1))

	api := NewPublicBlockChainAPI(b)
	assert.NotPanics(t, func() {
		res := api.GetHeaderByHash(ctx, common.HexToHash("0x1"))
		assert.NotEmpty(t, res)
	})
}
func TestPublicBlockChainAPI_GetHeaderByNumber(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()
	b.Returned("HeaderByNumber", &evmcore.EvmHeader{
		Number:     big.NewInt(1),
	})
	b.Returned("GetTd", big.NewInt(1))

	api := NewPublicBlockChainAPI(b)
	assert.NotPanics(t, func() {
		res, err := api.GetHeaderByNumber(ctx, rpc.BlockNumber(1))
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
func TestPublicBlockChainAPI_GetProof(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()
	stateDB := SetBackendStateDB(b)
	stateDB.AddBalance(common.Address{1}, big.NewInt(10))

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
	stateDB := SetBackendStateDB(b)
	stateDB.AddBalance(common.Address{1}, big.NewInt(10))

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
	b.Returned("GetBlock", &evmcore.EvmBlock{
		EvmHeader:    evmcore.EvmHeader{
			Number:     big.NewInt(1),
			Hash:       common.Hash{2},
			ParentHash: common.Hash{3},
			Root:       common.Hash{4},
			TxHash:     common.Hash{5},
			Time:       6,
			Coinbase:   common.Address{7},
			GasLimit:   8,
			GasUsed:    9,
		},
		Transactions: nil,
	})
	b.Returned("GetTd", big.NewInt(1))

	api := NewPublicBlockChainAPI(b)
	assert.NotPanics(t, func() {
		_, err := api.GetUncleByBlockHashAndIndex(ctx, common.Hash{1}, hexutil.Uint(1))
		assert.NoError(t, err)
	})
}
func TestPublicBlockChainAPI_GetUncleByBlockNumberAndIndex(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()
	b.Returned("BlockByNumber", &evmcore.EvmBlock{
		EvmHeader:    evmcore.EvmHeader{
			Number:     big.NewInt(1),
			Hash:       common.Hash{2},
			ParentHash: common.Hash{3},
			Root:       common.Hash{4},
			TxHash:     common.Hash{5},
			Time:       6,
			Coinbase:   common.Address{7},
			GasLimit:   8,
			GasUsed:    9,
		},
		Transactions: nil,
	})
	b.Returned("GetTd", big.NewInt(1))

	api := NewPublicBlockChainAPI(b)
	assert.NotPanics(t, func() {
		_, err := api.GetUncleByBlockNumberAndIndex(ctx, rpc.BlockNumber(1), hexutil.Uint(1))
		assert.NoError(t, err)
	})
}
func TestPublicBlockChainAPI_GetUncleCountByBlockHash(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()
	b.Returned("GetBlock", &evmcore.EvmBlock{
		EvmHeader:    evmcore.EvmHeader{
			Number:     big.NewInt(1),
			Hash:       common.Hash{2},
			ParentHash: common.Hash{3},
			Root:       common.Hash{4},
			TxHash:     common.Hash{5},
			Time:       6,
			Coinbase:   common.Address{7},
			GasLimit:   8,
			GasUsed:    9,
		},
		Transactions: nil,
	})
	b.Returned("GetTd", big.NewInt(1))

	api := NewPublicBlockChainAPI(b)
	assert.NotPanics(t, func() {
		api.GetUncleCountByBlockHash(ctx, common.Hash{1})
	})
}
func TestPublicBlockChainAPI_GetUncleCountByBlockNumber(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()
	b.Returned("BlockByNumber", &evmcore.EvmBlock{
		EvmHeader:    evmcore.EvmHeader{
			Number:     big.NewInt(1),
			Hash:       common.Hash{2},
			ParentHash: common.Hash{3},
			Root:       common.Hash{4},
			TxHash:     common.Hash{5},
			Time:       6,
			Coinbase:   common.Address{7},
			GasLimit:   8,
			GasUsed:    9,
		},
		Transactions: nil,
	})
	b.Returned("GetTd", big.NewInt(1))

	api := NewPublicBlockChainAPI(b)
	assert.NotPanics(t, func() {
		api.GetUncleCountByBlockNumber(ctx, rpc.BlockNumber(1))
	})
}
