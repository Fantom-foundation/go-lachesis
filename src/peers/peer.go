package peers

import (
	"encoding/hex"
	"sync"

	"github.com/Fantom-foundation/go-lachesis/src/common"
)

const (
	jsonPeerPath = "peers.json"
)

// PeerNIL is used for nil peer id
const PeerNIL uint64 = 0


/* PeerMessage type */

// Equals checks peer messages for equality
func (pm *PeerMessage) Equals(cmp *PeerMessage) bool {
	return pm.NetAddr == cmp.NetAddr &&
		pm.PubKeyHex == cmp.PubKeyHex
}

// PubKeyBytes returns the public key bytes for a peer
func (pm *PeerMessage) PubKeyBytes() ([]byte, error) {
	return hex.DecodeString(pm.PubKeyHex[2:])
}

// Address returns the address for a peerMessage
// TODO: hash of publickey
func (pm *PeerMessage) Address() (a common.Address) {
	bytes, err := pm.PubKeyBytes()
	if err != nil {
		panic(err)
	}
	copy(a[:], bytes)
	return
}

/* Peer type */

type Peer struct {
	sync.RWMutex
	Message   *PeerMessage
	ID        uint64
	Used      int64
	height    int64
	inDegree  int64
	weight    uint64
}

// NewPeer creates a new peer based on public key and network address
func NewPeer(pubKeyHex, netAddr string) *Peer {
	peer := &Peer{
		Message: &PeerMessage{
			PubKeyHex: pubKeyHex,
			NetAddr:   netAddr,
		},
		Used:      0,
		height:    -1,
		inDegree:  0,
	}

	if err := peer.computeID(); err != nil {
		panic(err)
	}

	return peer
}


// Equals checks peers for equality
func (p *Peer) Equals(cmp *Peer) bool {
	return p.ID == cmp.ID && p.Message.Equals(cmp.Message)
}

// PubKeyBytes returns the public key bytes for a peer
func (p *Peer) PubKeyBytes() ([]byte, error) {
	return (p.Message).PubKeyBytes()
//	return hex.DecodeString(p.Message.PubKeyHex[2:])
}

func (p *Peer) computeID() error {
	// TODO: Use the decoded bytes from hex
	pubKey, err := p.PubKeyBytes()

	if err != nil {
		return err
	}

	p.ID = common.Hash64(pubKey)

	return nil
}

// Address returns the address for a peer
func (p *Peer) Address() (a common.Address) {
	return (p.Message).Address()
}

// SetHeight() set the value for the height of the peer
func (p *Peer) SetHeight(height int64) {
	p.Lock()
	defer p.Unlock()
	p.height = height
}

// GetHeight() returns current value of the height of the peer
func (p *Peer) GetHeight() int64 {
	p.RLock()
	defer p.RUnlock()
	return p.height
}

// NextHeight() increase the current height by 1 and returns new value
func (p *Peer) NextHeight() int64 {
	p.Lock()
	defer p.Unlock()
	p.height++
	return p.height
}

// SetInDegree() set the value of the inDegree of the peer
func (p *Peer) SetInDegree(inDegree int64) {
	p.Lock()
	defer p.Unlock()
	p.inDegree = inDegree
}

// GetInDegree() returns the current value of the inDegree of the peer
func (p *Peer) GetInDegree() int64 {
	p.RLock()
	defer p.RUnlock()
	return p.inDegree
}

// IncInDegree() increse the current value if the inDegree by 1
func (p *Peer) IncInDegree() {
	p.Lock()
	defer p.Unlock()
	p.inDegree++
}

// GetWeight() returns the current weight of the peer for PoS calculation
func (p *Peer) GetWeight() uint64 {
	p.RLock()
	defer p.RUnlock()
	return p.weight
}

// SetWeight() set the weight of the peer for PoS calculation
func (p *Peer) SetWeight(w uint64) {
	p.Lock()
	defer p.Unlock()
	p.weight = w
}

// PeerStore provides an interface for persistent storage and
// retrieval of peers.
type PeerStore interface {
	// GetPeers returns the list of known peers.
	GetPeers() (*Peers, error)

	// GetPeers returns the list of known peers written as peer.PeerMessages.
	GetPeersFromMessages() (*Peers, error)

	// SetPeers sets the list of known peers. This is invoked when a peer is
	// added or removed.
	SetPeers([]*Peer) error
}

// ExcludePeer is used to exclude a single peer from a list of peers.
func ExcludePeer(peers []*Peer, peer string) (int, []*Peer) {
	index := -1
	otherPeers := make([]*Peer, 0, len(peers))
	for i, p := range peers {
		if p.Message.NetAddr != peer && p.Message.PubKeyHex != peer {
			otherPeers = append(otherPeers, p)
		} else {
			index = i
		}
	}
	return index, otherPeers
}

// ExcludePeers is used to exclude multiple peers from a list of peers.
func ExcludePeers(peers []*Peer, local string, last string) []*Peer {
	otherPeers := make([]*Peer, 0, len(peers))
	for _, p := range peers {
		if p.Message.NetAddr != local &&
			p.Message.PubKeyHex != local &&
			p.Message.NetAddr != last &&
			p.Message.PubKeyHex != last {
			otherPeers = append(otherPeers, p)
		}
	}
	return otherPeers
}
