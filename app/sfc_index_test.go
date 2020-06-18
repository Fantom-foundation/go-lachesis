package app

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/Fantom-foundation/go-lachesis/app/contract"
	"github.com/Fantom-foundation/go-lachesis/utils"
)

func Test(t *testing.T) {
	env := newTestEnv(3)
	defer env.Close()

	t.Run("Transfer", func(t *testing.T) {
		require := require.New(t)

		b0 := utils.ToFtm(1e18)
		require.Equal(b0, env.State().GetBalance(env.Address(0)))
		require.Equal(b0, env.State().GetBalance(env.Address(1)))
		require.Equal(b0, env.State().GetBalance(env.Address(2)))

		env.ApplyBlock(
			env.Tx(0, 1, utils.ToFtm(100)),
		)
		env.ApplyBlock(
			env.Tx(1, 2, utils.ToFtm(100)),
		)
		env.ApplyBlock(
			env.Tx(2, 0, utils.ToFtm(100)),
		)

		gas := big.NewInt(0).Mul(big.NewInt(21000), env.GasPrice)
		b1 := big.NewInt(0).Sub(b0, gas)
		require.Equal(b1, env.State().GetBalance(env.Address(0)))
		require.Equal(b1, env.State().GetBalance(env.Address(1)))
		require.Equal(b1, env.State().GetBalance(env.Address(2)))

	})

	t.Run("SFC deploy", func(t *testing.T) {
		require := require.New(t)

		mainContractBinV2 := hexutil.MustDecode(contract.StoreBin)
		r := env.ApplyBlock(
			env.Contract(0, utils.ToFtm(0), mainContractBinV2),
		)
		require.Equal(r[0].Status, eth.ReceiptStatusSuccessful, "tx failed")

		contract2, err := contract.NewStore(r[0].ContractAddress, env)
		require.NoError(err)

		epoch, err := contract2.StoreCaller.CurrentEpoch(&bind.CallOpts{})
		require.NoError(err)
		t.Logf("Epoch: %d", epoch.Uint64())

	})
}
