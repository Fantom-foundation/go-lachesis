package node

import (
	"crypto/ecdsa"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/crypto"
	"github.com/Fantom-foundation/go-lachesis/src/peers"
	"github.com/Fantom-foundation/go-lachesis/src/poset"
)

func initCores(n int, t *testing.T) ([]*Core,
	map[uint32]*ecdsa.PrivateKey, map[string]string) {
	cacheSize := 1000

	var cores []*Core
	index := make(map[string]string)
	participantKeys := map[uint32]*ecdsa.PrivateKey{}

	var tempPeers []*peers.Peer
	for i := 0; i < n; i++ {
		key, _ := crypto.GenerateECDSAKey()
		pubHex := fmt.Sprintf("0x%X",
			crypto.FromECDSAPub(&key.PublicKey))
		peer := peers.NewPeer(pubHex, "")
		tempPeers = append(tempPeers, peer)
		participantKeys[peer.ID] = key
	}

	peerSet := peers.NewPeerSet(tempPeers)

	for i, peer := range peerSet.Peers {
		core := NewCore(uint32(i),
			participantKeys[peer.ID],
			peerSet,
			poset.NewInmemStore(peerSet, cacheSize),
			nil,
			common.NewTestLogger(t))

		selfParent := fmt.Sprintf("Root%d", peer.ID)

		flagTable := make(map[string]int64)
		flagTable[selfParent] = 1

		// Create and save the first Event
		initialEvent := poset.NewEvent([][]byte(nil),
			[]*poset.InternalTransaction{},
			nil,
			[]string{selfParent, ""}, core.PubKey(), 0, flagTable)
		err := core.SignAndInsertSelfEvent(initialEvent)
		if err != nil {
			t.Fatal(err)
		}

		core.RunConsensus()

		cores = append(cores, core)
		index[fmt.Sprintf("e%d", i)] = core.head
	}

	return cores, participantKeys, index
}

/*
|  e12  |
|   | \ |
|   |   e20
|   | / |
|   /   |
| / |   |
e01 |   |
| \ |   |
e0  e1  e2
0   1   2
*/
func initPoset(t *testing.T, cores []*Core, keys map[uint32]*ecdsa.PrivateKey,
	index map[string]string, participant uint32) {
	for i := uint32(0); i < uint32(len(cores)); i++ {
		if i != participant {
			event, err := cores[i].GetEventBlock(index[fmt.Sprintf("e%d", i)])
			if err != nil {
				t.Fatal(err)
			}
			err = cores[participant].InsertEvent(event, true)
			if err != nil {
				t.Fatalf("error inserting %s: %s\n",
					getName(index, event.Hex()), err)
			}
		}
	}

	// Get flag tables from parents
	event0, err := cores[0].poset.Store.GetEventBlock(index["e0"])
	if err != nil {
		t.Fatalf("failed to get parent: %s", err)
	}
	event1, err := cores[0].poset.Store.GetEventBlock(index["e1"])
	if err != nil {
		t.Fatalf("failed to get parent: %s", err)
	}

	event1ft, _ := event1.GetFlagTable()
	event01ft, _ := event0.MergeFlagTable(event1ft)

	event01 := poset.NewEvent([][]byte{},
		[]*poset.InternalTransaction{},
		nil,
		[]string{index["e0"], index["e1"]}, // e0 and e1
		cores[0].PubKey(), 1, event01ft)
	if err := insertEvent(cores, keys, index, event01, "e01", participant,
		common.Hash32(cores[0].pubKey)); err != nil {
		t.Fatalf("error inserting e01: %s\n", err)
	}

	// Get flag tables from parents
	event2, err := cores[2].poset.Store.GetEventBlock(index["e2"])
	if err != nil {
		t.Fatalf("failed to get parent: %s", err)
	}

	event20ft, _ := event2.MergeFlagTable(event01ft)

	event20 := poset.NewEvent([][]byte{},
		[]*poset.InternalTransaction{},
		nil,
		[]string{index["e2"], index["e01"]}, // e2 and e01
		cores[2].PubKey(), 1, event20ft)
	if err := insertEvent(cores, keys, index, event20, "e20", participant,
		common.Hash32(cores[2].pubKey)); err != nil {
		fmt.Printf("error inserting e20: %s\n", err)
	}

	event12ft, _ := event1.MergeFlagTable(event20ft)

	event12 := poset.NewEvent([][]byte{},
		[]*poset.InternalTransaction{},
		nil,
		[]string{index["e1"], index["e20"]}, // e1 and e20
		cores[1].PubKey(), 1, event12ft)
	if err := insertEvent(cores, keys, index, event12, "e12", participant,
		common.Hash32(cores[1].pubKey)); err != nil {
		fmt.Printf("error inserting e12: %s\n", err)
	}
}

