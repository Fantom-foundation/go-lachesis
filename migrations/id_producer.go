package migrations

import (
	"github.com/pkg/errors"

	"github.com/Fantom-foundation/go-lachesis/kvdb"
)

// Flushable id implementation
type FlushableIdProducer struct {
	dbKey    kvdb.KeyValueStore		`table:"I"`
	keyBytes []byte
}

func NewFlushableIdProducer(db kvdb.KeyValueStore, key string) *FlushableIdProducer {
	return &FlushableIdProducer{
		dbKey:    db,
		keyBytes: []byte(key),
	}
}

func (p *FlushableIdProducer) GetId() (string, error) {
	id, err := p.dbKey.Get(p.keyBytes)
	if err != nil {
		return "", errors.Wrap(err, "FlushableIdProduser: GetId")
	}
	return string(id), nil
}

func (p *FlushableIdProducer) SetId(id string) error {
	err := p.dbKey.Put(p.keyBytes, []byte(id))
	if err != nil {
		return errors.Wrap(err, "FlushableIdProduser: SetId")
	}
	return nil
}
