package poset

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"

	"github.com/andrecronje/lachesis/src/crypto"
)

/*******************************************************************************
EventBody
*******************************************************************************/

type EventBody struct {
	Transactions    [][]byte         //the payload
	Parents         []string         //hashes of the event's parents, self-parent first
	Creator         []byte           //creator's public key
	Index           int              //index in the sequence of events created by Creator
	BlockSignatures []BlockSignature //list of Block signatures signed by the Event's Creator ONLY

	//wire
	//It is cheaper to send ints then hashes over the wire
	selfParentIndex      int
	otherParentCreatorID int
	otherParentIndex     int
	creatorID            int
}

//json encoding of body only
func (e *EventBody) Marshal() ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b) //will write to b
	if err := enc.Encode(e); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (e *EventBody) Unmarshal(data []byte) error {
	b := bytes.NewBuffer(data)
	dec := json.NewDecoder(b) //will read from b
	if err := dec.Decode(e); err != nil {
		return err
	}
	return nil
}

func (e *EventBody) Hash() ([]byte, error) {
	hashBytes, err := e.Marshal()
	if err != nil {
		return nil, err
	}
	return crypto.SHA256(hashBytes), nil
}

/*******************************************************************************
Event
*******************************************************************************/

type EventCoordinates struct {
	hash  string
	index int
}

type OrderedEventCoordinates []Index

func (o OrderedEventCoordinates) GetIDIndex(id int) int {
	for i, idx := range o {
		if idx.participantId == id {
			return i
		}
	}

	return -1
}

func (o OrderedEventCoordinates) GetByID(id int) (Index, bool) {
	for _, idx := range o {
		if idx.participantId == id {
			return idx, true
		}
	}

	return Index{}, false
}

func (o *OrderedEventCoordinates) Add(id int, event EventCoordinates) {
	*o = append(*o, Index{
		participantId: id,
		event:         event,
	})
}

type Index struct {
	participantId int
	event         EventCoordinates
}

// -----

type Event struct {
	Body      EventBody
	Signature string //creator's digital signature of body

	topologicalIndex int

	//used for sorting
	round            *int
	lamportTimestamp *int

	roundReceived *int

	lastAncestors    OrderedEventCoordinates //[participant fake id] => last ancestor
	firstDescendants OrderedEventCoordinates //[participant fake id] => first descendant

	creator string
	hash    []byte
	hex     string

	flagTable map[string]bool // flagTable stores connection information
	flags     int             // flags stores the number of connections
}

func NewEvent(transactions [][]byte,
	blockSignatures []BlockSignature,
	parents []string,
	creator []byte,
	index int,
	flagTable map[string]bool,
	flags int) Event {

	body := EventBody{
		Transactions:    transactions,
		BlockSignatures: blockSignatures,
		Parents:         parents,
		Creator:         creator,
		Index:           index,
	}
	return Event{
		Body:      body,
		flagTable: flagTable,
		flags:     flags,
	}
}

func (e *Event) Creator() string {
	if e.creator == "" {
		e.creator = fmt.Sprintf("0x%X", e.Body.Creator)
	}
	return e.creator
}

func (e *Event) SelfParent() string {
	return e.Body.Parents[0]
}

func (e *Event) OtherParent() string {
	return e.Body.Parents[1]
}

func (e *Event) Transactions() [][]byte {
	return e.Body.Transactions
}

func (e *Event) Index() int {
	return e.Body.Index
}

func (e *Event) BlockSignatures() []BlockSignature {
	return e.Body.BlockSignatures
}

//True if Event contains a payload or is the initial Event of its creator
func (e *Event) IsLoaded() bool {
	if e.Body.Index == 0 {
		return true
	}

	hasTransactions := e.Body.Transactions != nil &&
		len(e.Body.Transactions) > 0

	return hasTransactions
}

//ecdsa sig
func (e *Event) Sign(privKey *ecdsa.PrivateKey) error {
	signBytes, err := e.Body.Hash()
	if err != nil {
		return err
	}
	R, S, err := crypto.Sign(privKey, signBytes)
	if err != nil {
		return err
	}
	e.Signature = crypto.EncodeSignature(R, S)
	return err
}

func (e *Event) Verify() (bool, error) {
	pubBytes := e.Body.Creator
	pubKey := crypto.ToECDSAPub(pubBytes)

	signBytes, err := e.Body.Hash()
	if err != nil {
		return false, err
	}

	r, s, err := crypto.DecodeSignature(e.Signature)
	if err != nil {
		return false, err
	}

	return crypto.Verify(pubKey, signBytes, r, s), nil
}

//json encoding of body and signature
func (e *Event) Marshal() ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	if err := enc.Encode(e); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (e *Event) Unmarshal(data []byte) error {
	b := bytes.NewBuffer(data)
	dec := json.NewDecoder(b) //will read from b
	return dec.Decode(e)
}

