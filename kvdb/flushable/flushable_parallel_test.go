package flushable

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

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

	tableMutable1 := table.New(flushableDb, []byte("1"))
	tableImmutable := table.New(flushableDb, []byte("2"))
	tableMutable2 := table.New(flushableDb, []byte("3"))

	// fill data
	for i := uint64(0); i < testPairsNum; i++ {
		_ = tableImmutable.Put(bigendian.Int64ToBytes(i), bigendian.Int64ToBytes(i))
		if i == testPairsNum/2 { // a half of data is flushed, other half isn't
			_ = flushableDb.Flush()
		}
	}

	stop := make(chan struct{})
	stopped := func() bool {
		select {
		case <-stop:
			return true
		default:
			return false
		}
	}

	work := sync.WaitGroup{}
	work.Add(2)
	go func() {
		defer work.Done()
		require := require.New(t)
		for !stopped() {
			// iterate over tableImmutable and check its content
			it := tableImmutable.NewIterator()
			defer it.Release()
			i := uint64(0)
			for ; it.Next(); i++ {
				require.Equal(bigendian.Int64ToBytes(i), it.Key(), i)
				require.Equal(bigendian.Int64ToBytes(i), it.Value(), i)

				require.NoError(it.Error(), i)
			}
			require.Equal(testPairsNum, i)
		}
	}()

	go func() {
		defer work.Done()
		r := rand.New(rand.NewSource(0))
		for !stopped() {
			// try to spoil data in tableImmutable by updating other tables
			_ = tableMutable1.Put(bigendian.Int64ToBytes(r.Uint64()%testPairsNum), bigendian.Int64ToBytes(r.Uint64()))
			_ = tableMutable2.Put(bigendian.Int64ToBytes(r.Uint64() % testPairsNum)[:7], bigendian.Int64ToBytes(r.Uint64()))
			if r.Int63n(100) == 0 {
				_ = flushableDb.Flush() // flush with 1% chance
			}
		}
	}()

	time.Sleep(testDuration)
	close(stop)
	work.Wait()
}
