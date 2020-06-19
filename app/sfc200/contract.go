// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package sfc200

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ContractABI is the input ABI used to generate the binding from.
const ContractABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isDelegation\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BurntRewardStash\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fromEpoch\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"untilEpoch\",\"type\":\"uint256\"}],\"name\":\"ClaimedDelegationReward\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fromEpoch\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"untilEpoch\",\"type\":\"uint256\"}],\"name\":\"ClaimedValidatorReward\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"toStakerID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"CreatedDelegation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"dagSfcAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"CreatedStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"auth\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"wrID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"delegation\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"CreatedWithdrawRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"}],\"name\":\"DeactivatedDelegation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"}],\"name\":\"DeactivatedStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"diff\",\"type\":\"uint256\"}],\"name\":\"IncreasedDelegation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"diff\",\"type\":\"uint256\"}],\"name\":\"IncreasedStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fromEpoch\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"}],\"name\":\"LockingDelegation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fromEpoch\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"}],\"name\":\"LockingStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"auth\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"wrID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"delegation\",\"type\":\"bool\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"penalty\",\"type\":\"uint256\"}],\"name\":\"PartialWithdrawnByRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"}],\"name\":\"PreparedToWithdrawDelegation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"}],\"name\":\"PreparedToWithdrawStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"auth\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rewards\",\"type\":\"uint256\"}],\"name\":\"UnstashedRewards\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"UpdatedBaseRewardPerSec\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"oldStakerID\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"newStakerID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"UpdatedDelegation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"short\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"long\",\"type\":\"uint256\"}],\"name\":\"UpdatedGasPowerAllocationRate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"delegatedMe\",\"type\":\"uint256\"}],\"name\":\"UpdatedStake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"}],\"name\":\"UpdatedStakerMetadata\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldSfcAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newSfcAddress\",\"type\":\"address\"}],\"name\":\"UpdatedStakerSfcAddress\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"penalty\",\"type\":\"uint256\"}],\"name\":\"WithdrawnDelegation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"penalty\",\"type\":\"uint256\"}],\"name\":\"WithdrawnStake\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bool\",\"name\":\"isDelegation\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"}],\"name\":\"_rewardsBurnableOnDeactivation\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"sfcAddress\",\"type\":\"address\"}],\"name\":\"_sfcAddressToStakerID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"}],\"name\":\"_syncDelegator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"}],\"name\":\"_syncStaker\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"}],\"name\":\"_upgradeStakerStorage\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"bondedRatio\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"bondedTargetPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"bondedTargetRewardUnlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"bondedTargetStart\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fromEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxEpochs\",\"type\":\"uint256\"}],\"name\":\"calcDelegationRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_fromEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxEpochs\",\"type\":\"uint256\"}],\"name\":\"calcValidatorRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maxEpochs\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"}],\"name\":\"claimDelegationRewards\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"maxEpochs\",\"type\":\"uint256\"}],\"name\":\"claimValidatorRewards\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"contractCommission\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"to\",\"type\":\"uint256\"}],\"name\":\"createDelegation\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"name\":\"createStake\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"dagAdrress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sfcAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"name\":\"createStakeWithAddresses\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentEpoch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentSealedEpoch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"delegationLockPeriodEpochs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"delegationLockPeriodTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"delegations\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"createdEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deactivatedEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deactivatedTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"paidUntilEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"toStakerID\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"delegationsNum\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"delegationsTotalAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"delegations_v2\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"createdEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deactivatedEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deactivatedTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"paidUntilEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"toStakerID\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"epochSnapshots\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epochFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalBaseRewardWeight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalTxRewardWeight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseRewardPerSecond\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stakeTotalAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegationsTotalAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalLockedAmount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"e\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"v\",\"type\":\"uint256\"}],\"name\":\"epochValidator\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"stakeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegatedMe\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseRewardWeight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"txRewardWeight\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"firstLockedUpEpoch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getStakerID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"to\",\"type\":\"uint256\"}],\"name\":\"increaseDelegation\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"increaseStake\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"lockDuration\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"}],\"name\":\"lockUpDelegation\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"lockDuration\",\"type\":\"uint256\"}],\"name\":\"lockUpStake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"lockedDelegations\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"fromEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"lockedStakes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"fromEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxDelegatedRatio\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxStakerMetadataSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minDelegation\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minDelegationDecrease\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minDelegationIncrease\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minStakeDecrease\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minStakeIncrease\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"wrID\",\"type\":\"uint256\"}],\"name\":\"partialWithdrawByRequest\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"}],\"name\":\"prepareToWithdrawDelegation\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"wrID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"prepareToWithdrawDelegationPartial\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"prepareToWithdrawStake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"wrID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"prepareToWithdrawStakePartial\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"rewardsAllowed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"rewardsStash\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"slashedDelegationsTotalAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"slashedStakeTotalAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"stakeLockPeriodEpochs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"stakeLockPeriodTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"stakeTotalAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"stakerMetadata\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"stakers\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"status\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deactivatedEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deactivatedTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stakeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"paidUntilEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"delegatedMe\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"dagAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sfcAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"stakersLastID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"stakersNum\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epochNum\",\"type\":\"uint256\"}],\"name\":\"startLockedUp\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"unbondingStartDate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"unbondingUnlockPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"unlockedRatio\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unstashRewards\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"updateBaseRewardPerSec\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"short\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"long\",\"type\":\"uint256\"}],\"name\":\"updateGasPowerAllocationRate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"name\":\"updateStakerMetadata\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newSfcAddress\",\"type\":\"address\"}],\"name\":\"updateStakerSfcAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"validatorCommission\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"}],\"name\":\"withdrawDelegation\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"withdrawStake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"withdrawalRequests\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"stakerID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"delegation\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// ContractBin is the compiled bytecode used for deploying new contracts.
var ContractBin = "0x60806040819052600080546001600160a01b031916339081178255918291907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a35061556b806100536000396000f3fe6080604052600436106105075760003560e01c80638da5cb5b11610294578063c41b64051161015e578063dd099bb6116100d6578063f2fde38b1161008a578063f8b18d8a1161006f578063f8b18d8a14611149578063f99837e614611173578063fd5e6dd1146111a357610507565b8063f2fde38b146110ec578063f3ae5b1a1461111f57610507565b8063df4f49d4116100bb578063df4f49d414611098578063eac3baf2146110c2578063ec6a7f1c1461081257610507565b8063dd099bb61461101c578063df0e307a1461106e57610507565b8063cc8c21201161012d578063d845fc9011610112578063d845fc9014610fb2578063d9e257ef14610ff7578063dbd3aa8a14610fff57610507565b8063cc8c212014610ef7578063ce5aa00014610f9d57610507565b8063c41b640514610eb8578063c4b5dd7e1461050c578063c9400d4f14610ecd578063cb1c4e671461050c57610507565b8063ab2273c01161020c578063bb03a4bd116101c0578063bffe3486116101a5578063bffe348614610e35578063c312eb0714610e68578063c3d74f1a14610e8557610507565b8063bb03a4bd14610dea578063bed9d86114610e2057610507565b8063b1e64339116101f1578063b1e6433914610d37578063b42cb58d14610d61578063b9029d5014610d9457610507565b8063ab2273c014610cf0578063b1a3ebfa14610d0557610507565b806398ec2de511610263578063a70da4d211610248578063a70da4d214610c55578063a742b3d714610c6a578063a778651514610cdb57610507565b806398ec2de514610b86578063a4b89fab14610c2557610507565b80638da5cb5b14610a2d5780638f32d59b14610a5e57806390475ae414610a7357806396060e7114610b3257610507565b80633fee10a8116103d55780636a1cf4001161034d5780637b8c6b021161030157806381d9dc7a116102e657806381d9dc7a146109da5780638447c4df146109ef578063876f7e2a14610a1857610507565b80637b8c6b02146109b05780637cacb1d6146109c557610507565b80636f498663116103325780636f4986631461094d578063715018a614610986578063766718081461099b57610507565b80636a1cf400146109235780636e1a767a1461093857610507565b806354d77ed2116103a457806360c7e37f1161038957806360c7e37f1461050c57806363321e27146108f057806365cca35d146106a257610507565b806354d77ed2146105ce5780635dc03f1f146108b757610507565b80633fee10a8146108125780634bd202dc146108275780634e5a23281461083c57806353a72586146108a257610507565b80632265f2841161048357806330fa992911610437578063375b3c0a1161041c578063375b3c0a146107d35780633a0af4d4146107e85780633d0317fe146107fd57610507565b806330fa99291461070b57806333a149121461072057610507565b80632709275e116104685780632709275e146106a257806328dca8ff146106b7578063295cccba146106e157610507565b80632265f2841461065d57806326682c711461067257610507565b8063119e351a116104da5780631b593d8a116104bf5780631b593d8a146105a45780631d58179c146105ce5780631e8a6956146105e357610507565b8063119e351a1461057257806319ddb54f1461050c57610507565b8063029859921461050c578063041d97a31461053357806308728f6e146105485780630a29180c1461055d575b600080fd5b34801561051857600080fd5b50610521611229565b60408051918252519081900360200190f35b34801561053f57600080fd5b50610521611236565b34801561055457600080fd5b506105216112b1565b34801561056957600080fd5b506105216112b7565b34801561057e57600080fd5b506105a26004803603604081101561059557600080fd5b50803590602001356112bd565b005b3480156105b057600080fd5b506105a2600480360360208110156105c757600080fd5b5035611355565b3480156105da57600080fd5b506105216113e4565b3480156105ef57600080fd5b5061060d6004803603602081101561060657600080fd5b50356113e9565b604080519a8b5260208b0199909952898901979097526060890195909552608088019390935260a087019190915260c086015260e085015261010084015261012083015251908190036101400190f35b34801561066957600080fd5b5061052161143d565b34801561067e57600080fd5b506105a26004803603604081101561069557600080fd5b5080359060200135611444565b3480156106ae57600080fd5b506105216117cd565b3480156106c357600080fd5b506105a2600480360360208110156106da57600080fd5b50356117dd565b3480156106ed57600080fd5b506105a26004803603602081101561070457600080fd5b5035611895565b34801561071757600080fd5b50610521611949565b34801561072c57600080fd5b506105a26004803603602081101561074357600080fd5b81019060208101813564010000000081111561075e57600080fd5b82018360208201111561077057600080fd5b8035906020019184600183028401116401000000008311171561079257600080fd5b91908080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525092955061194f945050505050565b3480156107df57600080fd5b50610521611a11565b3480156107f457600080fd5b50610521611a20565b34801561080957600080fd5b50610521611a27565b34801561081e57600080fd5b50610521611a2d565b34801561083357600080fd5b50610521611a34565b34801561084857600080fd5b506108756004803603604081101561085f57600080fd5b506001600160a01b038135169060200135611a3a565b60408051958652602086019490945284840192909252606084015215156080830152519081900360a00190f35b3480156108ae57600080fd5b50610521611a77565b3480156108c357600080fd5b506105a2600480360360408110156108da57600080fd5b506001600160a01b038135169060200135611a7f565b3480156108fc57600080fd5b506105216004803603602081101561091357600080fd5b50356001600160a01b0316611af4565b34801561092f57600080fd5b50610521611b13565b34801561094457600080fd5b50610521611b7c565b34801561095957600080fd5b506105216004803603604081101561097057600080fd5b506001600160a01b038135169060200135611b82565b34801561099257600080fd5b506105a2611b9f565b3480156109a757600080fd5b50610521611c4f565b3480156109bc57600080fd5b50610521611c58565b3480156109d157600080fd5b50610521611c60565b3480156109e657600080fd5b50610521611c66565b3480156109fb57600080fd5b50610a04611c6c565b604080519115158252519081900360200190f35b348015610a2457600080fd5b506105a2611ca0565b348015610a3957600080fd5b50610a42611e0d565b604080516001600160a01b039092168252519081900360200190f35b348015610a6a57600080fd5b50610a04611e1c565b6105a260048036036060811015610a8957600080fd5b6001600160a01b038235811692602081013590911691810190606081016040820135640100000000811115610abd57600080fd5b820183602082011115610acf57600080fd5b80359060200191846001830284011164010000000083111715610af157600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550611e2d945050505050565b348015610b3e57600080fd5b50610b6860048036036060811015610b5557600080fd5b5080359060208101359060400135611eae565b60408051938452602084019290925282820152519081900360600190f35b348015610b9257600080fd5b50610bb060048036036020811015610ba957600080fd5b5035611f50565b6040805160208082528351818301528351919283929083019185019080838360005b83811015610bea578181015183820152602001610bd2565b50505050905090810190601f168015610c175780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b348015610c3157600080fd5b506105a260048036036040811015610c4857600080fd5b5080359060200135611feb565b348015610c6157600080fd5b506105216122b0565b348015610c7657600080fd5b50610ca360048036036040811015610c8d57600080fd5b506001600160a01b0381351690602001356122b6565b604080519788526020880196909652868601949094526060860192909252608085015260a084015260c0830152519081900360e00190f35b348015610ce757600080fd5b50610521612305565b348015610cfc57600080fd5b50610521612312565b348015610d1157600080fd5b50610a0460048036036040811015610d2857600080fd5b50803515159060200135612318565b348015610d4357600080fd5b506105a260048036036020811015610d5a57600080fd5b503561236f565b348015610d6d57600080fd5b5061052160048036036020811015610d8457600080fd5b50356001600160a01b03166124f3565b348015610da057600080fd5b50610dc460048036036040811015610db757600080fd5b5080359060200135612548565b604080519485526020850193909352838301919091526060830152519081900360800190f35b348015610df657600080fd5b506105a260048036036060811015610e0d57600080fd5b508035906020810135906040013561257b565b348015610e2c57600080fd5b506105a261288f565b348015610e4157600080fd5b50610ca360048036036020811015610e5857600080fd5b50356001600160a01b0316612b7a565b6105a260048036036020811015610e7e57600080fd5b5035612bb7565b348015610e9157600080fd5b506105a260048036036020811015610ea857600080fd5b50356001600160a01b0316612ea7565b348015610ec457600080fd5b506105a26130c5565b348015610ed957600080fd5b506105a260048036036020811015610ef057600080fd5b5035613231565b6105a260048036036020811015610f0d57600080fd5b810190602081018135640100000000811115610f2857600080fd5b820183602082011115610f3a57600080fd5b80359060200191846001830284011164010000000083111715610f5c57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550613348945050505050565b348015610fa957600080fd5b50610521613357565b348015610fbe57600080fd5b50610b6860048036036080811015610fd557600080fd5b506001600160a01b038135169060208101359060408101359060600135613365565b6105a26134df565b6105a26004803603602081101561101557600080fd5b50356135db565b34801561102857600080fd5b506110556004803603604081101561103f57600080fd5b506001600160a01b03813516906020013561387a565b6040805192835260208301919091528051918290030190f35b34801561107a57600080fd5b506105a26004803603602081101561109157600080fd5b503561389e565b3480156110a457600080fd5b50611055600480360360208110156110bb57600080fd5b5035613b63565b3480156110ce57600080fd5b506105a2600480360360208110156110e557600080fd5b5035613b7c565b3480156110f857600080fd5b506105a26004803603602081101561110f57600080fd5b50356001600160a01b0316613bdb565b34801561112b57600080fd5b506105a26004803603602081101561114257600080fd5b5035613c3d565b34801561115557600080fd5b506105a26004803603602081101561116c57600080fd5b5035613e42565b34801561117f57600080fd5b506105a26004803603604081101561119657600080fd5b5080359060200135614299565b3480156111af57600080fd5b506111cd600480360360208110156111c657600080fd5b5035614381565b604080519a8b5260208b0199909952898901979097526060890195909552608088019390935260a087019190915260c086015260e08501526001600160a01b039081166101008501521661012083015251908190036101400190f35b670de0b6b3a76400005b90565b601e546000908152601f60205260408120600901548061125a576000915050611233565b601e546000908152601f6020526040812060088101546007909101546112859163ffffffff6143de16565b90506112aa8261129e83620f424063ffffffff61443816565b9063ffffffff61449116565b9250505090565b60235481565b60285481565b6112c5611e1c565b611316576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b604080518381526020810183905281517f95ae5488127de4bc98492f4487556e7af9f37eb4b6d5e94f6d849e03ff76cc7c929181900390910190a15050565b61135d611e1c565b6113ae576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b6040805182815290517f8cd9dae1bbea2bc8a5e80ffce2c224727a25925130a03ae100619a8861ae23969181900360200190a150565b600390565b601f60205280600052604060002060009150905080600101549080600201549080600301549080600401549080600501549080600601549080600701549080600801549080600901549080600a015490508a565b62e4e1c090565b336000611450826124f3565b905061145b816144d3565b601e546000828152602080526040902060060154146114c1576040805162461bcd60e51b815260206004820152601760248201527f6e6f7420616c6c207265776172647320636c61696d6564000000000000000000604482015290519081900360640190fd5b60008181526030602052604090205415806114ec575060008181526030602052604090206001015442115b61153d576040805162461bcd60e51b815260206004820152600f60248201527f7374616b65206973206c6f636b65640000000000000000000000000000000000604482015290519081900360640190fd5b611545611229565b831015611599576040805162461bcd60e51b815260206004820152601060248201527f746f6f20736d616c6c20616d6f756e7400000000000000000000000000000000604482015290519081900360640190fd5b6000818152602080526040902060050154806115b3611a11565b85011115611608576040805162461bcd60e51b815260206004820152601c60248201527f6d757374206c65617665206174206c65617374206d696e5374616b6500000000604482015290519081900360640190fd5b6000828152602080526040902060070154848203906116268261453f565b1015611679576040805162461bcd60e51b815260206004820152601460248201527f746f6f206d7563682064656c65676174696f6e73000000000000000000000000604482015290519081900360640190fd5b6001600160a01b0384166000908152602d60209081526040808320898452909152902060030154156116f2576040805162461bcd60e51b815260206004820152601360248201527f7772494420616c72656164792065786973747300000000000000000000000000604482015290519081900360640190fd5b611700600084868886614560565b600083815260208080526040808320600501805489900390556001600160a01b0387168352602d82528083208984529091529020838155600301859055611745611c4f565b6001600160a01b0385166000818152602d602090815260408083208b8452825280832060018101959095554260029095019490945583518a8152908101919091528083018890529151859282917fde2d2a87af2fa2de55bde86f04143144eb632fa6be266dc224341a371fb8916d9181900360600190a46117c583613b7c565b505050505050565b600060646301c9c3805b04905090565b60008181526020805260409020600901546001600160a01b031615611849576040805162461bcd60e51b815260206004820152600b60248201527f6e6f742075706461746564000000000000000000000000000000000000000000604482015290519081900360640190fd5b61185281614670565b6000908152602080526040902060088101546009909101805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b03909216919091179055565b3360006118a1826124f3565b90506118ac81614670565b60008060006118bd84600088611eae565b600087815260208052604090206006015492955090935091506118e19083836146d2565b600084815260208052604090206006018190556118fe85846147d2565b6040805184815260208101849052808201839052905185917f2ea54c2b22a07549d19fb5eb8e4e48ebe1c653117215e94d5468c5612750d35c919081900360600190a2505050505050565b60265481565b600061195a336124f3565b905061196581614670565b61196d612312565b825111156119c2576040805162461bcd60e51b815260206004820152601060248201527f746f6f20626967206d6574616461746100000000000000000000000000000000604482015290519081900360640190fd5b6000818152602b6020908152604090912083516119e19285019061536e565b5060405181907fb7a99a0df6a9e15c2689e6a55811ef76cdb514c67d4a0e37fcb125ada0e3cd8390600090a25050565b6a02a055184a310c1260000090565b62ed4e0090565b60245481565b62093a8090565b60255481565b602d602090815260009283526040808420909152908252902080546001820154600283015460038401546004909401549293919290919060ff1685565b635e0580f890565b611a8882614882565b611a92828261494a565b6001600160a01b0382166000818152602e6020908152604080832085845282529182902060040154825190815291518493849390927f19b46b9014e4dc8ca74f505b8921797c6a8a489860217d15b3c7d741637dfcff92918290030190a45050565b6001600160a01b0381166000908152602160205260409020545b919050565b600080611b2e611b21611a77565b429063ffffffff6149c216565b90506000611b50611b3d611c58565b61129e620f42408563ffffffff61443816565b9050611b5a613357565b8110611b6b57600092505050611233565b80611b74613357565b039250505090565b602f5481565b602c60209081526000928352604080842090915290825290205481565b611ba7611e1c565b611bf8576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b600080546040516001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a36000805473ffffffffffffffffffffffffffffffffffffffff19169055565b601e5460010190565b63039ada0090565b601e5481565b60225481565b6000611c76611a20565b611c7e611a77565b0142101580611c9b5750611c90611b13565b611c98611236565b10155b905090565b336000818152602c60209081526040808320838052909152902054819080611d0f576040805162461bcd60e51b815260206004820152600a60248201527f6e6f207265776172647300000000000000000000000000000000000000000000604482015290519081900360640190fd5b611d17611c6c565b611d68576040805162461bcd60e51b815260206004820152601c60248201527f6265666f7265206d696e696d756d20756e6c6f636b20706572696f6400000000604482015290519081900360640190fd5b6001600160a01b038084166000908152602c60209081526040808320838052909152808220829055519184169183156108fc0291849190818181858888f19350505050158015611dbc573d6000803e3d6000fd5b50816001600160a01b0316836001600160a01b03167f80b36a0e929d7e7925087e54acfeecf4c6043e451b9d71ac5e908b66f9e5d126836040518082815260200191505060405180910390a3505050565b6000546001600160a01b031690565b6000546001600160a01b0316331490565b6001600160a01b038316151580611e4c57506001600160a01b03821615155b611e9d576040805162461bcd60e51b815260206004820152600f60248201527f696e76616c696420616464726573730000000000000000000000000000000000604482015290519081900360640190fd5b611ea983833484614a04565b505050565b600080600080611ed786602060008a815260200190815260200160002060060154600101614c6b565b60008881526020805260409020600601549091508111611eff57600093509150829050611f47565b600080825b601e548111158015611f17575087840181105b15611f3d57611f2e8a82611f29612305565b614c80565b90920191905060018101611f04565b5090945090925090505b93509350939050565b602b6020908152600091825260409182902080548351601f600260001961010060018616150201909316929092049182018490048402810184019094528084529091830182828015611fe35780601f10611fb857610100808354040283529160200191611fe3565b820191906000526020600020905b815481529060010190602001808311611fc657829003601f168201915b505050505081565b602f54158015906120115750601e5461200b90600163ffffffff6143de16565b602f5411155b612062576040805162461bcd60e51b815260206004820152601960248201527f6665617475726520776173206e6f742061637469766174656400000000000000604482015290519081900360640190fd5b3361206d818361494a565b6000828152602080526040902054156120cd576040805162461bcd60e51b815260206004820152601760248201527f7374616b65722073686f756c6420626520616374697665000000000000000000604482015290519081900360640190fd5b6212750083101580156120e457506301e133808311155b612135576040805162461bcd60e51b815260206004820152601260248201527f696e636f7272656374206475726174696f6e0000000000000000000000000000604482015290519081900360640190fd5b6000612147428563ffffffff6143de16565b60008481526030602052604090206001015490915081111561219a5760405162461bcd60e51b81526004018080602001828103825260228152602001806154ae6022913960400191505060405180910390fd5b6001600160a01b03821660009081526031602090815260408083208684529091529020600101548111612214576040805162461bcd60e51b815260206004820152601160248201527f616c7265616479206c6f636b6564207570000000000000000000000000000000604482015290519081900360640190fd5b601e5460009061222b90600163ffffffff6143de16565b60408051808201825282815260208082018681526001600160a01b0388166000818152603184528581208b8252845285902093518455905160019093019290925582518481529081018690528251939450879391927f823f252f996e1f519fd0215db7eb4d5a688d78587bf03bfb03d77bfca939806d92918290030190a35050505050565b60275481565b602e602052816000526040600020602052806000526040600020600091509150508060000154908060010154908060020154908060030154908060040154908060050154908060060154905087565b6000606462e4e1c06117d7565b61010090565b6000821580612366575060008281526020805260409020600501541580159061234c57506000828152602080526040902054155b801561236657506000828152602080526040902060040154155b90505b92915050565b3361237981614882565b6123838183614db3565b601e546001600160a01b0382166000908152602e60209081526040808320868452909152902060050154146123ff576040805162461bcd60e51b815260206004820152601760248201527f6e6f7420616c6c207265776172647320636c61696d6564000000000000000000604482015290519081900360640190fd5b6001600160a01b0381166000908152602e60209081526040808320858452909152902060040154612437906001908490849080614560565b61243f611c4f565b6001600160a01b0382166000908152602e602090815260408083208684528252808320600281019490945542600385015560049093015490805291902060050154156124b85760008381526020805260409020600701546124a6908263ffffffff6149c216565b60008481526020805260409020600701555b60405183906001600160a01b038416907f912c4125a208704a342cbdc4726795d26556b0170b7fc95bc706d5cb1f50646990600090a3505050565b6001600160a01b0381166000908152602160205260408120548061251b576000915050611b0e565b60008181526020805260409020600901546001600160a01b03848116911614612369576000915050611b0e565b6000918252601f602090815260408084209284529190529020805460018201546002830154600390930154919390929190565b3361258581614882565b61258f8184614db3565b601e546001600160a01b0382166000908152602e602090815260408083208784529091529020600501541461260b576040805162461bcd60e51b815260206004820152601760248201527f6e6f7420616c6c207265776172647320636c61696d6564000000000000000000604482015290519081900360640190fd5b612613611229565b821015612667576040805162461bcd60e51b815260206004820152601060248201527f746f6f20736d616c6c20616d6f756e7400000000000000000000000000000000604482015290519081900360640190fd5b6001600160a01b0381166000908152602e6020908152604080832086845290915290206004015480612697611229565b840111156126d65760405162461bcd60e51b81526004018080602001828103825260218152602001806154d06021913960400191505060405180910390fd5b6001600160a01b0382166000908152602d602090815260408083208884529091529020600301541561274f576040805162461bcd60e51b815260206004820152601360248201527f7772494420616c72656164792065786973747300000000000000000000000000604482015290519081900360640190fd5b61275d600185848685614560565b6000848152602080526040902060050154156127a6576000848152602080526040902060070154612794908463ffffffff6149c216565b60008581526020805260409020600701555b6001600160a01b0382166000818152602e60209081526040808320888452825280832060040180548890039055928252602d81528282208883529052208481556003018390556127f4611c4f565b6001600160a01b0383166000818152602d602090815260408083208a8452825291829020600180820195909555426002820155600401805460ff19168517905581518981529081019390935282810186905251869282917fde2d2a87af2fa2de55bde86f04143144eb632fa6be266dc224341a371fb8916d9181900360600190a461287f8285611a7f565b61288884613b7c565b5050505050565b33600061289b826124f3565b6000818152602080526040902060040154909150612900576040805162461bcd60e51b815260206004820152601960248201527f7374616b6572207761736e277420646561637469766174656400000000000000604482015290519081900360640190fd5b612908611a2d565b60008281526020805260409020600401540142101561296e576040805162461bcd60e51b815260206004820152601660248201527f6e6f7420656e6f7567682074696d652070617373656400000000000000000000604482015290519081900360640190fd5b6129766113e4565b600082815260208052604090206003015401612990611c4f565b10156129e3576040805162461bcd60e51b815260206004820152601860248201527f6e6f7420656e6f7567682065706f636873207061737365640000000000000000604482015290519081900360640190fd5b600081815260208080526040808320600881018054600583018054845488865560018087018a9055600287018a9055600387018a9055600487018a905592899055600686018990556007860189905573ffffffffffffffffffffffffffffffffffffffff1980851690955560099095018054909416909355602b9095529285206001600160a01b039093169490939092908216151590612a8390846153ec565b6001600160a01b038088166000908152602160205260408082208290559187168152908120558115612ac057600086815260208052604090208290555b60238054600019019055602454612add908563ffffffff6149c216565b60245580612b21576040516001600160a01b0388169085156108fc029086906000818181858888f19350505050158015612b1b573d6000803e3d6000fd5b50612b25565b8392505b602854612b38908463ffffffff6143de16565b60285560408051848152905187917f8c6548258f8f12a9d4b593fa89a223417ed901d4ee9712ba09beb4d56f5262b6919081900360200190a250505050505050565b6029602052600090815260409020805460018201546002830154600384015460048501546005860154600690960154949593949293919290919087565b33612bc182614e36565b612bc9611229565b341015612c1d576040805162461bcd60e51b815260206004820152601360248201527f696e73756666696369656e7420616d6f756e7400000000000000000000000000604482015290519081900360640190fd5b6001600160a01b0381166000908152602e6020908152604080832085845290915290206004015415612c96576040805162461bcd60e51b815260206004820152601960248201527f64656c65676174696f6e20616c72656164792065786973747300000000000000604482015290519081900360640190fd5b6001600160a01b03811660009081526021602052604090205415612d01576040805162461bcd60e51b815260206004820152600f60248201527f616c7265616479207374616b696e670000000000000000000000000000000000604482015290519081900360640190fd5b6000828152602080526040902060070154612d22903463ffffffff6143de16565b6000838152602080526040902060050154612d3c9061453f565b1015612d8f576040805162461bcd60e51b815260206004820152601a60248201527f7374616b65722773206c696d6974206973206578636565646564000000000000604482015290519081900360640190fd5b612d97615430565b612d9f611c4f565b8152426020808301918252346080840181815260c08501878152601e5460a087019081526001600160a01b0388166000908152602e865260408082208b8352875280822089518155975160018901558089015160028901556060890151600389015593516004880155905160058701559051600690950194909455918052912060070154612e329163ffffffff6143de16565b6000848152602080526040902060070155602580546001019055602654612e5f903463ffffffff6143de16565b60265560408051348152905184916001600160a01b038516917ffd8c857fb9acd6f4ad59b8621a2a77825168b7b4b76de9586d08e00d4ed462be9181900360200190a3505050565b33612eb182614882565b612eba81614882565b816001600160a01b0316816001600160a01b03161415612f21576040805162461bcd60e51b815260206004820152601060248201527f7468652073616d65206164647265737300000000000000000000000000000000604482015290519081900360640190fd5b6000612f2c826124f3565b9050612f3781614670565b6001600160a01b0383166000908152602160205260409020541580612f7357506001600160a01b03831660009081526021602052604090205481145b612fc4576040805162461bcd60e51b815260206004820152601460248201527f6164647265737320616c72656164792075736564000000000000000000000000604482015290519081900360640190fd5b60008181526020808052604080832060098101805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b03898116918217909255878216808752602186528487208790559086528386208790556008909201541684528184208590558352602c82528083208380529091529020541561307f576001600160a01b038281166000908152602c60208181526040808420848052808352818520958916855292825280842084805282528320845490555290555b826001600160a01b0316826001600160a01b0316827f7cc102ee500cbca85691c9642080562e8f012b04d27f5b7f389453672b20694660405160405180910390a4505050565b3360006130d1826124f3565b90506130dc816144d3565b601e54600082815260208052604090206006015414613142576040805162461bcd60e51b815260206004820152601760248201527f6e6f7420616c6c207265776172647320636c61696d6564000000000000000000604482015290519081900360640190fd5b600081815260306020526040902054158061316d575060008181526030602052604090206001015442115b6131be576040805162461bcd60e51b815260206004820152600f60248201527f7374616b65206973206c6f636b65640000000000000000000000000000000000604482015290519081900360640190fd5b60008181526020805260408120600501546131de91908390859080614560565b6131e6611c4f565b6000828152602080526040808220600381019390935542600490930192909255905182917ff7c308d0d978cce3aec157d1b34e355db4636b4e71ce91b4f5ec9e7a4f5cdc6091a25050565b613239611e1c565b61328a576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b601e5481116132e0576040805162461bcd60e51b815260206004820152601760248201527f63616e277420737461727420696e207468652070617374000000000000000000604482015290519081900360640190fd5b602f5415806132f25750601e54602f54115b613343576040805162461bcd60e51b815260206004820152601360248201527f6665617475726520776173207374617274656400000000000000000000000000604482015290519081900360640190fd5b602f55565b61335433333484614a04565b50565b600060646304c4b4006117d7565b6001600160a01b03841660009081526029602052604081206004015481908190156133c15760405162461bcd60e51b81526004018080602001828103825260258152602001806154f16025913960400191505060405180910390fd5b6001600160a01b0387166000908152602e602090815260408083208984529091528120600501546133f6908790600101614c6b565b6001600160a01b0389166000908152602e602090815260408083208b84529091529020600301549091501561342757fe5b6001600160a01b0388166000908152602e602090815260408083208a84529091529020600501548111613462576000935091508290506134d5565b600080825b601e54811115801561347a575087840181105b156134cb576001600160a01b038b166000908152602e602090815260408083208d84529091529020600401546134bc908b9083906134b6612305565b8f614e9f565b90920191905060018101613467565b5090945090925090505b9450945094915050565b60006134ea336124f3565b90506134f4611229565b341015613548576040805162461bcd60e51b815260206004820152601360248201527f696e73756666696369656e7420616d6f756e7400000000000000000000000000604482015290519081900360640190fd5b61355181614e36565b6000818152602080526040812060050154613572903463ffffffff6143de16565b6000838152602080526040902060050181905560245490915061359b903463ffffffff6143de16565b60245560408051828152346020820152815184927fa1d93e9a2a16bf4c2d0cdc6f47fe0fa054c741c96b3dac1297c79eaca31714e9928290030190a25050565b336135e581614882565b6135ef8183614db3565b601e546001600160a01b0382166000908152602e602090815260408083208684529091529020600501541461366b576040805162461bcd60e51b815260206004820152601760248201527f6e6f7420616c6c207265776172647320636c61696d6564000000000000000000604482015290519081900360640190fd5b6001600160a01b0381166000908152602e6020908152604082209184905252613692611229565b3410156136e6576040805162461bcd60e51b815260206004820152601360248201527f696e73756666696369656e7420616d6f756e7400000000000000000000000000604482015290519081900360640190fd5b6000828152602080526040902060070154613707903463ffffffff6143de16565b60008381526020805260409020600501546137219061453f565b1015613774576040805162461bcd60e51b815260206004820152601a60248201527f7374616b65722773206c696d6974206973206578636565646564000000000000604482015290519081900360640190fd5b61377d82614e36565b6001600160a01b0381166000908152602e602090815260408083208584529091528120600401546137b4903463ffffffff6143de16565b6001600160a01b0383166000908152602e6020908152604080832087845282528083206004018490559080529020600701549091506137f9903463ffffffff6143de16565b600084815260208052604090206007015560265461381d903463ffffffff6143de16565b60265560408051828152346020820152815185926001600160a01b038616927f4ca781bfe171e588a2661d5a7f2f5f59df879c53489063552fbad2145b707fc1929081900390910190a36138718284611a7f565b611ea983613b7c565b60316020908152600092835260408084209091529082529020805460019091015482565b336138a881614882565b6001600160a01b0381166000908152602e60209081526040808320858452909152902060030154613920576040805162461bcd60e51b815260206004820152601d60248201527f64656c65676174696f6e207761736e2774206465616374697661746564000000604482015290519081900360640190fd5b600082815260208052604090206005015415613a455761393e611a2d565b6001600160a01b0382166000908152602e60209081526040808320868452909152902060030154014210156139ba576040805162461bcd60e51b815260206004820152601660248201527f6e6f7420656e6f7567682074696d652070617373656400000000000000000000604482015290519081900360640190fd5b6139c26113e4565b6001600160a01b0382166000908152602e60209081526040808320868452909152902060020154016139f2611c4f565b1015613a45576040805162461bcd60e51b815260206004820152601860248201527f6e6f7420656e6f7567682065706f636873207061737365640000000000000000604482015290519081900360640190fd5b600082815260208080526040808320546001600160a01b0385168452602e835281842086855290925282206004810180548483556001808401869055600284018690556003840186905591859055600583018590556006909201849055602580546000190190556026549216151591613abe90826149c2565b60265581613b02576040516001600160a01b0385169082156108fc029083906000818181858888f19350505050158015613afc573d6000803e3d6000fd5b50613b06565b8092505b602754613b19908463ffffffff6143de16565b60275560408051848152905186916001600160a01b038716917f87e86b3710b72c10173ca52c6a9f9cf2df27e77ed177741a8b4feb12bb7a606f9181900360200190a35050505050565b6030602052600090815260409020805460019091015482565b613b8581614670565b600081815260208080526040918290206005810154600790910154835191825291810191909152815183927f509404fa75ce234a1273cf9f7918bcf54e0ef19f2772e4f71b6526606a723b7c928290030190a250565b613be3611e1c565b613c34576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b61335481614fff565b602f5415801590613c635750601e54613c5d90600163ffffffff6143de16565b602f5411155b613cb4576040805162461bcd60e51b815260206004820152601960248201527f6665617475726520776173206e6f742061637469766174656400000000000000604482015290519081900360640190fd5b6000613cbf336124f3565b9050613cca816144d3565b621275008210158015613ce157506301e133808211155b613d32576040805162461bcd60e51b815260206004820152601260248201527f696e636f7272656374206475726174696f6e0000000000000000000000000000604482015290519081900360640190fd5b613d42428363ffffffff6143de16565b60008281526030602052604090206001015410613da6576040805162461bcd60e51b815260206004820152601160248201527f616c7265616479206c6f636b6564207570000000000000000000000000000000604482015290519081900360640190fd5b601e54600090613dbd90600163ffffffff6143de16565b90506000613dd1428563ffffffff6143de16565b604080518082018252848152602080820184815260008881526030835284902092518355516001909201919091558151858152908101839052815192935085927f71f8e76b11dde805fa567e857c4beba340500f4ca9da003520a25014f542162b929181900390910190a250505050565b336000818152602d602090815260408083208584529091529020600201548190613eb3576040805162461bcd60e51b815260206004820152601560248201527f7265717565737420646f65736e27742065786973740000000000000000000000604482015290519081900360640190fd5b6001600160a01b0382166000908152602d6020908152604080832086845290915290206004810154905460ff90911690818015613eff5750600081815260208052604090206005015415155b1561401857613f0c611a2d565b6001600160a01b0385166000908152602d6020908152604080832089845290915290206002015401421015613f88576040805162461bcd60e51b815260206004820152601660248201527f6e6f7420656e6f7567682074696d652070617373656400000000000000000000604482015290519081900360640190fd5b613f906113e4565b6001600160a01b0385166000908152602d6020908152604080832089845290915290206001015401613fc0611c4f565b1015614013576040805162461bcd60e51b815260206004820152601860248201527f6e6f7420656e6f7567682065706f636873207061737365640000000000000000604482015290519081900360640190fd5b61412c565b8161412c57614025611a2d565b6001600160a01b0385166000908152602d60209081526040808320898452909152902060020154014210156140a1576040805162461bcd60e51b815260206004820152601660248201527f6e6f7420656e6f7567682074696d652070617373656400000000000000000000604482015290519081900360640190fd5b6140a96113e4565b6001600160a01b0385166000908152602d60209081526040808320898452909152902060010154016140d9611c4f565b101561412c576040805162461bcd60e51b815260206004820152601860248201527f6e6f7420656e6f7567682065706f636873207061737365640000000000000000604482015290519081900360640190fd5b600081815260208080526040808320546001600160a01b0388168452602d83528184208985529092528220600381018054848355600180840186905560028401869055918590556004909201805460ff1916905590911615159084156141a75760265461419f908263ffffffff6149c216565b6026556141be565b6024546141ba908263ffffffff6149c216565b6024555b816141ff576040516001600160a01b0387169082156108fc029083906000818181858888f193505050501580156141f9573d6000803e3d6000fd5b50614203565b8092505b84156142245760275461421c908463ffffffff6143de16565b60275561423b565b602854614237908463ffffffff6143de16565b6028555b604080518981528615156020820152808201859052905185916001600160a01b03808a1692908b16917fd5304dabc5bd47105b6921889d1b528c4b2223250248a916afd129b1c0512ddd919081900360600190a45050505050505050565b336142a381614882565b6142ad8183614db3565b60008060006142bf8486600089613365565b6001600160a01b0387166000908152602e602090815260408083208b845290915290206005015492955090935091506142f99083836146d2565b6001600160a01b0384166000908152602e60209081526040808320888452909152902060050181905561432c84846147d2565b6040805184815260208101849052808201839052905186916001600160a01b038716917f2676e1697cf4731b93ddb4ef54e0e5a98c06cccbbbb2202848a3c6286595e6ce9181900360600190a3505050505050565b60208052600090815260409020805460018201546002830154600384015460048501546005860154600687015460078801546008890154600990990154979896979596949593949293919290916001600160a01b0391821691168a565b600082820183811015612366576040805162461bcd60e51b815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b60008261444757506000612369565b8282028284828161445457fe5b04146123665760405162461bcd60e51b81526004018080602001828103825260218152602001806155166021913960400191505060405180910390fd5b600061236683836040518060400160405280601a81526020017f536166654d6174683a206469766973696f6e206279207a65726f0000000000008152506150ac565b6144dc81614670565b600081815260208052604090206004015415613354576040805162461bcd60e51b815260206004820152601560248201527f7374616b65722069732064656163746976617465640000000000000000000000604482015290519081900360640190fd5b6000612369620f424061129e61455361143d565b859063ffffffff61443816565b61456a8585612318565b15612888576000614581828463ffffffff6149c216565b6001600160a01b0385166000908152602c602090815260408083208380529091528120549192506145bc8461129e848663ffffffff61443816565b9050806145ec576001600160a01b0386166000908152602c60209081526040808320838052909152812055614613565b6001600160a01b0386166000908152602c6020908152604080832083805290915290208190555b818114614666576040805189151581528284036020820152815189926001600160a01b038a16927f0ea92567e76d40ddc52d2c1d74a521a59329a38b50411451de6ad2e565466d0f929081900390910190a35b5050505050505050565b6000818152602080526040902060050154613354576040805162461bcd60e51b815260206004820152601460248201527f7374616b657220646f65736e2774206578697374000000000000000000000000604482015290519081900360640190fd5b818310614726576040805162461bcd60e51b815260206004820152601560248201527f65706f636820697320616c726561647920706169640000000000000000000000604482015290519081900360640190fd5b601e5482111561477d576040805162461bcd60e51b815260206004820152600c60248201527f6675747572652065706f63680000000000000000000000000000000000000000604482015290519081900360640190fd5b81811015611ea9576040805162461bcd60e51b815260206004820152601160248201527f6e6f2065706f63687320636c61696d6564000000000000000000000000000000604482015290519081900360640190fd5b806147dc5761487e565b6147e4611c6c565b15614825576040516001600160a01b0383169082156108fc029083906000818181858888f1935050505015801561481f573d6000803e3d6000fd5b5061487e565b6001600160a01b0382166000908152602c60209081526040808320838052909152902054614859908263ffffffff6143de16565b6001600160a01b0383166000908152602c602090815260408083208380529091529020555b5050565b6001600160a01b03811660009081526029602052604090206004015415613354576001600160a01b03166000818152602960208181526040808420602e8352818520600680830180548852918552928620825481556001808401805491830191909155600280850180549184019190915560038086018054918501919091556004808701805491860191909155600580880180549187019190915586549590980194909455998952969095529186905592859055928490559383905590829055918190559055565b6001600160a01b0382166000908152602e6020908152604080832084845290915290206004015461487e576040805162461bcd60e51b815260206004820152601860248201527f64656c65676174696f6e20646f65736e27742065786973740000000000000000604482015290519081900360640190fd5b600061236683836040518060400160405280601e81526020017f536166654d6174683a207375627472616374696f6e206f766572666c6f77000081525061514e565b6001600160a01b0384166000908152602160205260409020541580614a3f57506001600160a01b038316600090815260216020526040902054155b614a90576040805162461bcd60e51b815260206004820152601560248201527f7374616b657220616c7265616479206578697374730000000000000000000000604482015290519081900360640190fd5b614a98611a11565b821015614aec576040805162461bcd60e51b815260206004820152601360248201527f696e73756666696369656e7420616d6f756e7400000000000000000000000000604482015290519081900360640190fd5b60228054600101908190556001600160a01b0380861660009081526021602090815260408083208590559287168252828220849055838252805220600501839055614b35611c4f565b600082815260208052604090206001808201929092554260028201556008810180546001600160a01b03808a1673ffffffffffffffffffffffffffffffffffffffff199283161790925560098301805492891692909116919091179055601e54600690910155602380549091019055602454614bb7908463ffffffff6143de16565b6024556040805184815290516001600160a01b0387169183917f0697dfe5062b9db8108e4b31254f47a912ae6bbb78837667b2e923a6f5160d399181900360200190a3815115614c0a57614c0a8261194f565b836001600160a01b0316856001600160a01b03161461288857836001600160a01b0316856001600160a01b0316827f7cc102ee500cbca85691c9642080562e8f012b04d27f5b7f389453672b20694660405160405180910390a45050505050565b600082614c79575080612369565b5090919050565b600080614c8d85856151a8565b6000858152601f60209081526040808320898452909152812080546001909101549293509190614cbd83836143de565b905080614cd1576000945050505050614dac565b6000614cfa614ced620f424061129e868b63ffffffff61443816565b859063ffffffff6143de16565b9050614d108261129e878463ffffffff61443816565b95506000602f54118015614d265750602f548810155b15614da6576000898152603060205260409020548810801590614d845750601f6000614d598a600163ffffffff6149c216565b815260200190815260200160002060010154603060008b815260200190815260200160002060010154115b15614da657614da3614d96858a6152d3565b879063ffffffff6143de16565b95505b50505050505b9392505050565b614dbd828261494a565b6001600160a01b0382166000908152602e602090815260408083208484529091529020600301541561487e576040805162461bcd60e51b815260206004820152601960248201527f64656c65676174696f6e20697320646561637469766174656400000000000000604482015290519081900360640190fd5b614e3f816144d3565b600081815260208052604090205415613354576040805162461bcd60e51b815260206004820152601760248201527f7374616b65722073686f756c6420626520616374697665000000000000000000604482015290519081900360640190fd5b600080614eac87876151a8565b6000878152601f602090815260408083208b8452909152812080546001909101549293509190614edc83836143de565b905080614ef0576000945050505050614ff6565b6000614f19620f424061129e614f0c828c63ffffffff6149c216565b8c9063ffffffff61443816565b9050614f2f8261129e878463ffffffff61443816565b95506000602f54118015614f455750602f548a10155b15614ff0576001600160a01b03871660009081526031602090815260408083208e84529091529020548a10801590614fdb5750601f6000614f8d8c600163ffffffff6149c216565b81526020019081526020016000206001015460316000896001600160a01b03166001600160a01b0316815260200190815260200160002060008d815260200190815260200160002060010154115b15614ff057614fed614d968a8c6152d3565b95505b50505050505b95945050505050565b6001600160a01b0381166150445760405162461bcd60e51b81526004018080602001828103825260268152602001806154886026913960400191505060405180910390fd5b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a36000805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0392909216919091179055565b600081836151385760405162461bcd60e51b81526004018080602001828103825283818151815260200191508051906020019080838360005b838110156150fd5781810151838201526020016150e5565b50505050905090810190601f16801561512a5780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b50600083858161514457fe5b0495945050505050565b600081848411156151a05760405162461bcd60e51b81526020600482018181528351602484015283519092839260449091019190850190808383600083156150fd5781810151838201526020016150e5565b505050900390565b6000818152601f60209081526040808320600481015486855292819052908320600281015460059092015460039091015484831561525d576000878152601f6020526040812060068101546002909101546152089163ffffffff61443816565b90506000602f5411801561521e5750602f548810155b1561524557615242620f424061129e6152356117cd565b849063ffffffff61443816565b90505b6152598661129e838863ffffffff61443816565b9150505b600082156152b6576000888152601f602052604090206003015461528d90859061129e908663ffffffff61443816565b90506152b3620f424061129e6152a16117cd565b8490620f42400363ffffffff61443816565b90505b6152c6828263ffffffff6143de16565b9998505050505050505050565b6000818152601f60205260408120600a015415615365576000828152601f6020526040812060068101546002909101546153129163ffffffff61443816565b6000848152601f60205260409020600a015490915061535d9061129e86615351620f42408361533f6117cd565b8890620f42400363ffffffff61443816565b9063ffffffff61443816565b915050612369565b50600092915050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106153af57805160ff19168380011785556153dc565b828001600101855582156153dc579182015b828111156153dc5782518255916020019190600101906153c1565b506153e892915061546d565b5090565b50805460018160011615610100020316600290046000825580601f106154125750613354565b601f016020900490600052602060002090810190613354919061546d565b6040518060e00160405280600081526020016000815260200160008152602001600081526020016000815260200160008152602001600081525090565b61123391905b808211156153e8576000815560010161547356fe4f776e61626c653a206e6577206f776e657220697320746865207a65726f20616464726573737374616b65722773206c6f636b696e672077696c6c2066696e6973682066697273746d757374206c65617665206174206c65617374206d696e44656c65676174696f6e6f6c642076657273696f6e2064656c65676174696f6e2c20706c6561736520757064617465536166654d6174683a206d756c7469706c69636174696f6e206f766572666c6f77a265627a7a72315820f922f230fa8274fcf3870de13a369115bfc06eec092a903702b2fe204eaca22064736f6c634300050c0032"

