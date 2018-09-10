package node

import (
	"math/rand"

	"github.com/andrecronje/lachesis/net"
)

type PeerSelector interface {
	Peers() []net.Peer
	UpdateLast(peer string)
	UpdateLastN(peer []string)
	Next() net.Peer
	NextN(n int)
}

//+++++++++++++++++++++++++++++++++++++++
//RANDOM

type RandomPeerSelector struct {
	peers []net.Peer
	last  string
	lastN []string
}

func NewRandomPeerSelector(participants []net.Peer, localAddr string) *RandomPeerSelector {
	_, peers := net.ExcludePeer(participants, localAddr)
	return &RandomPeerSelector{
		peers: peers,
	}
}

func (ps *RandomPeerSelector) Peers() []net.Peer {
	return ps.peers
}

func (ps *RandomPeerSelector) UpdateLast(peer string) {
	ps.last = peer
}

func (ps *RandomPeerSelector) UpdateLastN(peers []string) {
	ps.lastN = peers
}

func (ps *RandomPeerSelector) Next() net.Peer {
	selectablePeers := ps.peers
	if len(selectablePeers) > 1 {
		_, selectablePeers = net.ExcludePeer(selectablePeers, ps.last)
	}
	i := rand.Intn(len(selectablePeers))
	peer := selectablePeers[i]
	return peer
}

func (ps *RandomPeerSelector) NextN(n int) net.Peer {
	selectablePeers := ps.peers
	if len(selectablePeers) > n {
		_, selectablePeers = net.ExcludePeers(selectablePeers, ps.lastN)
	}
	i := rand.Intn(len(selectablePeers))
	peer := selectablePeers[i]
	return peer
}
