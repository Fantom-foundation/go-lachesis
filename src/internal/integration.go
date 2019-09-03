package internal

import (
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/p2p/simulations/adapters"

	"github.com/Fantom-foundation/go-lachesis/src/gossip"
	"github.com/Fantom-foundation/go-lachesis/src/lachesis"
	"github.com/Fantom-foundation/go-lachesis/src/poset"
	"github.com/Fantom-foundation/go-lachesis/src/signer"
	"github.com/Fantom-foundation/go-lachesis/src/utils"
)

func NewIntegration(cfg *adapters.NodeConfig, network lachesis.Config) *gossip.Service {
	makeDb := dbProducer(cfg.DataDir)
	gdb, cdb := makeStorages(makeDb)

	err := cdb.ApplyGenesis(&network.Genesis)
	if err != nil {
		panic(err)
	}

	c := poset.New(cdb, gdb)

	config := gossip.DefaultConfig(network)

	testSigner := signer.NewSignerTestManager(utils.NewTempDir("lachesis-config"))

	svc, err := gossip.NewService(config, new(event.TypeMux), gdb, c, testSigner)
	if err != nil {
		panic(err)
	}

	return svc
}