// DeployContract deploys a new Ethereum contract, binding an instance of Contract to it.
func DeployContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Contract, error) {
	parsed, err := abi.JSON(strings.NewReader(ContractABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// Contract is an auto generated Go binding around an Ethereum contract.
type Contract struct {
	ContractCaller     // Read-only binding to the contract
	ContractTransactor // Write-only binding to the contract
	ContractFilterer   // Log filterer for contract events
}

// ContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractSession struct {
	Contract     *Contract         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractCallerSession struct {
	Contract *ContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractTransactorSession struct {
	Contract     *ContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractRaw struct {
	Contract *Contract // Generic contract binding to access the raw methods on
}

// ContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractCallerRaw struct {
	Contract *ContractCaller // Generic read-only contract binding to access the raw methods on
}

// ContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractTransactorRaw struct {
	Contract *ContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContract creates a new instance of Contract, bound to a specific deployed contract.
func NewContract(address common.Address, backend bind.ContractBackend) (*Contract, error) {
	contract, err := bindContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// NewContractCaller creates a new read-only instance of Contract, bound to a specific deployed contract.
func NewContractCaller(address common.Address, caller bind.ContractCaller) (*ContractCaller, error) {
	contract, err := bindContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractCaller{contract: contract}, nil
}

// NewContractTransactor creates a new write-only instance of Contract, bound to a specific deployed contract.
func NewContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractTransactor, error) {
	contract, err := bindContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractTransactor{contract: contract}, nil
}

// NewContractFilterer creates a new log filterer instance of Contract, bound to a specific deployed contract.
func NewContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractFilterer, error) {
	contract, err := bindContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractFilterer{contract: contract}, nil
}

// bindContract binds a generic wrapper to an already deployed contract.
func bindContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.ContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transact(opts, method, params...)
}

// RewardsBurnableOnDeactivation is a free data retrieval call binding the contract method 0xb1a3ebfa.
//
// Solidity: function _rewardsBurnableOnDeactivation(bool isDelegation, uint256 stakerID) constant returns(bool)
func (_Contract *ContractCaller) RewardsBurnableOnDeactivation(opts *bind.CallOpts, isDelegation bool, stakerID *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "_rewardsBurnableOnDeactivation", isDelegation, stakerID)
	return *ret0, err
}

// RewardsBurnableOnDeactivation is a free data retrieval call binding the contract method 0xb1a3ebfa.
//
// Solidity: function _rewardsBurnableOnDeactivation(bool isDelegation, uint256 stakerID) constant returns(bool)
func (_Contract *ContractSession) RewardsBurnableOnDeactivation(isDelegation bool, stakerID *big.Int) (bool, error) {
	return _Contract.Contract.RewardsBurnableOnDeactivation(&_Contract.CallOpts, isDelegation, stakerID)
}

// RewardsBurnableOnDeactivation is a free data retrieval call binding the contract method 0xb1a3ebfa.
//
// Solidity: function _rewardsBurnableOnDeactivation(bool isDelegation, uint256 stakerID) constant returns(bool)
func (_Contract *ContractCallerSession) RewardsBurnableOnDeactivation(isDelegation bool, stakerID *big.Int) (bool, error) {
	return _Contract.Contract.RewardsBurnableOnDeactivation(&_Contract.CallOpts, isDelegation, stakerID)
}

// SfcAddressToStakerID is a free data retrieval call binding the contract method 0xb42cb58d.
//
// Solidity: function _sfcAddressToStakerID(address sfcAddress) constant returns(uint256)
func (_Contract *ContractCaller) SfcAddressToStakerID(opts *bind.CallOpts, sfcAddress common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "_sfcAddressToStakerID", sfcAddress)
	return *ret0, err
}

