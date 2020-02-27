package app

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/sfctype"
	"github.com/Fantom-foundation/go-lachesis/lachesis"
	"github.com/Fantom-foundation/go-lachesis/logger"
)

type (
	// App is a prototype of Tendermint ABCI Application
	App struct {
		config lachesis.Config
		store  *Store
		*blockContext

		logger.Instance
	}

	blockContext struct {
		statedb      *state.StateDB
		evmProcessor *evmcore.StateProcessor
		sealEpoch    bool
	}
)

// New is a constructor
func New(cfg lachesis.Config, s *Store) *App {
	return &App{
		config: cfg,
		store:  s,

		Instance: logger.MakeInstance(),
	}
}

// BeginBlock is a prototype of ABCIApplication.BeginBlock
func (a *App) BeginBlock(block *inter.Block, cheaters inter.Cheaters, stateHash common.Hash, stateReader evmcore.DummyChain) bool {
	startBlock, startTime := a.store.GetLastVoting()
	sealEpoch := (block.Index - startBlock) >= idx.Block(a.config.Dag.MaxEpochBlocks)
	sealEpoch = sealEpoch || (block.Time-startTime) >= inter.Timestamp(a.config.Dag.MaxEpochDuration)
	sealEpoch = sealEpoch || cheaters.Len() > 0

	a.blockContext = &blockContext{
		statedb:      a.store.StateDB(stateHash),
		evmProcessor: evmcore.NewStateProcessor(a.config.EvmChainConfig(), stateReader),
		sealEpoch:    sealEpoch,
	}

	return sealEpoch
}

// DeliverTxs includes a set of ABCIApplication.DeliverTx() calls
// It execs ordered txns of new block on state.
func (a *App) DeliverTxs(
	block *inter.Block,
	evmBlock *evmcore.EvmBlock,
) (
	*inter.Block,
	*evmcore.EvmBlock,
	*big.Int,
	types.Receipts,
) {
	// Process txs
	receipts, _, gasUsed, totalFee, skipped, err := a.blockContext.evmProcessor.
		Process(evmBlock, a.blockContext.statedb, vm.Config{}, false)
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
		a.store.IndexLogs(r.Logs...)
	}

	return block, evmBlock, totalFee, receipts
}

// EndBlock is a prototype of ABCIApplication.EndBlock
func (a *App) EndBlock(
	epoch idx.Epoch,
	block *inter.Block,
	receipts types.Receipts,
	cheaters inter.Cheaters,
	stats *sfctype.EpochStats,
) common.Hash {

	a.processSfc(epoch, block, receipts, a.blockContext.sealEpoch, cheaters, stats)
	newStateHash, err := a.blockContext.statedb.Commit(true)
	if err != nil {
		a.Log.Crit("Failed to commit state", "err", err)
	}

	if a.blockContext.sealEpoch {
		a.store.SetLastVoting(block.Index, block.Time)
	}

	// free resources
	a.blockContext = nil

	return newStateHash
}

func filterSkippedTxs(block *inter.Block, evmBlock *evmcore.EvmBlock) *evmcore.EvmBlock {
	// Filter skipped transactions. Receipts are filtered already
	skipCount := 0
	filteredTxs := make(types.Transactions, 0, len(evmBlock.Transactions))
	for i, tx := range evmBlock.Transactions {
		if skipCount < len(block.SkippedTxs) && block.SkippedTxs[skipCount] == uint(i) {
			skipCount++
		} else {
			filteredTxs = append(filteredTxs, tx)
		}
	}
	evmBlock.Transactions = filteredTxs
	return evmBlock
}
