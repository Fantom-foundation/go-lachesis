package flushable

import (
	"bytes"
	"errors"
	"sync"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"

	"github.com/Fantom-foundation/go-lachesis/kvdb"
)

var (
	errClosed = errors.New("database closed")
)

// Flushable is a ethdb.Database wrapper around any Database.
// On reading, it looks in memory cache first. If not found, it looks in a parent DB.
// On writing, it writes only in cache. To flush the cache into parent DB, call Flush().
type Flushable struct {
	onDrop     func()
	underlying kvdb.KeyValueStore

	modified       *rbt.Tree // modified, comparing to parent, pairs. deleted values are nil
	sizeEstimation *int

	lock *sync.Mutex // we have no guarantees that rbt.Tree works with concurrent reads, so we can't use MutexRW
}

// Wrap underlying db.
// All the writes into the cache won't be written in parent until .Flush() is called.
func Wrap(parent kvdb.KeyValueStore) *Flushable {
	if parent == nil {
		panic("nil parent")
	}

	return WrapWithDrop(parent, func() { parent.Drop() })
}

// WrapWithDrop is the same as Wrap, but defines onDrop callback.
func WrapWithDrop(parent kvdb.KeyValueStore, drop func()) *Flushable {
	if parent == nil {
		panic("nil parent")
	}

	return &Flushable{
		underlying:     parent,
		onDrop:         drop,
		modified:       rbt.NewWithStringComparator(),
		lock:           new(sync.Mutex),
		sizeEstimation: new(int),
	}
}

/*
 * Database interface implementation
 */

// Put puts key-value pair into the cache.
func (w *Flushable) Put(key []byte, value []byte) error {
	w.lock.Lock()
	defer w.lock.Unlock()

	return w.put(key, value)
}

func (w *Flushable) put(key []byte, value []byte) error {
	if value == nil || key == nil {
		return errors.New("Flushable: key or value is nil")
	}
	if w.modified == nil {
		return errClosed
	}

	w.modified.Put(string(key), common.CopyBytes(value))
	*w.sizeEstimation += len(key) + len(value)
	return nil
}

// Has checks if key is in the exists. Looks in cache first, then - in DB.
func (w *Flushable) Has(key []byte) (bool, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.modified == nil {
		return false, errClosed
	}

	val, ok := w.modified.Get(string(key))
	if ok {
		return val != nil, nil
	}

	return w.underlying.Has(key)
}

// Get returns key-value pair by key. Looks in cache first, then - in DB.
func (w *Flushable) Get(key []byte) ([]byte, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.modified == nil {
		return nil, errClosed
	}

	if entry, ok := w.modified.Get(string(key)); ok {
		if entry == nil {
			return nil, nil
		}
		return common.CopyBytes(entry.([]byte)), nil
	}

	return w.underlying.Get(key)
}

// Delete removes key-value pair by key. In parent DB, key won't be deleted until .Flush() is called.
func (w *Flushable) Delete(key []byte) error {
	w.lock.Lock()
	defer w.lock.Unlock()

	return w.delete(key)
}

func (w *Flushable) delete(key []byte) error {
	w.modified.Put(string(key), nil)
	*w.sizeEstimation += len(key) // it should be (len(key) - len(old value)), but we'd need to read old value
	return nil
}

// DropNotFlushed drops all the not flushed keys.
// After this call, the state of parent DB is identical to the state of this DB.
func (w *Flushable) DropNotFlushed() {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.dropNotFlushed()
}

func (w *Flushable) dropNotFlushed() {
	w.modified.Clear()
	*w.sizeEstimation = 0
}

// Close leaves underlying database.
func (w *Flushable) Close() error {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.dropNotFlushed()
	w.modified = nil

	return w.underlying.Close()
}

// Drop whole database.
func (w *Flushable) Drop() {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.modified != nil {
		panic("close db first")
	}

	if w.onDrop != nil {
		w.onDrop()
	}
}

