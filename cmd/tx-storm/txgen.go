package main

import (
	"math"
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

	instances uint
	accs      []*Acc
	offset    uint
	position  uint

	work sync.WaitGroup
	done chan struct{}
	sync.Mutex

	logger.Instance
}

func NewTxGenerator(cfg *Config, num, ofTotal uint) *Generator {
	accs := cfg.Accs.Count / ofTotal
	offset := cfg.Accs.Offset + accs*(num-1)
	g := &Generator{
		chainId:   uint(cfg.ChainId),
		instances: ofTotal,
		accs:      make([]*Acc, accs),
		offset:    offset,

		Instance: logger.MakeInstance(),
	}

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

	g.Log.Info("will use", "accounts", len(g.accs), "from", g.offset, "to", uint(len(g.accs))+g.offset)
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
}

func (g *Generator) getTPS() float64 {
	tps := atomic.LoadUint32(&g.tps)
	return float64(tps)
}

func (g *Generator) SetTPS(tps float64) {
	x := uint32(math.Ceil(tps / float64(g.instances)))
	atomic.StoreUint32(&g.tps, x)
}

func (g *Generator) background(output chan<- *Transaction) {
	defer g.work.Done()
	defer close(output)

	g.Log.Info("started")
	defer g.Log.Info("stopped")

	for {
		begin := time.Now()
		var (
			generating time.Duration
			sending    time.Duration
		)

		tps := g.getTPS()
		for count := tps; count > 0; count-- {
			begin := time.Now()
			tx := g.Yield()
			generating += time.Since(begin)

			begin = time.Now()
			select {
			case output <- tx:
				sending += time.Since(begin)
				continue
			case <-g.done:
				return
			}
		}

		spent := time.Since(begin)
		if spent >= time.Second {
			g.Log.Warn("exceeded performance", "tps", tps, "generating", generating, "sending", sending)
			continue
		}

		select {
		case <-time.After(time.Second - spent):
			continue
		case <-g.done:
			return
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

	nonce := position / count
	amount := big.NewInt(1e6)

	tx := &Transaction{
		Raw:  from.TransactionTo(to, nonce, amount, g.chainId),
		Info: meta.NewInfo(a, b),
	}

	// g.Log.Info("regular tx", "from", a, "to", b, "amount", amount, "nonce", nonce)
	return tx
}
