package poset

import (
	"fmt"
	"os"
	
	cm "github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/peers"
	"github.com/dgraph-io/badger"
	"github.com/golang/protobuf/proto"
)

const (
	participantPrefix   = "participant"
	rootSuffix          = "root"
	roundCreatedPrefix  = "roundCreated"
	roundReceivedPrefix = "roundReceived"
	topoPrefix          = "topo"
	blockPrefix         = "block"
	framePrefix         = "frame"
)

type BadgerStore struct {
	inmemStore   *InmemStore
	db           *badger.DB
	path         string
	needBoostrap bool
}

// NewBadgerStore creates a brand new Store with a new database
func NewBadgerStore(peerSet *peers.PeerSet, cacheSize int, path string) (*BadgerStore, error) {
	inmemStore := NewInmemStore(peerSet, cacheSize)
	opts := badger.DefaultOptions
	opts.Dir = path
	opts.ValueDir = path
	opts.SyncWrites = false
	handle, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	store := &BadgerStore{
		inmemStore: inmemStore,
		db:         handle,
		path:       path,
	}
	if err := store.dbSetPeerSet(peerSet); err != nil {
		return nil, err
	}
	if err := store.dbSetRoots(inmemStore.rootsByParticipant); err != nil {
		return nil, err
	}
	if err := store.dbSetRootEvents(inmemStore.rootsByParticipant); err != nil {
		return nil, err
	}
	return store, nil
}

// LoadBadgerStore creates a Store from an existing database
func LoadBadgerStore(cacheSize int, path string) (*BadgerStore, error) {

	if _, err := os.Stat(path); err != nil {
		return nil, err
	}

	opts := badger.DefaultOptions
	opts.Dir = path
	opts.ValueDir = path
	opts.SyncWrites = false
	handle, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	store := &BadgerStore{
		db:           handle,
		path:         path,
		needBoostrap: true,
	}

	peerSet, err := store.dbGetPeerSet()
	if err != nil {
		return nil, err
	}

	inmemStore := NewInmemStore(peerSet, cacheSize)

	// read roots from db and put them in InmemStore
	roots := make(map[string]*Root)
	for p := range peerSet.ByPubKey {
		root, err := store.dbGetRoot(p)
		if err != nil {
			return nil, err
		}
		roots[p] = root
	}

	// XXX temp
	// Get last frame
	// Should there be an initila frame?

	frame := &Frame{
		Round: 0,
		Peers: peerSet.Peers,
		Roots: roots,
	}

	if err := inmemStore.Reset(frame); err != nil {
		return nil, err
	}

	store.inmemStore = inmemStore

	return store, nil
}

