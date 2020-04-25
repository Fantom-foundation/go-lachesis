package app

import (
	"math"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/tendermint/tendermint/abci/types"

	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/sfctype"
)

const (
	txIsSkipped = 1
)

// InitChain implements ABCIApplication.InitChain.
// It should be called once upon genesis.
func (a *App) InitChain(req types.RequestInitChain) types.ResponseInitChain {
	chain := a.config.Net.ChainInfo()
	if !reflect.DeepEqual(req, chain) {
		panic("incompatible chain")
	}

	a.Bootstrap()
	return types.ResponseInitChain{}
}

// BeginBlock signals the beginning of a block.
// It implements ABCIApplication.BeginBlock (prototype).
func (a *App) BeginBlock(
	req types.RequestBeginBlock,
	blockParticipated map[idx.StakerID]bool,
) types.ResponseBeginBlock {
	evmHeader := extractEvmHeader(req)
	stateRoot := extractStateRoot(req)
	cheaters := extractCheaters(req)

	a.beginBlock(evmHeader, stateRoot, cheaters, blockParticipated)

	return types.ResponseBeginBlock{}
}

// DeliverTx for full processing.
// It implements ABCIApplication.DeliverTx.
func (a *App) DeliverTx(req types.RequestDeliverTx) types.ResponseDeliverTx {
	const strict = false

	dagTx := BytesToDagTx(req.Tx)
	tx := dagTx.Transaction

	receipt, fee, skip, err := a.ctx.evmProcessor.
		ProcessTx(tx, a.ctx.txCount, a.ctx.gp, &a.ctx.header.GasUsed, a.ctx.header, a.ctx.statedb, vm.Config{}, strict)
	a.ctx.txCount++
	if !strict && (skip || err != nil) {
		return types.ResponseDeliverTx{
			Code:      txIsSkipped,
			Info:      "skipped",
			GasWanted: int64(tx.Gas()),
			GasUsed:   0,
		}
	}

	a.ctx.txs = append(a.ctx.txs, tx)
	a.ctx.receipts = append(a.ctx.receipts, receipt)
	a.ctx.totalFee.Add(a.ctx.totalFee, fee)
	a.store.AddDirtyOriginationScore(dagTx.Originator, fee)

	return types.ResponseDeliverTx{
		Info:      "ok",
		GasWanted: int64(tx.Gas()),
		GasUsed:   int64(receipt.GasUsed),
	}
}

// EndBlock signals the end of a block, returns changes to the validator set.
// It implements ABCIApplication.EndBlock.
func (a *App) EndBlock(req types.RequestEndBlock) (types.ResponseEndBlock, bool) {
	if a.ctx.block.Index != idx.Block(req.Height) {
		a.Log.Crit("missed block", "current", a.ctx.block.Index, "got", req.Height)
	}

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

	// TODO: replace sealEpoch with response validator set changes.
	return types.ResponseEndBlock{}, sealEpoch
}

// Commit the state and return the application Merkle root hash.
// It implements ABCIApplication.Commit.
func (a *App) Commit() types.ResponseCommit {
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

	return types.ResponseCommit{
		Data: root.Bytes(),
	}
}

func extractEvmHeader(req types.RequestBeginBlock) evmcore.EvmHeader {
	return evmcore.EvmHeader{
		Number:     big.NewInt(req.Header.Height),
		Time:       inter.TimeToStamp(req.Header.Time),
		Hash:       common.BytesToHash(req.Header.ConsensusHash),
		ParentHash: common.BytesToHash(req.Header.LastBlockId.Hash),
		GasLimit:   math.MaxUint64,
	}
}

func extractCheaters(req types.RequestBeginBlock) inter.Cheaters {
	cheaters := make(inter.Cheaters, len(req.ByzantineValidators))
	for i, evil := range req.ByzantineValidators {
		cheaters[i] = idx.StakerID(uint32(evil.Height))
	}

	return cheaters
}

func extractStateRoot(req types.RequestBeginBlock) common.Hash {
	return common.BytesToHash(
		req.Header.LastCommitHash)
}
