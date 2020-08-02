package gossip

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	tendermint "github.com/tendermint/tendermint/abci/types"

	"github.com/Fantom-foundation/go-lachesis/app"

	"github.com/Fantom-foundation/go-lachesis/inter"
)

// applyNewState moves the state according to new block (txs execution, SFC logic, epoch sealing)
func (s *Service) applyNewState(
	abci tendermint.Application,
	block *inter.Block,
	cheaters inter.Cheaters,
) (
	*inter.Block,
	map[common.Hash]*app.TxPosition,
	types.Transactions,
	bool,
) {
	// s.engineMu is locked here

	start := time.Now()

	events, allTxs := s.usedEvents(block)
	unusedCount := len(block.Events) - len(events)
	block.Events = block.Events[unusedCount:]

	// memorize position of each tx, for indexing and origination scores
	txsPositions := make(map[common.Hash]*app.TxPosition)
	for _, e := range events {
		for i, tx := range e.Transactions {
			// if tx was met in multiple events, then assign to first ordered event
			if _, ok := txsPositions[tx.Hash()]; ok {
				continue
			}
			txsPositions[tx.Hash()] = &app.TxPosition{
				Event:       e.Hash(),
				Creator:     e.Creator,
				EventOffset: uint32(i),
			}
		}
	}

	epoch := s.engine.GetEpoch()
	stateRoot := s.store.GetBlock(block.Index - 1).Root

	var sealEpoch bool
	resp := abci.BeginBlock(
		beginBlockRequest(cheaters, stateRoot, block, s.blockParticipated))
	for _, appEvent := range resp.Events {
		switch appEvent.Type {
		case "epoch sealed":
			sealEpoch = true
		}
	}

	okTxs := make(types.Transactions, 0, len(allTxs))
	block.SkippedTxs = make([]uint, 0, len(allTxs))
	for i, tx := range allTxs {
		originator := txsPositions[tx.Hash()].Creator
		req := deliverTxRequest(tx, originator)
		resp := abci.DeliverTx(req)
		block.GasUsed += uint64(resp.GasUsed)

		if resp.Code != txIsFullyValid {
			block.SkippedTxs = append(block.SkippedTxs, uint(i))
			continue
		}
		okTxs = append(okTxs, tx)
		txsPositions[tx.Hash()].Block = block.Index
		notUsedGas := resp.GasWanted - resp.GasUsed
		s.store.IncGasPowerRefund(epoch, originator, notUsedGas)

		if resp.Log != "" {
			s.Log.Info("tx processed", "log", resp.Log)
		}
	}

	abci.EndBlock(endBlockRequest(block.Index))

	commit := abci.Commit()
	block.Root = common.BytesToHash(commit.Data)
	block.TxHash = types.DeriveSha(okTxs)

	// process new epoch
	if sealEpoch {
		// prune not needed gas power records
		s.store.DelGasPowerRefunds(epoch - 1)
		s.onEpochSealed(block, cheaters)
	}

	log.Info("New block",
		"index", block.Index,
		"atropos", block.Atropos,
		"gasUsed", block.GasUsed,
		"skipped_txs", len(block.SkippedTxs),
		"txs", len(okTxs),
		"t", time.Since(start))

	return block, txsPositions, okTxs, sealEpoch
}
