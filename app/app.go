package app

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/lachesis"
	"github.com/Fantom-foundation/go-lachesis/logger"
)

// Application logic: EVM and Economy.
type Application struct {
	config *lachesis.Config
	store  *Store

	engineMu *sync.RWMutex

	logger.Instance

	Gossip gossip
}

// NewApplication constructor.
func NewApplication(cfg *lachesis.Config, store *Store) *Application {
	app := &Application{
		config: cfg,
		store:  store,
	}

	return app
}

// ApplyNewState moves the state according to new block (txs execution, SFC logic, epoch sealing)
func (a *Application) ApplyNewState(
	prev, block *inter.Block,
	evmBlock *evmcore.EvmBlock,
	txPositions map[common.Hash]TxPosition,
	epoch idx.Epoch,
	sealEpoch bool,
	cheaters inter.Cheaters,
) (
	*inter.Block,
	*evmcore.EvmBlock,
	types.Receipts,
	map[common.Hash]TxPosition,
	*big.Int,
) {
	// Get stateDB
	statedb := a.store.StateDB(prev.Root)

	// Process EVM txs
	block, evmBlock, totalFee, receipts := a.executeEvmTransactions(block, evmBlock, statedb)

	// memorize block position of each tx, for indexing and origination scores
	for i, tx := range evmBlock.Transactions {
		// not skipped txs only
		position := txPositions[tx.Hash()]
		position.Block = block.Index
		position.BlockOffset = uint32(i)
		txPositions[tx.Hash()] = position
	}

	// Process PoI/score changes
	a.updateOriginationScores(block, evmBlock, receipts, txPositions, epoch, sealEpoch)
	a.updateValidationScores(prev, block, sealEpoch)
	a.updateUsersPOI(block, evmBlock, receipts, totalFee, sealEpoch)
	a.updateStakersPOI(prev, block, sealEpoch)

	// Process SFC contract transactions
	a.processSfc(block, receipts, totalFee, epoch, sealEpoch, cheaters, statedb)

	// Get state root
	newStateHash, err := statedb.Commit(true)
	if err != nil {
		a.Log.Crit("Failed to commit state", "err", err)
	}
	block.Root = newStateHash
	*evmBlock = evmcore.EvmBlock{
		EvmHeader:    *evmcore.ToEvmHeader(block),
		Transactions: evmBlock.Transactions,
	}

	return block, evmBlock, receipts, txPositions, totalFee
}

// executeTransactions execs ordered txns of new block on state.
func (a *Application) executeEvmTransactions(
	block *inter.Block,
	evmBlock *evmcore.EvmBlock,
	statedb *state.StateDB,
) (
	*inter.Block,
	*evmcore.EvmBlock,
	*big.Int,
	types.Receipts,
) {
	// s.engineMu is locked here

	evmProcessor := evmcore.NewStateProcessor(a.config.EvmChainConfig(), a.Gossip.GetDummyChainReader())

	// Process txs
	receipts, _, gasUsed, totalFee, skipped, err := evmProcessor.Process(evmBlock, statedb, vm.Config{}, false)
	if err != nil {
		a.Log.Crit("Shouldn't happen ever because it's not strict", "err", err)
	}
	block.SkippedTxs = skipped
	block.GasUsed = gasUsed

	// Filter skipped transactions
	evmBlock = filterSkippedTxs(block, evmBlock)

	block.TxHash = types.DeriveSha(evmBlock.Transactions)
	*evmBlock = evmcore.EvmBlock{
		EvmHeader:    *evmcore.ToEvmHeader(block),
		Transactions: evmBlock.Transactions,
	}

	for _, r := range receipts {

		err := a.store.table.EvmLogs.Push(r.Logs...)
		if err != nil {
			a.Log.Crit("DB logs index", "err", err)
		}
	}

	return block, evmBlock, totalFee, receipts
}
