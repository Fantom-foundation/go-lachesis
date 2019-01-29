package model_test

import (
	"bytes"
	"github.com/Fantom-foundation/go-lachesis/src/poset2/model"
	"reflect"

	"testing"
)

func TestFrameEncodeDecode(t *testing.T) {
	fxt := newFixture()
	codec := model.NewDefaultCodec()
	frame := model.NewFrame(1001, []*model.Event{fxt.event})
	raw, err := frame.Encode(codec)
	if err != nil {
		t.Fatal(err)
	}

	frame2, err := model.DecodeFrame(codec, raw)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(frame, frame2) {
		t.Fatal("frames not equals")
	}
}

func TestFrameHash(t *testing.T) {
	fxt := newFixture()
	data := newTestEventData()
	event1 := model.NewEvent(data.creator, data.flags,
		data.index+1, data.parents, data.transactions)

	txs := append(data.transactions, []byte("txX"))
	event2 := model.NewEvent(data.creator, data.flags,
		data.index, data.parents, txs)

	codec := model.NewDefaultCodec()
	frame := model.NewFrame(1001, []*model.Event{fxt.event})
	wrongFrame1 := model.NewFrame(1001, []*model.Event{event1})
	wrongFrame2 := model.NewFrame(1001, []*model.Event{event2})

	hash1 := frame.Hash(codec)
	hash2 := wrongFrame1.Hash(codec)
	hash3 := wrongFrame2.Hash(codec)

	if bytes.Equal(hash1.Bytes(), hash2.Bytes()) ||
		bytes.Equal(hash1.Bytes(), hash3.Bytes()) {
		t.Fatal("hashes should be different")
	}
}
