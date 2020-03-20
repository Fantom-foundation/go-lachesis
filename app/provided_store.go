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

// HasEpochValidator provides store's method.
func (a App) HasEpochValidator(epoch idx.Epoch, stakerID idx.StakerID) bool {
	return a.store.HasEpochValidator(epoch, stakerID)
}

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
