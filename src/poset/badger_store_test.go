package poset

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/Fantom-foundation/go-lachesis/src/crypto"
	"github.com/Fantom-foundation/go-lachesis/src/peers"
)

func initBadgerStore(cacheSize int, t *testing.T) (*BadgerStore, []pub) {
	n := 3
	var participantPubs []pub
	var tempPeers []*peers.Peer
	for i := 0; i < n; i++ {
		key, _ := crypto.GenerateECDSAKey()
		pubKey := crypto.FromECDSAPub(&key.PublicKey)
		peer := peers.NewPeer(fmt.Sprintf("0x%X", pubKey), "")
		tempPeers = append(tempPeers, peer)
		participantPubs = append(participantPubs,
			pub{peer.ID, key, pubKey, peer.PubKeyHex})
	}

	os.RemoveAll("test_data")
	os.Mkdir("test_data", os.ModeDir|0777)
	dir, err := ioutil.TempDir("test_data", "badger")
	if err != nil {
		log.Fatal(err)
	}

	store, err := NewBadgerStore(peers.NewPeerSet(tempPeers), cacheSize, dir)
	if err != nil {
		t.Fatal(err)
	}

	return store, participantPubs
}

func removeBadgerStore(store *BadgerStore, t *testing.T) {
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(store.path); err != nil {
		t.Fatal(err)
	}
}

func createTestDB(dir string, t *testing.T) *BadgerStore {
	peerSet := peers.NewPeerSet([]*peers.Peer{
		peers.NewPeer("0xAA", ""),
		peers.NewPeer("0xBB", ""),
		peers.NewPeer("0xCC", ""),
	})

	cacheSize := 100

	store, err := NewBadgerStore(peerSet, cacheSize, dir)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return store
}

