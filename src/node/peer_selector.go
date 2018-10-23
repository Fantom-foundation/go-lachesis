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
	UpdateLast(peers []*peers.Peer)
	Next(n int) []*peers.Peer
}

// RandomPeerSelector is a randomized peer selection struct
type RandomPeerSelector struct {
	peers     *peers.Peers
	localAddr string
	last      []string
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
func (ps *RandomPeerSelector) UpdateLast(peers []*peers.Peer) {
	ps.last = make([]string, len(peers))
	for i, p := range peers {
		ps.last[i] = p.NetAddr
	}
}

// Next returns the next randomly selected peer(s) to communicate with
func (ps *RandomPeerSelector) Next(n int) []*peers.Peer {
	slice := ps.peers.ToPeerSlice()
	var lastused []*peers.Peer
	var selected []*peers.Peer
	for _, p := range  slice {
		if p.NetAddr == ps.localAddr {
			continue
		}
		if contains(ps.last, p.NetAddr) || contains(ps.last, p.PubKeyHex) {
			lastused = append(lastused, p)
			continue
		}
		selected = append(selected, p)
	}
	if len(selected) < n {
		selected = append(selected, lastused...)
	}

	if len(selected) == n {
		return selected
	}

	rand.Shuffle(len(selected), func(i, j int) {
		selected[i], selected[j] = selected[j], selected[i]
	})

	return selected[:n]
}

func contains(stringSlice []string, searchString string) bool {
	for _, value := range stringSlice {
		if value == searchString {
			return true
		}
	}
	return false
}
