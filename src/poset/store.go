// +build !debug

package poset

import "github.com/Fantom-foundation/go-lachesis/src/peers"

// Store provides an interface for persistent and non-persistent stores
// to store key lachesis consensus information on a node.
type Store interface {
	CacheSize() int
	GetLastPeerSet() (*peers.PeerSet, error)
	GetPeerSet(int64) (*peers.PeerSet, error)
	SetPeerSet(int64, *peers.PeerSet) error
	RepertoireByPubKey() map[string]*peers.Peer
	RepertoireByID() map[uint64]*peers.Peer
	RootsBySelfParent() map[EventHash]*Root
	GetEventBlock(EventHash) (*Event, error)
	SetEvent(*Event) error
	ParticipantEvents(string, int64) (EventHashes, error)
	ParticipantEvent(string, int64) (EventHash, error)
	LastEventFrom(string) (EventHash, bool, error)
	LastConsensusEventFrom(string) (EventHash, bool, error)
	KnownEvents() map[uint64]int64
	ConsensusEvents() EventHashes
	ConsensusEventsCount() int64
	AddConsensusEvent(*Event) error
	GetRoundCreated(int64) (*RoundCreated, error)
	SetRoundCreated(int64, *RoundCreated) error
	GetRoundReceived(int64) (*RoundReceived, error)
	SetRoundReceived(int64, *RoundReceived) error
	LastRound() int64
	RoundClothos(int64) EventHashes
	RoundEvents(int64) int
	GetRoot(string) (*Root, error)
	GetBlock(int64) (*Block, error)
	SetBlock(*Block) error
	LastBlockIndex() int64
	GetFrame(int64) (*Frame, error)
	SetFrame(*Frame) error
	Reset(*Frame) error
	Close() error
	NeedBoostrap() bool // Was the store loaded from existing db
	StorePath() string
}
