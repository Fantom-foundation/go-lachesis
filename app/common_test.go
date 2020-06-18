package app

/*
docker run --rm -v $PWD:/tmp  ethereum/solc:0.5.12 -o /tmp/build --optimize --optimize-runs=2000 --bin --abi --bin-runtime --allow-paths /tmp/contracts --overwrite /tmp/contracts/sfc/Staker.sol
go run go-ethereum/cmd/abigen --bin=./build/Stakers.bin --abi=./build/Stakers.abi --pkg=contract --out=contracts/contract.go
*/

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/lachesis"
	"github.com/Fantom-foundation/go-lachesis/lachesis/genesis"
	"github.com/Fantom-foundation/go-lachesis/lachesis/params"
	"github.com/Fantom-foundation/go-lachesis/utils"
)

type testEnv struct {
	App   *App
	Store *Store

	GasPrice *big.Int

	vaccs  genesis.VAccounts
	signer eth.Signer

	lastBlock   idx.Block
	lastState   common.Hash
	originators []idx.StakerID

	nonces map[common.Address]uint64
}

func newTestEnv(validators int) *testEnv {
	vaccs := genesis.FakeAccounts(1, validators, utils.ToFtm(1e18), utils.ToFtm(3175000))
	cfg := Config{
		Net: lachesis.FakeNetConfig(vaccs),
	}
	cfg.Net.Dag.MaxEpochBlocks = 1

	s := NewMemStore()
	_, _, err := s.ApplyGenesis(&cfg.Net)
	if err != nil {
		panic(err)
	}

	a := New(cfg, s)
	_ = a.InitChain(cfg.Net.ChainInfo())

	originators := make([]idx.StakerID, len(vaccs.Validators))
	for i := range originators {
		originators[i] = idx.StakerID(i)
	}

	return &testEnv{
		App:   a,
		Store: s,

		GasPrice: params.MinGasPrice,

		vaccs:  vaccs,
		signer: eth.NewEIP155Signer(big.NewInt(int64(cfg.Net.NetworkID))),

		lastBlock:   0,
		lastState:   s.GetBlock(0).Root,
		originators: originators,

		nonces: make(map[common.Address]uint64),
	}
}

func (e *testEnv) Close() {
	e.App.Close()
	e.Store.Close()
}

func (env *testEnv) ApplyBlock(txs ...*eth.Transaction) eth.Receipts {
	env.lastBlock++

	evmHeader := evmcore.EvmHeader{
		Number:   big.NewInt(int64(env.lastBlock)),
		Root:     env.lastState,
		Time:     inter.TimeToStamp(time.Now()),
		GasLimit: math.MaxUint64,
	}

	blockParticipated := make(map[idx.StakerID]bool)
	for _, p := range env.originators {
		blockParticipated[p] = true
	}

	env.App.beginBlock(evmHeader, env.lastState, inter.Cheaters{}, blockParticipated)

	receipts := make(eth.Receipts, len(txs))
	for i, tx := range txs {
		originator := env.originators[i%len(env.originators)]
		receipt, err := env.App.deliverTx(tx, originator)
		if err != nil {
			panic(err)
		}
		receipts[i] = receipt
		evmHeader.GasUsed += uint64(receipt.GasUsed)
	}

	env.App.endBlock(env.lastBlock)

	env.lastState = env.App.commit()

	return receipts
}

func (env *testEnv) Tx(from int, to int, amount *big.Int) *eth.Transaction {
	const gasLimit = uint64(21000)

	nonce, _ := env.PendingNonceAt(nil, env.Address(from))
	key := env.privateKey(from)
	receiver := env.Address(to)
	tx := eth.NewTransaction(nonce, receiver, amount, gasLimit, env.GasPrice, nil)
	tx, err := eth.SignTx(tx, env.signer, key)
	if err != nil {
		panic(err)
	}

	return tx
}

func (env *testEnv) Contract(from int, amount *big.Int, data []byte) *eth.Transaction {
	const gasLimit = uint64(250000000)

	nonce, _ := env.PendingNonceAt(nil, env.Address(from))
	key := env.privateKey(from)
	tx := eth.NewContractCreation(nonce, amount, gasLimit, env.GasPrice, data)
	tx, err := eth.SignTx(tx, env.signer, key)
	if err != nil {
		panic(err)
	}

	return tx
}

func (env *testEnv) privateKey(n int) *ecdsa.PrivateKey {
	acc := env.vaccs.Validators[n].Address
	key := env.vaccs.Accounts[acc].PrivateKey
	return key
}

func (env *testEnv) Staker(n int) idx.StakerID {
	return idx.StakerID(n + 1)
}

func (env *testEnv) Address(n int) common.Address {
	return env.vaccs.Validators[n].Address
}

