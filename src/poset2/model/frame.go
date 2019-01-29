package model

import (
	"bytes"

	"github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/crypto"
)

type Frame struct {
	data *FrameData
}

type FrameData struct {
	Round  int64
	Events []*eventExternalData
}

func NewFrame(round int64, events []*Event) *Frame {
	var eventsData []*eventExternalData
	for k := range events {
		eventsData = append(eventsData, events[k].data.External)
	}

	return &Frame{
		data: &FrameData{
			Round:  round,
			Events: eventsData,
		},
	}
}

// NOT effective
func (f *Frame) Hash(codec Codec) common.Hash {
	data, err := f.encode(codec)
	if err != nil {
		return common.Hash{}
	}

	return crypto.Keccak256Hash(data)
}

func (f *Frame) Encode(codec Codec) ([]byte, error) {
	return f.encode(codec)
}

func (f *Frame) encode(codec Codec) ([]byte, error) {
	var data []byte
	buf := bytes.NewBuffer(data)
	err := codec.Encode(buf, f.data)

	return buf.Bytes(), err
}

func DecodeFrame(codec Codec, data []byte) (*Frame, error) {
	buf := bytes.NewBuffer(data)
	var result *FrameData
	if err := codec.Decode(buf, &result); err != nil {
		return nil, err
	}

	return &Frame{data: result}, nil
}
