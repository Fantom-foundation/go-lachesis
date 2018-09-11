package poset

import (
	"bytes"
	"fmt"

	"github.com/ugorji/go/codec"
)

/*
Roots constitute the base of a Poset. Each Participant is assigned a Root on
top of which Events will be added. The first Event of a participant must have a
Self-Parent and n Other-Parents that match its Root respectively.

This construction allows us to initialize Posets where the first Events are
taken from the middle of another Poset*/

//RootEvent contains enough information about an Event and its direct descendant
//to allow inserting Events on top of it.
type RootEvent struct {
	Hash             string
	CreatorID        int
	Index            int
	LamportTimestamp int
	Round            int
}

//NewBaseRootEvent creates a RootEvent corresponding to the the very beginning
//of a Poset.
func NewBaseRootEvent(creatorID int) RootEvent {
	res := RootEvent{
		Hash:             fmt.Sprintf("Root%d", creatorID),
		CreatorID:        creatorID,
		Index:            -1,
		LamportTimestamp: -1,
		Round:            -1,
	}
	return res
}

//Root forms a base on top of which a participant's Events can be inserted. In
//contains the SelfParent of the first descendant of the Root, as well as other
//Events, belonging to a past before the Root, which might be referenced
//in future Events. NextRound corresponds to a proposed value for the child's
//Round; it is only used if the child's OtherParent is empty or NOT in the
//Root's Others.
type Root struct {
	NextRound  int
	SelfParent RootEvent
	Others     map[string][]RootEvent
}

//NewBaseRoot initializes a Root object for a fresh Poset.
func NewBaseRoot(creatorID int) Root {
	res := Root{
		NextRound:  0,
		SelfParent: NewBaseRootEvent(creatorID),
		Others:     map[string][]RootEvent{},
	}
	return res
}

//The JSON encoding of a Root must be DETERMINISTIC because it is itself
//included in the JSON encoding of a Frame. The difficulty is that Roots contain
//go maps for which one should not expect a de facto order of entries; we cannot
//use the builtin JSON codec within overriding something. Instead, we are using
//a third party library (ugorji/codec) that enables deterministic encoding of
//golang maps.
func (root *Root) Marshal() ([]byte, error) {

	b := new(bytes.Buffer)
	jh := new(codec.JsonHandle)
	jh.Canonical = true
	enc := codec.NewEncoder(b, jh)

	if err := enc.Encode(root); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (root *Root) Unmarshal(data []byte) error {

	b := bytes.NewBuffer(data)
	jh := new(codec.JsonHandle)
	jh.Canonical = true
	dec := codec.NewDecoder(b, jh)

	return dec.Decode(root)
}
