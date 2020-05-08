package flushable

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Fantom-foundation/go-lachesis/common/bigendian"
	"github.com/Fantom-foundation/go-lachesis/kvdb/leveldb"
	"github.com/Fantom-foundation/go-lachesis/kvdb/table"
)

func TestFlushableParallel(t *testing.T) {

	testDuration := 2 * time.Second
	testPairsNum := uint64(1000)

	dir, err := ioutil.TempDir("", "test-flushable")
	if err != nil {
		panic(fmt.Sprintf("can't create temporary directory %s: %v", dir, err))
	}
	disk := leveldb.NewProducer(dir)

	// open raw databases
	ldb := disk.OpenDb("1")
	defer ldb.Drop()
	defer ldb.Close()

	flushableDb := Wrap(ldb)

	tableImmutable := table.New(flushableDb, []byte("2"))
	tableMutable1 := table.New(flushableDb, []byte("1"))
	tableMutable2 := table.New(flushableDb, []byte("3"))

	// fill data
	for i := uint64(0); i < testPairsNum; i++ {

		if i == testPairsNum/2 { // a half of data is flushed, other half isn't
			_ = flushableDb.Flush()
		}
		_ = tableImmutable.Put(bigendian.Int64ToBytes(i), bigendian.Int64ToBytes(i))
	}

	stopped := false

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		assertar := assert.New(t)
		for !stopped {
			// iterate over tableImmutable and check its content
			it := tableImmutable.NewIterator()

			defer it.Release()

			i := uint64(0)
			for ; it.Next(); i++ {
				assertar.Equal(bigendian.Int64ToBytes(i), it.Key(), i)
				assertar.Equal(bigendian.Int64ToBytes(i), it.Value(), i)

				assertar.NoError(it.Error(), i)
			}
			assertar.Equal(testPairsNum, i)
		}

		wg.Done()
	}()

	go func() {
		r := rand.New(rand.NewSource(0))
		for !stopped {
			// try to spoil data in tableImmutable by updating other tables
			_ = tableMutable1.Put(bigendian.Int64ToBytes(r.Uint64()%testPairsNum), bigendian.Int64ToBytes(r.Uint64()))
			_ = tableMutable2.Put(bigendian.Int64ToBytes(r.Uint64() % testPairsNum)[:7], bigendian.Int64ToBytes(r.Uint64()))

			if r.Int63n(100) == 0 {
				_ = flushableDb.Flush() // flush with 1% chance
			}
		}

		wg.Done()
	}()

	time.Sleep(testDuration)
	stopped = true
	wg.Wait()

}

func TestFlushableParallelCh(t *testing.T) {

	testDuration := 2 * time.Second
	testPairsNum := uint64(1000)

	dir, err := ioutil.TempDir("", "test-flushable")
	if err != nil {
		panic(fmt.Sprintf("can't create temporary directory %s: %v", dir, err))
	}
	disk := leveldb.NewProducer(dir)

	// open raw databases
	ldb := disk.OpenDb("1")
	defer ldb.Drop()
	defer ldb.Close()

	flushableDb := Wrap(ldb)

	tableImmutable := table.New(flushableDb, []byte("2"))
	tableMutable1 := table.New(flushableDb, []byte("1"))
	tableMutable2 := table.New(flushableDb, []byte("3"))

	// fill data
	for i := uint64(0); i < testPairsNum; i++ {

		if i == testPairsNum/2 { // a half of data is flushed, other half isn't
			_ = flushableDb.Flush()
		}
		_ = tableImmutable.Put(bigendian.Int64ToBytes(i), bigendian.Int64ToBytes(i))
	}

	loop_chan := make(chan bool)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func(ch chan bool) {

		assertar := assert.New(t)

	loop:
		for {
			select {

			default:

				// iterate over tableImmutable and check its content
				it := tableImmutable.NewIterator()

				defer it.Release()

				i := uint64(0)
				for ; it.Next(); i++ {
					assertar.Equal(bigendian.Int64ToBytes(i), it.Key(), i)
					assertar.Equal(bigendian.Int64ToBytes(i), it.Value(), i)

					assertar.NoError(it.Error(), i)
				}
				assertar.Equal(testPairsNum, i)

			case _, ok := <-ch:
				if !ok {
					break loop
				}
			}
		}
		wg.Done()
	}(loop_chan)

	go func(ch chan bool) {
		r := rand.New(rand.NewSource(0))
	loop:
		for {

			select {
			default:
				// try to spoil data in tableImmutable by updating other tables
				_ = tableMutable1.Put(bigendian.Int64ToBytes(r.Uint64()%testPairsNum), bigendian.Int64ToBytes(r.Uint64()))
				_ = tableMutable2.Put(bigendian.Int64ToBytes(r.Uint64() % testPairsNum)[:7], bigendian.Int64ToBytes(r.Uint64()))

				if r.Int63n(10) == 0 {
					_ = flushableDb.Flush() // flush with 1% chance
				}
			case _, ok := <-ch:
				if !ok {
					break loop
				}
			}
		}
		wg.Done()
	}(loop_chan)

	time.Sleep(testDuration)
	close(loop_chan)
	wg.Wait()
}
