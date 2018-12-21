package proxy

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/poset"
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
	block := poset.NewBlock(0, 1, []byte{}, transactions)

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

		stateHash, err := proxy.CommitBlock(block)
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

func (p *TestProxy) CommitHandler(block poset.Block) ([]byte, error) {
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
	return goldStateHash(), nil
}

func goldStateHash() []byte {
	return []byte("statehash")
}

func goldSnapshot() []byte {
	return []byte("snapshot")
}
