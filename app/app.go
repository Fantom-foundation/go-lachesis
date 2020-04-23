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
		config     Config
		store      *Store
		ctx        *blockContext
		checkpoint Checkpoint

		Feed
		logger.Instance
	}

	blockContext struct {
		block        *BlockInfo
		header       *evmcore.EvmHeader
		txs          types.Transactions
		statedb      *state.StateDB
		evmProcessor *evmcore.StateProcessor
		cheaters     inter.Cheaters
		sealEpoch    bool
		totalFee     *big.Int
		receipts     types.Receipts
		gp           *evmcore.GasPool
		txCount      uint
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

// BeginBlock signals the beginning of a block.
// It implements ABCIApplication.BeginBlock (prototype).
func (a *App) BeginBlock(
	evmHeader evmcore.EvmHeader,
	cheaters inter.Cheaters,
	stateHash common.Hash,
	blockParticipated map[idx.StakerID]bool,
) {
	info := blockInfo(&evmHeader)
	a.store.SetBlock(info)
	epoch := a.GetEpoch()

	a.ctx = &blockContext{
		block:        info,
		header:       &evmHeader,
		statedb:      a.store.StateDB(stateHash),
		evmProcessor: evmcore.NewStateProcessor(a.config.Net.EvmChainConfig(), a.BlockChain()),
		cheaters:     cheaters,
		sealEpoch:    a.shouldSealEpoch(info, cheaters),
		gp:           new(evmcore.GasPool),
		totalFee:     big.NewInt(0),
		txCount:      0,
	}
	a.ctx.header.GasUsed = 0
	a.ctx.gp.AddGas(evmHeader.GasLimit)

	a.updateValidationScores(epoch, info.Index, blockParticipated)
}

// endBlock is a prototype of ABCIApplication.EndBlock
func (a *App) endBlock() (types.Receipts, bool) {

	sealEpoch := a.ctx.sealEpoch || sfctype.EpochIsForceSealed(a.ctx.receipts)

	for _, r := range a.ctx.receipts {
		a.store.IndexLogs(r.Logs...)
	}

	if a.config.TxIndex && a.ctx.receipts.Len() > 0 {
		a.store.SetReceipts(a.ctx.block.Index, a.ctx.receipts)
	}

	// Process PoI/score changes
	a.updateOriginationScores(sealEpoch)
	a.updateUsersPOI(a.ctx.block, a.ctx.txs, a.ctx.receipts)
	a.updateStakersPOI(a.ctx.block)

	// Process SFC contract transactions
	epoch := a.GetEpoch()
	stats := a.updateEpochStats(epoch, a.ctx.block.Time, a.ctx.totalFee, sealEpoch)
	a.processSfc(epoch, a.ctx.block, a.ctx.receipts, a.ctx.cheaters, stats)

	a.incLastBlock()
	if sealEpoch {
		a.SetLastVoting(a.ctx.block.Index, a.ctx.block.Time)
		a.incEpoch()
	}

	return a.ctx.receipts, sealEpoch
}

func (a *App) shouldSealEpoch(block *BlockInfo, cheaters inter.Cheaters) bool {
	startBlock, startTime := a.GetLastVoting()
	seal := (block.Index - startBlock) >= idx.Block(a.config.Net.Dag.MaxEpochBlocks)
	seal = seal || (block.Time-startTime) >= inter.Timestamp(a.config.Net.Dag.MaxEpochDuration)
	seal = seal || cheaters.Len() > 0

	return seal
}

// blockTime by block number
func (a *App) blockTime(n idx.Block) inter.Timestamp {
	return a.store.GetBlock(n).Time
}
