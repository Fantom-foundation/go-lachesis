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

func (a App) GetEpochValidators(epoch idx.Epoch) []sfctype.SfcStakerAndID {
	return a.store.GetEpochValidators(epoch)
}

func (a App) GetSfcConstants(epoch idx.Epoch) SfcConstants {
	return a.store.GetSfcConstants(epoch)
}

func (a App) SetReceipts(n idx.Block, receipts types.Receipts) {
	a.store.SetReceipts(n, receipts)
}

func (a App) GetReceipts(n idx.Block) types.Receipts {
	return a.store.GetReceipts(n)
}

func (a App) StateDB(from common.Hash) *state.StateDB {
	return a.store.StateDB(from)
}

func (a App) EvmTable() ethdb.Database {
	return a.store.EvmTable()
}

func (a App) EvmLogs() *topicsdb.Index {
	return a.store.EvmLogs()
}

func (a App) HasSfcStaker(stakerID idx.StakerID) bool {
	return a.store.HasSfcStaker(stakerID)
}

func (a App) HasEpochValidator(epoch idx.Epoch, stakerID idx.StakerID) bool {
	return a.store.HasEpochValidator(epoch, stakerID)
}

func (a App) GetActiveValidationScore(stakerID idx.StakerID) *big.Int {
	return a.store.GetActiveValidationScore(stakerID)
}

func (a App) GetActiveOriginationScore(stakerID idx.StakerID) *big.Int {
	return a.store.GetActiveOriginationScore(stakerID)
}

func (a App) GetStakerPOI(stakerID idx.StakerID) *big.Int {
	return a.store.GetStakerPOI(stakerID)
}

func (a App) SetStakerPOI(stakerID idx.StakerID, poi *big.Int) {
	a.store.SetStakerPOI(stakerID, poi)
}

func (a App) GetBlocksMissed(stakerID idx.StakerID) BlocksMissed {
	return a.store.GetBlocksMissed(stakerID)
}

func (a App) GetSfcStaker(stakerID idx.StakerID) *sfctype.SfcStaker {
	return a.store.GetSfcStaker(stakerID)
}

func (a App) ForEachSfcStaker(do func(sfctype.SfcStakerAndID)) {
	a.store.ForEachSfcStaker(do)
}

func (a App) ForEachSfcDelegator(do func(sfctype.SfcDelegatorAndAddr)) {
	a.store.ForEachSfcDelegator(do)
}

func (a App) GetSfcDelegator(addr common.Address) *sfctype.SfcDelegator {
	return a.store.GetSfcDelegator(addr)
}

func (a App) GetDelegatorClaimedRewards(addr common.Address) *big.Int {
	return a.store.GetDelegatorClaimedRewards(addr)
}

func (a App) GetStakerClaimedRewards(stakerID idx.StakerID) *big.Int {
	return a.store.GetStakerClaimedRewards(stakerID)
}

func (a App) GetStakerDelegatorsClaimedRewards(stakerID idx.StakerID) *big.Int {
	return a.store.GetStakerDelegatorsClaimedRewards(stakerID)
}

// POI

func (a App) AddPoiFee(poiPeriod uint64, diff *big.Int) {
	a.store.AddPoiFee(poiPeriod, diff)
}

func (a App) GetAddressLastTxTime(addr common.Address) inter.Timestamp {
	return a.store.GetAddressLastTxTime(addr)
}

func (a App) SetAddressLastTxTime(addr common.Address, t inter.Timestamp) {
	a.store.SetAddressLastTxTime(addr, t)
}

func (a App) GetAddressFee(addr common.Address, poiPeriod uint64) *big.Int {
	return a.store.GetAddressFee(addr, poiPeriod)
}

func (a App) SetAddressFee(addr common.Address, poiPeriod uint64, val *big.Int) {
	a.store.SetAddressFee(addr, poiPeriod, val)
}

func (a App) GetWeightedDelegatorsFee(stakerID idx.StakerID) *big.Int {
	return a.store.GetWeightedDelegatorsFee(stakerID)
}

func (a App) SetWeightedDelegatorsFee(stakerID idx.StakerID, val *big.Int) {
	a.store.SetWeightedDelegatorsFee(stakerID, val)
}

func (a App) GetActiveSfcStakers() []sfctype.SfcStakerAndID {
	return a.store.GetActiveSfcStakers()
}

func (a App) DelAllWeightedDelegatorsFee() {
	a.store.DelAllWeightedDelegatorsFee()
}

// scores

func (a App) AddDirtyOriginationScore(stakerID idx.StakerID, v *big.Int) {
	a.store.AddDirtyOriginationScore(stakerID, v)
}

func (a App) DelActiveOriginationScore(stakerID idx.StakerID) {
	a.store.DelActiveOriginationScore(stakerID)
}

func (a App) DelAllActiveOriginationScores() {
	a.store.DelAllActiveOriginationScores()
}

func (a App) AddDirtyValidationScore(stakerID idx.StakerID, v *big.Int) {
	a.store.AddDirtyValidationScore(stakerID, v)
}

func (a App) ResetBlocksMissed(stakerID idx.StakerID) {
	a.store.ResetBlocksMissed(stakerID)
}

func (a App) MoveDirtyOriginationScoresToActive() {
	a.store.MoveDirtyOriginationScoresToActive()
}

func (a App) IncBlocksMissed(stakerID idx.StakerID, periodDiff inter.Timestamp) {
	a.store.IncBlocksMissed(stakerID, periodDiff)
}

func (a App) DelActiveValidationScore(stakerID idx.StakerID) {
	a.store.DelActiveValidationScore(stakerID)
}

func (a App) DelAllActiveValidationScores() {
	a.store.DelAllActiveValidationScores()
}

func (a App) MoveDirtyValidationScoresToActive() {
	a.store.MoveDirtyValidationScoresToActive()
}
