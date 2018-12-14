package peers

import (
	"fmt"
	"math"

	"github.com/Fantom-foundation/go-lachesis/src/crypto"
)

// XXX exclude peers should be in here

/* Constructors */

func NewEmptyPeerSet() *PeerSet {
	return &PeerSet{
		ByPubKey: make(map[string]*Peer),
		ByID:     make(map[uint32]*Peer),
	}
}

// NewPeerSet creates a new PeerSet from a list of PeerSet
func NewPeerSet(peers []*Peer) *PeerSet {
	peerSet := NewEmptyPeerSet()
	for _, peer := range peers {
		if peer.ID == 0 {
			peer.ComputeID()
		}

		peerSet.ByPubKey[peer.PubKeyHex] = peer
		peerSet.ByID[peer.ID] = peer
	}

	peerSet.Peers = peers

	return peerSet
}

func (ps *PeerSet) Lock() {

}

func (ps *PeerSet) Unlock() {

}

// WithNewPeer returns a new PeerSet with a list of peers including the new one.
func (ps *PeerSet) WithNewPeer(peer *Peer) *PeerSet {
	peers := append(ps.Peers, peer)
	newPeerSet := NewPeerSet(peers)
	return newPeerSet
}

// WithRemovedPeer returns a new PeerSet with a list of peers exluding the
// provided one
func (ps *PeerSet) WithRemovedPeer(peer *Peer) *PeerSet {
	var peers []*Peer
	for _, p := range ps.Peers {
		if p.PubKeyHex != peer.PubKeyHex {
			peers = append(peers, p)
		}
	}
	newPeerSet := NewPeerSet(peers)
	return newPeerSet
}

/* ToSlice Methods */

// PubKeys returns the PeerSet's slice of public keys
func (ps *PeerSet) PubKeys() []string {
	var res []string

	for _, peer := range ps.Peers {
		res = append(res, peer.PubKeyHex)
	}

	return res
}

// IDs returns the PeerSet's slice of IDs
func (ps *PeerSet) IDs() []uint32 {
	var res []uint32

	for _, peer := range ps.Peers {
		res = append(res, peer.ID)
	}

	return res
}

/* Utilities */

// Len returns the number of PeerSet in the PeerSet
func (ps *PeerSet) Len() int {
	return len(ps.ByPubKey)
}

// Hash uniquely identifies a PeerSet. It is computed by sorting the peers set
// by ID, and hashing (SHA256) their public keys together, one by one.
func (ps *PeerSet) Hash() ([]byte, error) {
	if len(ps.Hash_) == 0 {
		var hash []byte
		for _, p := range ps.Peers {
			pk, _ := p.PubKeyBytes()
			hash = crypto.SimpleHashFromTwoHashes(hash, pk)
		}
		ps.Hash_ = hash
	}
	return ps.Hash_, nil
}

// Hex is the hexadecimal representation of Hash
func (ps *PeerSet) Hex() string {
	if len(ps.Hex_) == 0 {
		hash, _ := ps.Hash()
		ps.Hex_ = fmt.Sprintf("0x%X", hash)
	}
	return ps.Hex_
}

// SuperMajority return the number of peers that forms a strong majortiy (+2/3)
// in the PeerSet
func (ps *PeerSet) SuperMajority() int64 {
	if ps.SuperMajority_ == 0 {
		val := int64(2*ps.Len()/3 + 1)
		ps.SuperMajority_ = val
	}
	return ps.SuperMajority_
}

func (ps *PeerSet) TrustCount() int64 {
	if ps.TrustCount_ == 0 {
		val := int64(math.Ceil(float64(ps.Len()) / float64(3)))
		ps.TrustCount_ = val
	}
	return ps.TrustCount_
}

func (ps *PeerSet) clearCache() {
	ps.Hash_ = []byte{}
	ps.Hex_ = ""
	ps.SuperMajority_ = 0
}
