package store

import (
	"fmt"
	"github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/pkg/errors"
)

// Errors.
var (
	ErrNotFound = errors.New("not found")
)

type PendingRound struct {
	Index   uint64
	Decided bool
}

// TODO temp storage
type Store struct {
	events        map[common.Hash][]byte
	rounds        map[uint64][]byte
	lastRound     uint64
	pendingEvents []common.Hash
	pendingRounds []*PendingRound
	lastEvents    map[string]common.Hash
}

func NewStore() *Store {
	return &Store{
		events:     make(map[common.Hash][]byte),
		rounds:     make(map[uint64][]byte),
		lastEvents: make(map[string]common.Hash),
	}
}

func (s *Store) GetEvent(hash common.Hash) []byte {
	return s.events[hash]
}

func (s *Store) SetEvent(hash common.Hash, data []byte) {
	s.events[hash] = data
}

func (s *Store) GetRound(index uint64) ([]byte, error) {
	data, ok := s.rounds[index]
	if !ok {
		return nil, ErrNotFound
	}

	return data, nil
}

func (s *Store) SetRound(index uint64, data []byte) {
	s.rounds[index] = data

	if s.lastRound < index {
		s.lastRound = index
	}
}

// TODO
func (s *Store) LastRound() uint64 {
	return s.lastRound
}

// TODO
func (s *Store) LastEvent(publicKey []byte) common.Hash {
	creator := fmt.Sprintf("0x%X", publicKey)
	return s.lastEvents[creator]
}

// TODO remove it
func (s *Store) SetLastEvent(publicKey []byte, hash common.Hash) {
	creator := fmt.Sprintf("0x%X", publicKey)
	s.lastEvents[creator] = hash
}

func (s *Store) PendingEvents() []common.Hash {
	return s.pendingEvents
}

func (s *Store) AddPendingEvent(hash common.Hash) {
	s.pendingEvents = append(s.pendingEvents, hash)
}

func (s *Store) PendingRounds() []*PendingRound {
	return s.pendingRounds
}

func (s *Store) AddPendingRound(info *PendingRound) {
	s.pendingRounds = append(s.pendingRounds, info)
}
