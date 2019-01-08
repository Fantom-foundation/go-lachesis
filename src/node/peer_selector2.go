package node

import (
	"math"
	"math/rand"

	"github.com/Fantom-foundation/go-lachesis/src/peers"
)

// PeerSelector provides an interface for the lachesis node to
// update the last peer it gossiped with and select the next peer
// to gossip with
// type PeerSelector interface {
// 	PeerSet() *peers.PeerSet
// 	UpdateLast(peer string)
// 	Next() *peers.Peer
// }

// +++++++++++++++++++++++++++++++++++++++
// Selection based on FlagTable of a randomly chosen undermined event

// SmartPeerSelector flag table based smart selection struct
type SmartPeerSelector struct {
	peerSet      *peers.PeerSet
	localAddr    string
	last         string
	GetFlagTable func() (map[string]int64, error)
}

// NewSmartPeerSelector creates a new smart peer selection struct
func NewSmartPeerSelector(peerSet *peers.PeerSet,
	localAddr string,
	GetFlagTable func() (map[string]int64, error)) *SmartPeerSelector {

	return &SmartPeerSelector{
		localAddr:    localAddr,
		peerSet:      peerSet,
		GetFlagTable: GetFlagTable,
	}
}

// Peers returns all known peers
func (ps *SmartPeerSelector) PeerSet() *peers.PeerSet {
	return ps.peerSet
}

// UpdateLast sets the last peer communicated with (avoid double talk)
func (ps *SmartPeerSelector) UpdateLast(peer string) {
	ps.last = peer
}

// Next returns the next peer based on the flag table cost function selection
func (ps *SmartPeerSelector) Next() *peers.Peer {
	flagTable, err := ps.GetFlagTable()
	if err != nil {
		flagTable = nil
	}


	sortedSrc := ps.peerSet.ToPeerByUsedSlice()
	n := int(2 * len(sortedSrc) / 3 + 1)
	if n < len(sortedSrc) {
		sortedSrc = sortedSrc[0 : n]
	}
	selected := make([]*peers.Peer, len(sortedSrc))
	sCount := 0
	flagged := make([]*peers.Peer, len(sortedSrc))
	fCount := 0
	minUsedIdx := 0
	minUsedVal := int64(math.MaxInt64)
	var lastused []*peers.Peer

	for _, p := range sortedSrc {
		if p.NetAddr == ps.localAddr {
			continue
		}
		if p.NetAddr == ps.last || p.PubKeyHex == ps.last {
			lastused = append(lastused, p)
			continue
		}

		if f, ok := flagTable[p.PubKeyHex]; ok && f == 1 {
			flagged[fCount] = p
			fCount += 1
			continue
		}

		if p.Used < minUsedVal {
			minUsedVal = p.Used
			minUsedIdx = sCount
		}
		selected[sCount] = p
		sCount += 1
	}

	selected = selected[minUsedIdx:sCount]
	if len(selected) < 1 {
		selected = flagged[0:fCount]
	}
	if len(selected) < 1 {
		selected = lastused
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
