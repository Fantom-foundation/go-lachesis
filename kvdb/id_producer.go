package kvdb

import (
	"crypto/sha256"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
)

// IdProducer stores id
type IdProducer struct {
	table KeyValueStore
	key   []byte
}

// NewIdProducer constructor
func NewIdProducer(table KeyValueStore) *IdProducer {
	return &IdProducer{
		table: table,
		key:   []byte("id"),
	}
}

// GetId is a getter
func (p *IdProducer) GetId() string {
	id, err := p.table.Get(p.key)
	if err != nil {
		log.Crit("Failed to get key-value", "err", err)
	}

	if id == nil {
		return ""
	}
	return string(id)
}

// SetId is a setter
func (p *IdProducer) SetId(id string) {
	hashed := p.hashedId(id)
	err := p.table.Put(p.key, []byte(hashed))
	if err != nil {
		log.Crit("Failed to put key-value", "err", err)
	}
}

// IsCurrent return true if saved id equal from param
func (p *IdProducer) IsCurrent(id string) bool {
	currentId := p.GetId()
	return p.hashedId(id) == currentId
}

func (p *IdProducer) hashedId(id string) string {
	digest := sha256.New()

	digest.Write([]byte(id))

	bytes := digest.Sum(nil)
	return fmt.Sprintf("%x", bytes)
}
