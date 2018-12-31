package peers

import (
	"sort"
)

// PubKeyPeers map of peers sorted by public key
type PubKeyPeers map[string]*Peer
// IDPeers map of peers sorted by ID
type IDPeers map[uint32]*Peer
// Listener for listening for new peers joining
type Listener func(*Peer)


// ToPeerByUsedSlice converts PeerSet ByPubKey to a slice of peers
func (ps *PeerSet) ToPeerByUsedSlice() []*Peer {
	var res []*Peer

	for _, p := range ps.ByPubKey {
		res = append(res, p)
	}

	sort.Sort(ByUsed(res))
	return res
}

// ToIDSlice convers PeerSet to a slice of Peer id's
func (ps *PeerSet) ToIDSlice() []uint32 {
	var res []uint32

	for _, peer := range ps.Peers {
		res = append(res, peer.ID)
	}

	return res
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
