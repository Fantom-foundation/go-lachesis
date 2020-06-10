package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"gopkg.in/urfave/cli.v1"

	"github.com/Fantom-foundation/go-lachesis/gossip"
	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
)

func importChain(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		utils.Fatalf("This command requires an argument.")
	}

	cfg := makeAllConfigs(ctx)
	err := importToNode(ctx, &cfg, ctx.Args()...)
	if err != nil {
		return err
	}

	return nil
}

func importToNode(ctx *cli.Context, cfg *config, files ...string) error {
	cfg.Lachesis.Emitter.Validator = common.Address{}
	cfg.Lachesis.TxPool.Journal = ""
	cfg.Node.IPCPath = ""
	cfg.Node.HTTPHost = ""
	cfg.Node.WSHost = ""
	cfg.Node.NoUSB = true
	cfg.Node.P2P.ListenAddr = ""
	cfg.Node.P2P.NoDiscovery = true
	cfg.Node.P2P.BootstrapNodes = nil
	cfg.Node.P2P.DiscoveryV5 = false
	cfg.Node.P2P.BootstrapNodesV5 = nil
	cfg.Node.P2P.StaticNodes = nil
	cfg.Node.P2P.TrustedNodes = nil

	node := makeNode(ctx, cfg)
	defer node.Close()
	startNode(ctx, node)

	var srv *gossip.Service
	if err := node.Service(&srv); err != nil {
		return err
	}

	fmt.Printf("Import starting.\n")
	start := time.Now()
	for _, fn := range files {
		if err := importFile(srv, fn); err != nil {
			log.Error("Import error", "file", fn, "err", err)
		}
	}
	fmt.Printf("Import done in %v.\n", time.Since(start))
	return nil
}

func importFile(srv *gossip.Service, fn string) error {
	// Watch for Ctrl-C while the import is running.
	// If a signal is received, the import will stop.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	log.Info("importing events", "file", fn)

	// Open the file handle and potentially unwrap the gzip stream
	fh, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer fh.Close()

	var reader io.Reader = fh
	if strings.HasSuffix(fn, ".gz") {
		if reader, err = gzip.NewReader(reader); err != nil {
			return err
		}
	}

	stream := rlp.NewStream(reader, 0)

	var h common.Hash
	if err = stream.Decode(&h); err == io.EOF {
		return nil
	}

	genesis := srv.GetEvmStateReader().GetDagBlock(hash.Event{}, 0)
	if genesis == nil {
		return fmt.Errorf("cann't init db")
	}
	if genesis.Hash != h {
		log.Warn("Incompatible genesis event", "current", genesis.Hash.String(), "want", h.String())
		return fmt.Errorf("incompatible genesis event")
	}

	for {
		select {
		case <-interrupt:
			return fmt.Errorf("interrupted")
		default:
		}

		var e inter.Event
		if err = stream.Decode(&e); err == io.EOF {
			break
		}

		if err = srv.ImportEvent(&e); err != nil {
			return err
		}
		log.Debug("event inserted", "event", e.Hash())
	}

	return nil
}
