package peers

import (
	"sort"
	"sync"
)

// PubKeyPeers map of peers sorted by public key
type PubKeyPeers map[string]*Peer
// IDPeers map of peers sorted by ID
type IDPeers map[int64]*Peer
// Listener for listening for new peers joining
type Listener func(*Peer)

// Peers struct for all known peers for this node
type Peers struct {
	sync.RWMutex
	sorted    []*Peer
	byPubKey  PubKeyPeers
	byId      IDPeers
	listeners []Listener
}

/* Constructors */

// NewPeers creates a new peers struct
func NewPeers() *Peers {
	return &Peers{
		byPubKey: make(PubKeyPeers),
		byId:     make(IDPeers),
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
		peer.computeID()
	}

	p.byPubKey[peer.PubKeyHex] = peer
	p.byId[peer.ID] = peer
}

func (p *Peers) internalSort() {
	res := []*Peer{}

	for _, p := range p.byPubKey {
		res = append(res, p)
	}

	sort.Sort(ByID(res))

	p.sorted = res
}

// AddPeer adds a peer to the peers struct
func (p *Peers) AddPeer(peer *Peer) {
	p.Lock()
	p.addPeerRaw(peer)
	p.internalSort()
	p.Unlock()
	p.EmitNewPeer(peer)
}

/* Remove Methods */

// RemovePeer removes a peer from the peers struct
func (p *Peers) RemovePeer(peer *Peer) {
	p.Lock()
	defer p.Unlock()

	if _, ok := p.byPubKey[peer.PubKeyHex]; !ok {
		return
	}

	delete(p.byPubKey, peer.PubKeyHex)
	delete(p.byId, peer.ID)

	p.internalSort()
}

// RemovePeerByPubKey removes a peer by their public key
func (p *Peers) RemovePeerByPubKey(pubKey string) {
	p.Lock()
	defer p.Unlock()

	peer, ok := p.byPubKey[pubKey]
	if !ok {
		return
	}
	peerID := peer.ID

	delete(p.byPubKey, pubKey)
	delete(p.byId, peerID)

	p.internalSort()
}

// RemovePeerByID removes a peer based on their ID
func (p *Peers) RemovePeerById(id int64) {
	p.Lock()
	defer p.Unlock()

	peer, ok := p.byId[id]
	if !ok {
		return
	}
	pubKey := peer.PubKeyHex

	delete(p.byPubKey, pubKey)
	delete(p.byId, id)

	p.internalSort()
}

func (p *Peers) GetByPubKeys() PubKeyPeers {
	p.RLock()
	defer p.RUnlock()

	res := PubKeyPeers{}

	for k, p := range p.byPubKey {
		res[k] = p
	}

	return res
}

func (p *Peers) GetByPubKey(PubKey string) (Peer, bool) {
	p.RLock()
	defer p.RUnlock()
	peer, ok := p.byPubKey[PubKey]
	if peer == nil {
		return Peer{}, false
	}
	return *peer, ok
}

func (p *Peers) GetByIds() IDPeers {
	p.RLock()
	defer p.RUnlock()

	res := IDPeers{}

	for k, p := range p.byId {
		res[k] = p
	}

	return res
}

func (p *Peers) GetById(Id int64) (Peer, bool) {
	p.RLock()
	defer p.RUnlock()
	peer, ok := p.byId[Id]
	return *peer, ok
}

/* ToSlice Methods */

// ToPeerSlice returns a slice of peers sorted
func (p *Peers) ToPeerSlice() []Peer {
	p.RLock()
	defer p.RUnlock()

	res := []Peer{}
	for _, p := range p.sorted {
		res = append(res, *p)
	}
	return res
}

// ToPeerByUsedSlice sorted peers list
func (p *Peers) ToPeerByUsedSlice() []Peer {
	p.RLock()
	defer p.RUnlock()

	res := []Peer{}

	for _, p := range p.byPubKey {
		res = append(res, *p)
	}

	sort.Sort(ByUsed(res))
	return res
}

// ToPubKeySlice peers struct by public key
func (p *Peers) ToPubKeySlice() []string {
	p.RLock()
	defer p.RUnlock()

	res := []string{}

	for _, peer := range p.sorted {
		res = append(res, peer.PubKeyHex)
	}

	return res
}

// ToIDSlice peers struct by ID
func (p *Peers) ToIDSlice() []int64 {
	p.RLock()
	defer p.RUnlock()

	res := []int64{}

	for _, peer := range p.sorted {
		res = append(res, peer.ID)
	}

	return res
}

/* EventListener */

// OnNewPeer on new peer joined event trigger listener
func (p *Peers) OnNewPeer(cb func(*Peer)) {
	p.Lock()
	defer p.Unlock()

	p.listeners = append(p.listeners, cb)
}

// EmitNewPeer emits an event for all listeners as soon as a peer joins
func (p *Peers) EmitNewPeer(peer *Peer) {
	p.RLock()
	defer p.RUnlock()

	for _, listener := range p.listeners {
		listener(peer)
	}
}


/* Utilities */

// Len returns the length of peers
func (p *Peers) Len() int {
	p.RLock()
	defer p.RUnlock()

	return len(p.sorted)
}

func (p *Peers) IncUsed(id int64) {
	p.Lock()
	defer p.Unlock()
	p.byId[id].Used++
}

// ByPubHex implements sort.Interface for Peers based on
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
type ByUsed []Peer

func (a ByUsed) Len() int      { return len(a) }
func (a ByUsed) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByUsed) Less(i, j int) bool {
	ai := a[i].Used
	aj := a[j].Used
	return ai > aj
}
