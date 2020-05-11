package gossip

import (
	"github.com/ethereum/go-ethereum/metrics"

	"github.com/Fantom-foundation/go-lachesis/cmd/tx-storm/meta"
)

/*
+ Up time
+ Epoch
+ Time Since Last Epoch 8:36:36 ago
+ Epoch Fees 0.000484135 FTM
+ Total Transactions 64327
+ Total Delegators 2523
+ Total Validators 32
+ Total Staked : 234,336,199
+ Total Delegated : 915,447,526 (NEED current value)
+ Blocks generated => Blocks
+ Txns generated => Txns included in events
+ Total Staked per week (or per month): to plot a chart
+ Total Delegated per week (or per month): to plot a chart
+ Total Supply : 2,265,643,481
Total Accounts 3633
Capacity
Pending withdrawn
*/

var (
	txTtfMeter         = metrics.NewRegisteredHistogram("tx_ttf", nil, metrics.NewUniformSample(500))

	epochGauge         = metrics.NewRegisteredGauge("epoch", nil)
	epochFeeGauge      = metrics.NewRegisteredGauge("epoch/fee", nil)
	confirmBlocksMeter = metrics.NewRegisteredCounter("confirm/blocks", nil)
	confirmTxnsMeter   = metrics.NewRegisteredCounter("confirm/transactions", nil)

	stakersCountGauge  = metrics.NewRegisteredGauge("stakers", nil)
	stakersStakeGauge  = metrics.NewRegisteredGauge("stakers/stake", nil)
	appUptimeGauge     = metrics.NewRegisteredGauge("uptime", nil)
	epochTimeGauge	   = metrics.NewRegisteredGauge("epoch/time", nil)
	delegatorsCountGauge = metrics.NewRegisteredGauge("delegators", nil)
	delegatorsAmountGauge = metrics.NewRegisteredGauge("delegators/amount", nil)
	validatorsCountGauge = metrics.NewRegisteredGauge("validators", nil)
	totalSupplyGauge   = metrics.NewRegisteredGauge("total_supply", nil)
)

var txLatency = meta.NewTxs()
