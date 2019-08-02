package kvdb

// IdealBatchSize was determined empirically.
// Code using batches should try to add this much data to the batch.
const IdealBatchSize = 100 * 1024

// Putter wraps the database write operation supported by both batches and regular databases.
type Putter interface {
	Put(key []byte, value []byte) error
}

// Deleter wraps the database delete operation supported by both batches and regular databases.
type Deleter interface {
	Delete(key []byte) error
}

// Database wraps all database operations. All methods are safe for concurrent use.
type Database interface {
	NewTable(prefix []byte) Database
	Putter
	Deleter
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	ForEach(prefix []byte, do func(key, val []byte) bool) error
	Close()
	NewBatch() Batch
}

type FlushableDatabase interface {
	NewTableFlushable(prefix []byte) FlushableDatabase
	Database
	NotFlushedPairs() int
	Flush() error
	ClearNotFlushed()
}

// Batch is a write-only database that commits changes to its host database
// when Write is called. Batch cannot be used concurrently.
type Batch interface {
	Putter
	Deleter
	ValueSize() int // amount of data in the batch
	Write() error
	// Reset resets the batch for reuse
	Reset()
}
