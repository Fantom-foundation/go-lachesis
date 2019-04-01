package peers

import (
	"fmt"
	"io/ioutil"
//	"os"
	"testing"

	"crypto/ecdsa"

	"reflect"

	scrypto "github.com/Fantom-foundation/go-lachesis/src/crypto"
)

func TestJSONPeers(t *testing.T) {
	// Create a test dir
	dir, err := ioutil.TempDir("", "lachesis")
	if err != nil {
		t.Fatalf("err: %v ", err)
	}
//	defer func() {
//		if err := os.RemoveAll(dir); err != nil {
//			t.Fatal(err)
//		}
//	}()

	// Create the store
	store := NewJSONPeers(dir)

	// Try a read, should get nothing
	peers, err := store.GetPeers()
	if err == nil {
		t.Fatalf("store.Peers() should generate an error")
	}
	if peers != nil {
		t.Fatalf("peers: %v", peers)
	}

	keys := map[string]*ecdsa.PrivateKey{}
	newPeers := NewPeers()
	for i := 0; i < 3; i++ {
		key, _ := scrypto.GenerateECDSAKey()
		peer := Peer{
			Message: &PeerMessage {
				NetAddr:   fmt.Sprintf("addr%d", i),
				PubKeyHex: fmt.Sprintf("0x%X", scrypto.FromECDSAPub(&key.PublicKey)),
			},
			weight: 1,
		}
		newPeers.AddPeer(&peer)
		keys[peer.Message.NetAddr] = key
	}

	newPeersSlice := newPeers.ToPeerSlice()

	if err := store.SetPeers(newPeersSlice); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Try a read, should find 3 peers
	peers, err = store.GetPeers()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if peers.Len() != 3 {
		t.Fatalf("peers: %v", peers)
	}

	peersSlice := peers.ToPeerSlice()

	for i := 0; i < 3; i++ {
		if peersSlice[i].Message.NetAddr != newPeersSlice[i].Message.NetAddr {
			t.Fatalf("peers[%d] NetAddr should be %s, not %s", i,
				newPeersSlice[i].Message.NetAddr, peersSlice[i].Message.NetAddr)
		}
		if peersSlice[i].Message.PubKeyHex != newPeersSlice[i].Message.PubKeyHex {
			t.Fatalf("peers[%d] PubKeyHex should be %s, not %s", i,
				newPeersSlice[i].Message.PubKeyHex, peersSlice[i].Message.PubKeyHex)
		}
		pubKeyBytes, err := peersSlice[i].PubKeyBytes()
		if err != nil {
			t.Fatal(err)
		}
		pubKey := scrypto.ToECDSAPub(pubKeyBytes)
		if !reflect.DeepEqual(*pubKey, keys[peersSlice[i].Message.NetAddr].PublicKey) {
			t.Fatalf("peers[%d] PublicKey not parsed correctly", i)
		}
	}
}
