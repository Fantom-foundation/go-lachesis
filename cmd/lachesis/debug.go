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
		Usage:     "Prints root hashes for states of blockchain",
		ArgsUsage: "[<epochFrom> [<epochTo>]]",
		Category:  "MISCELLANEOUS COMMANDS",
		Description: `
For debug purpose.
Optional first and second arguments control the first and last epoch to export.`,
	}
)

func printStateRoots(ctx *cli.Context) error {
	cfg := makeAllConfigs(ctx)

	// check errlock file
	errlock.SetDefaultDatadir(cfg.Node.DataDir)
	errlock.Check()

	gdb := makeGossipStore(cfg.Node.DataDir, &cfg.Lachesis)
	defer gdb.Close()

	from := idx.Epoch(0)
	if len(ctx.Args()) > 0 {
		n, err := strconv.ParseUint(ctx.Args().Get(1), 10, 32)
		if err != nil {
			return err
		}
		from = idx.Epoch(n)
	}
	to := idx.Epoch(0)
	if len(ctx.Args()) > 1 {
		n, err := strconv.ParseUint(ctx.Args().Get(2), 10, 32)
		if err != nil {
			return err
		}
		to = idx.Epoch(n)
	}

	fmt.Printf("from %d epoch to %d\n", from, to)
	// TODO: print

	return nil
}
