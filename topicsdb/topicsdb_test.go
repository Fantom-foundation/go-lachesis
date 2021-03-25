package topicsdb

import (
	"context"
	"math/rand"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/kvdb/memorydb"
	"github.com/Fantom-foundation/go-lachesis/logger"
)

// FindInBlocksAsync returns all log records of block range by pattern. 1st pattern element is an address.
// Fetches log's body async.
func (tt *Index) FindInBlocksAsync(ctx context.Context, from, to idx.Block, pattern [][]common.Hash) (logs []*types.Log, err error) {
	if from > to {
		return
	}

	err = checkPattern(pattern)
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	ready := make(chan *logrec)
	defer close(ready)

	go func() {
		failed := false
		for rec := range ready {
			wg.Done()
			if failed {
				continue
			}
			if rec.err != nil {
				err = rec.err
				failed = true
				continue
			}
			logs = append(logs, rec.result)
		}
	}()

	onMatched := func(rec *logrec) (gonext bool, err error) {
		if rec.ID.BlockNumber() > uint64(to) {
			return
		}

		wg.Add(1)
		go func() {
			rec.fetch(tt.table.Logrec)
			ready <- rec
		}()

		gonext = true
		return
	}

	err = tt.searchLazy(ctx, pattern, uintToBytes(uint64(from)), onMatched)
	wg.Wait()

	return
}

func TestIndexSearchMultyVariants(t *testing.T) {
	logger.SetTestMode(t)
	var (
		hash1 = common.BytesToHash([]byte("topic1"))
		hash2 = common.BytesToHash([]byte("topic2"))
		hash3 = common.BytesToHash([]byte("topic3"))
		hash4 = common.BytesToHash([]byte("topic4"))
		addr1 = randAddress()
		addr2 = randAddress()
		addr3 = randAddress()
		addr4 = randAddress()
	)
	testdata := []*types.Log{{
		BlockNumber: 1,
		Address:     addr1,
		Topics:      []common.Hash{hash1, hash1, hash1},
	}, {
		BlockNumber: 3,
		Address:     addr2,
		Topics:      []common.Hash{hash2, hash2, hash2},
	}, {
		BlockNumber: 998,
		Address:     addr3,
		Topics:      []common.Hash{hash3, hash3, hash3},
	}, {
		BlockNumber: 999,
		Address:     addr4,
		Topics:      []common.Hash{hash4, hash4, hash4},
	},
	}

	index := New(memorydb.New())

	for _, l := range testdata {
		err := index.Push(l)
		require.NoError(t, err)
	}

	// require.ElementsMatchf(testdata, got, "") doesn't work properly here,
	// so use check()
	check := func(require *require.Assertions, got []*types.Log) {
		count := 0
		for _, a := range got {
			for _, b := range testdata {
				if b.Address == a.Address {
					require.ElementsMatch(a.Topics, b.Topics)
					count++
					break
				}
			}
		}
	}

	for dsc, method := range map[string]func(context.Context, idx.Block, idx.Block, [][]common.Hash) ([]*types.Log, error){
		"sync":  index.FindInBlocks,
		"async": index.FindInBlocksAsync,
	} {
		t.Run(dsc, func(t *testing.T) {

			t.Run("With no addresses", func(t *testing.T) {
				require := require.New(t)
				got, err := method(nil, 0, 1000, [][]common.Hash{
					{},
					{hash1, hash2, hash3, hash4},
					{},
					{hash1, hash2, hash3, hash4},
				})
				require.NoError(err)
				require.Equal(4, len(got))
				check(require, got)
			})

			t.Run("With addresses", func(t *testing.T) {
				require := require.New(t)
				got, err := method(nil, 0, 1000, [][]common.Hash{
					{addr1.Hash(), addr2.Hash(), addr3.Hash(), addr4.Hash()},
					{hash1, hash2, hash3, hash4},
					{},
					{hash1, hash2, hash3, hash4},
				})
				require.NoError(err)
				require.Equal(4, len(got))
				check(require, got)
			})

			t.Run("With block range", func(t *testing.T) {
				require := require.New(t)
				got, err := method(nil, 2, 998, [][]common.Hash{
					{addr1.Hash(), addr2.Hash(), addr3.Hash(), addr4.Hash()},
					{hash1, hash2, hash3, hash4},
					{},
					{hash1, hash2, hash3, hash4},
				})
				require.NoError(err)
				require.Equal(2, len(got))
				check(require, got)
			})

		})
	}
}

