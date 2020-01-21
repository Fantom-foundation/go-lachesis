package app

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"

	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/pos"
	"github.com/Fantom-foundation/go-lachesis/inter/sfctype"
	"github.com/Fantom-foundation/go-lachesis/topicsdb"
)

// TODO: refactor to remove bad dependences.

// Gossip dependences.
type gossip interface {
	GetDummyChainReader() evmcore.DummyChain

	GetBlock(idx.Block) *inter.Block
	GetValidators() *pos.Validators
	GetEventHeader(epoch idx.Epoch, h hash.Event) *inter.EventHeaderData
	BlockParticipated(idx.StakerID) bool

	GetDirtyEpochStats() *sfctype.EpochStats
	SetDirtyEpochStats(val *sfctype.EpochStats)
	GetEpochStats(epoch idx.Epoch) *sfctype.EpochStats
	SetEpochStats(epoch idx.Epoch, val *sfctype.EpochStats)
}

func (a *Application) ForEachSfcDelegator(do func(sfctype.SfcDelegatorAndAddr)) {
	a.store.ForEachSfcDelegator(do)
}

func (a *Application) ForEachSfcStaker(do func(sfctype.SfcStakerAndID)) {
	a.store.ForEachSfcStaker(do)
}

func (a *Application) GetBlocksMissed(stakerID idx.StakerID) BlocksMissed {
	return a.store.GetBlocksMissed(stakerID)
}

func (a *Application) GetSfcConstants(epoch idx.Epoch) SfcConstants {
	return a.store.GetSfcConstants(epoch)
}

func (a *Application) GetGasPowerRefunds(epoch idx.Epoch) map[idx.StakerID]uint64 {
	return a.store.GetGasPowerRefunds(epoch)
}

func (a *Application) GetEpochValidators(epoch idx.Epoch) []sfctype.SfcStakerAndID {
	return a.store.GetEpochValidators(epoch)
}

func (a *Application) TableEvm() ethdb.Database {
	return a.store.table.Evm
}

func (a *Application) TableEvmLogs() *topicsdb.Index {
	return a.store.table.EvmLogs
}

func (a *Application) HasSfcStaker(stakerID idx.StakerID) bool {
	return a.store.HasSfcStaker(stakerID)
}

func (a *Application) GetActiveValidationScore(stakerID idx.StakerID) *big.Int {
	return a.store.GetActiveValidationScore(stakerID)
}

func (a *Application) GetActiveOriginationScore(stakerID idx.StakerID) *big.Int {
	return a.store.GetActiveOriginationScore(stakerID)
}

func (a *Application) GetSfcStaker(stakerID idx.StakerID) *sfctype.SfcStaker {
	return a.store.GetSfcStaker(stakerID)
}

func (a *Application) GetStakerPOI(stakerID idx.StakerID) *big.Int {
	return a.store.GetStakerPOI(stakerID)
}

func (a *Application) GetDelegatorClaimedRewards(addr common.Address) *big.Int {
	return a.store.GetDelegatorClaimedRewards(addr)
}

func (a *Application) GetStakerClaimedRewards(stakerID idx.StakerID) *big.Int {
	return a.store.GetStakerClaimedRewards(stakerID)
}

func (a *Application) GetStakerDelegatorsClaimedRewards(stakerID idx.StakerID) *big.Int {
	return a.store.GetStakerDelegatorsClaimedRewards(stakerID)
}

func (a *Application) GetSfcDelegator(addr common.Address) *sfctype.SfcDelegator {
	return a.store.GetSfcDelegator(addr)
}
