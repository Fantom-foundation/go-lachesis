package app

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
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

func (a *App) beginBlock(
	evmHeader evmcore.EvmHeader,
	stateRoot common.Hash,
	cheaters inter.Cheaters,
	blockParticipated map[idx.StakerID]bool,
) {
	block := blockInfo(&evmHeader)
	epoch := a.GetEpoch()

	prev := a.store.GetBlock(block.Index - 1)
	if prev.Root != stateRoot {
		//panic("inconsistent state db")
	}

	a.ctx = &blockContext{
		block:        block,
		header:       &evmHeader,
		statedb:      a.store.StateDB(stateRoot),
		evmProcessor: evmcore.NewStateProcessor(a.config.Net.EvmChainConfig(), a.BlockChain()),
		cheaters:     cheaters,
		sealEpoch:    a.shouldSealEpoch(block, cheaters),
		gp:           new(evmcore.GasPool),
		totalFee:     big.NewInt(0),
		txCount:      0,
	}
	a.ctx.header.GasUsed = 0
	a.ctx.gp.AddGas(evmHeader.GasLimit)

	a.updateValidationScores(epoch, block.Index, blockParticipated)
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
	if a.ctx.block.Index == n {
		return a.ctx.block.Time
	}

	return a.store.GetBlock(n).Time
}
