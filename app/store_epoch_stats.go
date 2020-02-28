package app

import (
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
)

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
