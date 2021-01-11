package main

import (
	"math/big"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"

	"github.com/Fantom-foundation/go-lachesis/inter/pos"
	"github.com/Fantom-foundation/go-lachesis/kvdb/flushable"
	"github.com/Fantom-foundation/go-lachesis/kvdb/leveldb"
	"github.com/Fantom-foundation/go-lachesis/kvdb/nokeyiserr"
	"github.com/Fantom-foundation/go-lachesis/lachesis/genesis"
)

func main() {
	go func() {
		http.ListenAndServe("127.0.0.1:8080", nil)
	}()

	dbs, db := stateDB()

	var stateRoot common.Hash
	accs := genesis.FakeValidators(1000, big.NewInt(10), pos.StakeToBalance(1))
	data := make([]byte, 1024)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-sigs:
			return
		default:
		}

		statedb, err := state.New(stateRoot, db, nil)
		if err != nil {
			panic(err)
		}

		for addr, acc := range accs.Accounts {
			nonce := statedb.GetNonce(addr)
			statedb.AddBalance(addr, acc.Balance)
			statedb.SetNonce(addr, nonce+1)
			rand.Read(data)
			statedb.SetCode(addr, data)
		}

		stateRoot, err := statedb.Commit(true)
		if err != nil {
			panic(err)
		}

		// flush
		err = db.TrieDB().Commit(stateRoot, false, nil)
		if err != nil {
			panic(err)
		}

		if dbs.IsFlushNeeded() {
			dbs.Flush([]byte("flag"))
		}

	}

}

func stateDB() (*flushable.SyncedPool, state.Database) {
	err := os.Mkdir("dbs", 0755)
	if err != nil {
		panic(err)
	}
	disk := leveldb.NewProducer("dbs")
	dbs := flushable.NewSyncedPool(disk)

	table := dbs.GetDb("main")
	// defer table.Close()

	db := state.NewDatabaseWithCache(
		rawdb.NewDatabase(
			nokeyiserr.Wrap(table)),
		16,
		"")

	return dbs, db
}
