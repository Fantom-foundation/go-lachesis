package vector

import (
	"github.com/Fantom-foundation/go-ethereum/rlp"

	"github.com/Fantom-foundation/go-lachesis/kvdb"
)

func (vi *Index) setRlp(table kvdb.KeyValueStore, key []byte, val interface{}) {
	buf, err := rlp.EncodeToBytes(val)
	if err != nil {
		vi.Log.Crit("Failed to encode rlp", "err", err)
	}

	if err := table.Put(key, buf); err != nil {
		vi.Log.Crit("Failed to put key-value", "err", err)
	}
}

func (vi *Index) getRlp(table kvdb.KeyValueStore, key []byte, to interface{}) interface{} {
	buf, err := table.Get(key)
	if err != nil {
		vi.Log.Crit("Failed to get key-value", "err", err)
	}
	if buf == nil {
		return nil
	}

	err = rlp.DecodeBytes(buf, to)
	if err != nil {
		vi.Log.Crit("Failed to decode rlp", "err", err, "size", len(buf))
	}
	return to
}