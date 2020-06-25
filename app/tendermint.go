package app

import (
	"math"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/abci/types"

	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
)

const (
	txIsSkipped = 1
)

// InitChain should be called once upon genesis.
// Wraps Bootstrap() to implement ABCIApplication.InitChain.
func (a *App) InitChain(req types.RequestInitChain) types.ResponseInitChain {
	chain := a.config.Net.ChainInfo()
	if !reflect.DeepEqual(req, chain) {
		panic("incompatible chain")
	}

	a.Bootstrap()
	return types.ResponseInitChain{}
}

// CheckTx validates a tx for the mempool.
// It implements ABCIApplication.CheckTx.
func (a *App) CheckTx(req types.RequestCheckTx) types.ResponseCheckTx {
	return types.ResponseCheckTx{}
}

// BeginBlock signals the beginning of a block.
// Wraps beginBlock() to implement ABCIApplication.BeginBlock.
func (a *App) BeginBlock(req types.RequestBeginBlock) types.ResponseBeginBlock {
	evmHeader := extractEvmHeader(req)
	stateRoot := extractStateRoot(req)
	cheaters := extractCheaters(req)
	blockParticipated := extractParticipated(req)

	a.beginBlock(evmHeader, stateRoot, cheaters, blockParticipated)

	return types.ResponseBeginBlock{}
}

// DeliverTx for full processing.
// Wraps deliverTx() to implement ABCIApplication.DeliverTx.
func (a *App) DeliverTx(req types.RequestDeliverTx) types.ResponseDeliverTx {
	dagTx := BytesToDagTx(req.Tx)
	tx := dagTx.Transaction

	receipt, err := a.deliverTx(tx, dagTx.Originator)
	if err != nil {
		return types.ResponseDeliverTx{
			Code:      txIsSkipped,
			Info:      "skipped",
			GasWanted: int64(tx.Gas()),
			GasUsed:   0,
		}
	}

	return types.ResponseDeliverTx{
		Info:      "ok",
		GasWanted: int64(tx.Gas()),
		GasUsed:   int64(receipt.GasUsed),
	}
}

// EndBlock signals the end of a block, returns changes to the validator set.
// Wraps endBlock() to implement ABCIApplication.EndBlock.
func (a *App) EndBlock(req types.RequestEndBlock) types.ResponseEndBlock {
	n := idx.Block(req.Height)

	sealEpoch := a.endBlock(n)

	res := types.ResponseEndBlock{}
	if sealEpoch {
		res.Events = []types.Event{
			{Type: "epoch sealed"},
		}
	}
	return res
}

// Commit the state and return the application Merkle root hash.
// Wraps commit() to implement ABCIApplication.Commit.
func (a *App) Commit() types.ResponseCommit {
	root := a.commit()
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

func extractParticipated(req types.RequestBeginBlock) map[idx.StakerID]bool {
	res := make(map[idx.StakerID]bool, len(req.LastCommitInfo.Votes))
	for _, v := range req.LastCommitInfo.Votes {
		staker := idx.BytesToStakerID(v.Validator.Address)
		res[staker] = true
	}
	return res
}

// Info returns application info.
// It implements ABCIApplication.Info.
func (a *App) Info(types.RequestInfo) types.ResponseInfo {
	return types.ResponseInfo{}
}

// SetOption sets application option.
// It implements ABCIApplication.SetOption.
func (a *App) SetOption(types.RequestSetOption) types.ResponseSetOption {
	return types.ResponseSetOption{}
}

// Query for state.
// It implements ABCIApplication.Query.
func (a *App) Query(types.RequestQuery) types.ResponseQuery {
	return types.ResponseQuery{}
}
