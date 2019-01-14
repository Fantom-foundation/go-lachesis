package poset

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/peers"
	"github.com/hashicorp/golang-lru"
)

// InmemStore struct
type InmemStore struct {
	cacheSize              int
	eventCache             *lru.Cache           // hash => Event
	roundCreatedCache      *lru.Cache           // round number => RoundCreated
	roundReceivedCache     *lru.Cache           // round received number => RoundReceived
	blockCache             *lru.Cache           // index => Block
	frameCache             *lru.Cache           // round received => Frame
	consensusCache         *common.RollingIndex // consensus index => hash
	totConsensusEvents     int64
	peerSetCache           *PeerSetCache // start round => PeerSet
	repertoireByPubKey     map[string]*peers.Peer
	repertoireByID         map[uint64]*peers.Peer
	participantEventsCache *ParticipantEventsCache // pubkey => Events
	rootsByParticipant     map[string]*Root        // [participant] => Root
	rootsBySelfParent      map[EventHash]*Root      // [Root.SelfParent.Hash] => Root
	lastRound              int64
	lastConsensusEvents    map[string]EventHash // [participant] => hex() of last consensus event
	lastBlock              int64

	lastRoundLocker          sync.RWMutex
	lastBlockLocker          sync.RWMutex
	totConsensusEventsLocker sync.RWMutex
}

// NewInmemStore constructor
func NewInmemStore(peerSet *peers.PeerSet, cacheSize int) *InmemStore {
	rootsByParticipant := make(map[string]*Root)

	for pk, pid := range peerSet.ByPubKey {
		root := NewBaseRoot(pid.ID)
		rootsByParticipant[pk] = root
	}

	eventCache, err := lru.New(cacheSize)
	if err != nil {
		fmt.Println("Unable to init InmemStore.eventCache:", err)
		os.Exit(31)
	}
	roundCreatedCache, err := lru.New(cacheSize)
	if err != nil {
		fmt.Println("Unable to init InmemStore.roundCreatedCache:", err)
		os.Exit(32)
	}
	roundReceivedCache, err := lru.New(cacheSize)
	if err != nil {
		fmt.Println("Unable to init InmemStore.roundReceivedCache:", err)
		os.Exit(35)
	}
	blockCache, err := lru.New(cacheSize)
	if err != nil {
		fmt.Println("Unable to init InmemStore.blockCache:", err)
		os.Exit(33)
	}
	frameCache, err := lru.New(cacheSize)
	if err != nil {
		fmt.Println("Unable to init InmemStore.frameCache:", err)
		os.Exit(34)
	}

	store := &InmemStore{
		cacheSize: cacheSize,
		eventCache:             eventCache,
		roundCreatedCache:      roundCreatedCache,
		roundReceivedCache:     roundReceivedCache,
		blockCache:             blockCache,
		frameCache:             frameCache,
		consensusCache:         common.NewRollingIndex("ConsensusCache", cacheSize),
		peerSetCache:           NewPeerSetCache(),
		repertoireByPubKey:     make(map[string]*peers.Peer),
		repertoireByID:         make(map[uint64]*peers.Peer),
		participantEventsCache: NewParticipantEventsCache(cacheSize, peerSet),
		rootsByParticipant:     rootsByParticipant,
		lastRound:              -1,
		lastBlock:              -1,
		lastConsensusEvents:    map[string]EventHash{},
	}

	err = store.SetPeerSet(0, peerSet)
	if err != nil {
		fmt.Println("Unable to init InmemStore.peerSetCache", err)
		os.Exit(36)
	}
	return store
}

// CacheSize size of cache
func (s *InmemStore) CacheSize() int {
	return s.cacheSize
}

// GetLastPeerSet returns last PeerSet
func (s *InmemStore) GetLastPeerSet() (*peers.PeerSet, error) {
	return s.peerSetCache.GetLast()
}

// GetPeerSet gets the PeerSet for the given round
func (s *InmemStore) GetPeerSet(round int64) (*peers.PeerSet, error) {
	return s.peerSetCache.Get(round)
}

// SetPeerSet stores the PeerSet given for the round provided
func (s *InmemStore) SetPeerSet(round int64, peerSet *peers.PeerSet) error {
	// Update PeerSetCache
	err := s.peerSetCache.Set(round, peerSet)
	if err != nil {
		return err
	}

	// Extend PartipantEventsCache and Roots with new peers
	for id, p := range peerSet.ByID {
		if _, ok := s.participantEventsCache.participants.ByID[id]; !ok {
			if err := s.participantEventsCache.AddPeer(p); err != nil {
				return err
			}
		}

		if _, ok := s.rootsByParticipant[p.PubKeyHex]; !ok {
			root := NewBaseRoot(p.ID)
			s.rootsByParticipant[p.PubKeyHex] = root
			s.rootsBySelfParent = nil
			s.RootsBySelfParent()
		}

		s.repertoireByPubKey[p.PubKeyHex] = p
		s.repertoireByID[p.ID] = p
	}

	return nil
}

