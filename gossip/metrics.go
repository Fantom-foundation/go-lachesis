package gossip

import (
	"time"
	
	"github.com/ethereum/go-ethereum/metrics"

	"github.com/Fantom-foundation/go-lachesis/cmd/tx-storm/meta"
)

var (
	confirmBlocksMeter = metrics.NewRegisteredCounter("confirm/blocks", nil)
	confirmTxnsMeter   = metrics.NewRegisteredCounter("confirm/transactions", nil)
	txTtfMeter         = metrics.NewRegisteredHistogram("tx_ttf", nil, metrics.NewUniformSample(500))

	nodeStartTime time.Time = time.Now()
	
	nodeBenchTimeMeter = metrics.NewRegisteredTimer("node/benchtime", nil)
	nodeInputTpsMeter = metrics.NewRegisteredGaugeFloat64("node/inputtps", nil)
	nodeTpsMeter = metrics.NewRegisteredGaugeFloat64("node/tps", nil)
	nodeBlockTpsMeter = metrics.NewRegisteredGaugeFloat64("node/blocktps", nil)
	nodeTtfMeter = metrics.NewRegisteredGaugeFloat64("node/ttf", nil)
)

var txLatency = meta.NewTxs()
