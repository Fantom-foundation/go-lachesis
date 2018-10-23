package node

import (
	"math"
	"math/rand"

	"github.com/Fantom-foundation/go-lachesis/src/peers"
)

// SmartPeerSelector provides selection based on FlagTable of a randomly chosen undermined event
type SmartPeerSelector struct {
	peers        *peers.Peers
	localAddr    string
	last         []string
	GetFlagTable func() (map[string]int64, error)
}

// NewSmartPeerSelector creates a new smart peer selection struct
func NewSmartPeerSelector(participants *peers.Peers,
	localAddr string,
	GetFlagTable func() (map[string]int64, error)) *SmartPeerSelector {

	return &SmartPeerSelector{
		localAddr:    localAddr,
		peers:        participants,
		GetFlagTable: GetFlagTable,
	}
}

// Peers returns all known peers
func (ps *SmartPeerSelector) Peers() *peers.Peers {
	return ps.peers
}

// UpdateLast sets the last peer communicated with (avoid double talk)
func (ps *SmartPeerSelector) UpdateLast(peers []*peers.Peer) {
	ps.last = make([]string, len(peers))
	for i, p := range peers {
		ps.last[i] = p.NetAddr
	}
}

// Next returns the next peer based on the flag table cost function selection
func (ps *SmartPeerSelector) Next(n int) []*peers.Peer {
	flagTable, err := ps.GetFlagTable()
	if err != nil {
		flagTable = nil
	}

	ps.peers.Lock()
	defer ps.peers.Unlock()

	sortedSrc := ps.peers.ToPeerByUsedSlice()
	m := int(2*len(sortedSrc)/3 + 1)
	if m < len(sortedSrc) && m < n {
		sortedSrc = sortedSrc[0:m]
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
		if contains(ps.last, p.NetAddr) || contains(ps.last, p.PubKeyHex) {
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
	if len(selected) < n {
		selected = flagged[0:fCount]
	}
	if len(selected) < n {
		selected = lastused
	}
	// if len(selected) == n {
	// 	selected[0].Used++
	// 	return selected[0]
	// }
	if len(selected) < n {
		return nil
	}

	rand.Shuffle(len(selected), func(i, j int) {
		selected[i], selected[j] = selected[j], selected[i]
	})
	for i := range selected[:n] {
		selected[i].Used++
	}

	return selected[:n]
}
