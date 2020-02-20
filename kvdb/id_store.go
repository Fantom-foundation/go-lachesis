package kvdb

import (
	"github.com/ethereum/go-ethereum/log"
)

// IdStore stores id
type IdStore struct {
	table KeyValueStore
	key   []byte
}

// NewIdStore constructor
func NewIdStore(table KeyValueStore) *IdStore {
	return &IdStore{
		table: table,
		key:   []byte("id"),
	}
}

// GetId is a getter
func (p *IdStore) GetId() string {
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
func (p *IdStore) SetId(id string) {
	err := p.table.Put(p.key, []byte(id))
	if err != nil {
		log.Crit("Failed to put key-value", "err", err)
	}
}
