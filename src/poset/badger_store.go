package poset

import (
	"fmt"
	"os"

	"github.com/dgraph-io/badger"
	"github.com/1lann/cete"

	"github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/common/hexutil"
//	"github.com/Fantom-foundation/go-lachesis/src/kvdb"
	"github.com/Fantom-foundation/go-lachesis/src/peers"
	"github.com/Fantom-foundation/go-lachesis/src/pos"
	"github.com/Fantom-foundation/go-lachesis/src/state"
)

const (
	participantPrefix   = "participant"
	rootSuffix          = "root"
	roundCreatedPrefix  = "roundCreated"
	roundReceivedPrefix = "roundReceived"
	topoPrefix          = "topo"
	blockPrefix         = "block"
	framePrefix         = "frame"
	statePrefix         = "state"
	EVENTS_TBL          = "events"
	TOPO_IDX            = "Message.TopologicalIndex"
	CREATOR_IDX         = "Message.Body.Creator,Message.Body.Index"
	FRAMERECEIVED_IDX   = "FrameReceived"
	SORT_IDX            = "Frame,LamportTimestamp,AtroposTimestamp,Message.Hash"
	FRAMEFINALITY_IDX   = "FrameReceived,Frame" // WIP: finality for frame: no records with FrameReceived=0
	CLOTHOCHK_TBL       = "clotho_chk"
	CLOTHOCREATORCHK_TBL= "clotho_creator_chk"
	TIMETABLE_TBL       = "time_table"
	PEERS_TBL           = "peers"
)

// BadgerStore struct for badger config data
type BadgerStore struct {
	participants  *peers.Peers
	inmemStore    *InmemStore
	db            *cete.DB
	path          string
	needBootstrap bool

	states    state.Database
	stateRoot common.Hash
}

// NewBadgerStore creates a brand new Store with a new database
func NewBadgerStore(participants *peers.Peers, cacheSize int, path string, posConf *pos.Config) (*BadgerStore, error) {
	inmemStore := NewInmemStore(participants, cacheSize, posConf)
	opts := badger.DefaultOptions
//	opts.Dir = path
//	opts.ValueDir = path
	opts.SyncWrites = false
	handle, err := cete.Open(path, opts)
	if err != nil {
		return nil, err
	}
	store := &BadgerStore{
		participants: participants,
		inmemStore:   inmemStore,
		db:           handle,
		path:         path,
//		states: state.NewDatabase(
//			kvdb.NewTable(
//				kvdb.NewBadgerDatabase(
//					handle), statePrefix)),
	}
	if err := store.db.NewTable(EVENTS_TBL); err != nil {
		return nil, err
	}
	if err:= store.db.Table(EVENTS_TBL).NewIndex(TOPO_IDX); err != nil {
		return nil, err
	}
	if err:= store.db.Table(EVENTS_TBL).NewIndex(CREATOR_IDX); err != nil {
		return nil, err
	}
	if err:= store.db.Table(EVENTS_TBL).NewIndex(FRAMERECEIVED_IDX); err != nil {
		return nil, err
	}
	if err:= store.db.Table(EVENTS_TBL).NewIndex(SORT_IDX); err != nil {
		return nil, err
	}
	if err:= store.db.Table(EVENTS_TBL).NewIndex(FRAMEFINALITY_IDX); err != nil {
		return nil, err
	}

	if err := store.db.NewTable(CLOTHOCHK_TBL); err != nil {
		return nil, err
	}

	if err := store.db.NewTable(CLOTHOCREATORCHK_TBL); err != nil {
		return nil, err
	}

	if err := store.db.NewTable(TIMETABLE_TBL); err != nil {
		return nil, err
	}

	if err := store.db.NewTable(PEERS_TBL); err != nil {
		return nil, err
	}

	if err := store.dbSetParticipants(participants); err != nil {
		return nil, err
	}
	if err := store.dbSetRoots(inmemStore.rootsByParticipant); err != nil {
		return nil, err
	}

	// TODO: replace with real genesis
	//store.stateRoot, err = pos.FakeGenesis(participants, posConf, store.states)
	//if err != nil {
	//	return nil, err
	//}

	return store, nil
}

