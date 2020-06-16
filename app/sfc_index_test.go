package app

import (
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/lachesis"
	"github.com/Fantom-foundation/go-lachesis/lachesis/genesis"
	"github.com/Fantom-foundation/go-lachesis/lachesis/params"
	"github.com/Fantom-foundation/go-lachesis/utils"
)

func Test(t *testing.T) {
	env := newTestEnv(3)
	defer env.Close()

	t.Run("Transfer", func(t *testing.T) {
		require := require.New(t)

		b0 := utils.ToFtm(1e10)
		require.Equal(b0, env.State().GetBalance(env.Address(0)))
		require.Equal(b0, env.State().GetBalance(env.Address(1)))
		require.Equal(b0, env.State().GetBalance(env.Address(2)))

		env.ApplyBlock(
			env.Tx(0, env.Staker(0), env.Address(1), utils.ToFtm(100)),
		)
		env.ApplyBlock(
			env.Tx(0, env.Staker(1), env.Address(2), utils.ToFtm(100)),
		)
		env.ApplyBlock(
			env.Tx(0, env.Staker(2), env.Address(0), utils.ToFtm(100)),
		)

		gas := big.NewInt(0).Mul(big.NewInt(21000), env.GasPrice)
		b1 := big.NewInt(0).Sub(b0, gas)
		require.Equal(b1, env.State().GetBalance(env.Address(0)))
		require.Equal(b1, env.State().GetBalance(env.Address(1)))
		require.Equal(b1, env.State().GetBalance(env.Address(2)))
	})

	t.Run("SFC deploy", func(t *testing.T) {
		require := require.New(t)

		r := env.ApplyBlock(
			env.Contract(1, env.Staker(0), utils.ToFtm(0), MainContractBinV2),
		)

		t.Logf("%+v", r[0])

		require.NotNil(r[0])

	})

}

/*
 * test env:
 */

type testEnv struct {
	App   *App
	Store *Store

	GasPrice *big.Int

	vaccs  genesis.VAccounts
	signer eth.Signer

	lastBlock   idx.Block
	lastState   common.Hash
	originators []idx.StakerID
}

func newTestEnv(validators int) *testEnv {
	vaccs := genesis.FakeAccounts(1, validators, utils.ToFtm(1e10), utils.ToFtm(3175000))
	cfg := Config{
		Net: lachesis.FakeNetConfig(vaccs),
	}
	cfg.Net.Dag.MaxEpochBlocks = 3

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

func (env *testEnv) Tx(nonce uint64, from idx.StakerID, to common.Address, amount *big.Int) *eth.Transaction {
	var (
		gasLimit = uint64(21000)
	)

	sender := env.vaccs.Validators[int(from-1)].Address
	key := env.vaccs.Accounts[sender].PrivateKey

	tx := eth.NewTransaction(nonce, to, amount, gasLimit, env.GasPrice, nil)
	tx, err := eth.SignTx(tx, env.signer, key)
	if err != nil {
		panic(err)
	}

	return tx
}

func (env *testEnv) Contract(nonce uint64, from idx.StakerID, amount *big.Int, data []byte) *eth.Transaction {
	var (
		gasLimit = uint64(2500000)
	)

	sender := env.vaccs.Validators[int(from-1)].Address
	key := env.vaccs.Accounts[sender].PrivateKey

	tx := eth.NewContractCreation(nonce, amount, gasLimit, env.GasPrice, data)
	tx, err := eth.SignTx(tx, env.signer, key)
	if err != nil {
		panic(err)
	}

	return tx
}

func (env *testEnv) Staker(n int) idx.StakerID {
	return idx.StakerID(n + 1)
}

func (env *testEnv) Address(n int) common.Address {
	// 	to := crypto.PubkeyToAddress(crypto.FakeKey(0).PublicKey)
	return env.vaccs.Validators[n].Address
}

func (env *testEnv) State() *state.StateDB {
	return env.Store.StateDB(env.lastState)
}