// SfcAddressToStakerID is a free data retrieval call binding the contract method 0xb42cb58d.
//
// Solidity: function _sfcAddressToStakerID(address sfcAddress) constant returns(uint256)
func (_Contract *ContractSession) SfcAddressToStakerID(sfcAddress common.Address) (*big.Int, error) {
	return _Contract.Contract.SfcAddressToStakerID(&_Contract.CallOpts, sfcAddress)
}

// SfcAddressToStakerID is a free data retrieval call binding the contract method 0xb42cb58d.
//
// Solidity: function _sfcAddressToStakerID(address sfcAddress) constant returns(uint256)
func (_Contract *ContractCallerSession) SfcAddressToStakerID(sfcAddress common.Address) (*big.Int, error) {
	return _Contract.Contract.SfcAddressToStakerID(&_Contract.CallOpts, sfcAddress)
}

// BondedRatio is a free data retrieval call binding the contract method 0x041d97a3.
//
// Solidity: function bondedRatio() constant returns(uint256)
func (_Contract *ContractCaller) BondedRatio(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "bondedRatio")
	return *ret0, err
}

// BondedRatio is a free data retrieval call binding the contract method 0x041d97a3.
//
// Solidity: function bondedRatio() constant returns(uint256)
func (_Contract *ContractSession) BondedRatio() (*big.Int, error) {
	return _Contract.Contract.BondedRatio(&_Contract.CallOpts)
}

// BondedRatio is a free data retrieval call binding the contract method 0x041d97a3.
//
// Solidity: function bondedRatio() constant returns(uint256)
func (_Contract *ContractCallerSession) BondedRatio() (*big.Int, error) {
	return _Contract.Contract.BondedRatio(&_Contract.CallOpts)
}

// BondedTargetPeriod is a free data retrieval call binding the contract method 0x7b8c6b02.
//
// Solidity: function bondedTargetPeriod() constant returns(uint256)
func (_Contract *ContractCaller) BondedTargetPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "bondedTargetPeriod")
	return *ret0, err
}

// BondedTargetPeriod is a free data retrieval call binding the contract method 0x7b8c6b02.
//
// Solidity: function bondedTargetPeriod() constant returns(uint256)
func (_Contract *ContractSession) BondedTargetPeriod() (*big.Int, error) {
	return _Contract.Contract.BondedTargetPeriod(&_Contract.CallOpts)
}

// BondedTargetPeriod is a free data retrieval call binding the contract method 0x7b8c6b02.
//
// Solidity: function bondedTargetPeriod() constant returns(uint256)
func (_Contract *ContractCallerSession) BondedTargetPeriod() (*big.Int, error) {
	return _Contract.Contract.BondedTargetPeriod(&_Contract.CallOpts)
}

// BondedTargetRewardUnlock is a free data retrieval call binding the contract method 0x6a1cf400.
//
// Solidity: function bondedTargetRewardUnlock() constant returns(uint256)
func (_Contract *ContractCaller) BondedTargetRewardUnlock(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "bondedTargetRewardUnlock")
	return *ret0, err
}

// BondedTargetRewardUnlock is a free data retrieval call binding the contract method 0x6a1cf400.
//
// Solidity: function bondedTargetRewardUnlock() constant returns(uint256)
func (_Contract *ContractSession) BondedTargetRewardUnlock() (*big.Int, error) {
	return _Contract.Contract.BondedTargetRewardUnlock(&_Contract.CallOpts)
}

// BondedTargetRewardUnlock is a free data retrieval call binding the contract method 0x6a1cf400.
//
// Solidity: function bondedTargetRewardUnlock() constant returns(uint256)
func (_Contract *ContractCallerSession) BondedTargetRewardUnlock() (*big.Int, error) {
	return _Contract.Contract.BondedTargetRewardUnlock(&_Contract.CallOpts)
}

// BondedTargetStart is a free data retrieval call binding the contract method 0xce5aa000.
//
// Solidity: function bondedTargetStart() constant returns(uint256)
func (_Contract *ContractCaller) BondedTargetStart(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "bondedTargetStart")
	return *ret0, err
}

// BondedTargetStart is a free data retrieval call binding the contract method 0xce5aa000.
//
// Solidity: function bondedTargetStart() constant returns(uint256)
func (_Contract *ContractSession) BondedTargetStart() (*big.Int, error) {
	return _Contract.Contract.BondedTargetStart(&_Contract.CallOpts)
}

// BondedTargetStart is a free data retrieval call binding the contract method 0xce5aa000.
//
// Solidity: function bondedTargetStart() constant returns(uint256)
func (_Contract *ContractCallerSession) BondedTargetStart() (*big.Int, error) {
	return _Contract.Contract.BondedTargetStart(&_Contract.CallOpts)
}

// CalcDelegationRewards is a free data retrieval call binding the contract method 0xd845fc90.
//
// Solidity: function calcDelegationRewards(address delegator, uint256 stakerID, uint256 _fromEpoch, uint256 maxEpochs) constant returns(uint256, uint256, uint256)
func (_Contract *ContractCaller) CalcDelegationRewards(opts *bind.CallOpts, delegator common.Address, stakerID *big.Int, _fromEpoch *big.Int, maxEpochs *big.Int) (*big.Int, *big.Int, *big.Int, error) {
	var (
		ret0 = new(*big.Int)
		ret1 = new(*big.Int)
		ret2 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}
	err := _Contract.contract.Call(opts, out, "calcDelegationRewards", delegator, stakerID, _fromEpoch, maxEpochs)
	return *ret0, *ret1, *ret2, err
}

// CalcDelegationRewards is a free data retrieval call binding the contract method 0xd845fc90.
//
// Solidity: function calcDelegationRewards(address delegator, uint256 stakerID, uint256 _fromEpoch, uint256 maxEpochs) constant returns(uint256, uint256, uint256)
func (_Contract *ContractSession) CalcDelegationRewards(delegator common.Address, stakerID *big.Int, _fromEpoch *big.Int, maxEpochs *big.Int) (*big.Int, *big.Int, *big.Int, error) {
	return _Contract.Contract.CalcDelegationRewards(&_Contract.CallOpts, delegator, stakerID, _fromEpoch, maxEpochs)
}

// CalcDelegationRewards is a free data retrieval call binding the contract method 0xd845fc90.
//
// Solidity: function calcDelegationRewards(address delegator, uint256 stakerID, uint256 _fromEpoch, uint256 maxEpochs) constant returns(uint256, uint256, uint256)
func (_Contract *ContractCallerSession) CalcDelegationRewards(delegator common.Address, stakerID *big.Int, _fromEpoch *big.Int, maxEpochs *big.Int) (*big.Int, *big.Int, *big.Int, error) {
	return _Contract.Contract.CalcDelegationRewards(&_Contract.CallOpts, delegator, stakerID, _fromEpoch, maxEpochs)
}

// CalcValidatorRewards is a free data retrieval call binding the contract method 0x96060e71.
//
// Solidity: function calcValidatorRewards(uint256 stakerID, uint256 _fromEpoch, uint256 maxEpochs) constant returns(uint256, uint256, uint256)
func (_Contract *ContractCaller) CalcValidatorRewards(opts *bind.CallOpts, stakerID *big.Int, _fromEpoch *big.Int, maxEpochs *big.Int) (*big.Int, *big.Int, *big.Int, error) {
	var (
		ret0 = new(*big.Int)
		ret1 = new(*big.Int)
		ret2 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}
	err := _Contract.contract.Call(opts, out, "calcValidatorRewards", stakerID, _fromEpoch, maxEpochs)
	return *ret0, *ret1, *ret2, err
}

// CalcValidatorRewards is a free data retrieval call binding the contract method 0x96060e71.
//
// Solidity: function calcValidatorRewards(uint256 stakerID, uint256 _fromEpoch, uint256 maxEpochs) constant returns(uint256, uint256, uint256)
func (_Contract *ContractSession) CalcValidatorRewards(stakerID *big.Int, _fromEpoch *big.Int, maxEpochs *big.Int) (*big.Int, *big.Int, *big.Int, error) {
	return _Contract.Contract.CalcValidatorRewards(&_Contract.CallOpts, stakerID, _fromEpoch, maxEpochs)
}

// CalcValidatorRewards is a free data retrieval call binding the contract method 0x96060e71.
//
// Solidity: function calcValidatorRewards(uint256 stakerID, uint256 _fromEpoch, uint256 maxEpochs) constant returns(uint256, uint256, uint256)
func (_Contract *ContractCallerSession) CalcValidatorRewards(stakerID *big.Int, _fromEpoch *big.Int, maxEpochs *big.Int) (*big.Int, *big.Int, *big.Int, error) {
	return _Contract.Contract.CalcValidatorRewards(&_Contract.CallOpts, stakerID, _fromEpoch, maxEpochs)
}

// ContractCommission is a free data retrieval call binding the contract method 0x2709275e.
//
// Solidity: function contractCommission() constant returns(uint256)
func (_Contract *ContractCaller) ContractCommission(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "contractCommission")
	return *ret0, err
}

// ContractCommission is a free data retrieval call binding the contract method 0x2709275e.
//
// Solidity: function contractCommission() constant returns(uint256)
func (_Contract *ContractSession) ContractCommission() (*big.Int, error) {
	return _Contract.Contract.ContractCommission(&_Contract.CallOpts)
}

// ContractCommission is a free data retrieval call binding the contract method 0x2709275e.
//
// Solidity: function contractCommission() constant returns(uint256)
func (_Contract *ContractCallerSession) ContractCommission() (*big.Int, error) {
	return _Contract.Contract.ContractCommission(&_Contract.CallOpts)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() constant returns(uint256)
func (_Contract *ContractCaller) CurrentEpoch(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "currentEpoch")
	return *ret0, err
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() constant returns(uint256)
func (_Contract *ContractSession) CurrentEpoch() (*big.Int, error) {
	return _Contract.Contract.CurrentEpoch(&_Contract.CallOpts)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() constant returns(uint256)
func (_Contract *ContractCallerSession) CurrentEpoch() (*big.Int, error) {
	return _Contract.Contract.CurrentEpoch(&_Contract.CallOpts)
}

// CurrentSealedEpoch is a free data retrieval call binding the contract method 0x7cacb1d6.
//
// Solidity: function currentSealedEpoch() constant returns(uint256)
func (_Contract *ContractCaller) CurrentSealedEpoch(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "currentSealedEpoch")
	return *ret0, err
}

// CurrentSealedEpoch is a free data retrieval call binding the contract method 0x7cacb1d6.
//
// Solidity: function currentSealedEpoch() constant returns(uint256)
func (_Contract *ContractSession) CurrentSealedEpoch() (*big.Int, error) {
	return _Contract.Contract.CurrentSealedEpoch(&_Contract.CallOpts)
}

// CurrentSealedEpoch is a free data retrieval call binding the contract method 0x7cacb1d6.
//
// Solidity: function currentSealedEpoch() constant returns(uint256)
func (_Contract *ContractCallerSession) CurrentSealedEpoch() (*big.Int, error) {
	return _Contract.Contract.CurrentSealedEpoch(&_Contract.CallOpts)
}

// DelegationLockPeriodEpochs is a free data retrieval call binding the contract method 0x1d58179c.
//
// Solidity: function delegationLockPeriodEpochs() constant returns(uint256)
func (_Contract *ContractCaller) DelegationLockPeriodEpochs(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "delegationLockPeriodEpochs")
	return *ret0, err
}

// DelegationLockPeriodEpochs is a free data retrieval call binding the contract method 0x1d58179c.
//
// Solidity: function delegationLockPeriodEpochs() constant returns(uint256)
func (_Contract *ContractSession) DelegationLockPeriodEpochs() (*big.Int, error) {
	return _Contract.Contract.DelegationLockPeriodEpochs(&_Contract.CallOpts)
}

// DelegationLockPeriodEpochs is a free data retrieval call binding the contract method 0x1d58179c.
//
// Solidity: function delegationLockPeriodEpochs() constant returns(uint256)
func (_Contract *ContractCallerSession) DelegationLockPeriodEpochs() (*big.Int, error) {
	return _Contract.Contract.DelegationLockPeriodEpochs(&_Contract.CallOpts)
}

// DelegationLockPeriodTime is a free data retrieval call binding the contract method 0xec6a7f1c.
//
// Solidity: function delegationLockPeriodTime() constant returns(uint256)
func (_Contract *ContractCaller) DelegationLockPeriodTime(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "delegationLockPeriodTime")
	return *ret0, err
}

// DelegationLockPeriodTime is a free data retrieval call binding the contract method 0xec6a7f1c.
//
// Solidity: function delegationLockPeriodTime() constant returns(uint256)
func (_Contract *ContractSession) DelegationLockPeriodTime() (*big.Int, error) {
	return _Contract.Contract.DelegationLockPeriodTime(&_Contract.CallOpts)
}

// DelegationLockPeriodTime is a free data retrieval call binding the contract method 0xec6a7f1c.
//
// Solidity: function delegationLockPeriodTime() constant returns(uint256)
func (_Contract *ContractCallerSession) DelegationLockPeriodTime() (*big.Int, error) {
	return _Contract.Contract.DelegationLockPeriodTime(&_Contract.CallOpts)
}

// Delegations is a free data retrieval call binding the contract method 0xbffe3486.
//
// Solidity: function delegations(address ) constant returns(uint256 createdEpoch, uint256 createdTime, uint256 deactivatedEpoch, uint256 deactivatedTime, uint256 amount, uint256 paidUntilEpoch, uint256 toStakerID)
func (_Contract *ContractCaller) Delegations(opts *bind.CallOpts, arg0 common.Address) (struct {
	CreatedEpoch     *big.Int
	CreatedTime      *big.Int
	DeactivatedEpoch *big.Int
	DeactivatedTime  *big.Int
	Amount           *big.Int
	PaidUntilEpoch   *big.Int
	ToStakerID       *big.Int
}, error) {
	ret := new(struct {
		CreatedEpoch     *big.Int
		CreatedTime      *big.Int
		DeactivatedEpoch *big.Int
		DeactivatedTime  *big.Int
		Amount           *big.Int
		PaidUntilEpoch   *big.Int
		ToStakerID       *big.Int
	})
	out := ret
	err := _Contract.contract.Call(opts, out, "delegations", arg0)
	return *ret, err
}

// Delegations is a free data retrieval call binding the contract method 0xbffe3486.
//
// Solidity: function delegations(address ) constant returns(uint256 createdEpoch, uint256 createdTime, uint256 deactivatedEpoch, uint256 deactivatedTime, uint256 amount, uint256 paidUntilEpoch, uint256 toStakerID)
func (_Contract *ContractSession) Delegations(arg0 common.Address) (struct {
	CreatedEpoch     *big.Int
	CreatedTime      *big.Int
	DeactivatedEpoch *big.Int
	DeactivatedTime  *big.Int
	Amount           *big.Int
	PaidUntilEpoch   *big.Int
	ToStakerID       *big.Int
}, error) {
	return _Contract.Contract.Delegations(&_Contract.CallOpts, arg0)
}

// Delegations is a free data retrieval call binding the contract method 0xbffe3486.
//
// Solidity: function delegations(address ) constant returns(uint256 createdEpoch, uint256 createdTime, uint256 deactivatedEpoch, uint256 deactivatedTime, uint256 amount, uint256 paidUntilEpoch, uint256 toStakerID)
func (_Contract *ContractCallerSession) Delegations(arg0 common.Address) (struct {
	CreatedEpoch     *big.Int
	CreatedTime      *big.Int
	DeactivatedEpoch *big.Int
	DeactivatedTime  *big.Int
	Amount           *big.Int
	PaidUntilEpoch   *big.Int
	ToStakerID       *big.Int
}, error) {
	return _Contract.Contract.Delegations(&_Contract.CallOpts, arg0)
}

// DelegationsNum is a free data retrieval call binding the contract method 0x4bd202dc.
//
// Solidity: function delegationsNum() constant returns(uint256)
func (_Contract *ContractCaller) DelegationsNum(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "delegationsNum")
	return *ret0, err
}

// DelegationsNum is a free data retrieval call binding the contract method 0x4bd202dc.
//
// Solidity: function delegationsNum() constant returns(uint256)
func (_Contract *ContractSession) DelegationsNum() (*big.Int, error) {
	return _Contract.Contract.DelegationsNum(&_Contract.CallOpts)
}

// DelegationsNum is a free data retrieval call binding the contract method 0x4bd202dc.
//
// Solidity: function delegationsNum() constant returns(uint256)
func (_Contract *ContractCallerSession) DelegationsNum() (*big.Int, error) {
	return _Contract.Contract.DelegationsNum(&_Contract.CallOpts)
}

// DelegationsTotalAmount is a free data retrieval call binding the contract method 0x30fa9929.
//
// Solidity: function delegationsTotalAmount() constant returns(uint256)
func (_Contract *ContractCaller) DelegationsTotalAmount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "delegationsTotalAmount")
	return *ret0, err
}

// DelegationsTotalAmount is a free data retrieval call binding the contract method 0x30fa9929.
//
// Solidity: function delegationsTotalAmount() constant returns(uint256)
func (_Contract *ContractSession) DelegationsTotalAmount() (*big.Int, error) {
	return _Contract.Contract.DelegationsTotalAmount(&_Contract.CallOpts)
}

// DelegationsTotalAmount is a free data retrieval call binding the contract method 0x30fa9929.
//
// Solidity: function delegationsTotalAmount() constant returns(uint256)
func (_Contract *ContractCallerSession) DelegationsTotalAmount() (*big.Int, error) {
	return _Contract.Contract.DelegationsTotalAmount(&_Contract.CallOpts)
}

// DelegationsV2 is a free data retrieval call binding the contract method 0xa742b3d7.
//
// Solidity: function delegations_v2(address , uint256 ) constant returns(uint256 createdEpoch, uint256 createdTime, uint256 deactivatedEpoch, uint256 deactivatedTime, uint256 amount, uint256 paidUntilEpoch, uint256 toStakerID)
func (_Contract *ContractCaller) DelegationsV2(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (struct {
	CreatedEpoch     *big.Int
	CreatedTime      *big.Int
	DeactivatedEpoch *big.Int
	DeactivatedTime  *big.Int
	Amount           *big.Int
	PaidUntilEpoch   *big.Int
	ToStakerID       *big.Int
}, error) {
	ret := new(struct {
		CreatedEpoch     *big.Int
		CreatedTime      *big.Int
		DeactivatedEpoch *big.Int
		DeactivatedTime  *big.Int
		Amount           *big.Int
		PaidUntilEpoch   *big.Int
		ToStakerID       *big.Int
	})
	out := ret
	err := _Contract.contract.Call(opts, out, "delegations_v2", arg0, arg1)
	return *ret, err
}

// DelegationsV2 is a free data retrieval call binding the contract method 0xa742b3d7.
//
// Solidity: function delegations_v2(address , uint256 ) constant returns(uint256 createdEpoch, uint256 createdTime, uint256 deactivatedEpoch, uint256 deactivatedTime, uint256 amount, uint256 paidUntilEpoch, uint256 toStakerID)
func (_Contract *ContractSession) DelegationsV2(arg0 common.Address, arg1 *big.Int) (struct {
	CreatedEpoch     *big.Int
	CreatedTime      *big.Int
	DeactivatedEpoch *big.Int
	DeactivatedTime  *big.Int
	Amount           *big.Int
	PaidUntilEpoch   *big.Int
	ToStakerID       *big.Int
}, error) {
	return _Contract.Contract.DelegationsV2(&_Contract.CallOpts, arg0, arg1)
}

// DelegationsV2 is a free data retrieval call binding the contract method 0xa742b3d7.
//
// Solidity: function delegations_v2(address , uint256 ) constant returns(uint256 createdEpoch, uint256 createdTime, uint256 deactivatedEpoch, uint256 deactivatedTime, uint256 amount, uint256 paidUntilEpoch, uint256 toStakerID)
func (_Contract *ContractCallerSession) DelegationsV2(arg0 common.Address, arg1 *big.Int) (struct {
	CreatedEpoch     *big.Int
	CreatedTime      *big.Int
	DeactivatedEpoch *big.Int
	DeactivatedTime  *big.Int
	Amount           *big.Int
	PaidUntilEpoch   *big.Int
	ToStakerID       *big.Int
}, error) {
	return _Contract.Contract.DelegationsV2(&_Contract.CallOpts, arg0, arg1)
}

// EpochSnapshots is a free data retrieval call binding the contract method 0x1e8a6956.
//
// Solidity: function epochSnapshots(uint256 ) constant returns(uint256 endTime, uint256 duration, uint256 epochFee, uint256 totalBaseRewardWeight, uint256 totalTxRewardWeight, uint256 baseRewardPerSecond, uint256 stakeTotalAmount, uint256 delegationsTotalAmount, uint256 totalSupply, uint256 totalLockedAmount)
func (_Contract *ContractCaller) EpochSnapshots(opts *bind.CallOpts, arg0 *big.Int) (struct {
	EndTime                *big.Int
	Duration               *big.Int
	EpochFee               *big.Int
	TotalBaseRewardWeight  *big.Int
	TotalTxRewardWeight    *big.Int
	BaseRewardPerSecond    *big.Int
	StakeTotalAmount       *big.Int
	DelegationsTotalAmount *big.Int
	TotalSupply            *big.Int
	TotalLockedAmount      *big.Int
}, error) {
	ret := new(struct {
		EndTime                *big.Int
		Duration               *big.Int
		EpochFee               *big.Int
		TotalBaseRewardWeight  *big.Int
		TotalTxRewardWeight    *big.Int
		BaseRewardPerSecond    *big.Int
		StakeTotalAmount       *big.Int
		DelegationsTotalAmount *big.Int
		TotalSupply            *big.Int
		TotalLockedAmount      *big.Int
	})
	out := ret
	err := _Contract.contract.Call(opts, out, "epochSnapshots", arg0)
	return *ret, err
}

// EpochSnapshots is a free data retrieval call binding the contract method 0x1e8a6956.
//
// Solidity: function epochSnapshots(uint256 ) constant returns(uint256 endTime, uint256 duration, uint256 epochFee, uint256 totalBaseRewardWeight, uint256 totalTxRewardWeight, uint256 baseRewardPerSecond, uint256 stakeTotalAmount, uint256 delegationsTotalAmount, uint256 totalSupply, uint256 totalLockedAmount)
func (_Contract *ContractSession) EpochSnapshots(arg0 *big.Int) (struct {
	EndTime                *big.Int
	Duration               *big.Int
	EpochFee               *big.Int
	TotalBaseRewardWeight  *big.Int
	TotalTxRewardWeight    *big.Int
	BaseRewardPerSecond    *big.Int
	StakeTotalAmount       *big.Int
	DelegationsTotalAmount *big.Int
	TotalSupply            *big.Int
	TotalLockedAmount      *big.Int
}, error) {
	return _Contract.Contract.EpochSnapshots(&_Contract.CallOpts, arg0)
}

// EpochSnapshots is a free data retrieval call binding the contract method 0x1e8a6956.
//
// Solidity: function epochSnapshots(uint256 ) constant returns(uint256 endTime, uint256 duration, uint256 epochFee, uint256 totalBaseRewardWeight, uint256 totalTxRewardWeight, uint256 baseRewardPerSecond, uint256 stakeTotalAmount, uint256 delegationsTotalAmount, uint256 totalSupply, uint256 totalLockedAmount)
func (_Contract *ContractCallerSession) EpochSnapshots(arg0 *big.Int) (struct {
	EndTime                *big.Int
	Duration               *big.Int
	EpochFee               *big.Int
	TotalBaseRewardWeight  *big.Int
	TotalTxRewardWeight    *big.Int
	BaseRewardPerSecond    *big.Int
	StakeTotalAmount       *big.Int
	DelegationsTotalAmount *big.Int
	TotalSupply            *big.Int
	TotalLockedAmount      *big.Int
}, error) {
	return _Contract.Contract.EpochSnapshots(&_Contract.CallOpts, arg0)
}

// EpochValidator is a free data retrieval call binding the contract method 0xb9029d50.
//
// Solidity: function epochValidator(uint256 e, uint256 v) constant returns(uint256 stakeAmount, uint256 delegatedMe, uint256 baseRewardWeight, uint256 txRewardWeight)
func (_Contract *ContractCaller) EpochValidator(opts *bind.CallOpts, e *big.Int, v *big.Int) (struct {
	StakeAmount      *big.Int
	DelegatedMe      *big.Int
	BaseRewardWeight *big.Int
	TxRewardWeight   *big.Int
}, error) {
	ret := new(struct {
		StakeAmount      *big.Int
		DelegatedMe      *big.Int
		BaseRewardWeight *big.Int
		TxRewardWeight   *big.Int
	})
	out := ret
	err := _Contract.contract.Call(opts, out, "epochValidator", e, v)
	return *ret, err
}

// EpochValidator is a free data retrieval call binding the contract method 0xb9029d50.
//
// Solidity: function epochValidator(uint256 e, uint256 v) constant returns(uint256 stakeAmount, uint256 delegatedMe, uint256 baseRewardWeight, uint256 txRewardWeight)
func (_Contract *ContractSession) EpochValidator(e *big.Int, v *big.Int) (struct {
	StakeAmount      *big.Int
	DelegatedMe      *big.Int
	BaseRewardWeight *big.Int
	TxRewardWeight   *big.Int
}, error) {
	return _Contract.Contract.EpochValidator(&_Contract.CallOpts, e, v)
}

// EpochValidator is a free data retrieval call binding the contract method 0xb9029d50.
//
// Solidity: function epochValidator(uint256 e, uint256 v) constant returns(uint256 stakeAmount, uint256 delegatedMe, uint256 baseRewardWeight, uint256 txRewardWeight)
func (_Contract *ContractCallerSession) EpochValidator(e *big.Int, v *big.Int) (struct {
	StakeAmount      *big.Int
	DelegatedMe      *big.Int
	BaseRewardWeight *big.Int
	TxRewardWeight   *big.Int
}, error) {
	return _Contract.Contract.EpochValidator(&_Contract.CallOpts, e, v)
}

// FirstLockedUpEpoch is a free data retrieval call binding the contract method 0x6e1a767a.
//
// Solidity: function firstLockedUpEpoch() constant returns(uint256)
func (_Contract *ContractCaller) FirstLockedUpEpoch(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "firstLockedUpEpoch")
	return *ret0, err
}

