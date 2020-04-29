package app

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/inter/sfctype"
	"github.com/Fantom-foundation/go-lachesis/lachesis"
)

// ApplyGenesis writes initial state.
func (s *Store) ApplyGenesis(net *lachesis.Config) (stateRoot common.Hash, isNew bool, err error) {
	s.migrate()

	stored := s.GetBlock(0)

	if stored != nil {
		isNew = false
		stateRoot, err = calcGenesis(net)
		if err != nil {
			return
		}

		if stateRoot != stored.Root {
			err = fmt.Errorf("database contains incompatible state hash (have %s, new %s)", stored.Root, stateRoot)
		}

		return
	}

	// if we'here, then it's first time genesis is applied
	isNew = true
	stateRoot, err = s.applyGenesis(net)
	if err != nil {
		return
	}

	block := evmcore.GenesisBlock(net, stateRoot)
	info := blockInfo(&block.EvmHeader)
	s.SetBlock(info)

	// init stats
	s.SetEpochStats(0, &sfctype.EpochStats{
		Start:    net.Genesis.Time,
		End:      net.Genesis.Time,
		TotalFee: new(big.Int),
	})
	s.SetDirtyEpochStats(&sfctype.EpochStats{
		Start:    net.Genesis.Time,
		TotalFee: new(big.Int),
	})
	s.SetCheckpoint(Checkpoint{
		BlockN:     0,
		EpochN:     1,
		EpochBlock: 0,
		EpochStart: net.Genesis.Time,
	})

	// calc total pre-minted supply
	totalSupply := big.NewInt(0)
	for _, account := range net.Genesis.Alloc.Accounts {
		totalSupply.Add(totalSupply, account.Balance)
	}
	s.SetTotalSupply(totalSupply)

	validatorsArr := []sfctype.SfcStakerAndID{}
	for _, validator := range net.Genesis.Alloc.Validators {
		staker := &sfctype.SfcStaker{
			Address:      validator.Address,
			CreatedEpoch: 0,
			CreatedTime:  net.Genesis.Time,
			StakeAmount:  validator.Stake,
			DelegatedMe:  big.NewInt(0),
		}
		s.SetSfcStaker(validator.ID, staker)
		validatorsArr = append(validatorsArr, sfctype.SfcStakerAndID{
			StakerID: validator.ID,
			Staker:   staker,
		})
	}
	s.SetEpochValidators(1, validatorsArr)

	s.FlushState()
	return
}

// calcGenesis calcs hash of genesis state.
func calcGenesis(net *lachesis.Config) (common.Hash, error) {
	s := NewMemStore()
	defer s.Close()

	s.Log.SetHandler(log.DiscardHandler())

	return s.applyGenesis(net)
}

func (s *Store) applyGenesis(net *lachesis.Config) (stateRoot common.Hash, err error) {
	stateRoot, err = evmcore.ApplyGenesis(s.table.Evm, net)
	if err != nil {
		return
	}

	return
}
