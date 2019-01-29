package model

import (
	"bytes"
	"crypto/ecdsa"

	"github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/crypto"
)

type Event struct {
	data *eventData
}

type eventData struct {
	External *eventExternalData
	// Internal data, not intended to be transmitted over a network.
	Internal *eventInternalData
}

type eventInternalData struct {
	Round   uint64
	Lamport uint64
}

type eventExternalData struct {
	RootProof    []common.Hash
	Creator      []byte
	Flags        []common.Hash
	Index        uint64
	Parents      []common.Hash
	Transactions [][]byte
	Signature    string
}

func NewEvent(creator []byte,
	flags []common.Hash, index uint64, parents []common.Hash,
	transactions [][]byte) *Event {
	return &Event{
		data: &eventData{
			External: &eventExternalData{
				Creator:      creator,
				Flags:        flags,
				Index:        index,
				Parents:      parents,
				Transactions: transactions,
			},
			Internal: &eventInternalData{},
		},
	}
}

func NewLeafEvent(codec Codec, key *ecdsa.PrivateKey) (*Event, error) {
	var (
		pubKey = crypto.FromECDSAPub(&key.PublicKey)
	)

	event := NewEvent(pubKey, nil, 0, nil, nil)
	if err := event.Sign(key, codec); err != nil {
		return nil, err
	}
	event.data.Internal.Lamport = 1
	return event, nil
}

func IsLeafEvent(codec Codec, event *Event) bool {
	valid, err := event.VerifySignature(codec)
	if err != nil {
		return false
	}

	// TODO more checks

	return valid && event.Index() == 0
}

func (e *Event) RootProof() []common.Hash {
	return e.data.External.RootProof
}

func (e *Event) SetRootProof(roots []common.Hash) {
	e.data.External.RootProof = roots
}

func (e *Event) Creator() []byte {
	return e.data.External.Creator
}

func (e *Event) SetFlags(flags []common.Hash) {
	e.data.External.Flags = flags
}

func (e *Event) Flags() []common.Hash {
	return e.data.External.Flags
}

func (e *Event) Index() uint64 {
	return e.data.External.Index
}

func (e *Event) Parents() []common.Hash {
	return e.data.External.Parents
}

func (e *Event) Transactions() [][]byte {
	return e.data.External.Transactions
}

func (e *Event) SelfParent() common.Hash {
	if len(e.data.External.Parents) == 0 {
		return common.Hash{}
	}
	return e.data.External.Parents[0]
}

func (e *Event) OtherParent() common.Hash {
	if len(e.data.External.Parents) < 2 {
		return common.Hash{}
	}
	return e.data.External.Parents[1]
}

func (e *Event) SetRound(round uint64) {
	e.data.Internal.Round = round
}

func (e *Event) Round() uint64 {
	return e.data.Internal.Round
}

func (e *Event) Lamport() uint64 {
	return e.data.Internal.Lamport
}

func (e *Event) SetLamport(lamport uint64) {
	e.data.Internal.Lamport = lamport
}

func (e *Event) Hash(codec Codec) common.Hash {
	var data []byte
	buf := bytes.NewBuffer(data)
	err := codec.Encode(buf, []interface{}{
		e.data.External.RootProof,
		e.data.External.Creator,
		e.data.External.Flags,
		e.data.External.Index,
		e.data.External.Parents,
		e.data.External.Transactions})
	if err != nil {
		return common.Hash{}
	}

	return crypto.Keccak256Hash(buf.Bytes())
}

func (e *Event) VerifySignature(codec Codec) (bool, error) {
	pubKey := crypto.ToECDSAPub(e.data.External.Creator)
	r, s, err := crypto.DecodeSignature(string(e.data.External.Signature))
	if err != nil {
		return false, err
	}

	return crypto.Verify(pubKey, e.Hash(codec).Bytes(), r, s), nil
}

func (e *Event) Sign(privateKey *ecdsa.PrivateKey, codec Codec) error {
	r, s, err := crypto.Sign(privateKey, e.Hash(codec).Bytes())
	if err != nil {
		return err
	}
	e.data.External.Signature = crypto.EncodeSignature(r, s)

	return nil
}

func (e *Event) Equals(other *Event, codec Codec) bool {
	hash1 := e.Hash(codec)
	hash2 := other.Hash(codec)

	return bytes.Equal(hash1.Bytes(), hash2.Bytes()) &&
		e.data.External.Signature == other.data.External.Signature
}

func (e *Event) IsLoaded() bool {
	return len(e.data.External.Transactions) != 0
}

func (e *Event) Encode(codec Codec) ([]byte, error) {
	var data []byte
	buf := bytes.NewBuffer(data)
	err := codec.Encode(buf, e.data)

	return buf.Bytes(), err
}

func DecodeEvent(codec Codec, data []byte) (*Event, error) {
	buf := bytes.NewBuffer(data)
	var eventData *eventData
	err := codec.Decode(buf, &eventData)

	return &Event{data: eventData}, err
}

func (e *Event) ToWire(codec Codec) ([]byte, error) {
	var data []byte
	buf := bytes.NewBuffer(data)
	err := codec.Encode(buf, e.data.External)

	return buf.Bytes(), err
}

func EventFromWire(codec Codec, data []byte) (*Event, error) {
	buf := bytes.NewBuffer(data)
	var extEventData *eventExternalData
	err := codec.Decode(buf, &extEventData)
	if err != nil {
		return nil, err
	}

	return &Event{
		data: &eventData{
			External: extEventData,
			Internal: &eventInternalData{},
		},
	}, nil
}

func (e *eventExternalData) Hash(codec Codec) common.Hash {
	var data []byte
	buf := bytes.NewBuffer(data)
	err := codec.Encode(buf, []interface{}{e.RootProof, e.Creator,
		e.Flags, e.Index, e.Parents, e.Transactions})
	if err != nil {
		return common.Hash{}
	}

	return crypto.Keccak256Hash(buf.Bytes())
}
