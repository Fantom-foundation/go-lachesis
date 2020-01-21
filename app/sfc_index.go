package app

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/sfctype"
	"github.com/Fantom-foundation/go-lachesis/lachesis"
	"github.com/Fantom-foundation/go-lachesis/lachesis/genesis/sfc"
	"github.com/Fantom-foundation/go-lachesis/lachesis/genesis/sfc/sfcpos"
	"github.com/Fantom-foundation/go-lachesis/utils"
)

// SfcConstants are constants which may be changed by SFC contract
type SfcConstants struct {
	ShortGasPowerAllocPerSec uint64
	LongGasPowerAllocPerSec  uint64
	BaseRewardPerSec         *big.Int
}

// GetActiveSfcStakers returns stakers which will become validators in next epoch
func (a *Application) GetActiveSfcStakers() []sfctype.SfcStakerAndID {
	stakers := make([]sfctype.SfcStakerAndID, 0, 200)
	a.store.ForEachSfcStaker(func(it sfctype.SfcStakerAndID) {
		if it.Staker.Ok() {
			stakers = append(stakers, it)
		}
	})
	return stakers
}

func (a *Application) delAllStakerData(stakerID idx.StakerID) {
	a.store.DelSfcStaker(stakerID)
	a.store.ResetBlocksMissed(stakerID)
	a.store.DelActiveValidationScore(stakerID)
	a.store.DelDirtyValidationScore(stakerID)
	a.store.DelActiveOriginationScore(stakerID)
	a.store.DelDirtyOriginationScore(stakerID)
	a.store.DelWeightedDelegatorsFee(stakerID)
	a.store.DelStakerPOI(stakerID)
	a.store.DelStakerClaimedRewards(stakerID)
	a.store.DelStakerDelegatorsClaimedRewards(stakerID)
}

func (a *Application) delAllDelegatorData(address common.Address) {
	a.store.DelSfcDelegator(address)
	a.store.DelDelegatorClaimedRewards(address)
}

var (
	max128 = new(big.Int).Sub(math.BigPow(2, 128), common.Big1)
)

func (a *Application) calcRewardWeights(stakers []sfctype.SfcStakerAndID, _epochDuration inter.Timestamp) (baseRewardWeights []*big.Int, txRewardWeights []*big.Int) {
	validationScores := make([]*big.Int, 0, len(stakers))
	originationScores := make([]*big.Int, 0, len(stakers))
	pois := make([]*big.Int, 0, len(stakers))
	stakes := make([]*big.Int, 0, len(stakers))

	if _epochDuration == 0 {
		_epochDuration = 1
	}
	epochDuration := new(big.Int).SetUint64(uint64(_epochDuration))

	for _, it := range stakers {
		stake := it.Staker.CalcTotalStake()
		poi := a.store.GetStakerPOI(it.StakerID)
		validationScore := a.store.GetActiveValidationScore(it.StakerID)
		originationScore := a.store.GetActiveOriginationScore(it.StakerID)

		stakes = append(stakes, stake)
		validationScores = append(validationScores, validationScore)
		originationScores = append(originationScores, originationScore)
		pois = append(pois, poi)
	}

	txRewardWeights = make([]*big.Int, 0, len(stakers))
	for i := range stakers {
		// txRewardWeight = ({origination score} + {CONST} * {PoI}) * {validation score}
		// origination score is roughly proportional to {validation score} * {stake}, so the whole formula is roughly
		// {stake} * {validation score} ^ 2
		poiWithRatio := new(big.Int).Mul(pois[i], a.config.Economy.TxRewardPoiImpact)
		poiWithRatio.Div(poiWithRatio, lachesis.PercentUnit)

		txRewardWeight := new(big.Int).Add(originationScores[i], poiWithRatio)
		txRewardWeight.Mul(txRewardWeight, validationScores[i])
		txRewardWeight.Div(txRewardWeight, epochDuration)
		if txRewardWeight.Cmp(max128) > 0 {
			txRewardWeight = new(big.Int).Set(max128) // never going to get here
		}

		txRewardWeights = append(txRewardWeights, txRewardWeight)
	}

	baseRewardWeights = make([]*big.Int, 0, len(stakers))
	for i := range stakers {
		// baseRewardWeight = {stake} * {validationScore ^ 2}
		baseRewardWeight := new(big.Int).Set(stakes[i])
		for pow := 0; pow < 2; pow++ {
			baseRewardWeight.Mul(baseRewardWeight, validationScores[i])
			baseRewardWeight.Div(baseRewardWeight, epochDuration)
		}
		if baseRewardWeight.Cmp(max128) > 0 {
			baseRewardWeight = new(big.Int).Set(max128) // never going to get here
		}

		baseRewardWeights = append(baseRewardWeights, baseRewardWeight)
	}

	return baseRewardWeights, txRewardWeights
}

