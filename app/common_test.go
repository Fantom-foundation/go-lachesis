package app

//go:generate mkdir -p solc
// NOTE: assumed that SFC-repo is in the same dir than lachesis-repo
// 1.0.0 (genesis)
//go:generate bash -c "cd ../../fantom-sfc && git checkout 1.0.0 && docker run --rm -v $(pwd):/src -v $(pwd)/../go-lachesis/app:/dst ethereum/solc:0.5.12 -o /dst/solc/ --optimize --optimize-runs=2000 --bin --abi --allow-paths /src/contracts --overwrite /src/contracts/upgradeability/UpgradeabilityProxy.sol"
//go:generate mkdir -p sfcproxy
//go:generate abigen --bin=./solc/UpgradeabilityProxy.bin --abi=./solc/UpgradeabilityProxy.abi --pkg=sfcproxy --type=Contract --out=sfcproxy/contract.go
// 1.1.0-rc1
//go:generate bash -c "cd ../../fantom-sfc && git checkout 1.1.0-rc1 && docker run --rm -v $(pwd):/src -v $(pwd)/../go-lachesis/app:/dst ethereum/solc:0.5.12 -o /dst/solc/ --optimize --optimize-runs=2000 --bin --abi --allow-paths /src/contracts --overwrite /src/contracts/sfc/Staker.sol"
//go:generate mkdir -p sfc110
//go:generate abigen --bin=./solc/Stakers.bin --abi=./solc/Stakers.abi --pkg=sfc110 --type=Contract --out=sfc110/contract.go
// 2.0.1-rc2
//go:generate bash -c "cd ../../fantom-sfc && git checkout 2.0.1-rc2 && docker run --rm -v $(pwd):/src -v $(pwd)/../go-lachesis/app:/dst ethereum/solc:0.5.12 -o /dst/solc/ --optimize --optimize-runs=2000 --bin --abi --allow-paths /src/contracts --overwrite /src/contracts/sfc/Staker.sol"
//go:generate mkdir -p sfc201
//go:generate abigen --bin=./solc/Stakers.bin --abi=./solc/Stakers.abi --pkg=sfc201 --type=Contract --out=sfc201/contract.go
// clean
//go:generate rm -fr ./solc

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/state"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/Fantom-foundation/go-lachesis/crypto"
	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/lachesis"
	"github.com/Fantom-foundation/go-lachesis/lachesis/genesis"
	"github.com/Fantom-foundation/go-lachesis/lachesis/params"
	"github.com/Fantom-foundation/go-lachesis/utils"
)

const (
	gasLimit       = uint64(21000)
	genesisStakers = 3
	genesisBalance = 1e18
	genesisStake   = 2 * 4e6
	epochDuration  = time.Hour

	sameEpoch = time.Hour / 1000
	nextEpoch = time.Hour
)

type testEnv struct {
	App   *App
	Store *Store

	GasPrice *big.Int

	signer eth.Signer

	lastBlock     idx.Block
	lastBlockTime time.Time
	lastState     common.Hash
	validators    []idx.StakerID
	delegators    []common.Address

	nonces map[common.Address]uint64
}

func newTestEnv() *testEnv {
	vaccs := genesis.FakeAccounts(1, genesisStakers, utils.ToFtm(genesisBalance), utils.ToFtm(genesisStake))
	cfg := Config{
		Net: lachesis.FakeNetConfig(vaccs),
	}
	cfg.Net.Dag.MaxEpochDuration = epochDuration

	s := NewMemStore()
	_, _, err := s.ApplyGenesis(&cfg.Net)
	if err != nil {
		panic(err)
	}

	a := New(cfg, s)
	_ = a.InitChain(cfg.Net.ChainInfo())

	validators := make([]idx.StakerID, len(vaccs.Validators))
	for i := range validators {
		validators[i] = idx.StakerID(i + 1)
	}

	return &testEnv{
		App:   a,
		Store: s,

		GasPrice: params.MinGasPrice,

		signer: eth.NewEIP155Signer(big.NewInt(int64(cfg.Net.NetworkID))),

		lastBlock:     0,
		lastBlockTime: cfg.Net.Genesis.Time.Time(),
		lastState:     s.GetBlock(0).Root,
		validators:    validators,

		nonces: make(map[common.Address]uint64),
	}
}

func (e *testEnv) Close() {
	e.App.Close()
	e.Store.Close()
}

