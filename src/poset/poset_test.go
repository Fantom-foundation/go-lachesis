package poset

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/crypto"
	"github.com/Fantom-foundation/go-lachesis/src/peers"
)

const (
	e0  = "e0"
	e1  = "e1"
	e2  = "e2"
	e10 = "e10"
	e21 = "e21"
	e02 = "e02"
	f1  = "f1"
	f0  = "f0"
	f2  = "f2"
	g1  = "g1"
	g0  = "g0"
	g2  = "g2"
	g10 = "g10"
	h0  = "h0"
	h2  = "h2"
	h10 = "h10"
	h21 = "h21"
	i1  = "i1"
	i0  = "i0"
	i2  = "i2"
	e01 = "e01"
	s20 = "s20"
	s10 = "s10"
	s00 = "s00"
	e20 = "e20"
	e12 = "e12"
	a   = "a"
	s11 = "s11"
	w00 = "w00"
	w01 = "w01"
	w02 = "w02"
	w03 = "w03"
	a23 = "a23"
	a00 = "a00"
	a12 = "a12"
	a10 = "a10"
	a21 = "a21"
	w13 = "w13"
	w12 = "w12"
	w11 = "w11"
	w10 = "w10"
	b21 = "b21"
	w23 = "w23"
	b00 = "b00"
	w21 = "w21"
	c10 = "c10"
	w22 = "w22"
	w20 = "w20"
	w31 = "w31"
	w32 = "w32"
	w33 = "w33"
	w30 = "w30"
	d13 = "d13"
	w40 = "w40"
	w41 = "w41"
	w42 = "w42"
	w43 = "w43"
	e23 = "e23"
	w51 = "w51"
	e32 = "e32"
	g13 = "g13"
	f01 = "f01"
	i32 = "i32"
	r0  = "r0"
	r1  = "r1"
	r2  = "r2"
	f2b = "f2b"
	g0x = "g0x"
	h0b = "h0b"
	j2  = "j2"
	j0  = "j0"
	j1  = "j1"
	k0  = "k0"
	k2  = "k2"
	k10 = "k10"
	l2  = "l2"
	l0  = "l0"
	l1  = "l1"
	m0  = "m0"
	m2  = "m2"
)

var (
	cacheSize = 100
	n         = 3
	badgerDir = "test_data/badger"
)

type TestNode struct {
	ID     int
	Pub    []byte
	PubHex string
	Key    *ecdsa.PrivateKey
	Events []Event
}

func NewTestNode(key *ecdsa.PrivateKey, id int) TestNode {
	pub := crypto.FromECDSAPub(&key.PublicKey)
	ID := common.Hash32(pub)
	node := TestNode{
		ID:     ID,
		Key:    key,
		Pub:    pub,
		PubHex: fmt.Sprintf("0x%X", pub),
		Events: []Event{},
	}
	return node
}

func (node *TestNode) signAndAddEvent(event Event, name string,
	index map[string]string, orderedEvents *[]Event) {
	event.Sign(node.Key)
	node.Events = append(node.Events, event)
	index[name] = event.Hex()
	*orderedEvents = append(*orderedEvents, event)
}

type ancestryItem struct {
	descendant, ancestor string
	val                  bool
	err                  bool
}

type roundItem struct {
	event string
	round int64
}

type play struct {
	to          int
	index       int64
	selfParent  string
	otherParent string
	name        string
	txPayload   [][]byte
	sigPayload  []BlockSignature
	knownRoots  []string
}

func testLogger(t testing.TB) *logrus.Entry {
	return common.NewTestLogger(t).WithField("id", "test")
}

/* Initialisation functions */

func initPosetNodes(n int) ([]TestNode, map[string]string,
	*[]Event, *peers.Peers) {
	var (
		participants  = peers.NewPeers()
		orderedEvents = &[]Event{}
		nodes         = make([]TestNode, 0)
		index         = make(map[string]string)
		keys          = make(map[string]*ecdsa.PrivateKey)
	)

	for i := 0; i < n; i++ {
		key, _ := crypto.GenerateECDSAKey()
		pub := crypto.FromECDSAPub(&key.PublicKey)
		pubHex := fmt.Sprintf("0x%X", pub)
		participants.AddPeer(peers.NewPeer(pubHex, ""))
		keys[pubHex] = key
	}

	for i, peer := range participants.ToPeerSlice() {
		nodes = append(nodes, NewTestNode(keys[peer.PubKeyHex], i))
	}

	return nodes, index, orderedEvents, participants
}

func playEvents(plays []play, nodes []TestNode,
	index map[string]string, orderedEvents *[]Event) {
	for _, p := range plays {
		ft := make(map[string]int64)
		for k := range p.knownRoots {
			ft[index[p.knownRoots[k]]] = 1
		}

		e := NewEvent(p.txPayload, nil,
			p.sigPayload,
			[]string{index[p.selfParent], index[p.otherParent]},
			nodes[p.to].Pub, p.index, ft)

		nodes[p.to].signAndAddEvent(e, p.name, index, orderedEvents)
	}
}

func createPoset(t testing.TB, db bool, orderedEvents *[]Event,
	participants *peers.Peers,
	logger *logrus.Entry) *Poset {
	var store Store
	if db {
		var err error
		store, err = NewBadgerStore(participants, cacheSize, badgerDir)
		if err != nil {
			t.Fatal("ERROR creating badger store", err)
		}
	} else {
		store = NewInmemStore(participants, cacheSize)
	}

	poset := NewPoset(participants, store, nil, logger)

	for i, ev := range *orderedEvents {
		if err := poset.InsertEvent(ev, true); err != nil {
			t.Fatalf("failed to insert event %d: %s", i, err)
		}
	}

	return poset
}

func initPosetFull(t testing.TB, plays []play, db bool, n int,
	logger *logrus.Entry) (*Poset, map[string]string, *[]Event, []TestNode) {
	nodes, index, orderedEvents, participants := initPosetNodes(n)

	// Needed to have sorted nodes based on participants hash32
	for i, peer := range participants.ToPeerSlice() {
		event := NewEvent(nil, nil, nil, []string{rootSelfParent(peer.ID), ""},
			nodes[i].Pub, 0, map[string]int64{rootSelfParent(peer.ID): 1})
		nodes[i].signAndAddEvent(event, fmt.Sprintf("e%d", i),
			index, orderedEvents)
	}

	playEvents(plays, nodes, index, orderedEvents)

	poset := createPoset(t, db, orderedEvents, participants, logger)

	// Add reference to each participants' root event
	for i, peer := range participants.ToPeerSlice() {
		root, err := poset.Store.GetRoot(peer.PubKeyHex)
		if err != nil {
			panic(err)
		}
		index["r"+strconv.Itoa(i)] = root.SelfParent.Hash
	}

	return poset, index, orderedEvents, nodes
}

/*  */

/*
|  e12  |
|   | \ |
|  s10 e20
|   | / |
|   /   |
| / |   |
s00 |  s20
|   |   |
e01 |   |
| \ |   |
e0  e1  e2
|   |   |
r0  r1  r2
0   1   2
*/
func initPoset(t *testing.T) (*Poset, map[string]string) {
	plays := []play{
		{0, 1, e0, e1, e01, nil, nil, []string{e0, e1}},
		{2, 1, e2, "", s20, nil, nil, []string{e2}},
		{1, 1, e1, "", s10, nil, nil, []string{e1}},
		{0, 2, e01, "", s00, nil, nil, []string{e0, e1}},
		{2, 2, s20, s00, e20, nil, nil, []string{e0, e1, e2}},
		{1, 2, s10, e20, e12, nil, nil, []string{e0, e1, e2}},
	}

	p, index, orderedEvents, _ := initPosetFull(t, plays, false, n,
		testLogger(t))

	for i, ev := range *orderedEvents {
		if err := p.Store.SetEvent(ev); err != nil {
			t.Fatalf("%d: %s", i, err)
		}
	}

	return p, index
}

func TestAncestor(t *testing.T) {
	p, index := initPoset(t)

	expected := []ancestryItem{
		// first generation
		{e01, e0, true, false},
		{e01, e1, true, false},
		{s00, e01, true, false},
		{s20, e2, true, false},
		{e20, s00, true, false},
		{e20, s20, true, false},
		{e12, e20, true, false},
		{e12, s10, true, false},
		// second generation
		{s00, e0, true, false},
		{s00, e1, true, false},
		{e20, e01, true, false},
		{e20, e2, true, false},
		{e12, e1, true, false},
		{e12, s20, true, false},
		// third generation
		{e20, e0, true, false},
		{e20, e1, true, false},
		{e20, e2, true, false},
		{e12, e01, true, false},
		{e12, e0, true, false},
		{e12, e1, true, false},
		{e12, e2, true, false},
		// false positive
		{e01, e2, false, false},
		{s00, e2, false, false},
		{e0, "", false, true},
		{s00, "", false, true},
		{e12, "", false, true},
		// root events
		{e1, r1, true, false},
		{e20, r1, true, false},
		{e12, r0, true, false},
		{s20, r1, false, false},
		{r0, r1, false, false},
	}

	for _, exp := range expected {
		a, err := p.ancestor(index[exp.descendant], index[exp.ancestor])
		if err != nil && !exp.err {
			t.Fatalf("Error computing ancestor(%s, %s). Err: %v",
				exp.descendant, exp.ancestor, err)
		}
		if a != exp.val {
			t.Fatalf("ancestor(%s, %s) should be %v, not %v",
				exp.descendant, exp.ancestor, exp.val, a)
		}
	}
}

func TestSelfAncestor(t *testing.T) {
	p, index := initPoset(t)

	expected := []ancestryItem{
		// 1 generation
		{e01, e0, true, false},
		{s00, e01, true, false},
		// 1 generation false negative
		{e01, e1, false, false},
		{e12, e20, false, false},
		{s20, e1, false, false},
		{s20, "", false, true},
		// 2 generations
		{e20, e2, true, false},
		{e12, e1, true, false},
		// 2 generations false negatives
		{e20, e0, false, false},
		{e12, e2, false, false},
		{e20, e01, false, false},
		// roots
		{e20, r2, true, false},
		{e1, r1, true, false},
		{e1, r0, false, false},
		{r1, r0, false, false},
	}

	for _, exp := range expected {
		a, err := p.selfAncestor(index[exp.descendant], index[exp.ancestor])
		if err != nil && !exp.err {
			t.Fatalf("Error computing selfAncestor(%s, %s). Err: %v",
				exp.descendant, exp.ancestor, err)
		}
		if a != exp.val {
			t.Fatalf("selfAncestor(%s, %s) should be %v, not %v",
				exp.descendant, exp.ancestor, exp.val, a)
		}
	}
}

func TestSee(t *testing.T) {
	p, index := initPoset(t)

	expected := []ancestryItem{
		{e01, e0, true, false},
		{e01, e1, true, false},
		{e20, e0, true, false},
		{e20, e01, true, false},
		{e12, e01, true, false},
		{e12, e0, true, false},
		{e12, e1, true, false},
		{e12, s20, true, false},
	}

	for _, exp := range expected {
		a, err := p.see(index[exp.descendant], index[exp.ancestor])
		if err != nil && !exp.err {
			t.Fatalf("Error computing see(%s, %s). Err: %v",
				exp.descendant, exp.ancestor, err)
		}
		if a != exp.val {
			t.Fatalf("see(%s, %s) should be %v, not %v",
				exp.descendant, exp.ancestor, exp.val, a)
		}
	}
}

func TestLamportTimestamp(t *testing.T) {
	p, index := initPoset(t)

	expectedTimestamps := map[string]int64{
		e0:  0,
		e1:  0,
		e2:  0,
		e01: 1,
		s10: 1,
		s20: 1,
		s00: 2,
		e20: 3,
		e12: 4,
	}

	for e, ets := range expectedTimestamps {
		ts, err := p.lamportTimestamp(index[e])
		if err != nil {
			t.Fatalf("Error computing lamportTimestamp(%s). Err: %s", e, err)
		}
		if ts != ets {
			t.Fatalf("%s LamportTimestamp should be %d, not %d", e, ets, ts)
		}
	}
}

/*
|    |    e20
|    |   / |
|    | /   |
|    /     |
|  / |     |
e01  |     |
| \  |     |
|   \|     |
|    |\    |
|    |  \  |
e0   e1 (a)e2
0    1     2

Node 2 Forks; events a and e2 are both created by node2, they are not
self-parent sand yet they are both ancestors of event e20
*/
func TestFork(t *testing.T) {
	index := make(map[string]string)
	var nodes []TestNode
	participants := peers.NewPeers()

	for i := 0; i < n; i++ {
		key, _ := crypto.GenerateECDSAKey()
		node := NewTestNode(key, i)
		nodes = append(nodes, node)
		participants.AddPeer(peers.NewPeer(node.PubHex, ""))
	}

	store := NewInmemStore(participants, cacheSize)
	poset := NewPoset(participants, store, nil, testLogger(t))

	for i, node := range nodes {
		event := NewEvent(nil, nil, nil, []string{"", ""}, node.Pub, 0, nil)
		event.Sign(node.Key)
		index[fmt.Sprintf("e%d", i)] = event.Hex()
		poset.InsertEvent(event, true)
	}

	//a and e2 need to have different hashes
	eventA := NewEvent([][]byte{[]byte("yo")}, nil, nil, []string{"", ""}, nodes[2].Pub, 0, nil)
	eventA.Sign(nodes[2].Key)
	index["a"] = eventA.Hex()
	if err := poset.InsertEvent(eventA, true); err == nil {
		t.Fatal("InsertEvent should return error for 'a'")
	}

	event01 := NewEvent(nil, nil, nil,
		[]string{index[e0], index[a]}, //e0 and a
		nodes[0].Pub, 1, nil)
	event01.Sign(nodes[0].Key)
	index[e01] = event01.Hex()
	if err := poset.InsertEvent(event01, true); err == nil {
		t.Fatalf("InsertEvent should return error for %s", e01)
	}

	event20 := NewEvent(nil, nil, nil,
		[]string{index[e2], index[e01]}, //e2 and e01
		nodes[2].Pub, 1, nil)
	event20.Sign(nodes[2].Key)
	index[e20] = event20.Hex()
	if err := poset.InsertEvent(event20, true); err == nil {
		t.Fatalf("InsertEvent should return error for %s", e20)
	}
}

