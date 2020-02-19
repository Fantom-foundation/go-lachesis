package kvdb

import (
	"crypto/sha256"
	"fmt"
	"github.com/Fantom-foundation/go-lachesis/utils/migration"
	"github.com/ethereum/go-ethereum/log"
)

// IdProducer stores id
type IdProducer struct {
	table KeyValueStore
	migrationChain *migration.Migration
	key   []byte
}

// NewIdProducer constructor
func NewIdProducer(table KeyValueStore, chain *migration.Migration) *IdProducer {
	return &IdProducer{
		table: table,
		migrationChain: chain,
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
	// Search prev migration name
	prevName := ""
	prev := p.migrationChain.PrevByName(id)
	if prev != nil {
		prevName = prev.Name()
	}

	hashed := p.hashedId(id, prevName)
	err := p.table.Put(p.key, []byte(hashed))
	if err != nil {
		log.Crit("Failed to put key-value", "err", err)
	}
}

// IsCurrent return true if saved id equal from param
func (p *IdProducer) IsCurrent(id string) bool {
	currentId := p.GetId()

	prevName := ""
	prev := p.migrationChain.PrevByName(id)
	if prev != nil {
		prevName = prev.Name()
	}

	return p.hashedId(id, prevName) == currentId
}

func (p *IdProducer) hashedId(id, prev string) string {
	digest := sha256.New()

	digest.Write([]byte(prev))
	digest.Write([]byte(id))

	bytes := digest.Sum(nil)
	return fmt.Sprintf("%x", bytes)
}
