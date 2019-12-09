package main

import (
	"github.com/ethereum/go-ethereum/metrics"
	cli "gopkg.in/urfave/cli.v1"

	"github.com/Fantom-foundation/go-lachesis/metrics/prometheus"
)

var MetricsPrometheusEndpointFlag = cli.StringFlag{
	Name:  "metrics.prometheus.endpoint",
	Usage: "Prometheus API endpoint to report metrics to",
	Value: ":19090",
}

var (
	reg = metrics.NewRegistry()

	txCountSentMeter = metrics.NewRegisteredCounter("tx_count_sent", reg)
	txCountGotMeter  = metrics.NewRegisteredCounter("tx_count_got", reg)
	txLatencyMeter   = metrics.NewRegisteredHistogram("tx_ttf", reg, metrics.NewUniformSample(500))

	txBenchTimeMeter = metrics.NewRegisteredTimer("tx_bench_time", reg)
	txInputTpsMeter = metrics.NewRegisteredGaugeFloat64("tx_input_tps", reg)
	txTpsMeter = metrics.NewRegisteredGaugeFloat64("tx_avg_tps", reg)
	txTtfMeter = metrics.NewRegisteredGaugeFloat64("tx_avg_ttf", reg)
)

func SetupPrometheus(ctx *cli.Context) {
	if !metrics.Enabled {
		return
	}

	prometheus.SetNamespace("txstorm")
	var endpoint = ctx.GlobalString(MetricsPrometheusEndpointFlag.Name)
	prometheus.ListenTo(endpoint, reg)
}