// FirstLockedUpEpoch is a free data retrieval call binding the contract method 0x6e1a767a.
//
// Solidity: function firstLockedUpEpoch() constant returns(uint256)
func (_Contract *ContractSession) FirstLockedUpEpoch() (*big.Int, error) {
	return _Contract.Contract.FirstLockedUpEpoch(&_Contract.CallOpts)
}

// FirstLockedUpEpoch is a free data retrieval call binding the contract method 0x6e1a767a.
//
// Solidity: function firstLockedUpEpoch() constant returns(uint256)
func (_Contract *ContractCallerSession) FirstLockedUpEpoch() (*big.Int, error) {
	return _Contract.Contract.FirstLockedUpEpoch(&_Contract.CallOpts)
}

// GetStakerID is a free data retrieval call binding the contract method 0x63321e27.
//
// Solidity: function getStakerID(address addr) constant returns(uint256)
func (_Contract *ContractCaller) GetStakerID(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "getStakerID", addr)
	return *ret0, err
}

// GetStakerID is a free data retrieval call binding the contract method 0x63321e27.
//
// Solidity: function getStakerID(address addr) constant returns(uint256)
func (_Contract *ContractSession) GetStakerID(addr common.Address) (*big.Int, error) {
	return _Contract.Contract.GetStakerID(&_Contract.CallOpts, addr)
}

// GetStakerID is a free data retrieval call binding the contract method 0x63321e27.
//
// Solidity: function getStakerID(address addr) constant returns(uint256)
func (_Contract *ContractCallerSession) GetStakerID(addr common.Address) (*big.Int, error) {
	return _Contract.Contract.GetStakerID(&_Contract.CallOpts, addr)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Contract *ContractCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Contract *ContractSession) IsOwner() (bool, error) {
	return _Contract.Contract.IsOwner(&_Contract.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Contract *ContractCallerSession) IsOwner() (bool, error) {
	return _Contract.Contract.IsOwner(&_Contract.CallOpts)
}

// LockedDelegations is a free data retrieval call binding the contract method 0xdd099bb6.
//
// Solidity: function lockedDelegations(address , uint256 ) constant returns(uint256 fromEpoch, uint256 endTime)
func (_Contract *ContractCaller) LockedDelegations(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (struct {
	FromEpoch *big.Int
	EndTime   *big.Int
}, error) {
	ret := new(struct {
		FromEpoch *big.Int
		EndTime   *big.Int
	})
	out := ret
	err := _Contract.contract.Call(opts, out, "lockedDelegations", arg0, arg1)
	return *ret, err
}

// LockedDelegations is a free data retrieval call binding the contract method 0xdd099bb6.
//
// Solidity: function lockedDelegations(address , uint256 ) constant returns(uint256 fromEpoch, uint256 endTime)
func (_Contract *ContractSession) LockedDelegations(arg0 common.Address, arg1 *big.Int) (struct {
	FromEpoch *big.Int
	EndTime   *big.Int
}, error) {
	return _Contract.Contract.LockedDelegations(&_Contract.CallOpts, arg0, arg1)
}

// LockedDelegations is a free data retrieval call binding the contract method 0xdd099bb6.
//
// Solidity: function lockedDelegations(address , uint256 ) constant returns(uint256 fromEpoch, uint256 endTime)
func (_Contract *ContractCallerSession) LockedDelegations(arg0 common.Address, arg1 *big.Int) (struct {
	FromEpoch *big.Int
	EndTime   *big.Int
}, error) {
	return _Contract.Contract.LockedDelegations(&_Contract.CallOpts, arg0, arg1)
}

// LockedStakes is a free data retrieval call binding the contract method 0xdf4f49d4.
//
// Solidity: function lockedStakes(uint256 ) constant returns(uint256 fromEpoch, uint256 endTime)
func (_Contract *ContractCaller) LockedStakes(opts *bind.CallOpts, arg0 *big.Int) (struct {
	FromEpoch *big.Int
	EndTime   *big.Int
}, error) {
	ret := new(struct {
		FromEpoch *big.Int
		EndTime   *big.Int
	})
	out := ret
	err := _Contract.contract.Call(opts, out, "lockedStakes", arg0)
	return *ret, err
}

// LockedStakes is a free data retrieval call binding the contract method 0xdf4f49d4.
//
// Solidity: function lockedStakes(uint256 ) constant returns(uint256 fromEpoch, uint256 endTime)
func (_Contract *ContractSession) LockedStakes(arg0 *big.Int) (struct {
	FromEpoch *big.Int
	EndTime   *big.Int
}, error) {
	return _Contract.Contract.LockedStakes(&_Contract.CallOpts, arg0)
}

// LockedStakes is a free data retrieval call binding the contract method 0xdf4f49d4.
//
// Solidity: function lockedStakes(uint256 ) constant returns(uint256 fromEpoch, uint256 endTime)
func (_Contract *ContractCallerSession) LockedStakes(arg0 *big.Int) (struct {
	FromEpoch *big.Int
	EndTime   *big.Int
}, error) {
	return _Contract.Contract.LockedStakes(&_Contract.CallOpts, arg0)
}

// MaxDelegatedRatio is a free data retrieval call binding the contract method 0x2265f284.
//
// Solidity: function maxDelegatedRatio() constant returns(uint256)
func (_Contract *ContractCaller) MaxDelegatedRatio(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "maxDelegatedRatio")
	return *ret0, err
}

// MaxDelegatedRatio is a free data retrieval call binding the contract method 0x2265f284.
//
// Solidity: function maxDelegatedRatio() constant returns(uint256)
func (_Contract *ContractSession) MaxDelegatedRatio() (*big.Int, error) {
	return _Contract.Contract.MaxDelegatedRatio(&_Contract.CallOpts)
}

// MaxDelegatedRatio is a free data retrieval call binding the contract method 0x2265f284.
//
// Solidity: function maxDelegatedRatio() constant returns(uint256)
func (_Contract *ContractCallerSession) MaxDelegatedRatio() (*big.Int, error) {
	return _Contract.Contract.MaxDelegatedRatio(&_Contract.CallOpts)
}

// MaxStakerMetadataSize is a free data retrieval call binding the contract method 0xab2273c0.
//
// Solidity: function maxStakerMetadataSize() constant returns(uint256)
func (_Contract *ContractCaller) MaxStakerMetadataSize(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "maxStakerMetadataSize")
	return *ret0, err
}

// MaxStakerMetadataSize is a free data retrieval call binding the contract method 0xab2273c0.
//
// Solidity: function maxStakerMetadataSize() constant returns(uint256)
func (_Contract *ContractSession) MaxStakerMetadataSize() (*big.Int, error) {
	return _Contract.Contract.MaxStakerMetadataSize(&_Contract.CallOpts)
}

// MaxStakerMetadataSize is a free data retrieval call binding the contract method 0xab2273c0.
//
// Solidity: function maxStakerMetadataSize() constant returns(uint256)
func (_Contract *ContractCallerSession) MaxStakerMetadataSize() (*big.Int, error) {
	return _Contract.Contract.MaxStakerMetadataSize(&_Contract.CallOpts)
}

// MinDelegation is a free data retrieval call binding the contract method 0x02985992.
//
// Solidity: function minDelegation() constant returns(uint256)
func (_Contract *ContractCaller) MinDelegation(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "minDelegation")
	return *ret0, err
}

// MinDelegation is a free data retrieval call binding the contract method 0x02985992.
//
// Solidity: function minDelegation() constant returns(uint256)
func (_Contract *ContractSession) MinDelegation() (*big.Int, error) {
	return _Contract.Contract.MinDelegation(&_Contract.CallOpts)
}

// MinDelegation is a free data retrieval call binding the contract method 0x02985992.
//
// Solidity: function minDelegation() constant returns(uint256)
func (_Contract *ContractCallerSession) MinDelegation() (*big.Int, error) {
	return _Contract.Contract.MinDelegation(&_Contract.CallOpts)
}

// MinDelegationDecrease is a free data retrieval call binding the contract method 0xcb1c4e67.
//
// Solidity: function minDelegationDecrease() constant returns(uint256)
func (_Contract *ContractCaller) MinDelegationDecrease(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "minDelegationDecrease")
	return *ret0, err
}

// MinDelegationDecrease is a free data retrieval call binding the contract method 0xcb1c4e67.
//
// Solidity: function minDelegationDecrease() constant returns(uint256)
func (_Contract *ContractSession) MinDelegationDecrease() (*big.Int, error) {
	return _Contract.Contract.MinDelegationDecrease(&_Contract.CallOpts)
}

// MinDelegationDecrease is a free data retrieval call binding the contract method 0xcb1c4e67.
//
// Solidity: function minDelegationDecrease() constant returns(uint256)
func (_Contract *ContractCallerSession) MinDelegationDecrease() (*big.Int, error) {
	return _Contract.Contract.MinDelegationDecrease(&_Contract.CallOpts)
}

// MinDelegationIncrease is a free data retrieval call binding the contract method 0x60c7e37f.
//
// Solidity: function minDelegationIncrease() constant returns(uint256)
func (_Contract *ContractCaller) MinDelegationIncrease(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "minDelegationIncrease")
	return *ret0, err
}

// MinDelegationIncrease is a free data retrieval call binding the contract method 0x60c7e37f.
//
// Solidity: function minDelegationIncrease() constant returns(uint256)
func (_Contract *ContractSession) MinDelegationIncrease() (*big.Int, error) {
	return _Contract.Contract.MinDelegationIncrease(&_Contract.CallOpts)
}

// MinDelegationIncrease is a free data retrieval call binding the contract method 0x60c7e37f.
//
// Solidity: function minDelegationIncrease() constant returns(uint256)
func (_Contract *ContractCallerSession) MinDelegationIncrease() (*big.Int, error) {
	return _Contract.Contract.MinDelegationIncrease(&_Contract.CallOpts)
}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() constant returns(uint256)
func (_Contract *ContractCaller) MinStake(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "minStake")
	return *ret0, err
}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() constant returns(uint256)
func (_Contract *ContractSession) MinStake() (*big.Int, error) {
	return _Contract.Contract.MinStake(&_Contract.CallOpts)
}

// MinStake is a free data retrieval call binding the contract method 0x375b3c0a.
//
// Solidity: function minStake() constant returns(uint256)
func (_Contract *ContractCallerSession) MinStake() (*big.Int, error) {
	return _Contract.Contract.MinStake(&_Contract.CallOpts)
}

// MinStakeDecrease is a free data retrieval call binding the contract method 0x19ddb54f.
//
// Solidity: function minStakeDecrease() constant returns(uint256)
func (_Contract *ContractCaller) MinStakeDecrease(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "minStakeDecrease")
	return *ret0, err
}

// MinStakeDecrease is a free data retrieval call binding the contract method 0x19ddb54f.
//
// Solidity: function minStakeDecrease() constant returns(uint256)
func (_Contract *ContractSession) MinStakeDecrease() (*big.Int, error) {
	return _Contract.Contract.MinStakeDecrease(&_Contract.CallOpts)
}

// MinStakeDecrease is a free data retrieval call binding the contract method 0x19ddb54f.
//
// Solidity: function minStakeDecrease() constant returns(uint256)
func (_Contract *ContractCallerSession) MinStakeDecrease() (*big.Int, error) {
	return _Contract.Contract.MinStakeDecrease(&_Contract.CallOpts)
}

// MinStakeIncrease is a free data retrieval call binding the contract method 0xc4b5dd7e.
//
// Solidity: function minStakeIncrease() constant returns(uint256)
func (_Contract *ContractCaller) MinStakeIncrease(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "minStakeIncrease")
	return *ret0, err
}

// MinStakeIncrease is a free data retrieval call binding the contract method 0xc4b5dd7e.
//
// Solidity: function minStakeIncrease() constant returns(uint256)
func (_Contract *ContractSession) MinStakeIncrease() (*big.Int, error) {
	return _Contract.Contract.MinStakeIncrease(&_Contract.CallOpts)
}

// MinStakeIncrease is a free data retrieval call binding the contract method 0xc4b5dd7e.
//
// Solidity: function minStakeIncrease() constant returns(uint256)
func (_Contract *ContractCallerSession) MinStakeIncrease() (*big.Int, error) {
	return _Contract.Contract.MinStakeIncrease(&_Contract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Contract *ContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Contract *ContractSession) Owner() (common.Address, error) {
	return _Contract.Contract.Owner(&_Contract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Contract *ContractCallerSession) Owner() (common.Address, error) {
	return _Contract.Contract.Owner(&_Contract.CallOpts)
}

// RewardsAllowed is a free data retrieval call binding the contract method 0x8447c4df.
//
// Solidity: function rewardsAllowed() constant returns(bool)
func (_Contract *ContractCaller) RewardsAllowed(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "rewardsAllowed")
	return *ret0, err
}

// RewardsAllowed is a free data retrieval call binding the contract method 0x8447c4df.
//
// Solidity: function rewardsAllowed() constant returns(bool)
func (_Contract *ContractSession) RewardsAllowed() (bool, error) {
	return _Contract.Contract.RewardsAllowed(&_Contract.CallOpts)
}

// RewardsAllowed is a free data retrieval call binding the contract method 0x8447c4df.
//
// Solidity: function rewardsAllowed() constant returns(bool)
func (_Contract *ContractCallerSession) RewardsAllowed() (bool, error) {
	return _Contract.Contract.RewardsAllowed(&_Contract.CallOpts)
}

// RewardsStash is a free data retrieval call binding the contract method 0x6f498663.
//
// Solidity: function rewardsStash(address , uint256 ) constant returns(uint256 amount)
func (_Contract *ContractCaller) RewardsStash(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "rewardsStash", arg0, arg1)
	return *ret0, err
}

// RewardsStash is a free data retrieval call binding the contract method 0x6f498663.
//
// Solidity: function rewardsStash(address , uint256 ) constant returns(uint256 amount)
func (_Contract *ContractSession) RewardsStash(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _Contract.Contract.RewardsStash(&_Contract.CallOpts, arg0, arg1)
}

// RewardsStash is a free data retrieval call binding the contract method 0x6f498663.
//
// Solidity: function rewardsStash(address , uint256 ) constant returns(uint256 amount)
func (_Contract *ContractCallerSession) RewardsStash(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _Contract.Contract.RewardsStash(&_Contract.CallOpts, arg0, arg1)
}

// SlashedDelegationsTotalAmount is a free data retrieval call binding the contract method 0xa70da4d2.
//
// Solidity: function slashedDelegationsTotalAmount() constant returns(uint256)
func (_Contract *ContractCaller) SlashedDelegationsTotalAmount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "slashedDelegationsTotalAmount")
	return *ret0, err
}

// SlashedDelegationsTotalAmount is a free data retrieval call binding the contract method 0xa70da4d2.
//
// Solidity: function slashedDelegationsTotalAmount() constant returns(uint256)
func (_Contract *ContractSession) SlashedDelegationsTotalAmount() (*big.Int, error) {
	return _Contract.Contract.SlashedDelegationsTotalAmount(&_Contract.CallOpts)
}

// SlashedDelegationsTotalAmount is a free data retrieval call binding the contract method 0xa70da4d2.
//
// Solidity: function slashedDelegationsTotalAmount() constant returns(uint256)
func (_Contract *ContractCallerSession) SlashedDelegationsTotalAmount() (*big.Int, error) {
	return _Contract.Contract.SlashedDelegationsTotalAmount(&_Contract.CallOpts)
}

// SlashedStakeTotalAmount is a free data retrieval call binding the contract method 0x0a29180c.
//
// Solidity: function slashedStakeTotalAmount() constant returns(uint256)
func (_Contract *ContractCaller) SlashedStakeTotalAmount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "slashedStakeTotalAmount")
	return *ret0, err
}

// SlashedStakeTotalAmount is a free data retrieval call binding the contract method 0x0a29180c.
//
// Solidity: function slashedStakeTotalAmount() constant returns(uint256)
func (_Contract *ContractSession) SlashedStakeTotalAmount() (*big.Int, error) {
	return _Contract.Contract.SlashedStakeTotalAmount(&_Contract.CallOpts)
}

// SlashedStakeTotalAmount is a free data retrieval call binding the contract method 0x0a29180c.
//
// Solidity: function slashedStakeTotalAmount() constant returns(uint256)
func (_Contract *ContractCallerSession) SlashedStakeTotalAmount() (*big.Int, error) {
	return _Contract.Contract.SlashedStakeTotalAmount(&_Contract.CallOpts)
}

// StakeLockPeriodEpochs is a free data retrieval call binding the contract method 0x54d77ed2.
//
// Solidity: function stakeLockPeriodEpochs() constant returns(uint256)
func (_Contract *ContractCaller) StakeLockPeriodEpochs(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "stakeLockPeriodEpochs")
	return *ret0, err
}

// StakeLockPeriodEpochs is a free data retrieval call binding the contract method 0x54d77ed2.
//
// Solidity: function stakeLockPeriodEpochs() constant returns(uint256)
func (_Contract *ContractSession) StakeLockPeriodEpochs() (*big.Int, error) {
	return _Contract.Contract.StakeLockPeriodEpochs(&_Contract.CallOpts)
}

// StakeLockPeriodEpochs is a free data retrieval call binding the contract method 0x54d77ed2.
//
// Solidity: function stakeLockPeriodEpochs() constant returns(uint256)
func (_Contract *ContractCallerSession) StakeLockPeriodEpochs() (*big.Int, error) {
	return _Contract.Contract.StakeLockPeriodEpochs(&_Contract.CallOpts)
}

// StakeLockPeriodTime is a free data retrieval call binding the contract method 0x3fee10a8.
//
// Solidity: function stakeLockPeriodTime() constant returns(uint256)
func (_Contract *ContractCaller) StakeLockPeriodTime(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "stakeLockPeriodTime")
	return *ret0, err
}

// StakeLockPeriodTime is a free data retrieval call binding the contract method 0x3fee10a8.
//
// Solidity: function stakeLockPeriodTime() constant returns(uint256)
func (_Contract *ContractSession) StakeLockPeriodTime() (*big.Int, error) {
	return _Contract.Contract.StakeLockPeriodTime(&_Contract.CallOpts)
}

// StakeLockPeriodTime is a free data retrieval call binding the contract method 0x3fee10a8.
//
// Solidity: function stakeLockPeriodTime() constant returns(uint256)
func (_Contract *ContractCallerSession) StakeLockPeriodTime() (*big.Int, error) {
	return _Contract.Contract.StakeLockPeriodTime(&_Contract.CallOpts)
}

// StakeTotalAmount is a free data retrieval call binding the contract method 0x3d0317fe.
//
// Solidity: function stakeTotalAmount() constant returns(uint256)
func (_Contract *ContractCaller) StakeTotalAmount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "stakeTotalAmount")
	return *ret0, err
}

// StakeTotalAmount is a free data retrieval call binding the contract method 0x3d0317fe.
//
// Solidity: function stakeTotalAmount() constant returns(uint256)
func (_Contract *ContractSession) StakeTotalAmount() (*big.Int, error) {
	return _Contract.Contract.StakeTotalAmount(&_Contract.CallOpts)
}

// StakeTotalAmount is a free data retrieval call binding the contract method 0x3d0317fe.
//
// Solidity: function stakeTotalAmount() constant returns(uint256)
func (_Contract *ContractCallerSession) StakeTotalAmount() (*big.Int, error) {
	return _Contract.Contract.StakeTotalAmount(&_Contract.CallOpts)
}

// StakerMetadata is a free data retrieval call binding the contract method 0x98ec2de5.
//
// Solidity: function stakerMetadata(uint256 ) constant returns(bytes)
func (_Contract *ContractCaller) StakerMetadata(opts *bind.CallOpts, arg0 *big.Int) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "stakerMetadata", arg0)
	return *ret0, err
}

// StakerMetadata is a free data retrieval call binding the contract method 0x98ec2de5.
//
// Solidity: function stakerMetadata(uint256 ) constant returns(bytes)
func (_Contract *ContractSession) StakerMetadata(arg0 *big.Int) ([]byte, error) {
	return _Contract.Contract.StakerMetadata(&_Contract.CallOpts, arg0)
}

// StakerMetadata is a free data retrieval call binding the contract method 0x98ec2de5.
//
// Solidity: function stakerMetadata(uint256 ) constant returns(bytes)
func (_Contract *ContractCallerSession) StakerMetadata(arg0 *big.Int) ([]byte, error) {
	return _Contract.Contract.StakerMetadata(&_Contract.CallOpts, arg0)
}

// Stakers is a free data retrieval call binding the contract method 0xfd5e6dd1.
//
// Solidity: function stakers(uint256 ) constant returns(uint256 status, uint256 createdEpoch, uint256 createdTime, uint256 deactivatedEpoch, uint256 deactivatedTime, uint256 stakeAmount, uint256 paidUntilEpoch, uint256 delegatedMe, address dagAddress, address sfcAddress)
func (_Contract *ContractCaller) Stakers(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Status           *big.Int
	CreatedEpoch     *big.Int
	CreatedTime      *big.Int
	DeactivatedEpoch *big.Int
	DeactivatedTime  *big.Int
	StakeAmount      *big.Int
	PaidUntilEpoch   *big.Int
	DelegatedMe      *big.Int
	DagAddress       common.Address
	SfcAddress       common.Address
}, error) {
	ret := new(struct {
		Status           *big.Int
		CreatedEpoch     *big.Int
		CreatedTime      *big.Int
		DeactivatedEpoch *big.Int
		DeactivatedTime  *big.Int
		StakeAmount      *big.Int
		PaidUntilEpoch   *big.Int
		DelegatedMe      *big.Int
		DagAddress       common.Address
		SfcAddress       common.Address
	})
	out := ret
	err := _Contract.contract.Call(opts, out, "stakers", arg0)
	return *ret, err
}

// Stakers is a free data retrieval call binding the contract method 0xfd5e6dd1.
//
// Solidity: function stakers(uint256 ) constant returns(uint256 status, uint256 createdEpoch, uint256 createdTime, uint256 deactivatedEpoch, uint256 deactivatedTime, uint256 stakeAmount, uint256 paidUntilEpoch, uint256 delegatedMe, address dagAddress, address sfcAddress)
func (_Contract *ContractSession) Stakers(arg0 *big.Int) (struct {
	Status           *big.Int
	CreatedEpoch     *big.Int
	CreatedTime      *big.Int
	DeactivatedEpoch *big.Int
	DeactivatedTime  *big.Int
	StakeAmount      *big.Int
	PaidUntilEpoch   *big.Int
	DelegatedMe      *big.Int
	DagAddress       common.Address
	SfcAddress       common.Address
}, error) {
	return _Contract.Contract.Stakers(&_Contract.CallOpts, arg0)
}

