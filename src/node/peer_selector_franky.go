package node

import (
	"sort"
	
	"github.com/Fantom-foundation/go-lachesis/src/peers"
)

// FrankyPeerSelector provides selection based on FlagTable of a randomly chosen undermined event
type FrankyPeerSelector struct {
	peers        *peers.Peers
	localAddr    string
	last         string
	GetFlagTable GetFlagTableFn
}

// FrankyPeerSelectorCreationFnArgs specifies which additional arguments are required to create a FrankyPeerSelector
type FrankyPeerSelectorCreationFnArgs struct {
	GetFlagTable GetFlagTableFn
	LocalAddr    string
}

// NewFrankyPeerSelector creates a new smart peer selection struct
func NewFrankyPeerSelector(participants *peers.Peers, args FrankyPeerSelectorCreationFnArgs) *FrankyPeerSelector {

	return &FrankyPeerSelector{
		localAddr:    args.LocalAddr,
		peers:        participants,
		GetFlagTable: args.GetFlagTable,
	}
}

// NewFrankyPeerSelectorWrapper implements SelectorCreationFn to allow dynamic creation of FrankyPeerSelector ie NewNode
func NewFrankyPeerSelectorWrapper(participants *peers.Peers, args interface{}) PeerSelector {
	return NewFrankyPeerSelector(participants, args.(FrankyPeerSelectorCreationFnArgs))
}

// Peers returns all known peers
func (ps *FrankyPeerSelector) Peers() *peers.Peers {
	return ps.peers
}

// UpdateLast sets the last peer communicated with (avoid double talk)
func (ps *FrankyPeerSelector) UpdateLast(peer string) {
	// We need an exclusive access to ps.last for writing;
	// let use peers' lock instead of adding additional lock.
	// ps.last is accessed for read under peers' lock
	ps.peers.Lock()
	defer ps.peers.Unlock()

	ps.last = peer
}

// Next returns the next peer based on the flag table cost function selection
func (ps *FrankyPeerSelector) Next() *peers.Peer {
	ps.peers.Lock()
	defer ps.peers.Unlock()

	sorted := ps.peers.ToPeerSlice()
	sort.Stable(peers.ByNetAddr(sorted))

	n := len(sorted)
	q := 2 * n / 3 + 1
	var next int
	
	idx := sort.Search(n, func(i int) bool { return sorted[i].Message.NetAddr >= ps.localAddr})

	if idx <= q {
		next = (idx + 1) % q
	} else {
		idx = idx - q
		next = q + (idx + 1) % (n - q)
	}
	
	return sorted[next]
}
