package model_test

import (
	"github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/poset2/model"
	"reflect"
	"testing"
)

//func TestRoundEvents(t *testing.T) {
//	hash1 := common.HexToHash("0x1")
//	hash2 := common.HexToHash("0x2")
//
//	round := model.NewRound()
//	round.AddEvent(hash1, true)
//	round.AddEvent(hash2, false)
//	round.SetConsensusEvent(hash1)
//
//	consensusHashes := round.Events(true)
//	nonConsensusHashes := round.Events(false)
//
//	if len(consensusHashes) != 1 {
//		t.Fatal("wrong number of consensus hashes")
//	}
//
//	if len(nonConsensusHashes) != 1 {
//		t.Fatal("wrong number of non consensus hashes")
//	}
//
//	if !bytes.Equal(hash1.Bytes(), consensusHashes[0].Bytes()) {
//		t.Fatalf("expected hash %s, got %s",
//			hash1.String(), consensusHashes[0].String())
//	}
//
//	if !bytes.Equal(hash2.Bytes(), nonConsensusHashes[0].Bytes()) {
//		t.Fatalf("expected hash %s, got %s",
//			hash2.String(), nonConsensusHashes[0].String())
//	}
//}

func TestRoundQueued(t *testing.T) {
	round := model.NewRound()
	if round.IsQueued() {
		t.Fatalf("round must be unqueued")
	}

	round.SetQueued(true)

	if !round.IsQueued() {
		t.Fatalf("round must be queued")
	}
}

//func TestDecodeRound(t *testing.T) {
//	hash1 := common.HexToHash("0x1")
//	round := model.NewRound()
//	round.AddEvent(hash1, true)
//
//	if round.IsDecided() {
//		t.Fatalf("round must be undecided")
//	}
//
//	round.SetAtropos(hash1, true)
//
//	if !round.IsDecided() {
//		t.Fatalf("round must be decided")
//	}
//}
//
//func TestClothosDecided(t *testing.T) {
//	hash1 := common.HexToHash("0x1")
//	round := model.NewRound()
//	round.AddEvent(hash1, true)
//
//	if round.ClothosDecided() {
//		t.Fatalf("clothos is undecided")
//	}
//
//	round.SetAtropos(hash1, true)
//
//	if !round.ClothosDecided() {
//		t.Fatalf("clothos is decided")
//	}
//}

func TestRoundEncodeDecode(t *testing.T) {
	hash1 := common.HexToHash("0x1")

	round := model.NewRound()
	codec := model.NewDefaultCodec()
	round.AddEvent(hash1, true)

	raw, err := round.Encode(codec)
	if err != nil {
		t.Fatal(err)
	}

	round2, err := model.DecodeRound(codec, raw)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(round, round2) {
		t.Fatal("rounds not equals")
	}
}