func TestIndexSearchSingleVariant(t *testing.T) {
	logger.SetTestMode(t)

	topics, recs, topics4rec := genTestData(100)

	index := New(memorydb.New())

	for _, rec := range recs {
		err := index.Push(rec)
		require.NoError(t, err)
	}

	for dsc, method := range map[string]func(context.Context, idx.Block, idx.Block, [][]common.Hash) ([]*types.Log, error){
		"sync":  index.FindInBlocks,
		"async": index.FindInBlocksAsync,
	} {
		t.Run(dsc, func(t *testing.T) {
			require := require.New(t)

			for i := 0; i < len(topics); i++ {
				from, to := topics4rec(i)
				tt := topics[from : to-1]

				qq := make([][]common.Hash, len(tt)+1)
				for pos, t := range tt {
					qq[pos+1] = []common.Hash{t}
				}

				got, err := method(nil, 0, 1000, qq)
				require.NoError(err)

				var expect []*types.Log
				for j, rec := range recs {
					if f, t := topics4rec(j); f != from || t != to {
						continue
					}
					expect = append(expect, rec)
				}

				require.ElementsMatchf(expect, got, "step %d", i)
			}

		})
	}
}

func TestIndexSearchSimple(t *testing.T) {
	logger.SetTestMode(t)

	var (
		hash1 = common.BytesToHash([]byte("topic1"))
		hash2 = common.BytesToHash([]byte("topic2"))
		hash3 = common.BytesToHash([]byte("topic3"))
		hash4 = common.BytesToHash([]byte("topic4"))
		addr  = randAddress()
	)
	testdata := []*types.Log{{
		BlockNumber: 1,
		Address:     addr,
		Topics:      []common.Hash{hash1},
	}, {
		BlockNumber: 2,
		Address:     addr,
		Topics:      []common.Hash{hash2},
	}, {
		BlockNumber: 998,
		Address:     addr,
		Topics:      []common.Hash{hash3},
	}, {
		BlockNumber: 999,
		Address:     addr,
		Topics:      []common.Hash{hash4},
	},
	}

	index := New(memorydb.New())

	for _, l := range testdata {
		err := index.Push(l)
		require.NoError(t, err)
	}

	var (
		got []*types.Log
		err error
	)

	for dsc, method := range map[string]func(context.Context, idx.Block, idx.Block, [][]common.Hash) ([]*types.Log, error){
		"sync":  index.FindInBlocks,
		"async": index.FindInBlocksAsync,
	} {
		t.Run(dsc, func(t *testing.T) {
			require := require.New(t)

			got, err = method(nil, 0, 0xffffffff, [][]common.Hash{
				{addr.Hash()},
				{hash1},
			})
			require.NoError(err)
			require.Equal(1, len(got))

			got, err = method(nil, 0, 0xffffffff, [][]common.Hash{
				{addr.Hash()},
				{hash2},
			})
			require.NoError(err)
			require.Equal(1, len(got))

			got, err = method(nil, 0, 0xffffffff, [][]common.Hash{
				{addr.Hash()},
				{hash3},
			})
			require.NoError(err)
			require.Equal(1, len(got))
		})
	}

}

func genTestData(count int) (
	topics []common.Hash,
	recs []*types.Log,
	topics4rec func(rec int) (from, to int),
) {
	const (
		period = 5
	)

	topics = make([]common.Hash, period)
	for i := range topics {
		topics[i] = hash.FakeHash(int64(i))
	}

	topics4rec = func(rec int) (from, to int) {
		from = rec % (period - 3)
		to = from + 3
		return
	}

	recs = make([]*types.Log, count)
	for i := range recs {
		from, to := topics4rec(i)
		r := &types.Log{
			BlockNumber: uint64(i / period),
			BlockHash:   hash.FakeHash(int64(i / period)),
			TxHash:      hash.FakeHash(int64(i % period)),
			Index:       uint(i % period),
			Address:     randAddress(),
			Topics:      topics[from:to],
			Data:        make([]byte, i),
		}
		_, _ = rand.Read(r.Data)
		recs[i] = r
	}

	return
}

func randAddress() (addr common.Address) {
	n, err := rand.Read(addr[:])
	if err != nil {
		panic(err)
	}
	if n != common.AddressLength {
		panic("address is not filled")
	}
	return
}
