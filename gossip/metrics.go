package gossip

import (
	"github.com/Fantom-foundation/go-ethereum/metrics"
)

var (
	confirmBlocksMeter = metrics.NewRegisteredCounter("confirm/blocks", nil)
	confirmTxnsMeter   = metrics.NewRegisteredCounter("confirm/transactions", nil)
	//confirmTimeMeter = metrics.NewRegisteredHistogram("confirm/seconds", nil)
)