/*
|  s11  |
|   |   |
|   f1  |
|  /|   |
| / s10 |
|/  |   |
e02 |   |
| \ |   |
|   \   |
|   | \ |
s00 |  e21
|   | / |
|  e10  s20
| / |   |
e0  e1  e2
0   1    2
*/

func initRoundPoset(t *testing.T) (*Poset, map[string]string, []TestNode) {
	plays := []play{
		{1, 1, e1, e0, e10, nil, nil, []string{e0, e1}},
		{2, 1, e2, "", s20, nil, nil, []string{e2}},
		{0, 1, e0, "", s00, nil, nil, []string{e0}},
		{2, 2, s20, e10, e21, nil, nil, []string{e0, e1, e2}},
		{0, 2, s00, e21, e02, nil, nil, []string{e0, e21}},
		{1, 2, e10, "", s10, nil, nil, []string{e0, e1}},
		{1, 3, s10, e02, f1, nil, nil, []string{e21, e02, e1}},
		{1, 4, f1, "", s11, [][]byte{[]byte("abc")}, nil,
			[]string{e21, e02, f1}},
	}

	p, index, _, nodes := initPosetFull(t, plays, false, n, testLogger(t))

	return p, index, nodes
}

func TestInsertEvent(t *testing.T) {
	p, index, _ := initRoundPoset(t)

	checkParents := func(e, selfAncestor, ancestor string) bool {
		ev, err := p.Store.GetEvent(index[e])
		if err != nil {
			t.Fatal(err)
		}
		return ev.SelfParent() == selfAncestor && ev.OtherParent() == ancestor
	}

	t.Run("Check Event Coordinates", func(t *testing.T) {

		e0Event, err := p.Store.GetEvent(index[e0])
		if err != nil {
			t.Fatal(err)
		}

		if !(e0Event.Message.SelfParentIndex == -1 &&
			e0Event.Message.OtherParentCreatorID == -1 &&
			e0Event.Message.OtherParentIndex == -1 &&
			e0Event.Message.CreatorID == p.Participants.ByPubKey[e0Event.Creator()].ID) {
			t.Fatalf("Invalid wire info on %s", e0)
		}

		e21Event, err := p.Store.GetEvent(index[e21])
		if err != nil {
			t.Fatal(err)
		}

		e10Event, err := p.Store.GetEvent(index[e10])
		if err != nil {
			t.Fatal(err)
		}

		if !(e21Event.Message.SelfParentIndex == 1 &&
			e21Event.Message.OtherParentCreatorID == p.Participants.ByPubKey[e10Event.Creator()].ID &&
			e21Event.Message.OtherParentIndex == 1 &&
			e21Event.Message.CreatorID == p.Participants.ByPubKey[e21Event.Creator()].ID) {
			t.Fatalf("Invalid wire info on %s", e21)
		}

		f1Event, err := p.Store.GetEvent(index[f1])
		if err != nil {
			t.Fatal(err)
		}

		if !(f1Event.Message.SelfParentIndex == 2 &&
			f1Event.Message.OtherParentCreatorID == p.Participants.ByPubKey[e0Event.Creator()].ID &&
			f1Event.Message.OtherParentIndex == 2 &&
			f1Event.Message.CreatorID == p.Participants.ByPubKey[f1Event.Creator()].ID) {
			t.Fatalf("Invalid wire info on %s", f1)
		}

		e0CreatorID := strconv.FormatInt(p.Participants.ByPubKey[e0Event.Creator()].ID, 10)

		type Hierarchy struct {
			ev, selfAncestor, ancestor string
		}

		toCheck := []Hierarchy{
			{e0, "Root" + e0CreatorID, ""},
			{e10, index[e1], index[e0]},
			{e21, index[s20], index[e10]},
			{e02, index[s00], index[e21]},
			{f1, index[s10], index[e02]},
		}

		for _, v := range toCheck {
			if !checkParents(v.ev, v.selfAncestor, v.ancestor) {
				t.Fatal(v.ev + " selfParent not good")
			}
		}
	})

	t.Run("Check UndeterminedEvents", func(t *testing.T) {

		expectedUndeterminedEvents := []string{
			index[e0],
			index[e1],
			index[e2],
			index[e10],
			index[s20],
			index[s00],
			index[e21],
			index[e02],
			index[s10],
			index[f1],
			index[s11]}

		for i, eue := range expectedUndeterminedEvents {
			if ue := p.UndeterminedEvents[i]; ue != eue {
				t.Fatalf("UndeterminedEvents[%d] should be %s, not %s",
					i, eue, ue)
			}
		}

		// Pending loaded Events
		// 3 Events with index 0,
		// 1 Event with non-empty Transactions
		// = 4 Loaded Events
		if ple := p.PendingLoadedEvents; ple != 4 {
			t.Fatalf("PendingLoadedEvents should be 4, not %d", ple)
		}
	})
}

func TestReadWireInfo(t *testing.T) {
	p, index, _ := initRoundPoset(t)

	for k, evh := range index {
		if k[0] == 'r' {
			continue
		}
		ev, err := p.Store.GetEvent(evh)
		if err != nil {
			t.Fatal(err)
		}

		evWire := ev.ToWire()

		evFromWire, err := p.ReadWireInfo(evWire)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(ev.Message.Body.BlockSignatures,
			evFromWire.Message.Body.BlockSignatures) {
			t.Fatalf("Error converting %s.Body.BlockSignatures"+
				" from light wire", k)
		}

		if !ev.Message.Body.Equals(evFromWire.Message.Body) {
			t.Fatalf("Error converting %s.Body from light wire", k)
		}

		if !reflect.DeepEqual(ev.Message.Signature,
			evFromWire.Message.Signature) {
			t.Fatalf("Error converting %s.Signature from light wire", k)
		}

		ok, err := evFromWire.Verify()
		if !ok {
			t.Fatalf("Error verifying signature for %s from ligh wire: %v",
				k, err)
		}
	}
}

func TestStronglySee(t *testing.T) {
	p, index, _ := initRoundPoset(t)

	expected := []ancestryItem{
		{e21, e0, true, false},
		{e02, e10, true, false},
		{e02, e0, true, false},
		{e02, e1, true, false},
		{f1, e21, true, false},
		{f1, e10, true, false},
		{f1, e0, true, false},
		{f1, e1, true, false},
		{f1, e2, true, false},
		{s11, e2, true, false},
		// false negatives
		{e10, e0, false, false},
		{e21, e1, false, false},
		{e21, e2, false, false},
		{e02, e2, false, false},
		{s11, e02, false, false},
		{s11, "", false, true},
		// root events
		{s11, r1, true, false},
		{e21, r0, true, false},
		{e21, r1, false, false},
		{e10, r0, false, false},
		{s20, r2, false, false},
		{e02, r2, false, false},
		{e21, r2, false, false},
	}

	for _, exp := range expected {
		a, err := p.stronglySee(index[exp.descendant], index[exp.ancestor])
		if err != nil && !exp.err {
			t.Fatalf("Error computing stronglySee(%s, %s). Err: %v",
				exp.descendant, exp.ancestor, err)
		}
		if a != exp.val {
			t.Fatalf("stronglySee(%s, %s) should be %v, not %v",
				exp.descendant, exp.ancestor, exp.val, a)
		}
	}
}

func TestWitness(t *testing.T) {
	p, index, _ := initRoundPoset(t)

	round0Witnesses := make(map[string]*RoundEvent)
	round0Witnesses[index[e0]] = &RoundEvent{
		Witness: true, Famous: Trilean_UNDEFINED}
	round0Witnesses[index[e1]] = &RoundEvent{
		Witness: true, Famous: Trilean_UNDEFINED}
	round0Witnesses[index[e2]] = &RoundEvent{
		Witness: true, Famous: Trilean_UNDEFINED}
	p.Store.SetRoundCreated(0, RoundCreated{
		Message: RoundCreatedMessage{Events: round0Witnesses}})

	round1Witnesses := make(map[string]*RoundEvent)
	round1Witnesses[index[f1]] = &RoundEvent{
		Witness: true, Famous: Trilean_UNDEFINED}
	p.Store.SetRoundCreated(1, RoundCreated{
		Message: RoundCreatedMessage{Events: round1Witnesses}})

	expected := []ancestryItem{
		{"", e0, true, false},
		{"", e1, true, false},
		{"", e2, true, false},
		{"", f1, true, false},
		{"", e10, false, false},
		{"", e21, true, false},
		{"", e02, true, false},
	}

	for _, exp := range expected {
		a, err := p.witness(index[exp.ancestor])
		if err != nil {
			t.Fatalf("Error computing witness(%s). Err: %v",
				exp.ancestor, err)
		}
		if a != exp.val {
			t.Fatalf("witness(%s) should be %v, not %v",
				exp.ancestor, exp.val, a)
		}
	}
}

func TestRound(t *testing.T) {
	p, index, _ := initRoundPoset(t)

	round0Witnesses := make(map[string]*RoundEvent)
	round0Witnesses[index[e0]] = &RoundEvent{
		Witness: true, Famous: Trilean_UNDEFINED}
	round0Witnesses[index[e1]] = &RoundEvent{
		Witness: true, Famous: Trilean_UNDEFINED}
	round0Witnesses[index[e2]] = &RoundEvent{
		Witness: true, Famous: Trilean_UNDEFINED}
	p.Store.SetRoundCreated(0, RoundCreated{Message: RoundCreatedMessage{
		Events: round0Witnesses}})

	round1Witnesses := make(map[string]*RoundEvent)
	round1Witnesses[index[e21]] = &RoundEvent{
		Witness: true, Famous: Trilean_UNDEFINED}
	round1Witnesses[index[e02]] = &RoundEvent{
		Witness: true, Famous: Trilean_UNDEFINED}
	round1Witnesses[index[f1]] = &RoundEvent{
		Witness: true, Famous: Trilean_UNDEFINED}
	p.Store.SetRoundCreated(1, RoundCreated{
		Message: RoundCreatedMessage{Events: round1Witnesses}})

	expected := []roundItem{
		{e0, 0},
		{e1, 0},
		{e2, 0},
		{s00, 0},
		{e10, 0},
		{s20, 0},
		{e21, 1},
		{e02, 1},
		{s10, 0},
		{f1, 1},
		{s11, 2},
	}

	for _, exp := range expected {
		r, err := p.round(index[exp.event])
		if err != nil {
			t.Fatalf("Error computing round(%s). Err: %v", exp.event, err)
		}
		if r != exp.round {
			t.Fatalf("round(%s) should be %v, not %v", exp.event, exp.round, r)
		}
	}
}

func TestRoundDiff(t *testing.T) {
	p, index, _ := initRoundPoset(t)

	round0Witnesses := make(map[string]*RoundEvent)
	round0Witnesses[index[e0]] = &RoundEvent{
		Witness: true, Famous: Trilean_UNDEFINED}
	round0Witnesses[index[e1]] = &RoundEvent{
		Witness: true, Famous: Trilean_UNDEFINED}
	round0Witnesses[index[e2]] = &RoundEvent{
		Witness: true, Famous: Trilean_UNDEFINED}
	p.Store.SetRoundCreated(0, RoundCreated{
		Message: RoundCreatedMessage{Events: round0Witnesses}})

	round1Witnesses := make(map[string]*RoundEvent)
	round1Witnesses[index[e21]] = &RoundEvent{
		Witness: true, Famous: Trilean_UNDEFINED}
	round1Witnesses[index[e02]] = &RoundEvent{
		Witness: true, Famous: Trilean_UNDEFINED}
	round1Witnesses[index[f1]] = &RoundEvent{
		Witness: true, Famous: Trilean_UNDEFINED}
	p.Store.SetRoundCreated(1,
		RoundCreated{Message: RoundCreatedMessage{Events: round1Witnesses}})

	if d, err := p.roundDiff(index[s11], index[e21]); d != 1 {
		if err != nil {
			t.Fatalf("RoundDiff(%s, %s) returned an error: %s", s11, e02, err)
		}
		t.Fatalf("RoundDiff(%s, %s) should be 1 not %d", s11, e02, d)
	}

	if d, err := p.roundDiff(index[f1], index[s11]); d != -1 {
		if err != nil {
			t.Fatalf("RoundDiff(%s, %s) returned an error: %s", s11, f1, err)
		}
		t.Fatalf("RoundDiff(%s, %s) should be -1 not %d", s11, f1, d)
	}
	if d, err := p.roundDiff(index[e02], index[e21]); d != 0 {
		if err != nil {
			t.Fatalf("RoundDiff(%s, %s) returned an error: %s", e20, e21, err)
		}
		t.Fatalf("RoundDiff(%s, %s) should be 0 not %d", e20, e21, d)
	}
}

