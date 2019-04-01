package peers

import (
	"math"
	"sort"
	"sync"

	"github.com/Fantom-foundation/go-lachesis/src/common"
)

// PubKeyPeers map of peers sorted by public key
type PubKeyPeers map[string]*Peer

// IDPeers map of peers sorted by ID
type IDPeers map[uint64]*Peer

// AddressPeers maps address to peer
type AddressPeers map[common.Address]*Peer

type NetAddrPeers map[string]*Peer

// Listener for listening for new peers joining
type Listener func(*Peer)

// Peers struct for all known peers for this node
type Peers struct {
	sync.RWMutex
	Sorted    []*Peer
	ByPubKey  PubKeyPeers
	ByID      IDPeers
	ByAddress AddressPeers
	ByNetAddr NetAddrPeers
	Listeners []Listener
	Stake     uint64
	SuperMajority uint64
	TrustCount    uint64
}

/* Constructors */

// NewPeers creates a new peers struct
func NewPeers() *Peers {
	return &Peers{
		ByPubKey:  make(PubKeyPeers),
		ByID:      make(IDPeers),
		ByAddress: make(AddressPeers),
		ByNetAddr: make(NetAddrPeers),
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

// NewPeersFromMessageSlice create a new peers struct from a subset of peer messages
func NewPeersFromMessageSlice(source []*PeerMessage) *Peers {
	peers := NewPeers()

	for _, pm := range source {
		peer := NewPeer(pm.PubKeyHex, pm.NetAddr)
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
		if err := peer.computeID(); err != nil {
			panic(err)
		}
	}
	p.Stake = p.Stake + peer.GetWeight()
	p.SuperMajority = uint64(2 * p.Stake / 3 + 1)
	p.TrustCount = uint64(math.Ceil(float64(p.Stake) / float64(3)))

	p.ByPubKey[peer.Message.PubKeyHex] = peer
	p.ByID[peer.ID] = peer
	p.ByAddress[peer.Address()] = peer
	p.ByNetAddr[peer.Message.NetAddr] = peer
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

	if _, ok := p.ByPubKey[peer.Message.PubKeyHex]; !ok {
		return
	}

	delete(p.ByPubKey, peer.Message.PubKeyHex)
	delete(p.ByID, peer.ID)
	delete(p.ByAddress, peer.Address())
	delete(p.ByNetAddr, peer.Message.NetAddr)

	p.internalSort()
}

// RemovePeerByPubKey removes a peer by their public key
func (p *Peers) RemovePeerByPubKey(pubKey string) {
	p.RemovePeer(p.ByPubKey[pubKey])
}

// RemovePeerByID removes a peer based on their ID
func (p *Peers) RemovePeerByID(id uint64) {
	p.RemovePeer(p.ByID[id])
}

/* ToSlice Methods */

// ToPeerSlice returns a slice of peers sorted
func (p *Peers) ToPeerSlice() []*Peer {
	return append([]*Peer(nil), p.Sorted...)
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

// ToPubKeySlice peers struct by public key
func (p *Peers) ToPubKeySlice() []string {
	p.RLock()
	defer p.RUnlock()

	res := []string{}

	for _, peer := range p.Sorted {
		res = append(res, peer.Message.PubKeyHex)
	}

	return res
}

// ToIDSlice peers struct by ID
func (p *Peers) ToIDSlice() []uint64 {
	p.RLock()
	defer p.RUnlock()

	res := []uint64{}

	for _, peer := range p.Sorted {
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

func (p *Peers) ReadByPubKey(key string) (Peer, bool) {
	p.RLock()
	defer p.RUnlock()
	peer, ok := p.ByPubKey[key]
	return *peer, ok
}

func (p *Peers) ReadByID(key uint64) (Peer, bool) {
	p.RLock()
	defer p.RUnlock()
	peer, ok := p.ByID[key]
	return *peer, ok
}

func (p *Peers) ReadByAddress(key common.Address) (Peer, bool) {
	p.RLock()
	defer p.RUnlock()
	peer, ok := p.ByAddress[key]
	return *peer, ok
}

func (p *Peers) ReadByNetAddr(key string) (Peer, bool) {
	p.RLock()
	defer p.RUnlock()
	peer, ok := p.ByNetAddr[key]
	return *peer, ok
}

func (p *Peers) SetHeightByPubKeyHex(key string, height int64) {
	p.Lock()
	defer p.Unlock()
	(p.ByPubKey[key]).SetHeight(height)
}

func (p *Peers) GetHeightByPubKeyHex(key string) int64 {
	p.RLock()
	defer p.RUnlock()
	return (p.ByPubKey[key]).GetHeight()
}

func (p *Peers) NextHeightByPubKeyHex(key string) int64 {
	p.Lock()
	defer p.Unlock()
	return (p.ByPubKey[key]).NextHeight()
}

func (p *Peers) SetInDegreeByPubKeyHex(key string, inDegree int64) {
	p.Lock()
	defer p.Unlock()
	(p.ByPubKey[key]).SetInDegree(inDegree)
}

func (p *Peers) IncInDegreeByPubKeyHex(key string) {
	p.Lock()
	defer p.Unlock()
	(p.ByPubKey[key]).IncInDegree()
}

// Set new weight to a peer and recalculate PoS values
func (p *Peers) SetPeerWeight(peer *Peer, w uint64) {
	p.Lock()
	defer p.Unlock()
	oldW := peer.GetWeight()
	p.Stake = p.Stake - oldW
	peer.SetWeight(w)
	p.Stake = p.Stake + w
	p.SuperMajority = uint64(2 * p.Stake / 3 + 1)
	p.TrustCount = uint64(math.Ceil(float64(p.Stake) / float64(3)))
}

// Set new weight to a peer by their public Key
func (p*Peers) SetPeerWeightByPubKey(pubKey string, w uint64) {
	p.SetPeerWeight(p.ByPubKey[pubKey], w)
}

// RemovePeerByID removes a peer based on their ID
func (p *Peers) SetPeerWeightByID(id uint64, w uint64) {
	p.SetPeerWeight(p.ByID[id], w)
}

// GetSuperMajority() return the current value of SuperMajority
func (p *Peers) GetSuperMajority() uint64 {
	p.RLock()
	defer p.RUnlock()
	return p.SuperMajority
}

// GetTrustCount() return the current value of TrustCount
func (p *Peers) GetTrustCount() uint64 {
	p.RLock()
	defer p.RUnlock()
	return p.TrustCount
}

// ByPubHex implements sort.Interface for Peers based on
// the PubKeyHex field.
type ByPubHex []*Peer

func (a ByPubHex) Len() int      { return len(a) }
func (a ByPubHex) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByPubHex) Less(i, j int) bool {
	ai := a[i].Message.PubKeyHex
	aj := a[j].Message.PubKeyHex
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

// ByUsed Sorted by Used count
type ByUsed []*Peer

func (a ByUsed) Len() int      { return len(a) }
func (a ByUsed) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByUsed) Less(i, j int) bool {
	ai := a[i].Used
	aj := a[j].Used
	return ai > aj
}

// ByNetAddr Sorted by Used count
type ByNetAddr []*Peer

func (a ByNetAddr) Len() int      { return len(a) }
func (a ByNetAddr) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByNetAddr) Less(i, j int) bool {
	ai := a[i].Message.NetAddr
	aj := a[j].Message.NetAddr
	return ai > aj
}