func insertEvent(cores []*Core, keys map[uint32]*ecdsa.PrivateKey,
	index map[string]string, event *poset.Event, name string, participant uint32,
	creator uint32) error {

	if participant == creator {
		if err := cores[participant].SignAndInsertSelfEvent(event); err != nil {
			return err
		}
		// event is not signed because passed by value
		index[name] = cores[participant].head
	} else {
		event.Sign(keys[creator])
		if err := cores[participant].InsertEvent(event, true); err != nil {
			return err
		}
		index[name] = event.Hex()
	}
	return nil
}

func checkHeights(
	cores []*Core, expectedHeights []map[string]uint64, t *testing.T) {
	for i, core := range cores {
		heights := core.Heights()
		if !reflect.DeepEqual(heights, expectedHeights[i]) {
			t.Errorf("Cores[%d].Heights() should be %v, not %v",
				i, expectedHeights[i], heights)
		}
	}
}

func TestEventDiff(t *testing.T) {
	cores, keys, index := initCores(3, t)

	initPoset(t, cores, keys, index, 0)

	/*
	  P0 knows

	  |  e12  |
	  |   | \ |
	  |   |   e20
	  |   | / |
	  |   /   |
	  | / |   |
	  e01 |   |        P1 knows
	  | \ |   |
	  e0  e1  e2       |   e1  |
	  0   1   2        0   1   2
	*/

	knownBy1 := cores[1].KnownEvents()
	unknownBy1, err := cores[0].EventDiff(knownBy1)
	if err != nil {
		t.Fatal(err)
	}

	if l := len(unknownBy1); l != 5 {
		t.Fatalf("length of unknown should be 5, not %d", l)
	}

	expectedOrder := []string{"e0", "e2", "e01", "e20", "e12"}
	for i, e := range unknownBy1 {
		if name := getName(index, e.Hex()); name != expectedOrder[i] {
			t.Fatalf("element %d should be %s, not %s",
				i, expectedOrder[i], name)
		}
	}

}

