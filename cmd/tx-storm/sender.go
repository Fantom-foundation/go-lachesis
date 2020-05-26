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
	url    string
	input  chan *Transaction
	blocks chan big.Int

	done chan struct{}
	work sync.WaitGroup

	logger.Instance
}

func NewSender(url string) *Sender {
	s := &Sender{
		url:    url,
		input:  make(chan *Transaction, 1),
		blocks: make(chan big.Int, 1),
		done:   make(chan struct{}),

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
	close(s.input)
	s.input = nil
}

func (s *Sender) Send(tx *Transaction) {
	s.input <- tx
}

func (s *Sender) Notify(bnum big.Int) {
	select {
	case s.blocks <- bnum:
	default:
	}
}

func (s *Sender) background() {
	defer s.work.Done()
	defer s.Log.Info("Stopped")

	var (
		client *ethclient.Client
		err    error
		tx     *Transaction
		info   string
		sbscr  ethereum.Subscription
		blocks = make(chan *types.Header, 1)
	)

	for {
		select {
		case tx = <-s.input:
			info = tx.Info.String()
		case <-s.done:
			return
		}

	connecting:
		for client == nil {
			client, err = ethclient.Dial(s.url)
			if err != nil {
				client = nil
				s.Log.Error("Connect to", "url", s.url, "err", err)
				select {
				case <-time.After(time.Second):
					continue connecting
				case <-s.done:
					return
				}
			}

			sbscr, err = c.SubscribeNewHead(context.Background(), blocks)
			if err != nil {
				sbscr = nil
				client.Close()
				client = nil
				s.Log.Error("Subscribe to", "url", s.url, "err", err)
				select {
				case <-time.After(time.Second):
					continue connecting
				case <-s.done:
					return
				}
			}
			defer sbscr.Unsubscribe()

		}

	sending:
		for {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			err = client.SendTransaction(ctx, tx.Raw)
			cancel()
			if err == nil {
				s.Log.Info("Tx sending ok", "info", info, "amount", tx.Raw.Value(), "nonce", tx.Raw.Nonce())
				break sending
			}

			switch err.Error() {
			case fmt.Sprintf("known transaction: %x", tx.Raw.Hash()):
				s.Log.Info("Tx sending skip", "info", info, "amount", tx.Raw.Value(), "cause", err, "nonce", tx.Raw.Nonce())
				break sending
			case evmcore.ErrNonceTooLow.Error():
				s.Log.Info("Tx sending skip", "info", info, "amount", tx.Raw.Value(), "cause", err, "nonce", tx.Raw.Nonce())
				break sending
			case evmcore.ErrReplaceUnderpriced.Error():
				s.Log.Info("Tx sending skip", "info", info, "amount", tx.Raw.Value(), "cause", err, "nonce", tx.Raw.Nonce())
				break sending
			default:
				s.Log.Error("Tx sending err", "info", info, "amount", tx.Raw.Value(), "cause", err, "nonce", tx.Raw.Nonce())
				select {
				case <-s.blocks:
					s.Log.Error("Try to send tx again", "info", info, "amount", tx.Raw.Value(), "nonce", tx.Raw.Nonce())
				case <-s.done:
					return
				}
			}
		}

	}
}
