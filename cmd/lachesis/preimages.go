package main

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"gopkg.in/urfave/cli.v1"
)

// importPreimages imports preimage data from the specified file.
func importPreimages(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		utils.Fatalf("This command requires an argument.")
	}

	cfg := makeAllConfigs(ctx)

	gdb, cdb := makeStores(cfg.Node.DataDir, &cfg.Lachesis)
	defer gdb.Close()
	defer cdb.Close()

	start := time.Now()

	if err := utils.ImportPreimages(gdb.App().EvmTable(), ctx.Args().First()); err != nil {
		utils.Fatalf("Import error: %v\n", err)
	}
	err := gdb.Commit(nil, true)
	if err != nil {
		utils.Fatalf("DB flushing error: %v\n", err)
	}
	fmt.Printf("Import done in %v\n", time.Since(start))
	return nil
}

// exportPreimages dumps the preimage data to specified json file in streaming way.
func exportPreimages(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		utils.Fatalf("This command requires an argument.")
	}

	cfg := makeAllConfigs(ctx)

	gdb, _ := makeStores(cfg.Node.DataDir, &cfg.Lachesis)
	defer gdb.Close()

	start := time.Now()

	if err := utils.ExportPreimages(gdb.App().EvmTable(), ctx.Args().First()); err != nil {
		utils.Fatalf("Export error: %v\n", err)
	}
	fmt.Printf("Export done in %v\n", time.Since(start))
	return nil
}