// LoadBadgerStore creates a Store from an existing database
func LoadBadgerStore(cacheSize int, path string) (*BadgerStore, error) {

	if _, err := os.Stat(path); err != nil {
		return nil, err
	}

	opts := badger.DefaultOptions
//	opts.Dir = path
//	opts.ValueDir = path
	opts.SyncWrites = false
	handle, err := cete.Open(path, opts)
	if err != nil {
		return nil, err
	}
	store := &BadgerStore{
		db:            handle,
		path:          path,
		needBootstrap: true,
//		states: state.NewDatabase(
//			kvdb.NewTable(
//				kvdb.NewBadgerDatabase(
//					handle), statePrefix)),
	}

	participants, err := store.dbGetParticipants()
	if err != nil {
		return nil, err
	}

	inmemStore := NewInmemStore(participants, cacheSize, nil)

	// read roots from db and put them in InmemStore
	roots := make(map[string]Root)
	for p := range participants.ByPubKey {
		root, err := store.dbGetRoot(p)
		if err != nil {
			return nil, err
		}
		roots[p] = root
	}

	if err := inmemStore.Reset(roots); err != nil {
		return nil, err
	}

	store.participants = participants
	store.inmemStore = inmemStore

	return store, nil
}

// LoadOrCreateBadgerStore load or create a new badger store
func LoadOrCreateBadgerStore(participants *peers.Peers, cacheSize int, path string, posConf *pos.Config) (*BadgerStore, error) {
	store, err := LoadBadgerStore(cacheSize, path)

	if err != nil {
		fmt.Println("Could not load store - creating new")
		store, err = NewBadgerStore(participants, cacheSize, path, posConf)

		if err != nil {
			return nil, err
		}
	}

	return store, nil
}

// ==============================================================================
// Keys

func topologicalEventKey(index int64) []byte {
	return []byte(fmt.Sprintf("%s_%09d", topoPrefix, index))
}

func participantKey(participant string) []byte {
	return []byte(fmt.Sprintf("%s_%s", participantPrefix, participant))
}

func participantEventKey(participant string, index int64) []byte {
	return []byte(fmt.Sprintf("%s__event_%09d", participant, index))
}

func participantRootKey(participant string) []byte {
	return []byte(fmt.Sprintf("%s_%s", participant, rootSuffix))
}

func roundCreatedKey(index int64) []byte {
	return []byte(fmt.Sprintf("%s_%09d", roundCreatedPrefix, index))
}

func roundReceivedKey(index int64) []byte {
	return []byte(fmt.Sprintf("%s_%09d", roundReceivedPrefix, index))
}
func blockKey(index int64) []byte {
	return []byte(fmt.Sprintf("%s_%09d", blockPrefix, index))
}

func frameKey(index int64) []byte {
	return []byte(fmt.Sprintf("%s_%09d", framePrefix, index))
}

func checkClothoKey(frame int64, hash EventHash) string {
	return string(fmt.Sprintf("%09d_%s", frame, hash.String()))
}

func checkClothoCreatorKey(frame int64, creatorID uint64) string {
	return string(fmt.Sprintf("%09d_%d", frame, creatorID))
}

func timeTableKey(hash EventHash) string {
	return string(fmt.Sprintf("timeTable_%s", hash.String()))
}

/*
 * Store interface implementation:
 */

// TopologicalEvents returns event in topological order.
func (s *BadgerStore) TopologicalEvents() ([]Event, error) {
	var res []Event

	r := s.db.Table(EVENTS_TBL).Index(TOPO_IDX).Between(cete.MinValue, cete.MaxValue)
	for r.Next() {
		var result Event
		r.Decode(&result)
		res = append(res, result)
	}
	if r.Error() != cete.ErrEndOfRange {
		return res, fmt.Errorf("%v", r.Error())
	}
	return res, nil
}

// CacheSize returns the cache size for the store
func (s *BadgerStore) CacheSize() int {
	return s.inmemStore.CacheSize()
}

// Participants returns all participants in the store
func (s *BadgerStore) Participants() (*peers.Peers, error) {
	return s.participants, nil
}

// RootsBySelfParent returns Self Parent's EventHash map of the roots
func (s *BadgerStore) RootsBySelfParent() map[EventHash]Root {
	return s.inmemStore.RootsBySelfParent()
}

