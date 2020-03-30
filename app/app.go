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
		ctx    *blockContext

		logger.Instance
	}

	blockContext struct {
		statedb      *state.StateDB
		evmProcessor *evmcore.StateProcessor
		sealEpoch    bool
		totalFee     *big.Int
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
func (a *App) BeginBlock(block *inter.Block, cheaters inter.Cheaters, stateHash common.Hash, stateReader evmcore.DummyChain) {
	a.store.SetBlock(blockInfo(block))
	a.ctx = &blockContext{
		statedb:      a.store.StateDB(stateHash),
		evmProcessor: evmcore.NewStateProcessor(a.config.EvmChainConfig(), stateReader),
		sealEpoch:    a.shouldSealEpoch(block, cheaters),
	}
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
	bool,
) {
	// Process txs
	receipts, _, gasUsed, totalFee, skipped, err := a.ctx.evmProcessor.
		Process(evmBlock, a.ctx.statedb, vm.Config{}, false)
	if err != nil {
		a.Log.Crit("Shouldn't happen ever because it's not strict", "err", err)
	}
	a.ctx.totalFee = totalFee
	block.SkippedTxs = skipped
	block.GasUsed = gasUsed

	a.ctx.sealEpoch = a.ctx.sealEpoch || sfctype.EpochIsForceSealed(receipts)

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

	a.store.SetReceipts(block.Index, receipts)

	return block, evmBlock, a.ctx.totalFee, receipts, a.ctx.sealEpoch
}

// EndBlock is a prototype of ABCIApplication.EndBlock
func (a *App) EndBlock(
	epoch idx.Epoch,
	block *inter.Block,
	evmBlock *evmcore.EvmBlock,
	receipts types.Receipts,
	cheaters inter.Cheaters,
	stats *sfctype.EpochStats,
	txPositions map[common.Hash]TxPosition,
	blockParticipated map[idx.StakerID]bool,
) common.Hash {
	// Process PoI/score changes
	a.updateOriginationScores(epoch, evmBlock, receipts, txPositions)
	a.updateValidationScores(epoch, block, blockParticipated)
	a.updateUsersPOI(block, evmBlock, receipts)
	a.updateStakersPOI(block)

	a.processSfc(epoch, block, receipts, cheaters, stats)
	newStateHash, err := a.ctx.statedb.Commit(true)
	if err != nil {
		a.Log.Crit("Failed to commit state", "err", err)
	}

	if a.ctx.sealEpoch {
		a.store.SetLastVoting(block.Index, block.Time)
	}

	// free resources
	a.ctx = nil
	a.store.FlushState()

	return newStateHash
}

func (a *App) shouldSealEpoch(block *inter.Block, cheaters inter.Cheaters) bool {
	startBlock, startTime := a.store.GetLastVoting()
	seal := (block.Index - startBlock) >= idx.Block(a.config.Dag.MaxEpochBlocks)
	seal = seal || (block.Time-startTime) >= inter.Timestamp(a.config.Dag.MaxEpochDuration)
	seal = seal || cheaters.Len() > 0

	return seal
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

// blockTime by block number
func (a *App) blockTime(n idx.Block) inter.Timestamp {
	return a.store.GetBlock(n).Time
}
