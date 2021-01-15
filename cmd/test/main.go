package main

import (
	"fmt"
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
	"github.com/Fantom-foundation/go-lachesis/kvdb"
	"github.com/Fantom-foundation/go-lachesis/kvdb/flushable"
	"github.com/Fantom-foundation/go-lachesis/kvdb/leveldb"
	"github.com/Fantom-foundation/go-lachesis/kvdb/nokeyiserr"
	"github.com/Fantom-foundation/go-lachesis/lachesis/genesis"
)

func main() {
	go func() {
		http.ListenAndServe("127.0.0.1:8080", nil)
	}()

	var (
		err       error
		dbs       = stateDB()
		table     kvdb.KeyValueStore
		db        state.Database
		stateRoot common.Hash
		accs      = genesis.FakeValidators(100, big.NewInt(10), pos.StakeToBalance(1))
		data      = make([]byte, 512)
	)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
mainloop:
	for i := 0; true; i++ {
		select {
		case <-sigs:
			break mainloop
		default:
		}

		// flush to disk
		if i%10e1 == 0 && table != nil {
			err = db.TrieDB().Commit(stateRoot, false, nil)
			if err != nil {
				panic(err)
			}

			if dbs.IsFlushNeeded() {
				dbs.Flush([]byte(fmt.Sprintf("point%d", i)))
			}
		}

		// free disk space
		if i%10e4 == 0 {
			if table != nil {
				err = table.Close()
				if err != nil {
					panic(err)
				}
				table.Drop()
			}
			table = dbs.GetDb(fmt.Sprintf("main%d", i))

			db = state.NewDatabaseWithCache(
				rawdb.NewDatabase(
					nokeyiserr.Wrap(table)),
				16,
				"")

			stateRoot = common.Hash{}
		}

		// change the state
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
		stateRoot, err = statedb.Commit(true)
		if err != nil {
			panic(err)
		}

	}

	if table != nil {
		table.Close()
	}
}

func stateDB() *flushable.SyncedPool {
	err := os.Mkdir("dbs", 0755)
	if err != nil {
		panic(err)
	}
	disk := leveldb.NewProducer("dbs")
	dbs := flushable.NewSyncedPool(disk)

	return dbs
}
