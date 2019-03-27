package node

import (
	"math"
	"math/rand"

	"github.com/Fantom-foundation/go-lachesis/src/peers"
)

const EPS float64 = 1e-32

// UnfairPeerSelector provides selection to prevent lazy node creation
type UnfairPeerSelector struct {
	// kPeerSize uint64
	last      string
	localAddr string
	peers     *peers.Peers
}

// UnfairPeerSelectorCreationFnArgs specifies which additional arguments are require to create a UnfairPeerSelector
type UnfairPeerSelectorCreationFnArgs struct {
	KPeerSize uint64
	LocalAddr string
}

// NewUnfairPeerSelector creates a new fair peer selection struct
func NewUnfairPeerSelector(participants *peers.Peers, args UnfairPeerSelectorCreationFnArgs) *UnfairPeerSelector {
	return &UnfairPeerSelector{
		localAddr: args.LocalAddr,
		peers:     participants,
		// kPeerSize: args.KPeerSize,
	}
}

// NewUnfairPeerSelectorWrapper implements SelectorCreationFn to allow dynamic creation of UnfairPeerSelector ie NewNode
func NewUnfairPeerSelectorWrapper(participants *peers.Peers, args interface{}) PeerSelector {
	return NewUnfairPeerSelector(participants, args.(UnfairPeerSelectorCreationFnArgs))
}

// Peers returns all known peers
func (ps *UnfairPeerSelector) Peers() *peers.Peers {
	return ps.peers
}

// UpdateLast sets the last peer communicated with (avoid double talk)
func (ps *UnfairPeerSelector) UpdateLast(peer string) {
	// We need exclusive access to ps.last for writing;
	// let use peers' lock instead of adding an additional lock.
	// ps.last is accessed for read under peers' lock
	ps.peers.Lock()
	defer ps.peers.Unlock()

	ps.last = peer
}

//func fairCostFunction(peer *peers.Peer) float64 {
//	if peer.GetHeight() == 0 {
//		return 0
//	}
//	return float64(peer.GetInDegree()) / float64(2 + peer.GetHeight())
//}

// Next returns the next peer based on the work cost function selection
func (ps *UnfairPeerSelector) Next() *peers.Peer {
	// Maximum number of peers to select/return. In case configurable KPeerSize is implemented.
	// maxPeers := ps.kPeerSize
	// if maxPeers == 0 {
	// 	maxPeers = 1
	// }

	ps.peers.Lock()
	defer ps.peers.Unlock()

	sortedSrc := ps.peers.ToPeerByUsedSlice()
	var lastUsed []*peers.Peer

	maxCost := math.Inf(-1)
	selected := make([]*peers.Peer, 0)
	for _, p := range sortedSrc {
		if p.Message.NetAddr == ps.localAddr {
			continue
		}
		if p.Message.NetAddr == ps.last || p.Message.PubKeyHex == ps.last {
			lastUsed = append(lastUsed, p)
			continue
		}

		cost := fairCostFunction(p)
		if math.Abs(maxCost - cost) < EPS {
			selected = append(selected, p)
		} else if maxCost < cost {
			maxCost = cost
			selected = make([]*peers.Peer, 1)
			selected[0] = p
		}

	}

	if len(selected) < 1 {
		selected = lastUsed
	}
	if len(selected) == 1 {
		selected[0].Used++
		return selected[0]
	}
	if len(selected) < 1 {
		return nil
	}

	i := rand.Intn(len(selected))
	selected[i].Used++
	return selected[i]
}
