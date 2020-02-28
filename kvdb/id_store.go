package kvdb

import (
	"github.com/ethereum/go-ethereum/log"
)

// IDStore stores id
type IDStore struct {
	table KeyValueStore
	key   []byte
}

// NewIDStore constructor
func NewIDStore(table KeyValueStore) *IDStore {
	return &IDStore{
		table: table,
		key:   []byte("id"),
	}
}

// GetID is a getter
func (p *IDStore) GetID() string {
	id, err := p.table.Get(p.key)
	if err != nil {
		log.Crit("Failed to get key-value", "err", err)
	}

	if id == nil {
		return ""
	}
	return string(id)
}

// SetID is a setter
func (p *IDStore) SetID(id string) {
	err := p.table.Put(p.key, []byte(id))
	if err != nil {
		log.Crit("Failed to put key-value", "err", err)
	}
}