// NotFlushedPairs returns num of not flushed keys, including deleted keys.
func (w *Flushable) NotFlushedPairs() int {
	return w.modified.Size()
}

// NotFlushedSizeEst returns estimation of not flushed data, including deleted keys.
func (w *Flushable) NotFlushedSizeEst() int {
	return *w.sizeEstimation
}

// Flush current cache into parent DB.
func (w *Flushable) Flush() error {
	w.lock.Lock()
	defer w.lock.Unlock()

	return w.flush()
}

func (w *Flushable) flush() error {
	if w.modified == nil {
		return errClosed
	}

	batch := w.underlying.NewBatch()
	for it := w.modified.Iterator(); it.Next(); {
		var err error

		if it.Value() == nil {
			err = batch.Delete([]byte(it.Key().(string)))
		} else {
			err = batch.Put([]byte(it.Key().(string)), it.Value().([]byte))
		}

		if err != nil {
			return err
		}

		if batch.ValueSize() > ethdb.IdealBatchSize {
			err = batch.Write()
			if err != nil {
				return err
			}
			batch.Reset()
		}
	}
	w.modified.Clear()
	*w.sizeEstimation = 0

	return batch.Write()
}

// Stat returns a particular internal stat of the database.
func (w *Flushable) Stat(property string) (string, error) {
	return w.underlying.Stat(property)
}

// Compact flattens the underlying data store for the given key range.
func (w *Flushable) Compact(start []byte, limit []byte) error {
	return w.underlying.Compact(start, limit)
}

/*
 * Iterator
 */

type iterator struct {
	lock *sync.Mutex

	tree *rbt.Tree

	key, val []byte
	prevKey  []byte

	parentIt ethdb.Iterator
	parentOk bool

	treeNode *rbt.Node
	treeOk   bool

	start, prefix []byte

	inited bool
}

// returns the smallest node which is > than specified node
func nextNode(tree *rbt.Tree, node *rbt.Node) (next *rbt.Node, ok bool) {
	origin := node
	if node.Right != nil {
		node = node.Right
		for node.Left != nil {
			node = node.Left
		}
		return node, node != nil
	}
	if node.Parent != nil {
		for node.Parent != nil {
			node = node.Parent
			if tree.Comparator(origin.Key, node.Key) <= 0 {
				return node, node != nil
			}
		}
	}

	return nil, false
}

func castToPair(node *rbt.Node) (key, val []byte) {
	if node == nil {
		return nil, nil
	}
	key = []byte(node.Key.(string))
	if node.Value == nil {
		val = nil // deleted key
	} else {
		val = node.Value.([]byte) // putted value
	}
	return key, val
}

// Next scans key-value pair by key in lexicographic order. Looks in cache first, then - in DB.
func (it *iterator) Next() bool {
	it.lock.Lock()
	defer it.lock.Unlock()

	if it.Error() != nil {
		return false
	}

	isSuitable := func(key, prevKey []byte) (ok, continue_ bool) {
		// check that prefixed. stop current iterating if it isn't
		if it.prefix != nil && !bytes.HasPrefix(key, it.prefix) {
			return false, false
		}
		return prevKey == nil || bytes.Compare(key, prevKey) > 0, true
	}

	if !it.inited {
		it.inited = true
		it.parentOk = it.parentIt.Next()
		if it.start != nil {
			it.treeNode, it.treeOk = it.tree.Ceiling(string(it.start)) // not strict >=
		} else {
			it.treeNode = it.tree.Left() // lowest key
			it.treeOk = it.treeNode != nil
		}
	}

	for it.treeOk || it.parentOk {
		// tree has priority, so check it first
		if it.treeOk {
			treeKey, treeVal := castToPair(it.treeNode)
			for it.treeOk && (!it.parentOk || bytes.Compare(treeKey, it.parentIt.Key()) <= 0) {
				// it's not possible that treeKey isn't bigger than prevKey
				// treeVal may be nil (i.e. deleted). move to next tree's key if it is
				var ok bool
				if treeVal != nil {
					ok, it.treeOk = isSuitable(treeKey, it.prevKey)
				} else {
					it.prevKey = treeKey // next key must be greater than deleted, even if from parent
				}

				if ok {
					it.key, it.val = treeKey, treeVal
					it.prevKey = it.key
				}
				if it.treeOk {
					it.treeNode, it.treeOk = nextNode(it.tree, it.treeNode) // strict >
					treeKey, treeVal = castToPair(it.treeNode)
				}
				if ok {
					return true
				}
			}
		}

		if it.parentOk {
			var ok bool
			ok, it.parentOk = isSuitable(it.parentIt.Key(), it.prevKey)

			if ok {
				it.key = common.CopyBytes(it.parentIt.Key()) // leveldb's iterator may use the same memory
				it.val = common.CopyBytes(it.parentIt.Value())
				it.prevKey = it.key
			}
			if it.parentOk {
				it.parentOk = it.parentIt.Next()
			}
			if ok {
				return true
			}
		}
	}

	return false
}

