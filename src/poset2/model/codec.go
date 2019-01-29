package model

import (
	"encoding/gob"
	"github.com/Fantom-foundation/go-lachesis/src/common"
	"io"
)

type Codec interface {
	Decode(r io.Reader, val interface{}) error
	Encode(w io.Writer, val interface{}) error
}

type DefaultCodec struct{}

func NewDefaultCodec() *DefaultCodec {
	gob.Register(common.Hash{})
	gob.Register([]common.Hash{})
	gob.Register([][]byte{})
	return &DefaultCodec{}
}

func (c *DefaultCodec) Decode(r io.Reader, val interface{}) error {
	dec := gob.NewDecoder(r)
	return dec.Decode(val)
}

func (c *DefaultCodec) Encode(w io.Writer, val interface{}) error {
	enc := gob.NewEncoder(w)
	return enc.Encode(val)
}