// Stakers is a free data retrieval call binding the contract method 0xfd5e6dd1.
//
// Solidity: function stakers(uint256 ) constant returns(uint256 status, uint256 createdEpoch, uint256 createdTime, uint256 deactivatedEpoch, uint256 deactivatedTime, uint256 stakeAmount, uint256 paidUntilEpoch, uint256 delegatedMe, address dagAddress, address sfcAddress)
func (_Contract *ContractCallerSession) Stakers(arg0 *big.Int) (struct {
	Status           *big.Int
	CreatedEpoch     *big.Int
	CreatedTime      *big.Int
	DeactivatedEpoch *big.Int
	DeactivatedTime  *big.Int
	StakeAmount      *big.Int
	PaidUntilEpoch   *big.Int
	DelegatedMe      *big.Int
	DagAddress       common.Address
	SfcAddress       common.Address
}, error) {
	return _Contract.Contract.Stakers(&_Contract.CallOpts, arg0)
}

// StakersLastID is a free data retrieval call binding the contract method 0x81d9dc7a.
//
// Solidity: function stakersLastID() constant returns(uint256)
func (_Contract *ContractCaller) StakersLastID(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "stakersLastID")
	return *ret0, err
}

// StakersLastID is a free data retrieval call binding the contract method 0x81d9dc7a.
//
// Solidity: function stakersLastID() constant returns(uint256)
func (_Contract *ContractSession) StakersLastID() (*big.Int, error) {
	return _Contract.Contract.StakersLastID(&_Contract.CallOpts)
}

// StakersLastID is a free data retrieval call binding the contract method 0x81d9dc7a.
//
// Solidity: function stakersLastID() constant returns(uint256)
func (_Contract *ContractCallerSession) StakersLastID() (*big.Int, error) {
	return _Contract.Contract.StakersLastID(&_Contract.CallOpts)
}

// StakersNum is a free data retrieval call binding the contract method 0x08728f6e.
//
// Solidity: function stakersNum() constant returns(uint256)
func (_Contract *ContractCaller) StakersNum(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "stakersNum")
	return *ret0, err
}

// StakersNum is a free data retrieval call binding the contract method 0x08728f6e.
//
// Solidity: function stakersNum() constant returns(uint256)
func (_Contract *ContractSession) StakersNum() (*big.Int, error) {
	return _Contract.Contract.StakersNum(&_Contract.CallOpts)
}

// StakersNum is a free data retrieval call binding the contract method 0x08728f6e.
//
// Solidity: function stakersNum() constant returns(uint256)
func (_Contract *ContractCallerSession) StakersNum() (*big.Int, error) {
	return _Contract.Contract.StakersNum(&_Contract.CallOpts)
}

// UnbondingStartDate is a free data retrieval call binding the contract method 0x53a72586.
//
// Solidity: function unbondingStartDate() constant returns(uint256)
func (_Contract *ContractCaller) UnbondingStartDate(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "unbondingStartDate")
	return *ret0, err
}

// UnbondingStartDate is a free data retrieval call binding the contract method 0x53a72586.
//
// Solidity: function unbondingStartDate() constant returns(uint256)
func (_Contract *ContractSession) UnbondingStartDate() (*big.Int, error) {
	return _Contract.Contract.UnbondingStartDate(&_Contract.CallOpts)
}

// UnbondingStartDate is a free data retrieval call binding the contract method 0x53a72586.
//
// Solidity: function unbondingStartDate() constant returns(uint256)
func (_Contract *ContractCallerSession) UnbondingStartDate() (*big.Int, error) {
	return _Contract.Contract.UnbondingStartDate(&_Contract.CallOpts)
}

// UnbondingUnlockPeriod is a free data retrieval call binding the contract method 0x3a0af4d4.
//
// Solidity: function unbondingUnlockPeriod() constant returns(uint256)
func (_Contract *ContractCaller) UnbondingUnlockPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "unbondingUnlockPeriod")
	return *ret0, err
}

// UnbondingUnlockPeriod is a free data retrieval call binding the contract method 0x3a0af4d4.
//
// Solidity: function unbondingUnlockPeriod() constant returns(uint256)
func (_Contract *ContractSession) UnbondingUnlockPeriod() (*big.Int, error) {
	return _Contract.Contract.UnbondingUnlockPeriod(&_Contract.CallOpts)
}

// UnbondingUnlockPeriod is a free data retrieval call binding the contract method 0x3a0af4d4.
//
// Solidity: function unbondingUnlockPeriod() constant returns(uint256)
func (_Contract *ContractCallerSession) UnbondingUnlockPeriod() (*big.Int, error) {
	return _Contract.Contract.UnbondingUnlockPeriod(&_Contract.CallOpts)
}

// UnlockedRatio is a free data retrieval call binding the contract method 0x65cca35d.
//
// Solidity: function unlockedRatio() constant returns(uint256)
func (_Contract *ContractCaller) UnlockedRatio(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "unlockedRatio")
	return *ret0, err
}

// UnlockedRatio is a free data retrieval call binding the contract method 0x65cca35d.
//
// Solidity: function unlockedRatio() constant returns(uint256)
func (_Contract *ContractSession) UnlockedRatio() (*big.Int, error) {
	return _Contract.Contract.UnlockedRatio(&_Contract.CallOpts)
}

// UnlockedRatio is a free data retrieval call binding the contract method 0x65cca35d.
//
// Solidity: function unlockedRatio() constant returns(uint256)
func (_Contract *ContractCallerSession) UnlockedRatio() (*big.Int, error) {
	return _Contract.Contract.UnlockedRatio(&_Contract.CallOpts)
}

// ValidatorCommission is a free data retrieval call binding the contract method 0xa7786515.
//
// Solidity: function validatorCommission() constant returns(uint256)
func (_Contract *ContractCaller) ValidatorCommission(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contract.contract.Call(opts, out, "validatorCommission")
	return *ret0, err
}

// ValidatorCommission is a free data retrieval call binding the contract method 0xa7786515.
//
// Solidity: function validatorCommission() constant returns(uint256)
func (_Contract *ContractSession) ValidatorCommission() (*big.Int, error) {
	return _Contract.Contract.ValidatorCommission(&_Contract.CallOpts)
}

// ValidatorCommission is a free data retrieval call binding the contract method 0xa7786515.
//
// Solidity: function validatorCommission() constant returns(uint256)
func (_Contract *ContractCallerSession) ValidatorCommission() (*big.Int, error) {
	return _Contract.Contract.ValidatorCommission(&_Contract.CallOpts)
}

// WithdrawalRequests is a free data retrieval call binding the contract method 0x4e5a2328.
//
// Solidity: function withdrawalRequests(address , uint256 ) constant returns(uint256 stakerID, uint256 epoch, uint256 time, uint256 amount, bool delegation)
func (_Contract *ContractCaller) WithdrawalRequests(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (struct {
	StakerID   *big.Int
	Epoch      *big.Int
	Time       *big.Int
	Amount     *big.Int
	Delegation bool
}, error) {
	ret := new(struct {
		StakerID   *big.Int
		Epoch      *big.Int
		Time       *big.Int
		Amount     *big.Int
		Delegation bool
	})
	out := ret
	err := _Contract.contract.Call(opts, out, "withdrawalRequests", arg0, arg1)
	return *ret, err
}

// WithdrawalRequests is a free data retrieval call binding the contract method 0x4e5a2328.
//
// Solidity: function withdrawalRequests(address , uint256 ) constant returns(uint256 stakerID, uint256 epoch, uint256 time, uint256 amount, bool delegation)
func (_Contract *ContractSession) WithdrawalRequests(arg0 common.Address, arg1 *big.Int) (struct {
	StakerID   *big.Int
	Epoch      *big.Int
	Time       *big.Int
	Amount     *big.Int
	Delegation bool
}, error) {
	return _Contract.Contract.WithdrawalRequests(&_Contract.CallOpts, arg0, arg1)
}

// WithdrawalRequests is a free data retrieval call binding the contract method 0x4e5a2328.
//
// Solidity: function withdrawalRequests(address , uint256 ) constant returns(uint256 stakerID, uint256 epoch, uint256 time, uint256 amount, bool delegation)
func (_Contract *ContractCallerSession) WithdrawalRequests(arg0 common.Address, arg1 *big.Int) (struct {
	StakerID   *big.Int
	Epoch      *big.Int
	Time       *big.Int
	Amount     *big.Int
	Delegation bool
}, error) {
	return _Contract.Contract.WithdrawalRequests(&_Contract.CallOpts, arg0, arg1)
}

// SyncDelegator is a paid mutator transaction binding the contract method 0x5dc03f1f.
//
// Solidity: function _syncDelegator(address delegator, uint256 stakerID) returns()
func (_Contract *ContractTransactor) SyncDelegator(opts *bind.TransactOpts, delegator common.Address, stakerID *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "_syncDelegator", delegator, stakerID)
}

// SyncDelegator is a paid mutator transaction binding the contract method 0x5dc03f1f.
//
// Solidity: function _syncDelegator(address delegator, uint256 stakerID) returns()
func (_Contract *ContractSession) SyncDelegator(delegator common.Address, stakerID *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SyncDelegator(&_Contract.TransactOpts, delegator, stakerID)
}

// SyncDelegator is a paid mutator transaction binding the contract method 0x5dc03f1f.
//
// Solidity: function _syncDelegator(address delegator, uint256 stakerID) returns()
func (_Contract *ContractTransactorSession) SyncDelegator(delegator common.Address, stakerID *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SyncDelegator(&_Contract.TransactOpts, delegator, stakerID)
}

// SyncStaker is a paid mutator transaction binding the contract method 0xeac3baf2.
//
// Solidity: function _syncStaker(uint256 stakerID) returns()
func (_Contract *ContractTransactor) SyncStaker(opts *bind.TransactOpts, stakerID *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "_syncStaker", stakerID)
}

// SyncStaker is a paid mutator transaction binding the contract method 0xeac3baf2.
//
// Solidity: function _syncStaker(uint256 stakerID) returns()
func (_Contract *ContractSession) SyncStaker(stakerID *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SyncStaker(&_Contract.TransactOpts, stakerID)
}

// SyncStaker is a paid mutator transaction binding the contract method 0xeac3baf2.
//
// Solidity: function _syncStaker(uint256 stakerID) returns()
func (_Contract *ContractTransactorSession) SyncStaker(stakerID *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SyncStaker(&_Contract.TransactOpts, stakerID)
}

// UpgradeStakerStorage is a paid mutator transaction binding the contract method 0x28dca8ff.
//
// Solidity: function _upgradeStakerStorage(uint256 stakerID) returns()
func (_Contract *ContractTransactor) UpgradeStakerStorage(opts *bind.TransactOpts, stakerID *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "_upgradeStakerStorage", stakerID)
}

// UpgradeStakerStorage is a paid mutator transaction binding the contract method 0x28dca8ff.
//
// Solidity: function _upgradeStakerStorage(uint256 stakerID) returns()
func (_Contract *ContractSession) UpgradeStakerStorage(stakerID *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.UpgradeStakerStorage(&_Contract.TransactOpts, stakerID)
}

// UpgradeStakerStorage is a paid mutator transaction binding the contract method 0x28dca8ff.
//
// Solidity: function _upgradeStakerStorage(uint256 stakerID) returns()
func (_Contract *ContractTransactorSession) UpgradeStakerStorage(stakerID *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.UpgradeStakerStorage(&_Contract.TransactOpts, stakerID)
}

// ClaimDelegationRewards is a paid mutator transaction binding the contract method 0xf99837e6.
//
// Solidity: function claimDelegationRewards(uint256 maxEpochs, uint256 stakerID) returns()
func (_Contract *ContractTransactor) ClaimDelegationRewards(opts *bind.TransactOpts, maxEpochs *big.Int, stakerID *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "claimDelegationRewards", maxEpochs, stakerID)
}

// ClaimDelegationRewards is a paid mutator transaction binding the contract method 0xf99837e6.
//
// Solidity: function claimDelegationRewards(uint256 maxEpochs, uint256 stakerID) returns()
func (_Contract *ContractSession) ClaimDelegationRewards(maxEpochs *big.Int, stakerID *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.ClaimDelegationRewards(&_Contract.TransactOpts, maxEpochs, stakerID)
}

// ClaimDelegationRewards is a paid mutator transaction binding the contract method 0xf99837e6.
//
// Solidity: function claimDelegationRewards(uint256 maxEpochs, uint256 stakerID) returns()
func (_Contract *ContractTransactorSession) ClaimDelegationRewards(maxEpochs *big.Int, stakerID *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.ClaimDelegationRewards(&_Contract.TransactOpts, maxEpochs, stakerID)
}

// ClaimValidatorRewards is a paid mutator transaction binding the contract method 0x295cccba.
//
// Solidity: function claimValidatorRewards(uint256 maxEpochs) returns()
func (_Contract *ContractTransactor) ClaimValidatorRewards(opts *bind.TransactOpts, maxEpochs *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "claimValidatorRewards", maxEpochs)
}

// ClaimValidatorRewards is a paid mutator transaction binding the contract method 0x295cccba.
//
// Solidity: function claimValidatorRewards(uint256 maxEpochs) returns()
func (_Contract *ContractSession) ClaimValidatorRewards(maxEpochs *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.ClaimValidatorRewards(&_Contract.TransactOpts, maxEpochs)
}

// ClaimValidatorRewards is a paid mutator transaction binding the contract method 0x295cccba.
//
// Solidity: function claimValidatorRewards(uint256 maxEpochs) returns()
func (_Contract *ContractTransactorSession) ClaimValidatorRewards(maxEpochs *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.ClaimValidatorRewards(&_Contract.TransactOpts, maxEpochs)
}

// CreateDelegation is a paid mutator transaction binding the contract method 0xc312eb07.
//
// Solidity: function createDelegation(uint256 to) returns()
func (_Contract *ContractTransactor) CreateDelegation(opts *bind.TransactOpts, to *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "createDelegation", to)
}

// CreateDelegation is a paid mutator transaction binding the contract method 0xc312eb07.
//
// Solidity: function createDelegation(uint256 to) returns()
func (_Contract *ContractSession) CreateDelegation(to *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.CreateDelegation(&_Contract.TransactOpts, to)
}

// CreateDelegation is a paid mutator transaction binding the contract method 0xc312eb07.
//
// Solidity: function createDelegation(uint256 to) returns()
func (_Contract *ContractTransactorSession) CreateDelegation(to *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.CreateDelegation(&_Contract.TransactOpts, to)
}

// CreateStake is a paid mutator transaction binding the contract method 0xcc8c2120.
//
// Solidity: function createStake(bytes metadata) returns()
func (_Contract *ContractTransactor) CreateStake(opts *bind.TransactOpts, metadata []byte) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "createStake", metadata)
}

// CreateStake is a paid mutator transaction binding the contract method 0xcc8c2120.
//
// Solidity: function createStake(bytes metadata) returns()
func (_Contract *ContractSession) CreateStake(metadata []byte) (*types.Transaction, error) {
	return _Contract.Contract.CreateStake(&_Contract.TransactOpts, metadata)
}

// CreateStake is a paid mutator transaction binding the contract method 0xcc8c2120.
//
// Solidity: function createStake(bytes metadata) returns()
func (_Contract *ContractTransactorSession) CreateStake(metadata []byte) (*types.Transaction, error) {
	return _Contract.Contract.CreateStake(&_Contract.TransactOpts, metadata)
}

// CreateStakeWithAddresses is a paid mutator transaction binding the contract method 0x90475ae4.
//
// Solidity: function createStakeWithAddresses(address dagAdrress, address sfcAddress, bytes metadata) returns()
func (_Contract *ContractTransactor) CreateStakeWithAddresses(opts *bind.TransactOpts, dagAdrress common.Address, sfcAddress common.Address, metadata []byte) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "createStakeWithAddresses", dagAdrress, sfcAddress, metadata)
}

// CreateStakeWithAddresses is a paid mutator transaction binding the contract method 0x90475ae4.
//
// Solidity: function createStakeWithAddresses(address dagAdrress, address sfcAddress, bytes metadata) returns()
func (_Contract *ContractSession) CreateStakeWithAddresses(dagAdrress common.Address, sfcAddress common.Address, metadata []byte) (*types.Transaction, error) {
	return _Contract.Contract.CreateStakeWithAddresses(&_Contract.TransactOpts, dagAdrress, sfcAddress, metadata)
}

// CreateStakeWithAddresses is a paid mutator transaction binding the contract method 0x90475ae4.
//
// Solidity: function createStakeWithAddresses(address dagAdrress, address sfcAddress, bytes metadata) returns()
func (_Contract *ContractTransactorSession) CreateStakeWithAddresses(dagAdrress common.Address, sfcAddress common.Address, metadata []byte) (*types.Transaction, error) {
	return _Contract.Contract.CreateStakeWithAddresses(&_Contract.TransactOpts, dagAdrress, sfcAddress, metadata)
}

// IncreaseDelegation is a paid mutator transaction binding the contract method 0xdbd3aa8a.
//
// Solidity: function increaseDelegation(uint256 to) returns()
func (_Contract *ContractTransactor) IncreaseDelegation(opts *bind.TransactOpts, to *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "increaseDelegation", to)
}

// IncreaseDelegation is a paid mutator transaction binding the contract method 0xdbd3aa8a.
//
// Solidity: function increaseDelegation(uint256 to) returns()
func (_Contract *ContractSession) IncreaseDelegation(to *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.IncreaseDelegation(&_Contract.TransactOpts, to)
}

// IncreaseDelegation is a paid mutator transaction binding the contract method 0xdbd3aa8a.
//
// Solidity: function increaseDelegation(uint256 to) returns()
func (_Contract *ContractTransactorSession) IncreaseDelegation(to *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.IncreaseDelegation(&_Contract.TransactOpts, to)
}

// IncreaseStake is a paid mutator transaction binding the contract method 0xd9e257ef.
//
// Solidity: function increaseStake() returns()
func (_Contract *ContractTransactor) IncreaseStake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "increaseStake")
}

// IncreaseStake is a paid mutator transaction binding the contract method 0xd9e257ef.
//
// Solidity: function increaseStake() returns()
func (_Contract *ContractSession) IncreaseStake() (*types.Transaction, error) {
	return _Contract.Contract.IncreaseStake(&_Contract.TransactOpts)
}

// IncreaseStake is a paid mutator transaction binding the contract method 0xd9e257ef.
//
// Solidity: function increaseStake() returns()
func (_Contract *ContractTransactorSession) IncreaseStake() (*types.Transaction, error) {
	return _Contract.Contract.IncreaseStake(&_Contract.TransactOpts)
}

// LockUpDelegation is a paid mutator transaction binding the contract method 0xa4b89fab.
//
// Solidity: function lockUpDelegation(uint256 lockDuration, uint256 stakerID) returns()
func (_Contract *ContractTransactor) LockUpDelegation(opts *bind.TransactOpts, lockDuration *big.Int, stakerID *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "lockUpDelegation", lockDuration, stakerID)
}

// LockUpDelegation is a paid mutator transaction binding the contract method 0xa4b89fab.
//
// Solidity: function lockUpDelegation(uint256 lockDuration, uint256 stakerID) returns()
func (_Contract *ContractSession) LockUpDelegation(lockDuration *big.Int, stakerID *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.LockUpDelegation(&_Contract.TransactOpts, lockDuration, stakerID)
}

// LockUpDelegation is a paid mutator transaction binding the contract method 0xa4b89fab.
//
// Solidity: function lockUpDelegation(uint256 lockDuration, uint256 stakerID) returns()
func (_Contract *ContractTransactorSession) LockUpDelegation(lockDuration *big.Int, stakerID *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.LockUpDelegation(&_Contract.TransactOpts, lockDuration, stakerID)
}

// LockUpStake is a paid mutator transaction binding the contract method 0xf3ae5b1a.
//
// Solidity: function lockUpStake(uint256 lockDuration) returns()
func (_Contract *ContractTransactor) LockUpStake(opts *bind.TransactOpts, lockDuration *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "lockUpStake", lockDuration)
}

// LockUpStake is a paid mutator transaction binding the contract method 0xf3ae5b1a.
//
// Solidity: function lockUpStake(uint256 lockDuration) returns()
func (_Contract *ContractSession) LockUpStake(lockDuration *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.LockUpStake(&_Contract.TransactOpts, lockDuration)
}

// LockUpStake is a paid mutator transaction binding the contract method 0xf3ae5b1a.
//
// Solidity: function lockUpStake(uint256 lockDuration) returns()
func (_Contract *ContractTransactorSession) LockUpStake(lockDuration *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.LockUpStake(&_Contract.TransactOpts, lockDuration)
}

// PartialWithdrawByRequest is a paid mutator transaction binding the contract method 0xf8b18d8a.
//
// Solidity: function partialWithdrawByRequest(uint256 wrID) returns()
func (_Contract *ContractTransactor) PartialWithdrawByRequest(opts *bind.TransactOpts, wrID *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "partialWithdrawByRequest", wrID)
}

// PartialWithdrawByRequest is a paid mutator transaction binding the contract method 0xf8b18d8a.
//
// Solidity: function partialWithdrawByRequest(uint256 wrID) returns()
func (_Contract *ContractSession) PartialWithdrawByRequest(wrID *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.PartialWithdrawByRequest(&_Contract.TransactOpts, wrID)
}

// PartialWithdrawByRequest is a paid mutator transaction binding the contract method 0xf8b18d8a.
//
// Solidity: function partialWithdrawByRequest(uint256 wrID) returns()
func (_Contract *ContractTransactorSession) PartialWithdrawByRequest(wrID *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.PartialWithdrawByRequest(&_Contract.TransactOpts, wrID)
}

// PrepareToWithdrawDelegation is a paid mutator transaction binding the contract method 0xb1e64339.
//
// Solidity: function prepareToWithdrawDelegation(uint256 stakerID) returns()
func (_Contract *ContractTransactor) PrepareToWithdrawDelegation(opts *bind.TransactOpts, stakerID *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "prepareToWithdrawDelegation", stakerID)
}

// PrepareToWithdrawDelegation is a paid mutator transaction binding the contract method 0xb1e64339.
//
// Solidity: function prepareToWithdrawDelegation(uint256 stakerID) returns()
func (_Contract *ContractSession) PrepareToWithdrawDelegation(stakerID *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.PrepareToWithdrawDelegation(&_Contract.TransactOpts, stakerID)
}

// PrepareToWithdrawDelegation is a paid mutator transaction binding the contract method 0xb1e64339.
//
// Solidity: function prepareToWithdrawDelegation(uint256 stakerID) returns()
func (_Contract *ContractTransactorSession) PrepareToWithdrawDelegation(stakerID *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.PrepareToWithdrawDelegation(&_Contract.TransactOpts, stakerID)
}

// PrepareToWithdrawDelegationPartial is a paid mutator transaction binding the contract method 0xbb03a4bd.
//
// Solidity: function prepareToWithdrawDelegationPartial(uint256 wrID, uint256 stakerID, uint256 amount) returns()
func (_Contract *ContractTransactor) PrepareToWithdrawDelegationPartial(opts *bind.TransactOpts, wrID *big.Int, stakerID *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "prepareToWithdrawDelegationPartial", wrID, stakerID, amount)
}

// PrepareToWithdrawDelegationPartial is a paid mutator transaction binding the contract method 0xbb03a4bd.
//
// Solidity: function prepareToWithdrawDelegationPartial(uint256 wrID, uint256 stakerID, uint256 amount) returns()
func (_Contract *ContractSession) PrepareToWithdrawDelegationPartial(wrID *big.Int, stakerID *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.PrepareToWithdrawDelegationPartial(&_Contract.TransactOpts, wrID, stakerID, amount)
}

// PrepareToWithdrawDelegationPartial is a paid mutator transaction binding the contract method 0xbb03a4bd.
//
// Solidity: function prepareToWithdrawDelegationPartial(uint256 wrID, uint256 stakerID, uint256 amount) returns()
func (_Contract *ContractTransactorSession) PrepareToWithdrawDelegationPartial(wrID *big.Int, stakerID *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.PrepareToWithdrawDelegationPartial(&_Contract.TransactOpts, wrID, stakerID, amount)
}

// PrepareToWithdrawStake is a paid mutator transaction binding the contract method 0xc41b6405.
//
// Solidity: function prepareToWithdrawStake() returns()
func (_Contract *ContractTransactor) PrepareToWithdrawStake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "prepareToWithdrawStake")
}

// PrepareToWithdrawStake is a paid mutator transaction binding the contract method 0xc41b6405.
//
// Solidity: function prepareToWithdrawStake() returns()
func (_Contract *ContractSession) PrepareToWithdrawStake() (*types.Transaction, error) {
	return _Contract.Contract.PrepareToWithdrawStake(&_Contract.TransactOpts)
}

// PrepareToWithdrawStake is a paid mutator transaction binding the contract method 0xc41b6405.
//
// Solidity: function prepareToWithdrawStake() returns()
func (_Contract *ContractTransactorSession) PrepareToWithdrawStake() (*types.Transaction, error) {
	return _Contract.Contract.PrepareToWithdrawStake(&_Contract.TransactOpts)
}

// PrepareToWithdrawStakePartial is a paid mutator transaction binding the contract method 0x26682c71.
//
// Solidity: function prepareToWithdrawStakePartial(uint256 wrID, uint256 amount) returns()
func (_Contract *ContractTransactor) PrepareToWithdrawStakePartial(opts *bind.TransactOpts, wrID *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "prepareToWithdrawStakePartial", wrID, amount)
}

