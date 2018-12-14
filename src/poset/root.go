package poset

import (
	"fmt"

	"github.com/golang/protobuf/proto"
)

/*
Roots constitute the base of a Poset. Each Participant is assigned a Root on
top of which Events will be added. The first Event of a participant must have a
Self-Parent and an Other-Parent that match its Root X and Y respectively.

This construction allows us to initialize Posets where the first Events are
taken from the middle of another Poset

ex 1:

-----------------        -----------------       -----------------
- Event E0      -        - Event E1      -       - Event E2      -
- SP = ""       -        - SP = ""       -       - SP = ""       -
- OP = ""       -        - OP = ""       -       - OP = ""       -
-----------------        -----------------       -----------------
        |                        |                       |
-----------------		 -----------------		 -----------------
- Root 0        - 		 - Root 1        - 		 - Root 2        -
- X = Y = ""    - 		 - X = Y = ""    -		 - X = Y = ""    -
- Index= -1     -		 - Index= -1     -       - Index= -1     -
- Others= empty - 		 - Others= empty -       - Others= empty -
-----------------		 -----------------       -----------------

ex 2:

-----------------
- Event E02     -
- SP = E01      -
- OP = E_OLD    -
-----------------
       |
-----------------
- Event E01     -
- SP = E00      -
- OP = E10      -  \
-----------------    \
       |               \
-----------------        -----------------       -----------------
- Event E00     -        - Event E10     -       - Event E20     -
- SP = x0       -        - SP = x1       -       - SP = x2       -
- OP = y0       -        - OP = y1       -       - OP = y2       -
-----------------        -----------------       -----------------
        |                        |                       |
-----------------		 -----------------		 -----------------
- Root 0        - 		 - Root 1        - 		 - Root 2        -
- X: x0, Y: y0  - 		 - X: x1, Y: y1  - 		 - X: x2, Y: y2  -
- Index= i0     -		 - Index= i1     -       - Index= i2     -
- Others= {     - 		 - Others= empty -       - Others= empty -
-  E02: E_OLD   -        -----------------       -----------------
- }             -
-----------------
*/

//RootEvent contains enough information about an Event and its direct descendant
//to allow inserting Events on top of it.
//NewBaseRootEvent creates a RootEvent corresponding to the the very beginning
//of a Poset.
func NewBaseRootEvent(creatorID uint32) RootEvent {
	hash := fmt.Sprintf("Root%d", creatorID)
	res := RootEvent{
		Hash:             hash,
		CreatorID:        creatorID,
		Index:            -1,
		LamportTimestamp: -1,
		Round:            -1,
	}
	return res
}

func (this *RootEvent) Equals(that *RootEvent) bool {
	return this.Hash == that.Hash &&
		this.CreatorID == that.CreatorID &&
		this.Index == that.Index &&
		this.LamportTimestamp == that.LamportTimestamp &&
		this.Round == that.Round
}

//Root forms a base on top of which a participant's Events can be inserted. It
//contains the SelfParent of the first descendant of the Root, as well as other
//Events, belonging to a past before the Root, which might be referenced
//in future Events. NextRound corresponds to a proposed value for the child's
//Round; it is only used if the child's OtherParent is empty or NOT in the
//Root's Others.
//NewBaseRoot initializes a Root object for a fresh Poset.
func NewBaseRoot(creatorID uint32) *Root {
	rootEvent := NewBaseRootEvent(creatorID)
	res := &Root{
		NextRound:  0,
		SelfParent: &rootEvent,
		Others:     map[string]*RootEvent{},
	}
	return res
}


func EqualsMapStringRootEvent(this map[string]*RootEvent, that map[string]*RootEvent) bool {
	if len(this) != len(that) {
		return false
	}
	for k, v := range this {
		v2, ok := that[k]
		if !ok || !v2.Equals(v) {
			return false
		}
	}
	return true
}

func (this *Root) Equals(that *Root) bool {
	return this.NextRound == that.NextRound &&
		this.SelfParent.Equals(that.SelfParent) &&
		EqualsMapStringRootEvent(this.Others, that.Others)
}

func (root *Root) ProtoMarshal() ([]byte, error) {
	var bf proto.Buffer
	bf.SetDeterministic(true)
	if err := bf.Marshal(root); err != nil {
		return nil, err
	}
	return bf.Bytes(), nil
}

func (root *Root) ProtoUnmarshal(data []byte) error {
	return proto.Unmarshal(data, root)
}
