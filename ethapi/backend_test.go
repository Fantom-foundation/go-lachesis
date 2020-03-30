package ethapi

import (
	"context"
	"math/big"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	notify "github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"

	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/pos"
	"github.com/Fantom-foundation/go-lachesis/inter/sfctype"
)

func method() string {
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		parts := strings.Split(details.Name(), ".")
		name := parts[len(parts)-1]
		// fmt.Printf("Method name = '%s'\n", name)
		return name
	}
	return ""
}

func TestMethod(t *testing.T) {
	assert.Equal(t, method(), "TestMethod")
}

type testBackend struct {
	result struct {
		returned map[string][]interface{}
		err      map[string]error
		panic    map[string]string
	}
}

func NewTestBackend() *testBackend {
	return &testBackend{
		result: struct {
			returned map[string][]interface{}
			err      map[string]error
			panic    map[string]string
		}{
			returned: make(map[string][]interface{}),
			err: make(map[string]error),
			panic: make(map[string]string),
		},
	}
}

func (b *testBackend) Returned(method string, args ...interface{}) {
	b.result.err[method] = nil
	b.result.panic[method] = ""
	b.result.returned[method] = make([]interface{}, len(args), len(args))
	for i, v := range args {
		b.result.returned[method][i] = v
	}
}

func (b *testBackend) Error(method string, err error) {
	b.result.err[method] = err
}

func (b *testBackend) Panic(method string, msg string) {
	b.result.panic[method] = msg
}


func (b *testBackend) checkPanic(method string) {
	if b.result.panic[method] != "" {
		panic(b.result.panic[method])
	}
}





func (b *testBackend) ProtocolVersion() int {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(int)
}
func (b *testBackend) Progress() PeerProgress {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(PeerProgress)
}
func (b *testBackend) SuggestPrice(ctx context.Context) (*big.Int, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(*big.Int), b.result.err[method]
}
func (b *testBackend) ChainDb() ethdb.Database {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(ethdb.Database)
}
func (b *testBackend) AccountManager() *accounts.Manager {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(*accounts.Manager)
}
func (b *testBackend) ExtRPCEnabled() bool {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(bool)
}
func (b *testBackend) RPCGasCap() *big.Int {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(*big.Int)
} // global gas cap for eth_call over rpc: DoS protection

// Blockchain API
func (b *testBackend) HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*evmcore.EvmHeader, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(*evmcore.EvmHeader), b.result.err[method]
}
func (b *testBackend) HeaderByHash(ctx context.Context, hash common.Hash) (*evmcore.EvmHeader, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(*evmcore.EvmHeader), b.result.err[method]
}
func (b *testBackend) BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*evmcore.EvmBlock, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(*evmcore.EvmBlock), b.result.err[method]
}
func (b *testBackend) StateAndHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*state.StateDB, *evmcore.EvmHeader, error) {
	method := method()
	b.checkPanic(method)
	return 	b.result.returned[method][0].(*state.StateDB),
			b.result.returned[method][1].(*evmcore.EvmHeader),
			b.result.err[method]
}

//GetHeader(ctx context.Context, hash common.Hash) *evmcore.EvmHeader
func (b *testBackend) GetBlock(ctx context.Context, hash common.Hash) (*evmcore.EvmBlock, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(*evmcore.EvmBlock), b.result.err[method]
}
func (b *testBackend) GetReceiptsByNumber(ctx context.Context, number rpc.BlockNumber) (types.Receipts, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(types.Receipts), b.result.err[method]
}
func (b *testBackend) GetTd(hash common.Hash) *big.Int {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(*big.Int)
}
func (b *testBackend) GetEVM(ctx context.Context, msg evmcore.Message, state *state.StateDB, header *evmcore.EvmHeader) (*vm.EVM, func() error, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(*vm.EVM), b.result.returned[method][1].(func() error), b.result.err[method]
}

// Transaction pool API
func (b *testBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	method := method()
	b.checkPanic(method)
	return b.result.err[method]
}
func (b *testBackend) GetTransaction(ctx context.Context, txHash common.Hash) (*types.Transaction, uint64, uint64, error) {
	method := method()
	b.checkPanic(method)
	return 	b.result.returned[method][0].(*types.Transaction),
			b.result.returned[method][1].(uint64),
			b.result.returned[method][2].(uint64),
			b.result.err[method]
}
func (b *testBackend) GetPoolTransactions() (types.Transactions, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(types.Transactions), b.result.err[method]
}
func (b *testBackend) GetPoolTransaction(txHash common.Hash) *types.Transaction {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(*types.Transaction)
}
func (b *testBackend) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(uint64), b.result.err[method]
}
func (b *testBackend) Stats() (pending int, queued int) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(int), b.result.returned[method][1].(int)
}
func (b *testBackend) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	method := method()
	b.checkPanic(method)
	return 	b.result.returned[method][0].(map[common.Address]types.Transactions),
			b.result.returned[method][1].(map[common.Address]types.Transactions)
}
func (b *testBackend) SubscribeNewTxsNotify(chan<- evmcore.NewTxsNotify) notify.Subscription {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(notify.Subscription)
}