// PrepareToWithdrawStakePartial is a paid mutator transaction binding the contract method 0x26682c71.
//
// Solidity: function prepareToWithdrawStakePartial(uint256 wrID, uint256 amount) returns()
func (_Contract *ContractSession) PrepareToWithdrawStakePartial(wrID *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.PrepareToWithdrawStakePartial(&_Contract.TransactOpts, wrID, amount)
}

// PrepareToWithdrawStakePartial is a paid mutator transaction binding the contract method 0x26682c71.
//
// Solidity: function prepareToWithdrawStakePartial(uint256 wrID, uint256 amount) returns()
func (_Contract *ContractTransactorSession) PrepareToWithdrawStakePartial(wrID *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.PrepareToWithdrawStakePartial(&_Contract.TransactOpts, wrID, amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Contract *ContractTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Contract *ContractSession) RenounceOwnership() (*types.Transaction, error) {
	return _Contract.Contract.RenounceOwnership(&_Contract.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Contract *ContractTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Contract.Contract.RenounceOwnership(&_Contract.TransactOpts)
}

// StartLockedUp is a paid mutator transaction binding the contract method 0xc9400d4f.
//
// Solidity: function startLockedUp(uint256 epochNum) returns()
func (_Contract *ContractTransactor) StartLockedUp(opts *bind.TransactOpts, epochNum *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "startLockedUp", epochNum)
}

// StartLockedUp is a paid mutator transaction binding the contract method 0xc9400d4f.
//
// Solidity: function startLockedUp(uint256 epochNum) returns()
func (_Contract *ContractSession) StartLockedUp(epochNum *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.StartLockedUp(&_Contract.TransactOpts, epochNum)
}

// StartLockedUp is a paid mutator transaction binding the contract method 0xc9400d4f.
//
// Solidity: function startLockedUp(uint256 epochNum) returns()
func (_Contract *ContractTransactorSession) StartLockedUp(epochNum *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.StartLockedUp(&_Contract.TransactOpts, epochNum)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Contract *ContractTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Contract *ContractSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Contract.Contract.TransferOwnership(&_Contract.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Contract *ContractTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Contract.Contract.TransferOwnership(&_Contract.TransactOpts, newOwner)
}

// UnstashRewards is a paid mutator transaction binding the contract method 0x876f7e2a.
//
// Solidity: function unstashRewards() returns()
func (_Contract *ContractTransactor) UnstashRewards(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "unstashRewards")
}

// UnstashRewards is a paid mutator transaction binding the contract method 0x876f7e2a.
//
// Solidity: function unstashRewards() returns()
func (_Contract *ContractSession) UnstashRewards() (*types.Transaction, error) {
	return _Contract.Contract.UnstashRewards(&_Contract.TransactOpts)
}

// UnstashRewards is a paid mutator transaction binding the contract method 0x876f7e2a.
//
// Solidity: function unstashRewards() returns()
func (_Contract *ContractTransactorSession) UnstashRewards() (*types.Transaction, error) {
	return _Contract.Contract.UnstashRewards(&_Contract.TransactOpts)
}

// UpdateBaseRewardPerSec is a paid mutator transaction binding the contract method 0x1b593d8a.
//
// Solidity: function updateBaseRewardPerSec(uint256 value) returns()
func (_Contract *ContractTransactor) UpdateBaseRewardPerSec(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "updateBaseRewardPerSec", value)
}

// UpdateBaseRewardPerSec is a paid mutator transaction binding the contract method 0x1b593d8a.
//
// Solidity: function updateBaseRewardPerSec(uint256 value) returns()
func (_Contract *ContractSession) UpdateBaseRewardPerSec(value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.UpdateBaseRewardPerSec(&_Contract.TransactOpts, value)
}

// UpdateBaseRewardPerSec is a paid mutator transaction binding the contract method 0x1b593d8a.
//
// Solidity: function updateBaseRewardPerSec(uint256 value) returns()
func (_Contract *ContractTransactorSession) UpdateBaseRewardPerSec(value *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.UpdateBaseRewardPerSec(&_Contract.TransactOpts, value)
}

// UpdateGasPowerAllocationRate is a paid mutator transaction binding the contract method 0x119e351a.
//
// Solidity: function updateGasPowerAllocationRate(uint256 short, uint256 long) returns()
func (_Contract *ContractTransactor) UpdateGasPowerAllocationRate(opts *bind.TransactOpts, short *big.Int, long *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "updateGasPowerAllocationRate", short, long)
}

// UpdateGasPowerAllocationRate is a paid mutator transaction binding the contract method 0x119e351a.
//
// Solidity: function updateGasPowerAllocationRate(uint256 short, uint256 long) returns()
func (_Contract *ContractSession) UpdateGasPowerAllocationRate(short *big.Int, long *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.UpdateGasPowerAllocationRate(&_Contract.TransactOpts, short, long)
}

// UpdateGasPowerAllocationRate is a paid mutator transaction binding the contract method 0x119e351a.
//
// Solidity: function updateGasPowerAllocationRate(uint256 short, uint256 long) returns()
func (_Contract *ContractTransactorSession) UpdateGasPowerAllocationRate(short *big.Int, long *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.UpdateGasPowerAllocationRate(&_Contract.TransactOpts, short, long)
}

// UpdateStakerMetadata is a paid mutator transaction binding the contract method 0x33a14912.
//
// Solidity: function updateStakerMetadata(bytes metadata) returns()
func (_Contract *ContractTransactor) UpdateStakerMetadata(opts *bind.TransactOpts, metadata []byte) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "updateStakerMetadata", metadata)
}

// UpdateStakerMetadata is a paid mutator transaction binding the contract method 0x33a14912.
//
// Solidity: function updateStakerMetadata(bytes metadata) returns()
func (_Contract *ContractSession) UpdateStakerMetadata(metadata []byte) (*types.Transaction, error) {
	return _Contract.Contract.UpdateStakerMetadata(&_Contract.TransactOpts, metadata)
}

// UpdateStakerMetadata is a paid mutator transaction binding the contract method 0x33a14912.
//
// Solidity: function updateStakerMetadata(bytes metadata) returns()
func (_Contract *ContractTransactorSession) UpdateStakerMetadata(metadata []byte) (*types.Transaction, error) {
	return _Contract.Contract.UpdateStakerMetadata(&_Contract.TransactOpts, metadata)
}

// UpdateStakerSfcAddress is a paid mutator transaction binding the contract method 0xc3d74f1a.
//
// Solidity: function updateStakerSfcAddress(address newSfcAddress) returns()
func (_Contract *ContractTransactor) UpdateStakerSfcAddress(opts *bind.TransactOpts, newSfcAddress common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "updateStakerSfcAddress", newSfcAddress)
}

// UpdateStakerSfcAddress is a paid mutator transaction binding the contract method 0xc3d74f1a.
//
// Solidity: function updateStakerSfcAddress(address newSfcAddress) returns()
func (_Contract *ContractSession) UpdateStakerSfcAddress(newSfcAddress common.Address) (*types.Transaction, error) {
	return _Contract.Contract.UpdateStakerSfcAddress(&_Contract.TransactOpts, newSfcAddress)
}

// UpdateStakerSfcAddress is a paid mutator transaction binding the contract method 0xc3d74f1a.
//
// Solidity: function updateStakerSfcAddress(address newSfcAddress) returns()
func (_Contract *ContractTransactorSession) UpdateStakerSfcAddress(newSfcAddress common.Address) (*types.Transaction, error) {
	return _Contract.Contract.UpdateStakerSfcAddress(&_Contract.TransactOpts, newSfcAddress)
}

// WithdrawDelegation is a paid mutator transaction binding the contract method 0xdf0e307a.
//
// Solidity: function withdrawDelegation(uint256 stakerID) returns()
func (_Contract *ContractTransactor) WithdrawDelegation(opts *bind.TransactOpts, stakerID *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "withdrawDelegation", stakerID)
}

// WithdrawDelegation is a paid mutator transaction binding the contract method 0xdf0e307a.
//
// Solidity: function withdrawDelegation(uint256 stakerID) returns()
func (_Contract *ContractSession) WithdrawDelegation(stakerID *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.WithdrawDelegation(&_Contract.TransactOpts, stakerID)
}

// WithdrawDelegation is a paid mutator transaction binding the contract method 0xdf0e307a.
//
// Solidity: function withdrawDelegation(uint256 stakerID) returns()
func (_Contract *ContractTransactorSession) WithdrawDelegation(stakerID *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.WithdrawDelegation(&_Contract.TransactOpts, stakerID)
}

// WithdrawStake is a paid mutator transaction binding the contract method 0xbed9d861.
//
// Solidity: function withdrawStake() returns()
func (_Contract *ContractTransactor) WithdrawStake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "withdrawStake")
}

// WithdrawStake is a paid mutator transaction binding the contract method 0xbed9d861.
//
// Solidity: function withdrawStake() returns()
func (_Contract *ContractSession) WithdrawStake() (*types.Transaction, error) {
	return _Contract.Contract.WithdrawStake(&_Contract.TransactOpts)
}

// WithdrawStake is a paid mutator transaction binding the contract method 0xbed9d861.
//
// Solidity: function withdrawStake() returns()
func (_Contract *ContractTransactorSession) WithdrawStake() (*types.Transaction, error) {
	return _Contract.Contract.WithdrawStake(&_Contract.TransactOpts)
}

