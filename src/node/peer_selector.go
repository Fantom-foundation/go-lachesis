package node

import (
	"math/rand"

	"github.com/Fantom-foundation/go-lachesis/src/peers"
)

// PeerSelector provides an interface for the lachesis node to
// update the last peer it gossiped with and select the next peer
// to gossip with
type PeerSelector interface {
	PeerSet() *peers.PeerSet
	UpdateLast(peer string)
	Next() *peers.Peer
}

//+++++++++++++++++++++++++++++++++++++++
//RANDOM

// RandomPeerSelector is a randomized peer selection struct
type RandomPeerSelector struct {
	peers     *peers.PeerSet
	localAddr string
	last      string
}

// NewRandomPeerSelector creates a new random peer selector
func NewRandomPeerSelector(peerSet *peers.PeerSet, localAddr string) *RandomPeerSelector {
	return &RandomPeerSelector{
		localAddr: localAddr,
		peers:     peerSet,
	}
}

// PeerSet returns all known peers
func (ps *RandomPeerSelector) PeerSet() *peers.PeerSet {
	return ps.peers
}

// UpdateLast sets the last peer communicated with (to avoid double talk)
func (ps *RandomPeerSelector) UpdateLast(peer string) {
	ps.last = peer
}

// Next returns the next randomly selected peer(s) to communicate with
func (ps *RandomPeerSelector) Next() *peers.Peer {
	selectablePeers := ps.peers.Peers

	if len(selectablePeers) > 1 {
		_, selectablePeers = peers.ExcludePeer(selectablePeers, ps.localAddr)

		if len(selectablePeers) > 1 {
			_, selectablePeers = peers.ExcludePeer(selectablePeers, ps.last)
		}
	}

	i := rand.Intn(len(selectablePeers))

	peer := selectablePeers[i]

	return peer
}
