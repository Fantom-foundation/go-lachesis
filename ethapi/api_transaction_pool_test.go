package ethapi

import (
	"context"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

// PublicTransactionPoolAPI

func TestPublicTransactionPoolAPI_FillTransaction(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()
	SetBackendStateDB(b)
	b.Returned("SuggestPrice", big.NewInt(1))
	b.Returned("GetPoolNonce", uint64(1))
	b.Returned("RPCGasCap", big.NewInt(1))

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

}
func TestPublicTransactionPoolAPI_GetBlockTransactionCountByNumber(t *testing.T) {

}
func TestPublicTransactionPoolAPI_GetRawTransactionByBlockHashAndIndex(t *testing.T) {

}
func TestPublicTransactionPoolAPI_GetRawTransactionByBlockNumberAndIndex(t *testing.T) {

}
func TestPublicTransactionPoolAPI_GetRawTransactionByHash(t *testing.T) {

}
func TestPublicTransactionPoolAPI_GetTransactionByBlockHashAndIndex(t *testing.T) {

}
func TestPublicTransactionPoolAPI_GetTransactionByBlockNumberAndIndex(t *testing.T) {

}
func TestPublicTransactionPoolAPI_GetTransactionByHash(t *testing.T) {

}
func TestPublicTransactionPoolAPI_GetTransactionCount(t *testing.T) {

}
func TestPublicTransactionPoolAPI_GetTransactionReceipt(t *testing.T) {

}
func TestPublicTransactionPoolAPI_PendingTransactions(t *testing.T) {

}
func TestPublicTransactionPoolAPI_Resend(t *testing.T) {

}
func TestPublicTransactionPoolAPI_SendRawTransaction(t *testing.T) {

}
func TestPublicTransactionPoolAPI_SendTransaction(t *testing.T) {

}
func TestPublicTransactionPoolAPI_Sign(t *testing.T) {

}
func TestPublicTransactionPoolAPI_SignTransaction(t *testing.T) {

}

// SendTxArgs

func TestSendTxArgs(t *testing.T) {

}
