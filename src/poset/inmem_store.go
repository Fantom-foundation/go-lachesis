package poset

import (
	"strconv"
	"sync"

	cm "github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/peers"
)

// InmemStore struct
type InmemStore struct {
	cacheSize              int
	participants           *peers.Peers
	eventCache             *cm.Cache        // hash => Event
	roundCreatedCache      *cm.Cache        // round number => RoundCreated
	roundReceivedCache     *cm.Cache        // round received number => RoundReceived
	blockCache             *cm.Cache        // index => Block
	frameCache             *cm.Cache        // round received => Frame
	consensusCache         *cm.RollingIndex // consensus index => hash
	totConsensusEvents     int64
	repertoireByPubKey     map[string]*peers.Peer
	repertoireByID         map[int64]*peers.Peer
	participantEventsCache *ParticipantEventsCache // pubkey => Events
	rootsByParticipant     map[string]Root         // [participant] => Root
	rootsBySelfParent      map[string]Root         // [Root.SelfParent.Hash] => Root
	lastRound              int64
	lastConsensusEvents    map[string]string // [participant] => hex() of last consensus event
	lastBlock              int64

	lastRoundLocker          sync.RWMutex
	lastBlockLocker          sync.RWMutex
	totConsensusEventsLocker sync.RWMutex
}

// NewInmemStore constructor
func NewInmemStore(participants *peers.Peers, cacheSize int) *InmemStore {
	rootsByParticipant := make(map[string]Root)

	for pk, pid := range participants.ByPubKey {
		root := NewBaseRoot(pid.ID)
		rootsByParticipant[pk] = root
	}

	store := &InmemStore{
		cacheSize:              cacheSize,
		participants:           participants,
		eventCache:             cm.NewCache(cacheSize),
		roundCreatedCache:      cm.NewCache(cacheSize),
		roundReceivedCache:     cm.NewCache(cacheSize),
		blockCache:             cm.NewCache(cacheSize),
		frameCache:             cm.NewCache(cacheSize),
		consensusCache:         cm.NewRollingIndex("ConsensusCache", cacheSize),
		repertoireByPubKey:     make(map[string]*peers.Peer),
		repertoireByID:         make(map[int64]*peers.Peer),
		participantEventsCache: NewParticipantEventsCache(cacheSize, participants),
		rootsByParticipant:     rootsByParticipant,
		lastRound:              -1,
		lastBlock:              -1,
		lastConsensusEvents:    map[string]string{},
	}

	participants.OnNewPeer(func(peer *peers.Peer) {
		root := NewBaseRoot(peer.ID)
		store.rootsByParticipant[peer.PubKeyHex] = root
		store.repertoireByPubKey[peer.PubKeyHex] = peer
		store.repertoireByID[peer.ID] = peer
		store.rootsBySelfParent = nil
		store.RootsBySelfParent()
		old := store.participantEventsCache
		store.participantEventsCache = NewParticipantEventsCache(cacheSize, participants)
		store.participantEventsCache.Import(old)
	})

	store.setPeers(0, participants)
	return store
}

// CacheSize size of cache
func (s *InmemStore) CacheSize() int {
	return s.cacheSize
}

// Participants returns participants
func (s *InmemStore) Participants() (*peers.Peers, error) {
	return s.participants, nil
}

func (s *InmemStore) setPeers(round int64, participants *peers.Peers) {
	// Extend PartipantEventsCache and Roots with new peers
	for _, peer := range participants.ByID {
		s.repertoireByPubKey[peer.PubKeyHex] = peer
		s.repertoireByID[peer.ID] = peer
	}
}

// RepertoireByPubKey retrieves cached PubKey map of peers
func (s *InmemStore) RepertoireByPubKey() map[string]*peers.Peer {
	return s.repertoireByPubKey
}

// RepertoireByID retrieve cached ID map of peers
func (s *InmemStore) RepertoireByID() map[int64]*peers.Peer {
	return s.repertoireByID
}

// RootsBySelfParent TODO
func (s *InmemStore) RootsBySelfParent() (map[string]Root, error) {
	if s.rootsBySelfParent == nil {
		s.rootsBySelfParent = make(map[string]Root)
		for _, root := range s.rootsByParticipant {
			s.rootsBySelfParent[root.SelfParent.Hash] = root
		}
	}
	return s.rootsBySelfParent, nil
}

// GetEventBlock returns the event block for a key
func (s *InmemStore) GetEventBlock(key string) (Event, error) {
	res, ok := s.eventCache.Get(key)
	if !ok {
		return Event{}, cm.NewStoreErr("EventCache", cm.KeyNotFound, key)
	}

	return res.(Event), nil
}