func TestSync(t *testing.T) {
	cores, _, index := initCores(3, t)

	/*
	   core 0           core 1          core 2

	   e0  |   |        |   e1  |       |   |   e2
	   0   1   2        0   1   2       0   1   2
	*/

	expectedHeights := make([]map[string]uint64, 3)
	expectedHeights[0] = map[string]uint64{
		cores[0].hexID: 1,
		cores[1].hexID: 0,
		cores[2].hexID: 0,
	}
	expectedHeights[1] = map[string]uint64{
		cores[0].hexID: 0,
		cores[1].hexID: 1,
		cores[2].hexID: 0,
	}
	expectedHeights[2] = map[string]uint64{
		cores[0].hexID: 0,
		cores[1].hexID: 0,
		cores[2].hexID: 1,
	}
	checkHeights(cores, expectedHeights, t)

	// core 1 is going to tell core 0 everything it knows
	if err := synchronizeCores(cores, 1, 0, [][]byte{}); err != nil {
		t.Fatal(err)
	}

	/*
	   core 0           core 1          core 2

	   e01 |   |
	   | \ |   |
	   e0  e1  |        |   e1  |       |   |   e2
	   0   1   2        0   1   2       0   1   2
	*/

	expectedHeights[0] = map[string]uint64{
		cores[0].hexID: 2,
		cores[1].hexID: 1,
		cores[2].hexID: 0,
	}
	expectedHeights[1] = map[string]uint64{
		cores[0].hexID: 0,
		cores[1].hexID: 1,
		cores[2].hexID: 0,
	}
	expectedHeights[2] = map[string]uint64{
		cores[0].hexID: 0,
		cores[1].hexID: 0,
		cores[2].hexID: 1,
	}
	checkHeights(cores, expectedHeights, t)

	knownBy0 := cores[0].KnownEvents()
	if k := knownBy0[common.Hash32(cores[0].pubKey)]; k != 1 {
		t.Fatalf("core 0 should have last-index 1 for core 0, not %d", k)
	}
	if k := knownBy0[common.Hash32(cores[1].pubKey)]; k != 0 {
		t.Fatalf("core 0 should have last-index 0 for core 1, not %d", k)
	}
	if k := knownBy0[common.Hash32(cores[2].pubKey)]; k != -1 {
		t.Fatalf("core 0 should have last-index -1 for core 2, not %d", k)
	}
	core0Head, _ := cores[0].GetHead()
	if core0Head.SelfParent() != index["e0"] {
		t.Fatalf("core 0 head self-parent should be e0")
	}
	if core0Head.OtherParent() != index["e1"] {
		t.Fatalf("core 0 head other-parent should be e1")
	}
	if len(core0Head.Message.FlagTable) == 0 {
		t.Fatal("flag table is null")
	}
	index["e01"] = core0Head.Hex()

	// core 0 is going to tell core 2 everything it knows
	if err := synchronizeCores(cores, 0, 2, [][]byte{}); err != nil {
		t.Fatal(err)
	}

	/*

	   core 0           core 1          core 2

	                                    |   |  e20
	                                    |   | / |
	                                    |   /   |
	                                    | / |   |
	   e01 |   |                        e01 |   |
	   | \ |   |                        | \ |   |
	   e0  e1  |        |   e1  |       e0  e1  e2
	   0   1   2        0   1   2       0   1   2
	*/

	expectedHeights[0] = map[string]uint64{
		cores[0].hexID: 2,
		cores[1].hexID: 1,
		cores[2].hexID: 0,
	}
	expectedHeights[1] = map[string]uint64{
		cores[0].hexID: 0,
		cores[1].hexID: 1,
		cores[2].hexID: 0,
	}
	expectedHeights[2] = map[string]uint64{
		cores[0].hexID: 2,
		cores[1].hexID: 1,
		cores[2].hexID: 2,
	}
	checkHeights(cores, expectedHeights, t)

	knownBy2 := cores[2].KnownEvents()
	if k := knownBy2[common.Hash32(cores[0].pubKey)]; k != 1 {
		t.Fatalf("core 2 should have last-index 1 for core 0, not %d", k)
	}
	if k := knownBy2[common.Hash32(cores[1].pubKey)]; k != 0 {
		t.Fatalf("core 2 should have last-index 0 core 1, not %d", k)
	}
	if k := knownBy2[common.Hash32(cores[2].pubKey)]; k != 1 {
		t.Fatalf("core 2 should have last-index 1 for core 2, not %d", k)
	}
	core2Head, _ := cores[2].GetHead()
	if core2Head.SelfParent() != index["e2"] {
		t.Fatalf("core 2 head self-parent should be e2")
	}
	if core2Head.OtherParent() != index["e01"] {
		t.Fatalf("core 2 head other-parent should be e01")
	}
	index["e20"] = core2Head.Hex()

	// core 2 is going to tell core 1 everything it knows
	if err := synchronizeCores(cores, 2, 1, [][]byte{}); err != nil {
		t.Fatal(err)
	}

	/*
	   core 0           core 1          core 2

	                    |  e12  |
	                    |   | \ |
	                    |   |  e20      |   |  e20
	                    |   | / |       |   | / |
	                    |   /   |       |   /   |
	                    | / |   |       | / |   |
	   e01 |   |        e01 |   |       e01 |   |
	   | \ |   |        | \ |   |       | \ |   |
	   e0  e1  |        e0  e1  e2      e0  e1  e2
	   0   1   2        0   1   2       0   1   2
	*/

	expectedHeights[0] = map[string]uint64{
		cores[0].hexID: 2,
		cores[1].hexID: 1,
		cores[2].hexID: 0,
	}
	expectedHeights[1] = map[string]uint64{
		cores[0].hexID: 2,
		cores[1].hexID: 2,
		cores[2].hexID: 2,
	}
	expectedHeights[2] = map[string]uint64{
		cores[0].hexID: 2,
		cores[1].hexID: 1,
		cores[2].hexID: 2,
	}
	checkHeights(cores, expectedHeights, t)

	knownBy1 := cores[1].KnownEvents()
	if k := knownBy1[common.Hash32(cores[0].pubKey)]; k != 1 {
		t.Fatalf("core 1 should have last-index 1 for core 0, not %d", k)
	}
	if k := knownBy1[common.Hash32(cores[1].pubKey)]; k != 1 {
		t.Fatalf("core 1 should have last-index 1 for core 1, not %d", k)
	}
	if k := knownBy1[common.Hash32(cores[2].pubKey)]; k != 1 {
		t.Fatalf("core 1 should have last-index 1 for core 2, not %d", k)
	}
	core1Head, _ := cores[1].GetHead()
	if core1Head.SelfParent() != index["e1"] {
		t.Fatalf("core 1 head self-parent should be e1")
	}
	if core1Head.OtherParent() != index["e20"] {
		t.Fatalf("core 1 head other-parent should be e20")
	}
	index["e12"] = core1Head.Hex()

}

