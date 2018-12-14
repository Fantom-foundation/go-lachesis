package poset

import (
	"fmt"
	"reflect"
	"testing"

	cm "github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/peers"
)

func TestParticipantEventsCache(t *testing.T) {
	size := 10
	testSize := int64(25)
	peerSet := peers.NewPeerSet([]*peers.Peer{
		peers.NewPeer("0xaa", ""),
		peers.NewPeer("0xbb", ""),
		peers.NewPeer("0xcc", ""),
	})

	pec := NewParticipantEventsCache(size, peerSet)

	items := make(map[string][]string)
	for pk := range peerSet.ByPubKey {
		items[pk] = []string{}
	}

	for i := int64(0); i < testSize; i++ {
		for pk := range peerSet.ByPubKey {
			item := fmt.Sprintf("%s%d", pk, i)

			pec.Set(pk, item, i)

			pitems := items[pk]
			pitems = append(pitems, item)
			items[pk] = pitems
		}
	}

	// GET ITEM ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	for pk := range peerSet.ByPubKey {

		index1 := int64(9)
		_, err := pec.GetItem(pk, index1)
		if err == nil || !cm.Is(err, cm.TooLate) {
			t.Fatalf("Expected ErrTooLate")
		}

		index2 := int64(15)
		expected2 := items[pk][index2]
		actual2, err := pec.GetItem(pk, index2)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(expected2, actual2) {
			t.Fatalf("expected and cached not equal")
		}

		index3 := int64(27)
		actual3, err := pec.Get(pk, index3)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual([]string{}, actual3) {
			t.Fatalf("expected and cached not equal")
		}
	}

	//KNOWN ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	known := pec.Known()
	for p, k := range known {
		expectedLastIndex := testSize - 1
		if k != expectedLastIndex {
			t.Errorf("Known[%d] should be %d, not %d", p, expectedLastIndex, k)
		}
	}

	//GET ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	for pk := range peerSet.ByPubKey {
		if _, err := pec.Get(pk, 0); err != nil && !cm.Is(err, cm.TooLate) {
			t.Fatalf("Skipping 0 elements should return ErrTooLate")
		}

		skipIndex := int64(9)
		expected := items[pk][skipIndex+1:]
		cached, err := pec.Get(pk, skipIndex)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(expected, cached) {
			t.Fatalf("expected and cached not equal")
		}

		skipIndex2 := int64(15)
		expected2 := items[pk][skipIndex2+1:]
		cached2, err := pec.Get(pk, skipIndex2)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(expected2, cached2) {
			t.Fatalf("expected and cached not equal")
		}

		skipIndex3 := int64(27)
		cached3, err := pec.Get(pk, skipIndex3)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual([]string{}, cached3) {
			t.Fatalf("expected and cached not equal")
		}
	}
}

func TestParticipantEventsCacheEdge(t *testing.T) {
	size := 10
	testSize := int64(11)
	peerSet := peers.NewPeerSet([]*peers.Peer{
		peers.NewPeer("0xaa", ""),
		peers.NewPeer("0xbb", ""),
		peers.NewPeer("0xcc", ""),
	})

	pec := NewParticipantEventsCache(size, peerSet)

	items := make(map[string][]string)
	for pk := range peerSet.ByPubKey {
		items[pk] = []string{}
	}

	for i := int64(0); i < testSize; i++ {
		for pk := range peerSet.ByPubKey {
			item := fmt.Sprintf("%s%d", pk, i)

			pec.Set(pk, item, i)

			pitems := items[pk]
			pitems = append(pitems, item)
			items[pk] = pitems
		}
	}

	for pk := range peerSet.ByPubKey {
		expected := items[pk][size:]
		cached, err := pec.Get(pk, int64(size-1))
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(expected, cached) {
			t.Fatalf("expected (%#v) and cached (%#v) not equal", expected, cached)
		}
	}
}
