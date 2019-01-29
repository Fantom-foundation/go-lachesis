package model

import (
	"bytes"
	"github.com/Fantom-foundation/go-lachesis/src/common"
)

type Round struct {
	data *roundData
}

type roundData struct {
	Events map[common.Hash]*roundEventInfo
	Queued bool
}

type roundEventInfo struct {
	Root bool
}

func NewRound() *Round {
	return &Round{
		data: &roundData{
			Events: make(map[common.Hash]*roundEventInfo),
		},
	}
}

func (r *Round) AddEvent(event common.Hash, root bool) {
	if _, ok := r.data.Events[event]; !ok {
		r.data.Events[event] = &roundEventInfo{Root: root}
		return
	}
	r.data.Events[event].Root = root
}

func (r *Round) Roots() (result []common.Hash) {
	for hash := range r.data.Events {
		if r.data.Events[hash].Root {
			result = append(result, hash)
		}
	}

	return result
}

//func (r *Round) Events(consensus bool) (result []common.Hash) {
//	for hash := range r.data.Events {
//		if r.data.Events[hash].Consensus != consensus {
//			continue
//		}
//		result = append(result, hash)
//	}
//
//	return result
//}

func (r *Round) IsQueued() bool {
	return r.data.Queued
}

func (r *Round) SetQueued(queued bool) {
	r.data.Queued = queued
}

//func (r *Round) IsDecided() bool {
//	return r.data.Decided
//}

func (r *Round) Encode(codec Codec) ([]byte, error) {
	var data []byte
	buf := bytes.NewBuffer(data)
	err := codec.Encode(buf, r.data)

	return buf.Bytes(), err
}

func DecodeRound(codec Codec, data []byte) (*Round, error) {
	var roundData *roundData
	buf := bytes.NewBuffer(data)
	if err := codec.Decode(buf, &roundData); err != nil {
		return nil, err
	}

	return &Round{data: roundData}, nil
}
