package vector

import (
	"github.com/Fantom-foundation/go-ethereum/common"
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
	highestBeforeSeqCacheSize = 1000
	highestBeforeTimeCacheSize = 1000
	lowestAfterSeqCacheSize = 1000
	eventBranchCacheSize = 1000
	branchesInfoCacheSize = 1000
	dfsSubgraphVisitedCacheSize = 1000
)

// IndexCacheConfig - config for cache sizes of Index
type IndexCacheConfig struct {
	ForklessCause     int `json:"forklessCause"`
	HighestBeforeSeq  int `json:"highestBeforeSeq"`
	HighestBeforeTime int `json:"highestBeforeTime"`
	LowestAfterSeq    int `json:"lowestAfterSeq"`
	EventBranch       int `json:"eventBranch"`
	BranchesInfo      int `json:"branchesInfo"`
	DfsSubgraphVisited int `json:"dfsSubgraphVisited"`
}

// IndexConfig - Index config (cache sizes)
type IndexConfig struct {
	Caches IndexCacheConfig `json:"cacheSizes"`
}

// Index is a data to detect forkless-cause condition, calculate median timestamp, detect forks.
type Index struct {
	validators    pos.Validators
	validatorIdxs map[common.Address]idx.Validator

	getEvent func(hash.Event) *inter.EventHeaderData

	forklessCauseCache *lru.Cache

	cfg IndexConfig

	logger.Instance
}

// DefaultIndexConfig return default index config for tests
func DefaultIndexConfig() IndexConfig {
	return IndexConfig{
		Caches: IndexCacheConfig{
			ForklessCause:      forklessCauseCacheSize,
			HighestBeforeSeq:   highestBeforeSeqCacheSize,
			HighestBeforeTime:  highestBeforeTimeCacheSize,
			LowestAfterSeq:     lowestAfterSeqCacheSize,
			EventBranch:        eventBranchCacheSize,
			BranchesInfo:       branchesInfoCacheSize,
			DfsSubgraphVisited: dfsSubgraphVisitedCacheSize,
		},
	}
}

// NewIndex creates Index instance.
func NewIndex(config IndexConfig, validators pos.Validators, db kvdb.KeyValueStore, getEvent func(hash.Event) *inter.EventHeaderData) *Index {
	cache, _ := lru.New(config.Caches.ForklessCause)

	vi := &Index{
		Instance:           logger.MakeInstance(),
		forklessCauseCache: cache,
		cfg: config,
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

func (vi *Index) cleanCaches() {
}

// Add calculates vector clocks for the event and saves into DB.
func (vi *Index) Add(e *inter.EventHeaderData) {
	// TODO
}
