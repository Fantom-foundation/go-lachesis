package peers

import (
	"encoding/hex"

	"github.com/andrecronje/lachesis/src/common"
)

const (
	jsonPeerPath = "peers.json"
)

type Peer struct {
	ID        int `json:"-"`
	NetAddr   string
	PubKeyHex string
}

func NewPeer(pubKeyHex, netAddr string) *Peer {
	peer := &Peer{
		PubKeyHex: pubKeyHex,
		NetAddr:   netAddr,
	}

	peer.computeID()

	return peer
}

func (p *Peer) PubKeyBytes() ([]byte, error) {
	return hex.DecodeString(p.PubKeyHex[2:])
}

func (p *Peer) computeID() error {
	// TODO: Use the decoded bytes from hex
	pubKey, err := p.PubKeyBytes()

	if err != nil {
		return err
	}

	p.ID = common.Hash32(pubKey)

	return nil
}

// PeerStore provides an interface for persistent storage and
// retrieval of peers.
type PeerStore interface {
	// Peers returns the list of known peers.
	Peers() (*Peers, error)

	// SetPeers sets the list of known peers. This is invoked when a peer is
	// added or removed.
	SetPeers([]*Peer) error
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// ExcludePeers is used to exclude one or more peers from a list of peers.
func ExcludePeers(peers []*Peer, exclusionPeers []string) (int, []*Peer) {
	index := -1
	otherPeers := make([]*Peer, 0, len(peers))
	for i, p := range peers {
		if !stringInSlice(p.NetAddr, exclusionPeers) && !stringInSlice(p.PubKeyHex, exclusionPeers) {
			otherPeers = append(otherPeers, p)
		} else {
			index = i
		}
	}
	return index, otherPeers
}
