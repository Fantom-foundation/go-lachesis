package integration

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"

	"github.com/Fantom-foundation/go-lachesis/app"
	"github.com/Fantom-foundation/go-lachesis/gossip"
	"github.com/Fantom-foundation/go-lachesis/kvdb/flushable"
	"github.com/Fantom-foundation/go-lachesis/poset"
)

// MakeEngine makes consensus engine from config.
// MakeEngine makes consensus engine from config.
func MakeEngine(dataDir string, gossipCfg *gossip.Config, appCfg *app.StoreConfig) (*poset.Poset, *app.Store, *gossip.Store) {
	dbs := flushable.NewSyncedPool(DBProducer(dataDir))

	adb := app.NewStore(dbs, *appCfg)
	gdb := gossip.NewStore(dbs, gossipCfg.StoreConfig)
	cdb := poset.NewStore(dbs, poset.DefaultStoreConfig())

	adb.SetName("app-db")
	gdb.SetName("gossip-db")
	cdb.SetName("poset-db")

	// write genesis

	stateRoot, _, err := adb.ApplyGenesis(&gossipCfg.Net)
	if err != nil {
		utils.Fatalf("Failed to write App genesis state: %v", err)
	}

	atropos, appHash, isNew, err := gdb.ApplyGenesis(&gossipCfg.Net, stateRoot)
	if err != nil {
		utils.Fatalf("Failed to write Gossip genesis state: %v", err)
	}

	err = cdb.ApplyGenesis(&gossipCfg.Net.Genesis, atropos, appHash)
	if err != nil {
		utils.Fatalf("Failed to write Poset genesis state: %v", err)
	}

	err = dbs.Flush(atropos.Bytes())
	if err != nil {
		utils.Fatalf("Failed to flush genesis state: %v", err)
	}

	if isNew {
		log.Info("Applied genesis state", "hash", cdb.GetGenesisHash().String())
	} else {
		log.Info("Genesis state is already written", "hash", cdb.GetGenesisHash().String())
	}

	// create consensus
	engine := poset.New(gossipCfg.Net.Dag, cdb, gdb)

	return engine, adb, gdb
}

// SetAccountKey sets key into accounts manager and unlocks it with pswd.
func SetAccountKey(
	am *accounts.Manager, key *ecdsa.PrivateKey, pswd string,
) (
	acc accounts.Account,
) {
	kss := am.Backends(keystore.KeyStoreType)
	if len(kss) < 1 {
		log.Warn("Keystore is not found")
		return
	}
	ks := kss[0].(*keystore.KeyStore)

	acc = accounts.Account{
		Address: crypto.PubkeyToAddress(key.PublicKey),
	}

	imported, err := ks.ImportECDSA(key, pswd)
	if err == nil {
		acc = imported
	} else if err.Error() != "account already exists" {
		log.Crit("Failed to import key", "err", err)
	}

	err = ks.Unlock(acc, pswd)
	if err != nil {
		log.Crit("failed to unlock key", "err", err)
	}

	return
}