func LoadOrCreateBadgerStore(peerSet *peers.PeerSet, cacheSize int, path string) (*BadgerStore, error) {
	store, err := LoadBadgerStore(cacheSize, path)

	if err != nil {
		fmt.Println("Could not load store - creating new")
		store, err = NewBadgerStore(peerSet, cacheSize, path)

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

// ==============================================================================
// Implement the Store interface

func (s *BadgerStore) CacheSize() int {
	return s.inmemStore.CacheSize()
}

func (s *BadgerStore) GetLastPeerSet() (*peers.PeerSet, error) {
	return s.inmemStore.GetLastPeerSet()
}

func (s *BadgerStore) GetPeerSet(round int64) (*peers.PeerSet, error) {
	return s.inmemStore.GetPeerSet(round)
}

func (s *BadgerStore) SetPeerSet(round int64, peerSet *peers.PeerSet) error {
	return s.inmemStore.SetPeerSet(round, peerSet)
}

func (s *BadgerStore) RepertoireByPubKey() map[string]*peers.Peer {
	return s.inmemStore.RepertoireByPubKey()
}

func (s *BadgerStore) RepertoireByID() map[uint32]*peers.Peer {
	return s.inmemStore.RepertoireByID()
}
func (s *BadgerStore) RootsBySelfParent() map[string]*Root {
	return s.inmemStore.RootsBySelfParent()
}

func (s *BadgerStore) GetEvent(key string) (event *Event, err error) {
	// try to get it from cache
	event, err = s.inmemStore.GetEvent(key)
	// if not in cache, try to get it from db
	if err != nil {
		event, err = s.dbGetEvent(key)
	}
	return event, mapError(err, "Event", key)
}

func (s *BadgerStore) SetEvent(event *Event) error {
	// try to add it to the cache
	if err := s.inmemStore.SetEvent(event); err != nil {
		return err
	}
	// try to add it to the db
	return s.dbSetEvents([]*Event{event})
}

func (s *BadgerStore) ParticipantEvents(participant string, skip int64) ([]string, error) {
	res, err := s.inmemStore.ParticipantEvents(participant, skip)
	if err != nil {
		res, err = s.dbParticipantEvents(participant, skip)
	}
	return res, err
}

func (s *BadgerStore) ParticipantEvent(participant string, index int64) (string, error) {
	result, err := s.inmemStore.ParticipantEvent(participant, index)
	if err != nil {
		result, err = s.dbParticipantEvent(participant, index)
	}
	return result, mapError(err, "ParticipantEvent", string(participantEventKey(participant, index)))
}

func (s *BadgerStore) LastEventFrom(participant string) (last string, isRoot bool, err error) {
	return s.inmemStore.LastEventFrom(participant)
}

func (s *BadgerStore) LastConsensusEventFrom(participant string) (last string, isRoot bool, err error) {
	return s.inmemStore.LastConsensusEventFrom(participant)
}

func (s *BadgerStore) KnownEvents() map[uint32]int64 {
	known := make(map[uint32]int64)
	peers, _ := s.GetLastPeerSet()
	for p, pid := range peers.ByPubKey {
		index := int64(-1)
		last, isRoot, err := s.LastEventFrom(p)
		if err == nil {
			if isRoot {
				root, err := s.GetRoot(p)
				if err != nil {
					last = root.SelfParent.Hash
					index = root.SelfParent.Index
				}
			} else {
				lastEvent, err := s.GetEvent(last)
				if err == nil {
					index = lastEvent.Index()
				}
			}

		}
		known[pid.ID] = index
	}
	return known
}

func (s *BadgerStore) ConsensusEvents() []string {
	return s.inmemStore.ConsensusEvents()
}

func (s *BadgerStore) ConsensusEventsCount() int64 {
	return s.inmemStore.ConsensusEventsCount()
}

func (s *BadgerStore) AddConsensusEvent(event *Event) error {
	return s.inmemStore.AddConsensusEvent(event)
}

func (s *BadgerStore) GetRoundCreated(r int64) (*RoundCreated, error) {
	res, err := s.inmemStore.GetRoundCreated(r)
	if err != nil {
		res, err = s.dbGetRoundCreated(r)
	}
	return res, mapError(err, "RoundCreated", string(roundCreatedKey(r)))
}

func (s *BadgerStore) SetRoundCreated(r int64, round *RoundCreated) error {
	if err := s.inmemStore.SetRoundCreated(r, round); err != nil {
		return err
	}
	return s.dbSetRoundCreated(r, round)
}

func (s *BadgerStore) GetRoundReceived(r int64) (*RoundReceived, error) {
	res, err := s.inmemStore.GetRoundReceived(r)
	if err != nil {
		res, err = s.dbGetRoundReceived(r)
	}
	return res, mapError(err, "RoundReceived", string(roundReceivedKey(r)))
}

func (s *BadgerStore) SetRoundReceived(r int64, round *RoundReceived) error {
	if err := s.inmemStore.SetRoundReceived(r, round); err != nil {
		return err
	}
	return s.dbSetRoundReceived(r, round)
}

func (s *BadgerStore) LastRound() int64 {
	return s.inmemStore.LastRound()
}

func (s *BadgerStore) RoundWitnesses(r int64) []string {
	round, err := s.GetRoundCreated(r)
	if err != nil {
		return []string{}
	}
	return round.Witnesses()
}

func (s *BadgerStore) RoundEvents(r int64) int {
	round, err := s.GetRoundCreated(r)
	if err != nil {
		return 0
	}
	return len(round.Message.Events)
}

func (s *BadgerStore) GetRoot(participant string) (*Root, error) {
	root, err := s.inmemStore.GetRoot(participant)
	if err != nil {
		root, err = s.dbGetRoot(participant)
	}
	return root, mapError(err, "Root", string(participantRootKey(participant)))
}

func (s *BadgerStore) GetBlock(rr int64) (*Block, error) {
	res, err := s.inmemStore.GetBlock(rr)
	if err != nil {
		res, err = s.dbGetBlock(rr)
	}
	return res, mapError(err, "Block", string(blockKey(rr)))
}

func (s *BadgerStore) SetBlock(block *Block) error {
	if err := s.inmemStore.SetBlock(block); err != nil {
		return err
	}
	return s.dbSetBlock(block)
}

func (s *BadgerStore) LastBlockIndex() int64 {
	return s.inmemStore.LastBlockIndex()
}

func (s *BadgerStore) GetFrame(rr int64) (*Frame, error) {
	res, err := s.inmemStore.GetFrame(rr)
	if err != nil {
		res, err = s.dbGetFrame(rr)
	}
	return res, mapError(err, "Frame", string(frameKey(rr)))
}

func (s *BadgerStore) SetFrame(frame *Frame) error {
	if err := s.inmemStore.SetFrame(frame); err != nil {
		return err
	}
	return s.dbSetFrame(frame)
}

func (s *BadgerStore) Reset(frame *Frame) error {
	return s.inmemStore.Reset(frame)
}

func (s *BadgerStore) Close() error {
	if err := s.inmemStore.Close(); err != nil {
		return err
	}
	return s.db.Close()
}

func (s *BadgerStore) NeedBoostrap() bool {
	return s.needBoostrap
}

func (s *BadgerStore) StorePath() string {
	return s.path
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// DB Methods

func (s *BadgerStore) dbGetEvent(key string) (*Event, error) {
	var eventBytes []byte
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			eventBytes = val
			return nil
		})
		return err
	})

	if err != nil {
		return nil, err
	}

	event := new(Event)
	if err := event.ProtoUnmarshal(eventBytes); err != nil {
		return nil, err
	}

	return event, nil
}

