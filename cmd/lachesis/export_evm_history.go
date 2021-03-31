package main

import (
	"compress/gzip"
	"io"
	"os"
	"strings"

	"github.com/Fantom-foundation/lachesis-base/common/bigendian"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"gopkg.in/urfave/cli.v1"
)

func WriteDB(writer io.Writer, db ethdb.Iteratee) error {
	it := db.NewIterator(nil, nil)
	defer it.Release()
	for it.Next() {
		_, err := writer.Write(bigendian.Uint32ToBytes(uint32(len(it.Key()))))
		if err != nil {
			return err
		}
		_, err = writer.Write(it.Key())
		if err != nil {
			return err
		}
		_, err = writer.Write(bigendian.Uint32ToBytes(uint32(len(it.Value()))))
		if err != nil {
			return err
		}
		_, err = writer.Write(it.Value())
		if err != nil {
			return err
		}
	}
	return nil
}

func exportEvmStorage(ctx *cli.Context) error {
	if len(ctx.Args()) != 1 {
		utils.Fatalf("This command requires one argument.")
	}

	cfg := makeAllConfigs(ctx)

	gdb, _ := makeStores(cfg.Node.DataDir, &cfg.Lachesis)
	defer gdb.Close()

	fn := ctx.Args().First()

	// Open the file handle and potentially wrap with a gzip stream
	fh, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer fh.Close()

	var writer io.Writer = fh
	if strings.HasSuffix(fn, ".gz") {
		writer = gzip.NewWriter(writer)
		defer writer.(*gzip.Writer).Close()
	}

	log.Info("Exporting EVM storage to file", "file", fn)
	defer log.Info("Exported EVM storage to file", "file", fn)

	return WriteDB(writer, gdb.App().EvmTable())
}
