package gossip

import (
	"github.com/ethereum/go-ethereum/metrics"

	"github.com/Fantom-foundation/go-lachesis/cmd/tx-storm/meta"
)

var (
	epochGauge         = metrics.NewRegisteredGauge("epoch", nil)
	confirmBlocksMeter = metrics.NewRegisteredCounter("confirm/blocks", nil)
	confirmTxnsMeter   = metrics.NewRegisteredCounter("confirm/transactions", nil)
	txTtfMeter         = metrics.NewRegisteredHistogram("tx_ttf", nil, metrics.NewUniformSample(500))

	stakersCountGauge  = metrics.NewRegisteredGauge("stakers", nil)
	stakersStakeGauge  = metrics.NewRegisteredGauge("stakers/stake", nil)
)

var txLatency = meta.NewTxs()
