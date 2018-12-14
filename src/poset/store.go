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
	RepertoireByID() map[uint32]*peers.Peer
	RootsBySelfParent() map[string]*Root
	GetEvent(string) (*Event, error)
	SetEvent(*Event) error
	ParticipantEvents(string, int64) ([]string, error)
	ParticipantEvent(string, int64) (string, error)
	LastEventFrom(string) (string, bool, error)
	LastConsensusEventFrom(string) (string, bool, error)
	KnownEvents() map[uint32]int64
	ConsensusEvents() []string
	ConsensusEventsCount() int64
	AddConsensusEvent(*Event) error
	GetRoundCreated(int64) (*RoundCreated, error)
	SetRoundCreated(int64, *RoundCreated) error
	GetRoundReceived(int64) (*RoundReceived, error)
	SetRoundReceived(int64, *RoundReceived) error
	LastRound() int64
	RoundWitnesses(int64) []string
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
