package app

import (
	"bytes"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/hashicorp/golang-lru"

	"github.com/Fantom-foundation/go-lachesis/common/bigendian"
	"github.com/Fantom-foundation/go-lachesis/kvdb"
	"github.com/Fantom-foundation/go-lachesis/kvdb/flushable"
	"github.com/Fantom-foundation/go-lachesis/kvdb/memorydb"
	"github.com/Fantom-foundation/go-lachesis/kvdb/nokeyiserr"
	"github.com/Fantom-foundation/go-lachesis/kvdb/table"
	"github.com/Fantom-foundation/go-lachesis/logger"
	"github.com/Fantom-foundation/go-lachesis/topicsdb"
)

// Store is an app persistent storage working over physical key-value database.
type Store struct {
	dbs *flushable.SyncedPool

	mainDb kvdb.KeyValueStore
	table  struct {

		// score economy tables
		ActiveValidationScore  kvdb.KeyValueStore `table:"V"`
		DirtyValidationScore   kvdb.KeyValueStore `table:"v"`
		ActiveOriginationScore kvdb.KeyValueStore `table:"O"`
		DirtyOriginationScore  kvdb.KeyValueStore `table:"o"`
		BlockDowntime          kvdb.KeyValueStore `table:"m"`

		// PoI economy tables
		StakerPOIScore      kvdb.KeyValueStore `table:"s"`
		AddressPOIScore     kvdb.KeyValueStore `table:"a"`
		AddressFee          kvdb.KeyValueStore `table:"g"`
		StakerDelegatorsFee kvdb.KeyValueStore `table:"d"`
		AddressLastTxTime   kvdb.KeyValueStore `table:"X"`
		TotalPoiFee         kvdb.KeyValueStore `table:"U"`

		// SFC-related economy tables
		Validators   kvdb.KeyValueStore `table:"1"`
		Stakers      kvdb.KeyValueStore `table:"2"`
		Delegators   kvdb.KeyValueStore `table:"3"`
		SfcConstants kvdb.KeyValueStore `table:"4"`
		TotalSupply  kvdb.KeyValueStore `table:"5"`

		// gas power economy tables
		GasPowerRefund kvdb.KeyValueStore `table:"R"`

		// API-only tables
		DelegatorOldRewards        kvdb.KeyValueStore `table:"6"`
		StakerOldRewards           kvdb.KeyValueStore `table:"7"`
		StakerDelegatorsOldRewards kvdb.KeyValueStore `table:"8"`

		Evm      ethdb.Database
		EvmState state.Database
		EvmLogs  *topicsdb.Index
	}

	cache struct {
		Stakers       *lru.Cache `cache:"-"` // store by pointer
		Delegators    *lru.Cache `cache:"-"` // store by pointer
		BlockDowntime *lru.Cache `cache:"-"` // store by pointer
	}

	mutex struct {
		Inc sync.Mutex
	}

	logger.Instance
}

// NewMemStore creates store over memory map.
func NewMemStore() *Store {
	mems := memorydb.NewProducer("")
	dbs := flushable.NewSyncedPool(mems)

	return NewStore(dbs)
}

// NewStore creates store over key-value db.
func NewStore(dbs *flushable.SyncedPool) *Store {
	s := &Store{
		dbs:      dbs,
		mainDb:   dbs.GetDb("app-main"),
		Instance: logger.MakeInstance(),
	}

	table.MigrateTables(&s.table, s.mainDb)

	evmTable := nokeyiserr.Wrap(table.New(s.mainDb, []byte("M"))) // ETH expects that "not found" is an error
	s.table.Evm = rawdb.NewDatabase(evmTable)
	s.table.EvmState = state.NewDatabaseWithCache(s.table.Evm, 16)
	s.table.EvmLogs = topicsdb.New(table.New(s.mainDb, []byte("L")))

	s.initCache()

	return s
}

func (s *Store) initCache() {
	s.cache.Stakers = s.makeCache(100)    // TODO: s.cfg.StakersCacheSize
	s.cache.Delegators = s.makeCache(100) // TODO: s.cfg.DelegatorsCacheSize
	s.cache.BlockDowntime = s.makeCache(256)
}

// Close leaves underlying database.
func (s *Store) Close() {
	setnil := func() interface{} {
		return nil
	}

	table.MigrateTables(&s.table, nil)
	table.MigrateCaches(&s.cache, setnil)

	s.mainDb.Close()
}

// Commit changes.
func (s *Store) Commit(flushID []byte, immediately bool) error {
	if flushID == nil {
		// if flushId not specified, use current time
		buf := bytes.NewBuffer(nil)
		buf.Write([]byte{0xbe, 0xee})                                    // 0xbeee eyecatcher that flushed time
		buf.Write(bigendian.Int64ToBytes(uint64(time.Now().UnixNano()))) // current UnixNano time
		flushID = buf.Bytes()
	}

	if immediately || s.dbs.IsFlushNeeded() {
		// Flush trie on the DB
		err := s.table.EvmState.TrieDB().Cap(0)
		if err != nil {
			s.Log.Error("Failed to flush trie DB into main DB", "err", err)
		}
		// Flush the DBs
		return s.dbs.Flush(flushID)
	}

	return nil
}

// StateDB returns state database.
func (s *Store) StateDB(from common.Hash) *state.StateDB {
	db, err := state.New(common.Hash(from), s.table.EvmState)
	if err != nil {
		s.Log.Crit("Failed to open state", "err", err)
	}
	return db
}

/*
 * Utils:
 */

// set RLP value
func (s *Store) set(table kvdb.KeyValueStore, key []byte, val interface{}) {
	buf, err := rlp.EncodeToBytes(val)
	if err != nil {
		s.Log.Crit("Failed to encode rlp", "err", err)
	}

	if err := table.Put(key, buf); err != nil {
		s.Log.Crit("Failed to put key-value", "err", err)
	}
}

// get RLP value
func (s *Store) get(table kvdb.KeyValueStore, key []byte, to interface{}) interface{} {
	buf, err := table.Get(key)
	if err != nil {
		s.Log.Crit("Failed to get key-value", "err", err)
	}
	if buf == nil {
		return nil
	}

	err = rlp.DecodeBytes(buf, to)
	if err != nil {
		s.Log.Crit("Failed to decode rlp", "err", err, "size", len(buf))
	}
	return to
}

func (s *Store) has(table kvdb.KeyValueStore, key []byte) bool {
	res, err := table.Has(key)
	if err != nil {
		s.Log.Crit("Failed to get key", "err", err)
	}
	return res
}

func (s *Store) dropTable(it ethdb.Iterator, t kvdb.KeyValueStore) {
	keys := make([][]byte, 0, 500) // don't write during iteration

	for it.Next() {
		keys = append(keys, it.Key())
	}

	for i := range keys {
		err := t.Delete(keys[i])
		if err != nil {
			s.Log.Crit("Failed to erase key-value", "err", err)
		}
	}
}

func (s *Store) makeCache(size int) *lru.Cache {
	if size <= 0 {
		return nil
	}

	cache, err := lru.New(size)
	if err != nil {
		s.Log.Crit("Error create LRU cache", "err", err)
		return nil
	}
	return cache
}
