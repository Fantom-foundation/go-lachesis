package gossip

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	tendermint "github.com/tendermint/tendermint/abci/types"

	"github.com/Fantom-foundation/go-lachesis/app"
	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/tracing"
)

// applyNewStateAsync is an async option of Service.applyNewStateSync()
func (s *Service) applyNewStateAsync(
	abci tendermint.Application,
	block *inter.Block,
	cheaters inter.Cheaters,
	blockParticipated map[idx.StakerID]bool,
) (
	sealEpoch bool,
	blockTxHashes []common.Hash,
) {
	// s.engineMu is locked here

	start := time.Now()

	stateRoot := s.store.GetBlock(block.Index - 1).Root

	blockParticipatedCopy := make(map[idx.StakerID]bool, len(blockParticipated))
	for k, v := range blockParticipated {
		blockParticipatedCopy[k] = v
	}

	resp := abci.BeginBlock(
		beginBlockRequest(cheaters, stateRoot, block, blockParticipatedCopy))
	for _, appEvent := range resp.Events {
		switch appEvent.Type {
		case "epoch sealed":
			sealEpoch = true
		}
	}

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
	s.blockTxHashes = append(s.blockTxHashes, block.TxHash)

	// process new epoch
	if sealEpoch {
		// prune not needed gas power records
		s.store.DelGasPowerRefunds(epoch - 1)
		s.onEpochSealed(block, cheaters)
		blockTxHashes, s.blockTxHashes = s.blockTxHashes, nil
	}

	log.Info("New block",
		"index", block.Index,
		"atropos", block.Atropos,
		"gasUsed", block.GasUsed,
		"skipped_txs", len(block.SkippedTxs),
		"txs", len(okTxs),
		"t", time.Since(start))

	s.store.SetBlock(block)
	s.store.SetBlockIndex(block.Atropos, block.Index)

	// Build index for txs
	if s.config.TxIndex {
		var i uint32
		for txHash, txPos := range txsPositions {
			if txPos.Block <= 0 {
				continue
			}
			// not skipped txs only
			txPos.BlockOffset = i
			i++
			s.store.SetTxPosition(txHash, txPos)
		}
	}

	// Trace by which event this block was confirmed (only for API)
	if s.config.DecisiveEventsIndex {
		s.store.SetBlockDecidedBy(block.Index, s.currentEvent)
	}

	evmHeader := evmcore.ToEvmHeader(block)
	s.feed.newBlock.Send(evmcore.ChainHeadNotify{
		Block: &evmcore.EvmBlock{
			EvmHeader:    *evmHeader,
			Transactions: okTxs,
		}})

	// trace confirmed transactions
	confirmTxnsMeter.Inc(int64(len(txsPositions)))
	for tx := range txsPositions {
		tracing.FinishTx(tx, "Service.onNewBlock()")
		if latency, err := txLatency.Finish(tx); err == nil {
			txTtfMeter.Update(latency.Milliseconds())
		}
	}

	return
}
