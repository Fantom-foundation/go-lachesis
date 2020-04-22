package app

import (
	"math/big"
	"reflect"

	"github.com/Fantom-foundation/go-lachesis/inter/idx"

	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/tendermint/tendermint/abci/types"

	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/inter"
)

// InitChain implements ABCIApplication.InitChain.
// It should be Called once upon genesis.
func (a *App) InitChain(req types.RequestInitChain) types.ResponseInitChain {
	chain := a.config.Net.ChainInfo()
	if !reflect.DeepEqual(req, chain) {
		panic("incompatible chain")
	}

	a.Bootstrap()
	return types.ResponseInitChain{}
}

// DeliverTx implements ABCIApplication.DeliverTx
func (a *App) DeliverTx(req types.RequestDeliverTx) types.ResponseDeliverTx {
	const strict = false

	dagTx := BytesToDagTx(req.Tx)
	tx := dagTx.Transaction

	receipt, fee, skip, err := a.ctx.evmProcessor.
		ProcessTx(tx, a.ctx.txCount, a.ctx.gp, &a.ctx.block.GasUsed, a.ctx.evmBlock, a.ctx.statedb, vm.Config{}, strict)
	a.ctx.txCount++
	if !strict && (skip || err != nil) {
		a.ctx.block.SkippedTxs = append(a.ctx.block.SkippedTxs, a.ctx.txCount)
		return types.ResponseDeliverTx{
			Info:      "skipped",
			GasWanted: int64(tx.Gas()),
			GasUsed:   0,
		}
	}

	a.ctx.totalFee.Add(a.ctx.totalFee, fee)
	a.ctx.receipts = append(a.ctx.receipts, receipt)
	a.store.AddDirtyOriginationScore(dagTx.Originator, fee)

	return types.ResponseDeliverTx{
		Info:      "ok",
		GasWanted: int64(tx.Gas()),
		GasUsed:   int64(receipt.GasUsed),
	}
}

// EndBlock implements ABCIApplication.EndBlock
func (a *App) EndBlock(
	req types.RequestEndBlock,
) (
	// TODO: return only types.ResponseEndBlock
	block *inter.Block,
	evmBlock *evmcore.EvmBlock,
	receipts eth.Receipts,
	totalFee *big.Int,
	sealEpoch bool,
) {
	if a.ctx.block.Index != idx.Block(req.Height) {
		a.Log.Crit("missed block", "current", a.ctx.block.Index, "got", req.Height)
	}

	return a.endBlock()
	/*
		return types.ResponseEndBlock{
			ValidatorUpdates: types.ValidatorUpdates{
				types.ValidatorUpdate{},
				types.ValidatorUpdate{},
			},
		}
	*/
}