func (s *BadgerStore) dbSetEvents(events []*Event) error {
	tx := s.db.NewTransaction(true)
	defer tx.Discard()

	for _, event := range events {
		eventHex := event.Hex()
		val, err := event.ProtoMarshal()
		if err != nil {
			return err
		}
		// check if it already exists
		existent := false
		_, err = tx.Get([]byte(eventHex))
		if err != nil && isDBKeyNotFound(err) {
			existent = true
		}
		// insert [event hash] => [event bytes]
		if err := tx.Set([]byte(eventHex), val); err != nil {
			return err
		}

		if existent {
			// insert [topo_index] => [event hash]
			topoKey := topologicalEventKey(event.Message.TopologicalIndex)
			if err := tx.Set(topoKey, []byte(eventHex)); err != nil {
				return err
			}
			// insert [participant_index] => [event hash]
			peKey := participantEventKey(event.Creator(), event.Index())
			if err := tx.Set(peKey, []byte(eventHex)); err != nil {
				return err
			}
		}
	}
	return tx.Commit(nil)
}

func (s *BadgerStore) dbTopologicalEvents() ([]*Event, error) {
	var res []*Event
	var evKey string
	t := int64(-1)
	err := s.db.View(func(txn *badger.Txn) error {
		key := topologicalEventKey(t)
		item, errr := txn.Get(key)
		for errr == nil {
			errrr := item.Value(func(v []byte) error {
				evKey = string(v)
				return nil
			})
			if errrr != nil {
				break
			}

			eventItem, err := txn.Get([]byte(evKey))
			if err != nil {
				return err
			}
			err = eventItem.Value(func(eventBytes []byte) error {
				event := new(Event)
				if err := event.ProtoUnmarshal(eventBytes); err != nil {
					return err
				}
				res = append(res, event)
				return nil
			})
			if err != nil {
				return err
			}

			t++
			key = topologicalEventKey(t)
			item, errr = txn.Get(key)
		}

		if !isDBKeyNotFound(errr) {
			return errr
		}

		return nil
	})

	return res, err
}

