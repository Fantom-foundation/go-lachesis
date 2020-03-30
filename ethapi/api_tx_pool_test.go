package ethapi

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

// PublicTxPoolAPI

func TestPublicTxPoolAPI_Content(t *testing.T) {
	b := NewTestBackend()
	SetBackendStateDB(b)
	b.Returned("TxPoolContent",
		map[common.Address]types.Transactions{
			common.Address{1}: types.Transactions{
				types.NewTransaction(1, common.Address{1}, big.NewInt(1), 1, big.NewInt(0), []byte{}),
				types.NewTransaction(2, common.Address{2}, big.NewInt(2), 2, big.NewInt(0), []byte{}),
			},
		},
		map[common.Address]types.Transactions{
			common.Address{1}: types.Transactions{
				types.NewTransaction(3, common.Address{3}, big.NewInt(3), 3, big.NewInt(0), []byte{}),
				types.NewTransaction(4, common.Address{4}, big.NewInt(4), 4, big.NewInt(0), []byte{}),
			},
		},
	)

	api := NewPublicTxPoolAPI(b)
	assert.NotPanics(t, func() {
		res := api.Content()
		assert.NotEmpty(t, res)
	})
}
func TestPublicTxPoolAPI_Status(t *testing.T) {
	b := NewTestBackend()
	SetBackendStateDB(b)
	b.Returned("Stats", 1, 1)

	api := NewPublicTxPoolAPI(b)
	assert.NotPanics(t, func() {
		res := api.Status()
		assert.NotEmpty(t, res)
	})
}
func TestPublicTxPoolAPI_Inspect(t *testing.T) {
	b := NewTestBackend()
	SetBackendStateDB(b)
	b.Returned("Stats", 2, 2)
	b.Returned("TxPoolContent",
		map[common.Address]types.Transactions{
			common.Address{1}: types.Transactions{
				types.NewTransaction(1, common.Address{1}, big.NewInt(1), 1, big.NewInt(0), []byte{}),
				types.NewTransaction(2, common.Address{2}, big.NewInt(2), 2, big.NewInt(0), []byte{}),
			},
		},
		map[common.Address]types.Transactions{
			common.Address{1}: types.Transactions{
				types.NewTransaction(3, common.Address{3}, big.NewInt(3), 3, big.NewInt(0), []byte{}),
				types.NewTransaction(4, common.Address{4}, big.NewInt(4), 4, big.NewInt(0), []byte{}),
			},
		},
	)

	api := NewPublicTxPoolAPI(b)
	assert.NotPanics(t, func() {
		res := api.Inspect()
		assert.NotEmpty(t, res)
	})
}
