package main

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/Fantom-foundation/go-lachesis/app"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/kvdb"
	"github.com/Fantom-foundation/go-lachesis/kvdb/flushable"
	"github.com/Fantom-foundation/go-lachesis/poset"
)

func checkReceipts(db kvdb.DbProducer) {
	dbs := flushable.NewSyncedPool(db)

	s := app.NewStore(dbs, app.DefaultStoreConfig())
	defer s.Close()

	//g := gossip.NewStore(dbs, gossip.DefaultStoreConfig())
	//defer g.Close()

	p := poset.NewStore(dbs, poset.DefaultStoreConfig())
	defer p.Close()

	lastBlock := p.GetCheckpoint().LastBlockN
	fmt.Printf("Last block: %d\n", lastBlock)
	for i := idx.Block(0); i < 100000; i++ {
		n := lastBlock - i
		rr := s.GetReceipts(n)
		if rr == nil {
			continue
		}
		fmt.Printf("%d - %d\n", n, memSizeOf(rr))
	}

}

func memSizeOf(v interface{}) int {
	b := new(bytes.Buffer)
	if err := gob.NewEncoder(b).Encode(v); err != nil {
		panic(err)
	}
	return b.Len()
}