func (env *testEnv) State() *state.StateDB {
	return env.Store.StateDB(env.lastState)
}

/*
 * bind.ContractCaller interface
 */
var (
	errBlockNumberUnsupported = errors.New("simulatedBackend cannot access blocks other than the latest block")
)

// CodeAt returns the code of the given account. This is needed to differentiate
// between contract internal errors and the local chain being out of sync.
func (env *testEnv) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	if blockNumber != nil && idx.Block(blockNumber.Uint64()) != env.lastBlock {
		return nil, errBlockNumberUnsupported
	}

	code := env.State().GetCode(contract)
	return code, nil
}

// ContractCall executes an Ethereum contract call with the specified data as the
// input.
func (env *testEnv) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	if blockNumber != nil && idx.Block(blockNumber.Uint64()) != env.lastBlock {
		return nil, errBlockNumberUnsupported
	}

	block := &evmcore.EvmBlock{
		EvmHeader: *env.App.BlockChain().GetHeader(common.Hash{}, uint64(env.lastBlock)),
	}

	rval, _, _, err := env.callContract(ctx, call, block, env.State())
	return rval, err
}

// callContract implements common code between normal and pending contract calls.
// state is modified during execution, make sure to copy it if necessary.
func (env *testEnv) callContract(
	ctx context.Context, call ethereum.CallMsg, block *evmcore.EvmBlock, statedb *state.StateDB,
) (
	ret []byte, usedGas uint64, failed bool, err error,
) {
	// Ensure message is initialized properly.
	if call.GasPrice == nil {
		call.GasPrice = big.NewInt(1)
	}
	if call.Gas == 0 {
		call.Gas = 50000000
	}
	if call.Value == nil {
		call.Value = new(big.Int)
	}
	// Set infinite balance to the fake caller account.
	from := statedb.GetOrNewStateObject(call.From)
	from.SetBalance(big.NewInt(math.MaxInt64))

	msg := callmsg{call}

	evmContext := evmcore.NewEVMContext(msg, block.Header(), env.App.BlockChain(), &call.From)
	// Create a new environment which holds all relevant information
	// about the transaction and calling mechanisms.
	vmenv := vm.NewEVM(evmContext, statedb, env.App.config.Net.EvmChainConfig(), vm.Config{})
	gaspool := new(evmcore.GasPool).AddGas(math.MaxUint64)
	ret, usedGas, _, failed, err = evmcore.NewStateTransition(vmenv, msg, gaspool).TransitionDb()
	return
}

/*
 * ContractTransactor interface
 */

// PendingCodeAt returns the code of the given account in the pending state.
func (env *testEnv) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	panic("not implemented yet")
	return nil, nil
}

// PendingNonceAt retrieves the current pending nonce associated with an account.
func (env *testEnv) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	nonce := env.nonces[account]
	env.nonces[account] = nonce + 1
	return nonce, nil
}

// SuggestGasPrice retrieves the currently suggested gas price to allow a timely
// execution of a transaction.
func (env *testEnv) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return env.GasPrice, nil
}

// EstimateGas tries to estimate the gas needed to execute a specific
// transaction based on the current pending state of the backend blockchain.
// There is no guarantee that this is the true gas limit requirement as other
// transactions may be added or removed by miners, but it should provide a basis
// for setting a reasonable default.
func (env *testEnv) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	return uint64(2500000), nil
}

// SendTransaction injects the transaction into the pending pool for execution.
func (env *testEnv) SendTransaction(ctx context.Context, tx *eth.Transaction) error {
	env.ApplyBlock(tx)
	return nil
}

/*
 *  bind.ContractFilterer interface
 */
// FilterLogs executes a log filter operation, blocking during execution and
// returning all the results in one batch.
func (env *testEnv) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]eth.Log, error) {
	panic("not implemented yet")
	return nil, nil
}

// SubscribeFilterLogs creates a background log filtering operation, returning
// a subscription immediately, which can be used to stream the found events.
func (env *testEnv) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- eth.Log) (ethereum.Subscription, error) {
	panic("not implemented yet")
	return nil, nil
}

// callmsg implements core.Message to allow passing it as a transaction simulator.
type callmsg struct {
	ethereum.CallMsg
}

func (m callmsg) From() common.Address { return m.CallMsg.From }
func (m callmsg) Nonce() uint64        { return 0 }
func (m callmsg) CheckNonce() bool     { return false }
func (m callmsg) To() *common.Address  { return m.CallMsg.To }
func (m callmsg) GasPrice() *big.Int   { return m.CallMsg.GasPrice }
func (m callmsg) Gas() uint64          { return m.CallMsg.Gas }
func (m callmsg) Value() *big.Int      { return m.CallMsg.Value }
func (m callmsg) Data() []byte         { return m.CallMsg.Data }
