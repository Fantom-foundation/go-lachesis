package main

import (
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/Fantom-foundation/go-lachesis/cmd/tx-storm/meta"
	"github.com/Fantom-foundation/go-lachesis/logger"
)

type Transaction struct {
	Raw  *types.Transaction
	Info *meta.Info
}

type Generator struct {
	tps     uint32
	chainId uint

	accs     []*Acc
	offset   uint
	position uint

	work sync.WaitGroup
	done chan struct{}
	sync.Mutex

	logger.Instance
}

func NewTxGenerator(cfg *Config) *Generator {
	g := &Generator{
		chainId: uint(cfg.ChainId),
		accs:    make([]*Acc, cfg.Accs.Count),
		offset:  cfg.Accs.Offset,

		Instance: logger.MakeInstance(),
	}

	g.Log.Info("Will use", "accounts", len(g.accs), "from", g.offset, "to", uint(len(g.accs))+g.offset)
	return g
}

func (g *Generator) Start() (output chan *Transaction) {
	g.Lock()
	defer g.Unlock()

	if g.done != nil {
		return
	}
	g.done = make(chan struct{})

	output = make(chan *Transaction, 100)
	g.work.Add(1)
	go g.background(output)

	g.Log.Info("started")
	return
}

func (g *Generator) Stop() {
	g.Lock()
	defer g.Unlock()

	if g.done == nil {
		return
	}

	close(g.done)
	g.work.Wait()
	g.done = nil

	g.Log.Info("stopped")
}

func (g *Generator) GetTPS() uint {
	tps := atomic.LoadUint32(&g.tps)
	return uint(tps)
}

func (g *Generator) SetTPS(tps uint) {
	atomic.StoreUint32(&g.tps, uint32(tps))
}

func (g *Generator) background(output chan<- *Transaction) {
	defer g.work.Done()
	defer close(output)

	for {
		start := time.Now()

		for count := g.GetTPS(); count > 0; count-- {
			tx := g.Yield()
			select {
			case output <- tx:
				break
			case <-g.done:
				return
			}
		}

		spent := time.Since(start)
		if spent < time.Second {
			select {
			case <-time.After(time.Second - spent):
				break
			case <-g.done:
				return
			}
		}
	}
}

func (g *Generator) Yield() *Transaction {
	tx := g.generate(g.position)
	g.position++

	return tx
}

func (g *Generator) generate(position uint) *Transaction {
	var count = uint(len(g.accs))

	a := position % count
	b := (position + 1) % count

	from := g.accs[a]
	if from == nil {
		from = MakeAcc(a + g.offset)
		g.accs[a] = from
	}
	a += g.offset

	to := g.accs[b]
	if to == nil {
		to = MakeAcc(b + g.offset)
		g.accs[b] = to
	}
	b += g.offset

	nonce := position/count + 1
	amount := big.NewInt(1e6)

	tx := &Transaction{
		Raw:  from.TransactionTo(to, nonce, amount, g.chainId),
		Info: meta.NewInfo(a, b),
	}

	// g.Log.Info("regular tx", "from", a, "to", b, "amount", amount, "nonce", nonce)
	return tx
}
