package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/Fantom-foundation/go-opera/cmd/opera/launcher"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p/discv5"
	"github.com/ethereum/go-ethereum/params"
	"github.com/naoina/toml"
	"gopkg.in/urfave/cli.v1"

	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/gossip"
	"github.com/Fantom-foundation/go-lachesis/gossip/gasprice"
	"github.com/Fantom-foundation/go-lachesis/lachesis"
)

var (
	dumpConfigCommand = cli.Command{
		Action:      utils.MigrateFlags(dumpConfig),
		Name:        "dumpconfig",
		Usage:       "Show configuration values",
		ArgsUsage:   "",
		Flags:       append(nodeFlags, testFlags...),
		Category:    "MISCELLANEOUS COMMANDS",
		Description: `The dumpconfig command shows configuration values.`,
	}

	configFileFlag = cli.StringFlag{
		Name:  "legacy_config",
		Usage: "Lachesis TOML configuration file",
	}

	operaConfigFileFlag = cli.StringFlag{
		Name:  "config",
		Usage: "Opera TOML configuration file",
	}

	// DataDirFlag defines directory to store Lachesis state and user's wallets
	DataDirFlag = utils.DirectoryFlag{
		Name:  "datadir",
		Usage: "Data directory for the databases and keystore",
		Value: utils.DirectoryString(DefaultDataDir()),
	}
)

// These settings ensure that TOML keys use the same names as Go struct fields.
var tomlSettings = toml.Config{
	NormFieldName: func(rt reflect.Type, key string) string {
		return key
	},
	FieldToKey: func(rt reflect.Type, field string) string {
		return field
	},
	MissingField: func(rt reflect.Type, field string) error {
		return fmt.Errorf("field '%s' is not defined in %s", field, rt.String())
	},
}

type config struct {
	Node     node.Config
	Lachesis gossip.Config
}

func loadAllConfigs(file string, cfg *config) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tomlSettings.NewDecoder(bufio.NewReader(f)).Decode(cfg)
	// Add file name to errors that have a line number.
	if _, ok := err.(*toml.LineError); ok {
		err = errors.New(file + ", " + err.Error())
	}
	if err != nil {
		return errors.New(fmt.Sprintf("TOML config file error: %v.\n"+
			"Use 'dumpconfig' command to get an example config file.\n"+
			"If node was recently upgraded and a previous network config file is used, then check updates for the config file.", err))
	}
	return err
}

func defaultLachesisConfig(ctx *cli.Context) lachesis.Config {
	var cfg lachesis.Config

	switch {
	case ctx.GlobalIsSet(FakeNetFlag.Name):
		_, accs, err := parseFakeGen(ctx.GlobalString(FakeNetFlag.Name))
		if err != nil {
			log.Crit("Invalid flag", "flag", FakeNetFlag.Name, "err", err)
		}
		cfg = lachesis.FakeNetConfig(accs)
	case ctx.GlobalBool(utils.LegacyTestnetFlag.Name):
		cfg = lachesis.TestNetConfig()
	default:
		cfg = lachesis.MainNetConfig()
	}

	return cfg
}

func setBootnodes(ctx *cli.Context, urls []string, cfg *node.Config) {
	for _, url := range urls {
		if url != "" {
			node, err := discv5.ParseNode(url)
			if err != nil {
				log.Error("Bootstrap URL invalid", "enode", url, "err", err)
				continue
			}
			cfg.P2P.BootstrapNodesV5 = append(cfg.P2P.BootstrapNodesV5, node)
		}
	}
}

func setDataDir(ctx *cli.Context, cfg *node.Config) {
	defaultDataDir := DefaultDataDir()

	switch {
	case ctx.GlobalIsSet(utils.DataDirFlag.Name):
		cfg.DataDir = ctx.GlobalString(utils.DataDirFlag.Name)
	case ctx.GlobalIsSet(FakeNetFlag.Name):
		_, accs, err := parseFakeGen(ctx.GlobalString(FakeNetFlag.Name))
		if err != nil {
			log.Crit("Invalid flag", "flag", FakeNetFlag.Name, "err", err)
		}
		cfg.DataDir = filepath.Join(defaultDataDir, fmt.Sprintf("fakenet-%d", len(accs.Accounts)))
	case ctx.GlobalBool(utils.LegacyTestnetFlag.Name):
		cfg.DataDir = filepath.Join(defaultDataDir, "testnet")
	default:
		cfg.DataDir = defaultDataDir
	}
}

