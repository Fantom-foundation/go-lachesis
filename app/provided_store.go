package app

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"

	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/sfctype"
	"github.com/Fantom-foundation/go-lachesis/topicsdb"
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

//  provides store's method.
func (a App) GetReceipts(n idx.Block) types.Receipts {
	return a.store.GetReceipts(n)
}

// StateDB provides store's method.
func (a App) StateDB(from common.Hash) *state.StateDB {
	return a.store.StateDB(from)
}

// EvmTable provides store's method.
func (a App) EvmTable() ethdb.Database {
	return a.store.EvmTable()
}

//  provides store's method.
func (a App) EvmLogs() *topicsdb.Index {
	return a.store.EvmLogs()
}

// HasSfcStaker provides store's method.
func (a App) HasSfcStaker(stakerID idx.StakerID) bool {
	return a.store.HasSfcStaker(stakerID)
}

// HasEpochValidator provides store's method.
func (a App) HasEpochValidator(epoch idx.Epoch, stakerID idx.StakerID) bool {
	return a.store.HasEpochValidator(epoch, stakerID)
}

//  provides store's method.
func (a App) GetActiveValidationScore(stakerID idx.StakerID) *big.Int {
	return a.store.GetActiveValidationScore(stakerID)
}

// GetActiveOriginationScore provides store's method.
func (a App) GetActiveOriginationScore(stakerID idx.StakerID) *big.Int {
	return a.store.GetActiveOriginationScore(stakerID)
}

//  provides store's method.
func (a App) GetStakerPOI(stakerID idx.StakerID) *big.Int {
	return a.store.GetStakerPOI(stakerID)
}

// SetStakerPOI provides store's method.
func (a App) SetStakerPOI(stakerID idx.StakerID, poi *big.Int) {
	a.store.SetStakerPOI(stakerID, poi)
}

// GetBlocksMissed provides store's method.
func (a App) GetBlocksMissed(stakerID idx.StakerID) BlocksMissed {
	return a.store.GetBlocksMissed(stakerID)
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

// GetSfcDelegator provides store's method.
func (a App) GetSfcDelegator(addr common.Address) *sfctype.SfcDelegator {
	return a.store.GetSfcDelegator(addr)
}

// GetDelegatorClaimedRewards provides store's method.
func (a App) GetDelegatorClaimedRewards(addr common.Address) *big.Int {
	return a.store.GetDelegatorClaimedRewards(addr)
}

// GetStakerClaimedRewards provides store's method.
func (a App) GetStakerClaimedRewards(stakerID idx.StakerID) *big.Int {
	return a.store.GetStakerClaimedRewards(stakerID)
}

// GetStakerDelegatorsClaimedRewards provides store's method.
func (a App) GetStakerDelegatorsClaimedRewards(stakerID idx.StakerID) *big.Int {
	return a.store.GetStakerDelegatorsClaimedRewards(stakerID)
}

// POI

// AddPoiFee provides store's method.
func (a App) AddPoiFee(poiPeriod uint64, diff *big.Int) {
	a.store.AddPoiFee(poiPeriod, diff)
}

// GetAddressLastTxTime provides store's method.
func (a App) GetAddressLastTxTime(addr common.Address) inter.Timestamp {
	return a.store.GetAddressLastTxTime(addr)
}

// SetAddressLastTxTime provides store's method.
func (a App) SetAddressLastTxTime(addr common.Address, t inter.Timestamp) {
	a.store.SetAddressLastTxTime(addr, t)
}

// GetAddressFee provides store's method.
func (a App) GetAddressFee(addr common.Address, poiPeriod uint64) *big.Int {
	return a.store.GetAddressFee(addr, poiPeriod)
}

// SetAddressFee provides store's method.
func (a App) SetAddressFee(addr common.Address, poiPeriod uint64, val *big.Int) {
	a.store.SetAddressFee(addr, poiPeriod, val)
}

//  provides store's method.
func (a App) GetWeightedDelegatorsFee(stakerID idx.StakerID) *big.Int {
	return a.store.GetWeightedDelegatorsFee(stakerID)
}

//  provides store's method.
func (a App) SetWeightedDelegatorsFee(stakerID idx.StakerID, val *big.Int) {
	a.store.SetWeightedDelegatorsFee(stakerID, val)
}

// GetActiveSfcStakers provides store's method.
func (a App) GetActiveSfcStakers() []sfctype.SfcStakerAndID {
	return a.store.GetActiveSfcStakers()
}

// DelAllWeightedDelegatorsFee provides store's method.
func (a App) DelAllWeightedDelegatorsFee() {
	a.store.DelAllWeightedDelegatorsFee()
}
