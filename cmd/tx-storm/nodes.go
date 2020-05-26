package main

import (
	"fmt"
	"sync/atomic"

	"github.com/Fantom-foundation/go-lachesis/logger"
)

// Nodes pool.
type Nodes struct {
	tps   chan uint
	conns []*Sender
	logger.Instance
}

func NewNodes(cfg *Config, input <-chan *Transaction) *Nodes {
	n := &Nodes{
		tps:      make(chan uint),
		Instance: logger.MakeInstance(),
	}
	for _, url := range cfg.URLs {
		n.add(url)
	}

	go n.background(input)
	return n
}

func (n *Nodes) TPS() <-chan uint {
	return n.tps
}

func (n *Nodes) notifyTPS(tps uint) {
	select {
	case n.tps < -tps:
		break
	default:
	}
}

func (n *Nodes) add(url string) {
	c := NewSender(url)
	c.SetName(fmt.Sprintf("Node-%d", len(n.conns)))
	n.conns = append(n.conns, c)
}

func (n *Nodes) background(input <-chan *Transaction) {
	i := 0
	for tx := range input {
		if i >= len(n.conns) {
			continue
		}
		c := n.conns[i]
		c.Send(tx)
		i = (i + 1) % len(n.conns)
	}

	for _, c := range n.conns {
		c.Close()
	}
	n.conns = nil
}
