package app

import (
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
	evmHeader evmcore.EvmHeader,
	blockParticipated map[idx.StakerID]bool,
) types.ResponseBeginBlock {
	cheaters := extractCheaters(req)
	stateRoot := extractStateRoot(req)
	a.beginBlock(evmHeader, stateRoot, cheaters, blockParticipated)
	return types.ResponseBeginBlock{}
	/*
		INFO [04-25|00:09:52.603] New event                                id=327:8271:17410d  parents=7 by=29 frame=806:n txs=0 t=841.932µs
		INFO [04-25|00:09:52.604] New event                                id=327:8270:d1966b  parents=7 by=25 frame=806:n txs=0 t=1.077515ms
		INFO [04-25|00:09:52.605] New event                                id=327:8270:e08b56  parents=7 by=3  frame=806:n txs=0 t=631.612µs
		INFO [04-25|00:09:52.606] New event                                id=327:8271:161b7a  parents=7 by=21 frame=806:n txs=0 t=666.152µs
		INFO [04-25|00:09:52.607] New event                                id=327:8271:0ea083  parents=7 by=20 frame=806:n txs=0 t=649.753µs
		INFO [04-25|00:09:52.608] New event                                id=327:8270:eb571d  parents=7 by=14 frame=806:n txs=0 t=688.922µs
		INFO [04-25|00:09:52.609] New event                                id=327:8271:2787af  parents=7 by=15 frame=806:n txs=0 t=1.009614ms
		INFO [04-25|00:09:52.653] New block                                index=40625 atropos=327:8257:5a0dbe gasUsed=0 skipped_txs=1 txs=0 t=39.431757ms
		INFO [04-25|00:09:52.732] Allocated cache and file handles         database=222/poset-epoch-328-ldb  cache=64.00MiB handles=16
		INFO [04-25|00:09:52.745] Allocated cache and file handles         database=222/gossip-epoch-328-ldb cache=64.00MiB handles=16
		INFO [04-25|00:09:52.814] New event                                id=327:8272:cc109a  parents=7 by=3  frame=807:y txs=0 t=204.614082ms
		WARN [04-25|00:09:53.391] Incoming event rejected                  event=328:1:70a872 creator=5 err="event has wrong GasPowerLeft"
		WARN [04-25|00:09:54.163] Events request error                     peer=41459d863ead95a4 err="shutting down"
		WARN [04-25|00:09:55.380] Events request error                     peer=41459d863ead95a4 err="shutting down"
	*/
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

	// notify
	var logs []*eth.Log
	for _, r := range a.ctx.receipts {
		for _, l := range r.Logs {
			logs = append(logs, l)
		}
	}
	a.Feed.newBlock.Send(evmcore.ChainHeadNotify{
		Block: &evmcore.EvmBlock{
			EvmHeader:    *a.ctx.header,
			Transactions: a.ctx.txs,
		}})
	a.Feed.newTxs.Send(core.NewTxsEvent{Txs: a.ctx.txs})
	a.Feed.newLogs.Send(logs)

	// free resources
	a.ctx = nil
	a.store.FlushState()

	return types.ResponseCommit{
		Data: root.Bytes(),
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
	return common.BytesToHash(req.Header.LastCommitHash)
}
