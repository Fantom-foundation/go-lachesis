package model_test

import (
	"bytes"
	"crypto/ecdsa"
	"github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/poset2/model"
	"reflect"
	"testing"

	"github.com/Fantom-foundation/go-lachesis/src/crypto"
)

type fixture struct {
	key     *ecdsa.PrivateKey
	creator []byte
	codec   model.Codec
	event   *model.Event
}

func newFixture() *fixture {
	key, _ := crypto.GenerateECDSAKey()
	creator := crypto.FromECDSAPub(&key.PublicKey)
	event := newEvent(creator)
	codec := model.NewDefaultCodec()

	return &fixture{
		key:     key,
		creator: creator,
		codec:   codec,
		event:   event,
	}
}

func newEvent(creator []byte) *model.Event {
	flags := []common.Hash{common.HexToHash("0x0"),
		common.HexToHash("0x1")}

	var index uint64 = 100
	parents := []common.Hash{common.HexToHash("0x0"),
		common.HexToHash("0x1")}
	transactions := [][]byte{[]byte("t1"), []byte("t2")}

	return model.NewEvent(creator, flags, index, parents, transactions)

}

func TestEncodeDecodeEvent(t *testing.T) {
	fxt := newFixture()

	raw, err := fxt.event.Encode(fxt.codec)
	if err != nil {
		t.Fatal(err)
	}

	event2, err := model.DecodeEvent(fxt.codec, raw)
	if err != nil {
		t.Fatal(err)
	}

	if !fxt.event.Equals(event2, fxt.codec) {
		t.Fatal("hashes is not equal")
	}
}

func TestEventHash(t *testing.T) {
	fxt := newFixture()
	emptyHash := common.Hash{}
	hash := fxt.event.Hash(fxt.codec)

	if bytes.Equal(hash.Bytes(), emptyHash.Bytes()) {
		t.Fatal("hashes is not equal")
	}

	event2 := fxt.event
	if !bytes.Equal(hash.Bytes(), event2.Hash(fxt.codec).Bytes()) {
		t.Fatal("hashes is not equal")
	}
}

type testEventData struct {
	creator      []byte
	flags        []common.Hash
	index        uint64
	parents      []common.Hash
	transactions [][]byte
}

func newTestEventData() *testEventData {
	key, _ := crypto.GenerateECDSAKey()
	creator := crypto.FromECDSAPub(&key.PublicKey)
	flags := []common.Hash{common.HexToHash("0x0"),
		common.HexToHash("0x1")}

	var index uint64 = 100
	parents := []common.Hash{common.HexToHash("0x0"),
		common.HexToHash("0x1")}
	transactions := [][]byte{[]byte("t1"), []byte("t2")}

	return &testEventData{
		creator:      creator,
		flags:        flags,
		index:        index,
		parents:      parents,
		transactions: transactions,
	}
}

func TestExternalEventData(t *testing.T) {
	data := newTestEventData()

	event := model.NewEvent(data.creator, data.flags,
		data.index, data.parents, data.transactions)

	if !bytes.Equal(data.creator, event.Creator()) {
		t.Fatalf("expected %v, got %v", data.creator, event.Creator())
	}

	if !reflect.DeepEqual(data.flags, event.Flags()) {
		t.Fatalf("expected %v, got %v", data.flags, event.Flags())
	}

	if data.index != event.Lamport() {
		t.Fatalf("expected %v, got %v", data.index, event.Lamport())
	}

	if !reflect.DeepEqual(data.parents, event.Parents()) {
		t.Fatalf("expected %v, got %v", data.parents, event.Parents())
	}

	if !reflect.DeepEqual(data.transactions, event.Transactions()) {
		t.Fatalf("expected %v, got %v", data.transactions, event.Transactions())
	}
}

func TestEventVerify(t *testing.T) {
	fxt := newFixture()
	_, err := fxt.event.VerifySignature(fxt.codec)
	if err == nil {
		t.Fatal("error must be nil")
	}

	if err := fxt.event.Sign(fxt.key, fxt.codec); err != nil {
		t.Fatal(err)
	}

	valid, err := fxt.event.VerifySignature(fxt.codec)
	if err != nil {
		t.Fatal(err)
	}

	if !valid {
		t.Fatal("invalid event")
	}
}

func TestEventIsLoaded(t *testing.T) {
	data := newTestEventData()
	event := model.NewEvent(data.creator, data.flags,
		data.index, data.parents, data.transactions)
	if !event.IsLoaded() {
		t.Fatal("event must be loaded")
	}

	event2 := model.NewEvent(data.creator, data.flags,
		data.index, data.parents, nil)
	if event2.IsLoaded() {
		t.Fatal("event must not be loaded")
	}
}

func TestWireEvent(t *testing.T) {
	fxt := newFixture()
	_, err := model.EventFromWire(fxt.codec, []byte("123"))
	if err == nil {
		t.Fatal("event must not be nil")
	}

	if err := fxt.event.Sign(fxt.key, fxt.codec); err != nil {
		t.Fatal(err)
	}

	wireEvent, err := fxt.event.ToWire(fxt.codec)
	if err != nil {
		t.Fatal(err)
	}

	event, err := model.EventFromWire(fxt.codec, wireEvent)
	if err != nil {
		t.Fatal(err)
	}

	if !fxt.event.Equals(event, fxt.codec) {
		t.Fatal("hashes is not equal")
	}

	valid, err := event.VerifySignature(fxt.codec)
	if err != nil {
		t.Fatal(err)
	}

	if !valid {
		t.Fatal("invalid event")
	}
}
