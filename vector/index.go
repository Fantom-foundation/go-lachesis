package vector

import (
	"github.com/quan8/go-ethereum/common"
	"github.com/hashicorp/golang-lru"

	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/pos"
	"github.com/Fantom-foundation/go-lachesis/kvdb"
	"github.com/Fantom-foundation/go-lachesis/logger"
)

const (
	forklessCauseCacheSize = 5000
)

// Index is a data to detect forkless-cause condition, calculate median timestamp, detect forks.
type Index struct {
	validators    pos.Validators
	validatorIdxs map[common.Address]idx.Validator

	getEvent func(hash.Event) *inter.EventHeaderData

	forklessCauseCache *lru.Cache

	logger.Instance
}

// NewIndex creates Index instance.
func NewIndex(validators pos.Validators, db kvdb.KeyValueStore, getEvent func(hash.Event) *inter.EventHeaderData) *Index {
	cache, _ := lru.New(forklessCauseCacheSize)

	vi := &Index{
		Instance:           logger.MakeInstance(),
		forklessCauseCache: cache,
	}
	vi.Reset(validators, db, getEvent)

	return vi
}

// Reset resets buffers.
func (vi *Index) Reset(validators pos.Validators, db kvdb.KeyValueStore, getEvent func(hash.Event) *inter.EventHeaderData) {
	// we use wrapper to be able to drop failed events by dropping cache
	vi.getEvent = getEvent
	vi.validators = validators.Copy()
	vi.validatorIdxs = validators.Idxs()
}

// Add calculates vector clocks for the event and saves into DB.
func (vi *Index) Add(e *inter.EventHeaderData) {
	// TODO
}

