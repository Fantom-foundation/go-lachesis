package node

import (
	"math/rand"

	"github.com/Fantom-foundation/go-lachesis/src/peers"
)

// PeerSelector provides an interface for the lachesis node to
// update the last peer it gossiped with and select the next peer
// to gossip with
type PeerSelector interface {
	Peers() *peers.Peers
	UpdateLast(peer string)
	UpdateLastById(id int64)
	Next() peers.Peer
	ToPeerSlice() []peers.Peer
	Len() int
}

//+++++++++++++++++++++++++++++++++++++++
//RANDOM

// RandomPeerSelector is a randomized peer selection struct
type RandomPeerSelector struct {
	peers     *peers.Peers
	localAddr string
	last      string
}

// NewRandomPeerSelector creates a new random peer selector
func NewRandomPeerSelector(participants *peers.Peers, localAddr string) *RandomPeerSelector {
	return &RandomPeerSelector{
		localAddr: localAddr,
		peers:     participants,
	}
}

// Peers returns all known peers
func (ps *RandomPeerSelector) Peers() *peers.Peers {
	return ps.peers
}

// UpdateLast sets the last peer communicated with (to avoid double talk)
func (ps *RandomPeerSelector) UpdateLast(peer string) {
	ps.last = peer
}

func (ps *RandomPeerSelector) UpdateLastById(id int64) {
	peer, ok := ps.Peers().GetById(id)
	if ok {
		ps.last = peer.NetAddr
	}
}

// Next returns the next randomly selected peer(s) to communicate with
func (ps *RandomPeerSelector) Next() peers.Peer {
	selectablePeers := ps.peers.ToPeerSlice()

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

func (ps *RandomPeerSelector) Len() int {
	return ps.peers.Len()
}

func (ps *RandomPeerSelector) ToPeerSlice() []peers.Peer {
	return ps.peers.ToPeerSlice()
}