func (env *testEnv) AddValidator(v idx.StakerID) {
	for _, already := range env.validators {
		if v == already {
			return
		}
	}
	env.validators = append(env.validators, v)
}

func (env *testEnv) DelValidator(v idx.StakerID) {
	for i := 0; i < len(env.validators); i++ {
		if env.validators[i] != v {
			continue
		}
		env.validators = append(env.validators[:i], env.validators[i+1:]...)
		return
	}
}

func (env *testEnv) AddDelegator(v common.Address) {
	for _, already := range env.delegators {
		if v == already {
			return
		}
	}
	env.delegators = append(env.delegators, v)
}

func (env *testEnv) DelDelegator(v common.Address) {
	for i := 0; i < len(env.delegators); i++ {
		if env.delegators[i] != v {
			continue
		}
		env.delegators = append(env.delegators[:i], env.delegators[i+1:]...)
		return
	}
}

func (env *testEnv) ApplyBlock(spent time.Duration, txs ...*eth.Transaction) eth.Receipts {
	env.lastBlock++
	env.lastBlockTime = env.lastBlockTime.Add(spent)

	evmHeader := evmcore.EvmHeader{
		Number:   big.NewInt(int64(env.lastBlock)),
		Root:     env.lastState,
		Time:     inter.TimeToStamp(env.lastBlockTime),
		GasLimit: math.MaxUint64,
	}

	blockParticipated := make(map[idx.StakerID]bool)
	for _, p := range env.validators {
		blockParticipated[p] = true
	}

	env.App.beginBlock(evmHeader, env.lastState, inter.Cheaters{}, blockParticipated)

	receipts := make(eth.Receipts, len(txs))
	for i, tx := range txs {
		originator := env.validators[i%len(env.validators)]
		receipt, err := env.App.deliverTx(tx, originator)
		if err != nil {
			panic(err)
		}
		if receipt.Status != eth.ReceiptStatusSuccessful {
			panic("tx failed")
		}
		receipts[i] = receipt

		evmHeader.GasUsed += uint64(receipt.GasUsed)

		sender, err := eth.Sender(env.signer, tx)
		if err != nil {
			panic(err)
		}
		env.incNonce(sender)
	}

	env.App.endBlock(env.lastBlock)

	env.lastState = env.App.commit()

	return receipts
}

func (env *testEnv) Transfer(from int, to int, amount *big.Int) *eth.Transaction {
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

func (env *testEnv) Contract(from int, amount *big.Int, hex string) *eth.Transaction {
	data := hexutil.MustDecode(hex)
	nonce, _ := env.PendingNonceAt(nil, env.Address(from))
	key := env.privateKey(from)
	tx := eth.NewContractCreation(nonce, amount, gasLimit*10000, env.GasPrice, data)
	tx, err := eth.SignTx(tx, env.signer, key)
	if err != nil {
		panic(err)
	}

	return tx
}

func (env *testEnv) privateKey(n int) *ecdsa.PrivateKey {
	key := crypto.FakeKey(n)
	return key
}

func (env *testEnv) Address(n int) common.Address {
	key := crypto.FakeKey(n)
	addr := crypto.PubkeyToAddress(key.PublicKey)
	return addr
}

func (env *testEnv) Payer(n int, amounts ...*big.Int) *bind.TransactOpts {
	key := env.privateKey(n)
	t := bind.NewKeyedTransactor(key)
	nonce, _ := env.PendingNonceAt(nil, env.Address(n))
	t.Nonce = big.NewInt(int64(nonce))
	t.Value = big.NewInt(0)
	for _, amount := range amounts {
		t.Value.Add(t.Value, amount)
	}
	return t
}

func (env *testEnv) ReadOnly() *bind.CallOpts {
	return &bind.CallOpts{}
}

func (env *testEnv) State() *state.StateDB {
	return env.Store.StateDB(env.lastState)
}

func (env *testEnv) incNonce(account common.Address) {
	env.nonces[account] += 1
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
	code := env.State().GetCode(account)
	return code, nil
}

// PendingNonceAt retrieves the current pending nonce associated with an account.
func (env *testEnv) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	nonce := env.nonces[account]
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
	if call.To == nil {
		gas = gasLimit * 10000
	} else {
		gas = gasLimit * 10
	}
	return
}

// SendTransaction injects the transaction into the pending pool for execution.
func (env *testEnv) SendTransaction(ctx context.Context, tx *eth.Transaction) error {
	// do nothing to avoid executing by transactor, only generating needed
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
