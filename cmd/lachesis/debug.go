package main

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"gopkg.in/urfave/cli.v1"

	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/utils/errlock"
)

var (
	staterootsCommand = cli.Command{
		Action:    utils.MigrateFlags(printStateRoots),
		Name:      "stateroots",
		Usage:     "Prints root hashes for states of blockchain blocks",
		ArgsUsage: "[<blockFrom> [<blockTo>]]",
		Category:  "MISCELLANEOUS COMMANDS",
		Description: `
For debug purpose.
Optional first and second arguments control the first and last block to print.`,
	}
)

func printStateRoots(ctx *cli.Context) error {
	cfg := makeAllConfigs(ctx)

	// check errlock file
	errlock.SetDefaultDatadir(cfg.Node.DataDir)
	errlock.Check()

	gdb := makeGossipStore(cfg.Node.DataDir, &cfg.Lachesis)
	defer gdb.Close()

	from := idx.Block(0)
	if len(ctx.Args()) > 0 {
		n, err := strconv.ParseUint(ctx.Args().Get(0), 10, 32)
		if err != nil {
			return err
		}
		from = idx.Block(n)
	}
	to := idx.Block(0)
	if len(ctx.Args()) > 1 {
		n, err := strconv.ParseUint(ctx.Args().Get(1), 10, 32)
		if err != nil {
			return err
		}
		to = idx.Block(n)
		if to <= from {
			return fmt.Errorf("to block num should be greater than from")
		}
	}

	for i := from; i < to || to <= from; i++ {
		block := gdb.GetBlock(i)
		if block == nil {
			return nil
		}
		fmt.Printf("%d\t%s\n", i, block.Root.String())
	}

	return nil
}
