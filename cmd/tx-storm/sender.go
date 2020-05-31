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

	done chan struct{}
	work sync.WaitGroup

	logger.Instance
}

type Block struct {
	Number  *big.Int
	TxCount uint
}

func NewSender(url string) *Sender {
	s := &Sender{
		url:     url,
		input:   make(chan *Transaction, 1),
		headers: make(chan *types.Header, 1),
		output:  make(chan Block, 1),
		done:    make(chan struct{}),

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
			s.Log.Error("Disonnect from", "url", s.url)
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

		if sbscr == nil {
			sbscr = s.subscribe(client)
			if sbscr == nil {
				disconnect()
				continue
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		err = client.SendTransaction(ctx, tx.Raw)
		cancel()
		if err == nil {
			s.Log.Info("Tx sending ok", "info", info, "amount", tx.Raw.Value(), "nonce", tx.Raw.Nonce())
			tx = nil
			continue
		}

		switch err.Error() {
		case fmt.Sprintf("known transaction: %x", tx.Raw.Hash()):
			s.Log.Info("Tx sending skip", "info", info, "amount", tx.Raw.Value(), "cause", err, "nonce", tx.Raw.Nonce())
			tx = nil
			continue
		case evmcore.ErrNonceTooLow.Error():
			s.Log.Info("Tx sending skip", "info", info, "amount", tx.Raw.Value(), "cause", err, "nonce", tx.Raw.Nonce())
			tx = nil
			continue
		case evmcore.ErrReplaceUnderpriced.Error():
			s.Log.Info("Tx sending skip", "info", info, "amount", tx.Raw.Value(), "cause", err, "nonce", tx.Raw.Nonce())
			tx = nil
			continue
		default:
			s.Log.Error("Tx sending err", "info", info, "amount", tx.Raw.Value(), "cause", err, "nonce", tx.Raw.Nonce())
			s.delay()
			continue
		}
	}
}

func (s *Sender) connect() *ethclient.Client {
	client, err := ethclient.Dial(s.url)
	s.Log.Error("Connect to", "url", s.url, "err", err)
	if err != nil {
		s.delay()
		return nil
	}
	return client
}

func (s *Sender) subscribe(client *ethclient.Client) ethereum.Subscription {
	sbscr, err := client.SubscribeNewHead(context.Background(), s.headers)
	s.Log.Error("Subscribe to", "url", s.url, "err", err)
	if err != nil {
		s.delay()
		return nil
	}
	return sbscr
}

func (s *Sender) countTxs(client *ethclient.Client, b *types.Header) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	txs, err := client.TransactionCount(ctx, b.Hash())
	cancel()
	if err != nil {
		s.Log.Error("New block", "number", b.Number, "block", b.Hash(), "err", err)
		return err
	}
	s.Log.Info("New block", "number", b.Number, "hash", b.Hash(), "txs", txs)

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