func (b *testBackend) ChainConfig() *params.ChainConfig {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(*params.ChainConfig)
}
func (b *testBackend) CurrentBlock() *evmcore.EvmBlock {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(*evmcore.EvmBlock)
}

// Lachesis DAG API
func (b *testBackend) GetEvent(ctx context.Context, shortEventID string) (*inter.Event, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(*inter.Event), b.result.err[method]
}
func (b *testBackend) GetEventHeader(ctx context.Context, shortEventID string) (*inter.EventHeaderData, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(*inter.EventHeaderData), b.result.err[method]
}
func (b *testBackend) GetConsensusTime(ctx context.Context, shortEventID string) (inter.Timestamp, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(inter.Timestamp), b.result.err[method]
}
func (b *testBackend) GetHeads(ctx context.Context, epoch rpc.BlockNumber) (hash.Events, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(hash.Events), b.result.err[method]
}
func (b *testBackend) CurrentEpoch(ctx context.Context) idx.Epoch {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(idx.Epoch)
}
func (b *testBackend) GetEpochStats(ctx context.Context, requestedEpoch rpc.BlockNumber) (*sfctype.EpochStats, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(*sfctype.EpochStats), b.result.err[method]
}
func (b *testBackend) TtfReport(ctx context.Context, untilBlock rpc.BlockNumber, maxBlocks idx.Block, mode string) (map[hash.Event]time.Duration, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(map[hash.Event]time.Duration), b.result.err[method]
}
func (b *testBackend) ForEachEvent(ctx context.Context, epoch rpc.BlockNumber, onEvent func(event *inter.Event) bool) error {
	method := method()
	b.checkPanic(method)
	return b.result.err[method]
}
func (b *testBackend) ValidatorTimeDrifts(ctx context.Context, epoch rpc.BlockNumber, maxEvents idx.Event) (map[idx.StakerID]map[hash.Event]time.Duration, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(map[idx.StakerID]map[hash.Event]time.Duration), b.result.err[method]
}

// Lachesis SFC API
func (b *testBackend) GetValidators(ctx context.Context) *pos.Validators {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(*pos.Validators)
}
func (b *testBackend) GetValidationScore(ctx context.Context, stakerID idx.StakerID) (*big.Int, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(*big.Int), b.result.err[method]
}
func (b *testBackend) GetOriginationScore(ctx context.Context, stakerID idx.StakerID) (*big.Int, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(*big.Int), b.result.err[method]
}
func (b *testBackend) GetRewardWeights(ctx context.Context, stakerID idx.StakerID) (*big.Int, *big.Int, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(*big.Int), b.result.returned[method][1].(*big.Int), b.result.err[method]
}
func (b *testBackend) GetStakerPoI(ctx context.Context, stakerID idx.StakerID) (*big.Int, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(*big.Int), b.result.err[method]
}
func (b *testBackend) GetDowntime(ctx context.Context, stakerID idx.StakerID) (idx.Block, inter.Timestamp, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(idx.Block), b.result.returned[method][0].(inter.Timestamp), b.result.err[method]
}
func (b *testBackend) GetDelegatorClaimedRewards(ctx context.Context, addr common.Address) (*big.Int, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(*big.Int), b.result.err[method]
}
func (b *testBackend) GetStakerClaimedRewards(ctx context.Context, stakerID idx.StakerID) (*big.Int, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(*big.Int), b.result.err[method]
}
func (b *testBackend) GetStakerDelegatorsClaimedRewards(ctx context.Context, stakerID idx.StakerID) (*big.Int, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(*big.Int), b.result.err[method]
}
func (b *testBackend) GetStaker(ctx context.Context, stakerID idx.StakerID) (*sfctype.SfcStaker, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(*sfctype.SfcStaker), b.result.err[method]
}
func (b *testBackend) GetStakerID(ctx context.Context, addr common.Address) (idx.StakerID, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(idx.StakerID), b.result.err[method]
}
func (b *testBackend) GetStakers(ctx context.Context) ([]sfctype.SfcStakerAndID, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].([]sfctype.SfcStakerAndID), b.result.err[method]
}
func (b *testBackend) GetDelegatorsOf(ctx context.Context, stakerID idx.StakerID) ([]sfctype.SfcDelegatorAndAddr, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].([]sfctype.SfcDelegatorAndAddr), b.result.err[method]
}
func (b *testBackend) GetDelegator(ctx context.Context, addr common.Address) (*sfctype.SfcDelegator, error) {
	method := method()
	b.checkPanic(method)
	return b.result.returned[method][0].(*sfctype.SfcDelegator), b.result.err[method]
}
