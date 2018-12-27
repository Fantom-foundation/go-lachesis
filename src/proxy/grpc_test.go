package proxy

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/poset"
	"github.com/Fantom-foundation/go-lachesis/src/proxy/proto"
	"github.com/Fantom-foundation/go-lachesis/src/utils"
)

func TestGrpcCalls(t *testing.T) {

	const (
		timeout    = 1 * time.Second
		errTimeout = "time is over"
	)

	addr := utils.GetUnusedNetAddr(t)

	logger := common.NewTestLogger(t)

	s, err := NewGrpcAppProxy(addr, timeout, logger)
	assert.NoError(t, err)

	c, err := NewGrpcLachesisProxy(addr, logger)
	assert.NoError(t, err)

	t.Run("#1 Send tx", func(t *testing.T) {
		asserter := assert.New(t)
		gold := []byte("123456")

		err = c.SubmitTx(gold)
		asserter.NoError(err)

		select {
		case tx := <-s.SubmitCh():
			asserter.Equal(gold, tx)
		case <-time.After(timeout):
			asserter.Fail(errTimeout)
		}
	})

	t.Run("#2 Receive block", func(t *testing.T) {
		asserter := assert.New(t)
		block := poset.Block{}
		gold := proto.Response{StateHash: []byte("123456")}

		go func() {
			select {
			case event := <-c.CommitCh():
				asserter.Equal(block, event.Block)
				event.RespChan <- proto.CommitResponse{
					Response: gold,
					Error:    nil,
				}
			case <-time.After(timeout):
				asserter.Fail(errTimeout)
			}
		}()

		answ, err := s.CommitBlock(block)
		if asserter.NoError(err) {
			asserter.Equal(gold, answ)
		}
	})

	t.Run("#3 Receive snapshot query", func(t *testing.T) {
		asserter := assert.New(t)
		index := int64(1)
		gold := []byte("123456")

		go func() {
			select {
			case event := <-c.SnapshotRequestCh():
				asserter.Equal(index, event.BlockIndex)
				event.RespChan <- proto.SnapshotResponse{
					Snapshot: gold,
					Error:    nil,
				}
			case <-time.After(timeout):
				asserter.Fail(errTimeout)
			}
		}()

		answ, err := s.GetSnapshot(index)
		if asserter.NoError(err) {
			asserter.Equal(gold, answ)
		}
	})

	t.Run("#4 Receive restore command", func(t *testing.T) {
		asserter := assert.New(t)
		gold := []byte("123456")

		go func() {
			select {
			case event := <-c.RestoreCh():
				asserter.Equal(gold, event.Snapshot)
				event.RespChan <- proto.RestoreResponse{
					StateHash: gold,
					Error:     nil,
				}
			case <-time.After(timeout):
				asserter.Fail(errTimeout)
			}
		}()

		err := s.Restore(gold)
		asserter.NoError(err)
	})

	err = c.Close()
	assert.NoError(t, err)

	err = s.Close()
	assert.NoError(t, err)

}

func TestGrpcReConnection(t *testing.T) {

	const (
		timeout    = 1 * time.Second
		errTimeout = "time is over"
	)
	addr := utils.GetUnusedNetAddr(t)
	logger := common.NewTestLogger(t)

	c, err := NewGrpcLachesisProxy(addr, logger)
	if assert.NoError(t, err) {
		assert.NotNil(t, c)
	}

	s, err := NewGrpcAppProxy(addr, timeout, logger)
	assert.NoError(t, err)

	checkConnAndStopServer := func(t *testing.T) {
		asserter := assert.New(t)
		gold := []byte("123456")

		err := c.SubmitTx(gold)
		asserter.NoError(err)

		select {
		case tx := <-s.SubmitCh():
			asserter.Equal(gold, tx)
		case <-time.After(timeout):
			asserter.Fail(errTimeout)
		}

		err = s.Close()
		asserter.NoError(err)
	}

	t.Run("#1 Send tx after connection", checkConnAndStopServer)

	s, err = NewGrpcAppProxy(addr, timeout/2, logger)
	assert.NoError(t, err)

	<-time.After(timeout)

	t.Run("#2 Send tx after reconnection", checkConnAndStopServer)

	err = c.Close()
	assert.NoError(t, err)
}

/*
func TestGrpcMaxMsgSize(t *testing.T) {
	const (
		largeSize  = 100 * 1024 * 1024
		timeout    = 3 * time.Minute
		errTimeout = "time is over"

	)
	addr := utils.GetUnusedNetAddr(t);
	logger := common.NewTestLogger(t)

	s, err := NewGrpcAppProxy(addr, timeout, logger)
	assert.NoError(t, err)

	c, err := NewGrpcLachesisProxy(addr, logger)
	assert.NoError(t, err)

	largeData := make([]byte, largeSize)
	_, err = rand.Read(largeData)
	assert.NoError(t, err)

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
			Body: poset.BlockBody{
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

		answ, err := s.CommitBlock(block)
		if assert.NoError(err) {
			assert.Equal(hash, answ)
		}
	})

	err = c.Close()
	assert.NoError(t, err)

	err = s.Close()
	assert.NoError(t, err)
}
*/