func TestDivideRounds(t *testing.T) {
	p, index, _ := initRoundPoset(t)

	if err := p.DivideRounds(); err != nil {
		t.Fatal(err)
	}

	if l := p.Store.LastRound(); l != 2 {
		t.Fatalf("last round should be 2 not %d", l)
	}

	round0, err := p.Store.GetRoundCreated(0)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(round0.Witnesses()); l != 3 {
		t.Fatalf("round 0 should have 3 witnesses, not %d", l)
	}
	if !contains(round0.Witnesses(), index[e0]) {
		t.Fatalf("round 0 witnesses should contain %s", e0)
	}
	if !contains(round0.Witnesses(), index[e1]) {
		t.Fatalf("round 0 witnesses should contain %s", e1)
	}
	if !contains(round0.Witnesses(), index[e2]) {
		t.Fatalf("round 0 witnesses should contain %s", e2)
	}

	round1, err := p.Store.GetRoundCreated(1)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(round1.Witnesses()); l != 3 {
		t.Fatalf("round 1 should have 1 witness, not %d", l)
	}
	if !contains(round1.Witnesses(), index[f1]) {
		t.Fatalf("round 1 witnesses should contain %s", f1)
	}

	round2, err := p.Store.GetRoundCreated(2)
	if err != nil {
		t.Fatal(err)
	}

	if l := len(round2.Witnesses()); l != 1 {
		t.Fatalf("round 1 should have 1 witness, not %d", l)
	}

	expectedPendingRounds := []pendingRound{
		{
			Index:   0,
			Decided: false,
		},
		{
			Index:   1,
			Decided: false,
		}, {
			Index:   2,
			Decided: false,
		},
	}
	for i, pd := range p.PendingRounds {
		if !reflect.DeepEqual(*pd, expectedPendingRounds[i]) {
			t.Fatalf("pendingRounds[%d] should be %v, not %v",
				i, expectedPendingRounds[i], *pd)
		}
	}

	// [event] => {lamportTimestamp, round}
	type tr struct {
		t, r int64
	}
	expectedTimestamps := map[string]tr{
		e0:  {0, 0},
		e1:  {0, 0},
		e2:  {0, 0},
		s00: {1, 0},
		e10: {1, 0},
		s20: {1, 0},
		e21: {2, 1},
		e02: {3, 1},
		s10: {2, 0},
		f1:  {4, 1},
		s11: {5, 2},
	}

	for e, et := range expectedTimestamps {
		ev, err := p.Store.GetEvent(index[e])
		if err != nil {
			t.Fatal(err)
		}
		if r := ev.round; r == nil || *r != et.r {
			t.Fatalf("%s round should be %d, not %d", e, et.r, *r)
		}
		if ts := ev.lamportTimestamp; ts == nil || *ts != et.t {
			t.Fatalf("%s lamportTimestamp should be %d, not %d", e, et.t, *ts)
		}
	}

}

func TestCreateRoot(t *testing.T) {
	p, index, _ := initRoundPoset(t)
	p.DivideRounds()

	participants := p.Participants.ToPeerSlice()

	baseRoot := NewBaseRoot(participants[0].ID)

	expected := map[string]Root{
		e0: baseRoot,
		e02: {
			NextRound: 1,
			SelfParent: &RootEvent{Hash: index[s00],
				CreatorID: participants[0].ID, Index: 1,
				LamportTimestamp: 1, Round: 0},
			Others: map[string]*RootEvent{
				index[e02]: {Hash: index[e21], CreatorID: participants[2].ID,
					Index: 2, LamportTimestamp: 2, Round: 1},
			},
		},
		s10: {
			NextRound: 0,
			SelfParent: &RootEvent{Hash: index[e10],
				CreatorID: participants[1].ID, Index: 1,
				LamportTimestamp: 1, Round: 0},
			Others: map[string]*RootEvent{},
		},
		f1: {
			NextRound: 1,
			SelfParent: &RootEvent{Hash: index[s10],
				CreatorID: participants[1].ID, Index: 2,
				LamportTimestamp: 2, Round: 0},
			Others: map[string]*RootEvent{
				index[f1]: {Hash: index[e02], CreatorID: participants[0].ID,
					Index: 2, LamportTimestamp: 3, Round: 1},
			},
		},
	}

	for evh, expRoot := range expected {
		ev, err := p.Store.GetEvent(index[evh])
		if err != nil {
			t.Fatal(err)
		}
		root, err := p.createRoot(ev)
		if err != nil {
			t.Fatalf("Error creating %s Root: %v", evh, err)
		}
		if !reflect.DeepEqual(expRoot, root) {
			t.Fatalf("%s Root should be %+v, not %+v", evh, expRoot, root)
		}
	}

}

func contains(s []string, x string) bool {
	for _, e := range s {
		if e == x {
			return true
		}
	}
	return false
}

/*



e01  e12
 |   |  \
 e0  R1  e2
 |       |
 R0      R2

*/
func initDentedPoset(t *testing.T) (*Poset, map[string]string) {
	nodes, index, orderedEvents, participants := initPosetNodes(n)

	orderedPeers := participants.ToPeerSlice()

	for _, peer := range orderedPeers {
		index[rootSelfParent(peer.ID)] = rootSelfParent(peer.ID)
	}

	plays := []play{
		{0, 0, rootSelfParent(orderedPeers[0].ID), "", e0, nil, nil,
			[]string{}},
		{2, 0, rootSelfParent(orderedPeers[2].ID), "", e2, nil, nil,
			[]string{}},
		{0, 1, e0, "", e01, nil, nil, []string{}},
		{1, 0, rootSelfParent(orderedPeers[1].ID), e2, e12, nil, nil,
			[]string{}},
	}

	playEvents(plays, nodes, index, orderedEvents)

	poset := createPoset(t, false, orderedEvents, participants, testLogger(t))

	return poset, index
}

func TestCreateRootBis(t *testing.T) {
	p, index := initDentedPoset(t)

	participants := p.Participants.ToPeerSlice()

	root := NewBaseRootEvent(participants[1].ID)
	expected := map[string]Root{
		e12: {
			NextRound:  0,
			SelfParent: &root,
			Others: map[string]*RootEvent{
				index[e12]: {Hash: index[e2], CreatorID: participants[2].ID,
					Index: 0, LamportTimestamp: 0, Round: 0},
			},
		},
	}

	for evh, expRoot := range expected {
		ev, err := p.Store.GetEvent(index[evh])
		if err != nil {
			t.Fatal(err)
		}
		root, err := p.createRoot(ev)
		if err != nil {
			t.Fatalf("Error creating %s Root: %v", evh, err)
		}
		if !reflect.DeepEqual(expRoot, root) {
			t.Fatalf("%s Root should be %v, not %v", evh, expRoot, root)
		}
	}
}

/*

e0  e1  e2    Block (0, 1)
0   1    2
*/
func initBlockPoset(t *testing.T) (*Poset, []TestNode, map[string]string) {
	nodes, index, orderedEvents, participants := initPosetNodes(n)

	for i, peer := range participants.ToPeerSlice() {
		event := NewEvent(nil, nil, nil, []string{rootSelfParent(peer.ID), ""},
			nodes[i].Pub, 0, nil)
		nodes[i].signAndAddEvent(event, fmt.Sprintf("e%d", i),
			index, orderedEvents)
	}

	poset := NewPoset(participants, NewInmemStore(participants, cacheSize),
		nil, testLogger(t))

	//create a block and signatures manually
	block := NewBlock(0, 1, []byte("framehash"),
		[][]byte{[]byte("block tx")})
	err := poset.Store.SetBlock(block)
	if err != nil {
		t.Fatalf("error setting block. Err: %s", err)
	}

	for i, ev := range *orderedEvents {
		if err := poset.InsertEvent(ev, true); err != nil {
			fmt.Printf("error inserting event %d: %s\n", i, err)
		}
	}

	return poset, nodes, index
}

