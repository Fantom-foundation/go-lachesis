package gossip

import (
	"github.com/Fantom-foundation/go-lachesis/app"
	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/gossip/gasprice"
	"github.com/Fantom-foundation/go-lachesis/lachesis"
)

type (
	// ProtocolConfig is config for p2p protocol
	ProtocolConfig struct {
		// 0/M means "optimize only for throughput", N/0 means "optimize only for latency", N/M is a balanced mode

		LatencyImportance    int
		ThroughputImportance int
	}
	// UpgradeConfig defines behaviour for network upgrades
	UpgradeConfig struct {
		ShutDownIfNotUpgraded bool // shut down the node in a case of non-supported network upgrade
		WarningIfNotUpgraded  bool // show a warning in a case of non-supported network upgrade
	}
	// Config for the gossip service.
	Config struct {
		Net     lachesis.Config
		Emitter EmitterConfig
		TxPool  evmcore.TxPoolConfig
		StoreConfig

		TxIndex             bool // Whether to enable indexing transactions and receipts or not
		DecisiveEventsIndex bool // Whether to enable indexing events which decide blocks or not
		EventLocalTimeIndex bool // Whether to enable indexing arrival time of events or not

		Upgrade UpgradeConfig

		// Protocol options
		Protocol ProtocolConfig

		// Gas Price Oracle options
		GPO gasprice.Config

		// Enables tracking of SHA3 preimages in the VM
		EnablePreimageRecording bool // TODO

		// Type of the EWASM interpreter ("" for default)
		EWASMInterpreter string

		// Type of the EVM interpreter ("" for default)
		EVMInterpreter string // TODO custom interpreter

		// RPCGasCap is the global gas cap for eth-call variants.
		RPCGasCap uint64 `toml:",omitempty"`

		// RPCTxFeeCap is the global transaction fee(price * gaslimit) cap for
		// send-transction variants. The unit is ether.
		RPCTxFeeCap float64 `toml:",omitempty"`

		ExtRPCEnabled bool

		RPCLogsBloom bool
	}

	// StoreConfig is a config for store db.
	StoreConfig struct {
		// Cache size for Events.
		EventsCacheSize int
		// Cache size for EventHeaderData (Epoch db).
		EventsHeadersCacheSize int
		// Cache size for Block.
		BlockCacheSize int
		// Cache size for PackInfos.
		PackInfosCacheSize int
		// Cache size for TxPositions.
		TxPositionsCacheSize int
		// Cache size for EpochStats.
		EpochStatsCacheSize int

		// NOTE: fields for config-file back compatibility
		// Cache size for Receipts.
		ReceiptsCacheSize int
		// Cache size for Stakers.
		StakersCacheSize int
		// Cache size for Delegations.
		DelegationsCacheSize int
	}
)

// DefaultConfig returns the default configurations for the gossip service.
func DefaultConfig(network lachesis.Config) Config {
	cfg := Config{
		Net:         network,
		Emitter:     DefaultEmitterConfig(),
		TxPool:      evmcore.DefaultTxPoolConfig(),
		StoreConfig: DefaultStoreConfig(),

		Upgrade: UpgradeConfig{
			ShutDownIfNotUpgraded: false,
			WarningIfNotUpgraded:  true,
		},

		TxIndex:             true,
		DecisiveEventsIndex: false,

		Protocol: ProtocolConfig{
			LatencyImportance:    60,
			ThroughputImportance: 40,
		},

		GPO: gasprice.Config{
			Blocks:     20,
			Percentile: 60,
			MaxPrice:   gasprice.DefaultMaxPrice,
		},
		RPCGasCap:    25000000,
		RPCTxFeeCap:  1000, // 1000 FTM
		RPCLogsBloom: true,
	}

	if network.NetworkID == lachesis.FakeNetworkID {
		cfg.Emitter = FakeEmitterConfig()
		// disable self-fork protection if fakenet 1/1
		if len(network.Genesis.Alloc.Validators) == 1 {
			cfg.Emitter.EmitIntervals.SelfForkProtection = 0
		}
	}
	/*if network.NetworkId == lachesis.DevNetworkId { // TODO dev network
		cfg.TxPool = evmcore.FakeTxPoolConfig()
		cfg.Emitter = FakeEmitterConfig()
	}*/

	return cfg
}

// DefaultStoreConfig for product.
func DefaultStoreConfig() StoreConfig {
	return StoreConfig{
		EventsCacheSize:        300,
		EventsHeadersCacheSize: 1000,
		BlockCacheSize:         100,
		PackInfosCacheSize:     100,
		TxPositionsCacheSize:   1000,
		EpochStatsCacheSize:    100,
		ReceiptsCacheSize:      app.DefaultStoreConfig().ReceiptsCacheSize,
		StakersCacheSize:       app.DefaultStoreConfig().StakersCacheSize,
		DelegationsCacheSize:   app.DefaultStoreConfig().DelegationsCacheSize,
	}
}

// LiteStoreConfig is for tests or inmemory.
func LiteStoreConfig() StoreConfig {
	return StoreConfig{
		EventsCacheSize:        100,
		EventsHeadersCacheSize: 1000,
		BlockCacheSize:         100,
		PackInfosCacheSize:     100,
		TxPositionsCacheSize:   100,
		EpochStatsCacheSize:    100,
		ReceiptsCacheSize:      app.LiteStoreConfig().ReceiptsCacheSize,
		StakersCacheSize:       app.LiteStoreConfig().StakersCacheSize,
		DelegationsCacheSize:   app.LiteStoreConfig().DelegationsCacheSize,
	}
}