func setGPO(ctx *cli.Context, cfg *gasprice.Config) {
	if ctx.GlobalIsSet(utils.LegacyGpoBlocksFlag.Name) {
		cfg.Blocks = ctx.GlobalInt(utils.LegacyGpoBlocksFlag.Name)
	}
	if ctx.GlobalIsSet(utils.GpoBlocksFlag.Name) {
		cfg.Blocks = ctx.GlobalInt(utils.GpoBlocksFlag.Name)
	}
	if ctx.GlobalIsSet(utils.LegacyGpoPercentileFlag.Name) {
		cfg.Percentile = ctx.GlobalInt(utils.LegacyGpoPercentileFlag.Name)
	}
	if ctx.GlobalIsSet(utils.GpoPercentileFlag.Name) {
		cfg.Percentile = ctx.GlobalInt(utils.GpoPercentileFlag.Name)
	}
}

func setTxPool(ctx *cli.Context, cfg *evmcore.TxPoolConfig) {
	if ctx.GlobalIsSet(utils.TxPoolLocalsFlag.Name) {
		locals := strings.Split(ctx.GlobalString(utils.TxPoolLocalsFlag.Name), ",")
		for _, account := range locals {
			if trimmed := strings.TrimSpace(account); !common.IsHexAddress(trimmed) {
				utils.Fatalf("Invalid account in --txpool.locals: %s", trimmed)
			} else {
				cfg.Locals = append(cfg.Locals, common.HexToAddress(account))
			}
		}
	}
	if ctx.GlobalIsSet(utils.TxPoolNoLocalsFlag.Name) {
		cfg.NoLocals = ctx.GlobalBool(utils.TxPoolNoLocalsFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolJournalFlag.Name) {
		cfg.Journal = ctx.GlobalString(utils.TxPoolJournalFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolRejournalFlag.Name) {
		cfg.Rejournal = ctx.GlobalDuration(utils.TxPoolRejournalFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolPriceLimitFlag.Name) {
		cfg.PriceLimit = ctx.GlobalUint64(utils.TxPoolPriceLimitFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolPriceBumpFlag.Name) {
		cfg.PriceBump = ctx.GlobalUint64(utils.TxPoolPriceBumpFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolAccountSlotsFlag.Name) {
		cfg.AccountSlots = ctx.GlobalUint64(utils.TxPoolAccountSlotsFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolGlobalSlotsFlag.Name) {
		cfg.GlobalSlots = ctx.GlobalUint64(utils.TxPoolGlobalSlotsFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolAccountQueueFlag.Name) {
		cfg.AccountQueue = ctx.GlobalUint64(utils.TxPoolAccountQueueFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolGlobalQueueFlag.Name) {
		cfg.GlobalQueue = ctx.GlobalUint64(utils.TxPoolGlobalQueueFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolLifetimeFlag.Name) {
		cfg.Lifetime = ctx.GlobalDuration(utils.TxPoolLifetimeFlag.Name)
	}
}

func gossipConfigWithFlags(ctx *cli.Context, src gossip.Config) gossip.Config {
	cfg := src

	// Avoid conflicting network flags
	utils.CheckExclusive(ctx, FakeNetFlag, utils.DeveloperFlag, utils.LegacyTestnetFlag)
	utils.CheckExclusive(ctx, FakeNetFlag, utils.DeveloperFlag, utils.ExternalSignerFlag) // Can't use both ephemeral unlocked and external signer

	setGPO(ctx, &cfg.GPO)
	setTxPool(ctx, &cfg.TxPool)

	if ctx.GlobalIsSet(utils.NetworkIdFlag.Name) {
		cfg.Net.NetworkID = ctx.GlobalUint64(utils.NetworkIdFlag.Name)
	}
	// TODO cache config
	//if ctx.GlobalIsSet(utils.CacheFlag.Name) || ctx.GlobalIsSet(utils.CacheDatabaseFlag.Name) {
	//	cfg.DatabaseCache = ctx.GlobalInt(utils.CacheFlag.Name) * ctx.GlobalInt(utils.CacheDatabaseFlag.Name) / 100
	//}
	//if ctx.GlobalIsSet(utils.CacheFlag.Name) || ctx.GlobalIsSet(CacheTrieFlag.Name) {
	//	cfg.TrieCleanCache = ctx.GlobalInt(utils.CacheFlag.Name) * ctx.GlobalInt(CacheTrieFlag.Name) / 100
	//}
	//if ctx.GlobalIsSet(utils.CacheFlag.Name) || ctx.GlobalIsSet(CacheGCFlag.Name) {
	//	cfg.TrieDirtyCache = ctx.GlobalInt(utils.CacheFlag.Name) * ctx.GlobalInt(CacheGCFlag.Name) / 100
	//}

	if ctx.GlobalIsSet(utils.VMEnableDebugFlag.Name) {
		cfg.EnablePreimageRecording = ctx.GlobalBool(utils.VMEnableDebugFlag.Name)
	}

	if ctx.GlobalIsSet(utils.EWASMInterpreterFlag.Name) {
		cfg.EWASMInterpreter = ctx.GlobalString(utils.EWASMInterpreterFlag.Name)
	}

	if ctx.GlobalIsSet(utils.EVMInterpreterFlag.Name) {
		cfg.EVMInterpreter = ctx.GlobalString(utils.EVMInterpreterFlag.Name)
	}
	if ctx.GlobalIsSet(utils.RPCGlobalGasCap.Name) {
		cfg.RPCGasCap = ctx.GlobalUint64(utils.RPCGlobalGasCap.Name)
	}
	if ctx.GlobalIsSet(utils.RPCGlobalTxFeeCap.Name) {
		cfg.RPCTxFeeCap = ctx.GlobalFloat64(utils.RPCGlobalTxFeeCap.Name)
	}

	return cfg
}

func nodeConfigWithFlags(ctx *cli.Context, cfg node.Config) node.Config {
	utils.SetNodeConfig(ctx, &cfg)

	if !ctx.GlobalIsSet(FakeNetFlag.Name) {
		setBootnodes(ctx, Bootnodes, &cfg)
	}
	setDataDir(ctx, &cfg)
	return cfg
}

func makeAllConfigs(ctx *cli.Context) *config {
	// Defaults (low priority)
	net := defaultLachesisConfig(ctx)
	cfg := config{Lachesis: gossip.DefaultConfig(net), Node: defaultNodeConfig()}

	// Load config file (medium priority)
	if file := ctx.GlobalString(configFileFlag.Name); file != "" {
		if err := loadAllConfigs(file, &cfg); err != nil {
			utils.Fatalf("%v", err)
		}
	}

	// Apply flags (high priority)
	cfg.Lachesis = gossipConfigWithFlags(ctx, cfg.Lachesis)
	cfg.Node = nodeConfigWithFlags(ctx, cfg.Node)

	return &cfg
}

func defaultNodeConfig() node.Config {
	cfg := NodeDefaultConfig
	cfg.Name = clientIdentifier
	cfg.Version = params.VersionWithCommit(gitCommit, gitDate)
	cfg.HTTPModules = append(cfg.HTTPModules, "eth", "ftm", "sfc", "web3")
	cfg.WSModules = append(cfg.WSModules, "eth", "ftm", "sfc", "web3")
	cfg.IPCPath = "lachesis.ipc"
	cfg.DataDir = DefaultDataDir()
	return cfg
}

// dumpConfig is the dumpconfig command.
func dumpConfig(ctx *cli.Context) error {
	cfg := makeAllConfigs(ctx)
	comment := ""

	out, err := tomlSettings.Marshal(&cfg)
	if err != nil {
		return err
	}

	log.Info("Dumping legacy Lachesis config file")

	dump := os.Stdout
	if ctx.NArg() > 0 {
		dump, err = os.OpenFile(ctx.Args().Get(0), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer dump.Close()
	}
	dump.WriteString(comment)
	dump.Write(out)

	log.Info("Dumping Opera config file")

	operaDatadir := path.Join(cfg.Node.DataDir, "opera")
	operaArgs := excludeArg(os.Args, "--"+validatorFlag.Name, true)
	operaArgs = excludeArg(operaArgs, "--"+configFileFlag.Name, true)
	operaArgs = excludeArg(operaArgs, "--"+DataDirFlag.Name, true)
	operaArgs = addFrontArgs(operaArgs, []string{"--" + DataDirFlag.Name, operaDatadir})
	operaArgs = excludeArg(operaArgs, "--"+utils.LegacyTestnetFlag.Name, false)

	metrics.Enabled = false
	return launcher.Launch(operaArgs)
}
