// Package memorydb implements the key-value database layer based on memory maps.
package memorydb

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
)

var (
	// errMemorydbClosed is returned if a memory database was already closed at the
	// invocation of a data access operation.
	errMemorydbClosed = errors.New("database closed")
)

// Database is an ephemeral key-value store. Apart from basic data storage
// functionality it also supports batch writes and iterating over the keyspace in
// binary-alphabetical order.
type Database struct {
	db     map[string][]byte
	onDrop func()
	lock   sync.RWMutex

	lag time.Duration
}

// New returns a wrapped map with all the required database interface methods
// implemented.
func New() *Database {
	return &Database{
		db: make(map[string][]byte),
	}
}

// New returns a wrapped map with all the required database interface methods
// implemented and drop callback.
func NewWithDrop(drop func()) *Database {
	return &Database{
		db:     make(map[string][]byte),
		onDrop: drop,
	}
}

// NewWithCap returns a wrapped map pre-allocated to the provided capcity with
// all the required database interface methods implemented.
func NewWithCap(size int) *Database {
	return &Database{
		db: make(map[string][]byte, size),
	}
}

// Close deallocates the internal map and ensures any consecutive data access op
// failes with an error.
func (db *Database) Close() error {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.db = nil
	return nil
}

func (db *Database) Drop() {
	if db.db != nil {
		panic("Close database first!")
	}
	if db.onDrop != nil {
		db.onDrop()
	}
}

// SetDelay on database operations for bench purpose.
func (db *Database) SetDelay(d time.Duration) {
	db.lag = d
}

// Has retrieves if a key is present in the key-value store.
func (db *Database) Has(key []byte) (bool, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if db.db == nil {
		return false, errMemorydbClosed
	}

	if db.lag != 0 {
		time.Sleep(db.lag)
	}

	_, ok := db.db[string(key)]
	return ok, nil
}

// Get retrieves the given key if it's present in the key-value store.
func (db *Database) Get(key []byte) ([]byte, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if db.db == nil {
		return nil, errMemorydbClosed
	}

	if db.lag != 0 {
		time.Sleep(db.lag)
	}

	if entry, ok := db.db[string(key)]; ok {
		return common.CopyBytes(entry), nil
	}
	return nil, nil
}

// Put inserts the given value into the key-value store.
func (db *Database) Put(key []byte, value []byte) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.db == nil {
		return errMemorydbClosed
	}

	if db.lag != 0 {
		time.Sleep(db.lag)
	}

	db.db[string(key)] = common.CopyBytes(value)
	return nil
}

// Delete removes the key from the key-value store.
func (db *Database) Delete(key []byte) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.db == nil {
		return errMemorydbClosed
	}

	if db.lag != 0 {
		time.Sleep(db.lag)
	}

	delete(db.db, string(key))
	return nil
}

// NewBatch creates a write-only key-value store that buffers changes to its host
// database until a final write is called.
func (db *Database) NewBatch() ethdb.Batch {
	return &batch{
		db: db,
	}
}

// NewIterator creates a binary-alphabetical iterator over a subset
// of database content with a particular key prefix, starting at a particular
// initial key (or after, if it does not exist).
func (db *Database) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	db.lock.RLock()
	defer db.lock.RUnlock()

	var (
		pr     = string(prefix)
		st     = string(append(prefix, start...))
		keys   = make([]string, 0, len(db.db))
		values = make([][]byte, 0, len(db.db))
	)
	// Collect the keys from the memory database corresponding to the given prefix
	// and start
	for key := range db.db {
		if !strings.HasPrefix(key, pr) {
			continue
		}
		if key >= st {
			keys = append(keys, key)
		}
	}
	// Sort the items and retrieve the associated values
	sort.Strings(keys)
	for _, key := range keys {
		values = append(values, db.db[key])
	}
	return &iterator{
		keys:   keys,
		values: values,
	}
}

// Stat returns a particular internal stat of the database.
func (db *Database) Stat(property string) (string, error) {
	return "", errors.New("unknown property")
}

