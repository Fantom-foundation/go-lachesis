// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package evmcore

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/stretchr/testify/require"

	"github.com/Fantom-foundation/go-lachesis/inter/pos"
	"github.com/Fantom-foundation/go-lachesis/kvdb/memorydb"
	"github.com/Fantom-foundation/go-lachesis/kvdb/nokeyiserr"
	"github.com/Fantom-foundation/go-lachesis/kvdb/table"
	"github.com/Fantom-foundation/go-lachesis/lachesis"
	"github.com/Fantom-foundation/go-lachesis/lachesis/genesis"
	"github.com/Fantom-foundation/go-lachesis/logger"
)

func TestApplyGenesis(t *testing.T) {
	require := require.New(t)

	logger.SetTestMode(t)

	db1 := rawdb.NewDatabase(
		nokeyiserr.Wrap(
			table.New(
				memorydb.New(), []byte("evm1_"))))
	db2 := rawdb.NewDatabase(
		nokeyiserr.Wrap(
			table.New(
				memorydb.New(), []byte("evm2_"))))

	// no genesis
	_, err := ApplyGenesis(db1, nil)
	require.Error(err)

	// the same genesis
	accsA := genesis.FakeAccounts(0, 3, big.NewInt(10000000000), pos.StakeToBalance(1))
	netA := lachesis.FakeNetConfig(accsA)
	blockA1, err := ApplyGenesis(db1, &netA)
	require.NoError(err)
	blockA2, err := ApplyGenesis(db2, &netA)
	require.NoError(err)
	require.Equal(blockA1, blockA2)

	// different genesis
	accsB := genesis.FakeAccounts(0, 4, big.NewInt(10000000000), pos.StakeToBalance(1))
	netB := lachesis.FakeNetConfig(accsB)
	blockB, err := ApplyGenesis(db2, &netB)
	require.NotEqual(blockA1, blockB)
	require.NoError(err)
}