// RootsByParticipant returns PubKeyHex map of the roots
func (s *BadgerStore) RootsByParticipant() map[string]Root {
	return s.inmemStore.RootsByParticipant()
}

// GetEventBlock get specific event block by hash
func (s *BadgerStore) GetEventBlock(hash EventHash) (event Event, err error) {
	// try to get it from cache
	event, err = s.inmemStore.GetEventBlock(hash)
	// if not in cache, try to get it from db
	if err != nil {
		event, err = s.dbGetEventBlock(hash)
	}
	return event, mapError(err, "Event", hash.String())
}

// SetEvent set a specific event
func (s *BadgerStore) SetEvent(event Event) error {
	// try to add it to the cache
	if err := s.inmemStore.SetEvent(event); err != nil {
		return err
	}
	// try to add it to the db
	return s.dbSetEvents([]Event{event})
}

// ParticipantEvents return all participant events
func (s *BadgerStore) ParticipantEvents(participant string, skip int64) (EventHashes, error) {
	res, err := s.inmemStore.ParticipantEvents(participant, skip)
	if err != nil {
		res, err = s.dbParticipantEvents(participant, skip)
	}
	return res, err
}

// ParticipantEvent get specific participant event
func (s *BadgerStore) ParticipantEvent(participant string, index int64) (EventHash, error) {
	result, err := s.inmemStore.ParticipantEvent(participant, index)
	if err != nil {
		result, err = s.dbParticipantEvent(participant, index)
	}
	return result, mapError(err, "ParticipantEvent", string(participantEventKey(participant, index)))
}

// LastEventFrom returns the last event for a participant
func (s *BadgerStore) LastEventFrom(participant string) (last EventHash, isRoot bool, err error) {
	return s.inmemStore.LastEventFrom(participant)
}

// LastConsensusEventFrom returns the last consensus events for a participant
func (s *BadgerStore) LastConsensusEventFrom(participant string) (last EventHash, isRoot bool, err error) {
	return s.inmemStore.LastConsensusEventFrom(participant)
}

// ConsensusEvents returns all consensus events
func (s *BadgerStore) ConsensusEvents() EventHashes {
	return s.inmemStore.ConsensusEvents()
}

// ConsensusEventsCount returns the count for all known consensus events
func (s *BadgerStore) ConsensusEventsCount() int64 {
	return s.inmemStore.ConsensusEventsCount()
}

// AddConsensusEvent adds a consensus event to the store
func (s *BadgerStore) AddConsensusEvent(event Event) error {
	return s.inmemStore.AddConsensusEvent(event)
}

// GetRoundCreated gets the created round info for a given index
func (s *BadgerStore) GetRoundCreated(r int64) (RoundCreated, error) {
	res, err := s.inmemStore.GetRoundCreated(r)
	if err != nil {
		res, err = s.dbGetRoundCreated(r)
	}
	return res, mapError(err, "RoundCreated", string(roundCreatedKey(r)))
}

// SetRoundCreated sets the created round info for a given index
func (s *BadgerStore) SetRoundCreated(r int64, round RoundCreated) error {
	if err := s.inmemStore.SetRoundCreated(r, round); err != nil {
		return err
	}
	return s.dbSetRoundCreated(r, round)
}

// GetRoundReceived gets the received round for a given index
func (s *BadgerStore) GetRoundReceived(r int64) (RoundReceived, error) {
	res, err := s.inmemStore.GetRoundReceived(r)
	if err != nil {
		res, err = s.dbGetRoundReceived(r)
	}
	return res, mapError(err, "RoundReceived", string(roundReceivedKey(r)))
}

// SetRoundReceived sets the received round info for a given index
func (s *BadgerStore) SetRoundReceived(r int64, round RoundReceived) error {
	if err := s.inmemStore.SetRoundReceived(r, round); err != nil {
		return err
	}
	return s.dbSetRoundReceived(r, round)
}

// LastRound returns the last round for the store
func (s *BadgerStore) LastRound() int64 {
	return s.inmemStore.LastRound()
}

// RoundClothos returns all clothos for a round
func (s *BadgerStore) RoundClothos(r int64) EventHashes {
	round, err := s.GetRoundCreated(r)
	if err != nil {
		return EventHashes{}
	}
	return round.Clotho()
}

