package sfcpos

import (
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/utils"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

// Global variables

func CurrentSealedEpoch() common.Hash {
	return utils.U64to256(0)
}

func VStakersLastIdx() common.Hash {
	return utils.U64to256(3)
}

func VStakersNum() common.Hash {
	return utils.U64to256(4)
}

func VStakeTotalAmount() common.Hash {
	return utils.U64to256(5)
}

// VStake

type VStakePos struct {
	object
}

func VStake(vstaker common.Address) VStakePos {
	position := getMapValue(common.Hash{}, vstaker.Hash(), 2)

	return VStakePos{object{base: position.Big()}}
}

func (p *VStakePos) IsCheater() common.Hash {
	return p.Field(0)
}

func (p *VStakePos) StakerIdx() common.Hash {
	return p.Field(1)
}

func (p *VStakePos) CreatedEpoch() common.Hash {
	return p.Field(2)
}

func (p *VStakePos) CreatedTime() common.Hash {
	return p.Field(3)
}

func (p *VStakePos) StakeAmount() common.Hash {
	return p.Field(6)
}

// EpochSnapshot

type EpochSnapshotPos struct {
	object
}

func EpochSnapshot(epoch idx.Epoch) EpochSnapshotPos {
	position := getMapValue(common.Hash{}, utils.U64to256(uint64(epoch)), 1)

	return EpochSnapshotPos{object{base: position.Big()}}
}

func (p *EpochSnapshotPos) EndTime() common.Hash {
	return p.Field(1)
}

func (p *EpochSnapshotPos) Duration() common.Hash {
	return p.Field(2)
}

func (p *EpochSnapshotPos) EpochFee() common.Hash {
	return p.Field(3)
}

func (p *EpochSnapshotPos) TotalValidatingPower() common.Hash {
	return p.Field(4)
}

// ValidatorMerit

type ValidatorMeritPos struct {
	object
}

func (p *EpochSnapshotPos) ValidatorMerit(validator common.Address) ValidatorMeritPos {
	base := p.Field(0)

	position := getMapValue(base, validator.Hash(), 0)

	return ValidatorMeritPos{object{base: position.Big()}}
}

func (p *ValidatorMeritPos) ValidatingPower() common.Hash {
	return p.Field(0)
}

func (p *ValidatorMeritPos) StakeAmount() common.Hash {
	return p.Field(1)
}

func (p *ValidatorMeritPos) DelegatedMe() common.Hash {
	return p.Field(2)
}

func (p *ValidatorMeritPos) StakerIdx() common.Hash {
	return p.Field(3)
}

// Util

func getMapValue(base common.Hash, key common.Hash, mapIdx int64) common.Hash {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(key.Bytes())
	start := base.Big()
	hasher.Write(common.BytesToHash(start.Add(start, big.NewInt(mapIdx)).Bytes()).Bytes())

	return common.BytesToHash(hasher.Sum(nil))
}

type object struct {
	base *big.Int
}

func (p *object) Field(offset int64) common.Hash {
	if offset == 0 {
		return common.BytesToHash(p.base.Bytes())
	}

	start := new(big.Int).Set(p.base)

	return common.BytesToHash(start.Add(start, big.NewInt(offset)).Bytes())
}