package peers

import (
	"fmt"
	"math"

	"github.com/Fantom-foundation/go-lachesis/src/crypto"
)

// XXX exclude peers should be in here

/* Constructors */

// NewEmptyPeerSet PeerSet constructor
func NewEmptyPeerSet() *PeerSet {
	return &PeerSet{
		ByPubKey: make(map[string]*Peer),
		ByID:     make(map[uint64]*Peer),
		Cache:    &CachedValues{},
	}
}

// NewPeerSet creates a new PeerSet from a list of PeerSet
func NewPeerSet(peers []*Peer) *PeerSet {
	peerSet := NewEmptyPeerSet()
	for _, peer := range peers {
		if peer.ID == 0 {
			// TODO: Decide what, if anything, to do with the error
			_ = peer.ComputeID()
		}

		peerSet.ByPubKey[peer.PubKeyHex] = peer
		peerSet.ByID[peer.ID] = peer
	}

	peerSet.Peers = peers

	return peerSet
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
func (ps *PeerSet) IDs() []uint64 {
	var res []uint64

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
	if len(ps.Cache.Hash) == 0 {
		var hash []byte
		for _, p := range ps.Peers {
			pk, _ := p.PubKeyBytes()
			hash = crypto.SimpleHashFromTwoHashes(hash, pk)
		}
		ps.Cache.Hash = hash
	}
	return ps.Cache.Hash, nil
}

// Hex is the hexadecimal representation of Hash
func (ps *PeerSet) Hex() string {
	if len(ps.Cache.Hex) == 0 {
		hash, _ := ps.Hash()
		ps.Cache.Hex = fmt.Sprintf("0x%X", hash)
	}
	return ps.Cache.Hex
}

// SuperMajority return the number of peers that forms a strong majortiy (+2/3)
// in the PeerSet
func (ps *PeerSet) SuperMajority() int64 {
	if ps.Cache.SuperMajority == 0 {
		val := int64(2*ps.Len()/3 + 1)
		ps.Cache.SuperMajority = val
	}
	return ps.Cache.SuperMajority
}

// TrustCount calculates and returns the trust count
func (ps *PeerSet) TrustCount() int64 {
	if ps.Cache.TrustCount == 0 {
		val := int64(math.Ceil(float64(ps.Len()) / float64(3)))
		ps.Cache.TrustCount = val
	}
	return ps.Cache.TrustCount
}