func checkInDegree(
	cores []*Core, expectedInDegree []map[string]uint64, t *testing.T) {
	for i, core := range cores {
		inDegrees := core.InDegrees()
		if !reflect.DeepEqual(inDegrees, expectedInDegree[i]) {
			t.Errorf("Cores[%d].InDegrees() should be %v, not %v",
				i, expectedInDegree[i], inDegrees)
		}
	}
}

func TestInDegrees(t *testing.T) {
	cores, _, _ := initCores(3, t)

	/*
	   core 0           core 1          core 2

	   e0  |   |        |   e1  |       |   |   e2
	   0   1   2        0   1   2       0   1   2
	*/

	// core 1 is going to tell core 0 everything it knows
	if err := synchronizeCores(cores, 1, 0, [][]byte{}); err != nil {
		t.Fatal(err)
	}

	/*
	   core 0           core 1          core 2

	   e01 |   |
	   | \ |   |
	   e0  e1  |        |   e1  |       |   |   e2
	   0   1   2        0   1   2       0   1   2
	*/

	expectedHeights := make([]map[string]uint64, 3)
	expectedHeights[0] = map[string]uint64{
		cores[0].hexID: 2,
		cores[1].hexID: 1,
		cores[2].hexID: 0,
	}
	expectedHeights[1] = map[string]uint64{
		cores[0].hexID: 0,
		cores[1].hexID: 1,
		cores[2].hexID: 0,
	}
	expectedHeights[2] = map[string]uint64{
		cores[0].hexID: 0,
		cores[1].hexID: 0,
		cores[2].hexID: 1,
	}
	checkHeights(cores, expectedHeights, t)

	expectedInDegree := make([]map[string]uint64, 3)
	expectedInDegree[0] = map[string]uint64{
		cores[0].hexID: 0,
		cores[1].hexID: 1,
		cores[2].hexID: 0,
	}
	expectedInDegree[1] = map[string]uint64{
		cores[0].hexID: 0,
		cores[1].hexID: 0,
		cores[2].hexID: 0,
	}
	expectedInDegree[2] = map[string]uint64{
		cores[0].hexID: 0,
		cores[1].hexID: 0,
		cores[2].hexID: 0,
	}
	checkInDegree(cores, expectedInDegree, t)

	// core 1 is going to tell core 2 everything it knows
	if err := synchronizeCores(cores, 1, 2, [][]byte{}); err != nil {
		t.Fatal(err)
	}

	/*

	   core 0           core 1          core 2

	   e01 |   |                        |   |  e21
	   | \ |   |                        |   | / |
	   e0  e1  |        |   e1  |       e0  e1  e2
	   0   1   2        0   1   2       0   1   2
	*/

	expectedHeights[0] = map[string]uint64{
		cores[0].hexID: 2,
		cores[1].hexID: 1,
		cores[2].hexID: 0,
	}
	expectedHeights[1] = map[string]uint64{
		cores[0].hexID: 0,
		cores[1].hexID: 1,
		cores[2].hexID: 0,
	}
	expectedHeights[2] = map[string]uint64{
		cores[0].hexID: 0,
		cores[1].hexID: 1,
		cores[2].hexID: 2,
	}
	checkHeights(cores, expectedHeights, t)

	expectedInDegree[0] = map[string]uint64{
		cores[0].hexID: 0,
		cores[1].hexID: 1,
		cores[2].hexID: 0,
	}
	expectedInDegree[1] = map[string]uint64{
		cores[0].hexID: 0,
		cores[1].hexID: 0,
		cores[2].hexID: 0,
	}
	expectedInDegree[2] = map[string]uint64{
		cores[0].hexID: 0,
		cores[1].hexID: 1,
		cores[2].hexID: 0,
	}
	checkInDegree(cores, expectedInDegree, t)

	// core 0 is going to tell core 2 everything it knows
	if err := synchronizeCores(cores, 0, 2, [][]byte{}); err != nil {
		t.Fatal(err)
	}

	/*

	   core 0           core 1          core 2

	                                    |   |  e20
	                                    |   | / |
	                                    |   /   |
	                                    | / |   |
	   e01 |   |                        e01 |  e21
	   | \ |   |                        | \ | / |
	   e0  e1  |        |   e1  |       e0  e1  e2
	   0   1   2        0   1   2       0   1   2
	*/

	expectedHeights[0] = map[string]uint64{
		cores[0].hexID: 2,
		cores[1].hexID: 1,
		cores[2].hexID: 0,
	}
	expectedHeights[1] = map[string]uint64{
		cores[0].hexID: 0,
		cores[1].hexID: 1,
		cores[2].hexID: 0,
	}
	expectedHeights[2] = map[string]uint64{
		cores[0].hexID: 2,
		cores[1].hexID: 1,
		cores[2].hexID: 3,
	}
	checkHeights(cores, expectedHeights, t)

	expectedInDegree[0] = map[string]uint64{
		cores[0].hexID: 0,
		cores[1].hexID: 1,
		cores[2].hexID: 0,
	}
	expectedInDegree[1] = map[string]uint64{
		cores[0].hexID: 0,
		cores[1].hexID: 0,
		cores[2].hexID: 0,
	}
	expectedInDegree[2] = map[string]uint64{
		cores[0].hexID: 1,
		cores[1].hexID: 2,
		cores[2].hexID: 0,
	}
	checkInDegree(cores, expectedInDegree, t)

	// core 2 is going to tell core 1 everything it knows
	if err := synchronizeCores(cores, 2, 1, [][]byte{}); err != nil {
		t.Fatal(err)
	}

	/*

	   core 0           core 1          core 2

	                    |  e12  |
	                    |   | \ |
	                    |   |  e20      |   |  e20
	                    |   | / |       |   | / |
	                    |   /   |       |   /   |
	                    | / |   |       | / |   |
	   e01 |   |        e01 |  e21      e01 |  e21
	   | \ |   |        | \ | / |       | \ | / |
	   e0  e1  |        e0  e1  e2      e0  e1  e2
	   0   1   2        0   1   2       0   1   2
	*/

	expectedHeights[0] = map[string]uint64{
		cores[0].hexID: 2,
		cores[1].hexID: 1,
		cores[2].hexID: 0,
	}
	expectedHeights[1] = map[string]uint64{
		cores[0].hexID: 2,
		cores[1].hexID: 2,
		cores[2].hexID: 3,
	}
	expectedHeights[2] = map[string]uint64{
		cores[0].hexID: 2,
		cores[1].hexID: 1,
		cores[2].hexID: 3,
	}
	checkHeights(cores, expectedHeights, t)

	expectedInDegree[0] = map[string]uint64{
		cores[0].hexID: 0,
		cores[1].hexID: 1,
		cores[2].hexID: 0,
	}
	expectedInDegree[1] = map[string]uint64{
		cores[0].hexID: 1,
		cores[1].hexID: 0,
		cores[2].hexID: 1,
	}
	expectedInDegree[2] = map[string]uint64{
		cores[0].hexID: 1,
		cores[1].hexID: 2,
		cores[2].hexID: 0,
	}
	checkInDegree(cores, expectedInDegree, t)
}

