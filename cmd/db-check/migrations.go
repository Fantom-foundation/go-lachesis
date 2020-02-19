package main

import (
	"fmt"

	"github.com/Fantom-foundation/go-lachesis/kvdb"
	"github.com/Fantom-foundation/go-lachesis/kvdb/table"
)

func checkAfterMigration(p kvdb.DbProducer) {
	mainDb := p.OpenDb("gossip-main")
	defer mainDb.Close()

	old1 := table.New(mainDb, []byte("p"))
	printData("old1", old1, []byte("serverPool"))

	old2 := table.New(mainDb, []byte("Z"))
	printData("old2", old2, nil)

	servDb := p.OpenDb("gossip-serv")
	defer servDb.Close()

	dst := table.New(servDb, []byte("Z"))
	printData("dst", dst, nil)
}

func printData(dsc string, src kvdb.KeyValueStore, prefix []byte) error {
	fmt.Println(">>>> " + dsc)

	it := src.NewIteratorWithPrefix(prefix)
	defer it.Release()

	for i := 0; it.Next(); i++ {
		fmt.Printf("%d) %s - %#v\n", i, string(it.Key()), it.Value())
	}

	return nil
}
