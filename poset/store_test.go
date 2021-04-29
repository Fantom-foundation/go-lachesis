package poset

import (
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"

	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/pos"
	"github.com/Fantom-foundation/go-lachesis/kvdb/flushable"
	"github.com/Fantom-foundation/go-lachesis/kvdb/leveldb"
	"github.com/Fantom-foundation/go-lachesis/lachesis"
	"github.com/Fantom-foundation/go-lachesis/lachesis/genesis"
	"github.com/Fantom-foundation/go-lachesis/logger"
)

/*
 * bench:
 */

func BenchmarkStore(b *testing.B) {
	logger.SetTestMode(b)

	benchmarkStore(b)
}

func benchmarkStore(b *testing.B) {
	dir, err := ioutil.TempDir("", "poset-bench")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			panic(err)
		}
	}()

	lvl := leveldb.NewProducer(dir)
	dbs := flushable.NewSyncedPool(lvl)

	input := NewEventStore(dbs.GetDb("input"), 1000)
	defer input.Close()

	store := NewStore(dbs, DefaultStoreConfig())
	defer store.Close()

	nodes := inter.GenNodes(5)

	p := benchPoset(nodes, input, store)

	p.callback.SelectValidatorsGroup = func(oldEpoch, newEpoch idx.Epoch) *pos.Validators {
		if oldEpoch == 1 {
			validators := p.Validators.Builder()
			// move stake from node0 to node1
			validators.Set(nodes[0], 0)
			validators.Set(nodes[1], 2)
			return validators.Build()
		}
		return p.Validators
	}

	// run test with random DAG, N + 1 epochs long
	b.ResetTimer()
	maxEpoch := idx.Epoch(b.N) + 1
	for epoch := idx.Epoch(1); epoch <= maxEpoch; epoch++ {
		r := rand.New(rand.NewSource(int64((epoch))))
		_ = inter.ForEachRandEvent(nodes, int(p.dag.MaxEpochBlocks*3), 3, r, inter.ForEachEvent{
			Process: func(e *inter.Event, name string) {
				input.SetEvent(e)
				err := p.ProcessEvent(e)
				if err != nil {
					panic(err)
				}
				err = dbs.Flush(e.Hash().Bytes())
				if err != nil {
					panic(err)
				}
			},
			Build: func(e *inter.Event, name string) *inter.Event {
				e.Epoch = epoch
				if e.Seq%2 != 0 {
					e.Transactions = append(e.Transactions, &types.Transaction{})
				}
				e.TxHash = types.DeriveSha(e.Transactions, new(trie.Trie))
				return p.Prepare(e)
			},
		})
	}
}

func benchPoset(nodes []idx.StakerID, input EventSource, store *Store) *Poset {
	validators := make(pos.GValidators, 0, len(nodes))
	for _, v := range nodes {
		validators = append(validators, pos.GenesisValidator{
			ID:    v,
			Stake: pos.StakeToBalance(1),
		})
	}

	err := store.ApplyGenesis(&genesis.Genesis{
		Time: genesisTime,
		Alloc: genesis.VAccounts{
			Validators: validators,
			Accounts:   nil,
		},
	}, hash.Event{}, common.Hash{})
	if err != nil {
		panic(err)
	}

	dag := lachesis.FakeNetDagConfig()
	poset := New(dag, store, input)
	poset.Bootstrap(inter.ConsensusCallbacks{})

	return poset
}
