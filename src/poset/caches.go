package poset

import (
	"fmt"
	"strconv"

	cm "github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/peers"
)

type Key struct {
	x string
	y string
}

func (k Key) ToString() string {
	return fmt.Sprintf("{%s, %s}", k.x, k.y)
}

type ParentRoundInfo struct {
	round                     int
	isRoot                    bool
	rootStronglySeenWitnesses int
}

// ------------------------------------------------------------------------------

type ParticipantEventsCache struct {
	participants *peers.PeerSet
	rim          *cm.RollingIndexMap
}

func NewParticipantEventsCache(size int, participants *peers.PeerSet) *ParticipantEventsCache {
	return &ParticipantEventsCache{
		participants: participants,
		rim:          cm.NewRollingIndexMap("ParticipantEvents", size, participants.IDs()),
	}
}

type PeerSetCache struct {
	rounds   cm.Int64Slice
	peerSets map[int64]*peers.PeerSet
}

func NewPeerSetCache() *PeerSetCache {
	return &PeerSetCache{
		rounds:   cm.Int64Slice{},
		peerSets: make(map[int64]*peers.PeerSet),
	}
}

func (c *PeerSetCache) Set(round int64, peerSet *peers.PeerSet) error {
	if _, ok := c.peerSets[round]; ok {
		return cm.NewStoreErr("PeerSetCache", cm.KeyAlreadyExists, strconv.FormatInt(round, 10))
	}
	c.peerSets[round] = peerSet
	c.rounds = append(c.rounds, round)
	c.rounds.Sort()
	return nil

}

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

func (c *PeerSetCache) GetLast() (*peers.PeerSet, error) {
	if len(c.rounds) == 0 {
		return nil, cm.NewStoreErr("PeerSetCache", cm.NoPeerSet, "")
	}
	return c.peerSets[c.rounds[len(c.rounds)-1]], nil
}

func (pec *ParticipantEventsCache) AddPeer(peer *peers.Peer) error {
	pec.participants = pec.participants.WithNewPeer(peer)
	return pec.rim.AddKey(peer.ID)
}

func (pec *ParticipantEventsCache) participantID(participant string) (uint32, error) {
	peer, ok := pec.participants.ByPubKey[participant]

	if !ok {
		return 0, cm.NewStoreErr("ParticipantEvents", cm.UnknownParticipant, participant)
	}

	return peer.ID, nil
}

// return participant events with index > skip
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

func (pec *ParticipantEventsCache) Set(participant string, hash string, index int64) error {
	id, err := pec.participantID(participant)
	if err != nil {
		return err
	}
	return pec.rim.Set(id, hash, index)
}

// returns [participant id] => lastKnownIndex
func (pec *ParticipantEventsCache) Known() map[uint32]int64 {
	return pec.rim.Known()
}

func (pec *ParticipantEventsCache) Reset() error {
	return pec.rim.Reset()
}

func (pec *ParticipantEventsCache) Import(other *ParticipantEventsCache) {
	pec.rim.Import(other.rim)
}

// ------------------------------------------------------------------------------

type ParticipantBlockSignaturesCache struct {
	participants *peers.Peers
	rim          *cm.RollingIndexMap
}

func NewParticipantBlockSignaturesCache(size int, participants *peers.Peers) *ParticipantBlockSignaturesCache {
	return &ParticipantBlockSignaturesCache{
		participants: participants,
		rim:          cm.NewRollingIndexMap("ParticipantBlockSignatures", size, participants.ToIDSlice()),
	}
}

func (psc *ParticipantBlockSignaturesCache) participantID(participant string) (uint32, error) {
	peer, ok := psc.participants.ByPubKey[participant]

	if !ok {
		return 0, cm.NewStoreErr("ParticipantBlockSignatures", cm.UnknownParticipant, participant)
	}

	return peer.ID, nil
}

// return participant BlockSignatures where index > skip
func (psc *ParticipantBlockSignaturesCache) Get(participant string, skipIndex int64) ([]BlockSignature, error) {
	id, err := psc.participantID(participant)
	if err != nil {
		return []BlockSignature{}, err
	}

	ps, err := psc.rim.Get(id, skipIndex)
	if err != nil {
		return []BlockSignature{}, err
	}

	res := make([]BlockSignature, len(ps))
	for k := 0; k < len(ps); k++ {
		res[k] = ps[k].(BlockSignature)
	}
	return res, nil
}

func (psc *ParticipantBlockSignaturesCache) GetItem(participant string, index int64) (BlockSignature, error) {
	id, err := psc.participantID(participant)
	if err != nil {
		return BlockSignature{}, err
	}

	item, err := psc.rim.GetItem(id, index)
	if err != nil {
		return BlockSignature{}, err
	}
	return item.(BlockSignature), nil
}

func (psc *ParticipantBlockSignaturesCache) GetLast(participant string) (BlockSignature, error) {
	last, err := psc.rim.GetLast(psc.participants.ByPubKey[participant].ID)

	if err != nil {
		return BlockSignature{}, err
	}

	return last.(BlockSignature), nil
}

func (psc *ParticipantBlockSignaturesCache) Set(participant string, sig BlockSignature) error {
	id, err := psc.participantID(participant)
	if err != nil {
		return err
	}

	return psc.rim.Set(id, sig, sig.Index)
}

// returns [participant id] => last BlockSignature Index
func (psc *ParticipantBlockSignaturesCache) Known() map[uint32]int64 {
	return psc.rim.Known()
}

func (psc *ParticipantBlockSignaturesCache) Reset() error {
	return psc.rim.Reset()
}
