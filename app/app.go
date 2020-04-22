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
		info         *BlockInfo
		evmBlock     *evmcore.EvmBlock
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

// BeginBlock is a prototype of ABCIApplication.BeginBlock
func (a *App) BeginBlock(
	evmBlock *evmcore.EvmBlock,
	cheaters inter.Cheaters,
	stateHash common.Hash,
	blockParticipated map[idx.StakerID]bool,
) {
	info := blockInfo(&evmBlock.EvmHeader)
	a.store.SetBlock(info)
	epoch := a.GetEpoch()

	a.ctx = &blockContext{
		info:         info,
		evmBlock:     evmBlock,
		statedb:      a.store.StateDB(stateHash),
		evmProcessor: evmcore.NewStateProcessor(a.config.Net.EvmChainConfig(), a.BlockChain()),
		cheaters:     cheaters,
		sealEpoch:    a.shouldSealEpoch(info, cheaters),
		gp:           new(evmcore.GasPool),
		totalFee:     big.NewInt(0),
		txCount:      0,
	}
	a.ctx.evmBlock.GasUsed = 0
	a.ctx.gp.AddGas(evmBlock.GasLimit)

	a.updateValidationScores(epoch, info.Index, blockParticipated)
}

// endBlock is a prototype of ABCIApplication.EndBlock
func (a *App) endBlock() (
	evmBlock *evmcore.EvmBlock,
	receipts types.Receipts,
	sealEpoch bool,
) {
	var (
		epoch = a.GetEpoch()
		err   error
	)

	info := a.ctx.info
	evmBlock = a.ctx.evmBlock
	receipts = a.ctx.receipts
	sealEpoch = a.ctx.sealEpoch || sfctype.EpochIsForceSealed(receipts)

	for _, r := range receipts {
		a.store.IndexLogs(r.Logs...)
	}

	if a.config.TxIndex && receipts.Len() > 0 {
		a.store.SetReceipts(info.Index, receipts)
	}

	// Process PoI/score changes
	a.updateOriginationScores()
	a.updateUsersPOI(info, evmBlock, receipts)
	a.updateStakersPOI(info)

	// Process SFC contract transactions
	stats := a.updateEpochStats(epoch, info.Time, a.ctx.totalFee, sealEpoch)
	a.processSfc(epoch, a.ctx.info, receipts, a.ctx.cheaters, stats)
	evmBlock.Root, err = a.ctx.statedb.Commit(true)
	if err != nil {
		a.Log.Crit("Failed to commit state", "err", err)
	}

	a.incLastBlock()
	if sealEpoch {
		a.SetLastVoting(info.Index, info.Time)
		a.incEpoch()
	}

	// free resources
	a.ctx = nil
	a.store.FlushState()

	return
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