// RoundEvents returns all events for a round
func (s *BadgerStore) RoundEvents(r int64) int {
	round, err := s.GetRoundCreated(r)
	if err != nil {
		return 0
	}
	return len(round.Message.Events)
}

// GetRoot returns the root for a participant
func (s *BadgerStore) GetRoot(participant string) (Root, error) {
	root, err := s.inmemStore.GetRoot(participant)
	if err != nil {
		root, err = s.dbGetRoot(participant)
	}
	return root, mapError(err, "Root", string(participantRootKey(participant)))
}

// GetBlock returns the block for a given index
func (s *BadgerStore) GetBlock(rr int64) (Block, error) {
	res, err := s.inmemStore.GetBlock(rr)
	if err != nil {
		res, err = s.dbGetBlock(rr)
	}
	return res, mapError(err, "Block", string(blockKey(rr)))
}

// SetBlock add a block
func (s *BadgerStore) SetBlock(block Block) error {
	if err := s.inmemStore.SetBlock(block); err != nil {
		return err
	}
	return s.dbSetBlock(block)
}

// LastBlockIndex returns the last block index (height)
func (s *BadgerStore) LastBlockIndex() int64 {
	return s.inmemStore.LastBlockIndex()
}

// GetFrame returns a specific frame for the index
func (s *BadgerStore) GetFrame(rr int64) (Frame, error) {
	res, err := s.inmemStore.GetFrame(rr)
	if err != nil {
		res, err = s.dbGetFrame(rr)
	}
	return res, mapError(err, "Frame", string(frameKey(rr)))
}

// SetFrame add a frame
func (s *BadgerStore) SetFrame(frame Frame) error {
	if err := s.inmemStore.SetFrame(frame); err != nil {
		return err
	}
	return s.dbSetFrame(frame)
}

// Reset all roots
func (s *BadgerStore) Reset(roots map[string]Root) error {
	return s.inmemStore.Reset(roots)
}

// Close badger
func (s *BadgerStore) Close() error {
	if err := s.inmemStore.Close(); err != nil {
		return err
	}
	s.db.Close()
	return nil
}

// NeedBootstrap checks if bootstrapping is required
func (s *BadgerStore) NeedBootstrap() bool {
	return s.needBootstrap
}

// StorePath returns the path to the file on disk
func (s *BadgerStore) StorePath() string {
	return s.path
}

// StateDB returns state database
func (s *BadgerStore) StateDB() state.Database {
	return s.states
}

