package app

import (
	"math/big"

	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
)

// TxPosition includes block and event.
type TxPosition struct {
	Block       idx.Block
	Event       hash.Event
	Creator     idx.StakerID
	EventOffset uint32
	BlockOffset uint32
}

// BlocksMissed is information about missed blocks from a staker.
type BlocksMissed struct {
	Num    idx.Block
	Period inter.Timestamp
}

// updateOriginationScores calculates the origination scores
func (a *App) updateOriginationScores() {
	if !a.ctx.sealEpoch {
		return
	}
	a.store.DelAllActiveOriginationScores()
	a.store.MoveDirtyOriginationScoresToActive()
}

// updateValidationScores calculates the validation scores
func (a *App) updateValidationScores(
	epoch idx.Epoch,
	blockN idx.Block,
	blockParticipated map[idx.StakerID]bool,
) {
	blockTimeDiff := a.blockTime(blockN) - a.blockTime(blockN-1)

	// Calc validation scores
	for _, it := range a.store.GetActiveSfcStakers() {
		// validators only
		if !a.store.HasEpochValidator(epoch, it.StakerID) {
			continue
		}

		// Check if validator has confirmed events by this Atropos
		missedBlock := !blockParticipated[it.StakerID]

		// If have no confirmed events by this Atropos - just add missed blocks for validator
		if missedBlock {
			a.store.IncBlocksMissed(it.StakerID, blockTimeDiff)
			continue
		}

		missedNum := a.store.GetBlocksMissed(it.StakerID).Num
		if missedNum > a.config.Net.Economy.BlockMissedLatency {
			missedNum = a.config.Net.Economy.BlockMissedLatency
		}

		// Add score for previous blocks, but no more than FrameLatency prev blocks
		a.store.AddDirtyValidationScore(it.StakerID, new(big.Int).SetUint64(uint64(blockTimeDiff)))
		for i := idx.Block(1); i <= missedNum && i < blockN; i++ {
			curBlockTime := a.blockTime(blockN - i)
			prevBlockTime := a.blockTime(blockN - i - 1)
			timeDiff := curBlockTime - prevBlockTime
			a.store.AddDirtyValidationScore(it.StakerID, new(big.Int).SetUint64(uint64(timeDiff)))
		}
		a.store.ResetBlocksMissed(it.StakerID)
	}

	if !a.ctx.sealEpoch {
		return
	}

	if a.config.EpochDowntimeIndex {
		a.store.NewDowntimeSnapshotEpoch(epoch)
	}
	if a.config.EpochActiveValidationScoreIndex {
		a.store.NewScoreSnapshotEpoch(epoch)
	}

	a.store.DelAllActiveValidationScores()
	a.store.MoveDirtyValidationScoresToActive()
}