//sha256 hash of body
func (e *Event) Hash() ([]byte, error) {
	if len(e.hash) == 0 {
		hash, err := e.Body.Hash()
		if err != nil {
			return nil, err
		}
		e.hash = hash
	}
	return e.hash, nil
}

func (e *Event) Hex() string {
	if e.hex == "" {
		hash, _ := e.Hash()
		e.hex = fmt.Sprintf("0x%X", hash)
	}
	return e.hex
}

func (e *Event) SetRound(r int) {
	if e.round == nil {
		e.round = new(int)
	}
	*e.round = r
}

func (e *Event) SetLamportTimestamp(t int) {
	if e.lamportTimestamp == nil {
		e.lamportTimestamp = new(int)
	}
	*e.lamportTimestamp = t
}

func (e *Event) SetRoundReceived(rr int) {
	if e.roundReceived == nil {
		e.roundReceived = new(int)
	}
	*e.roundReceived = rr
}

func (e *Event) SetWireInfo(selfParentIndex,
	otherParentCreatorID,
	otherParentIndex,
	creatorID int) {
	e.Body.selfParentIndex = selfParentIndex
	e.Body.otherParentCreatorID = otherParentCreatorID
	e.Body.otherParentIndex = otherParentIndex
	e.Body.creatorID = creatorID
}

func (e *Event) WireBlockSignatures() []WireBlockSignature {
	if e.Body.BlockSignatures != nil {
		wireSignatures := make([]WireBlockSignature, len(e.Body.BlockSignatures))
		for i, bs := range e.Body.BlockSignatures {
			wireSignatures[i] = bs.ToWire()
		}

		return wireSignatures
	}
	return nil
}

func (e *Event) ToWire() WireEvent {

	return WireEvent{
		Body: WireBody{
			Transactions:         e.Body.Transactions,
			SelfParentIndex:      e.Body.selfParentIndex,
			OtherParentCreatorID: e.Body.otherParentCreatorID,
			OtherParentIndex:     e.Body.otherParentIndex,
			CreatorID:            e.Body.creatorID,
			Index:                e.Body.Index,
			BlockSignatures:      e.WireBlockSignatures(),
		},
		Signature: e.Signature,
	}
}

// FlagTable returns the flagTable and the number of flags of the event
func (e *Event) FlagTable() (map[string]bool, int) {
	return e.flagTable, e.flags
}

func rootSelfParent(participantID int) string {
	return fmt.Sprintf("Root%d", participantID)
}

/*******************************************************************************
Sorting
*******************************************************************************/

// ByTopologicalOrder implements sort.Interface for []Event based on
// the topologicalIndex field.
// THIS IS A PARTIAL ORDER
type ByTopologicalOrder []Event

func (a ByTopologicalOrder) Len() int      { return len(a) }
func (a ByTopologicalOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByTopologicalOrder) Less(i, j int) bool {
	return a[i].topologicalIndex < a[j].topologicalIndex
}

// ByLamportTimestamp implements sort.Interface for []Event based on
// the lamportTimestamp field.
// THIS IS A TOTAL ORDER
type ByLamportTimestamp []Event

func (a ByLamportTimestamp) Len() int      { return len(a) }
func (a ByLamportTimestamp) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByLamportTimestamp) Less(i, j int) bool {
	it, jt := -1, -1
	if a[i].lamportTimestamp != nil {
		it = *a[i].lamportTimestamp
	}
	if a[j].lamportTimestamp != nil {
		jt = *a[j].lamportTimestamp
	}
	if it != jt {
		return it < jt
	}

	wsi, _, _ := crypto.DecodeSignature(a[i].Signature)
	wsj, _, _ := crypto.DecodeSignature(a[j].Signature)
	return wsi.Cmp(wsj) < 0
}

/*******************************************************************************
 WireEvent
*******************************************************************************/

type WireBody struct {
	Transactions    [][]byte
	BlockSignatures []WireBlockSignature

	SelfParentIndex      int
	OtherParentCreatorID int
	OtherParentIndex     int
	CreatorID            int

	Index int
}

type WireEvent struct {
	Body      WireBody
	Signature string
}

func (we *WireEvent) BlockSignatures(validator []byte) []BlockSignature {
	if we.Body.BlockSignatures != nil {
		blockSignatures := make([]BlockSignature, len(we.Body.BlockSignatures))
		for k, bs := range we.Body.BlockSignatures {
			blockSignatures[k] = BlockSignature{
				Validator: validator,
				Index:     bs.Index,
				Signature: bs.Signature,
			}
		}
		return blockSignatures
	}
	return nil
}
