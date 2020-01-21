package app

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
)

// BlocksMissed is information about missed blocks from a staker
type BlocksMissed struct {
	Num    idx.Block
	Period inter.Timestamp
}

const (
	minGasPowerRefund = 800
)

// updateOriginationScores calculates the origination scores
func (a *Application) updateOriginationScores(block *inter.Block, evmBlock *evmcore.EvmBlock, receipts types.Receipts, txPositions map[common.Hash]TxPosition, epoch idx.Epoch, sealEpoch bool) {
	// Calc origination scores
	for i, tx := range evmBlock.Transactions {
		txEventPos := txPositions[receipts[i].TxHash]
		// sanity check
		if txEventPos.Block != block.Index {
			a.Log.Crit("Incorrect tx block position", "tx", receipts[i].TxHash,
				"block", txEventPos.Block, "block_got", block.Index)
		}

		txEvent := a.Gossip.GetEventHeader(txEventPos.Event.Epoch(), txEventPos.Event)
		// sanity check
		if txEvent == nil {
			a.Log.Crit("Incorrect tx event position", "tx", receipts[i].TxHash, "event", txEventPos.Event, "reason", "event has no transactions")
		}

		txFee := new(big.Int).Mul(new(big.Int).SetUint64(receipts[i].GasUsed), tx.GasPrice())

		a.store.AddDirtyOriginationScore(txEvent.Creator, txFee)

		{ // logic for gas power refunds
			if tx.Gas() < receipts[i].GasUsed {
				a.Log.Crit("Transaction gas used is higher than tx gas limit", "tx", receipts[i].TxHash, "event", txEventPos.Event)
			}
			notUsedGas := tx.Gas() - receipts[i].GasUsed
			if notUsedGas >= minGasPowerRefund { // do not refund if refunding is more costly than refunded value
				a.store.IncGasPowerRefund(epoch, txEvent.Creator, notUsedGas)
			}
		}
	}

	if sealEpoch {
		a.store.DelAllActiveOriginationScores()
		a.store.MoveDirtyOriginationScoresToActive()
		// prune not needed gas power records
		a.store.DelGasPowerRefunds(epoch - 1)
	}
}

// updateValidationScores calculates the validation scores
func (a *Application) updateValidationScores(prev, block *inter.Block, sealEpoch bool) {
	blockTimeDiff := block.Time - prev.Time

	// Calc validation scores
	for _, it := range a.GetActiveSfcStakers() {
		// validators only
		if !a.Gossip.GetValidators().Exists(it.StakerID) {
			continue
		}

		// Check if validator has confirmed events by this Atropos
		missedBlock := !a.Gossip.BlockParticipated(it.StakerID)

		// If have no confirmed events by this Atropos - just add missed blocks for validator
		if missedBlock {
			a.store.IncBlocksMissed(it.StakerID, blockTimeDiff)
			continue
		}

		missedNum := a.store.GetBlocksMissed(it.StakerID).Num
		if missedNum > a.config.Economy.BlockMissedLatency {
			missedNum = a.config.Economy.BlockMissedLatency
		}

		// Add score for previous blocks, but no more than FrameLatency prev blocks
		a.store.AddDirtyValidationScore(it.StakerID, new(big.Int).SetUint64(uint64(blockTimeDiff)))
		for i := idx.Block(1); i <= missedNum && i < block.Index; i++ {
			blockTime := a.Gossip.GetBlock(block.Index - i).Time
			prevBlockTime := a.Gossip.GetBlock(block.Index - i - 1).Time
			timeDiff := blockTime - prevBlockTime
			a.store.AddDirtyValidationScore(it.StakerID, new(big.Int).SetUint64(uint64(timeDiff)))
		}
		a.store.ResetBlocksMissed(it.StakerID)
	}

	if sealEpoch {
		a.store.DelAllActiveValidationScores()
		a.store.MoveDirtyValidationScoresToActive()
	}
}