func TestInsertEventsWithBlockSignatures(t *testing.T) {
	p, nodes, index := initBlockPoset(t)

	block, err := p.Store.GetBlock(0)
	if err != nil {
		t.Fatalf("error retrieving block 0. %s", err)
	}

	blockSigs := make([]BlockSignature, n)
	for k, n := range nodes {
		blockSigs[k], err = block.Sign(n.Key)
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Run("InsertingEventsWithValidSignatures", func(t *testing.T) {

		/*
			s00 |   |
			|   |   |
			|  e10  s20
			| / |   |
			e0  e1  e2
			0   1    2
		*/
		plays := []play{
			{1, 1, e1, e0, e10, nil, []BlockSignature{blockSigs[1]},
				[]string{}},
			{2, 1, e2, "", s20, nil, []BlockSignature{blockSigs[2]},
				[]string{}},
			{0, 1, e0, "", s00, nil, []BlockSignature{blockSigs[0]},
				[]string{}},
		}

		for _, pl := range plays {
			e := NewEvent(pl.txPayload,
				nil,
				pl.sigPayload,
				[]string{index[pl.selfParent], index[pl.otherParent]},
				nodes[pl.to].Pub,
				pl.index, nil)
			e.Sign(nodes[pl.to].Key)
			index[pl.name] = e.Hex()
			if err := p.InsertEvent(e, true); err != nil {
				t.Fatalf("error inserting event %s: %s\n", pl.name, err)
			}
		}

		// Check SigPool
		if l := len(p.SigPool); l != 3 {
			t.Fatalf("block signature pool should contain 3 signatures,"+
				" not %d", l)
		}

		// Process SigPool
		p.ProcessSigPool()

		// Check that the block contains 3 signatures
		block, _ := p.Store.GetBlock(0)
		if l := len(block.Signatures); l != 2 {
			t.Fatalf("block 0 should contain 2 signatures, not %d", l)
		}

		// Check that SigPool was cleared
		if l := len(p.SigPool); l != 0 {
			t.Fatalf("block signature pool should contain 0 signatures,"+
				" not %d", l)
		}
	})

	t.Run("InsertingEventsWithSignatureOfUnknownBlock",
		func(t *testing.T) {
			// The Event should be inserted
			// The block signature is simply ignored

			block1 := NewBlock(1, 2, []byte("framehash"), [][]byte{})
			sig, _ := block1.Sign(nodes[2].Key)

			// unknown block
			unknownBlockSig := BlockSignature{
				Validator: nodes[2].Pub,
				Index:     1,
				Signature: sig.Signature,
			}
			pl := play{2, 2, s20, e10, e21, nil, []BlockSignature{unknownBlockSig},
				[]string{}}

			e := NewEvent(nil,
				nil,
				pl.sigPayload,
				[]string{index[pl.selfParent], index[pl.otherParent]},
				nodes[pl.to].Pub,
				pl.index, nil)
			e.Sign(nodes[pl.to].Key)
			index[pl.name] = e.Hex()
			if err := p.InsertEvent(e, true); err != nil {
				t.Fatalf("ERROR inserting event %s: %s", pl.name, err)
			}

			// check that the event was recorded
			_, err := p.Store.GetEvent(index[e21])
			if err != nil {
				t.Fatalf("ERROR fetching Event %s: %s", e21, err)
			}

		})

	t.Run("InsertingEventsWithBlockSignatureNotFromCreator",
		func(t *testing.T) {
			// The Event should be inserted
			// The block signature is simply ignored

			// wrong validator
			// Validator should be same as Event creator (node 0)
			key, _ := crypto.GenerateECDSAKey()
			badNode := NewTestNode(key, 666)
			badNodeSig, _ := block.Sign(badNode.Key)

			pl := play{0, 2, s00, e21, e02, nil, []BlockSignature{badNodeSig},
				[]string{}}

			e := NewEvent(nil,
				nil,
				pl.sigPayload,
				[]string{index[pl.selfParent], index[pl.otherParent]},
				nodes[pl.to].Pub,
				pl.index, nil)
			e.Sign(nodes[pl.to].Key)
			index[pl.name] = e.Hex()
			if err := p.InsertEvent(e, true); err != nil {
				t.Fatalf("ERROR inserting event %s: %s\n", pl.name, err)
			}

			// check that the signature was not appended to the block
			block, _ := p.Store.GetBlock(0)
			if l := len(block.Signatures); l > 3 {
				t.Fatalf("Block 0 should contain 3 signatures, not %d", l)
			}
		})

}

/*
                   Round 8
      [m0]  | [m2]-----------------------------
		| \ | / |  Round 7
		|  <l1> |
		|  /|   |
	  <l0>  |   |
		| \ |   |
		|   \   |
		|   | \ |
		|   | <l2>-----------------------------
		|   | / |  Round 6
		| [k10] |
		| / |   |
	  [k0]  | [k2]-----------------------------
		| \ | / |  Round 5
		| <j1>  |
		|  /|   |
	  <j0>  |   |
		| \ |   |
		|   \   |
		|   | \ |
	    |   | <j2>-----------------------------
		|   | / |  Round 4
		| [i1]  |
		| / |   |
	  [i0]  | [i2]-----------------------------
		| \ | / |  Round 3
		| <h10> |
		|  /|   |
	   h0b  |   |
		|   |   |
	  <h0>  |   |
		| \ |   |
		|   \   |
		|   | \ |
	---g0x  | <h2>----------------------------- //g0x's other-parent is f2. This situation can happen with concurrency.
	|	|   | / |  Round 2
	|	|  g10  |
	|	| / |   |
	|  [g0] | [g2]
	|	| \ | / |
	|	| [g1]  | ------------------------------
	|	|   |   |  Round 1
	|	| <f1>  |
	|  	|  /|   |
	| <f0>  |   |
	|	| \ |   |
	|	|   \   |
	|	|   | \ |
	|   |   |  f2b
	|	|   |   |
	----------<f2>------------------------------
		|   | / |  Round 0
		|  e10  |
	    | / |   |
	   [e0][e1][e2]
		0   1    2
*/
func initConsensusPoset(db bool, t testing.TB) (*Poset, map[string]string) {
	plays := []play{
		{1, 1, e1, e0, e10, nil, nil, []string{e0, e1}},
		{2, 1, e2, e10, f2, [][]byte{[]byte(f2)}, nil, []string{e0, e1, e2}},
		{2, 2, f2, "", f2b, nil, nil, []string{f2}},
		{0, 1, e0, f2b, f0, nil, nil, []string{e0, f2}},
		{1, 2, e10, f0, f1, nil, nil, []string{f2, f0, e1}},
		{1, 3, f1, "", g1, [][]byte{[]byte(g1)}, nil, []string{f2, f0, f1}},
		{0, 2, f0, g1, g0, nil, nil, []string{g1, f0}},
		{2, 3, f2b, g1, g2, nil, nil, []string{g1, f2}},
		{1, 4, g1, g0, g10, nil, nil, []string{g1, f0}},
		{0, 3, g0, f2, g0x, nil, nil, []string{g0, g1, f2b}},
		{2, 4, g2, g10, h2, nil, nil, []string{g1, g0, g2}},
		{0, 4, g0x, h2, h0, nil, nil, []string{h2, g0, g1}},
		{0, 5, h0, "", h0b, [][]byte{[]byte(h0b)}, nil, []string{h0, h2}},
		{1, 5, g10, h0b, h10, nil, nil, []string{h0, h2, g1}},
		{0, 6, h0b, h10, i0, nil, nil, []string{h10, h0, h2}},
		{2, 5, h2, h10, i2, nil, nil, []string{h10, h0, h2}},
		{1, 6, h10, i0, i1, [][]byte{[]byte(i1)}, nil, []string{i0, h10, h0, h2}},
		{2, 6, i2, i1, j2, nil, nil, []string{i1, i0, i2}},
		{0, 7, i0, j2, j0, [][]byte{[]byte(j0)}, nil, []string{i0, j2}},
		{1, 7, i1, j0, j1, nil, nil, []string{i1, i0, j0, j2}},
		{0, 8, j0, j1, k0, nil, nil, []string{j1, j0, j2}},
		{2, 7, j2, j1, k2, nil, nil, []string{j1, j0, j2}},
		{1, 8, j1, k0, k10, nil, nil, []string{j1, j0, j2, k0}},
		{2, 8, k2, k10, l2, nil, nil, []string{k0, k10, k2}},
		{0, 9, k0, l2, l0, nil, nil, []string{k0, l2}},
		{1, 9, k10, l0, l1, nil, nil, []string{l0, l2, k10, k0}},
		{0, 10, l0, l1, m0, nil, nil, []string{l1, l0, l2}},
		{2, 9, l2, l1, m2, nil, nil, []string{l1, l0, l2}},
	}

	poset, index, _, _ := initPosetFull(t, plays, db, n, testLogger(t))

	return poset, index
}

func TestDivideRoundsBis(t *testing.T) {
	p, index := initConsensusPoset(false, t)

	if err := p.DivideRounds(); err != nil {
		t.Fatal(err)
	}

	// [event] => {lamportTimestamp, round}
	type tr struct {
		t, r int64
	}
	expectedTimestamps := map[string]tr{
		e0:  {0, 0},
		e1:  {0, 0},
		e2:  {0, 0},
		e10: {1, 0},
		f2:  {2, 1},
		f2b: {3, 1},
		f0:  {4, 1},
		f1:  {5, 1},
		g1:  {6, 2},
		g0:  {7, 2},
		g2:  {7, 2},
		g10: {8, 2},
		g0x: {8, 2},
		h2:  {9, 3},
		h0:  {10, 3},
		h0b: {11, 3},
		h10: {12, 3},
		i0:  {13, 4},
		i2:  {13, 4},
		i1:  {14, 4},
		j2:  {15, 5},
		j0:  {16, 5},
		j1:  {17, 5},
		k0:  {18, 6},
		k2:  {18, 6},
		k10: {19, 6},
		l2:  {20, 7},
		l0:  {21, 7},
		l1:  {22, 7},
		m0:  {23, 8},
		m2:  {23, 8},
	}

	for e, et := range expectedTimestamps {
		ev, err := p.Store.GetEvent(index[e])
		if err != nil {
			t.Fatal(err)
		}
		if r := ev.round; r == nil || *r != et.r {
			t.Fatalf("%s round should be %d, not %d", e, et.r, *r)
		}
		if ts := ev.lamportTimestamp; ts == nil || *ts != et.t {
			t.Fatalf("%s lamportTimestamp should be %d, not %d", e, et.t, *ts)
		}
	}

}

func TestDecideFame(t *testing.T) {
	p, index := initConsensusPoset(false, t)

	p.DivideRounds()
	if err := p.DecideFame(); err != nil {
		t.Fatal(err)
	}

	round0, err := p.Store.GetRoundCreated(0)
	if err != nil {
		t.Fatal(err)
	}
	if f := round0.Message.Events[index[e0]]; !(f.Witness &&
		f.Famous == Trilean_TRUE) {
		t.Fatalf("%s should be famous; got %v", e0, f)
	}
	if f := round0.Message.Events[index[e1]]; !(f.Witness &&
		f.Famous == Trilean_TRUE) {
		t.Fatalf("%s should be famous; got %v", e1, f)
	}
	if f := round0.Message.Events[index[e2]]; !(f.Witness &&
		f.Famous == Trilean_TRUE) {
		t.Fatalf("%s should be famous; got %v", e2, f)
	}

	round1, err := p.Store.GetRoundCreated(1)
	if err != nil {
		t.Fatal(err)
	}
	if f := round1.Message.Events[index[f2]]; !(f.Witness &&
		f.Famous == Trilean_TRUE) {
		t.Fatalf("%s should be famous; got %v", f2, f)
	}
	if f := round1.Message.Events[index[f0]]; !(f.Witness &&
		f.Famous == Trilean_TRUE) {
		t.Fatalf("%s should be famous; got %v", f0, f)
	}
	if f := round1.Message.Events[index[f1]]; !(f.Witness &&
		f.Famous == Trilean_TRUE) {
		t.Fatalf("%s should be famous; got %v", f1, f)
	}

	round2, err := p.Store.GetRoundCreated(2)
	if err != nil {
		t.Fatal(err)
	}
	if f := round2.Message.Events[index[g1]]; !(f.Witness &&
		f.Famous == Trilean_TRUE) {
		t.Fatalf("%s should be famous; got %v", g1, f)
	}
	if f := round2.Message.Events[index[g0]]; !(f.Witness &&
		f.Famous == Trilean_TRUE) {
		t.Fatalf("%s should be famous; got %v", g0, f)
	}
	if f := round2.Message.Events[index[g2]]; !(f.Witness &&
		f.Famous == Trilean_TRUE) {
		t.Fatalf("%s should be famous; got %v", g2, f)
	}

	round3, err := p.Store.GetRoundCreated(3)
	if err != nil {
		t.Fatal(err)
	}
	if f := round3.Message.Events[index[h2]]; !(f.Witness &&
		f.Famous == Trilean_TRUE) {
		t.Fatalf("%s should be famous; got %v", h2, f)
	}
	if f := round3.Message.Events[index[h0]]; !(f.Witness &&
		f.Famous == Trilean_TRUE) {
		t.Fatalf("%s should be famous; got %v", h0, f)
	}
	if f := round3.Message.Events[index[h10]]; !(f.Witness &&
		f.Famous == Trilean_TRUE) {
		t.Fatalf("%s should be famous; got %v", h10, f)
	}

	round4, err := p.Store.GetRoundCreated(4)
	if err != nil {
		t.Fatal(err)
	}
	if f := round4.Message.Events[index[i0]]; !(f.Witness &&
		f.Famous == Trilean_UNDEFINED) {
		t.Fatalf("%s should be famous; got %v", i0, f)
	}
	if f := round4.Message.Events[index[i2]]; !(f.Witness &&
		f.Famous == Trilean_UNDEFINED) {
		t.Fatalf("%s should be famous; got %v", i2, f)
	}
	if f := round4.Message.Events[index[i1]]; !(f.Witness &&
		f.Famous == Trilean_UNDEFINED) {
		t.Fatalf("%s should be famous; got %v", i1, f)
	}

	expectedPendingRounds := []pendingRound{
		{Index: 0, Decided: true},
		{Index: 1, Decided: true},
		{Index: 2, Decided: true},
		{Index: 3, Decided: true},
		{Index: 4, Decided: false},
		{Index: 5, Decided: true},
		{Index: 6, Decided: false},
		{Index: 7, Decided: false},
		{Index: 8, Decided: false},
	}
	for i, pd := range p.PendingRounds {
		if !reflect.DeepEqual(*pd, expectedPendingRounds[i]) {
			t.Fatalf("pendingRounds[%d] should be %v, not %v",
				i, expectedPendingRounds[i], *pd)
		}
	}
}

func TestDecideRoundReceived(t *testing.T) {
	p, index := initConsensusPoset(false, t)

	p.DivideRounds()
	p.DecideFame()
	if err := p.DecideRoundReceived(); err != nil {
		t.Fatal(err)
	}

	for name, hash := range index {
		e, _ := p.Store.GetEvent(hash)

		switch rune(name[0]) {
		case rune('e'):
			if r := *e.roundReceived; r != 1 {
				t.Fatalf("%s round received should be 1 not %d", name, r)
			}
		case rune('f'):
			if r := *e.roundReceived; r != 2 {
				t.Fatalf("%s round received should be 2 not %d", name, r)
			}
		}
	}

	round0, err := p.Store.GetRoundCreated(0)
	if err != nil {
		t.Fatalf("could not retrieve Round 0. %s", err)
	}
	if ce := len(round0.ConsensusEvents()); ce != 0 {
		t.Fatalf("round 0 should contain 0 ConsensusEvents, not %d", ce)
	}

	round1, err := p.Store.GetRoundCreated(1)
	if err != nil {
		t.Fatalf("could not retrieve Round 1. %s", err)
	}
	if ce := len(round1.ConsensusEvents()); ce != 4 {
		t.Fatalf("round 1 should contain 4 ConsensusEvents, not %d", ce)
	}

	round2, err := p.Store.GetRoundCreated(2)
	if err != nil {
		t.Fatalf("could not retrieve Round 2. %s", err)
	}
	if ce := len(round2.ConsensusEvents()); ce != 4 {
		t.Fatalf("round 2 should contain 9 ConsensusEvents, not %d", ce)
	}

	expectedUndeterminedEvents := []string{
		index[g0x],
		index[h2],
		index[h0],
		index[h0b],
		index[h10],
		index[j2],
		index[j0],
		index[j1],
		index[k0],
		index[k2],
		index[k10],
		index[l2],
		index[l0],
		index[l1],
		index[m0],
		index[m2],
	}

	for i, eue := range expectedUndeterminedEvents {
		if ue := p.UndeterminedEvents[i]; ue != eue {
			t.Fatalf("undetermined event %d should be %s, not %s", i, eue, ue)
		}
	}
}

func TestProcessDecidedRounds(t *testing.T) {
	p, index := initConsensusPoset(false, t)

	p.DivideRounds()
	p.DecideFame()
	p.DecideRoundReceived()
	if err := p.ProcessDecidedRounds(); err != nil {
		t.Fatal(err)
	}

	consensusEvents := p.Store.ConsensusEvents()

	for i, e := range consensusEvents {
		t.Logf("consensus[%d]: %s\n", i, getName(index, e))
	}

	if l := len(consensusEvents); l != 12 {
		t.Fatalf("length of consensus should be 12 not %d", l)
	}

	if ple := p.PendingLoadedEvents; ple != 3 {
		t.Fatalf("pending loaded events number should be 3, not %d", ple)
	}

	block0, err := p.Store.GetBlock(0)
	if err != nil {
		t.Fatalf("store should contain a block with Index 0: %v", err)
	}

	if ind := block0.Index(); ind != 0 {
		t.Fatalf("block0's index should be 0, not %d", ind)
	}

	if rr := block0.RoundReceived(); rr != 2 {
		t.Fatalf("block0's round received should be 2, not %d", rr)
	}

	if l := len(block0.Transactions()); l != 1 {
		t.Fatalf("block0 should contain 1 transaction, not %d", l)
	}
	if tx := block0.Transactions()[0]; !reflect.DeepEqual(tx, []byte(f2)) {
		t.Fatalf("transaction 0 from block0 should be '%s', not %s", f2, tx)
	}

	frame1, err := p.GetFrame(block0.RoundReceived())
	frame1Hash, err := frame1.Hash()
	if !reflect.DeepEqual(block0.GetFrameHash(), frame1Hash) {
		t.Fatalf("frame hash from block0 should be %v, not %v",
			frame1Hash, block0.GetFrameHash())
	}

	block1, err := p.Store.GetBlock(1)
	if err != nil {
		t.Fatalf("store should contain a block with Index 1: %v", err)
	}

	if ind := block1.Index(); ind != 1 {
		t.Fatalf("block1's index should be 1, not %d", ind)
	}

	if rr := block1.RoundReceived(); rr != 3 {
		t.Fatalf("block1's round received should be 3, not %d", rr)
	}

	if l := len(block1.Transactions()); l != 1 {
		t.Fatalf("block1 should contain 1 transactions, not %d", l)
	}

	if tx := block1.Transactions()[0]; !reflect.DeepEqual(tx, []byte(g1)) {
		t.Fatalf("transaction 0 from block1 should be '%s', not %s", g1, tx)
	}

	frame2, err := p.GetFrame(block1.RoundReceived())
	frame2Hash, err := frame2.Hash()
	if !reflect.DeepEqual(block1.GetFrameHash(), frame2Hash) {
		t.Fatalf("frame hash from block1 should be %v, not %v",
			frame2Hash, block1.GetFrameHash())
	}

	expRounds := []pendingRound{
		{Index: 4, Decided: false},
		{Index: 5, Decided: true},
		{Index: 6, Decided: false},
		{Index: 7, Decided: false},
		{Index: 8, Decided: false},
	}
	for i, pd := range p.PendingRounds {
		if !reflect.DeepEqual(*pd, expRounds[i]) {
			t.Fatalf("pending round %d should be %v, not %v", i,
				expRounds[i], *pd)
		}
	}

	if v := p.AnchorBlock; v != nil {
		t.Fatalf("anchor block should be nil, not %v", v)
	}

}

func BenchmarkConsensus(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// we do not want to benchmark the initialization code
		b.StopTimer()
		p, _ := initConsensusPoset(false, b)
		b.StartTimer()

		p.DivideRounds()
		p.DecideFame()
		p.DecideRoundReceived()
		p.ProcessDecidedRounds()
	}
}

