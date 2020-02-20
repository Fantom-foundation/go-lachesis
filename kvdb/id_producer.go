package kvdb

import (
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
	err := p.table.Put(p.key, []byte(id))
	if err != nil {
		log.Crit("Failed to put key-value", "err", err)
	}
}
