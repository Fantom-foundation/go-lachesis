package main

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"gopkg.in/urfave/cli.v1"

	"github.com/Fantom-foundation/go-lachesis/gossip"
	"github.com/Fantom-foundation/go-lachesis/integration"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/kvdb/flushable"
	"github.com/Fantom-foundation/go-lachesis/utils/errlock"
)

// statsReportLimit is the time limit during import and export after which we
// always print out progress. This avoids the user wondering what's going on.
const statsReportLimit = 8 * time.Second

func exportChain(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		utils.Fatalf("This command requires an argument.")
	}

	cfg := makeAllConfigs(ctx)

	// check errlock file
	errlock.SetDefaultDatadir(cfg.Node.DataDir)
	errlock.Check()

	gdb := makeGossipStore(cfg.Node.DataDir, &cfg.Lachesis)
	defer gdb.Close()

	start := time.Now()

	fn := ctx.Args().First()

	// Open the file handle and potentially wrap with a gzip stream
	fh, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer fh.Close()

	var writer io.Writer = fh
	if strings.HasSuffix(fn, ".gz") {
		writer = gzip.NewWriter(writer)
		defer writer.(*gzip.Writer).Close()
	}

	from := idx.Epoch(0)
	if len(ctx.Args()) > 1 {
		n, err := strconv.ParseUint(ctx.Args().Get(1), 10, 32)
		if err != nil {
			return err
		}
		from = idx.Epoch(n)
	}
	to := idx.Epoch(0)
	if len(ctx.Args()) > 2 {
		n, err := strconv.ParseUint(ctx.Args().Get(2), 10, 32)
		if err != nil {
			return err
		}
		to = idx.Epoch(n)
	}

	err = exportTo(writer, gdb, from, to)
	if err != nil {
		utils.Fatalf("Export error: %v\n", err)
	}

	fmt.Printf("Export done in %v\n", time.Since(start))
	return nil
}

func makeGossipStore(dataDir string, gossipCfg *gossip.Config) *gossip.Store {
	dbs := flushable.NewSyncedPool(integration.DBProducer(dataDir))
	gdb := gossip.NewStore(dbs, gossipCfg.StoreConfig)
	gdb.SetName("gossip-db")
	return gdb
}

// exportTo writer the active chain.
func exportTo(w io.Writer, gdb *gossip.Store, from, to idx.Epoch) error {
	log.Info("Exporting batch of events")
	start, reported := time.Now(), time.Now()

	var err error
	gdb.ForEachEvents(func(event *inter.Event) bool {
		if event == nil {
			err = errors.New("export failed, event not found")
			return false
		}

		if event.Epoch < from {
			return true
		}
		if to > from && event.Epoch > to {
			return (event.Epoch - to) < 1
		}

		err := event.EncodeRLP(w)
		if err != nil {
			err = fmt.Errorf("export failed, error: %v", err)
			return false
		}
		log.Debug("exported", "event", event.String())

		if time.Since(reported) >= statsReportLimit {
			log.Info("Exporting events", "exported", event.String(), "elapsed", common.PrettyDuration(time.Since(start)))
			reported = time.Now()
		}
		return true
	})

	return err
}