func TestKnown(t *testing.T) {
	p, _ := initConsensusPoset(false, t)

	participants := p.Participants.ToPeerSlice()

	expectedKnown := map[int64]int64{
		participants[0].ID: 10,
		participants[1].ID: 9,
		participants[2].ID: 9,
	}

	known := p.Store.KnownEvents()
	for i := range p.Participants.ToIDSlice() {
		if l := known[int64(i)]; l != expectedKnown[int64(i)] {
			t.Fatalf("known event %d should be %d, not %d", i,
				expectedKnown[int64(i)], l)
		}
	}
}

func TestGetFrame(t *testing.T) {
	p, index := initConsensusPoset(false, t)

	participants := p.Participants.ToPeerSlice()

	p.DivideRounds()
	p.DecideFame()
	p.DecideRoundReceived()
	p.ProcessDecidedRounds()

	t.Run("round 1", func(t *testing.T) {
		expRoots := make([]Root, n)
		expRoots[0] = NewBaseRoot(participants[0].ID)
		expRoots[1] = NewBaseRoot(participants[1].ID)
		expRoots[2] = NewBaseRoot(participants[2].ID)

		frame, err := p.GetFrame(1)
		if err != nil {
			t.Fatal(err)
		}

		for p, r := range frame.Roots {
			expRoot := expRoots[p]
			compareRootEvents(t, r.SelfParent, expRoot.SelfParent, index)
			compareOtherParents(t, r.Others, expRoot.Others, index)
		}

		var expEvents []Event

		hashes := []string{index[e0], index[e1], index[e2], index[e10]}
		for _, eh := range hashes {
			e, err := p.Store.GetEvent(eh)
			if err != nil {
				t.Fatal(err)
			}
			expEvents = append(expEvents, e)
		}

		sort.Sort(ByLamportTimestamp(expEvents))
		expEventMessages := make([]*EventMessage, len(expEvents))
		for k := range expEvents {
			expEventMessages[k] = &expEvents[k].Message
		}

		messages := frame.GetEvents()
		if len(expEventMessages) != len(messages) {
			t.Fatalf("expected number of other parents: %d, got: %d",
				len(expEventMessages), len(messages))
		}

		for k, msg := range expEventMessages {
			compareEventMessages(t, messages[k], msg, index)
		}
	})

	t.Run("round 2", func(t *testing.T) {
		expRoots := make([]Root, n)
		expRoots[0] = Root{
			NextRound: 1,
			SelfParent: &RootEvent{
				Hash:             index[e0],
				CreatorID:        participants[0].ID,
				Index:            0,
				LamportTimestamp: 0,
				Round:            0,
			},
			Others: map[string]*RootEvent{
				index[f0]: {
					Hash:             index[f2b],
					CreatorID:        participants[2].ID,
					Index:            2,
					LamportTimestamp: 3,
					Round:            1,
				},
			},
		}
		expRoots[1] = Root{
			NextRound: 1,
			SelfParent: &RootEvent{
				Hash:             index[e10],
				CreatorID:        participants[1].ID,
				Index:            1,
				LamportTimestamp: 1,
				Round:            0,
			},
			Others: map[string]*RootEvent{
				index[f1]: {
					Hash:             index[f0],
					CreatorID:        participants[0].ID,
					Index:            1,
					LamportTimestamp: 4,
					Round:            1,
				},
			},
		}
		expRoots[2] = Root{
			NextRound: 1,
			SelfParent: &RootEvent{
				Hash:             index[e2],
				CreatorID:        participants[2].ID,
				Index:            0,
				LamportTimestamp: 0,
				Round:            0,
			},
			Others: map[string]*RootEvent{
				index[f2]: {
					Hash:             index[e10],
					CreatorID:        participants[1].ID,
					Index:            1,
					LamportTimestamp: 1,
					Round:            0,
				},
			},
		}

		frame, err := p.GetFrame(2)
		if err != nil {
			t.Fatal(err)
		}

		for p, r := range frame.Roots {
			expRoot := expRoots[p]
			compareRootEvents(t, r.SelfParent, expRoot.SelfParent, index)
			compareOtherParents(t, r.Others, expRoot.Others, index)
		}

		expectedEventsHashes := []string{
			index[f2],
			index[f2b],
			index[f0],
			index[f1],
		}
		var expEvents []Event
		for _, eh := range expectedEventsHashes {
			e, err := p.Store.GetEvent(eh)
			if err != nil {
				t.Fatal(err)
			}
			expEvents = append(expEvents, e)
		}
		sort.Sort(ByLamportTimestamp(expEvents))
		expEventMessages := make([]*EventMessage, len(expEvents))
		for k := range expEvents {
			expEventMessages[k] = &expEvents[k].Message
		}

		messages := frame.GetEvents()
		if len(expEventMessages) != len(messages) {
			t.Fatalf("expected number of other parents: %d, got: %d",
				len(expEventMessages), len(messages))
		}

		for k, msg := range expEventMessages {
			compareEventMessages(t, messages[k], msg, index)
		}

		block0, err := p.Store.GetBlock(0)
		if err != nil {
			t.Fatalf("store should contain a block with Index 0: %v", err)
		}

		frameHash, err := frame.Hash()
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(block0.GetFrameHash(), frameHash) {
			t.Fatalf("frame hash (0x%X) from block 0 and frame hash"+
				" (0x%X) differ", block0.GetFrameHash(), frameHash)
		}
	})

}

func TestResetFromFrame(t *testing.T) {
	p, index := initConsensusPoset(false, t)

	participants := p.Participants.ToPeerSlice()

	p.DivideRounds()
	p.DecideFame()
	p.DecideRoundReceived()
	p.ProcessDecidedRounds()

	block, err := p.Store.GetBlock(1)
	if err != nil {
		t.Fatal(err)
	}

	frame, err := p.GetFrame(block.RoundReceived())
	if err != nil {
		t.Fatal(err)
	}

	// This operation clears the private fields which need to be recomputed
	// in the Events (round, roundReceived,etc)
	marshalledFrame, _ := frame.ProtoMarshal()
	unmarshalledFrame := new(Frame)
	unmarshalledFrame.ProtoUnmarshal(marshalledFrame)

	p2 := NewPoset(p.Participants,
		NewInmemStore(p.Participants, cacheSize),
		nil,
		testLogger(t))
	err = p2.Reset(block, *unmarshalledFrame)
	if err != nil {
		t.Fatal(err)
	}

	/*
		The poset should now look like this:

		   |   |  f10  |
		   |   | / |   |
		   |   f0  |   f2
		   |   | \ | / |
		   |   |  f1b  |
		   |   |   |   |
		   |   |   f1  |
		   |   |   |   |
		   +-- R0  R1  R2
	*/

	// Test Known
	expectedKnown := map[int64]int64{
		participants[0].ID: 2,
		participants[1].ID: 4,
		participants[2].ID: 3,
	}

	known := p2.Store.KnownEvents()
	for _, peer := range p2.Participants.ById {
		if l := known[peer.ID]; l != expectedKnown[peer.ID] {
			t.Fatalf("Known[%d] should be %d, not %d",
				peer.ID, expectedKnown[peer.ID], l)
		}
	}

	t.Run("TestDivideRounds", func(t *testing.T) {
		if err := p2.DivideRounds(); err != nil {
			t.Fatal(err)
		}

		pRound1, err := p.Store.GetRoundCreated(2)
		if err != nil {
			t.Fatal(err)
		}
		p2Round1, err := p2.Store.GetRoundCreated(2)
		if err != nil {
			t.Fatal(err)
		}

		// Check round 1 witnesses
		pWitnesses := pRound1.Witnesses()
		p2Witnesses := p2Round1.Witnesses()
		sort.Strings(pWitnesses)
		sort.Strings(p2Witnesses)
		if !reflect.DeepEqual(pWitnesses, p2Witnesses) {
			t.Fatalf("Reset Hg Round 1 witnesses should be %v, not %v",
				pWitnesses, p2Witnesses)
		}

		// check event rounds and lamport timestamps
		for _, em := range frame.Events {
			e := em.ToEvent()
			ev := &e
			p2r, err := p2.round(ev.Hex())
			if err != nil {
				t.Fatalf("Error computing %s Round: %d",
					getName(index, ev.Hex()), p2r)
			}
			hr, _ := p.round(ev.Hex())
			if p2r != hr {

				t.Fatalf("p2[%v].Round should be %d, not %d",
					getName(index, ev.Hex()), hr, p2r)
			}

			p2s, err := p2.lamportTimestamp(ev.Hex())
			if err != nil {
				t.Fatalf("Error computing %s LamportTimestamp: %d",
					getName(index, ev.Hex()), p2s)
			}
			hs, _ := p.lamportTimestamp(ev.Hex())
			if p2s != hs {
				t.Fatalf("p2[%v].LamportTimestamp should be %d, not %d",
					getName(index, ev.Hex()), hs, p2s)
			}
		}
	})

	t.Run("TestConsensus", func(t *testing.T) {
		p2.DecideFame()
		p2.DecideRoundReceived()
		p2.ProcessDecidedRounds()

		if lbi := p2.Store.LastBlockIndex(); lbi != block.Index() {
			t.Fatalf("LastBlockIndex should be %d, not %d",
				block.Index(), lbi)
		}

		if r := p2.LastConsensusRound; r == nil || *r != block.RoundReceived() {
			t.Fatalf("LastConsensusRound should be %d, not %d",
				block.RoundReceived(), *r)
		}

		if v := p2.AnchorBlock; v != nil {
			t.Fatalf("AnchorBlock should be nil, not %v", v)
		}
	})

	t.Run("TestContinueAfterReset", func(t *testing.T) {
		// Insert remaining Events into the Reset poset
		for r := int64(2); r <= int64(2); r++ {
			round, err := p.Store.GetRoundCreated(r)
			if err != nil {
				t.Fatal(err)
			}

			var events []Event
			for _, e := range round.RoundEvents() {
				ev, err := p.Store.GetEvent(e)
				if err != nil {
					t.Fatal(err)
				}
				events = append(events, ev)
			}

			sort.Sort(ByTopologicalOrder(events))

			for _, ev := range events {

				marshalledEv, _ := ev.ProtoMarshal()
				unmarshalledEv := new(Event)
				unmarshalledEv.ProtoUnmarshal(marshalledEv)
				p2.InsertEvent(*unmarshalledEv, true)
			}
		}

		p2.DivideRounds()
		p2.DecideFame()
		p2.DecideRoundReceived()
		p2.ProcessDecidedRounds()

		for r := int64(2); r <= 2; r++ {
			pRound, err := p.Store.GetRoundCreated(r)
			if err != nil {
				t.Fatal(err)
			}
			p2Round, err := p2.Store.GetRoundCreated(r)
			if err != nil {
				t.Fatal(err)
			}

			pWitnesses := pRound.Witnesses()
			p2Witnesses := p2Round.Witnesses()
			sort.Strings(pWitnesses)
			sort.Strings(p2Witnesses)

			if !reflect.DeepEqual(pWitnesses, p2Witnesses) {
				t.Fatalf("Reset Hg Round %d witnesses should be %v, not %v",
					r, pWitnesses, p2Witnesses)
			}
		}
	})
}

