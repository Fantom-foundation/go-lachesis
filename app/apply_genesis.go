package app

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/inter/sfctype"
	"github.com/Fantom-foundation/go-lachesis/lachesis"
)

func (s *Store) ApplyGenesis(net *lachesis.Config) (root common.Hash, err error) {
	evmBlock, err := evmcore.ApplyGenesis(s.table.Evm, net)
	if err != nil {
		return
	}

	// calc total pre-minted supply
	totalSupply := big.NewInt(0)
	for _, account := range net.Genesis.Alloc.Accounts {
		totalSupply.Add(totalSupply, account.Balance)
	}
	s.SetTotalSupply(totalSupply)

	for _, validator := range net.Genesis.Alloc.Validators {
		staker := &sfctype.SfcStaker{
			Address:      validator.Address,
			CreatedEpoch: 0,
			CreatedTime:  net.Genesis.Time,
			StakeAmount:  validator.Stake,
			DelegatedMe:  big.NewInt(0),
		}
		s.SetSfcStaker(validator.ID, staker)
		s.SetEpochValidator(1, validator.ID, staker)
	}

	root = evmBlock.Root
	return
}