// ContractBurntRewardStashIterator is returned from FilterBurntRewardStash and is used to iterate over the raw logs and unpacked data for BurntRewardStash events raised by the Contract contract.
type ContractBurntRewardStashIterator struct {
	Event *ContractBurntRewardStash // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractBurntRewardStashIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractBurntRewardStash)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractBurntRewardStash)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractBurntRewardStashIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractBurntRewardStashIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractBurntRewardStash represents a BurntRewardStash event raised by the Contract contract.
type ContractBurntRewardStash struct {
	Addr         common.Address
	StakerID     *big.Int
	IsDelegation bool
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterBurntRewardStash is a free log retrieval operation binding the contract event 0x0ea92567e76d40ddc52d2c1d74a521a59329a38b50411451de6ad2e565466d0f.
//
// Solidity: event BurntRewardStash(address indexed addr, uint256 indexed stakerID, bool isDelegation, uint256 amount)
func (_Contract *ContractFilterer) FilterBurntRewardStash(opts *bind.FilterOpts, addr []common.Address, stakerID []*big.Int) (*ContractBurntRewardStashIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "BurntRewardStash", addrRule, stakerIDRule)
	if err != nil {
		return nil, err
	}
	return &ContractBurntRewardStashIterator{contract: _Contract.contract, event: "BurntRewardStash", logs: logs, sub: sub}, nil
}

// WatchBurntRewardStash is a free log subscription operation binding the contract event 0x0ea92567e76d40ddc52d2c1d74a521a59329a38b50411451de6ad2e565466d0f.
//
// Solidity: event BurntRewardStash(address indexed addr, uint256 indexed stakerID, bool isDelegation, uint256 amount)
func (_Contract *ContractFilterer) WatchBurntRewardStash(opts *bind.WatchOpts, sink chan<- *ContractBurntRewardStash, addr []common.Address, stakerID []*big.Int) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}
	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "BurntRewardStash", addrRule, stakerIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractBurntRewardStash)
				if err := _Contract.contract.UnpackLog(event, "BurntRewardStash", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBurntRewardStash is a log parse operation binding the contract event 0x0ea92567e76d40ddc52d2c1d74a521a59329a38b50411451de6ad2e565466d0f.
//
// Solidity: event BurntRewardStash(address indexed addr, uint256 indexed stakerID, bool isDelegation, uint256 amount)
func (_Contract *ContractFilterer) ParseBurntRewardStash(log types.Log) (*ContractBurntRewardStash, error) {
	event := new(ContractBurntRewardStash)
	if err := _Contract.contract.UnpackLog(event, "BurntRewardStash", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractClaimedDelegationRewardIterator is returned from FilterClaimedDelegationReward and is used to iterate over the raw logs and unpacked data for ClaimedDelegationReward events raised by the Contract contract.
type ContractClaimedDelegationRewardIterator struct {
	Event *ContractClaimedDelegationReward // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractClaimedDelegationRewardIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractClaimedDelegationReward)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractClaimedDelegationReward)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractClaimedDelegationRewardIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractClaimedDelegationRewardIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractClaimedDelegationReward represents a ClaimedDelegationReward event raised by the Contract contract.
type ContractClaimedDelegationReward struct {
	From       common.Address
	StakerID   *big.Int
	Reward     *big.Int
	FromEpoch  *big.Int
	UntilEpoch *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterClaimedDelegationReward is a free log retrieval operation binding the contract event 0x2676e1697cf4731b93ddb4ef54e0e5a98c06cccbbbb2202848a3c6286595e6ce.
//
// Solidity: event ClaimedDelegationReward(address indexed from, uint256 indexed stakerID, uint256 reward, uint256 fromEpoch, uint256 untilEpoch)
func (_Contract *ContractFilterer) FilterClaimedDelegationReward(opts *bind.FilterOpts, from []common.Address, stakerID []*big.Int) (*ContractClaimedDelegationRewardIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "ClaimedDelegationReward", fromRule, stakerIDRule)
	if err != nil {
		return nil, err
	}
	return &ContractClaimedDelegationRewardIterator{contract: _Contract.contract, event: "ClaimedDelegationReward", logs: logs, sub: sub}, nil
}

// WatchClaimedDelegationReward is a free log subscription operation binding the contract event 0x2676e1697cf4731b93ddb4ef54e0e5a98c06cccbbbb2202848a3c6286595e6ce.
//
// Solidity: event ClaimedDelegationReward(address indexed from, uint256 indexed stakerID, uint256 reward, uint256 fromEpoch, uint256 untilEpoch)
func (_Contract *ContractFilterer) WatchClaimedDelegationReward(opts *bind.WatchOpts, sink chan<- *ContractClaimedDelegationReward, from []common.Address, stakerID []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "ClaimedDelegationReward", fromRule, stakerIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractClaimedDelegationReward)
				if err := _Contract.contract.UnpackLog(event, "ClaimedDelegationReward", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseClaimedDelegationReward is a log parse operation binding the contract event 0x2676e1697cf4731b93ddb4ef54e0e5a98c06cccbbbb2202848a3c6286595e6ce.
//
// Solidity: event ClaimedDelegationReward(address indexed from, uint256 indexed stakerID, uint256 reward, uint256 fromEpoch, uint256 untilEpoch)
func (_Contract *ContractFilterer) ParseClaimedDelegationReward(log types.Log) (*ContractClaimedDelegationReward, error) {
	event := new(ContractClaimedDelegationReward)
	if err := _Contract.contract.UnpackLog(event, "ClaimedDelegationReward", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractClaimedValidatorRewardIterator is returned from FilterClaimedValidatorReward and is used to iterate over the raw logs and unpacked data for ClaimedValidatorReward events raised by the Contract contract.
type ContractClaimedValidatorRewardIterator struct {
	Event *ContractClaimedValidatorReward // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractClaimedValidatorRewardIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractClaimedValidatorReward)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractClaimedValidatorReward)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractClaimedValidatorRewardIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractClaimedValidatorRewardIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractClaimedValidatorReward represents a ClaimedValidatorReward event raised by the Contract contract.
type ContractClaimedValidatorReward struct {
	StakerID   *big.Int
	Reward     *big.Int
	FromEpoch  *big.Int
	UntilEpoch *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterClaimedValidatorReward is a free log retrieval operation binding the contract event 0x2ea54c2b22a07549d19fb5eb8e4e48ebe1c653117215e94d5468c5612750d35c.
//
// Solidity: event ClaimedValidatorReward(uint256 indexed stakerID, uint256 reward, uint256 fromEpoch, uint256 untilEpoch)
func (_Contract *ContractFilterer) FilterClaimedValidatorReward(opts *bind.FilterOpts, stakerID []*big.Int) (*ContractClaimedValidatorRewardIterator, error) {

	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "ClaimedValidatorReward", stakerIDRule)
	if err != nil {
		return nil, err
	}
	return &ContractClaimedValidatorRewardIterator{contract: _Contract.contract, event: "ClaimedValidatorReward", logs: logs, sub: sub}, nil
}

// WatchClaimedValidatorReward is a free log subscription operation binding the contract event 0x2ea54c2b22a07549d19fb5eb8e4e48ebe1c653117215e94d5468c5612750d35c.
//
// Solidity: event ClaimedValidatorReward(uint256 indexed stakerID, uint256 reward, uint256 fromEpoch, uint256 untilEpoch)
func (_Contract *ContractFilterer) WatchClaimedValidatorReward(opts *bind.WatchOpts, sink chan<- *ContractClaimedValidatorReward, stakerID []*big.Int) (event.Subscription, error) {

	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "ClaimedValidatorReward", stakerIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractClaimedValidatorReward)
				if err := _Contract.contract.UnpackLog(event, "ClaimedValidatorReward", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseClaimedValidatorReward is a log parse operation binding the contract event 0x2ea54c2b22a07549d19fb5eb8e4e48ebe1c653117215e94d5468c5612750d35c.
//
// Solidity: event ClaimedValidatorReward(uint256 indexed stakerID, uint256 reward, uint256 fromEpoch, uint256 untilEpoch)
func (_Contract *ContractFilterer) ParseClaimedValidatorReward(log types.Log) (*ContractClaimedValidatorReward, error) {
	event := new(ContractClaimedValidatorReward)
	if err := _Contract.contract.UnpackLog(event, "ClaimedValidatorReward", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractCreatedDelegationIterator is returned from FilterCreatedDelegation and is used to iterate over the raw logs and unpacked data for CreatedDelegation events raised by the Contract contract.
type ContractCreatedDelegationIterator struct {
	Event *ContractCreatedDelegation // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractCreatedDelegationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractCreatedDelegation)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractCreatedDelegation)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractCreatedDelegationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractCreatedDelegationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractCreatedDelegation represents a CreatedDelegation event raised by the Contract contract.
type ContractCreatedDelegation struct {
	Delegator  common.Address
	ToStakerID *big.Int
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCreatedDelegation is a free log retrieval operation binding the contract event 0xfd8c857fb9acd6f4ad59b8621a2a77825168b7b4b76de9586d08e00d4ed462be.
//
// Solidity: event CreatedDelegation(address indexed delegator, uint256 indexed toStakerID, uint256 amount)
func (_Contract *ContractFilterer) FilterCreatedDelegation(opts *bind.FilterOpts, delegator []common.Address, toStakerID []*big.Int) (*ContractCreatedDelegationIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var toStakerIDRule []interface{}
	for _, toStakerIDItem := range toStakerID {
		toStakerIDRule = append(toStakerIDRule, toStakerIDItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "CreatedDelegation", delegatorRule, toStakerIDRule)
	if err != nil {
		return nil, err
	}
	return &ContractCreatedDelegationIterator{contract: _Contract.contract, event: "CreatedDelegation", logs: logs, sub: sub}, nil
}

// WatchCreatedDelegation is a free log subscription operation binding the contract event 0xfd8c857fb9acd6f4ad59b8621a2a77825168b7b4b76de9586d08e00d4ed462be.
//
// Solidity: event CreatedDelegation(address indexed delegator, uint256 indexed toStakerID, uint256 amount)
func (_Contract *ContractFilterer) WatchCreatedDelegation(opts *bind.WatchOpts, sink chan<- *ContractCreatedDelegation, delegator []common.Address, toStakerID []*big.Int) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var toStakerIDRule []interface{}
	for _, toStakerIDItem := range toStakerID {
		toStakerIDRule = append(toStakerIDRule, toStakerIDItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "CreatedDelegation", delegatorRule, toStakerIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractCreatedDelegation)
				if err := _Contract.contract.UnpackLog(event, "CreatedDelegation", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCreatedDelegation is a log parse operation binding the contract event 0xfd8c857fb9acd6f4ad59b8621a2a77825168b7b4b76de9586d08e00d4ed462be.
//
// Solidity: event CreatedDelegation(address indexed delegator, uint256 indexed toStakerID, uint256 amount)
func (_Contract *ContractFilterer) ParseCreatedDelegation(log types.Log) (*ContractCreatedDelegation, error) {
	event := new(ContractCreatedDelegation)
	if err := _Contract.contract.UnpackLog(event, "CreatedDelegation", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractCreatedStakeIterator is returned from FilterCreatedStake and is used to iterate over the raw logs and unpacked data for CreatedStake events raised by the Contract contract.
type ContractCreatedStakeIterator struct {
	Event *ContractCreatedStake // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractCreatedStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractCreatedStake)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractCreatedStake)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractCreatedStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractCreatedStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractCreatedStake represents a CreatedStake event raised by the Contract contract.
type ContractCreatedStake struct {
	StakerID      *big.Int
	DagSfcAddress common.Address
	Amount        *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterCreatedStake is a free log retrieval operation binding the contract event 0x0697dfe5062b9db8108e4b31254f47a912ae6bbb78837667b2e923a6f5160d39.
//
// Solidity: event CreatedStake(uint256 indexed stakerID, address indexed dagSfcAddress, uint256 amount)
func (_Contract *ContractFilterer) FilterCreatedStake(opts *bind.FilterOpts, stakerID []*big.Int, dagSfcAddress []common.Address) (*ContractCreatedStakeIterator, error) {

	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}
	var dagSfcAddressRule []interface{}
	for _, dagSfcAddressItem := range dagSfcAddress {
		dagSfcAddressRule = append(dagSfcAddressRule, dagSfcAddressItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "CreatedStake", stakerIDRule, dagSfcAddressRule)
	if err != nil {
		return nil, err
	}
	return &ContractCreatedStakeIterator{contract: _Contract.contract, event: "CreatedStake", logs: logs, sub: sub}, nil
}

// WatchCreatedStake is a free log subscription operation binding the contract event 0x0697dfe5062b9db8108e4b31254f47a912ae6bbb78837667b2e923a6f5160d39.
//
// Solidity: event CreatedStake(uint256 indexed stakerID, address indexed dagSfcAddress, uint256 amount)
func (_Contract *ContractFilterer) WatchCreatedStake(opts *bind.WatchOpts, sink chan<- *ContractCreatedStake, stakerID []*big.Int, dagSfcAddress []common.Address) (event.Subscription, error) {

	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}
	var dagSfcAddressRule []interface{}
	for _, dagSfcAddressItem := range dagSfcAddress {
		dagSfcAddressRule = append(dagSfcAddressRule, dagSfcAddressItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "CreatedStake", stakerIDRule, dagSfcAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractCreatedStake)
				if err := _Contract.contract.UnpackLog(event, "CreatedStake", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCreatedStake is a log parse operation binding the contract event 0x0697dfe5062b9db8108e4b31254f47a912ae6bbb78837667b2e923a6f5160d39.
//
// Solidity: event CreatedStake(uint256 indexed stakerID, address indexed dagSfcAddress, uint256 amount)
func (_Contract *ContractFilterer) ParseCreatedStake(log types.Log) (*ContractCreatedStake, error) {
	event := new(ContractCreatedStake)
	if err := _Contract.contract.UnpackLog(event, "CreatedStake", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractCreatedWithdrawRequestIterator is returned from FilterCreatedWithdrawRequest and is used to iterate over the raw logs and unpacked data for CreatedWithdrawRequest events raised by the Contract contract.
type ContractCreatedWithdrawRequestIterator struct {
	Event *ContractCreatedWithdrawRequest // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractCreatedWithdrawRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractCreatedWithdrawRequest)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractCreatedWithdrawRequest)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractCreatedWithdrawRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractCreatedWithdrawRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractCreatedWithdrawRequest represents a CreatedWithdrawRequest event raised by the Contract contract.
type ContractCreatedWithdrawRequest struct {
	Auth       common.Address
	Receiver   common.Address
	StakerID   *big.Int
	WrID       *big.Int
	Delegation bool
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCreatedWithdrawRequest is a free log retrieval operation binding the contract event 0xde2d2a87af2fa2de55bde86f04143144eb632fa6be266dc224341a371fb8916d.
//
// Solidity: event CreatedWithdrawRequest(address indexed auth, address indexed receiver, uint256 indexed stakerID, uint256 wrID, bool delegation, uint256 amount)
func (_Contract *ContractFilterer) FilterCreatedWithdrawRequest(opts *bind.FilterOpts, auth []common.Address, receiver []common.Address, stakerID []*big.Int) (*ContractCreatedWithdrawRequestIterator, error) {

	var authRule []interface{}
	for _, authItem := range auth {
		authRule = append(authRule, authItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "CreatedWithdrawRequest", authRule, receiverRule, stakerIDRule)
	if err != nil {
		return nil, err
	}
	return &ContractCreatedWithdrawRequestIterator{contract: _Contract.contract, event: "CreatedWithdrawRequest", logs: logs, sub: sub}, nil
}

// WatchCreatedWithdrawRequest is a free log subscription operation binding the contract event 0xde2d2a87af2fa2de55bde86f04143144eb632fa6be266dc224341a371fb8916d.
//
// Solidity: event CreatedWithdrawRequest(address indexed auth, address indexed receiver, uint256 indexed stakerID, uint256 wrID, bool delegation, uint256 amount)
func (_Contract *ContractFilterer) WatchCreatedWithdrawRequest(opts *bind.WatchOpts, sink chan<- *ContractCreatedWithdrawRequest, auth []common.Address, receiver []common.Address, stakerID []*big.Int) (event.Subscription, error) {

	var authRule []interface{}
	for _, authItem := range auth {
		authRule = append(authRule, authItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "CreatedWithdrawRequest", authRule, receiverRule, stakerIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractCreatedWithdrawRequest)
				if err := _Contract.contract.UnpackLog(event, "CreatedWithdrawRequest", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCreatedWithdrawRequest is a log parse operation binding the contract event 0xde2d2a87af2fa2de55bde86f04143144eb632fa6be266dc224341a371fb8916d.
//
// Solidity: event CreatedWithdrawRequest(address indexed auth, address indexed receiver, uint256 indexed stakerID, uint256 wrID, bool delegation, uint256 amount)
func (_Contract *ContractFilterer) ParseCreatedWithdrawRequest(log types.Log) (*ContractCreatedWithdrawRequest, error) {
	event := new(ContractCreatedWithdrawRequest)
	if err := _Contract.contract.UnpackLog(event, "CreatedWithdrawRequest", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractDeactivatedDelegationIterator is returned from FilterDeactivatedDelegation and is used to iterate over the raw logs and unpacked data for DeactivatedDelegation events raised by the Contract contract.
type ContractDeactivatedDelegationIterator struct {
	Event *ContractDeactivatedDelegation // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractDeactivatedDelegationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractDeactivatedDelegation)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractDeactivatedDelegation)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractDeactivatedDelegationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractDeactivatedDelegationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractDeactivatedDelegation represents a DeactivatedDelegation event raised by the Contract contract.
type ContractDeactivatedDelegation struct {
	Delegator common.Address
	StakerID  *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDeactivatedDelegation is a free log retrieval operation binding the contract event 0x912c4125a208704a342cbdc4726795d26556b0170b7fc95bc706d5cb1f506469.
//
// Solidity: event DeactivatedDelegation(address indexed delegator, uint256 indexed stakerID)
func (_Contract *ContractFilterer) FilterDeactivatedDelegation(opts *bind.FilterOpts, delegator []common.Address, stakerID []*big.Int) (*ContractDeactivatedDelegationIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "DeactivatedDelegation", delegatorRule, stakerIDRule)
	if err != nil {
		return nil, err
	}
	return &ContractDeactivatedDelegationIterator{contract: _Contract.contract, event: "DeactivatedDelegation", logs: logs, sub: sub}, nil
}

// WatchDeactivatedDelegation is a free log subscription operation binding the contract event 0x912c4125a208704a342cbdc4726795d26556b0170b7fc95bc706d5cb1f506469.
//
// Solidity: event DeactivatedDelegation(address indexed delegator, uint256 indexed stakerID)
func (_Contract *ContractFilterer) WatchDeactivatedDelegation(opts *bind.WatchOpts, sink chan<- *ContractDeactivatedDelegation, delegator []common.Address, stakerID []*big.Int) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "DeactivatedDelegation", delegatorRule, stakerIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractDeactivatedDelegation)
				if err := _Contract.contract.UnpackLog(event, "DeactivatedDelegation", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDeactivatedDelegation is a log parse operation binding the contract event 0x912c4125a208704a342cbdc4726795d26556b0170b7fc95bc706d5cb1f506469.
//
// Solidity: event DeactivatedDelegation(address indexed delegator, uint256 indexed stakerID)
func (_Contract *ContractFilterer) ParseDeactivatedDelegation(log types.Log) (*ContractDeactivatedDelegation, error) {
	event := new(ContractDeactivatedDelegation)
	if err := _Contract.contract.UnpackLog(event, "DeactivatedDelegation", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractDeactivatedStakeIterator is returned from FilterDeactivatedStake and is used to iterate over the raw logs and unpacked data for DeactivatedStake events raised by the Contract contract.
type ContractDeactivatedStakeIterator struct {
	Event *ContractDeactivatedStake // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractDeactivatedStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractDeactivatedStake)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractDeactivatedStake)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractDeactivatedStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractDeactivatedStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractDeactivatedStake represents a DeactivatedStake event raised by the Contract contract.
type ContractDeactivatedStake struct {
	StakerID *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDeactivatedStake is a free log retrieval operation binding the contract event 0xf7c308d0d978cce3aec157d1b34e355db4636b4e71ce91b4f5ec9e7a4f5cdc60.
//
// Solidity: event DeactivatedStake(uint256 indexed stakerID)
func (_Contract *ContractFilterer) FilterDeactivatedStake(opts *bind.FilterOpts, stakerID []*big.Int) (*ContractDeactivatedStakeIterator, error) {

	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "DeactivatedStake", stakerIDRule)
	if err != nil {
		return nil, err
	}
	return &ContractDeactivatedStakeIterator{contract: _Contract.contract, event: "DeactivatedStake", logs: logs, sub: sub}, nil
}

// WatchDeactivatedStake is a free log subscription operation binding the contract event 0xf7c308d0d978cce3aec157d1b34e355db4636b4e71ce91b4f5ec9e7a4f5cdc60.
//
// Solidity: event DeactivatedStake(uint256 indexed stakerID)
func (_Contract *ContractFilterer) WatchDeactivatedStake(opts *bind.WatchOpts, sink chan<- *ContractDeactivatedStake, stakerID []*big.Int) (event.Subscription, error) {

	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "DeactivatedStake", stakerIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractDeactivatedStake)
				if err := _Contract.contract.UnpackLog(event, "DeactivatedStake", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDeactivatedStake is a log parse operation binding the contract event 0xf7c308d0d978cce3aec157d1b34e355db4636b4e71ce91b4f5ec9e7a4f5cdc60.
//
// Solidity: event DeactivatedStake(uint256 indexed stakerID)
func (_Contract *ContractFilterer) ParseDeactivatedStake(log types.Log) (*ContractDeactivatedStake, error) {
	event := new(ContractDeactivatedStake)
	if err := _Contract.contract.UnpackLog(event, "DeactivatedStake", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractIncreasedDelegationIterator is returned from FilterIncreasedDelegation and is used to iterate over the raw logs and unpacked data for IncreasedDelegation events raised by the Contract contract.
type ContractIncreasedDelegationIterator struct {
	Event *ContractIncreasedDelegation // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractIncreasedDelegationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIncreasedDelegation)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractIncreasedDelegation)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractIncreasedDelegationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIncreasedDelegationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIncreasedDelegation represents a IncreasedDelegation event raised by the Contract contract.
type ContractIncreasedDelegation struct {
	Delegator common.Address
	StakerID  *big.Int
	NewAmount *big.Int
	Diff      *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterIncreasedDelegation is a free log retrieval operation binding the contract event 0x4ca781bfe171e588a2661d5a7f2f5f59df879c53489063552fbad2145b707fc1.
//
// Solidity: event IncreasedDelegation(address indexed delegator, uint256 indexed stakerID, uint256 newAmount, uint256 diff)
func (_Contract *ContractFilterer) FilterIncreasedDelegation(opts *bind.FilterOpts, delegator []common.Address, stakerID []*big.Int) (*ContractIncreasedDelegationIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "IncreasedDelegation", delegatorRule, stakerIDRule)
	if err != nil {
		return nil, err
	}
	return &ContractIncreasedDelegationIterator{contract: _Contract.contract, event: "IncreasedDelegation", logs: logs, sub: sub}, nil
}

// WatchIncreasedDelegation is a free log subscription operation binding the contract event 0x4ca781bfe171e588a2661d5a7f2f5f59df879c53489063552fbad2145b707fc1.
//
// Solidity: event IncreasedDelegation(address indexed delegator, uint256 indexed stakerID, uint256 newAmount, uint256 diff)
func (_Contract *ContractFilterer) WatchIncreasedDelegation(opts *bind.WatchOpts, sink chan<- *ContractIncreasedDelegation, delegator []common.Address, stakerID []*big.Int) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "IncreasedDelegation", delegatorRule, stakerIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIncreasedDelegation)
				if err := _Contract.contract.UnpackLog(event, "IncreasedDelegation", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseIncreasedDelegation is a log parse operation binding the contract event 0x4ca781bfe171e588a2661d5a7f2f5f59df879c53489063552fbad2145b707fc1.
//
// Solidity: event IncreasedDelegation(address indexed delegator, uint256 indexed stakerID, uint256 newAmount, uint256 diff)
func (_Contract *ContractFilterer) ParseIncreasedDelegation(log types.Log) (*ContractIncreasedDelegation, error) {
	event := new(ContractIncreasedDelegation)
	if err := _Contract.contract.UnpackLog(event, "IncreasedDelegation", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractIncreasedStakeIterator is returned from FilterIncreasedStake and is used to iterate over the raw logs and unpacked data for IncreasedStake events raised by the Contract contract.
type ContractIncreasedStakeIterator struct {
	Event *ContractIncreasedStake // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractIncreasedStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractIncreasedStake)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractIncreasedStake)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractIncreasedStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractIncreasedStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractIncreasedStake represents a IncreasedStake event raised by the Contract contract.
type ContractIncreasedStake struct {
	StakerID  *big.Int
	NewAmount *big.Int
	Diff      *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterIncreasedStake is a free log retrieval operation binding the contract event 0xa1d93e9a2a16bf4c2d0cdc6f47fe0fa054c741c96b3dac1297c79eaca31714e9.
//
// Solidity: event IncreasedStake(uint256 indexed stakerID, uint256 newAmount, uint256 diff)
func (_Contract *ContractFilterer) FilterIncreasedStake(opts *bind.FilterOpts, stakerID []*big.Int) (*ContractIncreasedStakeIterator, error) {

	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "IncreasedStake", stakerIDRule)
	if err != nil {
		return nil, err
	}
	return &ContractIncreasedStakeIterator{contract: _Contract.contract, event: "IncreasedStake", logs: logs, sub: sub}, nil
}

// WatchIncreasedStake is a free log subscription operation binding the contract event 0xa1d93e9a2a16bf4c2d0cdc6f47fe0fa054c741c96b3dac1297c79eaca31714e9.
//
// Solidity: event IncreasedStake(uint256 indexed stakerID, uint256 newAmount, uint256 diff)
func (_Contract *ContractFilterer) WatchIncreasedStake(opts *bind.WatchOpts, sink chan<- *ContractIncreasedStake, stakerID []*big.Int) (event.Subscription, error) {

	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "IncreasedStake", stakerIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractIncreasedStake)
				if err := _Contract.contract.UnpackLog(event, "IncreasedStake", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseIncreasedStake is a log parse operation binding the contract event 0xa1d93e9a2a16bf4c2d0cdc6f47fe0fa054c741c96b3dac1297c79eaca31714e9.
//
// Solidity: event IncreasedStake(uint256 indexed stakerID, uint256 newAmount, uint256 diff)
func (_Contract *ContractFilterer) ParseIncreasedStake(log types.Log) (*ContractIncreasedStake, error) {
	event := new(ContractIncreasedStake)
	if err := _Contract.contract.UnpackLog(event, "IncreasedStake", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractLockingDelegationIterator is returned from FilterLockingDelegation and is used to iterate over the raw logs and unpacked data for LockingDelegation events raised by the Contract contract.
type ContractLockingDelegationIterator struct {
	Event *ContractLockingDelegation // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractLockingDelegationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractLockingDelegation)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractLockingDelegation)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractLockingDelegationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractLockingDelegationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractLockingDelegation represents a LockingDelegation event raised by the Contract contract.
type ContractLockingDelegation struct {
	Delegator common.Address
	StakerID  *big.Int
	FromEpoch *big.Int
	EndTime   *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterLockingDelegation is a free log retrieval operation binding the contract event 0x823f252f996e1f519fd0215db7eb4d5a688d78587bf03bfb03d77bfca939806d.
//
// Solidity: event LockingDelegation(address indexed delegator, uint256 indexed stakerID, uint256 fromEpoch, uint256 endTime)
func (_Contract *ContractFilterer) FilterLockingDelegation(opts *bind.FilterOpts, delegator []common.Address, stakerID []*big.Int) (*ContractLockingDelegationIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "LockingDelegation", delegatorRule, stakerIDRule)
	if err != nil {
		return nil, err
	}
	return &ContractLockingDelegationIterator{contract: _Contract.contract, event: "LockingDelegation", logs: logs, sub: sub}, nil
}

// WatchLockingDelegation is a free log subscription operation binding the contract event 0x823f252f996e1f519fd0215db7eb4d5a688d78587bf03bfb03d77bfca939806d.
//
// Solidity: event LockingDelegation(address indexed delegator, uint256 indexed stakerID, uint256 fromEpoch, uint256 endTime)
func (_Contract *ContractFilterer) WatchLockingDelegation(opts *bind.WatchOpts, sink chan<- *ContractLockingDelegation, delegator []common.Address, stakerID []*big.Int) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "LockingDelegation", delegatorRule, stakerIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractLockingDelegation)
				if err := _Contract.contract.UnpackLog(event, "LockingDelegation", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseLockingDelegation is a log parse operation binding the contract event 0x823f252f996e1f519fd0215db7eb4d5a688d78587bf03bfb03d77bfca939806d.
//
// Solidity: event LockingDelegation(address indexed delegator, uint256 indexed stakerID, uint256 fromEpoch, uint256 endTime)
func (_Contract *ContractFilterer) ParseLockingDelegation(log types.Log) (*ContractLockingDelegation, error) {
	event := new(ContractLockingDelegation)
	if err := _Contract.contract.UnpackLog(event, "LockingDelegation", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractLockingStakeIterator is returned from FilterLockingStake and is used to iterate over the raw logs and unpacked data for LockingStake events raised by the Contract contract.
type ContractLockingStakeIterator struct {
	Event *ContractLockingStake // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractLockingStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractLockingStake)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractLockingStake)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractLockingStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractLockingStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractLockingStake represents a LockingStake event raised by the Contract contract.
type ContractLockingStake struct {
	StakerID  *big.Int
	FromEpoch *big.Int
	EndTime   *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterLockingStake is a free log retrieval operation binding the contract event 0x71f8e76b11dde805fa567e857c4beba340500f4ca9da003520a25014f542162b.
//
// Solidity: event LockingStake(uint256 indexed stakerID, uint256 fromEpoch, uint256 endTime)
func (_Contract *ContractFilterer) FilterLockingStake(opts *bind.FilterOpts, stakerID []*big.Int) (*ContractLockingStakeIterator, error) {

	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "LockingStake", stakerIDRule)
	if err != nil {
		return nil, err
	}
	return &ContractLockingStakeIterator{contract: _Contract.contract, event: "LockingStake", logs: logs, sub: sub}, nil
}

// WatchLockingStake is a free log subscription operation binding the contract event 0x71f8e76b11dde805fa567e857c4beba340500f4ca9da003520a25014f542162b.
//
// Solidity: event LockingStake(uint256 indexed stakerID, uint256 fromEpoch, uint256 endTime)
func (_Contract *ContractFilterer) WatchLockingStake(opts *bind.WatchOpts, sink chan<- *ContractLockingStake, stakerID []*big.Int) (event.Subscription, error) {

	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "LockingStake", stakerIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractLockingStake)
				if err := _Contract.contract.UnpackLog(event, "LockingStake", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseLockingStake is a log parse operation binding the contract event 0x71f8e76b11dde805fa567e857c4beba340500f4ca9da003520a25014f542162b.
//
// Solidity: event LockingStake(uint256 indexed stakerID, uint256 fromEpoch, uint256 endTime)
func (_Contract *ContractFilterer) ParseLockingStake(log types.Log) (*ContractLockingStake, error) {
	event := new(ContractLockingStake)
	if err := _Contract.contract.UnpackLog(event, "LockingStake", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Contract contract.
type ContractOwnershipTransferredIterator struct {
	Event *ContractOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractOwnershipTransferred represents a OwnershipTransferred event raised by the Contract contract.
type ContractOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Contract *ContractFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ContractOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ContractOwnershipTransferredIterator{contract: _Contract.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Contract *ContractFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ContractOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractOwnershipTransferred)
				if err := _Contract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Contract *ContractFilterer) ParseOwnershipTransferred(log types.Log) (*ContractOwnershipTransferred, error) {
	event := new(ContractOwnershipTransferred)
	if err := _Contract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractPartialWithdrawnByRequestIterator is returned from FilterPartialWithdrawnByRequest and is used to iterate over the raw logs and unpacked data for PartialWithdrawnByRequest events raised by the Contract contract.
type ContractPartialWithdrawnByRequestIterator struct {
	Event *ContractPartialWithdrawnByRequest // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractPartialWithdrawnByRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPartialWithdrawnByRequest)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractPartialWithdrawnByRequest)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractPartialWithdrawnByRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPartialWithdrawnByRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPartialWithdrawnByRequest represents a PartialWithdrawnByRequest event raised by the Contract contract.
type ContractPartialWithdrawnByRequest struct {
	Auth       common.Address
	Receiver   common.Address
	StakerID   *big.Int
	WrID       *big.Int
	Delegation bool
	Penalty    *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterPartialWithdrawnByRequest is a free log retrieval operation binding the contract event 0xd5304dabc5bd47105b6921889d1b528c4b2223250248a916afd129b1c0512ddd.
//
// Solidity: event PartialWithdrawnByRequest(address indexed auth, address indexed receiver, uint256 indexed stakerID, uint256 wrID, bool delegation, uint256 penalty)
func (_Contract *ContractFilterer) FilterPartialWithdrawnByRequest(opts *bind.FilterOpts, auth []common.Address, receiver []common.Address, stakerID []*big.Int) (*ContractPartialWithdrawnByRequestIterator, error) {

	var authRule []interface{}
	for _, authItem := range auth {
		authRule = append(authRule, authItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "PartialWithdrawnByRequest", authRule, receiverRule, stakerIDRule)
	if err != nil {
		return nil, err
	}
	return &ContractPartialWithdrawnByRequestIterator{contract: _Contract.contract, event: "PartialWithdrawnByRequest", logs: logs, sub: sub}, nil
}

// WatchPartialWithdrawnByRequest is a free log subscription operation binding the contract event 0xd5304dabc5bd47105b6921889d1b528c4b2223250248a916afd129b1c0512ddd.
//
// Solidity: event PartialWithdrawnByRequest(address indexed auth, address indexed receiver, uint256 indexed stakerID, uint256 wrID, bool delegation, uint256 penalty)
func (_Contract *ContractFilterer) WatchPartialWithdrawnByRequest(opts *bind.WatchOpts, sink chan<- *ContractPartialWithdrawnByRequest, auth []common.Address, receiver []common.Address, stakerID []*big.Int) (event.Subscription, error) {

	var authRule []interface{}
	for _, authItem := range auth {
		authRule = append(authRule, authItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}
	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "PartialWithdrawnByRequest", authRule, receiverRule, stakerIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPartialWithdrawnByRequest)
				if err := _Contract.contract.UnpackLog(event, "PartialWithdrawnByRequest", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePartialWithdrawnByRequest is a log parse operation binding the contract event 0xd5304dabc5bd47105b6921889d1b528c4b2223250248a916afd129b1c0512ddd.
//
// Solidity: event PartialWithdrawnByRequest(address indexed auth, address indexed receiver, uint256 indexed stakerID, uint256 wrID, bool delegation, uint256 penalty)
func (_Contract *ContractFilterer) ParsePartialWithdrawnByRequest(log types.Log) (*ContractPartialWithdrawnByRequest, error) {
	event := new(ContractPartialWithdrawnByRequest)
	if err := _Contract.contract.UnpackLog(event, "PartialWithdrawnByRequest", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractPreparedToWithdrawDelegationIterator is returned from FilterPreparedToWithdrawDelegation and is used to iterate over the raw logs and unpacked data for PreparedToWithdrawDelegation events raised by the Contract contract.
type ContractPreparedToWithdrawDelegationIterator struct {
	Event *ContractPreparedToWithdrawDelegation // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractPreparedToWithdrawDelegationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPreparedToWithdrawDelegation)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractPreparedToWithdrawDelegation)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractPreparedToWithdrawDelegationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPreparedToWithdrawDelegationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPreparedToWithdrawDelegation represents a PreparedToWithdrawDelegation event raised by the Contract contract.
type ContractPreparedToWithdrawDelegation struct {
	Delegator common.Address
	StakerID  *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterPreparedToWithdrawDelegation is a free log retrieval operation binding the contract event 0x5b1eea49e405ef6d509836aac841959c30bb0673b1fd70859bfc6ae5e4ee3df2.
//
// Solidity: event PreparedToWithdrawDelegation(address indexed delegator, uint256 indexed stakerID)
func (_Contract *ContractFilterer) FilterPreparedToWithdrawDelegation(opts *bind.FilterOpts, delegator []common.Address, stakerID []*big.Int) (*ContractPreparedToWithdrawDelegationIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "PreparedToWithdrawDelegation", delegatorRule, stakerIDRule)
	if err != nil {
		return nil, err
	}
	return &ContractPreparedToWithdrawDelegationIterator{contract: _Contract.contract, event: "PreparedToWithdrawDelegation", logs: logs, sub: sub}, nil
}

// WatchPreparedToWithdrawDelegation is a free log subscription operation binding the contract event 0x5b1eea49e405ef6d509836aac841959c30bb0673b1fd70859bfc6ae5e4ee3df2.
//
// Solidity: event PreparedToWithdrawDelegation(address indexed delegator, uint256 indexed stakerID)
func (_Contract *ContractFilterer) WatchPreparedToWithdrawDelegation(opts *bind.WatchOpts, sink chan<- *ContractPreparedToWithdrawDelegation, delegator []common.Address, stakerID []*big.Int) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "PreparedToWithdrawDelegation", delegatorRule, stakerIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPreparedToWithdrawDelegation)
				if err := _Contract.contract.UnpackLog(event, "PreparedToWithdrawDelegation", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePreparedToWithdrawDelegation is a log parse operation binding the contract event 0x5b1eea49e405ef6d509836aac841959c30bb0673b1fd70859bfc6ae5e4ee3df2.
//
// Solidity: event PreparedToWithdrawDelegation(address indexed delegator, uint256 indexed stakerID)
func (_Contract *ContractFilterer) ParsePreparedToWithdrawDelegation(log types.Log) (*ContractPreparedToWithdrawDelegation, error) {
	event := new(ContractPreparedToWithdrawDelegation)
	if err := _Contract.contract.UnpackLog(event, "PreparedToWithdrawDelegation", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractPreparedToWithdrawStakeIterator is returned from FilterPreparedToWithdrawStake and is used to iterate over the raw logs and unpacked data for PreparedToWithdrawStake events raised by the Contract contract.
type ContractPreparedToWithdrawStakeIterator struct {
	Event *ContractPreparedToWithdrawStake // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractPreparedToWithdrawStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractPreparedToWithdrawStake)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractPreparedToWithdrawStake)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractPreparedToWithdrawStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractPreparedToWithdrawStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractPreparedToWithdrawStake represents a PreparedToWithdrawStake event raised by the Contract contract.
type ContractPreparedToWithdrawStake struct {
	StakerID *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterPreparedToWithdrawStake is a free log retrieval operation binding the contract event 0x84244546a9da4942f506db48ff90ebc240c73bb399e3e47d58843c6bb60e7185.
//
// Solidity: event PreparedToWithdrawStake(uint256 indexed stakerID)
func (_Contract *ContractFilterer) FilterPreparedToWithdrawStake(opts *bind.FilterOpts, stakerID []*big.Int) (*ContractPreparedToWithdrawStakeIterator, error) {

	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "PreparedToWithdrawStake", stakerIDRule)
	if err != nil {
		return nil, err
	}
	return &ContractPreparedToWithdrawStakeIterator{contract: _Contract.contract, event: "PreparedToWithdrawStake", logs: logs, sub: sub}, nil
}

// WatchPreparedToWithdrawStake is a free log subscription operation binding the contract event 0x84244546a9da4942f506db48ff90ebc240c73bb399e3e47d58843c6bb60e7185.
//
// Solidity: event PreparedToWithdrawStake(uint256 indexed stakerID)
func (_Contract *ContractFilterer) WatchPreparedToWithdrawStake(opts *bind.WatchOpts, sink chan<- *ContractPreparedToWithdrawStake, stakerID []*big.Int) (event.Subscription, error) {

	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "PreparedToWithdrawStake", stakerIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractPreparedToWithdrawStake)
				if err := _Contract.contract.UnpackLog(event, "PreparedToWithdrawStake", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePreparedToWithdrawStake is a log parse operation binding the contract event 0x84244546a9da4942f506db48ff90ebc240c73bb399e3e47d58843c6bb60e7185.
//
// Solidity: event PreparedToWithdrawStake(uint256 indexed stakerID)
func (_Contract *ContractFilterer) ParsePreparedToWithdrawStake(log types.Log) (*ContractPreparedToWithdrawStake, error) {
	event := new(ContractPreparedToWithdrawStake)
	if err := _Contract.contract.UnpackLog(event, "PreparedToWithdrawStake", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractUnstashedRewardsIterator is returned from FilterUnstashedRewards and is used to iterate over the raw logs and unpacked data for UnstashedRewards events raised by the Contract contract.
type ContractUnstashedRewardsIterator struct {
	Event *ContractUnstashedRewards // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractUnstashedRewardsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractUnstashedRewards)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractUnstashedRewards)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractUnstashedRewardsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractUnstashedRewardsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractUnstashedRewards represents a UnstashedRewards event raised by the Contract contract.
type ContractUnstashedRewards struct {
	Auth     common.Address
	Receiver common.Address
	Rewards  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterUnstashedRewards is a free log retrieval operation binding the contract event 0x80b36a0e929d7e7925087e54acfeecf4c6043e451b9d71ac5e908b66f9e5d126.
//
// Solidity: event UnstashedRewards(address indexed auth, address indexed receiver, uint256 rewards)
func (_Contract *ContractFilterer) FilterUnstashedRewards(opts *bind.FilterOpts, auth []common.Address, receiver []common.Address) (*ContractUnstashedRewardsIterator, error) {

	var authRule []interface{}
	for _, authItem := range auth {
		authRule = append(authRule, authItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "UnstashedRewards", authRule, receiverRule)
	if err != nil {
		return nil, err
	}
	return &ContractUnstashedRewardsIterator{contract: _Contract.contract, event: "UnstashedRewards", logs: logs, sub: sub}, nil
}

// WatchUnstashedRewards is a free log subscription operation binding the contract event 0x80b36a0e929d7e7925087e54acfeecf4c6043e451b9d71ac5e908b66f9e5d126.
//
// Solidity: event UnstashedRewards(address indexed auth, address indexed receiver, uint256 rewards)
func (_Contract *ContractFilterer) WatchUnstashedRewards(opts *bind.WatchOpts, sink chan<- *ContractUnstashedRewards, auth []common.Address, receiver []common.Address) (event.Subscription, error) {

	var authRule []interface{}
	for _, authItem := range auth {
		authRule = append(authRule, authItem)
	}
	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "UnstashedRewards", authRule, receiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractUnstashedRewards)
				if err := _Contract.contract.UnpackLog(event, "UnstashedRewards", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnstashedRewards is a log parse operation binding the contract event 0x80b36a0e929d7e7925087e54acfeecf4c6043e451b9d71ac5e908b66f9e5d126.
//
// Solidity: event UnstashedRewards(address indexed auth, address indexed receiver, uint256 rewards)
func (_Contract *ContractFilterer) ParseUnstashedRewards(log types.Log) (*ContractUnstashedRewards, error) {
	event := new(ContractUnstashedRewards)
	if err := _Contract.contract.UnpackLog(event, "UnstashedRewards", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractUpdatedBaseRewardPerSecIterator is returned from FilterUpdatedBaseRewardPerSec and is used to iterate over the raw logs and unpacked data for UpdatedBaseRewardPerSec events raised by the Contract contract.
type ContractUpdatedBaseRewardPerSecIterator struct {
	Event *ContractUpdatedBaseRewardPerSec // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractUpdatedBaseRewardPerSecIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractUpdatedBaseRewardPerSec)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractUpdatedBaseRewardPerSec)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractUpdatedBaseRewardPerSecIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractUpdatedBaseRewardPerSecIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractUpdatedBaseRewardPerSec represents a UpdatedBaseRewardPerSec event raised by the Contract contract.
type ContractUpdatedBaseRewardPerSec struct {
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterUpdatedBaseRewardPerSec is a free log retrieval operation binding the contract event 0x8cd9dae1bbea2bc8a5e80ffce2c224727a25925130a03ae100619a8861ae2396.
//
// Solidity: event UpdatedBaseRewardPerSec(uint256 value)
func (_Contract *ContractFilterer) FilterUpdatedBaseRewardPerSec(opts *bind.FilterOpts) (*ContractUpdatedBaseRewardPerSecIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "UpdatedBaseRewardPerSec")
	if err != nil {
		return nil, err
	}
	return &ContractUpdatedBaseRewardPerSecIterator{contract: _Contract.contract, event: "UpdatedBaseRewardPerSec", logs: logs, sub: sub}, nil
}

// WatchUpdatedBaseRewardPerSec is a free log subscription operation binding the contract event 0x8cd9dae1bbea2bc8a5e80ffce2c224727a25925130a03ae100619a8861ae2396.
//
// Solidity: event UpdatedBaseRewardPerSec(uint256 value)
func (_Contract *ContractFilterer) WatchUpdatedBaseRewardPerSec(opts *bind.WatchOpts, sink chan<- *ContractUpdatedBaseRewardPerSec) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "UpdatedBaseRewardPerSec")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractUpdatedBaseRewardPerSec)
				if err := _Contract.contract.UnpackLog(event, "UpdatedBaseRewardPerSec", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpdatedBaseRewardPerSec is a log parse operation binding the contract event 0x8cd9dae1bbea2bc8a5e80ffce2c224727a25925130a03ae100619a8861ae2396.
//
// Solidity: event UpdatedBaseRewardPerSec(uint256 value)
func (_Contract *ContractFilterer) ParseUpdatedBaseRewardPerSec(log types.Log) (*ContractUpdatedBaseRewardPerSec, error) {
	event := new(ContractUpdatedBaseRewardPerSec)
	if err := _Contract.contract.UnpackLog(event, "UpdatedBaseRewardPerSec", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractUpdatedDelegationIterator is returned from FilterUpdatedDelegation and is used to iterate over the raw logs and unpacked data for UpdatedDelegation events raised by the Contract contract.
type ContractUpdatedDelegationIterator struct {
	Event *ContractUpdatedDelegation // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractUpdatedDelegationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractUpdatedDelegation)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractUpdatedDelegation)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractUpdatedDelegationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractUpdatedDelegationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractUpdatedDelegation represents a UpdatedDelegation event raised by the Contract contract.
type ContractUpdatedDelegation struct {
	Delegator   common.Address
	OldStakerID *big.Int
	NewStakerID *big.Int
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterUpdatedDelegation is a free log retrieval operation binding the contract event 0x19b46b9014e4dc8ca74f505b8921797c6a8a489860217d15b3c7d741637dfcff.
//
// Solidity: event UpdatedDelegation(address indexed delegator, uint256 indexed oldStakerID, uint256 indexed newStakerID, uint256 amount)
func (_Contract *ContractFilterer) FilterUpdatedDelegation(opts *bind.FilterOpts, delegator []common.Address, oldStakerID []*big.Int, newStakerID []*big.Int) (*ContractUpdatedDelegationIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var oldStakerIDRule []interface{}
	for _, oldStakerIDItem := range oldStakerID {
		oldStakerIDRule = append(oldStakerIDRule, oldStakerIDItem)
	}
	var newStakerIDRule []interface{}
	for _, newStakerIDItem := range newStakerID {
		newStakerIDRule = append(newStakerIDRule, newStakerIDItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "UpdatedDelegation", delegatorRule, oldStakerIDRule, newStakerIDRule)
	if err != nil {
		return nil, err
	}
	return &ContractUpdatedDelegationIterator{contract: _Contract.contract, event: "UpdatedDelegation", logs: logs, sub: sub}, nil
}

// WatchUpdatedDelegation is a free log subscription operation binding the contract event 0x19b46b9014e4dc8ca74f505b8921797c6a8a489860217d15b3c7d741637dfcff.
//
// Solidity: event UpdatedDelegation(address indexed delegator, uint256 indexed oldStakerID, uint256 indexed newStakerID, uint256 amount)
func (_Contract *ContractFilterer) WatchUpdatedDelegation(opts *bind.WatchOpts, sink chan<- *ContractUpdatedDelegation, delegator []common.Address, oldStakerID []*big.Int, newStakerID []*big.Int) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var oldStakerIDRule []interface{}
	for _, oldStakerIDItem := range oldStakerID {
		oldStakerIDRule = append(oldStakerIDRule, oldStakerIDItem)
	}
	var newStakerIDRule []interface{}
	for _, newStakerIDItem := range newStakerID {
		newStakerIDRule = append(newStakerIDRule, newStakerIDItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "UpdatedDelegation", delegatorRule, oldStakerIDRule, newStakerIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractUpdatedDelegation)
				if err := _Contract.contract.UnpackLog(event, "UpdatedDelegation", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpdatedDelegation is a log parse operation binding the contract event 0x19b46b9014e4dc8ca74f505b8921797c6a8a489860217d15b3c7d741637dfcff.
//
// Solidity: event UpdatedDelegation(address indexed delegator, uint256 indexed oldStakerID, uint256 indexed newStakerID, uint256 amount)
func (_Contract *ContractFilterer) ParseUpdatedDelegation(log types.Log) (*ContractUpdatedDelegation, error) {
	event := new(ContractUpdatedDelegation)
	if err := _Contract.contract.UnpackLog(event, "UpdatedDelegation", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractUpdatedGasPowerAllocationRateIterator is returned from FilterUpdatedGasPowerAllocationRate and is used to iterate over the raw logs and unpacked data for UpdatedGasPowerAllocationRate events raised by the Contract contract.
type ContractUpdatedGasPowerAllocationRateIterator struct {
	Event *ContractUpdatedGasPowerAllocationRate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractUpdatedGasPowerAllocationRateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractUpdatedGasPowerAllocationRate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractUpdatedGasPowerAllocationRate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractUpdatedGasPowerAllocationRateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractUpdatedGasPowerAllocationRateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractUpdatedGasPowerAllocationRate represents a UpdatedGasPowerAllocationRate event raised by the Contract contract.
type ContractUpdatedGasPowerAllocationRate struct {
	Short *big.Int
	Long  *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterUpdatedGasPowerAllocationRate is a free log retrieval operation binding the contract event 0x95ae5488127de4bc98492f4487556e7af9f37eb4b6d5e94f6d849e03ff76cc7c.
//
// Solidity: event UpdatedGasPowerAllocationRate(uint256 short, uint256 long)
func (_Contract *ContractFilterer) FilterUpdatedGasPowerAllocationRate(opts *bind.FilterOpts) (*ContractUpdatedGasPowerAllocationRateIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "UpdatedGasPowerAllocationRate")
	if err != nil {
		return nil, err
	}
	return &ContractUpdatedGasPowerAllocationRateIterator{contract: _Contract.contract, event: "UpdatedGasPowerAllocationRate", logs: logs, sub: sub}, nil
}

// WatchUpdatedGasPowerAllocationRate is a free log subscription operation binding the contract event 0x95ae5488127de4bc98492f4487556e7af9f37eb4b6d5e94f6d849e03ff76cc7c.
//
// Solidity: event UpdatedGasPowerAllocationRate(uint256 short, uint256 long)
func (_Contract *ContractFilterer) WatchUpdatedGasPowerAllocationRate(opts *bind.WatchOpts, sink chan<- *ContractUpdatedGasPowerAllocationRate) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "UpdatedGasPowerAllocationRate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractUpdatedGasPowerAllocationRate)
				if err := _Contract.contract.UnpackLog(event, "UpdatedGasPowerAllocationRate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpdatedGasPowerAllocationRate is a log parse operation binding the contract event 0x95ae5488127de4bc98492f4487556e7af9f37eb4b6d5e94f6d849e03ff76cc7c.
//
// Solidity: event UpdatedGasPowerAllocationRate(uint256 short, uint256 long)
func (_Contract *ContractFilterer) ParseUpdatedGasPowerAllocationRate(log types.Log) (*ContractUpdatedGasPowerAllocationRate, error) {
	event := new(ContractUpdatedGasPowerAllocationRate)
	if err := _Contract.contract.UnpackLog(event, "UpdatedGasPowerAllocationRate", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractUpdatedStakeIterator is returned from FilterUpdatedStake and is used to iterate over the raw logs and unpacked data for UpdatedStake events raised by the Contract contract.
type ContractUpdatedStakeIterator struct {
	Event *ContractUpdatedStake // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractUpdatedStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractUpdatedStake)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractUpdatedStake)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractUpdatedStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractUpdatedStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractUpdatedStake represents a UpdatedStake event raised by the Contract contract.
type ContractUpdatedStake struct {
	StakerID    *big.Int
	Amount      *big.Int
	DelegatedMe *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterUpdatedStake is a free log retrieval operation binding the contract event 0x509404fa75ce234a1273cf9f7918bcf54e0ef19f2772e4f71b6526606a723b7c.
//
// Solidity: event UpdatedStake(uint256 indexed stakerID, uint256 amount, uint256 delegatedMe)
func (_Contract *ContractFilterer) FilterUpdatedStake(opts *bind.FilterOpts, stakerID []*big.Int) (*ContractUpdatedStakeIterator, error) {

	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "UpdatedStake", stakerIDRule)
	if err != nil {
		return nil, err
	}
	return &ContractUpdatedStakeIterator{contract: _Contract.contract, event: "UpdatedStake", logs: logs, sub: sub}, nil
}

// WatchUpdatedStake is a free log subscription operation binding the contract event 0x509404fa75ce234a1273cf9f7918bcf54e0ef19f2772e4f71b6526606a723b7c.
//
// Solidity: event UpdatedStake(uint256 indexed stakerID, uint256 amount, uint256 delegatedMe)
func (_Contract *ContractFilterer) WatchUpdatedStake(opts *bind.WatchOpts, sink chan<- *ContractUpdatedStake, stakerID []*big.Int) (event.Subscription, error) {

	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "UpdatedStake", stakerIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractUpdatedStake)
				if err := _Contract.contract.UnpackLog(event, "UpdatedStake", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpdatedStake is a log parse operation binding the contract event 0x509404fa75ce234a1273cf9f7918bcf54e0ef19f2772e4f71b6526606a723b7c.
//
// Solidity: event UpdatedStake(uint256 indexed stakerID, uint256 amount, uint256 delegatedMe)
func (_Contract *ContractFilterer) ParseUpdatedStake(log types.Log) (*ContractUpdatedStake, error) {
	event := new(ContractUpdatedStake)
	if err := _Contract.contract.UnpackLog(event, "UpdatedStake", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractUpdatedStakerMetadataIterator is returned from FilterUpdatedStakerMetadata and is used to iterate over the raw logs and unpacked data for UpdatedStakerMetadata events raised by the Contract contract.
type ContractUpdatedStakerMetadataIterator struct {
	Event *ContractUpdatedStakerMetadata // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractUpdatedStakerMetadataIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractUpdatedStakerMetadata)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractUpdatedStakerMetadata)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractUpdatedStakerMetadataIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractUpdatedStakerMetadataIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractUpdatedStakerMetadata represents a UpdatedStakerMetadata event raised by the Contract contract.
type ContractUpdatedStakerMetadata struct {
	StakerID *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterUpdatedStakerMetadata is a free log retrieval operation binding the contract event 0xb7a99a0df6a9e15c2689e6a55811ef76cdb514c67d4a0e37fcb125ada0e3cd83.
//
// Solidity: event UpdatedStakerMetadata(uint256 indexed stakerID)
func (_Contract *ContractFilterer) FilterUpdatedStakerMetadata(opts *bind.FilterOpts, stakerID []*big.Int) (*ContractUpdatedStakerMetadataIterator, error) {

	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "UpdatedStakerMetadata", stakerIDRule)
	if err != nil {
		return nil, err
	}
	return &ContractUpdatedStakerMetadataIterator{contract: _Contract.contract, event: "UpdatedStakerMetadata", logs: logs, sub: sub}, nil
}

// WatchUpdatedStakerMetadata is a free log subscription operation binding the contract event 0xb7a99a0df6a9e15c2689e6a55811ef76cdb514c67d4a0e37fcb125ada0e3cd83.
//
// Solidity: event UpdatedStakerMetadata(uint256 indexed stakerID)
func (_Contract *ContractFilterer) WatchUpdatedStakerMetadata(opts *bind.WatchOpts, sink chan<- *ContractUpdatedStakerMetadata, stakerID []*big.Int) (event.Subscription, error) {

	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "UpdatedStakerMetadata", stakerIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractUpdatedStakerMetadata)
				if err := _Contract.contract.UnpackLog(event, "UpdatedStakerMetadata", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpdatedStakerMetadata is a log parse operation binding the contract event 0xb7a99a0df6a9e15c2689e6a55811ef76cdb514c67d4a0e37fcb125ada0e3cd83.
//
// Solidity: event UpdatedStakerMetadata(uint256 indexed stakerID)
func (_Contract *ContractFilterer) ParseUpdatedStakerMetadata(log types.Log) (*ContractUpdatedStakerMetadata, error) {
	event := new(ContractUpdatedStakerMetadata)
	if err := _Contract.contract.UnpackLog(event, "UpdatedStakerMetadata", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractUpdatedStakerSfcAddressIterator is returned from FilterUpdatedStakerSfcAddress and is used to iterate over the raw logs and unpacked data for UpdatedStakerSfcAddress events raised by the Contract contract.
type ContractUpdatedStakerSfcAddressIterator struct {
	Event *ContractUpdatedStakerSfcAddress // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractUpdatedStakerSfcAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractUpdatedStakerSfcAddress)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractUpdatedStakerSfcAddress)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractUpdatedStakerSfcAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractUpdatedStakerSfcAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractUpdatedStakerSfcAddress represents a UpdatedStakerSfcAddress event raised by the Contract contract.
type ContractUpdatedStakerSfcAddress struct {
	StakerID      *big.Int
	OldSfcAddress common.Address
	NewSfcAddress common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterUpdatedStakerSfcAddress is a free log retrieval operation binding the contract event 0x7cc102ee500cbca85691c9642080562e8f012b04d27f5b7f389453672b206946.
//
// Solidity: event UpdatedStakerSfcAddress(uint256 indexed stakerID, address indexed oldSfcAddress, address indexed newSfcAddress)
func (_Contract *ContractFilterer) FilterUpdatedStakerSfcAddress(opts *bind.FilterOpts, stakerID []*big.Int, oldSfcAddress []common.Address, newSfcAddress []common.Address) (*ContractUpdatedStakerSfcAddressIterator, error) {

	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}
	var oldSfcAddressRule []interface{}
	for _, oldSfcAddressItem := range oldSfcAddress {
		oldSfcAddressRule = append(oldSfcAddressRule, oldSfcAddressItem)
	}
	var newSfcAddressRule []interface{}
	for _, newSfcAddressItem := range newSfcAddress {
		newSfcAddressRule = append(newSfcAddressRule, newSfcAddressItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "UpdatedStakerSfcAddress", stakerIDRule, oldSfcAddressRule, newSfcAddressRule)
	if err != nil {
		return nil, err
	}
	return &ContractUpdatedStakerSfcAddressIterator{contract: _Contract.contract, event: "UpdatedStakerSfcAddress", logs: logs, sub: sub}, nil
}

// WatchUpdatedStakerSfcAddress is a free log subscription operation binding the contract event 0x7cc102ee500cbca85691c9642080562e8f012b04d27f5b7f389453672b206946.
//
// Solidity: event UpdatedStakerSfcAddress(uint256 indexed stakerID, address indexed oldSfcAddress, address indexed newSfcAddress)
func (_Contract *ContractFilterer) WatchUpdatedStakerSfcAddress(opts *bind.WatchOpts, sink chan<- *ContractUpdatedStakerSfcAddress, stakerID []*big.Int, oldSfcAddress []common.Address, newSfcAddress []common.Address) (event.Subscription, error) {

	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}
	var oldSfcAddressRule []interface{}
	for _, oldSfcAddressItem := range oldSfcAddress {
		oldSfcAddressRule = append(oldSfcAddressRule, oldSfcAddressItem)
	}
	var newSfcAddressRule []interface{}
	for _, newSfcAddressItem := range newSfcAddress {
		newSfcAddressRule = append(newSfcAddressRule, newSfcAddressItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "UpdatedStakerSfcAddress", stakerIDRule, oldSfcAddressRule, newSfcAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractUpdatedStakerSfcAddress)
				if err := _Contract.contract.UnpackLog(event, "UpdatedStakerSfcAddress", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpdatedStakerSfcAddress is a log parse operation binding the contract event 0x7cc102ee500cbca85691c9642080562e8f012b04d27f5b7f389453672b206946.
//
// Solidity: event UpdatedStakerSfcAddress(uint256 indexed stakerID, address indexed oldSfcAddress, address indexed newSfcAddress)
func (_Contract *ContractFilterer) ParseUpdatedStakerSfcAddress(log types.Log) (*ContractUpdatedStakerSfcAddress, error) {
	event := new(ContractUpdatedStakerSfcAddress)
	if err := _Contract.contract.UnpackLog(event, "UpdatedStakerSfcAddress", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractWithdrawnDelegationIterator is returned from FilterWithdrawnDelegation and is used to iterate over the raw logs and unpacked data for WithdrawnDelegation events raised by the Contract contract.
type ContractWithdrawnDelegationIterator struct {
	Event *ContractWithdrawnDelegation // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractWithdrawnDelegationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractWithdrawnDelegation)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractWithdrawnDelegation)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractWithdrawnDelegationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractWithdrawnDelegationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractWithdrawnDelegation represents a WithdrawnDelegation event raised by the Contract contract.
type ContractWithdrawnDelegation struct {
	Delegator common.Address
	StakerID  *big.Int
	Penalty   *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdrawnDelegation is a free log retrieval operation binding the contract event 0x87e86b3710b72c10173ca52c6a9f9cf2df27e77ed177741a8b4feb12bb7a606f.
//
// Solidity: event WithdrawnDelegation(address indexed delegator, uint256 indexed stakerID, uint256 penalty)
func (_Contract *ContractFilterer) FilterWithdrawnDelegation(opts *bind.FilterOpts, delegator []common.Address, stakerID []*big.Int) (*ContractWithdrawnDelegationIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "WithdrawnDelegation", delegatorRule, stakerIDRule)
	if err != nil {
		return nil, err
	}
	return &ContractWithdrawnDelegationIterator{contract: _Contract.contract, event: "WithdrawnDelegation", logs: logs, sub: sub}, nil
}

// WatchWithdrawnDelegation is a free log subscription operation binding the contract event 0x87e86b3710b72c10173ca52c6a9f9cf2df27e77ed177741a8b4feb12bb7a606f.
//
// Solidity: event WithdrawnDelegation(address indexed delegator, uint256 indexed stakerID, uint256 penalty)
func (_Contract *ContractFilterer) WatchWithdrawnDelegation(opts *bind.WatchOpts, sink chan<- *ContractWithdrawnDelegation, delegator []common.Address, stakerID []*big.Int) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "WithdrawnDelegation", delegatorRule, stakerIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractWithdrawnDelegation)
				if err := _Contract.contract.UnpackLog(event, "WithdrawnDelegation", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdrawnDelegation is a log parse operation binding the contract event 0x87e86b3710b72c10173ca52c6a9f9cf2df27e77ed177741a8b4feb12bb7a606f.
//
// Solidity: event WithdrawnDelegation(address indexed delegator, uint256 indexed stakerID, uint256 penalty)
func (_Contract *ContractFilterer) ParseWithdrawnDelegation(log types.Log) (*ContractWithdrawnDelegation, error) {
	event := new(ContractWithdrawnDelegation)
	if err := _Contract.contract.UnpackLog(event, "WithdrawnDelegation", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ContractWithdrawnStakeIterator is returned from FilterWithdrawnStake and is used to iterate over the raw logs and unpacked data for WithdrawnStake events raised by the Contract contract.
type ContractWithdrawnStakeIterator struct {
	Event *ContractWithdrawnStake // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractWithdrawnStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractWithdrawnStake)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractWithdrawnStake)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractWithdrawnStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractWithdrawnStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractWithdrawnStake represents a WithdrawnStake event raised by the Contract contract.
type ContractWithdrawnStake struct {
	StakerID *big.Int
	Penalty  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterWithdrawnStake is a free log retrieval operation binding the contract event 0x8c6548258f8f12a9d4b593fa89a223417ed901d4ee9712ba09beb4d56f5262b6.
//
// Solidity: event WithdrawnStake(uint256 indexed stakerID, uint256 penalty)
func (_Contract *ContractFilterer) FilterWithdrawnStake(opts *bind.FilterOpts, stakerID []*big.Int) (*ContractWithdrawnStakeIterator, error) {

	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "WithdrawnStake", stakerIDRule)
	if err != nil {
		return nil, err
	}
	return &ContractWithdrawnStakeIterator{contract: _Contract.contract, event: "WithdrawnStake", logs: logs, sub: sub}, nil
}

// WatchWithdrawnStake is a free log subscription operation binding the contract event 0x8c6548258f8f12a9d4b593fa89a223417ed901d4ee9712ba09beb4d56f5262b6.
//
// Solidity: event WithdrawnStake(uint256 indexed stakerID, uint256 penalty)
func (_Contract *ContractFilterer) WatchWithdrawnStake(opts *bind.WatchOpts, sink chan<- *ContractWithdrawnStake, stakerID []*big.Int) (event.Subscription, error) {

	var stakerIDRule []interface{}
	for _, stakerIDItem := range stakerID {
		stakerIDRule = append(stakerIDRule, stakerIDItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "WithdrawnStake", stakerIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractWithdrawnStake)
				if err := _Contract.contract.UnpackLog(event, "WithdrawnStake", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdrawnStake is a log parse operation binding the contract event 0x8c6548258f8f12a9d4b593fa89a223417ed901d4ee9712ba09beb4d56f5262b6.
//
// Solidity: event WithdrawnStake(uint256 indexed stakerID, uint256 penalty)
func (_Contract *ContractFilterer) ParseWithdrawnStake(log types.Log) (*ContractWithdrawnStake, error) {
	event := new(ContractWithdrawnStake)
	if err := _Contract.contract.UnpackLog(event, "WithdrawnStake", log); err != nil {
		return nil, err
	}
	return event, nil
}
