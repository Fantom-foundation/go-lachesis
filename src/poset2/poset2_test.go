package poset2_test

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/crypto"
	"github.com/Fantom-foundation/go-lachesis/src/poset2"
	"github.com/Fantom-foundation/go-lachesis/src/poset2/model"
	"github.com/Fantom-foundation/go-lachesis/src/poset2/store"
	"testing"
)

const (
	e0  = "e0"
	e1  = "e1"
	e2  = "e2"
	e10 = "e10"
	e21 = "e21"
	e02 = "e02"
	f1  = "f1"
	s20 = "s20"
	s10 = "s10"
	s00 = "s00"
	s11 = "s11"
)

type posetData struct {
	to          int
	index       uint64
	selfParent  string
	otherParent string
	name        string
	txs         [][]byte
	roots       []string
}

type dominatorData struct {
	newest string
	oldest string
	val    bool
	err    bool
}

type node struct {
	id     int
	key    *ecdsa.PrivateKey
	pub    []byte
	events []*model.Event
}

func newNode(index int, key *ecdsa.PrivateKey) *node {
	pub := crypto.FromECDSAPub(&key.PublicKey)
	return &node{
		id:  index,
		key: key,
		pub: pub,
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

func newPosetV1(t *testing.T, db *store.Store) (*poset2.Poset2, map[string]common.Hash) {
	data := []*posetData{
		{1, 1, e1, e0, e10, nil, []string{e0, e1}},
		{2, 1, e2, "", s20, nil, []string{e2}},
		{0, 1, e0, "", s00, nil, []string{e0}},
		{2, 2, s20, e10, e21, nil, []string{e0, e1, e2}},
		{0, 2, s00, e21, e02, nil, []string{e0, e21}},
		{1, 2, e10, "", s10, nil, []string{e0, e1}},
		{1, 3, s10, e02, f1, nil, []string{e21, e02, e1}},
		{1, 4, f1, "", s11, [][]byte{[]byte("abc")}, []string{e21, e02, f1}},
	}

	return initPoset(t, db, data, 3)
}

func initPoset(t *testing.T, db *store.Store, data []*posetData,
	numberNodes int) (*poset2.Poset2, map[string]common.Hash) {
	nodes := initNodes(numberNodes)
	index := make(map[string]common.Hash)
	codec := model.NewDefaultCodec()
	logger := common.NewTestLogger(t).WithField("id", "test")
	var orderedEvents []*model.Event

	// Creates genesis events.
	for k := range nodes {
		event, err := model.NewLeafEvent(codec, nodes[k].key)
		if err != nil {
			t.Fatal(err)
		}
		if err := event.Sign(nodes[k].key, codec); err != nil {
			t.Fatal(err)
		}

		nodes[k].events = append(nodes[k].events, event)
		orderedEvents = append(orderedEvents, event)

		index[fmt.Sprintf("e%d", k)] = event.Hash(codec)
	}

	// Creates events from poset data
	for k := range data {
		var flags []common.Hash
		for _, idx := range data[k].roots {
			flags = append(flags, index[idx])
		}

		parents := []common.Hash{index[data[k].selfParent],
			index[data[k].otherParent]}

		creator := nodes[data[k].to].pub

		event := model.NewEvent(
			creator, flags, data[k].index, parents, data[k].txs)
		if err := event.Sign(nodes[data[k].to].key, codec); err != nil {
			t.Fatal(err)
		}

		nodes[data[k].to].events = append(nodes[data[k].to].events, event)
		index[data[k].name] = event.Hash(codec)
		orderedEvents = append(orderedEvents, event)
	}

	validator := poset2.NewDefaultValidator(db, codec, len(nodes))
	poset := poset2.NewPoset2(db, validator, logger)
	for k := range orderedEvents {
		if err := poset.InsertEvent(orderedEvents[k]); err != nil {
			t.Fatal(err)
		}

	}

	return poset, index
}

func initNodes(number int) (nodes []*node) {
	for i := 0; i < number; i++ {
		key, _ := crypto.GenerateECDSAKey()
		nodes = append(nodes, newNode(i, key))
	}
	return nodes
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

func TestPosetDominated(t *testing.T) {
	db := store.NewStore()
	poset, index := newPosetV1(t, db)

	expected := []*dominatorData{
		{e0, e1, false, false},
		{e0, e2, false, false},
		{e1, e2, false, false},
		{e1, e1, false, false},

		{e10, e10, false, false},

		{e10, e0, true, false},
		{e10, e1, true, false},
		{e10, e2, false, false},
		{e10, s10, false, false},
		{e10, e02, false, false},

		{s20, e0, false, false},
		{s20, e1, false, false},
		{s20, e2, true, false},

		{s00, e0, true, false},
		{s00, e1, false, false},
		{s00, e10, false, false},
		{s00, s20, false, false},
		{s00, e21, false, false},

		{e21, e0, true, false},
		{e21, e1, true, false},
		{e21, e2, true, false},
		{e21, e10, true, false},
		{e21, s20, true, false},
		{e21, s00, false, false},
		{e21, s00, false, false},

		{s11, e0, true, false},
		{s11, e1, true, false},
		{s11, e2, true, false},
		{s11, e10, true, false},
		{s11, s20, true, false},
		{s11, s00, true, false},
		{s11, e21, true, false},

		{e1, s11, false, false},
		{e2, s11, false, false},
		{e10, s11, false, false},
		{s20, s11, false, false},
		{"", s11, false, true},
		{s20, "", false, true},
	}

	for _, exp := range expected {
		type resp struct {
			dominate bool
			err      error
		}

		results := make(chan *resp)

		go func() {
			var err2 error
			var dominate bool

			defer func() {
				results <- &resp{dominate, err2}
			}()

			newest, err := poset.GetEvent(index[exp.newest])
			if err != nil {
				err2 = err
				return
			}

			oldest, err := poset.GetEvent(index[exp.oldest])
			if err != nil {
				err2 = err
				return
			}

			dominate, err = poset.Validator.Dominated(newest, oldest)
			if err != nil {
				err2 = err
				return
			}
		}()

		result := <-results
		if (result.err == nil) != !exp.err {
			t.Fatalf("wrong error to eventblocks: %s, %s",
				exp.newest, exp.oldest)
		}
		if result.dominate != exp.val {
			t.Fatalf("expected %v, got %v to eventblocks: %s, %s",
				exp.val, result.dominate, exp.newest, exp.oldest)
		}
	}
}

func TestDivideRounds(t *testing.T) {
	db := store.NewStore()
	poset, _ := newPosetV1(t, db)
	expected := 2

	if err := poset.DivideRounds(); err != nil {
		t.Fatal(err)
	}

	if db.LastRound() != uint64(expected) {
		t.Fatal(db.LastRound())
	}

}