/*
h0  |   h2
| \ | / |
|   h1  |
|  /|   |--------------------
g02 |   | R2
| \ |   |
|   \   |
|   | \ |
|   |  g21
|   | / |
|  g10  |
| / |   |
g0  |   g2
| \ | / |
|   g1  |
|  /|   |--------------------
f02 |   | R1
| \ |   |
|   \   |
|   | \ |
|   |  f21
|   | / |
|  f10  |
| / |   |
f0  |   f2
| \ | / |
|   f1  |
|  /|   |--------------------
e02 |   | R0 Consensus
| \ |   |
|   \   |
|   | \ |
|   |  e21
|   | / |
|  e10  |
| / |   |
e0  e1  e2
0   1    2
*/
type play struct {
	from    int
	to      int
	payload [][]byte
}

func initConsensusPoset(t *testing.T) []*Core {
	cores, _, _ := initCores(3, t)
	playbook := []play{
		{from: 0, to: 1, payload: [][]byte{[]byte("e10")}},
		{from: 1, to: 2, payload: [][]byte{[]byte("e21")}},
		{from: 2, to: 0, payload: [][]byte{[]byte("e02")}},
		{from: 0, to: 1, payload: [][]byte{[]byte("f1")}},
		{from: 1, to: 0, payload: [][]byte{[]byte("f0")}},
		{from: 1, to: 2, payload: [][]byte{[]byte("f2")}},

		{from: 0, to: 1, payload: [][]byte{[]byte("f10")}},
		{from: 1, to: 2, payload: [][]byte{[]byte("f21")}},
		{from: 2, to: 0, payload: [][]byte{[]byte("f02")}},
		{from: 0, to: 1, payload: [][]byte{[]byte("g1")}},
		{from: 1, to: 0, payload: [][]byte{[]byte("g0")}},
		{from: 1, to: 2, payload: [][]byte{[]byte("g2")}},

		{from: 0, to: 1, payload: [][]byte{[]byte("g10")}},
		{from: 1, to: 2, payload: [][]byte{[]byte("g21")}},
		{from: 2, to: 0, payload: [][]byte{[]byte("g02")}},
		{from: 0, to: 1, payload: [][]byte{[]byte("h1")}},
		{from: 1, to: 0, payload: [][]byte{[]byte("h0")}},
		{from: 1, to: 2, payload: [][]byte{[]byte("h2")}},
	}

	for _, play := range playbook {
		if err := syncAndRunConsensus(
			cores, play.from, play.to, play.payload); err != nil {
			t.Fatal(err)
		}
	}
	return cores
}

