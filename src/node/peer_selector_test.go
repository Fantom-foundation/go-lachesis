package node

import (
	"fmt"

	"github.com/Fantom-foundation/go-lachesis/src/crypto"
	"github.com/Fantom-foundation/go-lachesis/src/peers"
)

/*
 * stuff
 */

func clonePeers(src *peers.PeerSet) *peers.PeerSet {
	dst := peers.NewPeerSet(src.Peers)
	return dst
}

func fakePeers(n int) *peers.PeerSet {
	var participants []*peers.Peer
	for i := 0; i < n; i++ {
		key, _ := crypto.GenerateECDSAKey()
		peer := peers.Peer{
			NetAddr:   fakeAddr(i),
			PubKeyHex: fmt.Sprintf("0x%X", crypto.FromECDSAPub(&key.PublicKey)),
		}
		participants = append(participants, &peer)
	}
	return peers.NewPeerSet(participants)
}

func fakeAddr(i int) string {
	return fmt.Sprintf("addr%d", i)
}