func TestBootstrap(t *testing.T) {

	// Initialize a first Poset with a DB backend
	// Add events and run consensus methods on it
	p, _ := initConsensusPoset(true, t)
	p.DivideRounds()
	p.DecideFame()
	p.DecideRoundReceived()
	p.ProcessDecidedRounds()

	p.Store.Close()
	defer os.RemoveAll(badgerDir)

	// Now we want to create a new Poset based on the database of the previous
	// Poset and see if we can boostrap it to the same state.
	recycledStore, err := LoadBadgerStore(cacheSize, badgerDir)
	np := NewPoset(recycledStore.participants,
		recycledStore,
		nil,
		logrus.New().WithField("id", "bootstrapped"))
	err = np.Bootstrap()
	if err != nil {
		t.Fatal(err)
	}

	hConsensusEvents := p.Store.ConsensusEvents()
	nhConsensusEvents := np.Store.ConsensusEvents()
	if len(hConsensusEvents) != len(nhConsensusEvents) {
		t.Fatalf("Bootstrapped poset should contain %d consensus events,"+
			"not %d", len(hConsensusEvents), len(nhConsensusEvents))
	}

	hKnown := p.Store.KnownEvents()
	nhKnown := np.Store.KnownEvents()
	if !reflect.DeepEqual(hKnown, nhKnown) {
		t.Fatalf("Bootstrapped poset's Known should be %#v, not %#v",
			hKnown, nhKnown)
	}

	if *p.LastConsensusRound != *np.LastConsensusRound {
		t.Fatalf("Bootstrapped poset's LastConsensusRound should be %#v,"+
			" not %#v", *p.LastConsensusRound, *np.LastConsensusRound)
	}

	if p.LastCommitedRoundEvents != np.LastCommitedRoundEvents {
		t.Fatalf("Bootstrapped poset's LastCommitedRoundEvents should be %#v,"+
			" not %#v", p.LastCommitedRoundEvents, np.LastCommitedRoundEvents)
	}

	if p.ConsensusTransactions != np.ConsensusTransactions {
		t.Fatalf("Bootstrapped poset's ConsensusTransactions should be %#v,"+
			" not %#v", p.ConsensusTransactions, np.ConsensusTransactions)
	}

	if p.PendingLoadedEvents != np.PendingLoadedEvents {
		t.Fatalf("Bootstrapped poset's PendingLoadedEvents should be %#v,"+
			" not %#v", p.PendingLoadedEvents, np.PendingLoadedEvents)
	}
}

/*

	|   <w51> |    |
    |    |  \ |    |
	|    |   <e23> |
	|	 |    |	\  |	   	ROUND 7
	|    |    |  <w43>----------------------
	|    |    | /  | 		ROUND 6
    |    |  [w42]  |
    |    | /  |    |
    |  [w41]  |    |
	| /  |    |    |
  [w40]  |    |    |------------------------
    | \  |    |    |		ROUND 5
    |  <d13>  |    |
    |    |  \ |    |
  <w30>  |    \    |
    | \  |    | \  |
    |   \     |  <w33>----------------------
    |    | \  |  / |		ROUND 4
    |    |  [w32]  |
    |    |  / |    |
	|  [w31]  |    |
    |  / |    |    |
   [w20] |    |    |------------------------
    |  \ |    |    | 		ROUND 3
    |    \    |    |
    |    | \  |    |
    |    |  <w22>  |
    |    | /  |    |
    |   c10   |    |
    | /  |    |    |
  <b00><w21>  |    |------------------------
    |    |  \ |    |		ROUND 2
    |    |    \    |
    |    |    | \  |
    |    |    |  [w23]
    |    |    | /  |
   [w10] |   b21   |
	| \  | /  |    |
    |  [w11]  |    |
    |    |  \ |    |
	|    |  [w12]  |------------------------
    |    |    | \  |		ROUND 1
    |    |    |  <w13>
    |    |    | /  |
    |   a10 <a21>  |
    |  / |  / |    |
    |/ <a12>  |    |------------------------
   a00   |  \ |    |		ROUND 0
	|    |   a23   |
    |    |    | \  |
  [w00][w01][w02][w03]
	0	 1	  2	   3
*/

func initFunkyPoset(t *testing.T, logger *logrus.Logger, full bool) (*Poset, map[string]string) {
	nodes, index, orderedEvents, participants := initPosetNodes(4)

	for i, peer := range participants.ToPeerSlice() {
		name := fmt.Sprintf("w0%d", i)
		event := NewEvent([][]byte{[]byte(name)}, nil,
			nil, []string{rootSelfParent(peer.ID), ""}, nodes[i].Pub, 0,
			map[string]int64{rootSelfParent(peer.ID): 1})
		nodes[i].signAndAddEvent(event, name, index, orderedEvents)
	}

	plays := []play{
		{2, 1, w02, w03, a23, [][]byte{[]byte(a23)},
			nil, []string{w02, w03}},
		{1, 1, w01, a23, a12, [][]byte{[]byte(a12)},
			nil, []string{w01, w02, w03}},
		{0, 1, w00, "", a00, [][]byte{[]byte(a00)},
			nil, []string{w00}},
		{1, 2, a12, a00, a10, [][]byte{[]byte(a10)},
			nil, []string{w00, a12}},
		{2, 2, a23, a12, a21, [][]byte{[]byte(a21)},
			nil, []string{a12, w02, w03}},
		{3, 1, w03, a21, w13, [][]byte{[]byte(w13)},
			nil, []string{a12, a21, w03}},
		{2, 3, a21, w13, w12, [][]byte{[]byte(w12)},
			nil, []string{a12, a21, w13}},
		{1, 3, a10, w12, w11, [][]byte{[]byte(w11)},
			nil, []string{w12, a12}},
		{0, 2, a00, w11, w10, [][]byte{[]byte(w10)},
			nil, []string{w11, w12, w00}},
		{2, 4, w12, w11, b21, [][]byte{[]byte(b21)},
			nil, []string{w11, w12}},
		{3, 2, w13, b21, w23, [][]byte{[]byte(w23)},
			nil, []string{w11, w12, w13}},
		{1, 4, w11, w23, w21, [][]byte{[]byte(w21)},
			nil, []string{w11, w12, w23}},
		{0, 3, w10, "", b00, [][]byte{[]byte(b00)},
			nil, []string{w10, w11, w12}},
		{1, 5, w21, b00, c10, [][]byte{[]byte(c10)},
			nil, []string{b00, w21}},
		{2, 5, b21, c10, w22, [][]byte{[]byte(w22)},
			nil, []string{b00, w21, w11, w12}},
		{0, 4, b00, w22, w20, [][]byte{[]byte(w20)},
			nil, []string{b00, w21, w22}},
		{1, 6, c10, w20, w31, [][]byte{[]byte(w31)},
			nil, []string{w20, b00, w21}},
		{2, 6, w22, w31, w32, [][]byte{[]byte(w32)},
			nil, []string{w31, w20, w22, b00, w21}},
		{0, 5, w20, w32, w30, [][]byte{[]byte(w30)},
			nil, []string{w32, w31, w20}},
		{3, 3, w23, w32, w33, [][]byte{[]byte(w33)},
			nil, []string{w23, w11, w12, w32, w31, w20}},
		{1, 7, w31, w33, d13, [][]byte{[]byte(d13)},
			nil, []string{w33, w31, w20}},
		{0, 6, w30, d13, w40, [][]byte{[]byte(w40)},
			nil, []string{w30, d13, w33}},
		{1, 8, d13, w40, w41, [][]byte{[]byte(w41)},
			nil, []string{w40, d13, w33}},
		{2, 7, w32, w41, w42, [][]byte{[]byte(w42)},
			nil, []string{w41, w40, w32, w31, w20}},
		{3, 4, w33, w42, w43, [][]byte{[]byte(w43)},
			nil, []string{w42, w41, w40, w33}},
	}
	if full {
		newPlays := []play{
			{2, 8, w42, w43, e23, [][]byte{[]byte(e23)},
				nil, []string{w43, w42, w41, w40}},
			{1, 9, w41, e23, w51, [][]byte{[]byte(w51)},
				nil, []string{e23, w43, w41, w40}},
		}
		plays = append(plays, newPlays...)
	}

	playEvents(plays, nodes, index, orderedEvents)

	poset := createPoset(t, false, orderedEvents, participants, logger.WithField("test", 6))

	return poset, index
}

func TestFunkyPosetFame(t *testing.T) {
	p, index := initFunkyPoset(t, common.NewTestLogger(t), false)

	if err := p.DivideRounds(); err != nil {
		t.Fatal(err)
	}
	if err := p.DecideFame(); err != nil {
		t.Fatal(err)
	}

	l := p.Store.LastRound()
	if l != 7 {
		t.Fatalf("last round should be 7 not %d", l)
	}

	for r := int64(0); r < l+1; r++ {
		round, err := p.Store.GetRoundCreated(r)
		if err != nil {
			t.Fatal(err)
		}
		var witnessNames []string
		for _, w := range round.Witnesses() {
			witnessNames = append(witnessNames, getName(index, w))
		}
		t.Logf("round %d witnesses: %v", r, witnessNames)
	}

	expPendingRounds := []pendingRound{
		{Index: 0, Decided: true},
		{Index: 1, Decided: true},
		{Index: 2, Decided: true},
		{Index: 3, Decided: true},
		{Index: 4, Decided: true},
		{Index: 5, Decided: false},
		{Index: 6, Decided: false},
		{Index: 7, Decided: false},
	}

	for i, pd := range p.PendingRounds {
		if !reflect.DeepEqual(*pd, expPendingRounds[i]) {
			t.Fatalf("pending round %d should be %v, not %v", i,
				expPendingRounds[i], *pd)
		}
	}

	if err := p.DecideRoundReceived(); err != nil {
		t.Fatal(err)
	}
	if err := p.ProcessDecidedRounds(); err != nil {
		t.Fatal(err)
	}

	remainingPendingRounds := expPendingRounds[5:]
	for i := 0; i < len(p.PendingRounds); i++ {
		if !reflect.DeepEqual(*p.PendingRounds[i], remainingPendingRounds[i]) {
			t.Fatalf("remaining pending round %d should be %v, not %v", i,
				remainingPendingRounds[i], *p.PendingRounds[i])
		}
	}
}

func TestFunkyPosetBlocks(t *testing.T) {
	p, index := initFunkyPoset(t, common.NewTestLogger(t), true)

	if err := p.DivideRounds(); err != nil {
		t.Fatal(err)
	}
	if err := p.DecideFame(); err != nil {
		t.Fatal(err)
	}
	if err := p.DecideRoundReceived(); err != nil {
		t.Fatal(err)
	}
	if err := p.ProcessDecidedRounds(); err != nil {
		t.Fatal(err)
	}

	l := p.Store.LastRound()
	if l != 7 {
		t.Fatalf("last round should be 7 not %d", l)
	}

	for r := int64(0); r < l+1; r++ {
		round, err := p.Store.GetRoundCreated(r)
		if err != nil {
			t.Fatal(err)
		}
		var witnessNames []string
		for _, w := range round.Witnesses() {
			witnessNames = append(witnessNames, getName(index, w))
		}
		t.Logf("round %d witnesses: %v", r, witnessNames)
	}

	// Rounds 0,1,2,3,4 and 5 should be decided.
	expPendingRounds := []pendingRound{
		{Index: 6, Decided: false},
		{Index: 7, Decided: false},
	}
	for i, pd := range p.PendingRounds {
		if !reflect.DeepEqual(*pd, expPendingRounds[i]) {
			t.Fatalf("pending round %d should be %v, not %v",
				i, expPendingRounds[i], *pd)
		}
	}

	expBlockTxCounts := map[int64]int64{0: 4, 1: 3, 2: 5, 3: 7, 4: 3}

	for bi := int64(0); bi < 5; bi++ {
		b, err := p.Store.GetBlock(bi)
		if err != nil {
			t.Fatal(err)
		}
		for i, tx := range b.Transactions() {
			t.Logf("block %d, tx %d: %s", bi, i, string(tx))
		}
		if txs := int64(len(b.Transactions())); txs != expBlockTxCounts[bi] {
			t.Fatalf("Blocks[%d] should contain %d transactions, not %d", bi,
				expBlockTxCounts[bi], txs)
		}
	}
}

