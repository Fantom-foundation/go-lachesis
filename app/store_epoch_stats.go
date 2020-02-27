package app

import (
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
)

type voting struct {
	Block idx.Block
	Start inter.Timestamp
}

// GetLastVoting says when last voting was (epoch sealed)
func (s *Store) GetLastVoting() (block idx.Block, start inter.Timestamp) {
	key := []byte("last")

	w, _ := s.get(s.table.Voting, key, &voting{}).(*voting)

	block, start = w.Block, w.Start
	return
}

// SetLastVoting saves when last voting was (epoch sealed)
func (s *Store) SetLastVoting(block idx.Block, start inter.Timestamp) {
	key := []byte("last")

	s.set(s.table.Voting, key, &voting{
		block, start,
	})
}
