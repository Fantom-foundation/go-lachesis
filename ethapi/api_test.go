package ethapi

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"

	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/kvdb/nokeyiserr"
	"github.com/Fantom-foundation/go-lachesis/kvdb/table"
)

func SetBackendStateDB(b *testBackend) *state.StateDB {
	db1 := rawdb.NewDatabase(
		nokeyiserr.Wrap(
			table.New(
				memorydb.New(), []byte("evm1_"))))
	stateDB, _ := state.New(common.HexToHash("0x0"), state.NewDatabase(db1))
	b.Returned("StateAndHeaderByNumber", stateDB, &evmcore.EvmHeader{})
	vmCtx := vm.Context{}
	evm := vm.NewEVM(vmCtx, stateDB, &params.ChainConfig{}, vm.Config{})
	b.Returned("GetEVM", evm, func()error{return nil})
	b.Returned("AccountManager", &accounts.Manager{})

	return stateDB
}

// TODO: for no error
func TestDoCall(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()
	SetBackendStateDB(b)

	nonce := hexutil.Uint64(1)
	assert.NotPanics(t, func() {
		gas := hexutil.Uint64(1)
		DoCall(ctx, b, CallArgs{Gas: &gas, }, rpc.BlockNumber(1), map[common.Address]account{
			common.HexToAddress("0x0"): account{Nonce: &nonce},
		}, vm.Config{}, 100*time.Second, big.NewInt(100000))
		// assert.NoError(t, err)
	})
}
// TODO: for no error
func TestDoEstimateGas(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()
	SetBackendStateDB(b)

	assert.NotPanics(t, func() {
		gas := hexutil.Uint64(21000)
		DoEstimateGas(ctx, b, CallArgs{Gas: &gas, }, rpc.BlockNumber(1), big.NewInt(22000))
		// assert.NoError(t, err)
	})
}
func TestFormatLogs(t *testing.T) {
	assert.NotPanics(t, func() {
		stack := make([]*big.Int, 0, 1024)
		stack = append(stack, big.NewInt(1))
		mem := make([]byte, 32, 32)
		storage := make(map[common.Hash]common.Hash, 3)
		log := vm.StructLog{
			Pc:            1,
			Op:            2,
			Gas:           3,
			GasCost:       4,
			Memory:        mem,
			MemorySize:    5,
			Stack:         stack,
			Storage:       storage,
			Depth:         6,
			RefundCounter: 7,
			Err:           nil,
		}
		res := FormatLogs([]vm.StructLog{log})
		assert.NotEmpty(t, res)
	})
}
func TestRPCMarshalBlock(t *testing.T) {
	assert.NotPanics(t, func() {
		res, err := RPCMarshalBlock(&evmcore.EvmBlock{
			EvmHeader:    evmcore.EvmHeader{
				Number:     big.NewInt(1),
				Hash:       common.Hash{2},
				ParentHash: common.Hash{1},
				Root:       common.Hash{0},
				TxHash:     common.Hash{1},
				Time:       1,
				Coinbase:   common.Address{0},
				GasLimit:   0,
				GasUsed:    0,
			},
			Transactions: types.Transactions{
				types.NewTransaction(1, common.Address{1}, big.NewInt(1), 1, big.NewInt(0), []byte{}),
				types.NewTransaction(2, common.Address{2}, big.NewInt(2), 2, big.NewInt(0), []byte{}),
			},
		}, true, true)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
func TestRPCMarshalEvent(t *testing.T) {
	assert.NotPanics(t, func() {
		res, err := RPCMarshalEvent(
			&inter.Event{
				EventHeader:  inter.EventHeader{
					EventHeaderData: inter.EventHeaderData{
						Version:       0,
						Epoch:         0,
						Seq:           0,
						Frame:         0,
						IsRoot:        false,
						Creator:       0,
						PrevEpochHash: common.Hash{},
						Parents:       nil,
						GasPowerLeft:  inter.GasPowerLeft{},
						GasPowerUsed:  0,
						Lamport:       0,
						ClaimedTime:   0,
						MedianTime:    0,
						TxHash:        common.Hash{},
						Extra:         nil,
					},
					Sig:             nil,
				},
				Transactions: types.Transactions{
					types.NewTransaction(1, common.Address{1}, big.NewInt(1), 1, big.NewInt(0), []byte{}),
					types.NewTransaction(2, common.Address{2}, big.NewInt(2), 2, big.NewInt(0), []byte{}),
				},
			},
			true, true,
		)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
func TestRPCMarshalEventHeader(t *testing.T) {
	assert.NotPanics(t, func() {
		res := RPCMarshalEventHeader(
			&inter.EventHeaderData{
				Version:       0,
				Epoch:         0,
				Seq:           0,
				Frame:         0,
				IsRoot:        false,
				Creator:       0,
				PrevEpochHash: common.Hash{},
				Parents:       nil,
				GasPowerLeft:  inter.GasPowerLeft{},
				GasPowerUsed:  0,
				Lamport:       0,
				ClaimedTime:   0,
				MedianTime:    0,
				TxHash:        common.Hash{},
				Extra:         nil,
			},
		)
		assert.NotEmpty(t, res)
	})
}
func TestRPCMarshalHeader(t *testing.T) {
	assert.NotPanics(t, func() {
		res := RPCMarshalHeader(
			&evmcore.EvmHeader{
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
		)
		assert.NotEmpty(t, res)
	})
}
// TODO
func TestSubmitTransaction(t *testing.T) {

}