func TestConsensus(t *testing.T) {
	cores := initConsensusPoset(t)

	if l := len(cores[0].GetConsensusEvents()); l != 4 {
		t.Fatalf("length of consensus should be 4 not %d", l)
	}

	core0Consensus := cores[0].GetConsensusEvents()
	core1Consensus := cores[1].GetConsensusEvents()
	core2Consensus := cores[2].GetConsensusEvents()

	for i, e := range core0Consensus {
		if core1Consensus[i] != e {
			t.Fatalf("core 1 consensus[%d] does not match core 0's", i)
		}
		if core2Consensus[i] != e {
			t.Fatalf("core 2 consensus[%d] does not match core 0's", i)
		}
	}
}

func TestOverSyncLimit(t *testing.T) {
	cores := initConsensusPoset(t)

	// positive
	known := map[uint32]int64{
		common.Hash32(cores[0].pubKey): 1,
		common.Hash32(cores[1].pubKey): 1,
		common.Hash32(cores[2].pubKey): 1,
	}

	syncLimit := int64(10)

	if !cores[0].OverSyncLimit(known, syncLimit) {
		t.Fatalf("OverSyncLimit(%v, %v) should return true", known, syncLimit)
	}

	// negative
	known = map[uint32]int64{
		common.Hash32(cores[0].pubKey): 6,
		common.Hash32(cores[1].pubKey): 6,
		common.Hash32(cores[2].pubKey): 6,
	}

	if cores[0].OverSyncLimit(known, syncLimit) {
		t.Fatalf("OverSyncLimit(%v, %v) should return false", known, syncLimit)
	}

	// edge
	known = map[uint32]int64{
		common.Hash32(cores[0].pubKey): 2,
		common.Hash32(cores[1].pubKey): 3,
		common.Hash32(cores[2].pubKey): 3,
	}
	if cores[0].OverSyncLimit(known, syncLimit) {
		t.Fatalf("OverSyncLimit(%v, %v) should return false", known, syncLimit)
	}

}