// StateRoot returns genesis state hash.
func (s *BadgerStore) StateRoot() common.Hash {
	return s.stateRoot
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// DB Methods

func (s *BadgerStore) dbGetEventBlock(hash EventHash) (Event, error) {
	var event Event
	_, err := s.db.Table(EVENTS_TBL).Get(hash.String(), &event)
	if err != nil {
		return Event{}, err
	}
	return event, nil
}

func (s *BadgerStore) dbSetEvents(events []Event) error {

	for _, event := range events {
		eventHash := event.Hash()
		err := s.db.Table(EVENTS_TBL).Set(eventHash.String(), event)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *BadgerStore) dbParticipantEvents(participant string, skip int64) (res EventHashes, err error) {

	creator, err := hexutil.Decode(participant)
	if err != nil {
		return nil, err
	}

	r := s.db.Table(EVENTS_TBL).Index(CREATOR_IDX).Between(
		[]interface{}{creator, cete.MinValue}, []interface{}{creator, cete.MaxValue})

	for r.Next() {
		var result Event
		r.Decode(&result)
		res = append(res, result.Hash())
	}
	if r.Error() != cete.ErrEndOfRange {
		return res, fmt.Errorf("%v", r.Error())
	}
	return res, nil
}

func (s *BadgerStore) dbParticipantEvent(participant string, index int64) (hash EventHash, err error) {

	creator, err := hexutil.Decode(participant)
	if err != nil {
		return
	}
	var result Event

	_, _, err = s.db.Table(EVENTS_TBL).Index(CREATOR_IDX).One(
		[]interface{}{creator, index}, &result)

	if err == nil {
		hash = result.Hash()
	}
	return
}

func (s *BadgerStore) dbSetRoots(roots map[string]Root) error {
//	tx := s.db.NewTransaction(true)
//	defer tx.Discard()
//	for participant, root := range roots {
//		val, err := root.ProtoMarshal()
//		if err != nil {
//			return err
//		}
//		key := participantRootKey(participant)
//		// fmt.Println("Setting root", participant, "->", key)
//		// insert [participant_root] => [root bytes]
//		if err := tx.Set(key, val); err != nil {
//			return err
//		}
//	}
//	return tx.Commit(nil)
	return nil
}

func (s *BadgerStore) dbGetRoot(participant string) (Root, error) {
//	var rootBytes []byte
//	key := participantRootKey(participant)
//	err := s.db.View(func(txn *badger.Txn) error {
//		item, err := txn.Get(key)
//		if err != nil {
//			return err
//		}
//		err = item.Value(func(val []byte) error {
//			rootBytes = val
//			return nil
//		})
//		return err
//	})

//	if err != nil {
//		return Root{}, err
//	}

	root := new(Root)
//	if err := root.ProtoUnmarshal(rootBytes); err != nil {
//		return Root{}, err
//	}

	return *root, nil
}

func (s *BadgerStore) dbGetRoundCreated(index int64) (RoundCreated, error) {
//	var roundBytes []byte
//	key := roundCreatedKey(index)
//	err := s.db.View(func(txn *badger.Txn) error {
//		item, err := txn.Get(key)
//		if err != nil {
//			return err
//		}
//		err = item.Value(func(val []byte) error {
//			roundBytes = val
//			return nil
//		})
//		return err
//	})

//	if err != nil {
//		return *NewRoundCreated(), err
//	}

	roundInfo := new(RoundCreated)
//	if err := roundInfo.ProtoUnmarshal(roundBytes); err != nil {
//		return *NewRoundCreated(), err
//	}
	// In the current design, Queued field must be re-calculated every time for
	// each round. When retrieving a round info from a database, this field
	// should be ignored.
	roundInfo.Message.Queued = false

	return *roundInfo, nil
}

func (s *BadgerStore) dbSetRoundCreated(index int64, round RoundCreated) error {
//	tx := s.db.NewTransaction(true)
//	defer tx.Discard()
//
//	key := roundCreatedKey(index)
//	val, err := round.ProtoMarshal()
//	if err != nil {
//		return err
//	}

	// insert [round_index] => [round bytes]
//	if err := tx.Set(key, val); err != nil {
//		return err
//	}

//	return tx.Commit(nil)
	return nil
}

func (s *BadgerStore) dbGetRoundReceived(index int64) (RoundReceived, error) {
//	var roundBytes []byte
//	key := roundReceivedKey(index)
//	err := s.db.View(func(txn *badger.Txn) error {
//		item, err := txn.Get(key)
//		if err != nil {
//			return err
//		}
//		err = item.Value(func(val []byte) error {
//			roundBytes = val
//			return nil
//		})
//		return err
//	})
//
//	if err != nil {
//		return *NewRoundReceived(), err
//	}

	roundInfo := new(RoundReceived)
//	if err := roundInfo.ProtoUnmarshal(roundBytes); err != nil {
//		return *NewRoundReceived(), err
//	}

	return *roundInfo, nil
}

func (s *BadgerStore) dbSetRoundReceived(index int64, round RoundReceived) error {
//	tx := s.db.NewTransaction(true)
//	defer tx.Discard()
//
//	key := roundReceivedKey(index)
//	val, err := round.ProtoMarshal()
//	if err != nil {
//		return err
//	}

	// insert [round_index] => [round bytes]
//	if err := tx.Set(key, val); err != nil {
//		return err
//	}

//	return tx.Commit(nil)
	return nil
}

func (s *BadgerStore) dbGetParticipants() (*peers.Peers, error) {
	res := peers.NewPeers()

	r := s.db.Table(PEERS_TBL).All()
	for r.Next() {
		var result peers.Peer
		r.Decode(&result)
		res.AddPeer(&result)
	}
	if r.Error() != cete.ErrEndOfRange {
		return res, fmt.Errorf("%v", r.Error())
	}
	return res, nil
}

func (s *BadgerStore) dbSetParticipants(participants *peers.Peers) error {
	participants.RLock()
	defer participants.RUnlock()
	for pubKey, peer := range participants.ByPubKey {
		err := s.db.Table(PEERS_TBL).Set(pubKey, peer)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *BadgerStore) dbGetBlock(index int64) (Block, error) {
//	var blockBytes []byte
//	key := blockKey(index)
//	err := s.db.View(func(txn *badger.Txn) error {
//		item, err := txn.Get(key)
//		if err != nil {
//			return err
//		}
//		err = item.Value(func(val []byte) error {
//			blockBytes = val
//			return nil
//		})
//		return err
//	})

//	if err != nil {
//		return Block{}, err
//	}

	block := new(Block)
//	if err := block.ProtoUnmarshal(blockBytes); err != nil {
//		return Block{}, err
//	}

	return *block, nil
}

func (s *BadgerStore) dbSetBlock(block Block) error {
//	tx := s.db.NewTransaction(true)
//	defer tx.Discard()
//
//	key := blockKey(block.Index())
//	val, err := block.ProtoMarshal()
//	if err != nil {
//		return err
//	}
//
	// insert [index] => [block bytes]
//	if err := tx.Set(key, val); err != nil {
//		return err
//	}

//	return tx.Commit(nil)
	return nil
}

func (s *BadgerStore) dbGetFrame(index int64) (Frame, error) {
//	var frameBytes []byte
//	key := frameKey(index)
//	err := s.db.View(func(txn *badger.Txn) error {
//		item, err := txn.Get(key)
//		if err != nil {
//			return err
//		}
//		err = item.Value(func(val []byte) error {
//			frameBytes = val
//			return nil
//		})
//		return err
//	})

//	if err != nil {
//		return Frame{}, err
//	}

	frame := new(Frame)
//	if err := frame.ProtoUnmarshal(frameBytes); err != nil {
//		return Frame{}, err
//	}

	return *frame, nil
}

func (s *BadgerStore) dbSetFrame(frame Frame) error {
//	tx := s.db.NewTransaction(true)
//	defer tx.Discard()
//
//	key := frameKey(frame.Round)
//	val, err := frame.ProtoMarshal()
//	if err != nil {
//		return err
//	}

	// insert [index] => [block bytes]
//	if err := tx.Set(key, val); err != nil {
//		return err
//	}

//	return tx.Commit(nil)
	return nil
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func isDBKeyNotFound(err error) bool {
	return err == badger.ErrKeyNotFound || err == cete.ErrNotFound
}

func mapError(err error, name, key string) error {
	if err != nil {
		if isDBKeyNotFound(err) {
			return common.NewStoreErr(name, common.KeyNotFound, key)
		}
	}
	return err
}


// GetClothoCheck retrieves EventHash by frame + EventHash
func (s *BadgerStore) GetClothoCheck(frame int64, hash EventHash) (EventHash, error) {
	res, err := s.inmemStore.GetClothoCheck(frame, hash)
	if err != nil {
		res, err = s.dbGetClothoCheck(frame, hash)
	}
	return res, err // mapError(err, "badger_store GetClothoCheck", string(checkClothoKey(frame, hash)))
}

// GetClothoCreatorCheck retrieves EventHash by frame + creator
func (s *BadgerStore) GetClothoCreatorCheck(frame int64, creatorID uint64) (EventHash, error) {
	res, err := s.inmemStore.GetClothoCreatorCheck(frame, creatorID)
	if err != nil {
		res, err = s.dbGetClothoCreatorCheck(frame, creatorID)
	}
	return res, err // mapError(err, "badger_store GetClothoCreatorCheck", string(checkClothoCreatorKey(frame, hash)))
}

// AddClothoCheck to store
func (s *BadgerStore) AddClothoCheck(frame int64, creatorID uint64, hash EventHash) error {
	if err := s.inmemStore.AddClothoCheck(frame, creatorID, hash); err != nil {
		return err
	}
	return s.dbAddClothoCheck(frame, creatorID, hash)
}

func (s *BadgerStore) dbGetClothoCheck(frame int64, phash EventHash) (EventHash, error) {
	var hash EventHash
	key := checkClothoKey(frame, phash)
	_, err := s.db.Table(CLOTHOCHK_TBL).Get(key, &hash)
	if err != nil {
		return EventHash{}, err
	}
	return hash, nil
}

func (s *BadgerStore) dbGetClothoCreatorCheck(frame int64, creatorID uint64) (EventHash, error) {
	var hash EventHash
	key := checkClothoCreatorKey(frame, creatorID)
	_, err := s.db.Table(CLOTHOCREATORCHK_TBL).Get(key, &hash)
	if err != nil {
		return EventHash{}, err
	}
	return hash, nil
}

func (s *BadgerStore) dbAddClothoCheck(frame int64, creatorID uint64, hash EventHash) error {

	key := checkClothoKey(frame, hash)

	// insert [frame EventHash] => [EventHash]
	if err := s.db.Table(CLOTHOCHK_TBL).Set(key, hash); err != nil {
		return err
	}

	key = checkClothoCreatorKey(frame, creatorID)

	// insert [frame EventHash] => [EventHash]
	if err := s.db.Table(CLOTHOCREATORCHK_TBL).Set(key, hash); err != nil {
		return err
	}

	return nil
}

// AddTimeTable adds lamport timestamp for pair of events for voting in atropos time selection
func (s *BadgerStore) AddTimeTable(hashTo EventHash, hashFrom EventHash, lamportTime int64) error {
	if err := s.inmemStore.AddTimeTable(hashTo, hashFrom, lamportTime); err != nil {
		return err
	}
	return s.dbAddTimeTable(hashTo, hashFrom, lamportTime)
}

// GetTimeTable retrieve FlagTable with lamport time votes in atropos time selection for specified EventHash
func (s *BadgerStore) GetTimeTable(hash EventHash) (FlagTable, error) {
	res, err := s.inmemStore.GetTimeTable(hash)
	if err != nil {
		res, err = s.dbGetTimeTable(hash)
	}
	return res, err // mapError(err, "badger_store GetTimeTable", string(timeTableKey(hash)))
}

func (s *BadgerStore) dbAddTimeTable(hashTo EventHash, hashFrom EventHash, lamportTime int64) error {
	ft := NewFlagTable()
	key := timeTableKey(hashTo)
	_, err := s.db.Table(TIMETABLE_TBL).Get(key, &ft)
	if err != cete.ErrNotFound {
		return err
	}

	ft[hashFrom] = lamportTime

	if err := s.db.Table(TIMETABLE_TBL).Set(key, ft); err != nil {
		return err
	}

	return nil
}

func (s *BadgerStore) dbGetTimeTable(hash EventHash) (FlagTable, error) {
	ft := NewFlagTable()
	key := timeTableKey(hash)
	_, err := s.db.Table(TIMETABLE_TBL).Get(key, &ft)
	if err != cete.ErrNotFound {
		return nil, err
	}
	return ft, nil
}

// CheckFrameFinality checks if a frame is ready to push out in consensus order
func (s *BadgerStore) CheckFrameFinality(frame int64) bool {
	_, _, err := s.db.Table(EVENTS_TBL).Index(FRAMEFINALITY_IDX).One(
		[]interface{}{0, frame}, nil)
	if err == cete.ErrNotFound {
		return true
	}
	return false
}

func (s *BadgerStore) ProcessOutFrame(frame int64, address string) error {
	file, err := os.OpenFile(fmt.Sprintf("Node_%v.finality", address), os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("*** Open  err: %v", err)
		return err
	}
	defer file.Close()

	r := s.db.Table(EVENTS_TBL).Index(SORT_IDX).Between(
		[]interface{}{frame, cete.MinValue, cete.MinValue, cete.MinValue},
		[]interface{}{frame, cete.MaxValue, cete.MaxValue, cete.MaxValue})
	for r.Next() {
		var ev Event
		r.Decode(&ev)
		if ev.IsLoaded() {
			hash := ev.Hash()
			fmt.Fprintf(file, "%v:%v:%v:%v:%v\n",
				hash.String(), ev.Frame, ev.FrameReceived, ev.LamportTimestamp, ev.AtroposTimestamp)
		}
	}
	if r.Error() != cete.ErrEndOfRange {
		return fmt.Errorf("%v", r.Error())
	}
	return nil
}