// getRewardPerSec returns current rewardPerSec, depending on config and value provided by SFC
func (a *Application) getRewardPerSec(epoch idx.Epoch) *big.Int {
	rewardPerSecond := a.store.GetSfcConstants(epoch - 1).BaseRewardPerSec
	if rewardPerSecond == nil || rewardPerSecond.Sign() == 0 {
		rewardPerSecond = a.config.Economy.InitialRewardPerSecond
	}
	if rewardPerSecond.Cmp(a.config.Economy.MaxRewardPerSecond) > 0 {
		rewardPerSecond = a.config.Economy.MaxRewardPerSecond
	}
	return new(big.Int).Set(rewardPerSecond)
}

// processSfc applies the new SFC state
func (a *Application) processSfc(block *inter.Block, receipts types.Receipts, blockFee *big.Int, epoch idx.Epoch, sealEpoch bool, cheaters inter.Cheaters, statedb *state.StateDB) {
	// a.engineMu is locked here

	// process SFC contract logs
	for _, receipt := range receipts {
		for _, l := range receipt.Logs {
			if l.Address != sfc.ContractAddress {
				continue
			}
			// Add new stakers
			if l.Topics[0] == sfcpos.Topics.CreatedStake && len(l.Topics) > 2 && len(l.Data) >= 32 {
				stakerID := idx.StakerID(new(big.Int).SetBytes(l.Topics[1][:]).Uint64())
				address := common.BytesToAddress(l.Topics[2][12:])
				amount := new(big.Int).SetBytes(l.Data[0:32])

				a.store.SetSfcStaker(stakerID, &sfctype.SfcStaker{
					Address:      address,
					CreatedEpoch: epoch,
					CreatedTime:  block.Time,
					StakeAmount:  amount,
					DelegatedMe:  big.NewInt(0),
				})
			}

			// Increase stakes
			if l.Topics[0] == sfcpos.Topics.IncreasedStake && len(l.Topics) > 1 && len(l.Data) >= 32 {
				stakerID := idx.StakerID(new(big.Int).SetBytes(l.Topics[1][:]).Uint64())
				newAmount := new(big.Int).SetBytes(l.Data[0:32])

				staker := a.store.GetSfcStaker(stakerID)
				if staker == nil {
					a.Log.Error("Internal SFC index isn't synced with SFC contract")
					continue
				}
				staker.StakeAmount = newAmount
				a.store.SetSfcStaker(stakerID, staker)
			}

			// Add new delegators
			if l.Topics[0] == sfcpos.Topics.CreatedDelegation && len(l.Topics) > 1 && len(l.Data) >= 32 {
				address := common.BytesToAddress(l.Topics[1][12:])
				toStakerID := idx.StakerID(new(big.Int).SetBytes(l.Topics[2][12:]).Uint64())
				amount := new(big.Int).SetBytes(l.Data[0:32])

				staker := a.store.GetSfcStaker(toStakerID)
				if staker == nil {
					a.Log.Error("Internal SFC index isn't synced with SFC contract")
					continue
				}
				staker.DelegatedMe.Add(staker.DelegatedMe, amount)

				a.store.SetSfcDelegator(address, &sfctype.SfcDelegator{
					ToStakerID:   toStakerID,
					CreatedEpoch: epoch,
					CreatedTime:  block.Time,
					Amount:       amount,
				})
				a.store.SetSfcStaker(toStakerID, staker)
			}

			// Deactivate stakes
			if l.Topics[0] == sfcpos.Topics.PreparedToWithdrawStake && len(l.Topics) > 1 {
				stakerID := idx.StakerID(new(big.Int).SetBytes(l.Topics[1][:]).Uint64())

				staker := a.store.GetSfcStaker(stakerID)
				staker.DeactivatedEpoch = epoch
				staker.DeactivatedTime = block.Time
				a.store.SetSfcStaker(stakerID, staker)
			}

			// Deactivate delegators
			if l.Topics[0] == sfcpos.Topics.PreparedToWithdrawDelegation && len(l.Topics) > 1 {
				address := common.BytesToAddress(l.Topics[1][12:])

				delegator := a.store.GetSfcDelegator(address)
				staker := a.store.GetSfcStaker(delegator.ToStakerID)
				if staker != nil {
					staker.DelegatedMe.Sub(staker.DelegatedMe, delegator.Amount)
					a.store.SetSfcStaker(delegator.ToStakerID, staker)
				}
				delegator.DeactivatedEpoch = epoch
				delegator.DeactivatedTime = block.Time
				a.store.SetSfcDelegator(address, delegator)
			}

			// Delete stakes
			if l.Topics[0] == sfcpos.Topics.WithdrawnStake && len(l.Topics) > 1 {
				stakerID := idx.StakerID(new(big.Int).SetBytes(l.Topics[1][:]).Uint64())
				a.delAllStakerData(stakerID)
			}

			// Delete delegators
			if l.Topics[0] == sfcpos.Topics.WithdrawnDelegation && len(l.Topics) > 1 {
				address := common.BytesToAddress(l.Topics[1][12:])
				a.delAllDelegatorData(address)
			}

			// Track changes of constants by SFC
			if l.Topics[0] == sfcpos.Topics.UpdatedBaseRewardPerSec && len(l.Data) >= 32 {
				baseRewardPerSec := new(big.Int).SetBytes(l.Data[0:32])
				constants := a.store.GetSfcConstants(epoch)
				constants.BaseRewardPerSec = baseRewardPerSec
				a.store.SetSfcConstants(epoch, constants)
			}
			if l.Topics[0] == sfcpos.Topics.UpdatedGasPowerAllocationRate && len(l.Data) >= 64 {
				shortAllocationRate := new(big.Int).SetBytes(l.Data[0:32])
				longAllocationRate := new(big.Int).SetBytes(l.Data[32:64])
				constants := a.store.GetSfcConstants(epoch)
				constants.ShortGasPowerAllocPerSec = shortAllocationRate.Uint64()
				constants.LongGasPowerAllocPerSec = longAllocationRate.Uint64()
				a.store.SetSfcConstants(epoch, constants)
			}

			// Track rewards (API-only)
			if l.Topics[0] == sfcpos.Topics.ClaimedValidatorReward && len(l.Topics) > 1 && len(l.Data) >= 32 {
				stakerID := idx.StakerID(new(big.Int).SetBytes(l.Topics[1][:]).Uint64())
				reward := new(big.Int).SetBytes(l.Data[0:32])

				a.store.IncStakerClaimedRewards(stakerID, reward)
			}
			if l.Topics[0] == sfcpos.Topics.ClaimedDelegationReward && len(l.Topics) > 2 && len(l.Data) >= 32 {
				address := common.BytesToAddress(l.Topics[1][12:])
				stakerID := idx.StakerID(new(big.Int).SetBytes(l.Topics[2][:]).Uint64())
				reward := new(big.Int).SetBytes(l.Data[0:32])

				a.store.IncDelegatorClaimedRewards(address, reward)
				a.store.IncStakerDelegatorsClaimedRewards(stakerID, reward)
			}
		}
	}

	// Update EpochStats
	stats := a.Gossip.GetDirtyEpochStats()
	stats.TotalFee = new(big.Int).Add(stats.TotalFee, blockFee)
	if sealEpoch {
		// dirty EpochStats becomes active
		stats.End = block.Time
		a.Gossip.SetEpochStats(epoch, stats)

		// new dirty EpochStats
		a.Gossip.SetDirtyEpochStats(&sfctype.EpochStats{
			Start:    block.Time,
			TotalFee: new(big.Int),
		})
	} else {
		a.Gossip.SetDirtyEpochStats(stats)
	}

	// Write cheaters
	for _, stakerID := range cheaters {
		staker := a.store.GetSfcStaker(stakerID)
		if staker.HasFork() {
			continue
		}
		// write into DB
		staker.Status |= sfctype.ForkBit
		a.store.SetSfcStaker(stakerID, staker)
		// write into SFC contract
		position := sfcpos.Staker(stakerID)
		statedb.SetState(sfc.ContractAddress, position.Status(), utils.U64to256(staker.Status))
	}

	if sealEpoch {
		if a.store.HasSfcConstants(epoch) {
			a.store.SetSfcConstants(epoch+1, a.store.GetSfcConstants(epoch))
		}

		// Write offline validators
		for _, it := range a.store.GetSfcStakers() {
			if it.Staker.Offline() {
				continue
			}

			gotMissed := a.store.GetBlocksMissed(it.StakerID)
			badMissed := a.config.Economy.OfflinePenaltyThreshold
			if gotMissed.Num >= badMissed.BlocksNum && gotMissed.Period >= inter.Timestamp(badMissed.Period) {
				// write into DB
				it.Staker.Status |= sfctype.OfflineBit
				a.store.SetSfcStaker(it.StakerID, it.Staker)
				// write into SFC contract
				position := sfcpos.Staker(it.StakerID)
				statedb.SetState(sfc.ContractAddress, position.Status(), utils.U64to256(it.Staker.Status))
			}
		}

		// Write epoch snapshot (for reward)
		cheatersSet := cheaters.Set()
		epochPos := sfcpos.EpochSnapshot(epoch)
		epochValidators := a.store.GetEpochValidators(epoch)
		baseRewardWeights, txRewardWeights := a.calcRewardWeights(epochValidators, stats.Duration())

		totalBaseRewardWeight := new(big.Int)
		totalTxRewardWeight := new(big.Int)
		totalStake := new(big.Int)
		totalDelegated := new(big.Int)
		for i, it := range epochValidators {
			baseRewardWeight := baseRewardWeights[i]
			txRewardWeight := txRewardWeights[i]
			totalStake.Add(totalStake, it.Staker.StakeAmount)
			totalDelegated.Add(totalDelegated, it.Staker.DelegatedMe)

			if _, ok := cheatersSet[it.StakerID]; ok {
				continue // don't give reward to cheaters
			}
			if baseRewardWeight.Sign() == 0 && txRewardWeight.Sign() == 0 {
				continue // don't give reward to offline validators
			}

			meritPos := epochPos.ValidatorMerit(it.StakerID)

			statedb.SetState(sfc.ContractAddress, meritPos.StakeAmount(), utils.BigTo256(it.Staker.StakeAmount))
			statedb.SetState(sfc.ContractAddress, meritPos.DelegatedMe(), utils.BigTo256(it.Staker.DelegatedMe))
			statedb.SetState(sfc.ContractAddress, meritPos.BaseRewardWeight(), utils.BigTo256(baseRewardWeight))
			statedb.SetState(sfc.ContractAddress, meritPos.TxRewardWeight(), utils.BigTo256(txRewardWeight))

			totalBaseRewardWeight.Add(totalBaseRewardWeight, baseRewardWeight)
			totalTxRewardWeight.Add(totalTxRewardWeight, txRewardWeight)
		}
		baseRewardPerSec := a.getRewardPerSec(epoch)

		// set total supply
		baseRewards := new(big.Int).Mul(big.NewInt(stats.Duration().Unix()), baseRewardPerSec)
		rewards := new(big.Int).Add(baseRewards, stats.TotalFee)
		totalSupply := new(big.Int).Add(a.store.GetTotalSupply(), rewards)
		statedb.SetState(sfc.ContractAddress, sfcpos.CurrentSealedEpoch(), utils.U64to256(uint64(epoch)))
		a.store.SetTotalSupply(totalSupply)

		statedb.SetState(sfc.ContractAddress, epochPos.TotalBaseRewardWeight(), utils.BigTo256(totalBaseRewardWeight))
		statedb.SetState(sfc.ContractAddress, epochPos.TotalTxRewardWeight(), utils.BigTo256(totalTxRewardWeight))
		statedb.SetState(sfc.ContractAddress, epochPos.EpochFee(), utils.BigTo256(stats.TotalFee))
		statedb.SetState(sfc.ContractAddress, epochPos.EndTime(), utils.U64to256(uint64(stats.End.Unix())))
		statedb.SetState(sfc.ContractAddress, epochPos.Duration(), utils.U64to256(uint64(stats.Duration().Unix())))
		statedb.SetState(sfc.ContractAddress, epochPos.BaseRewardPerSecond(), utils.BigTo256(baseRewardPerSec))
		statedb.SetState(sfc.ContractAddress, epochPos.StakeTotalAmount(), utils.BigTo256(totalStake))
		statedb.SetState(sfc.ContractAddress, epochPos.DelegationsTotalAmount(), utils.BigTo256(totalDelegated))
		statedb.SetState(sfc.ContractAddress, epochPos.TotalSupply(), utils.BigTo256(totalSupply))
		statedb.SetState(sfc.ContractAddress, sfcpos.CurrentSealedEpoch(), utils.U64to256(uint64(epoch)))

		// Add balance for SFC to pay rewards
		statedb.AddBalance(sfc.ContractAddress, rewards)

		// Select new validators
		for _, it := range a.GetActiveSfcStakers() {
			// Note: cheaters are not active
			if _, ok := cheatersSet[it.StakerID]; ok {
				a.Log.Crit("Cheaters must be deactivated")
			}
			a.store.SetEpochValidator(epoch+1, it.StakerID, it.Staker)
		}
	}
}
