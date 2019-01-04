package poset

import (
	"fmt"
	"strconv"

	cm "github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/peers"
)

// Key struct
type Key struct {
	x string
	y string
}

// ToString converts key to string
func (k Key) ToString() string {
	return fmt.Sprintf("{%s, %s}", k.x, k.y)
}

// ParentRoundInfo struct
type ParentRoundInfo struct {
	round   int
	isRoot  bool
	Atropos int
}

// NewBaseParentRoundInfo constructor
func NewBaseParentRoundInfo() ParentRoundInfo {
	return ParentRoundInfo{
		round:  -1,
		isRoot: false,
	}
}

// ------------------------------------------------------------------------------

// ParticipantEventsCache struct
type ParticipantEventsCache struct {
	participants *peers.PeerSet
	rim          *cm.RollingIndexMap
}

// NewParticipantEventsCache constructor
func NewParticipantEventsCache(size int, participants *peers.PeerSet) *ParticipantEventsCache {
	return &ParticipantEventsCache{
		participants: participants,
		rim:          cm.NewRollingIndexMap("ParticipantEvents", size, participants.IDs()),
	}
}

// PeerSetCache struct holds map of PeerSet per round
type PeerSetCache struct {
	rounds   cm.Int64Slice
	peerSets map[int64]*peers.PeerSet
}

// NewPeerSetCache PeerSetCache constructor
func NewPeerSetCache() *PeerSetCache {
	return &PeerSetCache{
		rounds:   cm.Int64Slice{},
		peerSets: make(map[int64]*peers.PeerSet),
	}
}

// Set stores a rounds PeerSet
func (c *PeerSetCache) Set(round int64, peerSet *peers.PeerSet) error {
	if _, ok := c.peerSets[round]; ok {
		return cm.NewStoreErr("PeerSetCache", cm.KeyAlreadyExists, strconv.FormatInt(round, 10))
	}
	c.peerSets[round] = peerSet
	c.rounds = append(c.rounds, round)
	c.rounds.Sort()
	return nil
}

// Get retrieves a PeerSet for the given round
func (c *PeerSetCache) Get(round int64) (*peers.PeerSet, error) {
	// check if directly in peerSets
	ps, ok := c.peerSets[round]
	if ok {
		return ps, nil
	}

	// situate round in sorted rounds
	if len(c.rounds) == 0 {
		return nil, cm.NewStoreErr("PeerSetCache", cm.KeyNotFound, strconv.FormatInt(round, 10))
	}

	if len(c.rounds) == 1 {
		if round < c.rounds[0] {
			return nil, cm.NewStoreErr("PeerSetCache", cm.KeyNotFound, strconv.FormatInt(round, 10))
		}
		return c.peerSets[c.rounds[0]], nil
	}

	for i := 0; i < len(c.rounds)-1; i++ {
		if round >= c.rounds[i] && round < c.rounds[i+1] {
			return c.peerSets[c.rounds[i]], nil
		}
	}

	// return last PeerSet
	return c.peerSets[c.rounds[len(c.rounds)-1]], nil
}

// GetLast retrieves the last PeerSet stored
func (c *PeerSetCache) GetLast() (*peers.PeerSet, error) {
	if len(c.rounds) == 0 {
		return nil, cm.NewStoreErr("PeerSetCache", cm.NoPeerSet, "")
	}
	return c.peerSets[c.rounds[len(c.rounds)-1]], nil
}

// AddPeer adds peer to cache and rolling index map, returns error if it failed to add to map
func (pec *ParticipantEventsCache) AddPeer(peer *peers.Peer) error {
	pec.participants = pec.participants.WithNewPeer(peer)
	return pec.rim.AddKey(peer.ID)
}

func (pec *ParticipantEventsCache) participantID(participant string) (uint64, error) {
	peer, ok := pec.participants.ByPubKey[participant]

	if !ok {
		return 0, cm.NewStoreErr("ParticipantEvents", cm.UnknownParticipant, participant)
	}

	return peer.ID, nil
}

// Get return participant events with index > skip
func (pec *ParticipantEventsCache) Get(participant string, skipIndex int64) ([]string, error) {
	id, err := pec.participantID(participant)
	if err != nil {
		return []string{}, err
	}

	pe, err := pec.rim.Get(id, skipIndex)
	if err != nil {
		return []string{}, err
	}

	res := make([]string, len(pe))
	for k := 0; k < len(pe); k++ {
		res[k] = pe[k].(string)
	}
	return res, nil
}

// GetItem get event for participant at index
func (pec *ParticipantEventsCache) GetItem(participant string, index int64) (string, error) {
	id, err := pec.participantID(participant)
	if err != nil {
		return "", err
	}

	item, err := pec.rim.GetItem(id, index)
	if err != nil {
		return "", err
	}
	return item.(string), nil
}

// GetLast get last event for participant
func (pec *ParticipantEventsCache) GetLast(participant string) (string, error) {
	id, err := pec.participantID(participant)
	if err != nil {
		return "", err
	}

	last, err := pec.rim.GetLast(id)
	if err != nil {
		return "", err
	}
	return last.(string), nil
}

// GetLastConsensus get last consensus for participant
func (pec *ParticipantEventsCache) GetLastConsensus(participant string) (string, error) {
	id, err := pec.participantID(participant)
	if err != nil {
		return "", err
	}

	last, err := pec.rim.GetLast(id)
	if err != nil {
		return "", err
	}
	return last.(string), nil
}

// Set the event for the participant
func (pec *ParticipantEventsCache) Set(participant string, hash string, index int64) error {
	id, err := pec.participantID(participant)
	if err != nil {
		return err
	}
	return pec.rim.Set(id, hash, index)
}

// Known returns [participant id] => lastKnownIndex
func (pec *ParticipantEventsCache) Known() map[uint64]int64 {
	return pec.rim.Known()
}

// Reset resets the event cache
func (pec *ParticipantEventsCache) Reset() error {
	return pec.rim.Reset()
}

// Import from another event cache
func (pec *ParticipantEventsCache) Import(other *ParticipantEventsCache) {
	pec.rim.Import(other.rim)
}
