package main

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/logger"
)

type Sender struct {
	url     string
	input   chan *Transaction
	headers chan *types.Header
	output  chan Block

	cfg struct {
		listenToBlocks bool
		sendTrustedTx  bool
	}

	done chan struct{}
	work sync.WaitGroup

	logger.Instance
}

type Block struct {
	Number  *big.Int
	TxCount uint
}

func NewSender(url string, listenToBlocks, sendTrustedTx bool) *Sender {
	s := &Sender{
		url:     url,
		input:   make(chan *Transaction, 10),
		headers: make(chan *types.Header, 1),
		output:  make(chan Block, 1),
		done:    make(chan struct{}),

		cfg: struct {
			listenToBlocks bool
			sendTrustedTx  bool
		}{listenToBlocks, sendTrustedTx},

		Instance: logger.MakeInstance(),
	}

	s.work.Add(1)
	go s.background()

	return s
}

func (s *Sender) Close() {
	if s.done == nil {
		return
	}
	close(s.done)
	s.done = nil

	s.work.Wait()
	close(s.output)
	close(s.input)
}

func (s *Sender) Send(tx *Transaction) {
	s.input <- tx
}

func (s *Sender) Blocks() <-chan Block {
	return s.output
}

func (s *Sender) background() {
	defer s.work.Done()
	s.Log.Info("started")
	defer s.Log.Info("stopped")

	var (
		client *ethclient.Client
		err    error
		tx     *Transaction
		info   string
		sbscr  ethereum.Subscription
	)

	disconnect := func() {
		if sbscr != nil {
			sbscr.Unsubscribe()
			sbscr = nil
		}
		if client != nil {
			client.Close()
			client = nil
			s.Log.Error("disonnect from", "url", s.url)
		}
	}
	defer disconnect()

	for {

		for tx == nil {
			select {
			case tx = <-s.input:
				info = tx.Info.String()
			case b := <-s.headers:
				err = s.countTxs(client, b)
				if err != nil {
					disconnect()
				}
			case <-s.done:
				return
			}
		}

		for client == nil {
			client = s.connect()
		}

		if s.cfg.listenToBlocks && sbscr == nil {
			sbscr = s.subscribe(client)
			if sbscr == nil {
				disconnect()
				continue
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		if s.cfg.sendTrustedTx {
			err = client.SendTrustedTransaction(ctx, tx.Raw)
		} else {
			err = client.SendTransaction(ctx, tx.Raw)
		}
		cancel()
		if err == nil {
			txCountSentMeter.Inc(1)
			s.Log.Debug("tx sending ok", "info", info, "amount", tx.Raw.Value(), "nonce", tx.Raw.Nonce())
			tx = nil
			continue
		}
		switch err.Error() {
		case fmt.Sprintf("known transaction: %x", tx.Raw.Hash()),
			evmcore.ErrNonceTooLow.Error(),
			evmcore.ErrReplaceUnderpriced.Error():
			s.Log.Warn("tx sending skip", "info", info, "amount", tx.Raw.Value(), "cause", err, "nonce", tx.Raw.Nonce())
			tx = nil
			continue
		default:
			s.Log.Error("tx sending err", "info", info, "amount", tx.Raw.Value(), "cause", err, "nonce", tx.Raw.Nonce())
			disconnect()
			s.delay()
			continue
		}
	}
}

func (s *Sender) connect() *ethclient.Client {
	client, err := ethclient.Dial(s.url)
	if err != nil {
		s.Log.Error("connect to", "url", s.url, "err", err)
		s.delay()
		return nil
	}
	s.Log.Info("connect to", "url", s.url)
	return client
}

func (s *Sender) subscribe(client *ethclient.Client) ethereum.Subscription {
	sbscr, err := client.SubscribeNewHead(context.Background(), s.headers)
	if err != nil {
		s.Log.Error("subscribe to", "url", s.url, "err", err)
		s.delay()
		return nil
	}
	s.Log.Info("subscribe to", "url", s.url)
	return sbscr
}

func (s *Sender) countTxs(client *ethclient.Client, h *types.Header) error {
	b := evmcore.ConvertFromEthHeader(h)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	txs, err := client.TransactionCount(ctx, b.Hash)
	cancel()
	if err != nil {
		s.Log.Error("new block", "number", b.Number, "block", b.Hash, "err", err)
		return err
	}

	s.output <- Block{
		Number:  b.Number,
		TxCount: txs,
	}
	return nil
}

func (s *Sender) delay() {
	select {
	case <-time.After(2 * time.Second):
		return
	case <-s.done:
		return
	}
}
