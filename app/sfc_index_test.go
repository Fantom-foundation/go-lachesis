package app

import (
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/Fantom-foundation/go-lachesis/inter/idx"

	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/require"

	"github.com/Fantom-foundation/go-lachesis/app/sfc110"
	"github.com/Fantom-foundation/go-lachesis/app/sfc200"
	"github.com/Fantom-foundation/go-lachesis/app/sfcproxy"
	"github.com/Fantom-foundation/go-lachesis/lachesis/genesis/sfc"
	"github.com/Fantom-foundation/go-lachesis/lachesis/genesis/sfc/sfcpos"
	"github.com/Fantom-foundation/go-lachesis/utils"
)

func TestMain(m *testing.M) {

	log.Root().SetHandler(log.LvlFilterHandler(
		log.LvlTrace,
		log.StreamHandler(os.Stderr, log.TerminalFormat(false))))

	os.Exit(m.Run())
}

func TestSFC(t *testing.T) {
	env := newTestEnv()
	defer env.Close()

	sfcProxy, err := sfcproxy.NewContract(sfc.ContractAddress, env)
	require.NoError(t, err)

	var (
		sfc1 *sfc110.Contract
		sfc2 *sfc200.Contract
	)

	_ = true &&

		t.Run("Genesis v1.0.0", func(t *testing.T) {
			// nothing to do
		}) &&

		t.Run("Some transfers 1", func(t *testing.T) {
			cicleTransfers(t, env, 3)
		}) &&

		t.Run("Upgrade to v1.1.0-rc1", func(t *testing.T) {
			require := require.New(t)

			r := env.ApplyBlock(nextEpoch,
				env.Contract(1, utils.ToFtm(0), sfc110.ContractBin),
			)
			newImpl := r[0].ContractAddress

			admin := env.Payer(1)
			tx, err := sfcProxy.ContractTransactor.UpgradeTo(admin, newImpl)
			require.NoError(err)
			env.ApplyBlock(sameEpoch, tx)

			impl, err := sfcProxy.Implementation(env.ReadOnly())
			require.NoError(err)
			require.Equal(newImpl, impl, "SFC-proxy: implementation address")

			sfc1, err = sfc110.NewContract(sfc.ContractAddress, env)
			require.NoError(err)

			epoch, err := sfc1.ContractCaller.CurrentEpoch(env.ReadOnly())
			require.NoError(err)
			require.Equal(epoch.Cmp(big.NewInt(2)), 0, "current epoch")
		}) &&

		t.Run("Upgrade stakers storage", func(t *testing.T) {
			require := require.New(t)

			stakers, err := sfc1.StakersLastID(env.ReadOnly())
			require.NoError(err)
			txs := make([]*eth.Transaction, 0, int(stakers.Int64()))
			for i := stakers.Int64(); i > 0; i-- {
				tx, err := sfc1.UpgradeStakerStorage(env.Payer(int(i)), big.NewInt(i))
				require.NoError(err)
				txs = append(txs, tx)
			}
			env.ApplyBlock(sameEpoch, txs...)

		}) &&

		t.Run("Some transfers 2", func(t *testing.T) {
			cicleTransfers(t, env, 3)
		}) &&

		t.Run("Upgrade to SFC newmodel_exp", func(t *testing.T) {
			require := require.New(t)

			r := env.ApplyBlock(nextEpoch,
				env.Contract(1, utils.ToFtm(0), sfc200.ContractBin),
			)
			newImpl := r[0].ContractAddress

			admin := env.Payer(1)
			tx, err := sfcProxy.ContractTransactor.UpgradeTo(admin, newImpl)
			require.NoError(err)
			env.ApplyBlock(sameEpoch, tx)

			impl, err := sfcProxy.Implementation(env.ReadOnly())
			require.NoError(err)
			require.Equal(newImpl, impl, "SFC-proxy: implementation address")

			sfc2, err = sfc200.NewContract(sfc.ContractAddress, env)
			require.NoError(err)

			epoch, err := sfc2.ContractCaller.CurrentEpoch(env.ReadOnly())
			require.NoError(err)
			require.Equal(epoch.Cmp(big.NewInt(3)), 0, "current epoch")
		}) &&

		t.Run("Some transfers 3", func(t *testing.T) {
			cicleTransfers(t, env, 3)
		}) &&

		t.Run("Enable lockedup functionality", func(t *testing.T) {
			require := require.New(t)

			tx, err := sfc2.StartLockedUp(env.Payer(1), big.NewInt(3))
			require.NoError(err)
			env.ApplyBlock(nextEpoch, tx)

			epoch, err := sfc2.FirstLockedUpEpoch(env.ReadOnly())
			require.NoError(err)
			require.Equal(epoch.Cmp(big.NewInt(3)), 0, "1st locked-up epoch")

			raw := new(big.Int).SetBytes(env.State().GetState(
				sfc.ContractAddress,
				sfcpos.FirstLockedUpEpoch()).Bytes())
			require.Equal(epoch.Cmp(raw), 0, "raw 1st locked-up epoch")

			raw = new(big.Int).SetBytes(env.State().GetState(
				sfc.ContractAddress,
				sfcpos.CurrentSealedEpoch()).Bytes())
			require.Equal(epoch.Cmp(raw), 0, "raw last sealed epoch")
		}) &&

		t.Run("Make staker 4", func(t *testing.T) {
			require := require.New(t)

			newStake := utils.ToFtm(genesisStake / 2)
			minStake, err := sfc2.MinStake(env.ReadOnly())
			require.NoError(err)
			require.Greater(newStake.Cmp(minStake), 0,
				fmt.Sprintf("newStake(%s) < minStake(%s)", newStake, minStake))

			env.ApplyBlock(sameEpoch,
				env.Transfer(1, 4, big.NewInt(0).Add(newStake, utils.ToFtm(10))),
			)
			tx, err := sfc2.CreateStake(env.Payer(4, newStake), nil)
			require.NoError(err)
			env.ApplyBlock(nextEpoch, tx)
			newId, err := sfc2.SfcAddressToStakerID(env.ReadOnly(), env.Address(4))
			require.NoError(err)
			env.AddValidator(idx.StakerID(newId.Uint64()))

			env.ApplyBlock(nextEpoch)

			requireRewards(t, env, sfc2, []int64{2, 2, 2, 1})

			/*
				env.ApplyBlock(sameEpoch,
					env.Transfer(1, 4, utils.ToFtm(1001)),
				)

				staker, err := sfc2.SfcAddressToStakerID(env.ReadOnly(), env.Address(1))
				require.NoError(err)

				tx, err := sfc2.CreateDelegation(env.Payer(4, utils.ToFtm(1000)), staker)
				require.NoError(err)
				env.ApplyBlock(nextEpoch, tx)

				tx, err = sfc2.LockUpStake(env.Payer(1), big.NewInt(15*86400))
				require.NoError(err)
				env.ApplyBlock(sameEpoch, tx)

				tx, err = sfc2.LockUpDelegation(env.Payer(4), big.NewInt(14*86400), staker)
				require.NoError(err)
				env.ApplyBlock(sameEpoch, tx)
			*/
		})
}