func (s *BadgerStore) dbParticipantEvents(participant string, skip int64) ([]string, error) {
	var res []string
	err := s.db.View(func(txn *badger.Txn) error {
		i := skip + 1
		key := participantEventKey(participant, i)
		item, errr := txn.Get(key)
		for errr == nil {
			errrr := item.Value(func(v []byte) error {
				res = append(res, string(v))
				return nil
			})
			if errrr != nil {
				break
			}

			i++
			key = participantEventKey(participant, i)
			item, errr = txn.Get(key)
		}

		if !isDBKeyNotFound(errr) {
			return errr
		}

		return nil
	})
	return res, err
}

func (s *BadgerStore) dbParticipantEvent(participant string, index int64) (string, error) {
	var data []byte
	key := participantEventKey(participant, index)
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			data = val
			return nil
		})
		return err
	})
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *BadgerStore) dbSetRoots(roots map[string]*Root) error {
	tx := s.db.NewTransaction(true)
	defer tx.Discard()
	for participant, root := range roots {
		val, err := root.ProtoMarshal()
		if err != nil {
			return err
		}
		key := participantRootKey(participant)
		// 		fmt.Println("Setting root", participant, "->", key)
		// insert [participant_root] => [root bytes]
		if err := tx.Set(key, val); err != nil {
			return err
		}
	}
	return tx.Commit(nil)
}

func (s *BadgerStore) dbSetRootEvents(roots map[string]*Root) error {
	for participant, root := range roots {
		var creator []byte
		fmt.Sscanf(participant, "0x%X", &creator)
		flagTable := map[string]int64{root.SelfParent.Hash: 1}
		ft, _ := proto.Marshal(&FlagTableWrapper{Body: flagTable})
		body := EventBody{
			Creator: creator, /*s.participants.ByPubKey[participant].PubKey,*/
			Index:   root.SelfParent.Index,
			Parents: []string{"", ""},
		}
		event := &Event{
			Message: EventMessage{
				Hex:              root.SelfParent.Hash,
				CreatorID:        root.SelfParent.CreatorID,
				TopologicalIndex: -1,
				Body:             &body,
				FlagTable:        ft,
				LamportTimestamp: 0,
				Round:            0,
				RoundReceived:    0 /*RoundNIL*/,
				WitnessProof:     []string{root.SelfParent.Hash},
			},
		}
		if err := s.SetEvent(event); err != nil {
			return err
		}
	}
	return nil
}

func (s *BadgerStore) dbGetRoot(participant string) (*Root, error) {
	var rootBytes []byte
	key := participantRootKey(participant)
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			rootBytes = val
			return nil
		})
		return err
	})

	if err != nil {
		return nil, err
	}

	root := new(Root)
	if err := root.ProtoUnmarshal(rootBytes); err != nil {
		return nil, err
	}

	return root, nil
}

func (s *BadgerStore) dbGetRoundCreated(index int64) (*RoundCreated, error) {
	var roundBytes []byte
	key := roundCreatedKey(index)
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			roundBytes = val
			return nil
		})
		return err
	})

	if err != nil {
		return NewRoundCreated(), err
	}

	roundInfo := new(RoundCreated)
	if err := roundInfo.ProtoUnmarshal(roundBytes); err != nil {
		return NewRoundCreated(), err
	}

	return roundInfo, nil
}

