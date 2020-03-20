package app

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/sfctype"
)

/*
 * NOTE: all the methods are temporary and will be refactored during Tendermint implementation.
 */

// GetEpochValidators provides store's method.
func (a App) GetEpochValidators(epoch idx.Epoch) []sfctype.SfcStakerAndID {
	return a.store.GetEpochValidators(epoch)
}

// GetSfcConstants provides store's method.
func (a App) GetSfcConstants(epoch idx.Epoch) SfcConstants {
	return a.store.GetSfcConstants(epoch)
}

// SetReceipts provides store's method.
func (a App) SetReceipts(n idx.Block, receipts types.Receipts) {
	a.store.SetReceipts(n, receipts)
}

// StateDB provides store's method.
func (a App) StateDB(from common.Hash) *state.StateDB {
	return a.store.StateDB(from)
}

// HasSfcStaker provides store's method.
func (a App) HasSfcStaker(stakerID idx.StakerID) bool {
	return a.store.HasSfcStaker(stakerID)
}

// HasEpochValidator provides store's method.
func (a App) HasEpochValidator(epoch idx.Epoch, stakerID idx.StakerID) bool {
	return a.store.HasEpochValidator(epoch, stakerID)
}

// GetSfcStaker provides store's method.
func (a App) GetSfcStaker(stakerID idx.StakerID) *sfctype.SfcStaker {
	return a.store.GetSfcStaker(stakerID)
}

// ForEachSfcStaker provides store's method.
func (a App) ForEachSfcStaker(do func(sfctype.SfcStakerAndID)) {
	a.store.ForEachSfcStaker(do)
}

// ForEachSfcDelegator provides store's method.
func (a App) ForEachSfcDelegator(do func(sfctype.SfcDelegatorAndAddr)) {
	a.store.ForEachSfcDelegator(do)
}