// Error returns any accumulated error. Exhausting all the key/value pairs
// is not considered to be an error. A memory iterator cannot encounter errors.
func (it *iterator) Error() error {
	return it.parentIt.Error()
}

// Key returns the key of the current key/value pair, or nil if done. The caller
// should not modify the contents of the returned slice, and its contents may
// change on the next call to Next.
func (it *iterator) Key() []byte {
	return it.key
}

// Value returns the value of the current key/value pair, or nil if done. The
// caller should not modify the contents of the returned slice, and its contents
// may change on the next call to Next.
func (it *iterator) Value() []byte {
	return it.val
}

// Release releases associated resources. Release should always succeed and can
// be called multiple times without causing error.
func (it *iterator) Release() {
	it.parentIt.Release()
	*it = iterator{}
}

// NewIterator creates a binary-alphabetical iterator over a subset
// of database content with a particular key prefix, starting at a particular
// initial key (or after, if it does not exist).
func (w *Flushable) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	return &iterator{
		lock:     w.lock,
		tree:     w.modified,
		start:    append(common.CopyBytes(prefix), start...),
		prefix:   prefix,
		parentIt: w.underlying.NewIterator(prefix, start),
	}
}

/*
 * Batch
 */

// NewBatch creates new batch.
func (w *Flushable) NewBatch() ethdb.Batch {
	return &cacheBatch{db: w}
}

type kv struct {
	k, v []byte
}

// cacheBatch is a batch structure.
type cacheBatch struct {
	db     *Flushable
	writes []kv
	size   int
}

// Put adds "add key-value pair" operation into batch.
func (b *cacheBatch) Put(key, value []byte) error {
	b.writes = append(b.writes, kv{common.CopyBytes(key), common.CopyBytes(value)})
	b.size += len(value) + len(key)
	return nil
}

// Delete adds "remove key" operation into batch.
func (b *cacheBatch) Delete(key []byte) error {
	b.writes = append(b.writes, kv{common.CopyBytes(key), nil})
	b.size += len(key)
	return nil
}

// Write writes batch into db. Not atomic.
func (b *cacheBatch) Write() error {
	b.db.lock.Lock()
	defer b.db.lock.Unlock()
	for _, kv := range b.writes {
		var err error

		if kv.v == nil {
			err = b.db.delete(kv.k)
		} else {
			err = b.db.put(kv.k, kv.v)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

// ValueSize returns key-values sizes sum.
func (b *cacheBatch) ValueSize() int {
	return b.size
}

// Reset cleans whole batch.
func (b *cacheBatch) Reset() {
	b.writes = b.writes[:0]
	b.size = 0
}

// Replay replays the batch contents.
func (b *cacheBatch) Replay(w ethdb.KeyValueWriter) error {
	for _, kv := range b.writes {
		if kv.v == nil {
			if err := w.Delete(kv.k); err != nil {
				return err
			}
			continue
		}
		if err := w.Put(kv.k, kv.v); err != nil {
			return err
		}
	}
	return nil
}
