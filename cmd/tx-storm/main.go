package main

import (
	"fmt"
	"os"
	"os/signal"
	"sort"
	"syscall"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"gopkg.in/urfave/cli.v1"

	"github.com/Fantom-foundation/go-lachesis/flags"
	_ "github.com/Fantom-foundation/go-lachesis/version"
)

var (
	// Git SHA1 commit hash of the release (set via linker flags).
	gitCommit = ""
	gitDate   = ""
	// The App that holds all commands and flags.
	App = flags.NewApp(gitCommit, gitDate, "the transactions generator CLI")

	Flags []cli.Flag
)

// init the CLI app.
func init() {
	App.Action = generatorMain
	App.Version = params.VersionWithCommit(gitCommit, gitDate)

	App.Commands = []cli.Command{}
	sort.Sort(cli.CommandsByName(App.Commands))

	App.Flags = append(App.Flags,
		ConfigFileFlag,
		NumberFlag,
		TxnsRateFlag,
		utils.MetricsEnabledFlag,
		MetricsPrometheusEndpointFlag,
		VerbosityFlag,
	)
}

func main() {
	if err := App.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// generatorMain is the main entry point.
func generatorMain(ctx *cli.Context) error {
	glogger := log.NewGlogHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(false)))
	glogger.Verbosity(log.Lvl(ctx.GlobalInt(VerbosityFlag.Name)))
	log.Root().SetHandler(glogger)

	SetupPrometheus(ctx)

	cfg := OpenConfig(ctx)
	for i := 0; i < ctx.NArg(); i++ {
		cfg.URLs = append(cfg.URLs, ctx.Args().Get(i))
	}
	if ctx.GlobalIsSet(AccsStartFlag.Name) {
		cfg.Accs.Offset = uint(ctx.GlobalInt(AccsStartFlag.Name))
	}
	if ctx.GlobalIsSet(AccsCountFlag.Name) {
		cfg.Accs.Count = uint(ctx.GlobalInt(AccsCountFlag.Name))
	}

	num, ofTotal := getNumber(ctx)

	generator := NewTxGenerator(cfg, num, ofTotal)
	defer generator.Stop()
	generator.SetName(fmt.Sprintf("TxGen-%d", num))
	txs := generator.Start()

	nodes := NewNodes(cfg, txs)
	go func() {
		for tps := range nodes.TPS() {
			generator.SetTPS(tps + 50.0*float64(nodes.Count()))
		}
	}()

	waitForSignal()
	return nil
}

func waitForSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	<-sigs
}
