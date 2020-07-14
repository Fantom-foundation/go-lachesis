package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/Fantom-foundation/go-lachesis/logger"
	"github.com/Fantom-foundation/go-lachesis/utils"
)

// Nodes pool.
type Nodes struct {
	tps    chan float64
	conns  []*Sender
	blocks chan Block
	done   chan struct{}
	cfg    *Config
	logger.Instance
}

func NewNodes(cfg *Config, input <-chan *Transaction) *Nodes {
	n := &Nodes{
		tps:      make(chan float64, 1),
		blocks:   make(chan Block, 1),
		done:     make(chan struct{}),
		cfg:      cfg,
		Instance: logger.MakeInstance(),
	}

	was := make(map[string]struct{}, len(cfg.URLs))
	for _, url := range cfg.URLs {
		_, double := was[url]
		n.add(url, !double)
		was[url] = struct{}{}
	}

	n.notifyTPS(0)
	go n.background(input)
	go n.measureTPS()
	return n
}

func (n *Nodes) Count() int {
	return len(n.conns)
}

func (n *Nodes) TPS() <-chan float64 {
	return n.tps
}

func (n *Nodes) notifyTPS(tps float64) {
	select {
	case n.tps <- tps:
		break
	default:
	}
}

func (n *Nodes) measureTPS() {
	var (
		lastBlock *big.Int
		avgbuff   = utils.NewAvgBuff(10)
		start     = time.Unix(1, 0)
	)
	for b := range n.blocks {
		if lastBlock != nil && b.Number.Cmp(lastBlock) < 1 {
			continue
		}

		txCountGotMeter.Inc(int64(b.TxCount))

		dur := time.Since(start).Seconds()
		tps := float64(b.TxCount) / dur
		avgbuff.Push(float64(b.TxCount), dur)

		txTpsMeter.Update(int64(tps))

		start = time.Now()
		lastBlock = b.Number
		avg := avgbuff.Avg()
		n.notifyTPS(avg)
		n.Log.Info("TPS", "block", b.Number, "value", tps, "avg", avg)
	}
}

func (n *Nodes) add(url string, is1st bool) {
	c := NewSender(url, is1st, n.cfg.SendTrusted)
	c.SetName(fmt.Sprintf("Node-%d", len(n.conns)))
	n.conns = append(n.conns, c)

	go func() {
		defer n.stop()
		for b := range c.Blocks() {
			n.blocks <- b
		}
	}()
}

func (n *Nodes) stop() {
	// TODO: mutex
	close(n.blocks)
}

func (n *Nodes) background(input <-chan *Transaction) {
	if len(n.conns) < 1 {
		panic("no connections")
	}

	i := 0
	for tx := range input {
		c := n.conns[i]
		c.Send(tx)
		i = (i + 1) % len(n.conns)
	}

	for _, c := range n.conns {
		c.Close()
	}
}
