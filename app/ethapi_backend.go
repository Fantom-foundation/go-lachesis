package app

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/sfctype"
	"github.com/Fantom-foundation/go-lachesis/lachesis/genesis/sfc"
	"github.com/Fantom-foundation/go-lachesis/lachesis/genesis/sfc/sfcpos"
	"github.com/Fantom-foundation/go-lachesis/topicsdb"
)

// EthAPIBackend provides methods for ethapi.Backend
type EthAPIBackend struct {
	app *App
}

// EthAPIBackend getter
func (a *App) EthAPIBackend() *EthAPIBackend {
	return &EthAPIBackend{a}
}

func (b *EthAPIBackend) ChainDb() ethdb.Database {
	return b.app.store.EvmTable()
}

func (b *EthAPIBackend) EvmLogIndex() *topicsdb.Index {
	return b.app.store.EvmLogs()
}

func (b *EthAPIBackend) GetReceipts(n idx.Block) types.Receipts {
	return b.app.store.GetReceipts(n)
}

// GetValidationScore returns staker's ValidationScore for epoch.
// * When epoch is "pending" the validation score for latest epoch are returned.
// * When epoch is "latest" the validation score for latest sealed epoch are returned.
func (b *EthAPIBackend) GetValidationScore(
	ctx context.Context, stakerID idx.StakerID, epoch rpc.BlockNumber,
) (
	score *big.Int, err error,
) {
	idxEpoch, err := b.epochWithDefault(ctx, epoch)
	if err != nil {
		return
	}

	currentEpoch := b.app.GetEpoch()

	if idxEpoch != currentEpoch && !b.app.config.EpochActiveValidationScoreIndex {
		err = errors.New("pass 'pending' epoch if EpochActiveValidationScoreIndex is false")
		return
	}

	if idxEpoch == currentEpoch {
		score = b.app.store.GetActiveValidationScore(stakerID)
	} else {
		score = b.app.store.GetActiveValidationScoreEpoch(stakerID, idxEpoch)
	}

	return
}

// GetOriginationScore returns staker's OriginationScore.
func (b *EthAPIBackend) GetOriginationScore(ctx context.Context, stakerID idx.StakerID) (*big.Int, error) {
	if !b.app.store.HasSfcStaker(stakerID) {
		return nil, nil
	}
	return b.app.store.GetActiveOriginationScore(stakerID), nil
}

// GetStakerPoI returns staker's PoI.
func (b *EthAPIBackend) GetStakerPoI(ctx context.Context, stakerID idx.StakerID) (*big.Int, error) {
	if !b.app.store.HasSfcStaker(stakerID) {
		return nil, nil
	}
	return b.app.store.GetStakerPOI(stakerID), nil
}

// GetDowntime returns staker's Downtime for epoch.
// * When epoch is "pending" the downtime for latest epoch are returned.
// * When epoch is "latest" the downtime for latest sealed epoch are returned.
func (b *EthAPIBackend) GetDowntime(
	ctx context.Context, stakerID idx.StakerID, epoch rpc.BlockNumber,
) (
	num idx.Block, period inter.Timestamp, err error,
) {
	idxEpoch, err := b.epochWithDefault(ctx, epoch)
	if err != nil {
		return
	}

	currentEpoch := b.app.GetEpoch()

	if idxEpoch != currentEpoch && !b.app.config.EpochDowntimeIndex {
		err = errors.New("pass 'pending' epoch if EpochDowntimeIndex is false")
		return
	}

	var missed BlocksMissed
	if idxEpoch == currentEpoch {
		missed = b.app.store.GetBlocksMissed(stakerID)
	} else {
		missed = b.app.store.GetBlocksMissedEpoch(stakerID, idxEpoch)
	}

	num = missed.Num
	period = missed.Period
	return
}

// GetDelegator returns SFC delegator info
func (b *EthAPIBackend) GetDelegator(ctx context.Context, addr common.Address) (*sfctype.SfcDelegator, error) {
	return b.app.store.GetSfcDelegator(addr), nil
}

// GetDelegatorClaimedRewards returns sum of claimed rewards in past, by this delegator
func (b *EthAPIBackend) GetDelegatorClaimedRewards(ctx context.Context, addr common.Address) (*big.Int, error) {
	return b.app.store.GetDelegatorClaimedRewards(addr), nil
}

// GetStakerClaimedRewards returns sum of claimed rewards in past, by this staker
func (b *EthAPIBackend) GetStakerClaimedRewards(ctx context.Context, stakerID idx.StakerID) (*big.Int, error) {
	return b.app.store.GetStakerClaimedRewards(stakerID), nil
}

// GetStakerDelegatorsClaimedRewards returns sum of claimed rewards in past, by this delegators of this staker
func (b *EthAPIBackend) GetStakerDelegatorsClaimedRewards(ctx context.Context, stakerID idx.StakerID) (*big.Int, error) {
	return b.app.store.GetStakerDelegatorsClaimedRewards(stakerID), nil
}

// HasSfcStaker provides store's method.
func (b *EthAPIBackend) HasSfcStaker(stakerID idx.StakerID) bool {
	return b.app.store.HasSfcStaker(stakerID)
}

// GetSfcStaker provides store's method.
func (b *EthAPIBackend) GetSfcStaker(stakerID idx.StakerID) *sfctype.SfcStaker {
	return b.app.store.GetSfcStaker(stakerID)
}

// ForEachSfcStaker provides store's method.
func (b *EthAPIBackend) ForEachSfcStaker(do func(sfctype.SfcStakerAndID)) {
	b.app.store.ForEachSfcStaker(do)
}

// ForEachSfcDelegator provides store's method.
func (b *EthAPIBackend) ForEachSfcDelegator(do func(sfctype.SfcDelegatorAndAddr)) {
	b.app.store.ForEachSfcDelegator(do)
}

// GetEpochStats returns epoch statistics.
// * When epoch is "pending" the statistics for latest epoch are returned.
// * When epoch is "latest" the statistics for latest sealed epoch are returned.
func (b *EthAPIBackend) GetEpochStats(ctx context.Context, requestedEpoch rpc.BlockNumber) (*sfctype.EpochStats, error) {
	epoch, err := b.epochWithDefault(ctx, requestedEpoch)
	if err != nil {
		return nil, err
	}

	stats := b.app.store.GetEpochStats(epoch)
	if stats == nil {
		return nil, nil
	}
	stats.Epoch = epoch

	// read total reward weights from SFC contract
	root := b.app.LastBlock().Root
	statedb := b.app.StateDB(root)

	epochPosition := sfcpos.EpochSnapshot(epoch)
	stats.TotalBaseRewardWeight = statedb.GetState(sfc.ContractAddress, epochPosition.TotalBaseRewardWeight()).Big()
	stats.TotalTxRewardWeight = statedb.GetState(sfc.ContractAddress, epochPosition.TotalTxRewardWeight()).Big()

	return stats, nil
}

func (b *EthAPIBackend) epochWithDefault(ctx context.Context, epoch rpc.BlockNumber) (requested idx.Epoch, err error) {
	current := b.app.GetEpoch()

	switch {
	case epoch == rpc.PendingBlockNumber:
		requested = current
	case epoch == rpc.LatestBlockNumber:
		requested = current - 1
	case epoch >= 0 && idx.Epoch(epoch) <= current:
		requested = idx.Epoch(epoch)
	default:
		err = errors.New("epoch is not in range")
		return
	}
	return requested, nil
}
