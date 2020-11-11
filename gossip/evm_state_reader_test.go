package gossip

import (
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/Fantom-foundation/go-lachesis/app"
	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/pos"
	"github.com/Fantom-foundation/go-lachesis/lachesis"
	"github.com/Fantom-foundation/go-lachesis/lachesis/genesis"
	"github.com/Fantom-foundation/go-lachesis/logger"
)

func TestGetGenesisBlock(t *testing.T) {
	logger.SetTestMode(t)
	require := require.New(t)

	net := lachesis.FakeNetConfig(genesis.FakeAccounts(0, 5, big.NewInt(0), pos.StakeToBalance(1)))
	addrWithStorage := net.Genesis.Alloc.Accounts.Addresses()[0]
	accountWithCode := net.Genesis.Alloc.Accounts[addrWithStorage]
	accountWithCode.Code = []byte{1, 2, 3}
	accountWithCode.Storage = make(map[common.Hash]common.Hash)
	accountWithCode.Storage[common.Hash{}] = common.BytesToHash(common.Big1.Bytes())
	net.Genesis.Alloc.Accounts[addrWithStorage] = accountWithCode

	apps := app.NewMemStore()
	stateRoot, _, err := apps.ApplyGenesis(&net)
	require.NoError(err)

	store := NewMemStore()
	genesisHash, stateHash, _, err := store.ApplyGenesis(&net, stateRoot)
	require.NoError(err)

	require.NotEqual(common.Hash{}, genesisHash)
	require.NotEqual(common.Hash{}, stateHash)

	reader := EvmStateReader{
		store:    store,
		app:      apps,
		engineMu: new(sync.RWMutex),
	}
	genesisBlock := reader.GetBlock(common.Hash(genesisHash), 0)

	require.Equal(common.Hash(genesisHash), genesisBlock.Hash)
	require.Equal(net.Genesis.Time, genesisBlock.Time)
	require.Empty(genesisBlock.Transactions)

	statedb, err := reader.StateAt(genesisBlock.Root)
	require.NoError(err)
	for addr, account := range net.Genesis.Alloc.Accounts {
		require.Equal(account.Balance.String(), statedb.GetBalance(addr).String())
		require.Equal(account.Code, statedb.GetCode(addr))
		if len(account.Storage) == 0 {
			require.Equal(common.Hash{}, statedb.GetState(addr, common.Hash{}))
		} else {
			for key, val := range account.Storage {
				require.Equal(val, statedb.GetState(addr, key))
			}
		}
	}
}

func TestGetBlock(t *testing.T) {
	logger.SetTestMode(t)
	require := require.New(t)

	net := lachesis.FakeNetConfig(genesis.FakeAccounts(0, 5, big.NewInt(0), pos.StakeToBalance(1)))

	apps := app.NewMemStore()
	stateRoot, _, err := apps.ApplyGenesis(&net)
	require.NoError(err)

	store := NewMemStore()
	genesisHash, _, _, err := store.ApplyGenesis(&net, stateRoot)
	require.NoError(err)

	txs := types.Transactions{}
	key, err := crypto.GenerateKey()
	require.NoError(err)
	for i := 0; i < 6; i++ {
		tx, err := types.SignTx(types.NewTransaction(uint64(i), common.Address{}, big.NewInt(100), 0, big.NewInt(1), nil), types.HomesteadSigner{}, key)
		require.NoError(err)
		txs = append(txs, tx)
	}

	event1 := inter.NewEvent()
	event2 := inter.NewEvent()
	event1.Transactions = txs[:1]
	event1.Seq = 1
	event2.Transactions = txs[1:]
	event1.Seq = 2
	block := inter.NewBlock(1, 123, event2.Hash(), genesisHash, hash.Events{event1.Hash(), event2.Hash()})
	block.SkippedTxs = []uint{0, 2, 4}

	store.SetEvent(event1)
	store.SetEvent(event2)
	store.SetBlock(block)

	reader := EvmStateReader{
		store:    store,
		app:      apps,
		engineMu: new(sync.RWMutex),
	}
	evmBlock := reader.GetDagBlock(block.Atropos, block.Index)

	require.Equal(uint64(block.Index), evmBlock.Number.Uint64())
	require.Equal(common.Hash(block.Atropos), evmBlock.Hash)
	require.Equal(common.Hash(genesisHash), evmBlock.ParentHash)
	require.Equal(block.Time, evmBlock.Time)
	require.Equal(len(txs)-len(block.SkippedTxs), evmBlock.Transactions.Len())
	require.Equal(txs[1].Hash(), evmBlock.Transactions[0].Hash())
	require.Equal(txs[3].Hash(), evmBlock.Transactions[1].Hash())
	require.Equal(txs[5].Hash(), evmBlock.Transactions[2].Hash())
}
