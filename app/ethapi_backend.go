package app

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"

	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/sfctype"
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

// GetValidationScore returns staker's ValidationScore.
func (b *EthAPIBackend) GetValidationScore(ctx context.Context, stakerID idx.StakerID) (*big.Int, error) {
	if !b.app.store.HasSfcStaker(stakerID) {
		return nil, nil
	}
	return b.app.store.GetActiveValidationScore(stakerID), nil
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

// GetDowntime returns staker's Downtime.
func (b *EthAPIBackend) GetDowntime(ctx context.Context, stakerID idx.StakerID) (idx.Block, inter.Timestamp, error) {
	missed := b.app.store.GetBlocksMissed(stakerID)
	return missed.Num, missed.Period, nil
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
