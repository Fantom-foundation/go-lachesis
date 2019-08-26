package lachesis

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/Fantom-foundation/go-lachesis/src/crypto"
	"github.com/Fantom-foundation/go-lachesis/src/cryptoaddr"
	"github.com/Fantom-foundation/go-lachesis/src/inter"
	"github.com/Fantom-foundation/go-lachesis/src/inter/pos"
)

var (
	genesisTestTime = inter.Timestamp(1565000000 * time.Second)
)

type Genesis struct {
	Balances  map[common.Address]pos.Stake
	StateHash common.Hash
	Time      inter.Timestamp
}

// FakeGenesis generates fake genesis with n-nodes.
func FakeGenesis(n int) (Genesis, []common.Address, []*crypto.PrivateKey) {
	balances := make(map[common.Address]pos.Stake, n)
	keys := make([]*crypto.PrivateKey, n)
	ids := make([]common.Address, n)
	for i := 0; i < n; i++ {
		keys[i] = crypto.GenerateFakeKey(i)
		ids[i] = cryptoaddr.AddressOf(keys[i].Public())
		balances[ids[i]] = 1000000000
	}

	return Genesis{
		Balances: balances,
		Time:     genesisTestTime,
	}, ids, keys
}

// MainGenesis returns builtin genesis keys of mainnet.
func MainGenesis() Genesis {
	return Genesis{
		Time:     genesisTestTime,
		Balances: map[common.Address]pos.Stake{
			// TODO: fill with official keys and balances.
		},
	}
}

// TestGenesis returns builtin genesis keys of testnet.
func TestGenesis() Genesis {
	return Genesis{
		Time:     genesisTestTime,
		Balances: map[common.Address]pos.Stake{
			// TODO: fill with official keys and balances.
		},
	}
}