// SetEvent set event for event block
func (s *InmemStore) SetEvent(event Event) error {
	key := event.Hex()
	_, err := s.GetEventBlock(key)
	if err != nil && !cm.Is(err, cm.KeyNotFound) {
		return err
	}
	if cm.Is(err, cm.KeyNotFound) {
		if err := s.addParticpantEvent(event.GetCreator(), key, event.Index()); err != nil {
			return err
		}
	}

	// fmt.Println("Adding event to cache", event.Hex())
	s.eventCache.Add(key, event)

	return nil
}

func (s *InmemStore) addParticpantEvent(participant string, hash string, index int64) error {
	return s.participantEventsCache.Set(participant, hash, index)
}

// ParticipantEvents events for the participant
func (s *InmemStore) ParticipantEvents(participant string, skip int64) ([]string, error) {
	return s.participantEventsCache.Get(participant, skip)
}

// ParticipantEvent specific event
func (s *InmemStore) ParticipantEvent(participant string, index int64) (string, error) {
	ev, err := s.participantEventsCache.GetItem(participant, index)
	if err != nil {
		root, ok := s.rootsByParticipant[participant]
		if !ok {
			return "", cm.NewStoreErr("InmemStore.Roots", cm.NoRoot, participant)
		}
		if root.SelfParent.Index == index {
			ev = root.SelfParent.Hash
			err = nil
		}
	}
	return ev, err
}

// LastEventFrom participant
func (s *InmemStore) LastEventFrom(participant string) (last string, isRoot bool, err error) {
	// try to get the last event from this participant
	last, err = s.participantEventsCache.GetLast(participant)

	// if there is none, grab the root
	if err != nil && cm.Is(err, cm.Empty) {
		root, ok := s.rootsByParticipant[participant]
		if ok {
			last = root.SelfParent.Hash
			isRoot = true
			err = nil
		} else {
			err = cm.NewStoreErr("InmemStore.Roots", cm.NoRoot, participant)
		}
	}
	return
}

// LastConsensusEventFrom participant
func (s *InmemStore) LastConsensusEventFrom(participant string) (last string, isRoot bool, err error) {
	// try to get the last consensus event from this participant
	last, ok := s.lastConsensusEvents[participant]
	// if there is none, grab the root
	if !ok {
		root, ok := s.rootsByParticipant[participant]
		if ok {
			last = root.SelfParent.Hash
			isRoot = true
		} else {
			err = cm.NewStoreErr("InmemStore.Roots", cm.NoRoot, participant)
		}
	}
	return
}

// KnownEvents returns all known events
func (s *InmemStore) KnownEvents() map[int64]int64 {
	known := s.participantEventsCache.Known()
	for p, pid := range s.participants.ByPubKey {
		if known[pid.ID] == -1 {
			root, ok := s.rootsByParticipant[p]
			if ok {
				known[pid.ID] = root.SelfParent.Index
			}
		}
	}
	return known
}

// ConsensusEvents returns all consensus events
func (s *InmemStore) ConsensusEvents() []string {
	lastWindow, _ := s.consensusCache.GetLastWindow()
	res := make([]string, len(lastWindow))
	for i, item := range lastWindow {
		res[i] = item.(string)
	}
	return res
}

// ConsensusEventsCount returns count of all consnesus events
func (s *InmemStore) ConsensusEventsCount() int64 {
	s.totConsensusEventsLocker.RLock()
	defer s.totConsensusEventsLocker.RUnlock()
	return s.totConsensusEvents
}

// AddConsensusEvent to store
func (s *InmemStore) AddConsensusEvent(event Event) error {
	s.totConsensusEventsLocker.Lock()
	defer s.totConsensusEventsLocker.Unlock()
	s.consensusCache.Set(event.Hex(), s.totConsensusEvents)
	s.totConsensusEvents++
	s.lastConsensusEvents[event.GetCreator()] = event.Hex()
	return nil
}

// GetRoundCreated retrieves created round by ID
func (s *InmemStore) GetRoundCreated(r int64) (RoundCreated, error) {
	res, ok := s.roundCreatedCache.Get(r)
	if !ok {
		return *NewRoundCreated(), cm.NewStoreErr("RoundCreatedCache", cm.KeyNotFound, strconv.FormatInt(r, 10))
	}
	return res.(RoundCreated), nil
}

// SetRoundCreated stores created round by ID
func (s *InmemStore) SetRoundCreated(r int64, round RoundCreated) error {
	s.lastRoundLocker.Lock()
	defer s.lastRoundLocker.Unlock()
	s.roundCreatedCache.Add(r, round)
	if r > s.lastRound {
		s.lastRound = r
	}
	return nil
}