func requireRewards(t *testing.T, env *testEnv, sfc2 *sfc200.Contract, stakes []int64) {
	require := require.New(t)

	epoch, err := sfc2.CurrentSealedEpoch(env.ReadOnly())
	require.NoError(err)

	var reward0 *big.Int
	for i, id := range env.validators {
		rewardX, _, _, err := sfc2.CalcValidatorRewards(env.ReadOnly(), big.NewInt(int64(id)), epoch, big.NewInt(1))
		require.NoError(err)
		if i == 0 {
			reward0 = rewardX
			continue
		}

		require.Equal(
			new(big.Int).Mul(reward0, big.NewInt(stakes[i])).Cmp(
				new(big.Int).Mul(rewardX, big.NewInt(stakes[0]))),
			0, "reward#0: %s, reward#%d: %s", reward0, i, rewardX)
	}
}

func printValidators(t *testing.T, env *testEnv, sfc2 *sfc200.Contract) {
	require := require.New(t)

	max, err := sfc2.StakersLastID(env.ReadOnly())
	require.NoError(err)

	for id := big.NewInt(1); id.Cmp(max) <= 0; id.Add(id, big.NewInt(1)) {
		s, err := sfc2.Stakers(env.ReadOnly(), id)
		require.NoError(err)
		t.Logf("%s: %#v", id, s)
	}
}

func cicleTransfers(t *testing.T, env *testEnv, count uint64) {
	require := require.New(t)

	balances := make([]*big.Int, 3)
	for i := range balances {
		balances[i] = env.State().GetBalance(env.Address(i + 1))
	}

	for i := uint64(0); i < count; i++ {
		env.ApplyBlock(sameEpoch,
			env.Transfer(1, 2, utils.ToFtm(100)),
		)
		env.ApplyBlock(sameEpoch,
			env.Transfer(2, 3, utils.ToFtm(100)),
		)
		env.ApplyBlock(sameEpoch,
			env.Transfer(3, 1, utils.ToFtm(100)),
		)
	}

	gas := big.NewInt(0).Mul(big.NewInt(int64(count*gasLimit)), env.GasPrice)
	for i := range balances {
		require.Equal(
			big.NewInt(0).Sub(balances[i], gas),
			env.State().GetBalance(env.Address(i+1)),
			fmt.Sprintf("account%d", i),
		)
	}
}
