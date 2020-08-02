package app

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"

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

// beginBlock signals the beginning of a block.
func (a *App) beginBlock(
	evmHeader evmcore.EvmHeader,
	stateRoot common.Hash,
	cheaters inter.Cheaters,
	blockParticipated map[idx.StakerID]bool,
) (sealEpoch bool) {
	block := blockInfo(&evmHeader)
	epoch := a.GetEpoch()
	sealEpoch = a.shouldSealEpoch(block, cheaters)

	prev := a.store.GetBlock(block.Index - 1)
	if prev.Root != stateRoot {
		panic("inconsistent state db")
	}

	a.ctx = &blockContext{
		block:        block,
		header:       &evmHeader,
		statedb:      a.store.StateDB(stateRoot),
		evmProcessor: evmcore.NewStateProcessor(a.config.Net.EvmChainConfig(), a.BlockChain()),
		cheaters:     cheaters,
		sealEpoch:    sealEpoch,
		gp:           new(evmcore.GasPool),
		totalFee:     big.NewInt(0),
		txCount:      0,
	}
	a.ctx.header.GasUsed = 0
	a.ctx.gp.AddGas(evmHeader.GasLimit)

	a.updateValidationScores(epoch, block.Index, blockParticipated)

	return
}

// deliverTx for full processing.
func (a *App) deliverTx(tx *eth.Transaction, originator idx.StakerID) (*eth.Receipt, error) {
	const strict = false
	receipt, fee, skip, err := a.ctx.evmProcessor.
		ProcessTx(tx, a.ctx.txCount, a.ctx.gp, &a.ctx.header.GasUsed, a.ctx.header, a.ctx.statedb, vm.Config{}, strict)
	a.ctx.txCount++
	if !strict && err != nil {
		return nil, err
	}
	if !strict && skip {
		return nil, fmt.Errorf("skipped")
	}

	a.ctx.txs = append(a.ctx.txs, tx)
	a.ctx.receipts = append(a.ctx.receipts, receipt)
	a.ctx.totalFee.Add(a.ctx.totalFee, fee)
	a.store.AddDirtyOriginationScore(originator, fee)

	return receipt, nil
}

// endBlock signals the end of a block, returns changes to the validator set.
func (a *App) endBlock(n idx.Block) {
	if a.ctx.block.Index != n {
		a.Log.Crit("missed block", "current", a.ctx.block.Index, "got", n)
	}

	for _, r := range a.ctx.receipts {
		a.store.IndexLogs(r.Logs...)
	}

	if a.config.TxIndex && a.ctx.receipts.Len() > 0 {
		a.store.SetReceipts(a.ctx.block.Index, a.ctx.receipts)
	}

	// Process PoI/score changes
	a.updateOriginationScores(a.ctx.sealEpoch)
	a.updateUsersPOI(a.ctx.block, a.ctx.txs, a.ctx.receipts)
	a.updateStakersPOI(a.ctx.block)

	// Process SFC contract transactions
	epoch := a.GetEpoch()
	stats := a.updateEpochStats(epoch, a.ctx.block.Time, a.ctx.totalFee, a.ctx.sealEpoch)
	a.processSfc(epoch, a.ctx.block, a.ctx.receipts, a.ctx.cheaters, stats)

	a.incLastBlock()
	if a.ctx.sealEpoch {
		a.SetLastVoting(a.ctx.block.Index, a.ctx.block.Time)
		a.incEpoch()
	}
}

// commit the state and return the application Merkle root hash.
func (a *App) commit() common.Hash {
	root, err := a.ctx.statedb.Commit(true)
	if err != nil {
		a.Log.Crit("Failed to commit state", "err", err)
	}

	a.ctx.block.Root = root
	a.store.SetBlock(a.ctx.block)

	// notify
	var logs []*eth.Log
	for _, r := range a.ctx.receipts {
		for _, l := range r.Logs {
			logs = append(logs, l)
		}
	}

	a.Feed.newTxs.Send(core.NewTxsEvent{Txs: a.ctx.txs})
	a.Feed.newLogs.Send(logs)

	// free resources
	a.ctx = nil
	a.store.FlushState()

	return root
}

func (a *App) GetTotalSupply() *big.Int {
	return a.store.GetTotalSupply()
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