// Compact is not supported on a memory database.
func (db *Database) Compact(start []byte, limit []byte) error {
	return errors.New("unsupported operation")
}

// Len returns the number of entries currently present in the memory database.
//
// Note, this method is only used for testing (i.e. not public in general) and
// does not have explicit checks for closed-ness to allow simpler testing code.
func (db *Database) Len() int {
	db.lock.RLock()
	defer db.lock.RUnlock()

	return len(db.db)
}

// keyvalue is a key-value tuple tagged with a deletion field to allow creating
// memory-database write batches.
type keyvalue struct {
	key    []byte
	value  []byte
	delete bool
}

// batch is a write-only memory batch that commits changes to its host
// database when Write is called. A batch cannot be used concurrently.
type batch struct {
	db     *Database
	writes []keyvalue
	size   int
}

// Put inserts the given value into the batch for later committing.
func (b *batch) Put(key, value []byte) error {
	b.writes = append(b.writes, keyvalue{common.CopyBytes(key), common.CopyBytes(value), false})
	b.size += len(value)
	return nil
}

// Delete inserts the a key removal into the batch for later committing.
func (b *batch) Delete(key []byte) error {
	b.writes = append(b.writes, keyvalue{common.CopyBytes(key), nil, true})
	b.size++
	return nil
}

// ValueSize retrieves the amount of data queued up for writing.
func (b *batch) ValueSize() int {
	return b.size
}

// Write flushes any accumulated data to the memory database.
func (b *batch) Write() error {
	b.db.lock.Lock()
	defer b.db.lock.Unlock()

	for _, keyvalue := range b.writes {
		if keyvalue.delete {
			delete(b.db.db, string(keyvalue.key))
			continue
		}
		b.db.db[string(keyvalue.key)] = keyvalue.value
	}
	return nil
}

// Reset resets the batch for reuse.
func (b *batch) Reset() {
	b.writes = b.writes[:0]
	b.size = 0
}

// Replay replays the batch contents.
func (b *batch) Replay(w ethdb.KeyValueWriter) error {
	for _, keyvalue := range b.writes {
		if keyvalue.delete {
			if err := w.Delete(keyvalue.key); err != nil {
				return err
			}
			continue
		}
		if err := w.Put(keyvalue.key, keyvalue.value); err != nil {
			return err
		}
	}
	return nil
}

// iterator can walk over the (potentially partial) keyspace of a memory key
// value store. Internally it is a deep copy of the entire iterated state,
// sorted by keys.
type iterator struct {
	inited bool
	keys   []string
	values [][]byte
}

// Next moves the iterator to the next key/value pair. It returns whether the
// iterator is exhausted.
func (it *iterator) Next() bool {
	// If the iterator was not yet initialized, do it now
	if !it.inited {
		it.inited = true
		return len(it.keys) > 0
	}
	// Iterator already initialize, advance it
	if len(it.keys) > 0 {
		it.keys = it.keys[1:]
		it.values = it.values[1:]
	}
	return len(it.keys) > 0
}

// Error returns any accumulated error. Exhausting all the key/value pairs
// is not considered to be an error. A memory iterator cannot encounter errors.
func (it *iterator) Error() error {
	return nil
}

// Key returns the key of the current key/value pair, or nil if done. The caller
// should not modify the contents of the returned slice, and its contents may
// change on the next call to Next.
func (it *iterator) Key() []byte {
	if len(it.keys) > 0 {
		return []byte(it.keys[0])
	}
	return nil
}

// Value returns the value of the current key/value pair, or nil if done. The
// caller should not modify the contents of the returned slice, and its contents
// may change on the next call to Next.
func (it *iterator) Value() []byte {
	if len(it.values) > 0 {
		return it.values[0]
	}
	return nil
}

// Release releases associated resources. Release should always succeed and can
// be called multiple times without causing error.
func (it *iterator) Release() {
	it.keys, it.values = nil, nil
}