// RepertoireByPubKey retrieves cached PubKey map of peers
func (s *InmemStore) RepertoireByPubKey() map[string]*peers.Peer {
	return s.repertoireByPubKey
}

// RepertoireByID retrieve cached ID map of peers
func (s *InmemStore) RepertoireByID() map[uint64]*peers.Peer {
	return s.repertoireByID
}

// RootsBySelfParent TODO
func (s *InmemStore) RootsBySelfParent() map[EventHash]*Root {
	if s.rootsBySelfParent == nil {
		s.rootsBySelfParent = make(map[EventHash]*Root)
		for _, root := range s.rootsByParticipant {
			var hash EventHash
			hash.Set(root.SelfParent.Hash)
			s.rootsBySelfParent[hash] = root
		}
	}
	return s.rootsBySelfParent
}

// GetEventBlock gets specific event block by hash
func (s *InmemStore) GetEventBlock(hash EventHash) (*Event, error) {
	res, ok := s.eventCache.Get(hash)
	if !ok {
		return nil, common.NewStoreErr("EventCache", common.KeyNotFound, hash.String())
	}

	return res.(*Event), nil
}

// SetEvent set event for event block
func (s *InmemStore) SetEvent(event *Event) error {
	eventHash := event.Hash()
	_, err := s.GetEventBlock(eventHash)
	if err != nil && !common.Is(err, common.KeyNotFound) {
		return err
	}
	if common.Is(err, common.KeyNotFound) {
		if err := s.addParticpantEvent(event.GetCreator(), eventHash, event.Index()); err != nil {
			return err
		}
	}

	// fmt.Println("Adding event to cache", event.Hex())
	s.eventCache.Add(eventHash, event)

	return nil
}

func (s *InmemStore) addParticpantEvent(participant string, hash EventHash, index int64) error {
	return s.participantEventsCache.Set(participant, hash, index)
}

// ParticipantEvents events for the participant
func (s *InmemStore) ParticipantEvents(participant string, skip int64) (EventHashes, error) {
	return s.participantEventsCache.Get(participant, skip)
}

// ParticipantEvent specific event
func (s *InmemStore) ParticipantEvent(participant string, index int64) (hash EventHash, err error) {
	hash, err = s.participantEventsCache.GetItem(participant, index)
	if err == nil {
		return
	}

	root, ok := s.rootsByParticipant[participant]
	if !ok {
		err = common.NewStoreErr("InmemStore.Roots", common.NoRoot, participant)
		return
	}

	if root.SelfParent.Index == index {
		hash.Set(root.SelfParent.Hash)
		err = nil
	}
	return
}

// LastEventFrom participant
func (s *InmemStore) LastEventFrom(participant string) (last EventHash, isRoot bool, err error) {
	// try to get the last event from this participant
	last, err = s.participantEventsCache.GetLast(participant)
	if err == nil || !common.Is(err, common.Empty) {
		return
	}
	// if there is none, grab the root
	if root, ok := s.rootsByParticipant[participant]; ok {
		last.Set(root.SelfParent.Hash)
		isRoot = true
		err = nil
	} else {
		err = common.NewStoreErr("InmemStore.Roots", common.NoRoot, participant)
	}
	return
}

// LastConsensusEventFrom participant
func (s *InmemStore) LastConsensusEventFrom(participant string) (last EventHash, isRoot bool, err error) {
	// try to get the last consensus event from this participant
	last, ok := s.lastConsensusEvents[participant]
	if ok {
		return
	}
	// if there is none, grab the root
	root, ok := s.rootsByParticipant[participant]
	if ok {
		last.Set(root.SelfParent.Hash)
		isRoot = true
	} else {
		err = common.NewStoreErr("InmemStore.Roots", common.NoRoot, participant)
	}

	return
}