/*
    |   |   |   |-----------------
	|   w31 |   | R3
	|	| \ |   |
    |   |  w32  |
    |   |   | \ |
    |   |   |  w33
    |   |   | / |-----------------
    |   |  g21  | R2
	|   | / |   |
	|  w21  |   |
	|	| \ |   |
    |   |  w22  |
    |   |   | \ |
    |   |   |  w23
    |   |   | / |-----------------
    |   |  f21  | R1
	|   | / |   | LastConsensusRound
	|  w11  |   |
	|	| \ |   |
    |   |   \   |
    |   |   | \ |
	|   |   |  w13
	|   |   | / |
   FSE  |  w12  | FSE is only added after FastForward
    |\  | / |   | -----------------
    |  e13  |   | R0
	|	| \ |   |
    |   |   \   |
    |   |   | \ |
    |   |   |  e32
    |   |   | / |
    |   |  e21  | All Events in Round 0 are Consensus Events.
    |   | / |   |
    |   e1  e2  e3
    0	1	2	3
*/
func initFFPoset(cores []*Core, t *testing.T) {
	playbook := []play{
		{from: 1, to: 2, payload: [][]byte{[]byte("e21")}},
		{from: 2, to: 3, payload: [][]byte{[]byte("e32")}},
		{from: 3, to: 1, payload: [][]byte{[]byte("e13")}},
		{from: 1, to: 2, payload: [][]byte{[]byte("w12")}},
		{from: 2, to: 3, payload: [][]byte{[]byte("w13")}},
		{from: 3, to: 1, payload: [][]byte{[]byte("w11")}},
		{from: 1, to: 2, payload: [][]byte{[]byte("f21")}},
		{from: 2, to: 3, payload: [][]byte{[]byte("w23")}},
		{from: 3, to: 2, payload: [][]byte{[]byte("w22")}},
		{from: 2, to: 1, payload: [][]byte{[]byte("w21")}},
		{from: 1, to: 2, payload: [][]byte{[]byte("g21")}},
		{from: 2, to: 3, payload: [][]byte{[]byte("w33")}},
		{from: 3, to: 2, payload: [][]byte{[]byte("w32")}},
		{from: 2, to: 1, payload: [][]byte{[]byte("w31")}},
	}

	for k, play := range playbook {
		if err := syncAndRunConsensus(cores, play.from, play.to, play.payload); err != nil {
			t.Fatalf("play %d: %s", k, err)
		}
	}
}

func TestConsensusFF(t *testing.T) {
	cores, _, _ := initCores(4, t)
	initFFPoset(cores, t)

	if r := cores[1].GetLastConsensusRound(); r < 0 || r != 2 {
		disp := "nil"
		if r >= 0 {
			disp = strconv.FormatInt(r, 10)
		}
		t.Fatalf("Cores[1] last consensus Round should be 2, not %s", disp)
	}

	if l := len(cores[1].GetConsensusEvents()); l != 7 {
		t.Fatalf("Node 1 should have 7 consensus events, not %d", l)
	}

	core1Consensus := cores[1].GetConsensusEvents()
	core2Consensus := cores[2].GetConsensusEvents()
	core3Consensus := cores[3].GetConsensusEvents()

	for i, e := range core1Consensus {
		if core2Consensus[i] != e {
			t.Fatalf("Node 2 consensus[%d] does not match Node 1's", i)
		}
		if core3Consensus[i] != e {
			t.Fatalf("Node 3 consensus[%d] does not match Node 1's", i)
		}
	}
}

