package app

import (
	"reflect"

	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/tendermint/tendermint/abci/types"
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

	tx := new(eth.Transaction)
	err := rlp.DecodeBytes(req.Tx, tx)
	if err != nil {
		a.Log.Crit("tx decode error", "err", err)
	}

	receipt, fee, skip, err := a.ctx.evmProcessor.
		ProcessTx(tx, a.ctx.txN, a.ctx.gp, &a.ctx.block.GasUsed, a.ctx.evmBlock, a.ctx.statedb, vm.Config{}, strict)
	a.ctx.txN++
	if !strict && (skip || err != nil) {
		a.ctx.block.SkippedTxs = append(a.ctx.block.SkippedTxs, a.ctx.txN)
		return types.ResponseDeliverTx{
			Info:      "skipped",
			GasWanted: int64(tx.Gas()),
			GasUsed:   0,
		}
	}

	a.ctx.totalFee.Add(a.ctx.totalFee, fee)
	a.ctx.receipts = append(a.ctx.receipts, receipt)

	return types.ResponseDeliverTx{
		Info:      "ok",
		GasWanted: int64(tx.Gas()),
		GasUsed:   int64(receipt.GasUsed),
	}
}

/*
// EndBlock implements ABCIApplication.EndBlock
func (a *App) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	return types.ResponseEndBlock{
		ValidatorUpdates: types.ValidatorUpdates{
			types.ValidatorUpdate{},
			types.ValidatorUpdate{},
		},
	}
}
*/
