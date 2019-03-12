package poset

import (
	"fmt"
	"strconv"
	"github.com/golang/protobuf/proto"
)

// FlagTable is a dedicated type for the Events flags map.
type FlagTable map[uint64]int64

// NewFlagTable creates new empty FlagTable
func NewFlagTable() FlagTable {
	return FlagTable(make(map[uint64]int64))
}

// Marshal converts FlagTable to protobuf.
func (ft FlagTable) Marshal() []byte {
	body := make(map[string]int64, len(ft))
	for k, v := range ft {
		body[fmt.Sprintf("%v", k)] = v
	}

	wrapper := &FlagTableWrapper{Body: body}
	bytes, err := proto.Marshal(wrapper)
	if err != nil {
		panic(err)
	}

	return bytes
}

// Unmarshal reads protobuff into FlagTable.
func (ft FlagTable) Unmarshal(buf []byte) error {
	wrapper := new(FlagTableWrapper)
	err := proto.Unmarshal(buf, wrapper)
	if err != nil {
		return err
	}

	for k, v := range wrapper.Body {
		var creatorID uint64
		creatorID, err = strconv.ParseUint(k, 10, 64)
		if err != nil {
			return err
		}
		ft[creatorID] = v
	}
	return nil
}

func (ft FlagTable) Copy() FlagTable {
	res := NewFlagTable()
	for id, frame := range ft {
		res[id] = frame
	}
	return res
}