func TestFunkyPosetFrames(t *testing.T) {
	p, index := initFunkyPoset(t, common.NewTestLogger(t), true)

	participants := p.Participants.ToPeerSlice()

	if err := p.DivideRounds(); err != nil {
		t.Fatal(err)
	}
	if err := p.DecideFame(); err != nil {
		t.Fatal(err)
	}
	if err := p.DecideRoundReceived(); err != nil {
		t.Fatal(err)
	}
	if err := p.ProcessDecidedRounds(); err != nil {
		t.Fatal(err)
	}

	for bi := int64(0); bi < 5; bi++ {
		block, err := p.Store.GetBlock(bi)
		if err != nil {
			t.Fatal(err)
		}

		frame, err := p.GetFrame(block.RoundReceived())
		for k, em := range frame.Events {
			e := em.ToEvent()
			ev := &e
			r, _ := p.round(ev.Hex())
			t.Logf("frame %d events %d: %s, round %d",
				frame.Round, k, getName(index, ev.Hex()), r)
		}
		for k, r := range frame.Roots {
			t.Logf("frame %d root %d: next round %d, self parent: %v,"+
				" others: %v", frame.Round, k, r.NextRound,
				r.SelfParent, r.Others)
		}
	}

	expFrameRoots := map[int64][]Root{
		1: {
			NewBaseRoot(participants[0].ID),
			NewBaseRoot(participants[1].ID),
			NewBaseRoot(participants[2].ID),
			NewBaseRoot(participants[3].ID),
		},
		2: {
			NewBaseRoot(participants[0].ID),
			{
				NextRound: 1,
				SelfParent: &RootEvent{Hash: index[w01],
					CreatorID: participants[1].ID, Index: 0,
					LamportTimestamp: 0, Round: 0},
				Others: map[string]*RootEvent{
					index[a12]: {Hash: index[a23],
						CreatorID: participants[2].ID, Index: 1,
						LamportTimestamp: 1, Round: 0},
				},
			},
			{
				NextRound: 1,
				SelfParent: &RootEvent{Hash: index[a23],
					CreatorID: participants[2].ID, Index: 1,
					LamportTimestamp: 1, Round: 0},
				Others: map[string]*RootEvent{
					index[a21]: {Hash: index[a12],
						CreatorID: participants[1].ID, Index: 1,
						LamportTimestamp: 2, Round: 1},
				},
			},
			{
				NextRound: 1,
				SelfParent: &RootEvent{Hash: index[w03],
					CreatorID: participants[3].ID, Index: 0,
					LamportTimestamp: 0, Round: 0},
				Others: map[string]*RootEvent{
					index[w13]: {Hash: index[a21],
						CreatorID: participants[2].ID, Index: 2,
						LamportTimestamp: 3, Round: 1},
				},
			},
		},
		3: {
			NewBaseRoot(participants[0].ID),
			{
				NextRound: 1,
				SelfParent: &RootEvent{Hash: index[a12],
					CreatorID: participants[1].ID, Index: 1,
					LamportTimestamp: 2, Round: 1},
				Others: map[string]*RootEvent{
					index[a10]: {Hash: index[a00],
						CreatorID: participants[0].ID, Index: 1,
						LamportTimestamp: 1, Round: 0},
				},
			},
			{
				NextRound: 2,
				SelfParent: &RootEvent{Hash: index[a21],
					CreatorID: participants[2].ID, Index: 2,
					LamportTimestamp: 3, Round: 1},
				Others: map[string]*RootEvent{
					index[w12]: {Hash: index[w13],
						CreatorID: participants[3].ID, Index: 1,
						LamportTimestamp: 4, Round: 1},
				},
			},
			{
				NextRound: 1,
				SelfParent: &RootEvent{Hash: index[w03],
					CreatorID: participants[3].ID, Index: 0,
					LamportTimestamp: 0, Round: 0},
				Others: map[string]*RootEvent{
					index[w13]: {Hash: index[a21],
						CreatorID: participants[2].ID, Index: 2,
						LamportTimestamp: 3, Round: 1},
				},
			},
		},
		4: {
			{
				NextRound: 2,
				SelfParent: &RootEvent{Hash: index[a00],
					CreatorID: participants[0].ID, Index: 1,
					LamportTimestamp: 1, Round: 0},
				Others: map[string]*RootEvent{
					index[w10]: {Hash: index[w11],
						CreatorID: participants[1].ID, Index: 3,
						LamportTimestamp: 6, Round: 2},
				},
			},
			{
				NextRound: 3,
				SelfParent: &RootEvent{Hash: index[w11],
					CreatorID: participants[1].ID, Index: 3,
					LamportTimestamp: 6, Round: 2},
				Others: map[string]*RootEvent{
					index[w21]: {Hash: index[w23],
						CreatorID: participants[3].ID, Index: 2,
						LamportTimestamp: 8, Round: 2},
				},
			},
			{
				NextRound: 2,
				SelfParent: &RootEvent{Hash: index[w12],
					CreatorID: participants[2].ID, Index: 3,
					LamportTimestamp: 5, Round: 2},
				Others: map[string]*RootEvent{
					index[b21]: {Hash: index[w11],
						CreatorID: participants[1].ID, Index: 3,
						LamportTimestamp: 6, Round: 2},
				},
			},
			{
				NextRound: 2,
				SelfParent: &RootEvent{Hash: index[w13],
					CreatorID: participants[3].ID, Index: 1,
					LamportTimestamp: 4, Round: 1},
				Others: map[string]*RootEvent{
					index[w23]: {Hash: index[b21],
						CreatorID: participants[2].ID, Index: 4,
						LamportTimestamp: 7, Round: 2},
				},
			},
		},
		5: {
			{
				NextRound: 4,
				SelfParent: &RootEvent{Hash: index[b00],
					CreatorID: participants[0].ID, Index: 3,
					LamportTimestamp: 8, Round: 3},
				Others: map[string]*RootEvent{
					index[w20]: {Hash: index[w22],
						CreatorID: participants[2].ID, Index: 5,
						LamportTimestamp: 11, Round: 3},
				},
			},
			{
				NextRound: 4,
				SelfParent: &RootEvent{Hash: index[c10],
					CreatorID: participants[1].ID, Index: 5,
					LamportTimestamp: 10, Round: 3},
				Others: map[string]*RootEvent{
					index[w31]: {Hash: index[w20],
						CreatorID: participants[0].ID, Index: 4,
						LamportTimestamp: 12, Round: 4},
				},
			},
			{
				NextRound: 4,
				SelfParent: &RootEvent{Hash: index[w22],
					CreatorID: participants[2].ID, Index: 5,
					LamportTimestamp: 11, Round: 3},
				Others: map[string]*RootEvent{
					index[w32]: {Hash: index[w31],
						CreatorID: participants[1].ID, Index: 6,
						LamportTimestamp: 13, Round: 4},
				},
			},
			{
				NextRound: 2,
				SelfParent: &RootEvent{Hash: index[w13],
					CreatorID: participants[3].ID, Index: 1,
					LamportTimestamp: 4, Round: 1},
				Others: map[string]*RootEvent{
					index[w23]: {Hash: index[b21],
						CreatorID: participants[2].ID, Index: 4,
						LamportTimestamp: 7, Round: 2},
				},
			},
		},
	}

	for bi := int64(0); bi < 5; bi++ {
		block, err := p.Store.GetBlock(bi)
		if err != nil {
			t.Fatal(err)
		}

		frame, err := p.GetFrame(block.RoundReceived())
		if err != nil {
			t.Fatal(err)
		}

		for k, r := range frame.Roots {
			compareRoots(t, r, &expFrameRoots[frame.Round][k], index)
		}
	}
}

func TestFunkyPosetReset(t *testing.T) {
	p, index := initFunkyPoset(t, common.NewTestLogger(t), true)

	p.DivideRounds()
	p.DecideFame()
	p.DecideRoundReceived()
	p.ProcessDecidedRounds()

	for bi := int64(0); bi < 3; bi++ {
		block, err := p.Store.GetBlock(bi)
		if err != nil {
			t.Fatal(err)
		}

		frame, err := p.GetFrame(block.RoundReceived())
		if err != nil {
			t.Fatal(err)
		}

		// This operation clears the private fields which need to be recomputed
		// in the Events (round, roundReceived,etc)
		marshalledFrame, _ := frame.ProtoMarshal()
		unmarshalledFrame := new(Frame)
		unmarshalledFrame.ProtoUnmarshal(marshalledFrame)

		p2 := NewPoset(p.Participants,
			NewInmemStore(p.Participants, cacheSize),
			nil,
			testLogger(t))
		err = p2.Reset(block, *unmarshalledFrame)
		if err != nil {
			t.Fatal(err)
		}

		// Test continue after reset
		// Compute diff
		p2Known := p2.Store.KnownEvents()
		diff := getDiff(p, p2Known, t)

		wireDiff := make([]WireEvent, len(diff), len(diff))
		for i, e := range diff {
			wireDiff[i] = e.ToWire()
		}

		// Insert remaining Events into the Reset poset
		for i, wev := range wireDiff {
			ev, err := p2.ReadWireInfo(wev)
			if err != nil {
				t.Fatalf("Reading WireInfo for %s: %s",
					getName(index, diff[i].Hex()), err)
			}
			err = p2.InsertEvent(*ev, false)
			if err != nil {
				t.Fatal(err)
			}
		}

		p2.DivideRounds()
		p2.DecideFame()
		p2.DecideRoundReceived()
		p2.ProcessDecidedRounds()

		compareRoundWitnesses(p, p2, index, bi, true, t)
	}

}

/*

    |  <w51>  |    |
    |    |    \    |
	|    |    | \  |  7
    |    |    |  <i32>-
    |    |    | /  |  6
    |    |  [w42]  |
    |    |  / |    |
    |  [w41]  |    |
	|    |   \     |
	|    |    | \  |  6
    |    |    |  [w43]-
	|    |    | /  |  5
    |    | <h21>   |
    |    | /  |    |
    |  <w31>  |    |
	|    |   \     |
	|    |    | \  |  5
    |    |    |  <w33>-
    |    |    | /  |  4
    |    |  [w32]  |
	|    | /  |    |
    |  [g13]  |    |
	|    |   \     |
	|    |    | \  |  4
    |    |    |  [w23]-
    |    |    | /  |  3
    |    |  <w22>  |
    |    | /  |    |
    |  <w21>  |    |
	|	 |	 \	   |
	|    |      \  |  3
    |    |    |  <w13>-
    |    |    | /  |  2
    |    |  [w12]  |
    |     /   |    |
	|  / |    |    |
  [f01]  |    |    |
	| \  |    |    |  2
    |  [w11]  |    |-
    | /  |    |    |  1
  <w10>  |    |    |
    |    \    |    |
    |    |    \    |
    |    |    |  <e32>
    |    |    | /  |  1
    |    |  <e21>  |-
    |    | /  |    |  0
    |   e10   |    |
    |  / |    |    |
  [w00][w01][w02][w03]
	|    |    |    |
    R0   R1   R2   R3
	0	 1	  2	   3
*/

func initSparsePoset(
	t *testing.T, logger *logrus.Logger) (*Poset, map[string]string) {
	nodes, index, orderedEvents, participants := initPosetNodes(4)

	for i, peer := range participants.ToPeerSlice() {
		name := fmt.Sprintf("w0%d", i)
		event := NewEvent([][]byte{[]byte(name)}, nil,
			nil, []string{rootSelfParent(peer.ID), ""}, nodes[i].Pub, 0,
			map[string]int64{rootSelfParent(peer.ID): 1})
		nodes[i].signAndAddEvent(event, name, index, orderedEvents)
	}

	plays := []play{
		{1, 1, w01, w00, e10, [][]byte{[]byte(e10)},
			nil, []string{w00, w01}},
		{2, 1, w02, e10, e21, [][]byte{[]byte(e21)},
			nil, []string{w00, w01, w02}},
		{3, 1, w03, e21, e32, [][]byte{[]byte(e32)},
			nil, []string{e21, w03}},
		{0, 1, w00, e32, w10, [][]byte{[]byte(w10)},
			nil, []string{e21, e32, w00}},
		{1, 2, e10, w10, w11, [][]byte{[]byte(w11)},
			nil, []string{w10, e32, e21, w01, w00}},
		{0, 2, w10, w11, f01, [][]byte{[]byte(f01)},
			nil, []string{w11, w10, e32, e21}},
		{2, 2, e21, f01, w12, [][]byte{[]byte(w12)},
			nil, []string{f01, w11, e21}},
		{3, 2, e32, w12, w13, [][]byte{[]byte(w13)},
			nil, []string{w12, f01, w11, e32, e21}},
		{1, 3, w11, w13, w21, [][]byte{[]byte(w21)},
			nil, []string{w13, w11}},
		{2, 3, w12, w21, w22, [][]byte{[]byte(w22)},
			nil, []string{w21, w13, w12, f01, w11}},
		{3, 3, w13, w22, w23, [][]byte{[]byte(w23)},
			nil, []string{w22, w21, w13}},
		{1, 4, w21, w23, g13, [][]byte{[]byte(g13)},
			nil, []string{w23, w21, w13}},
		{2, 4, w22, g13, w32, [][]byte{[]byte(w32)},
			nil, []string{g13, w23, w22, w21, w13}},
		{3, 4, w23, w32, w33, [][]byte{[]byte(w33)},
			nil, []string{w32, g13, w23}},
		{1, 5, g13, w33, w31, [][]byte{[]byte(w31)},
			nil, []string{w33, g13, w23}},
		{2, 5, w32, w31, h21, [][]byte{[]byte(h21)},
			nil, []string{w31, w33, w32, g13, w23}},
		{3, 5, w33, h21, w43, [][]byte{[]byte(w43)},
			nil, []string{h21, w31, w33}},
		{1, 6, w31, w43, w41, [][]byte{[]byte(w41)},
			nil, []string{w43, w31, w33}},
		{2, 6, h21, w41, w42, [][]byte{[]byte(w42)},
			nil, []string{w41, w43, h21, w31, w33}},
		{3, 6, w43, w42, i32, [][]byte{[]byte(i32)},
			nil, []string{w42, w41, w43}},
		{1, 7, w41, i32, w51, [][]byte{[]byte(w51)},
			nil, []string{i32, w41, w43}},
	}

	playEvents(plays, nodes, index, orderedEvents)

	poset := createPoset(t, false, orderedEvents, participants,
		logger.WithField("test", 6))

	return poset, index
}

