package node

import (
	"math/rand"

	"github.com/Fantom-foundation/go-lachesis/src/peers"
)

// PeerSelector provides an interface for the lachesis node to 
// update the last peer it gossiped with and select the next peer
// to gossip with 
//type PeerSelector interface {
//	PeerSet() *peers.PeerSet
//	UpdateLast(peer string)
//	Next() *peers.Peer
//}

//+++++++++++++++++++++++++++++++++++++++
//Selection based on FlagTable of a randomly chosen undermined event

type SmartPeerSelector struct {
	peerSet      *peers.PeerSet
	localAddr    string
	last         string
	GetFlagTable func() (map[string]int64, error)
}

func NewSmartPeerSelector(peerSet *peers.PeerSet,
	localAddr string,
	GetFlagTable func() (map[string]int64, error)) *SmartPeerSelector {

	return &SmartPeerSelector{
		localAddr:    localAddr,
		peerSet:      peerSet,
		GetFlagTable: GetFlagTable,
	}
}

func (ps *SmartPeerSelector) PeerSet() *peers.PeerSet {
	return ps.peerSet
}

func (ps *SmartPeerSelector) UpdateLast(peer string) {
	ps.last = peer
}

func (ps *SmartPeerSelector) Next() *peers.Peer {
	//ps.peerSet.Lock()
	//defer ps.peerSet.Unlock()
	selectablePeers := ps.peerSet.ToPeerByUsedSlice() //[1:]
	if len(selectablePeers) > 1 {
		_, selectablePeers = peers.ExcludePeer(selectablePeers, ps.localAddr)
		if len(selectablePeers) > 1 {
			_, selectablePeers = peers.ExcludePeer(selectablePeers, ps.last)
			if len(selectablePeers) > 1 {
				var k int64
				minUsed := selectablePeers[len(selectablePeers) - 1].Used
				for k = 0; selectablePeers[k].Used > minUsed; k++ {}
				selectablePeers = selectablePeers[k:]
				if ft, err := ps.GetFlagTable(); err == nil {
					for id, flag := range ft {
						if flag == 1 && len(selectablePeers) > 1 {
							peers.ExcludePeer(selectablePeers, id)
						}
					}
				}
			}
		}
	}
	i := rand.Intn(len(selectablePeers))
	selectablePeers[i].Used++;
	return selectablePeers[i]
}

