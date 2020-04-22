package app

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/sfctype"
	"github.com/Fantom-foundation/go-lachesis/logger"
)

type (
	// App is a prototype of Tendermint ABCI Application
	App struct {
		config Config
		store  *Store
		ctx    *blockContext

		checkpoint Checkpoint

		logger.Instance
	}

	blockContext struct {
		block        *inter.Block
		evmBlock     *evmcore.EvmBlock
		statedb      *state.StateDB
		evmProcessor *evmcore.StateProcessor
		cheaters     inter.Cheaters
		sealEpoch    bool
		totalFee     *big.Int
		receipts     types.Receipts
		gp           *evmcore.GasPool
		txN          uint
	}
)

// New is a constructor
func New(cfg Config, s *Store) *App {
	return &App{
		config: cfg,
		store:  s,

		Instance: logger.MakeInstance(),
	}
}

// BeginBlock is a prototype of ABCIApplication.BeginBlock
func (a *App) BeginBlock(
	block *inter.Block, evmBlock *evmcore.EvmBlock, cheaters inter.Cheaters, stateHash common.Hash,
) {
	a.store.SetBlock(blockInfo(block))
	block.SkippedTxs = make([]uint, 0, len(evmBlock.Transactions))
	a.ctx = &blockContext{
		block:        block,
		evmBlock:     evmBlock,
		statedb:      a.store.StateDB(stateHash),
		evmProcessor: evmcore.NewStateProcessor(a.config.Net.EvmChainConfig(), a.BlockChain()),
		cheaters:     cheaters,
		sealEpoch:    a.shouldSealEpoch(block, cheaters),
		gp:           new(evmcore.GasPool),
		totalFee:     big.NewInt(0),
		txN:          0,
	}

	a.ctx.gp.AddGas(evmBlock.GasLimit)
}

// EndBlock is a prototype of ABCIApplication.EndBlock
func (a *App) EndBlock(
	txPositions map[common.Hash]TxPosition,
	blockParticipated map[idx.StakerID]bool,
) (
	block *inter.Block,
	evmBlock *evmcore.EvmBlock,
	receipts types.Receipts,
	totalFee *big.Int,
	sealEpoch bool,
) {
	var (
		epoch = a.GetEpoch()
		err   error
	)

	block = a.ctx.block
	evmBlock = a.ctx.evmBlock
	receipts = a.ctx.receipts
	totalFee = a.ctx.totalFee
	sealEpoch = a.ctx.sealEpoch || sfctype.EpochIsForceSealed(receipts)

	evmBlock.Transactions = filterSkippedTxs(block, evmBlock.Transactions)
	block.TxHash = types.DeriveSha(evmBlock.Transactions)
	*evmBlock = evmcore.EvmBlock{
		EvmHeader:    *evmcore.ToEvmHeader(block),
		Transactions: evmBlock.Transactions,
	}

	for _, r := range receipts {
		a.store.IndexLogs(r.Logs...)
	}

	if a.config.TxIndex && receipts.Len() > 0 {
		a.store.SetReceipts(block.Index, receipts)
	}

	// Process PoI/score changes
	a.updateOriginationScores(epoch, evmBlock, receipts, txPositions)
	a.updateValidationScores(epoch, block, blockParticipated)
	a.updateUsersPOI(block, evmBlock, receipts)
	a.updateStakersPOI(block)

	// Process SFC contract transactions
	stats := a.updateEpochStats(epoch, block.Time, totalFee, sealEpoch)
	a.processSfc(epoch, block, receipts, a.ctx.cheaters, stats)
	a.ctx.block.Root, err = a.ctx.statedb.Commit(true)
	if err != nil {
		a.Log.Crit("Failed to commit state", "err", err)
	}

	a.incLastBlock()
	if sealEpoch {
		a.SetLastVoting(block.Index, block.Time)
		a.incEpoch()
	}

	// free resources
	a.ctx = nil
	a.store.FlushState()

	return
}

func (a *App) shouldSealEpoch(block *inter.Block, cheaters inter.Cheaters) bool {
	startBlock, startTime := a.GetLastVoting()
	seal := (block.Index - startBlock) >= idx.Block(a.config.Net.Dag.MaxEpochBlocks)
	seal = seal || (block.Time-startTime) >= inter.Timestamp(a.config.Net.Dag.MaxEpochDuration)
	seal = seal || cheaters.Len() > 0

	return seal
}

func filterSkippedTxs(block *inter.Block, txs types.Transactions) types.Transactions {
	// receipts are filtered already
	skipCount := 0
	filteredTxs := make(types.Transactions, 0, len(txs))
	for i, tx := range txs {
		if skipCount < len(block.SkippedTxs) && block.SkippedTxs[skipCount] == uint(i) {
			skipCount++
		} else {
			filteredTxs = append(filteredTxs, tx)
		}
	}
	return filteredTxs
}

// blockTime by block number
func (a *App) blockTime(n idx.Block) inter.Timestamp {
	return a.store.GetBlock(n).Time
}
