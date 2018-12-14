package proxy

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/peers"
	"github.com/Fantom-foundation/go-lachesis/src/poset"
	"github.com/Fantom-foundation/go-lachesis/src/proxy/proto"
)

func TestInmemAppCalls(t *testing.T) {
	const (
		timeout    = 1 * time.Second
		errTimeout = "time is over"
	)

	proxy := NewTestProxy(t)

	transactions := [][]byte{
		[]byte("tx 1"),
		[]byte("tx 2"),
		[]byte("tx 3"),
	}
	block := poset.NewBlock(0, 1, []byte{},
		[]*peers.Peer{},
		transactions,
		[]*poset.InternalTransaction{
			poset.NewInternalTransaction(poset.TransactionType_PEER_ADD, *peers.NewPeer("peer1", "paris")),
			poset.NewInternalTransaction(poset.TransactionType_PEER_REMOVE, *peers.NewPeer("peer2", "london")),
		})

	t.Run("#1 Send tx", func(t *testing.T) {
		asserter := assert.New(t)

		txOrigin := []byte("the test transaction")

		go func() {
			select {
			case tx := <-proxy.SubmitCh():
				asserter.Equal(txOrigin, tx)
			case <-time.After(timeout):
				asserter.Fail(errTimeout)
			}
		}()

		proxy.SubmitTx(txOrigin)
	})

	t.Run("#2 Commit block", func(t *testing.T) {
		asserter := assert.New(t)

		stateHash, err := proxy.CommitBlock(*block)
		if asserter.NoError(err) {
			asserter.EqualValues(goldStateHash(), stateHash)
			asserter.EqualValues(transactions, proxy.transactions)
		}
	})

	t.Run("#3 Get snapshot", func(t *testing.T) {
		asserter := assert.New(t)

		snapshot, err := proxy.GetSnapshot(block.Index())
		if asserter.NoError(err) {
			asserter.Equal(goldSnapshot(), snapshot)
		}
	})

	t.Run("#4 Restore snapshot", func(t *testing.T) {
		asserter := assert.New(t)

		err := proxy.Restore(goldSnapshot())
		asserter.NoError(err)
	})
}

/*
 * staff
 */

type TestProxy struct {
	*InmemAppProxy
	transactions [][]byte
	logger       *logrus.Logger
}

func NewTestProxy(t *testing.T) *TestProxy {
	proxy := &TestProxy{
		transactions: [][]byte{},
		logger:       common.NewTestLogger(t),
	}

	proxy.InmemAppProxy = NewInmemAppProxy(proxy, proxy.logger)

	return proxy
}

func (p *TestProxy) CommitHandler(block poset.Block) (proto.Response, error) {
	p.logger.Debug("CommitBlock")
	p.transactions = append(p.transactions, block.Transactions()...)
	return goldStateHash(), nil
}

func (p *TestProxy) SnapshotHandler(blockIndex int64) ([]byte, error) {
	p.logger.Debug("GetSnapshot")
	return goldSnapshot(), nil
}

func (p *TestProxy) RestoreHandler(snapshot []byte) ([]byte, error) {
	p.logger.Debug("RestoreSnapshot")
	return goldStateHash().StateHash, nil
}

func goldStateHash() proto.Response {
	return proto.Response{StateHash: []byte("statehash")}
}

func goldSnapshot() []byte {
	return []byte("snapshot")
}
