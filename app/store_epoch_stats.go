package app

import (
	"math"

	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/sfctype"
)

const (
	pendingEpoch = idx.Epoch(math.MaxUint32 - 2)
)

// GetDirtyEpochStats returns EpochStats for current (not sealed) epoch
func (s *Store) GetDirtyEpochStats() *sfctype.EpochStats {
	return s.GetEpochStats(pendingEpoch)
}

// GetEpochStats returns EpochStats for an already sealed epoch
func (s *Store) GetEpochStats(epoch idx.Epoch) *sfctype.EpochStats {
	key := epoch.Bytes()

	// Get data from LRU cache first.
	if s.cache.EpochStats != nil {
		if c, ok := s.cache.EpochStats.Get(epoch); ok {
			if b, ok := c.(*sfctype.EpochStats); ok {
				return b
			}
		}
	}

	w, _ := s.get(s.table.EpochStats, key, &sfctype.EpochStats{}).(*sfctype.EpochStats)

	// Add to LRU cache.
	if w != nil && s.cache.EpochStats != nil {
		s.cache.EpochStats.Add(epoch, w)
	}

	return w
}

// SetDirtyEpochStats set EpochStats for current (not sealed) epoch
func (s *Store) SetDirtyEpochStats(value *sfctype.EpochStats) {
	s.SetEpochStats(pendingEpoch, value)
}

// SetEpochStats set EpochStats for an already sealed epoch
func (s *Store) SetEpochStats(epoch idx.Epoch, value *sfctype.EpochStats) {
	key := epoch.Bytes()

	s.set(s.table.EpochStats, key, value)

	// Add to LRU cache.
	if s.cache.EpochStats != nil {
		s.cache.EpochStats.Add(epoch, value)
	}
}

// GetLastVoting says when last voting was (sealed epoch)
func (s *Store) GetLastVoting() (block idx.Block, start inter.Timestamp) {
	key := []byte("last")

	buf, err := s.table.Voting.Get(key)
	if err != nil {
		s.Log.Crit("Failed to get key-value", "err", err)
	}
	if buf == nil {
		return
	}

	block = idx.BytesToBlock(buf[:8])
	start = inter.BytesToTimestamp(buf[8:])
	return
}

// SetLastVoting saves when last voting was (sealed epoch)
func (s *Store) SetLastVoting(block idx.Block, start inter.Timestamp) {
	key := []byte("last")
	val := append(block.Bytes(), start.Bytes()...)

	err := s.table.Voting.Put(key, val)
	if err != nil {
		s.Log.Crit("Failed to put key-value", "err", err)
	}
}