// KnownEvents returns all known events
func (s *InmemStore) KnownEvents() map[uint64]int64 {
	known := s.participantEventsCache.Known()
	lastPeerSet, _ := s.GetLastPeerSet()
	for p, pid := range lastPeerSet.ByPubKey {
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
func (s *InmemStore) ConsensusEvents() EventHashes {
	lastWindow, _ := s.consensusCache.GetLastWindow()
	res := make(EventHashes, len(lastWindow))
	for i, item := range lastWindow {
		res[i] = item.(EventHash)
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
func (s *InmemStore) AddConsensusEvent(event *Event) error {
	s.totConsensusEventsLocker.Lock()
	defer s.totConsensusEventsLocker.Unlock()
	err := s.consensusCache.Set(event.Hash(), s.totConsensusEvents)
	if err != nil {
		return err
	}
	s.totConsensusEvents++
	s.lastConsensusEvents[event.GetCreator()] = event.Hash()
	return nil
}

// GetRoundCreated retrieves created round by ID
func (s *InmemStore) GetRoundCreated(r int64) (*RoundCreated, error) {
	res, ok := s.roundCreatedCache.Get(r)
	if !ok {
		return NewRoundCreated(), common.NewStoreErr("RoundCreatedCache", common.KeyNotFound, strconv.FormatInt(r, 10))
	}
	return res.(*RoundCreated), nil
}

// SetRoundCreated stores created round by ID
func (s *InmemStore) SetRoundCreated(r int64, round *RoundCreated) error {
	s.lastRoundLocker.Lock()
	defer s.lastRoundLocker.Unlock()
	s.roundCreatedCache.Add(r, round)
	if r > s.lastRound {
		s.lastRound = r
	}
	return nil
}

// GetRoundReceived gets received round by ID
func (s *InmemStore) GetRoundReceived(r int64) (*RoundReceived, error) {
	res, ok := s.roundReceivedCache.Get(r)
	if !ok {
		return nil, common.NewStoreErr("RoundReceivedCache", common.KeyNotFound, strconv.FormatInt(r, 10))
	}
	return res.(*RoundReceived), nil
}

// SetRoundReceived stores received round by ID
func (s *InmemStore) SetRoundReceived(r int64, round *RoundReceived) error {
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
func (s *InmemStore) RoundClothos(r int64) EventHashes {
	round, err := s.GetRoundCreated(r)
	if err != nil {
		return EventHashes{}
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
func (s *InmemStore) GetRoot(participant string) (*Root, error) {
	res, ok := s.rootsByParticipant[participant]
	if !ok {
		return nil, common.NewStoreErr("RootCache", common.KeyNotFound, participant)
	}
	return res, nil
}

// GetBlock for index
func (s *InmemStore) GetBlock(index int64) (*Block, error) {
	res, ok := s.blockCache.Get(index)
	if !ok {
		return nil, common.NewStoreErr("BlockCache", common.KeyNotFound, strconv.FormatInt(index, 10))
	}
	return res.(*Block), nil
}

// SetBlock TODO
func (s *InmemStore) SetBlock(block *Block) error {
	s.lastBlockLocker.Lock()
	defer s.lastBlockLocker.Unlock()
	index := block.Index()
	_, err := s.GetBlock(index)
	if err != nil && !common.Is(err, common.KeyNotFound) {
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
func (s *InmemStore) GetFrame(index int64) (*Frame, error) {
	res, ok := s.frameCache.Get(index)
	if !ok {
		return nil, common.NewStoreErr("FrameCache", common.KeyNotFound, strconv.FormatInt(index, 10))
	}
	return res.(*Frame), nil
}

// SetFrame in the store
func (s *InmemStore) SetFrame(frame *Frame) error {
	index := frame.Round
	_, err := s.GetFrame(index)
	if err != nil && !common.Is(err, common.KeyNotFound) {
		return err
	}
	s.frameCache.Add(index, frame)
	return nil
}

// Reset resets the store
func (s *InmemStore) Reset(frame *Frame) error {
	eventCache, errr := lru.New(s.cacheSize)
	if errr != nil {
		fmt.Println("Unable to reset InmemStore.eventCache:", errr)
		os.Exit(41)
	}
	roundCache, errr := lru.New(s.cacheSize)
	if errr != nil {
		fmt.Println("Unable to reset InmemStore.roundCreatedCache:", errr)
		os.Exit(42)
	}
	roundReceivedCache, errr := lru.New(s.cacheSize)
	if errr != nil {
		fmt.Println("Unable to reset InmemStore.roundReceivedCache:", errr)
		os.Exit(45)
	}
	// FIXIT: Should we recreate blockCache, frameCache and participantEventsCache here as well
	//        and reset lastConsensusEvents ?
	// Reset Root caches
	s.rootsByParticipant = frame.Roots
	s.rootsBySelfParent = nil

	// Reset Peer caches
	peerSet := peers.NewPeerSet(frame.Peers)
	_ = s.peerSetCache.Set(frame.Round, peerSet) // ignore key already exists error
	s.participantEventsCache = NewParticipantEventsCache(s.cacheSize, peerSet)

	// Reset Event and Round caches
	s.eventCache = eventCache
	s.roundCreatedCache = roundCache
	s.roundReceivedCache = roundReceivedCache
	s.consensusCache = common.NewRollingIndex("ConsensusCache", s.cacheSize)

	s.lastRoundLocker.Lock()
	s.lastRound = -1
	s.lastRoundLocker.Unlock()
	s.lastBlockLocker.Lock()
	s.lastBlock = -1
	s.lastBlockLocker.Unlock()

	_ = s.RootsBySelfParent()

	return s.SetFrame(frame)
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