func TestSparsePosetFrames(t *testing.T) {
	p, index := initSparsePoset(t, common.NewTestLogger(t))

	participants := p.Participants.ToPeerSlice()

	if err := p.DivideRounds(); err != nil {
		t.Fatal(err)
	}
	if err := p.DecideFame(); err != nil {
		t.Fatal(err)
	}
	if err := p.DecideRoundReceived(); err != nil {
		t.Fatal(err)
	}
	if err := p.ProcessDecidedRounds(); err != nil {
		t.Fatal(err)
	}

	for bi := int64(0); bi < 5; bi++ {
		block, err := p.Store.GetBlock(bi)
		if err != nil {
			t.Fatal(err)
		}

		frame, err := p.GetFrame(block.RoundReceived())
		for k, ev := range frame.Events {
			ev.Body.Hash()
			hash, err := ev.Body.Hash()
			if err != nil {
				t.Fatal(err)
			}
			hex := fmt.Sprintf("0x%X", hash)
			r, err := p.round(hex)
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("frame %d event %d: %s, round %d",
				frame.Round, k, getName(index, hex), r)
		}
		for k, r := range frame.Roots {
			sp := getName(index, r.SelfParent.Hash)
			var ops []string
			for k := range r.Others {
				ops = append(ops, getName(index, k))
			}

			t.Logf("frame %d root %d: self parent index %s:"+
				" %v, others indexes %s: %v", frame.Round, k, sp,
				r.SelfParent, ops, r.Others)
		}
	}

	expectedFrameRoots := map[int64][]Root{
		1: {
			NewBaseRoot(participants[0].ID),
			NewBaseRoot(participants[1].ID),
			NewBaseRoot(participants[2].ID),
			NewBaseRoot(participants[3].ID),
		},
		2: {
			{
				NextRound: 1,
				SelfParent: &RootEvent{Hash: index[w00],
					CreatorID: participants[0].ID, Index: 0,
					LamportTimestamp: 0, Round: 0},
				Others: map[string]*RootEvent{
					index[w10]: {Hash: index[e32],
						CreatorID: participants[3].ID, Index: 1,
						LamportTimestamp: 3, Round: 1},
				},
			},
			{
				NextRound: 0,
				SelfParent: &RootEvent{Hash: index[w01],
					CreatorID: participants[1].ID, Index: 0,
					LamportTimestamp: 0, Round: 0},
				Others: map[string]*RootEvent{
					index[e10]: {Hash: index[w00],
						CreatorID: participants[0].ID, Index: 0,
						LamportTimestamp: 0, Round: 0},
				},
			},
			{
				NextRound: 1,
				SelfParent: &RootEvent{Hash: index[w02],
					CreatorID: participants[2].ID, Index: 0,
					LamportTimestamp: 0, Round: 0},
				Others: map[string]*RootEvent{
					index[e21]: {Hash: index[e10],
						CreatorID: participants[1].ID, Index: 1,
						LamportTimestamp: 1, Round: 0},
				},
			},
			NewBaseRoot(participants[3].ID),
		},
		3: {
			{
				NextRound: 2,
				SelfParent: &RootEvent{Hash: index[w10],
					CreatorID: participants[0].ID, Index: 1,
					LamportTimestamp: 4, Round: 1},
				Others: map[string]*RootEvent{
					index[f01]: {Hash: index[w11],
						CreatorID: participants[1].ID, Index: 2,
						LamportTimestamp: 5, Round: 2},
				},
			},
			{
				NextRound: 2,
				SelfParent: &RootEvent{Hash: index[e10],
					CreatorID: participants[1].ID, Index: 1,
					LamportTimestamp: 1, Round: 0},
				Others: map[string]*RootEvent{
					index[w11]: {Hash: index[w10],
						CreatorID: participants[0].ID, Index: 1,
						LamportTimestamp: 4, Round: 1},
				},
			},
			{
				NextRound: 2,
				SelfParent: &RootEvent{Hash: index[e21],
					CreatorID: participants[2].ID, Index: 1,
					LamportTimestamp: 2, Round: 1},
				Others: map[string]*RootEvent{
					index[w12]: {Hash: index[f01],
						CreatorID: participants[0].ID, Index: 2,
						LamportTimestamp: 6, Round: 2},
				},
			},
			{
				NextRound: 1,
				SelfParent: &RootEvent{Hash: index[w03],
					CreatorID: participants[3].ID, Index: 0,
					LamportTimestamp: 0, Round: 0},
				Others: map[string]*RootEvent{
					index[e32]: {Hash: index[e21],
						CreatorID: participants[2].ID, Index: 1,
						LamportTimestamp: 2, Round: 1},
				},
			},
		},
		4: {
			{
				NextRound: 2,
				SelfParent: &RootEvent{Hash: index[w10],
					CreatorID: participants[0].ID, Index: 1,
					LamportTimestamp: 4, Round: 1},
				Others: map[string]*RootEvent{
					index[f01]: {Hash: index[w11],
						CreatorID: participants[1].ID, Index: 2,
						LamportTimestamp: 5, Round: 2},
				},
			},
			{
				NextRound: 3,
				SelfParent: &RootEvent{Hash: index[w11],
					CreatorID: participants[1].ID, Index: 2,
					LamportTimestamp: 5, Round: 2},
				Others: map[string]*RootEvent{
					index[w21]: {Hash: index[w13],
						CreatorID: participants[3].ID, Index: 2,
						LamportTimestamp: 8, Round: 3},
				},
			},
			{
				NextRound: 3,
				SelfParent: &RootEvent{Hash: index[w12],
					CreatorID: participants[2].ID, Index: 2,
					LamportTimestamp: 7, Round: 2},
				Others: map[string]*RootEvent{
					index[w22]: {Hash: index[w21],
						CreatorID: participants[1].ID, Index: 3,
						LamportTimestamp: 9, Round: 3},
				},
			},
			{
				NextRound: 3,
				SelfParent: &RootEvent{Hash: index[e32],
					CreatorID: participants[3].ID, Index: 1,
					LamportTimestamp: 3, Round: 1},
				Others: map[string]*RootEvent{
					index[w13]: {Hash: index[w12],
						CreatorID: participants[2].ID, Index: 2,
						LamportTimestamp: 7, Round: 2},
				},
			},
		},
		5: {
			{
				NextRound: 2,
				SelfParent: &RootEvent{Hash: index[w10],
					CreatorID: participants[0].ID, Index: 1,
					LamportTimestamp: 4, Round: 1},
				Others: map[string]*RootEvent{
					index[f01]: {Hash: index[w11],
						CreatorID: participants[1].ID, Index: 2,
						LamportTimestamp: 5, Round: 2},
				},
			},
			{
				NextRound: 4,
				SelfParent: &RootEvent{Hash: index[w21],
					CreatorID: participants[1].ID, Index: 3,
					LamportTimestamp: 9, Round: 3},
				Others: map[string]*RootEvent{
					index[g13]: {Hash: index[w23],
						CreatorID: participants[3].ID, Index: 3,
						LamportTimestamp: 11, Round: 4},
				},
			},
			{
				NextRound: 4,
				SelfParent: &RootEvent{Hash: index[w22],
					CreatorID: participants[2].ID, Index: 3,
					LamportTimestamp: 10, Round: 3},
				Others: map[string]*RootEvent{
					index[w32]: {Hash: index[g13],
						CreatorID: participants[1].ID, Index: 4,
						LamportTimestamp: 12, Round: 4},
				},
			},
			{
				NextRound: 4,
				SelfParent: &RootEvent{Hash: index[w13],
					CreatorID: participants[3].ID, Index: 2,
					LamportTimestamp: 8, Round: 3},
				Others: map[string]*RootEvent{
					index[w23]: {Hash: index[w22],
						CreatorID: participants[2].ID, Index: 3,
						LamportTimestamp: 10, Round: 3},
				},
			},
		},
	}

	for bi := int64(0); bi < 5; bi++ {
		block, err := p.Store.GetBlock(bi)
		if err != nil {
			t.Fatal(err)
		}

		frame, err := p.GetFrame(block.RoundReceived())
		if err != nil {
			t.Fatal(err)
		}

		for k, r := range frame.Roots {
			compareRoots(t, r, &expectedFrameRoots[frame.Round][k], index)
		}
	}
}

func TestSparsePosetReset(t *testing.T) {
	p, index := initSparsePoset(t, common.NewTestLogger(t))

	p.DivideRounds()
	p.DecideFame()
	p.DecideRoundReceived()
	p.ProcessDecidedRounds()

	for bi := 0; bi < 5; bi++ {
		block, err := p.Store.GetBlock(int64(bi))
		if err != nil {
			t.Fatal(err)
		}

		frame, err := p.GetFrame(block.RoundReceived())
		if err != nil {
			t.Fatal(err)
		}

		// This operation clears the private fields which need to be recomputed
		// in the Events (round, roundReceived,etc)
		marshalledFrame, _ := frame.ProtoMarshal()
		unmarshalledFrame := new(Frame)
		unmarshalledFrame.ProtoUnmarshal(marshalledFrame)

		p2 := NewPoset(p.Participants,
			NewInmemStore(p.Participants, cacheSize),
			nil,
			testLogger(t))
		err = p2.Reset(block, *unmarshalledFrame)
		if err != nil {
			t.Fatal(err)
		}

		// Test continue after Reset

		// Compute diff
		p2Known := p2.Store.KnownEvents()
		diff := getDiff(p, p2Known, t)

		t.Logf("p2.Known: %v", p2Known)
		t.Logf("diff: %v", len(diff))

		wireDiff := make([]WireEvent, len(diff), len(diff))
		for i, e := range diff {
			wireDiff[i] = e.ToWire()
		}

		// Insert remaining Events into the Reset poset
		for i, wev := range wireDiff {
			eventName := getName(index, diff[i].Hex())
			ev, err := p2.ReadWireInfo(wev)
			if err != nil {
				t.Fatalf("ReadWireInfo(%s): %s", eventName, err)
			}
			compareEventMessages(t, &ev.Message, &diff[i].Message, index)
			err = p2.InsertEvent(*ev, false)
			if err != nil {
				t.Fatalf("InsertEvent(%s): %s", eventName, err)
			}
		}

		p2.DivideRounds()
		p2.DecideFame()
		p2.DecideRoundReceived()
		p2.ProcessDecidedRounds()

		compareRoundWitnesses(p, p2, index, int64(bi), true, t)
	}

}

func compareRoundWitnesses(p, p2 *Poset, index map[string]string, round int64, check bool, t *testing.T) {

	for i := round; i <= 5; i++ {
		pRound, err := p.Store.GetRoundCreated(i)
		if err != nil {
			t.Fatal(err)
		}
		p2Round, err := p2.Store.GetRoundCreated(i)
		if err != nil {
			t.Fatal(err)
		}

		//Check Round1 Witnesses
		pWitnesses := pRound.Witnesses()
		p2Witnesses := p2Round.Witnesses()
		sort.Strings(pWitnesses)
		sort.Strings(p2Witnesses)
		hwn := make([]string, len(pWitnesses))
		p2wn := make([]string, len(p2Witnesses))
		for _, w := range pWitnesses {
			hwn = append(hwn, getName(index, w))
		}
		for _, w := range p2Witnesses {
			p2wn = append(p2wn, getName(index, w))
		}

		if check && !reflect.DeepEqual(hwn, p2wn) {
			t.Fatalf("Reset Hg Round %d witnesses should be %v, not %v", i, hwn, p2wn)
		}
	}

}

func getDiff(p *Poset, known map[int64]int64, t *testing.T) []Event {
	var diff []Event
	for id, ct := range known {
		pk := p.Participants.ById[id].PubKeyHex
		// get participant Events with index > ct
		participantEvents, err := p.Store.ParticipantEvents(pk, ct)
		if err != nil {
			t.Fatal(err)
		}
		for _, e := range participantEvents {
			ev, err := p.Store.GetEvent(e)
			if err != nil {
				t.Fatal(err)
			}
			diff = append(diff, ev)
		}
	}
	sort.Sort(ByTopologicalOrder(diff))
	return diff
}

func getName(index map[string]string, hash string) string {
	for name, h := range index {
		if h == hash {
			return name
		}
	}
	return ""
}

func compareRootEvents(t *testing.T, x, exp *RootEvent,
	index map[string]string) {
	if x.Hash != exp.Hash || x.Index != exp.Index ||
		x.CreatorID != exp.CreatorID || x.Round != exp.Round ||
		x.LamportTimestamp != exp.LamportTimestamp {
		t.Fatalf("expected root event %s: %v, got %s: %v",
			getName(index, exp.Hash), exp, getName(index, x.Hash), x)
	}
}

func compareOtherParents(t *testing.T, x, exp map[string]*RootEvent,
	index map[string]string) {
	if len(x) != len(exp) {
		t.Fatalf("expected number of other parents: %d, got: %d",
			len(exp), len(x))
	}

	var others []string
	for k := range x {
		others = append(others, getName(index, k))
	}

	for k, v := range exp {
		root, ok := x[k]
		if !ok {
			t.Fatalf("root %v not exists in other roots: %s", v, others)
		}
		compareRootEvents(t, root, v, index)
	}
}

func compareRoots(t *testing.T, x, exp *Root, index map[string]string) {
	compareRootEvents(t, x.SelfParent, exp.SelfParent, index)
	compareOtherParents(t, x.Others, exp.Others, index)
	if exp.NextRound != x.NextRound {
		t.Fatalf("expected next round: %d, got: %d",
			exp.NextRound, x.NextRound)
	}
}

func compareEventMessages(t *testing.T, x, exp *EventMessage,
	index map[string]string) {
	if !reflect.DeepEqual(x.WitnessProof, exp.WitnessProof) ||
		!bytes.Equal(x.FlagTable, exp.FlagTable) ||
		x.Signature != exp.Signature {
		hash, _ := exp.Body.Hash()
		hex := fmt.Sprintf("0x%X", hash)
		t.Fatalf("expcted message to event %s: %v, got: %v",
			getName(index, hex), exp, x)
	}
	compareEventBody(t, x.Body, exp.Body)
}

func compareEventBody(t *testing.T, x, exp *EventBody) {
	if x.Index != exp.Index || !bytes.Equal(x.Creator, exp.Creator) ||
		!reflect.DeepEqual(x.BlockSignatures, exp.BlockSignatures) ||
		!reflect.DeepEqual(x.InternalTransactions, exp.InternalTransactions) ||
		!reflect.DeepEqual(x.Parents, exp.Parents) ||
		!reflect.DeepEqual(x.Transactions, exp.Transactions) {
		t.Fatalf("expcted event body: %v, got: %v", exp, x)
	}
}