// GetRoundReceived gets received round by ID
func (s *InmemStore) GetRoundReceived(r int64) (RoundReceived, error) {
	res, ok := s.roundReceivedCache.Get(r)
	if !ok {
		return *NewRoundReceived(), cm.NewStoreErr("RoundReceivedCache", cm.KeyNotFound, strconv.FormatInt(r, 10))
	}
	return res.(RoundReceived), nil
}

// SetRoundReceived stores received round by ID
func (s *InmemStore) SetRoundReceived(r int64, round RoundReceived) error {
	s.lastRoundLocker.Lock()
	defer s.lastRoundLocker.Unlock()
	s.roundReceivedCache.Add(r, round)
	if r > s.lastRound {
		s.lastRound = r
	}
	return nil
}

// LastRound getter
func (s *InmemStore) LastRound() int64 {
	s.lastRoundLocker.RLock()
	defer s.lastRoundLocker.RUnlock()
	return s.lastRound
}

// RoundClothos all clothos for the specified round
func (s *InmemStore) RoundClothos(r int64) []string {
	round, err := s.GetRoundCreated(r)
	if err != nil {
		return []string{}
	}
	return round.Clotho()
}

// RoundEvents returns events for the round
func (s *InmemStore) RoundEvents(r int64) int {
	round, err := s.GetRoundCreated(r)
	if err != nil {
		return 0
	}
	return len(round.Message.Events)
}

// GetRoot for participant
func (s *InmemStore) GetRoot(participant string) (Root, error) {
	res, ok := s.rootsByParticipant[participant]
	if !ok {
		return Root{}, cm.NewStoreErr("RootCache", cm.KeyNotFound, participant)
	}
	return res, nil
}

// GetBlock for index
func (s *InmemStore) GetBlock(index int64) (Block, error) {
	res, ok := s.blockCache.Get(index)
	if !ok {
		return Block{}, cm.NewStoreErr("BlockCache", cm.KeyNotFound, strconv.FormatInt(index, 10))
	}
	return res.(Block), nil
}

// SetBlock TODO
func (s *InmemStore) SetBlock(block Block) error {
	s.lastBlockLocker.Lock()
	defer s.lastBlockLocker.Unlock()
	index := block.Index()
	_, err := s.GetBlock(index)
	if err != nil && !cm.Is(err, cm.KeyNotFound) {
		return err
	}
	s.blockCache.Add(index, block)
	if index > s.lastBlock {
		s.lastBlock = index
	}
	return nil
}

// LastBlockIndex getter
func (s *InmemStore) LastBlockIndex() int64 {
	s.lastBlockLocker.RLock()
	defer s.lastBlockLocker.RUnlock()
	return s.lastBlock
}

// GetFrame by index
func (s *InmemStore) GetFrame(index int64) (Frame, error) {
	res, ok := s.frameCache.Get(index)
	if !ok {
		return Frame{}, cm.NewStoreErr("FrameCache", cm.KeyNotFound, strconv.FormatInt(index, 10))
	}
	return res.(Frame), nil
}

// SetFrame in the store
func (s *InmemStore) SetFrame(frame Frame) error {
	index := frame.Round
	_, err := s.GetFrame(index)
	if err != nil && !cm.Is(err, cm.KeyNotFound) {
		return err
	}
	s.frameCache.Add(index, frame)
	return nil
}

// Reset resets the store
func (s *InmemStore) Reset(roots map[string]Root) error {
	// FIXIT: Should we recreate blockCache, frameCache and participantEventsCache here as well
	//        and reset lastConsensusEvents ?
	s.rootsByParticipant = roots
	s.rootsBySelfParent = nil
	s.eventCache.Purge()
	s.roundCreatedCache.Purge()
	s.roundReceivedCache.Purge()
	s.consensusCache = cm.NewRollingIndex("ConsensusCache", s.cacheSize)
	err := s.participantEventsCache.Reset()
	s.lastRoundLocker.Lock()
	s.lastRound = -1
	s.lastRoundLocker.Unlock()
	s.lastBlockLocker.Lock()
	s.lastBlock = -1
	s.lastBlockLocker.Unlock()

	if _, err := s.RootsBySelfParent(); err != nil {
		return err
	}

	return err
}

// Close the store
func (s *InmemStore) Close() error {
	return nil
}

// NeedBoostrap for the store
func (s *InmemStore) NeedBoostrap() bool {
	return false
}

// StorePath getter
func (s *InmemStore) StorePath() string {
	return ""
}