func (s *BadgerStore) dbSetRoundCreated(index int64, round *RoundCreated) error {
	tx := s.db.NewTransaction(true)
	defer tx.Discard()

	key := roundCreatedKey(index)
	val, err := round.ProtoMarshal()
	if err != nil {
		return err
	}

	// insert [round_index] => [round bytes]
	if err := tx.Set(key, val); err != nil {
		return err
	}

	return tx.Commit(nil)
}

func (s *BadgerStore) dbGetRoundReceived(index int64) (*RoundReceived, error) {
	var roundBytes []byte
	key := roundReceivedKey(index)
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			roundBytes = val
			return nil
		})
		return err
	})

	if err != nil {
		return NewRoundReceived(), err
	}

	roundInfo := new(RoundReceived)
	if err := roundInfo.ProtoUnmarshal(roundBytes); err != nil {
		return NewRoundReceived(), err
	}

	return roundInfo, nil
}

func (s *BadgerStore) dbSetRoundReceived(index int64, round *RoundReceived) error {
	tx := s.db.NewTransaction(true)
	defer tx.Discard()

	key := roundReceivedKey(index)
	val, err := round.ProtoMarshal()
	if err != nil {
		return err
	}

	// insert [round_index] => [round bytes]
	if err := tx.Set(key, val); err != nil {
		return err
	}

	return tx.Commit(nil)
}

func (s *BadgerStore) dbGetPeerSet() (*peers.PeerSet, error) {
	var tempPeers []*peers.Peer

	err := s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(participantPrefix)

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := string(item.Key())

			pubKey := k[len(participantPrefix)+1:]
			tempPeers = append(tempPeers, peers.NewPeer(pubKey, ""))
		}

		return nil
	})

	res := peers.NewPeerSet(tempPeers)
	return res, err
}

func (s *BadgerStore) dbSetPeerSet(peerSet *peers.PeerSet) error {
	tx := s.db.NewTransaction(true)
	defer tx.Discard()

	for participant, peer := range peerSet.ByPubKey {
		key := participantKey(participant)
		val := []byte(fmt.Sprint(peer.ID))
		// insert [participant_participant] => [id]
		if err := tx.Set(key, val); err != nil {
			return err
		}
	}
	return tx.Commit(nil)
}

func (s *BadgerStore) dbGetBlock(index int64) (*Block, error) {
	var blockBytes []byte
	key := blockKey(index)
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			blockBytes = val
			return nil
		})
		return err
	})

	if err != nil {
		return nil, err
	}

	block := new(Block)
	if err := block.ProtoUnmarshal(blockBytes); err != nil {
		return nil, err
	}

	return block, nil
}

func (s *BadgerStore) dbSetBlock(block *Block) error {
	tx := s.db.NewTransaction(true)
	defer tx.Discard()

	key := blockKey(block.Index())
	val, err := block.ProtoMarshal()
	if err != nil {
		return err
	}

	// insert [index] => [block bytes]
	if err := tx.Set(key, val); err != nil {
		return err
	}

	return tx.Commit(nil)
}

func (s *BadgerStore) dbGetFrame(index int64) (*Frame, error) {
	var frameBytes []byte
	key := frameKey(index)
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			frameBytes = val
			return nil
		})
		return err
	})

	if err != nil {
		return nil, err
	}

	frame := new(Frame)
	if err := frame.ProtoUnmarshal(frameBytes); err != nil {
		return nil, err
	}

	return frame, nil
}

func (s *BadgerStore) dbSetFrame(frame *Frame) error {
	tx := s.db.NewTransaction(true)
	defer tx.Discard()

	key := frameKey(frame.Round)
	val, err := frame.ProtoMarshal()
	if err != nil {
		return err
	}

	// insert [index] => [block bytes]
	if err := tx.Set(key, val); err != nil {
		return err
	}

	return tx.Commit(nil)
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func isDBKeyNotFound(err error) bool {
	return err.Error() == badger.ErrKeyNotFound.Error()
}

func mapError(err error, name, key string) error {
	if err != nil {
		if isDBKeyNotFound(err) {
			return cm.NewStoreErr(name, cm.KeyNotFound, key)
		}
	}
	return err
}
