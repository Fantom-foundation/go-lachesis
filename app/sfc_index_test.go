package app

import (
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"

	// "github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/abci/types"

	"github.com/Fantom-foundation/go-lachesis/crypto"
	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/lachesis"
	"github.com/Fantom-foundation/go-lachesis/lachesis/genesis"
	"github.com/Fantom-foundation/go-lachesis/utils"
)

func Test(t *testing.T) {
	//require := require.New(t)

	accs := []common.Address{
		crypto.PubkeyToAddress(crypto.FakeKey(1).PublicKey),
		crypto.PubkeyToAddress(crypto.FakeKey(2).PublicKey),
		crypto.PubkeyToAddress(crypto.FakeKey(3).PublicKey),
	}

	a, s := initTestApp(len(accs))
	defer s.Close()
	defer a.Close()

	gasPrice := utils.ToFtm(1e8)
	gasLimit := uint64(1e18)
	state := s.GetHeader(common.Hash{}, 0).Root

	txs := eth.Transactions{
		eth.NewTransaction(0, accs[1], utils.ToFtm(1e19), gasLimit, gasPrice, nil),
	}

	state = applyTxs(a, 1, state, txs)
	// TODO: check SFC
}

func applyTxs(a *App, n idx.Block, state common.Hash, txs eth.Transactions) common.Hash {
	prev := a.store.GetHeader(common.Hash{}, uint64(n-1))

	block := &evmcore.EvmHeader{
		Hash:       common.Hash(hash.FakeEvent()),
		ParentHash: prev.ParentHash,
		Root:       state,
		Number:     big.NewInt(int64(n)),
		Time:       inter.TimeToStamp(time.Now()),
		GasLimit:   math.MaxUint64,
	}

	applyBlock(a, state, block, map[idx.StakerID]bool{
		1: true,
		2: true,
	})

	for _, t := range txs {
		tx := &DagTx{
			Transaction: t,
			Originator:  idx.StakerID(1),
		}
		a.DeliverTx(types.RequestDeliverTx{
			Tx: tx.Bytes(),
		})
	}

	a.EndBlock(types.RequestEndBlock{
		Height: int64(n),
	})

	state = common.BytesToHash(a.Commit().Data)
	return state
}

func applyBlock(a *App,
	state common.Hash,
	block *evmcore.EvmHeader,
	participated map[idx.StakerID]bool,
) {

	votes := make([]types.VoteInfo, 0, len(participated))
	for staker, ok := range participated {
		if !ok {
			continue
		}
		votes = append(votes, types.VoteInfo{
			Validator: types.Validator{
				Address: staker.Bytes(),
			},
		})
	}

	req := types.RequestBeginBlock{
		Hash: block.Hash.Bytes(),
		Header: types.Header{
			Height:        block.Number.Int64(),
			Time:          block.Time.Time(),
			ConsensusHash: block.Hash.Bytes(),
			LastBlockId: types.BlockID{
				Hash: block.ParentHash.Bytes(),
			},
			LastCommitHash: state.Bytes(),
		},
		LastCommitInfo: types.LastCommitInfo{
			Votes: votes,
		},
		ByzantineValidators: make([]types.Evidence, 0),
	}

	_ = a.BeginBlock(req)
}

func initTestApp(validators int) (*App, *Store) {
	s := NewMemStore()

	vaccs := genesis.FakeAccounts(1, validators, utils.ToFtm(1e10), utils.ToFtm(3175000))
	cfg := Config{
		Net: lachesis.FakeNetConfig(vaccs),
	}
	cfg.Net.Dag.MaxEpochBlocks = 3
	_, _, err := s.ApplyGenesis(&cfg.Net)
	if err != nil {
		panic(err)
	}

	a := New(cfg, s)
	_ = a.InitChain(
		cfg.Net.ChainInfo())

	return a, s
}
