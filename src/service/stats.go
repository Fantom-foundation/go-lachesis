package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Fantom-foundation/go-lachesis/src/peers"
	"github.com/Fantom-foundation/go-lachesis/src/poset"
)

// Poset interface.
type Poset interface {
	GetLastConsensusRound() int64
	GetConsensusTransactionsCount() uint64
	GetUndeterminedEvents() poset.EventHashes
	GetLastCommittedRoundEventsCount() int
}

// Store interface.
type Store interface {
	ConsensusEventsCount() int64
	LastBlockIndex() int64
	Participants() (*peers.Peers, error)
	ConsensusEvents() poset.EventHashes
	GetEventBlock(poset.EventHash) (poset.Event, error)
	LastEventFrom(string) (poset.EventHash, bool, error)
	GetRoundCreated(int64) (poset.RoundCreated, error)
	LastRound() int64
	RoundClothos(int64) poset.EventHashes
	RoundEvents(int64) int
	GetRoot(string) (poset.Root, error)
	GetBlock(int64) (poset.Block, error)
}

// Node interface.
type Node interface {
	ID() uint64
	State() string
	StartTime() time.Time
	SyncRate() float64
	SyncLimit() int64
	HeartbeatTimeout() time.Duration
	GetTransactionPoolCount() int64
	KnownEvents() map[uint64]int64
}

// Stats returns stats from postet and store.
type Stats struct {
	store Store
	poset Poset
	node  Node
}

// NewStats factory to initialize service
func NewStats(s Store, p Poset, n Node) *Stats {
	srv := Stats{
		store: s,
		poset: p,
		node:  n,
	}

	return &srv
}

// Stats returns processing stats.
func (s *Stats) Stats() map[string]string {
	toString := func(i int64) string {
		if i <= 0 {
			return "nil"
		}
		return strconv.FormatInt(i, 10)
	}

	st := s.node.StartTime()
	timeElapsed := time.Since(st)

	consensusEvents := s.store.ConsensusEventsCount()
	consensusEventsPerSecond := float64(consensusEvents) / timeElapsed.Seconds()
	consensusTransactions := s.poset.GetConsensusTransactionsCount()
	transactionsPerSecond := float64(consensusTransactions) / timeElapsed.Seconds()

	lastConsensusRound := s.poset.GetLastConsensusRound()
	var consensusRoundsPerSecond float64
	if lastConsensusRound > poset.RoundNIL {
		consensusRoundsPerSecond = float64(lastConsensusRound+1) / timeElapsed.Seconds()
	}

	peers, _ := s.store.Participants()

	result := map[string]string{
		"last_consensus_round":    toString(lastConsensusRound),
		"time_elapsed":            strconv.FormatFloat(timeElapsed.Seconds(), 'f', 2, 64),
		"heartbeat":               strconv.FormatFloat(s.node.HeartbeatTimeout().Seconds(), 'f', 2, 64),
		"node_current":            strconv.FormatInt(time.Now().Unix(), 10),
		"node_start":              strconv.FormatInt(st.Unix(), 10),
		"last_block_index":        strconv.FormatInt(s.store.LastBlockIndex(), 10),
		"consensus_events":        strconv.FormatInt(consensusEvents, 10),
		"sync_limit":              strconv.FormatInt(s.node.SyncLimit(), 10),
		"consensus_transactions":  strconv.FormatUint(consensusTransactions, 10),
		"undetermined_events":     strconv.Itoa(len(s.poset.GetUndeterminedEvents())),
		"transaction_pool":        strconv.FormatInt(s.node.GetTransactionPoolCount(), 10),
		"num_peers":               strconv.Itoa(peers.Len()),
		"sync_rate":               strconv.FormatFloat(s.node.SyncRate(), 'f', 2, 64),
		"transactions_per_second": strconv.FormatFloat(transactionsPerSecond, 'f', 2, 64),
		"events_per_second":       strconv.FormatFloat(consensusEventsPerSecond, 'f', 2, 64),
		"rounds_per_second":       strconv.FormatFloat(consensusRoundsPerSecond, 'f', 2, 64),
		"round_events":            strconv.Itoa(s.poset.GetLastCommittedRoundEventsCount()),
		"id":                      fmt.Sprint(s.node.ID()),
		"state":                   s.node.State(),
	}

	return result
}

// Participants returns all participants from store.
func (s *Stats) Participants() (*peers.Peers, error) {
	return s.store.Participants()
}

// EventBlock returns a specific event block for the given hash
func (s *Stats) EventBlock(event poset.EventHash) (poset.Event, error) {
	return s.store.GetEventBlock(event)
}

// LastEventFrom returns the last event block for a specific participant
func (s *Stats) LastEventFrom(participant string) (poset.EventHash, bool, error) {
	return s.store.LastEventFrom(participant)
}

// KnownEvents returns all known events
func (s *Stats) KnownEvents() map[uint64]int64 {
	return s.node.KnownEvents()
}

// ConsensusEvents returns all consensus events
func (s *Stats) ConsensusEvents() poset.EventHashes {
	return s.store.ConsensusEvents()
}

// Round returns the created round info for a given index
func (s *Stats) Round(roundIndex int64) (poset.RoundCreated, error) {
	return s.store.GetRoundCreated(roundIndex)
}

// ConsensusTransactionsCount get the count of finalized transactions
func (s *Stats) ConsensusTransactionsCount() uint64 {
	return s.poset.GetConsensusTransactionsCount()
}

// LastRound returns the last round
func (s *Stats) LastRound() int64 {
	return s.store.LastRound()
}

// RoundClothos returns all clotho for a given round index
func (s *Stats) RoundClothos(roundIndex int64) poset.EventHashes {
	return s.store.RoundClothos(roundIndex)
}

// RoundEvents returns all the round events for a given round index
func (s *Stats) RoundEvents(roundIndex int64) int {
	return s.store.RoundEvents(roundIndex)
}

// Root returns the chain root for the frame
func (s *Stats) Root(rootIndex string) (poset.Root, error) {
	return s.store.GetRoot(rootIndex)
}

// Block returns the block for a given index
func (s *Stats) Block(blockIndex int64) (poset.Block, error) {
	return s.store.GetBlock(blockIndex)
}
