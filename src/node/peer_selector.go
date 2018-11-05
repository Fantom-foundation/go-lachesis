package node

import (
	"math/rand"
	"time"

	"github.com/andrecronje/lachesis/src/peers"
)

// PeerSelector provides an interface for the lachesis node to 
// update the last peer it gossiped with and select the next peer
// to gossip with 
type PeerSelector interface {
	Peers() *peers.Peers
	UpdateLast(peers []string)
	Next()  []*peers.Peer
}

//+++++++++++++++++++++++++++++++++++++++
//RANDOM

type RandomPeerSelector struct {
	peers     *peers.Peers
	localAddr string
	last      []string
	nextSize  int
}

func NewRandomPeerSelector(participants *peers.Peers, localAddr string, nextPeerCount int) *RandomPeerSelector {
	nextSize := 1
	if nextPeerCount > 1 {
		nextSize = nextPeerCount
	}
	return &RandomPeerSelector{
		localAddr: localAddr,
		peers:     participants,
		nextSize:  nextSize,
	}
}

func (ps *RandomPeerSelector) Peers() *peers.Peers {
	return ps.peers
}

func (ps *RandomPeerSelector) UpdateLast(peer []string) {
	ps.last = peer
}

func (ps *RandomPeerSelector) Next() []*peers.Peer {
	selectablePeers := ps.peers.ToPeerSlice()
	currentLength := func () int {
		return len(selectablePeers)
	}

	if currentLength() > 1 {
		_, selectablePeers = peers.ExcludePeers(selectablePeers, []string{ps.localAddr})

		if currentLength() > 1 {
			_, selectablePeers = peers.ExcludePeers(selectablePeers, ps.last)
		}
	}

	available := currentLength()
	if available > ps.nextSize {
		available = ps.nextSize
	}

	peers := make([]*peers.Peer, 0, available)
	rand.Seed(time.Now().UnixNano())
	randomUniqueNumbers := rand.Perm(available)
	for _, i := range randomUniqueNumbers[:available] {
		peers = append(peers, selectablePeers[i])
	}
	//fmt.Printf("Selected Peer: Random... [count: %d] [random: %v] [peers: %v]\n", available, randomUniqueNumbers, peers)

	return peers
}
