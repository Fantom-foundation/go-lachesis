package vector

import (
	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/kvdb"
)

func (vi *Index) getBytes(table kvdb.KeyValueStore, id hash.Event) []byte {
	key := id.Bytes()
	b, err := table.Get(key)
	if err != nil {
		vi.Log.Crit("Failed to get key-value", "err", err)
	}
	return b
}

func (vi *Index) setBytes(table kvdb.KeyValueStore, id hash.Event, b []byte) {
	key := id.Bytes()
	err := table.Put(key, b)
	if err != nil {
		vi.Log.Crit("Failed to put key-value", "err", err)
	}
}
