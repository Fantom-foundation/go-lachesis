package ethapi

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
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
			From: common.Address{1},
			To:   &common.Address{2},
			Gas:  &gas,
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
	assert.NotPanics(t, func() {
		b.Returned("GetBlock", nil)
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
	assert.NotPanics(t, func() {
		b.Returned("GetBlock", nil)
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
	assert.NotPanics(t, func() {
		b.Returned("GetBlock", nil)
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
	assert.NotPanics(t, func() {
		b.Returned("BlockByNumber", nil)
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
	assert.NotPanics(t, func() {
		b.Returned("GetTransaction", nil)
		res, err := api.GetRawTransactionByHash(ctx, common.Hash{1})
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
	assert.NotPanics(t, func() {
		b.Returned("GetTransaction", nil)
		b.Returned("GetPoolTransaction", nil)
		res, err := api.GetRawTransactionByHash(ctx, common.Hash{1})
		assert.NoError(t, err)
		assert.Empty(t, res)
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
	assert.NotPanics(t, func() {
		b.Returned("GetBlock", nil)
		res := api.GetTransactionByBlockHashAndIndex(ctx, common.Hash{1}, hexutil.Uint(0))
		assert.Empty(t, res)
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
	assert.NotPanics(t, func() {
		b.Returned("BlockByNumber", nil)
		res := api.GetTransactionByBlockNumberAndIndex(ctx, rpc.BlockNumber(1), hexutil.Uint(0))
		assert.Empty(t, res)
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
	assert.NotPanics(t, func() {
		b.Returned("GetTransaction", nil)
		res, err := api.GetTransactionByHash(ctx, common.Hash{1})
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
	assert.NotPanics(t, func() {
		b.PrepareMethods()
		b.Error("GetTransaction", ErrBackendTest)
		res, err := api.GetTransactionByHash(ctx, common.Hash{1})
		assert.Error(t, err)
		assert.Empty(t, res)
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
	assert.NotPanics(t, func() {
		res, err := api.GetTransactionCount(ctx, common.Address{1}, rpc.BlockNumber(rpc.PendingBlockNumber))
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
	assert.NotPanics(t, func() {
		b.Error("StateAndHeaderByNumber", ErrBackendTest)
		res, err := api.GetTransactionCount(ctx, common.Address{1}, rpc.BlockNumber(1))
		assert.Error(t, err)
		assert.Empty(t, res)
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
	assert.NotPanics(t, func() {
		rec1 := types.NewReceipt([]byte{}, false, 100)
		rec1.PostState = []byte{1, 2, 3}
		rec1.ContractAddress = common.Address{1}
		b.Returned("GetReceiptsByNumber", types.Receipts{
			rec1,
		})
		res, err := api.GetTransactionReceipt(ctx, common.Hash{1})
		assert.NoError(t, err)
		assert.Empty(t, res)
	})
	assert.NotPanics(t, func() {
		rec1 := types.NewReceipt([]byte{}, false, 100)
		rec1.PostState = []byte{}
		rec1.ContractAddress = common.Address{1}
		b.Returned("GetReceiptsByNumber", types.Receipts{
			rec1,
		})
		res, err := api.GetTransactionReceipt(ctx, common.Hash{1})
		assert.NoError(t, err)
		assert.Empty(t, res)
	})
	assert.NotPanics(t, func() {
		b.Error("GetReceiptsByNumber", ErrBackendTest)
		b.Returned("GetReceiptsByNumber", nil)
		res, err := api.GetTransactionReceipt(ctx, common.Hash{1})
		assert.NoError(t, err)
		assert.Empty(t, res)
	})
	assert.NotPanics(t, func() {
		b.Error("HeaderByNumber", ErrBackendTest)
		b.Returned("HeaderByNumber", nil)
		res, err := api.GetTransactionReceipt(ctx, common.Hash{1})
		assert.NoError(t, err)
		assert.Empty(t, res)
	})
	assert.NotPanics(t, func() {
		b.Error("GetTransaction", ErrBackendTest)
		b.Returned("GetTransaction", nil)
		res, err := api.GetTransactionReceipt(ctx, common.Hash{1})
		assert.NoError(t, err)
		assert.Empty(t, res)
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
			From:  common.Address{1},
			To:    &common.Address{2},
			Nonce: &nonce,
			Gas:   &gas,
		}, &gasPrice, &gasLimit)
	})
}
func TestPublicTransactionPoolAPI_SignTransaction(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	nonceLock := new(AddrLocker)
	api := NewPublicTransactionPoolAPI(b, nonceLock)
	assert.NotPanics(t, func() {
		gas := hexutil.Uint64(0)
		gasPrice := hexutil.Big(*big.NewInt(0))
		nonce := hexutil.Uint64(1)
		res, err := api.SignTransaction(ctx, SendTxArgs{
			From:     common.Address{1},
			To:       &common.Address{2},
			Gas:      &gas,
			GasPrice: &gasPrice,
			Nonce:    &nonce,
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
	assert.NotPanics(t, func() {
		gasPrice := hexutil.Big(*big.NewInt(0))
		nonce := hexutil.Uint64(1)
		res, err := api.SignTransaction(ctx, SendTxArgs{
			From:     common.Address{1},
			To:       &common.Address{2},
			Gas:      nil,
			GasPrice: &gasPrice,
			Nonce:    &nonce,
		})
		assert.Error(t, err)
		assert.Empty(t, res)
	})
	assert.NotPanics(t, func() {
		gas := hexutil.Uint64(0)
		nonce := hexutil.Uint64(1)
		res, err := api.SignTransaction(ctx, SendTxArgs{
			From:     common.Address{1},
			To:       &common.Address{2},
			Gas:      &gas,
			GasPrice: nil,
			Nonce:    &nonce,
		})
		assert.Error(t, err)
		assert.Empty(t, res)
	})
	assert.NotPanics(t, func() {
		gas := hexutil.Uint64(0)
		gasPrice := hexutil.Big(*big.NewInt(0))
		res, err := api.SignTransaction(ctx, SendTxArgs{
			From:     common.Address{1},
			To:       &common.Address{2},
			Gas:      &gas,
			GasPrice: &gasPrice,
			Nonce:    nil,
		})
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}
func TestPublicTransactionPoolAPI_SendTransaction(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	nonceLock := new(AddrLocker)
	api := NewPublicTransactionPoolAPI(b, nonceLock)

	api2 := NewPrivateAccountAPI(b, nonceLock)
	api2.am = b.AM

	addr, _ := api2.NewAccount("1234")
	key, _ := api2.ImportRawKey("11223344556677889900aabbccddff0011223344556677889900aabbccddff00", "1234")
	d := uint64(1)
	_, _ = api2.UnlockAccount(ctx, key, "1234", &d)

	assert.NotPanics(t, func() {
		gas := hexutil.Uint64(0)
		gasPrice := hexutil.Big(*big.NewInt(0))
		nonce := hexutil.Uint64(1)
		res, err := api.SendTransaction(ctx, SendTxArgs{
			From:     addr,
			To:       &common.Address{2},
			Gas:      &gas,
			GasPrice: &gasPrice,
			Nonce:    &nonce,
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
	assert.NotPanics(t, func() {
		gas := hexutil.Uint64(0)
		gasPrice := hexutil.Big(*big.NewInt(0))
		res, err := api.SendTransaction(ctx, SendTxArgs{
			From:     addr,
			To:       &common.Address{2},
			Gas:      &gas,
			GasPrice: &gasPrice,
			Nonce:    nil,
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
func TestPublicTransactionPoolAPI_SendRawTransaction(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	nonceLock := new(AddrLocker)
	api := NewPublicTransactionPoolAPI(b, nonceLock)
	assert.NotPanics(t, func() {
		trx := types.NewTransaction(1, common.Address{1}, big.NewInt(1), 1, big.NewInt(0), []byte{})
		data, _ := rlp.EncodeToBytes(trx)
		res, err := api.SendRawTransaction(ctx, data)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
	assert.NotPanics(t, func() {
		data := hexutil.Bytes([]byte{1,2,3})
		res, err := api.SendRawTransaction(ctx, data)
		assert.Error(t, err)
		assert.NotEmpty(t, res)
	})
}
func TestPublicTransactionPoolAPI_Sign(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	nonceLock := new(AddrLocker)
	apiAM := NewPrivateAccountAPI(b, nonceLock)
	apiAM.am = b.AM

	addr, _ := apiAM.NewAccount("1234")
	key, _ := apiAM.ImportRawKey("11223344556677889900aabbccddff0011223344556677889900aabbccddff00", "1234")
	d := uint64(1)
	_, _ = apiAM.UnlockAccount(ctx, key, "1234", &d)

	api := NewPublicTransactionPoolAPI(b, nonceLock)
	assert.NotPanics(t, func() {
		trx := types.NewTransaction(1, common.Address{1}, big.NewInt(1), 1, big.NewInt(0), []byte{})
		data, _ := rlp.EncodeToBytes(trx)
		res, err := api.Sign(addr, data)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
