package app

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/lachesis"
)

// PoiPeriod calculate POI period from int64 unix time
func PoiPeriod(t inter.Timestamp, config *lachesis.EconomyConfig) uint64 {
	return uint64(t) / uint64(config.PoiPeriodDuration)
}

// UpdateAddressPOI calculate and save POI for user
func (a *Application) UpdateAddressPOI(address common.Address, senderTotalFee *big.Int, poiPeriod uint64) {
	/*if senderTotalFee.Sign() == 0 {
		s.store.SetAddressPOI(address, common.Big0)
		return // avoid division by 0
	}
	poi := new(big.Int).Mul(senderTotalFee, lachesis.PercentUnit)
	poi.Div(poi, s.store.GetPoiFee(poiPeriod)) // rebase user's PoI as <= 1.0 ratio
	s.store.SetAddressPOI(address, poi)*/
}

// updateUsersPOI calculates the Proof Of Importance weights for users
func (a *Application) updateUsersPOI(block *inter.Block, evmBlock *evmcore.EvmBlock, receipts types.Receipts, totalFee *big.Int, sealEpoch bool) {
	// User POI calculations
	poiPeriod := PoiPeriod(block.Time, &a.config.Economy)
	a.store.AddPoiFee(poiPeriod, totalFee)

	for i, tx := range evmBlock.Transactions {
		txFee := new(big.Int).Mul(new(big.Int).SetUint64(receipts[i].GasUsed), tx.GasPrice())

		signer := types.NewEIP155Signer(a.config.EvmChainConfig().ChainID)
		sender, err := signer.Sender(tx)
		if err != nil {
			a.Log.Crit("Failed to get sender from transaction", "err", err)
		}

		senderLastTxTime := a.store.GetAddressLastTxTime(sender)
		prevUserPoiPeriod := PoiPeriod(senderLastTxTime, &a.config.Economy)
		senderTotalFee := a.store.GetAddressFee(sender, prevUserPoiPeriod)

		delegator := a.store.GetSfcDelegator(sender)
		if delegator != nil {
			staker := a.store.GetSfcStaker(delegator.ToStakerID)
			prevWeightedTxFee := a.store.GetWeightedDelegatorsFee(delegator.ToStakerID)

			weightedTxFee := new(big.Int).Mul(txFee, delegator.Amount)
			weightedTxFee.Div(weightedTxFee, staker.CalcTotalStake())

			weightedTxFee.Add(weightedTxFee, prevWeightedTxFee)
			a.store.SetWeightedDelegatorsFee(delegator.ToStakerID, weightedTxFee)
		}

		if prevUserPoiPeriod != poiPeriod {
			a.UpdateAddressPOI(sender, senderTotalFee, prevUserPoiPeriod)
			senderTotalFee = big.NewInt(0)
		}

		a.store.SetAddressLastTxTime(sender, block.Time)
		senderTotalFee.Add(senderTotalFee, txFee)
		a.store.SetAddressFee(sender, poiPeriod, senderTotalFee)
	}

}

// UpdateStakerPOI calculate and save POI for staker
func (a *Application) UpdateStakerPOI(stakerID idx.StakerID, stakerAddress common.Address, poiPeriod uint64) {
	staker := a.store.GetSfcStaker(stakerID)

	vFee := a.store.GetAddressFee(stakerAddress, poiPeriod)
	weightedDFee := a.store.GetWeightedDelegatorsFee(stakerID)
	if vFee.Sign() == 0 && weightedDFee.Sign() == 0 {
		a.store.SetStakerPOI(stakerID, common.Big0)
		return // optimization
	}

	weightedVFee := new(big.Int).Mul(vFee, staker.StakeAmount)
	weightedVFee.Div(weightedVFee, staker.CalcTotalStake())

	weightedFee := new(big.Int).Add(weightedDFee, weightedVFee)

	if weightedFee.Sign() == 0 {
		a.store.SetStakerPOI(stakerID, common.Big0)
		return // avoid division by 0
	}
	poi := weightedFee // no need to rebase validator's PoI as <= 1.0 ratio
	/*poi := new(big.Int).Mul(weightedFee, lachesis.PercentUnit)
	poi.Div(poi, s.store.GetPoiFee(poiPeriod))*/
	a.store.SetStakerPOI(stakerID, poi)
}

// updateStakersPOI calculates the Proof Of Importance weights for stakers
func (a *Application) updateStakersPOI(prev, block *inter.Block, sealEpoch bool) {
	// Stakers POI calculations
	poiPeriod := PoiPeriod(block.Time, &a.config.Economy)
	prevBlockPoiPeriod := PoiPeriod(prev.Time, &a.config.Economy)

	if poiPeriod != prevBlockPoiPeriod {
		for _, it := range a.GetActiveSfcStakers() {
			a.UpdateStakerPOI(it.StakerID, it.Staker.Address, prevBlockPoiPeriod)
		}
		// clear StakersDelegatorsFee counters
		a.store.DelAllWeightedDelegatorsFee()
	}
}