func TestCoreFastForward(t *testing.T) {
	cores, _, _ := initCores(4, t)
	initFFPoset(cores, t)

	t.Run("Test no Anchor", func(t *testing.T) {
		// Test no anchor block
		_, _, err := cores[1].GetAnchorBlockWithFrame()
		if err == nil {
			t.Fatal("GetAnchorBlockWithFrame should throw an error" +
				" because there is no anchor block yet")
		}
	})

	block0, err := cores[1].poset.Store.GetBlock(0)
	if err != nil {
		t.Fatal(err)
	}

	// collect signatures
	signatures := make([]poset.BlockSignature, 3)
	for k, c := range cores[1:] {
		b, err := c.poset.Store.GetBlock(0)
		if err != nil {
			t.Fatal(err)
		}
		sig, err := c.SignBlock(b)
		if err != nil {
			t.Fatal(err)
		}
		signatures[k] = sig
	}

	t.Run("Test not enough signatures", func(t *testing.T) {
		// Append only 1 signatures
		if err := block0.SetSignature(signatures[0]); err != nil {
			t.Fatal(err)
		}

		// Save Block
		if err := cores[1].poset.Store.SetBlock(block0); err != nil {
			t.Fatal(err)
		}
		// Assign AnchorBlock
		cores[1].poset.AnchorBlock = new(int64)
		*cores[1].poset.AnchorBlock = 0

		// Now the function should find an AnchorBlock
		block, frame, err := cores[1].GetAnchorBlockWithFrame()
		if err != nil {
			t.Fatal(err)
		}

		err = cores[0].FastForward(cores[1].hexID, block, frame)
		// We should get an error because AnchorBlock doesnt contain enough
		// signatures
		if err == nil {
			t.Fatal("FastForward should throw an error because the Block" +
				" does not contain enough signatures")
		}
	})

	t.Run("Test positive", func(t *testing.T) {
		// Append the 2nd and 3rd signatures
		for i := 1; i < 3; i++ {
			if err := block0.SetSignature(signatures[i]); err != nil {
				t.Fatal(err)
			}
		}

		// Save Block
		if err := cores[1].poset.Store.SetBlock(block0); err != nil {
			t.Fatal(err)
		}

		block, frame, err := cores[1].GetAnchorBlockWithFrame()
		if err != nil {
			t.Fatal(err)
		}

		err = cores[0].FastForward(cores[1].hexID, block, frame)
		if err != nil {
			t.Fatal(err)
		}

		knownBy0 := cores[0].KnownEvents()
		if err != nil {
			t.Fatal(err)
		}

		expectedKnown := map[uint32]int64{
			common.Hash32(cores[0].pubKey): -1,
			common.Hash32(cores[1].pubKey): 0,
			common.Hash32(cores[2].pubKey): 1,
			common.Hash32(cores[3].pubKey): 0,
		}

		if !reflect.DeepEqual(knownBy0, expectedKnown) {
			t.Fatalf("Cores[0].Known should be %v, not %v",
				expectedKnown, knownBy0)
		}

		if r := cores[0].GetLastConsensusRound(); r < 0 || r != 1 {
			disp := "nil"
			if r >= 0 {
				disp = strconv.FormatInt(r, 10)
			}
			t.Fatalf("Cores[0] last consensus Round should be 1, not %s", disp)
		}

		if lbi := cores[0].poset.Store.LastBlockIndex(); lbi != 0 {
			t.Fatalf("Cores[0].poset.LastBlockIndex should be 0, not %d", lbi)
		}

		sBlock, err := cores[0].poset.Store.GetBlock(block.Index())
		if err != nil {
			t.Fatalf("Error retrieving latest Block from reset poset: %v", err)
		}
		if !reflect.DeepEqual(sBlock.Body, block.Body) {
			t.Fatalf("Blocks defer")
		}

		lastEventFrom0, _, err := cores[0].poset.Store.LastEventFrom(
			cores[0].hexID)
		if err != nil {
			t.Fatal(err)
		}
		if c0h := cores[0].Head(); c0h != lastEventFrom0 {
			t.Fatalf("Head should be %s, not %s", lastEventFrom0, c0h)
		}

		if c0s := cores[0].Seq; c0s != -1 {
			t.Fatalf("Seq should be %d, not %d", -1, c0s)
		}

	})

}

func synchronizeCores(cores []*Core, from int, to int, payload [][]byte) error {
	knownByTo := cores[to].KnownEvents()
	unknownByTo, err := cores[from].EventDiff(knownByTo)
	if err != nil {
		return err
	}

	unknownWire, err := cores[from].ToWire(unknownByTo)
	if err != nil {
		return err
	}

	err = cores[to].AddTransactions(payload)
	if err != nil {
		return err
	}

	return cores[to].Sync(unknownWire)
}

func syncAndRunConsensus(
	cores []*Core, from int, to int, payload [][]byte) error {
	if err := synchronizeCores(cores, from, to, payload); err != nil {
		return err
	}
	cores[to].RunConsensus()
	return nil
}

func getName(index map[string]string, hash string) string {
	for name, h := range index {
		if h == hash {
			return name
		}
	}
	return fmt.Sprintf("%s not found", hash)
}
