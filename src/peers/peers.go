package peers

import (
	"sort"
	"sync"
)

// PubKeyPeers map of peers sorted by public key
type PubKeyPeers map[string]*Peer
// IDPeers map of peers sorted by ID
type IDPeers map[uint32]*Peer
// Listener for listening for new peers joining
type Listener func(*Peer)

// Peers struct for all known peers for this node
type Peers struct {
	sync.RWMutex
	Sorted    []*Peer
	ByPubKey  PubKeyPeers
	ByID      IDPeers
	Listeners []Listener
}

/* Constructors */

// NewPeers creates a new peers struct
func NewPeers() *Peers {
	return &Peers{
		ByPubKey: make(PubKeyPeers),
		ByID:     make(IDPeers),
	}
}

// NewPeersFromSlice create a new peers struct from a subset of peers
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
	p.ByID[peer.ID] = peer
}

// AddPeer adds a peer to the peers struct
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

// RemovePeer removes a peer from the peers struct
func (p *Peers) RemovePeer(peer *Peer) {
	p.Lock()
	defer p.Unlock()

	if _, ok := p.ByPubKey[peer.PubKeyHex]; !ok {
		return
	}

	delete(p.ByPubKey, peer.PubKeyHex)
	delete(p.ByID, peer.ID)

	p.internalSort()
}

// RemovePeerByPubKey removes a peer by their public key
func (p *Peers) RemovePeerByPubKey(pubKey string) {
	p.RemovePeer(p.ByPubKey[pubKey])
}

// RemovePeerByID removes a peer based on their ID
func (p *Peers) RemovePeerByID(id uint32) {
	p.RemovePeer(p.ByID[id])
}

/* ToSlice Methods */

// ToPeerSlice returns a slice of peers sorted
func (p *Peers) ToPeerSlice() []*Peer {
	return p.Sorted
}

// ToPeerSlice TODO remove
func (ps *PeerSet) ToPeerSlice() []*Peer {
	return ps.Peers
}

// Sorted TODO remove
func (ps *PeerSet) Sorted() []*Peer {
	return ps.Peers
}

// ToPeerByUsedSlice sorted peers list
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

// ToPubKeySlice peers struct by public key
func (p *Peers) ToPubKeySlice() []string {
	p.RLock()
	defer p.RUnlock()

	res := []string{}

	for _, peer := range p.Sorted {
		res = append(res, peer.PubKeyHex)
	}

	return res
}

// ToIDSlice peers struct by ID
func (p *Peers) ToIDSlice() []uint32 {
	p.RLock()
	defer p.RUnlock()

	res := []uint32{}

	for _, peer := range p.Sorted {
		res = append(res, peer.ID)
	}

	return res
}

// TODO replace / remove
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

// OnNewPeer on new peer joined event trigger listener
func (p *Peers) OnNewPeer(cb func(*Peer)) {
	p.Listeners = append(p.Listeners, cb)
}
// EmitNewPeer emits an event for all listeners as soon as a peer joins
func (p *Peers) EmitNewPeer(peer *Peer) {
	for _, listener := range p.Listeners {
		listener(peer)
	}
}


/* Utilities */

// Len returns the length of peers
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

// ByID sorted by ID peers list
type ByID []*Peer

func (a ByID) Len() int      { return len(a) }
func (a ByID) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool {
	ai := a[i].ID
	aj := a[j].ID
	return ai < aj
}

// ByUsed TODO
type ByUsed []*Peer

func (a ByUsed) Len() int      { return len(a) }
func (a ByUsed) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByUsed) Less(i, j int) bool {
	ai := a[i].Used
	aj := a[j].Used
	return ai > aj
}
