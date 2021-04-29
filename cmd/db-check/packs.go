package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"

	"github.com/Fantom-foundation/go-lachesis/gossip"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/kvdb"
	"github.com/Fantom-foundation/go-lachesis/kvdb/table"
)

func checkPacks(db kvdb.KeyValueStore) {
	t := table.New(db, []byte("p"))

	it := t.NewIterator(nil, nil)
	defer it.Release()

	for it.Next() {
		buf := it.Key()
		w := it.Value()

		var info gossip.PackInfo
		err := rlp.DecodeBytes(w, &info)
		if err != nil {
			fmt.Printf(">>> %s\n ", string(buf))
			continue
		}

		epoch := idx.BytesToEpoch(buf[0:4])
		pack := idx.BytesToEpoch(buf[4:8])
		fmt.Printf("%d:%d %+v\n", epoch, pack, info)
	}
}
