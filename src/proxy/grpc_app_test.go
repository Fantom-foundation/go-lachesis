package proxy

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/network"
	"github.com/Fantom-foundation/go-lachesis/src/poset"
	"github.com/Fantom-foundation/go-lachesis/src/proxy/proto"
	"github.com/fortytw2/leaktest"
)

func TestGrpcAppCalls(t *testing.T) {
	t.Run("over TCP", func(t *testing.T) {
		testGrpcAppCalls(t, network.TCPListener)
	})

	t.Run("over Fake", func(t *testing.T) {
		dialer := network.FakeDialer("client.fake")
		testGrpcAppCalls(t, network.FakeListener, grpc.WithContextDialer(dialer))
	})
}

func TestGrpcAppReconnect(t *testing.T) {
	t.Run("over TCP", func(t *testing.T) {
		testGrpcAppReconnect(t, network.TCPListener)
	})

	t.Run("over Fake", func(t *testing.T) {
		dialer := network.FakeDialer("client.fake")
		testGrpcAppReconnect(t, network.FakeListener, grpc.WithContextDialer(dialer))
	})
}

func testGrpcAppCalls(t *testing.T, listen network.ListenFunc, opts ...grpc.DialOption) {
    ctx, _ := context.WithTimeout(context.Background(), time.Second)
    defer leaktest.CheckContext(ctx, t)()

	const (
		timeout    = 1 * time.Second
		errTimeout = "time is over"
	)

	logger := common.NewTestLogger(t)

	s, addr, err := NewGrpcAppProxy("127.0.0.1:", timeout, logger, listen)
	if !assert.NoError(t, err) {
		return
	}
	defer s.Close()

	c, err := NewGrpcLachesisProxy(addr, logger, opts...)
	if !assert.NoError(t, err) {
		return
	}
	defer c.Close()

	t.Run("#1 Send tx", func(t *testing.T) {
		assertar := assert.New(t)
		gold := []byte("123456")

		err := c.SubmitTx(gold)
		assertar.NoError(err)

		select {
		case tx := <-s.SubmitCh():
			assertar.Equal(gold, tx)
		case <-time.After(timeout):
			assertar.Fail(errTimeout)
		}
	})

	t.Run("#2 Receive block", func(t *testing.T) {
		assertar := assert.New(t)
		block := poset.Block{}
		gold := []byte("123456")

		go func() {
			select {
			case event := <-c.CommitCh():
				assertar.Equal(block, event.Block)
				event.RespChan <- proto.CommitResponse{
					StateHash: gold,
					Error:     nil,
				}
			case <-time.After(timeout):
				assertar.Fail(errTimeout)
			}
		}()

		answer, err := s.CommitBlock(block)
		if assertar.NoError(err) {
			assertar.Equal(gold, answer)
		}
	})

	t.Run("#3 Receive snapshot query", func(t *testing.T) {
		assertar := assert.New(t)
		index := int64(1)
		gold := []byte("123456")

		go func() {
			select {
			case event := <-c.SnapshotRequestCh():
				assertar.Equal(index, event.BlockIndex)
				event.RespChan <- proto.SnapshotResponse{
					Snapshot: gold,
					Error:    nil,
				}
			case <-time.After(timeout):
				assertar.Fail(errTimeout)
			}
		}()

		answer, err := s.GetSnapshot(index)
		if assertar.NoError(err) {
			assertar.Equal(gold, answer)
		}
	})

	t.Run("#4 Receive restore command", func(t *testing.T) {
		assertar := assert.New(t)
		gold := []byte("123456")

		go func() {
			select {
			case event := <-c.RestoreCh():
				assertar.Equal(gold, event.Snapshot)
				event.RespChan <- proto.RestoreResponse{
					StateHash: gold,
					Error:     nil,
				}
			case <-time.After(timeout):
				assertar.Fail(errTimeout)
			}
		}()

		err := s.Restore(gold)
		assertar.NoError(err)
	})
}

func testGrpcAppReconnect(t *testing.T, listen network.ListenFunc, opts ...grpc.DialOption) {
	const (
		timeout    = 1 * time.Second
		errTimeout = "time is over"
	)
	assertar := assert.New(t)

	logger := common.NewTestLogger(t)

	s, addr, err := NewGrpcAppProxy("127.0.0.1:", timeout, logger, listen)
	if !assertar.NoError(err) {
		return
	}

	c, err := NewGrpcLachesisProxy(addr, logger, opts...)
	if !assertar.NoError(err) {
		return
	}
	defer c.Close()

	checkConn := func(t *testing.T) {
		assertar := assert.New(t)

		gold := []byte("123456")

		err := c.SubmitTx(gold)
		if !assertar.NoError(err) {
			return
		}

		select {
		case tx := <-s.SubmitCh():
			assertar.Equal(gold, tx)
		case <-time.After(timeout):
			assertar.Fail(errTimeout)
		}
	}

	t.Run("#1 Send tx after connection", checkConn)

	s.Close()
	s, _, err = NewGrpcAppProxy(addr, timeout/2, logger, listen)
	if !assertar.NoError(err) {
		return
	}
	defer s.Close()

	<-time.After(timeout)
	t.Run("#2 Send tx after reconnection", checkConn)
}

// TODO: fix it
/*
func TestGrpcMaxMsgSize(t *testing.T) {
	const (
		largeSize  = 100 * 1024 * 1024
		timeout    = 3 * time.Minute
		errTimeout = "time is over"
	)

	logger := common.NewTestLogger(t)

	s, addr, err := NewGrpcAppProxy("127.0.0.1:", timeout, logger, network.TCPListener)
	if !assert.NoError(t, err) {
		return
	}
	defer s.Close()

	c, err := NewGrpcLachesisProxy(addr, logger)
	if !assert.NoError(t, err) {
		return
	}
	defer c.Close()

	largeData := make([]byte, largeSize)
	_, err = rand.Read(largeData)
	if !assert.NoError(t, err) {
		return
	}

	t.Run("#1 Send large tx", func(t *testing.T) {
		assert := assert.New(t)

		err = c.SubmitTx(largeData)
		assert.NoError(err)

		select {
		case tx := <-s.SubmitCh():
			assert.Equal(largeData, tx)
		case <-time.After(timeout):
			assert.Fail(errTimeout)
		}
	})

	t.Run("#2 Receive large block", func(t *testing.T) {
		assert := assert.New(t)
		block := poset.Block{
			Body: &poset.BlockBody{
				Transactions: [][]byte{
					largeData,
				},
			},
		}
		hash := largeData[:largeSize/10]

		go func() {
			select {
			case event := <-c.CommitCh():
				assert.EqualValues(block, event.Block)
				event.RespChan <- proto.CommitResponse{
					StateHash: hash,
					Error:     nil,
				}
			case <-time.After(timeout):
				assert.Fail(errTimeout)
			}
		}()

		answer, err := s.CommitBlock(block)
		if assert.NoError(err) {
			assert.Equal(hash, answer)
		}
	})
}
*/
