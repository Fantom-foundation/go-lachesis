package node

import (
	"math/rand"

	"github.com/andrecronje/lachesis/net"
)

type PeerSelector interface {
	Peers() []net.Peer
	UpdateLast(peer string)
	UpdateLastN(peer []net.Peer)
	Next() net.Peer
	NextN(n int) []net.Peer
}

//+++++++++++++++++++++++++++++++++++++++
//RANDOM

type RandomPeerSelector struct {
	peers []net.Peer
	last  string
	lastN []net.Peer
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

func (ps *RandomPeerSelector) UpdateLastN(peers []net.Peer) {
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

func (ps *RandomPeerSelector) NextN(n int) []net.Peer {
	selectablePeers := ps.peers
	if len(selectablePeers) > (n*2) {
		selectablePeers = net.ExcludePeers(selectablePeers, ps.lastN)
	}
	//do n times
	if len(selectablePeers) > n {
		returnPeers := make([]net.Peer, 0, n)
		for i := 0; i < n; i++ {
			found := false
			added := false
			for (!added) {
				randomPeer := selectablePeers[rand.Intn(len(selectablePeers))]
				for _, p := range returnPeers {
					if (p.NetAddr == randomPeer.NetAddr) {
						found = true
						break
					}
				}
				if (!found) {
					returnPeers = append(returnPeers, randomPeer)
					added = true
				}
			}
		}
		return returnPeers
	} else {
		return selectablePeers
	}
}