func TestNewBadgerStore(t *testing.T) {
	os.RemoveAll("test_data")
	os.Mkdir("test_data", os.ModeDir|0777)

	dbPath := "test_data/badger"
	store := createTestDB(dbPath, t)
	defer os.RemoveAll(store.path)

	if store.path != dbPath {
		t.Fatalf("unexpected path %q", store.path)
	}
	if _, err := os.Stat(dbPath); err != nil {
		t.Fatalf("err: %s", err)
	}

	// check roots
	inmemRoots := store.inmemStore.rootsByParticipant

	if len(inmemRoots) != 3 {
		t.Fatalf("DB root should have 3 items, not %d", len(inmemRoots))
	}

	for participant, root := range inmemRoots {
		dbRoot, err := store.dbGetRoot(participant)
		if err != nil {
			t.Fatalf("Error retrieving DB root for participant %s: %s", participant, err)
		}
		if !dbRoot.Equals(root) {
			t.Fatalf("%s DB root should be %#v, not %#v", participant, root, dbRoot)
		}
	}

	if err := store.Close(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestLoadBadgerStore(t *testing.T) {
	os.RemoveAll("test_data")
	os.Mkdir("test_data", os.ModeDir|0777)
	dbPath := "test_data/badger"

	// Create the test db
	tempStore := createTestDB(dbPath, t)
	defer os.RemoveAll(tempStore.path)
	tempStore.Close()

	badgerStore, err := LoadBadgerStore(cacheSize, tempStore.path)
	if err != nil {
		t.Fatal(err)
	}

	dbParticipants, err := badgerStore.dbGetPeerSet()
	if err != nil {
		t.Fatal(err)
	}

	peerSet, err := badgerStore.GetLastPeerSet()
	if err != nil {
		t.Fatal(err)
	}

	peerCount := peerSet.Len()
	if peerCount != 3 {
		t.Fatalf("store.participants  length should be %d items, not %d", 3, peerCount)
	}

	if peerCount != dbParticipants.Len() {
		t.Fatalf("store.participants should contain %d items, not %d",
			dbParticipants.Len(),
			peerCount)
	}

	for dbP, dbPeer := range dbParticipants.ByPubKey {
		peer, ok := peerSet.ByPubKey[dbP]
		if !ok {
			t.Fatalf("BadgerStore participants does not contains %s", dbP)
		}
		if peer.ID != dbPeer.ID {
			t.Fatalf("participant %s ID should be %d, not %d", dbP, dbPeer.ID, peer.ID)
		}
	}

}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Call DB methods directly

func TestDBEventMethods(t *testing.T) {
	cacheSize := 1 // Inmem_store's caches accept positive cacheSize only
	testSize := int64(100)
	store, participants := initBadgerStore(cacheSize, t)
	defer removeBadgerStore(store, t)

	// insert events in db directly
	events := make(map[string][]*Event)
	topologicalIndex := int64(0)
	var topologicalEvents []*Event
	for _, p := range participants {
		var items []*Event
		for k := int64(0); k < testSize; k++ {
			event := NewEvent(
				[][]byte{[]byte(fmt.Sprintf("%s_%d", p.hex[:5], k))},
				[]*InternalTransaction{},
				[]BlockSignature{{Validator: []byte("validator"), Index: 0, Signature: "r|s"}},
				[]string{"", ""},
				p.pubKey,
				k, nil)
			event.Sign(p.privKey)
			event.Message.TopologicalIndex = topologicalIndex
			topologicalIndex++
			topologicalEvents = append(topologicalEvents, event)

			items = append(items, event)
			err := store.dbSetEvents([]*Event{event})
			if err != nil {
				t.Fatal(err)
			}
		}
		events[p.hex] = items
	}

	// check events where correctly inserted and can be retrieved
	for p, evs := range events {
		for k, ev := range evs {
			rev, err := store.dbGetEventBlock(ev.Hex())
			if err != nil {
				t.Fatal(err)
			}
			if !ev.Message.Body.Equals(rev.Message.Body) {
				t.Fatalf("events[%s][%d].Body should be %#v, not %#v", p, k, ev.Message.Body, rev.Message.Body)
			}
			if !reflect.DeepEqual(ev.Message.Signature, rev.Message.Signature) {
				t.Fatalf("events[%s][%d].Signature should be %#v, not %#v", p, k, ev.Message.Signature, rev.Message.Signature)
			}
			if ver, err := rev.Verify(); err != nil && !ver {
				t.Fatalf("failed to verify signature. err: %s", err)
			}
		}
	}

	// check topological order of events was correctly created
	dbTopologicalEvents, err := store.dbTopologicalEvents()
	if err != nil {
		t.Fatal(err)
	}
	if len(dbTopologicalEvents) != len(topologicalEvents) {
		t.Fatalf("Length of dbTopologicalEvents should be %d, not %d",
			len(topologicalEvents), len(dbTopologicalEvents))
	}
	for i, dte := range dbTopologicalEvents {
		te := topologicalEvents[i]

		if dte.Hex() != te.Hex() {
			t.Fatalf("dbTopologicalEvents[%d].Hex should be %s, not %s", i,
				te.Hex(),
				dte.Hex())
		}
		if !te.Message.Body.Equals(dte.Message.Body) {
			t.Fatalf("dbTopologicalEvents[%d].Body should be %#v, not %#v", i,
				te.Message.Body,
				dte.Message.Body)
		}
		if !reflect.DeepEqual(te.Message.Signature, dte.Message.Signature) {
			t.Fatalf("dbTopologicalEvents[%d].Signature should be %#v, not %#v", i,
				te.Message.Signature,
				dte.Message.Signature)
		}

		if ver, err := dte.Verify(); err != nil && !ver {
			t.Fatalf("failed to verify signature. err: %s", err)
		}
	}

	// check that participant events where correctly added
	skipIndex := int64(-1) // do not skip any indexes
	for _, p := range participants {
		pEvents, err := store.dbParticipantEvents(p.hex, skipIndex)
		if err != nil {
			t.Fatal(err)
		}
		if l := int64(len(pEvents)); l != testSize {
			t.Fatalf("%s should have %d events, not %d", p.hex, testSize, l)
		}

		expectedEvents := events[p.hex][skipIndex+1:]
		for k, e := range expectedEvents {
			if e.Hex() != pEvents[k] {
				t.Fatalf("ParticipantEvents[%s][%d] should be %s, not %s",
					p.hex, k, e.Hex(), pEvents[k])
			}
		}
	}
}

func TestDBRoundMethods(t *testing.T) {
	cacheSize := 1 // Inmem_store's caches accept positive cacheSize only
	store, participants := initBadgerStore(cacheSize, t)
	defer removeBadgerStore(store, t)

	round := NewRoundCreated()
	events := make(map[string]*Event)
	for _, p := range participants {
		event := NewEvent([][]byte{},
			[]*InternalTransaction{},
			[]BlockSignature{},
			[]string{"", ""},
			p.pubKey,
			0, nil)
		events[p.hex] = event
		round.AddEvent(event.Hex(), true)
	}

	if err := store.dbSetRoundCreated(0, round); err != nil {
		t.Fatal(err)
	}

	storedRound, err := store.dbGetRoundCreated(0)
	if err != nil {
		t.Fatal(err)
	}

	if !round.Equals(storedRound) {
		t.Fatalf("Round and StoredRound do not match")
	}

	clothos := store.RoundClothos(0)
	expectedClothos := round.Clotho()
	if len(clothos) != len(expectedClothos) {
		t.Fatalf("There should be %d clothos, not %d", len(expectedClothos), len(clothos))
	}
	for _, w := range expectedClothos {
		if !contains(clothos, w) {
			t.Fatalf("Clothos should contain %s", w)
		}
	}
}

func TestDBParticipantMethods(t *testing.T) {
	cacheSize := 1 // Inmem_store's caches accept positive cacheSize only
	store, _ := initBadgerStore(cacheSize, t)
	defer removeBadgerStore(store, t)

	peerSet, err := store.GetLastPeerSet()
	if err != nil {
		t.Fatal(err)
	}

	if err := store.dbSetPeerSet(peerSet); err != nil {
		t.Fatal(err)
	}

	participantsFromDB, err := store.dbGetPeerSet()
	if err != nil {
		t.Fatal(err)
	}

	for p, peer := range peerSet.ByPubKey {
		dbPeer, ok := participantsFromDB.ByPubKey[p]
		if !ok {
			t.Fatalf("DB does not contain participant %s", p)
		}
		if peer.ID != dbPeer.ID {
			t.Fatalf("DB participant %s should have ID %d, not %d", p, peer.ID, dbPeer.ID)
		}
	}
}

func TestDBBlockMethods(t *testing.T) {
	cacheSize := 1 // Inmem_store's caches accept positive cacheSize only
	store, participants := initBadgerStore(cacheSize, t)
	defer removeBadgerStore(store, t)

	index := int64(0)
	roundReceived := int64(5)
	transactions := [][]byte{
		[]byte("tx1"),
		[]byte("tx2"),
		[]byte("tx3"),
		[]byte("tx4"),
		[]byte("tx5"),
	}
	internalTransactions := []*InternalTransaction{
		NewInternalTransaction(TransactionType_PEER_ADD, *peers.NewPeer("peer1", "paris")),
		NewInternalTransaction(TransactionType_PEER_REMOVE, *peers.NewPeer("peer2", "london")),
	}
	frameHash := []byte("this is the frame hash")

	block := NewBlock(index, roundReceived, frameHash, transactions, internalTransactions)

	sig1, err := block.Sign(participants[0].privKey)
	if err != nil {
		t.Fatal(err)
	}

	sig2, err := block.Sign(participants[1].privKey)
	if err != nil {
		t.Fatal(err)
	}

	block.SetSignature(sig1)
	block.SetSignature(sig2)

	t.Run("Store Block", func(t *testing.T) {
		if err := store.dbSetBlock(block); err != nil {
			t.Fatal(err)
		}

		storedBlock, err := store.dbGetBlock(index)
		if err != nil {
			t.Fatal(err)
		}

		if !storedBlock.Equals(block) {
			t.Fatalf("Block and StoredBlock do not match")
		}
	})

	t.Run("Check signatures in stored Block", func(t *testing.T) {
		storedBlock, err := store.dbGetBlock(index)
		if err != nil {
			t.Fatal(err)
		}

		val1Sig, ok := storedBlock.Signatures[participants[0].hex]
		if !ok {
			t.Fatalf("Validator1 signature not stored in block")
		}
		if val1Sig != sig1.Signature {
			t.Fatal("Validator1 block signatures differ")
		}

		val2Sig, ok := storedBlock.Signatures[participants[1].hex]
		if !ok {
			t.Fatalf("Validator2 signature not stored in block")
		}
		if val2Sig != sig2.Signature {
			t.Fatal("Validator2 block signatures differ")
		}
	})
}

func TestDBFrameMethods(t *testing.T) {
	cacheSize := 1 // Inmem_store's caches accept positive cacheSize only
	store, participants := initBadgerStore(cacheSize, t)
	defer removeBadgerStore(store, t)

	events := make([]*EventMessage, len(participants))
	roots := make(map[string]*Root)
	for id, p := range participants {
		event := NewEvent(
			[][]byte{[]byte(fmt.Sprintf("%s_%d", p.hex[:5], 0))},
			[]*InternalTransaction{},
			[]BlockSignature{{Validator: []byte("validator"), Index: 0, Signature: "r|s"}},
			[]string{"", ""},
			p.pubKey,
			0, nil)
		event.Sign(p.privKey)
		events[id] = event.Message

		root := NewBaseRoot(uint64(id))
		roots[p.hex] = root
	}
	frame := &Frame{
		Round:  1,
		Events: events,
		Roots:  roots,
	}

	t.Run("Store Frame", func(t *testing.T) {
		if err := store.dbSetFrame(frame); err != nil {
			t.Fatal(err)
		}

		storedFrame, err := store.dbGetFrame(frame.Round)
		if err != nil {
			t.Fatal(err)
		}

		if !storedFrame.Equals(frame) {
			t.Fatalf("Frame and StoredFrame do not match")
		}
	})
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// Check that the wrapper methods work
// These methods use the inmemStore as a cache on top of the DB

func TestBadgerEvents(t *testing.T) {
	// Insert more events than can fit in cache to test retrieving from db.
	cacheSize := 10
	testSize := int64(100)
	store, participants := initBadgerStore(cacheSize, t)
	defer removeBadgerStore(store, t)

	// insert event
	events := make(map[string][]*Event)
	for _, p := range participants {
		var items []*Event
		for k := int64(0); k < testSize; k++ {
			event := NewEvent([][]byte{[]byte(fmt.Sprintf("%s_%d", p.hex[:5], k))},
				[]*InternalTransaction{},
				[]BlockSignature{{Validator: []byte("validator"), Index: 0, Signature: "r|s"}},
				[]string{"", ""},
				p.pubKey,
				k, nil)
			items = append(items, event)
			err := store.SetEvent(event)
			if err != nil {
				t.Fatal(err)
			}
		}
		events[p.hex] = items
	}

	// check that events were correclty inserted
	for p, evs := range events {
		for k, ev := range evs {
			rev, err := store.GetEventBlock(ev.Hex())
			if err != nil {
				t.Fatal(err)
			}
			if !ev.Message.Body.Equals(rev.Message.Body) {
				t.Fatalf("events[%s][%d].Body should be %#v, not %#v", p, k, ev, rev)
			}
			if !reflect.DeepEqual(ev.Message.Signature, rev.Message.Signature) {
				t.Fatalf("events[%s][%d].Signature should be %#v, not %#v", p, k, ev.Message.Signature, rev.Message.Signature)
			}
		}
	}

	// check retrieving events per participant
	skipIndex := int64(-1) // do not skip any indexes
	for _, p := range participants {
		pEvents, err := store.ParticipantEvents(p.hex, skipIndex)
		if err != nil {
			t.Fatal(err)
		}
		if l := int64(len(pEvents)); l != testSize {
			t.Fatalf("%s should have %d events, not %d", p.hex, testSize, l)
		}

		expectedEvents := events[p.hex][skipIndex+1:]
		for k, e := range expectedEvents {
			if e.Hex() != pEvents[k] {
				t.Fatalf("ParticipantEvents[%s][%d] should be %s, not %s",
					p.hex, k, e.Hex(), pEvents[k])
			}
		}
	}

	// check retrieving participant last
	for _, p := range participants {
		last, _, err := store.LastEventFrom(p.hex)
		if err != nil {
			t.Fatal(err)
		}

		evs := events[p.hex]
		expectedLast := evs[len(evs)-1]
		if last != expectedLast.Hex() {
			t.Fatalf("%s last should be %s, not %s", p.hex, expectedLast.Hex(), last)
		}
	}

	expectedKnown := make(map[uint64]int64)
	for _, p := range participants {
		expectedKnown[p.id] = testSize - 1
	}
	known := store.KnownEvents()
	if !reflect.DeepEqual(expectedKnown, known) {
		t.Fatalf("Incorrect Known. Got %#v, expected %#v", known, expectedKnown)
	}

	for _, p := range participants {
		evs := events[p.hex]
		for _, ev := range evs {
			if err := store.AddConsensusEvent(ev); err != nil {
				t.Fatal(err)
			}
		}

	}
}

func TestBadgerRounds(t *testing.T) {
	cacheSize := 1 // Inmem_store's caches accept positive cacheSize only
	store, participants := initBadgerStore(cacheSize, t)
	defer removeBadgerStore(store, t)

	round := NewRoundCreated()
	events := make(map[string]*Event)
	for _, p := range participants {
		event := NewEvent([][]byte{},
			[]*InternalTransaction{},
			[]BlockSignature{},
			[]string{"", ""},
			p.pubKey,
			0, nil)
		events[p.hex] = event
		round.AddEvent(event.Hex(), true)
	}

	if err := store.SetRoundCreated(0, round); err != nil {
		t.Fatal(err)
	}

	if c := store.LastRound(); c != 0 {
		t.Fatalf("Store LastRound should be 0, not %d", c)
	}

	storedRound, err := store.GetRoundCreated(0)
	if err != nil {
		t.Fatal(err)
	}

	if !round.Equals(storedRound) {
		t.Fatalf("Round and StoredRound do not match")
	}

	clothos := store.RoundClothos(0)
	expectedClothos := round.Clotho()
	if len(clothos) != len(expectedClothos) {
		t.Fatalf("There should be %d clothos, not %d", len(expectedClothos), len(clothos))
	}
	for _, w := range expectedClothos {
		if !contains(clothos, w) {
			t.Fatalf("Clothos should contain %s", w)
		}
	}
}

func TestBadgerBlocks(t *testing.T) {
	cacheSize := 1 // Inmem_store's caches accept positive cacheSize only
	store, participants := initBadgerStore(cacheSize, t)
	defer removeBadgerStore(store, t)

	index := int64(0)
	roundReceived := int64(5)
	transactions := [][]byte{
		[]byte("tx1"),
		[]byte("tx2"),
		[]byte("tx3"),
		[]byte("tx4"),
		[]byte("tx5"),
	}
	internalTransactions := []*InternalTransaction{
		NewInternalTransaction(TransactionType_PEER_ADD, *peers.NewPeer("peer1", "paris")),
		NewInternalTransaction(TransactionType_PEER_REMOVE, *peers.NewPeer("peer2", "london")),
	}
	frameHash := []byte("this is the frame hash")
	block := NewBlock(index, roundReceived, frameHash, transactions, internalTransactions)

	sig1, err := block.Sign(participants[0].privKey)
	if err != nil {
		t.Fatal(err)
	}

	sig2, err := block.Sign(participants[1].privKey)
	if err != nil {
		t.Fatal(err)
	}

	block.SetSignature(sig1)
	block.SetSignature(sig2)

	t.Run("Store Block", func(t *testing.T) {
		if err := store.SetBlock(block); err != nil {
			t.Fatal(err)
		}

		storedBlock, err := store.GetBlock(index)
		if err != nil {
			t.Fatal(err)
		}

		if !storedBlock.Equals(block) {
			t.Fatalf("Block and StoredBlock do not match")
		}
	})

	t.Run("Check signatures in stored Block", func(t *testing.T) {
		storedBlock, err := store.GetBlock(index)
		if err != nil {
			t.Fatal(err)
		}

		val1Sig, ok := storedBlock.Signatures[participants[0].hex]
		if !ok {
			t.Fatalf("Validator1 signature not stored in block")
		}
		if val1Sig != sig1.Signature {
			t.Fatal("Validator1 block signatures differ")
		}

		val2Sig, ok := storedBlock.Signatures[participants[1].hex]
		if !ok {
			t.Fatalf("Validator2 signature not stored in block")
		}
		if val2Sig != sig2.Signature {
			t.Fatal("Validator2 block signatures differ")
		}
	})
}

func TestBadgerFrames(t *testing.T) {
	cacheSize := 1 // Inmem_store's caches accept positive cacheSize only
	store, participants := initBadgerStore(cacheSize, t)
	defer removeBadgerStore(store, t)

	events := make([]*EventMessage, len(participants))
	roots := make(map[string]*Root)
	for id, p := range participants {
		event := NewEvent(
			[][]byte{[]byte(fmt.Sprintf("%s_%d", p.hex[:5], 0))},
			[]*InternalTransaction{},
			[]BlockSignature{{Validator: []byte("validator"), Index: 0, Signature: "r|s"}},
			[]string{"", ""},
			p.pubKey,
			0, nil)
		event.Sign(p.privKey)
		events[id] = event.Message

		root := NewBaseRoot(uint64(id))
		roots[p.hex] = root
	}
	frame := &Frame{
		Round:  1,
		Events: events,
		Roots:  roots,
	}

	t.Run("Store Frame", func(t *testing.T) {
		if err := store.SetFrame(frame); err != nil {
			t.Fatal(err)
		}

		storedFrame, err := store.GetFrame(frame.Round)
		if err != nil {
			t.Fatal(err)
		}

		if !storedFrame.Equals(frame) {
			t.Fatalf("Frame and StoredFrame do not match")
		}
	})
}
