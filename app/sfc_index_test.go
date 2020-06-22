package app

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Fantom-foundation/go-lachesis/app/sfc110"
	"github.com/Fantom-foundation/go-lachesis/app/sfc200"
	"github.com/Fantom-foundation/go-lachesis/app/sfcproxy"
	"github.com/Fantom-foundation/go-lachesis/lachesis/genesis/sfc"
	"github.com/Fantom-foundation/go-lachesis/utils"
)

func Test(t *testing.T) {
	env := newTestEnv(3)
	defer env.Close()

	sfcProxy, err := sfcproxy.NewContract(sfc.ContractAddress, env)
	require.NoError(t, err)

	t.Run("SFC 1.0.0 (genesis)", func(t *testing.T) {
		// nothing to do
	})

	t.Run("Some transfers 1", func(t *testing.T) {
		require := require.New(t)

		balances := make([]*big.Int, 3)
		for i := range balances {
			balances[i] = env.State().GetBalance(env.Address(i))
		}

		const N = 2
		for i := 0; i < N; i++ {
			env.ApplyBlock(sameEpoch,
				env.Tx(0, 1, utils.ToFtm(100)),
			)
			env.ApplyBlock(sameEpoch,
				env.Tx(1, 2, utils.ToFtm(100)),
			)
			env.ApplyBlock(sameEpoch,
				env.Tx(2, 0, utils.ToFtm(100)),
			)
		}

		gas := big.NewInt(0).Mul(big.NewInt(int64(N*gasLimit)), env.GasPrice)
		for i := range balances {
			require.Equal(
				big.NewInt(0).Sub(balances[i], gas),
				env.State().GetBalance(env.Address(i)),
				fmt.Sprintf("account%d", i),
			)
		}
	})

	t.Run("Upgrade to SFC v1.1.0-rc1", func(t *testing.T) {
		require := require.New(t)

		r := env.ApplyBlock(nextEpoch,
			env.Contract(1, utils.ToFtm(0), sfc110.ContractBin),
		)
		newImpl := r[0].ContractAddress

		admin := env.transactor(0)
		tx, err := sfcProxy.ContractTransactor.UpgradeTo(admin, newImpl)
		require.NoError(err)
		env.ApplyBlock(sameEpoch, tx)

		impl, err := sfcProxy.Implementation(env.caller())
		require.NoError(err)
		require.Equal(newImpl, impl, "SFC-proxy: implementation address")

		newSfc, err := sfc110.NewContract(sfc.ContractAddress, env)
		require.NoError(err)

		epoch, err := newSfc.ContractCaller.CurrentEpoch(env.caller())
		require.NoError(err)
		require.Equal(uint64(2), epoch.Uint64(), "current epoch")
	})

	t.Run("Some transfers 2", func(t *testing.T) {
		require := require.New(t)

		balances := make([]*big.Int, 3)
		for i := range balances {
			balances[i] = env.State().GetBalance(env.Address(i))
		}

		const N = 2
		for i := 0; i < N; i++ {
			env.ApplyBlock(sameEpoch,
				env.Tx(0, 1, utils.ToFtm(100)),
			)
			env.ApplyBlock(sameEpoch,
				env.Tx(1, 2, utils.ToFtm(100)),
			)
			env.ApplyBlock(sameEpoch,
				env.Tx(2, 0, utils.ToFtm(100)),
			)
		}

		gas := big.NewInt(0).Mul(big.NewInt(int64(N*gasLimit)), env.GasPrice)
		for i := range balances {
			require.Equal(
				big.NewInt(0).Sub(balances[i], gas),
				env.State().GetBalance(env.Address(i)),
				fmt.Sprintf("account%d", i),
			)
		}
	})

	t.Run("Upgrade to SFC newmodel_exp", func(t *testing.T) {
		require := require.New(t)

		r := env.ApplyBlock(nextEpoch,
			env.Contract(1, utils.ToFtm(0), sfc200.ContractBin),
		)
		newImpl := r[0].ContractAddress

		admin := env.transactor(0)
		tx, err := sfcProxy.ContractTransactor.UpgradeTo(admin, newImpl)
		require.NoError(err)
		env.ApplyBlock(sameEpoch, tx)

		impl, err := sfcProxy.Implementation(env.caller())
		require.NoError(err)
		require.Equal(newImpl, impl, "SFC-proxy: implementation address")

		newSfc, err := sfc200.NewContract(sfc.ContractAddress, env)
		require.NoError(err)

		epoch, err := newSfc.ContractCaller.CurrentEpoch(env.caller())
		require.NoError(err)
		require.Equal(uint64(3), epoch.Uint64(), "current epoch")
	})

	t.Run("Some transfers 3", func(t *testing.T) {
		require := require.New(t)

		balances := make([]*big.Int, 3)
		for i := range balances {
			balances[i] = env.State().GetBalance(env.Address(i))
		}

		const N = 2
		for i := 0; i < N; i++ {
			env.ApplyBlock(sameEpoch,
				env.Tx(0, 1, utils.ToFtm(100)),
			)
			env.ApplyBlock(sameEpoch,
				env.Tx(1, 2, utils.ToFtm(100)),
			)
			env.ApplyBlock(sameEpoch,
				env.Tx(2, 0, utils.ToFtm(100)),
			)
		}

		gas := big.NewInt(0).Mul(big.NewInt(int64(N*gasLimit)), env.GasPrice)
		for i := range balances {
			require.Equal(
				big.NewInt(0).Sub(balances[i], gas),
				env.State().GetBalance(env.Address(i)),
				fmt.Sprintf("account%d", i),
			)
		}
	})

}