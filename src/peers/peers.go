package peers

import (
	"sort"
	"sync"
)

type PubKeyPeers map[string]*Peer
type IdPeers map[uint32]*Peer
type Listener func(*Peer)

type Peers struct {
	sync.RWMutex
	Sorted    []*Peer
	ByPubKey  PubKeyPeers
	ById      IdPeers
	Listeners []Listener
}

/* Constructors */

func NewPeers() *Peers {
	return &Peers{
		ByPubKey: make(PubKeyPeers),
		ById:     make(IdPeers),
	}
}

func NewPeersFromSlice(source []*Peer) *Peers {
	peers := NewPeers()

	for _, peer := range source {
		peers.addPeerRaw(peer)
	}

	peers.internalSort()

	return peers
}

/* Add Methods */

// Add a peer without sorting the set.
// Useful for adding a bunch of peers at the same time
// This method is private and is not protected by mutex.
// Handle with care
func (p *Peers) addPeerRaw(peer *Peer) {
	if peer.ID == 0 {
		peer.ComputeID()
	}

	p.ByPubKey[peer.PubKeyHex] = peer
	p.ById[peer.ID] = peer
}

func (p *Peers) AddPeer(peer *Peer) {
	p.Lock()
	p.addPeerRaw(peer)
	p.internalSort()
	p.Unlock()
 	p.EmitNewPeer(peer)
}

func (p *Peers) internalSort() {
	res := []*Peer{}

	for _, p := range p.ByPubKey {
		res = append(res, p)
	}

	sort.Sort(ByID(res))

	p.Sorted = res
}

/* Remove Methods */

func (p *Peers) RemovePeer(peer *Peer) {
	p.Lock()
	defer p.Unlock()

	if _, ok := p.ByPubKey[peer.PubKeyHex]; !ok {
		return
	}

	delete(p.ByPubKey, peer.PubKeyHex)
	delete(p.ById, peer.ID)

	p.internalSort()
}

func (p *Peers) RemovePeerByPubKey(pubKey string) {
	p.RemovePeer(p.ByPubKey[pubKey])
}

func (p *Peers) RemovePeerById(id uint32) {
	p.RemovePeer(p.ById[id])
}

/* ToSlice Methods */

func (p *Peers) ToPeerSlice() []*Peer {
	return p.Sorted
}

func (ps *PeerSet) ToPeerSlice() []*Peer {
	return ps.Peers
}

func (ps *PeerSet) Sorted() []*Peer {
	return ps.Peers
}

func (p *Peers) ToPeerByUsedSlice() []*Peer {
	res := []*Peer{}

	for _, p := range p.ByPubKey {
		res = append(res, p)
	}

	sort.Sort(ByUsed(res))
	return res
}

func (ps *PeerSet) ToPeerByUsedSlice() []*Peer {
	res := []*Peer{}

	for _, p := range ps.ByPubKey {
		res = append(res, p)
	}

	sort.Sort(ByUsed(res))
	return res
}

func (p *Peers) ToPubKeySlice() []string {
	p.RLock()
	defer p.RUnlock()

	res := []string{}

	for _, peer := range p.Sorted {
		res = append(res, peer.PubKeyHex)
	}

	return res
}

func (p *Peers) ToIDSlice() []uint32 {
	p.RLock()
	defer p.RUnlock()

	res := []uint32{}

	for _, peer := range p.Sorted {
		res = append(res, peer.ID)
	}

	return res
}

func (ps *PeerSet) ToIDSlice() []uint32 {
	//p.RLock()
	//defer p.RUnlock()

	res := []uint32{}

	for _, peer := range ps.Peers {
		res = append(res, peer.ID)
	}

	return res
}

/* EventListener */

func (p *Peers) OnNewPeer(cb func(*Peer)) {
	p.Listeners = append(p.Listeners, cb)
}
func (p *Peers) EmitNewPeer(peer *Peer) {
	for _, listener := range p.Listeners {
		listener(peer)
	}
}


/* Utilities */

func (p *Peers) Len() int {
	p.RLock()
	defer p.RUnlock()

	return len(p.ByPubKey)
}

// ByPubHex implements sort.Interface for PeerSet based on
// the PubKeyHex field.
type ByPubHex []*Peer

func (a ByPubHex) Len() int      { return len(a) }
func (a ByPubHex) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByPubHex) Less(i, j int) bool {
	ai := a[i].PubKeyHex
	aj := a[j].PubKeyHex
	return ai < aj
}

type ByID []*Peer

func (a ByID) Len() int      { return len(a) }
func (a ByID) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool {
	ai := a[i].ID
	aj := a[j].ID
	return ai < aj
}

type ByUsed []*Peer

func (a ByUsed) Len() int      { return len(a) }
func (a ByUsed) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByUsed) Less(i, j int) bool {
	ai := a[i].Used
	aj := a[j].Used
	return ai > aj
}
