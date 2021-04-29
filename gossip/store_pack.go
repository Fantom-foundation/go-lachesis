package gossip

/*
	In LRU cache data stored like value
*/

import (
	"bytes"

	"github.com/ethereum/go-ethereum/rlp"

	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
)

const (
	epochSize   = 4
	packSize    = 4
	eventIDSize = 32
)

func (s *Store) GetPackInfo(epoch idx.Epoch, idx idx.Pack) *PackInfo {
	key := bytes.Buffer{}
	key.Write(epoch.Bytes())
	key.Write(idx.Bytes())

	// Get data from LRU cache first.
	if s.cache.PackInfos != nil {
		if c, ok := s.cache.PackInfos.Get(key.String()); ok {
			if b, ok := c.(PackInfo); ok {
				return &b
			}
		}
	}

	w, _ := s.get(s.table.PackInfos, key.Bytes(), &PackInfo{}).(*PackInfo)

	// Add to LRU cache.
	if w != nil && s.cache.PackInfos != nil {
		s.cache.PackInfos.Add(key.String(), *w)
	}

	return w
}

// returns default value if not found
func (s *Store) GetPackInfoOrDefault(epoch idx.Epoch, idx idx.Pack) PackInfo {
	packInfo := s.GetPackInfo(epoch, idx)
	if packInfo == nil {
		return PackInfo{
			Index: idx,
		}
	}
	return *packInfo
}

func (s *Store) GetPackInfoRLP(epoch idx.Epoch, idx idx.Pack) rlp.RawValue {
	key := bytes.Buffer{}
	key.Write(epoch.Bytes())
	key.Write(idx.Bytes())

	w, _ := s.table.PackInfos.Get(key.Bytes())
	return w
}

func (s *Store) SetPackInfo(epoch idx.Epoch, idx idx.Pack, value PackInfo) {
	key := bytes.Buffer{}
	key.Write(epoch.Bytes())
	key.Write(idx.Bytes())

	s.set(s.table.PackInfos, key.Bytes(), value)

	// Add to LRU cache.
	if s.cache.PackInfos != nil {
		s.cache.PackInfos.Add(key.String(), value)
	}
}

func (s *Store) AddToPack(epoch idx.Epoch, idx idx.Pack, e hash.Event) {
	key := bytes.Buffer{}
	key.Write(epoch.Bytes())
	key.Write(idx.Bytes())
	key.Write(e.Bytes())

	err := s.table.Packs.Put(key.Bytes(), []byte{})
	if err != nil {
		s.Log.Crit("Failed to put key-value", "err", err)
	}
}

func (s *Store) GetPack(epoch idx.Epoch, idx idx.Pack) hash.Events {
	prefix := bytes.Buffer{}
	prefix.Write(epoch.Bytes())
	prefix.Write(idx.Bytes())

	res := make(hash.Events, 0, hardLimitItems)

	it := s.table.Packs.NewIterator(prefix.Bytes(), nil)
	defer it.Release()
	for it.Next() {
		if len(it.Key()) != epochSize+packSize+eventIDSize {
			s.Log.Crit("packs table: Incorrect key len", "len(key)", len(it.Key()))
		}
		res.Add(hash.BytesToEvent(it.Key()[epochSize+packSize:]))
	}
	if it.Error() != nil {
		s.Log.Crit("Failed to iterate keys", "err", it.Error())
	}

	if len(res) == 0 {
		return nil
	}
	return res
}

func (s *Store) GetPacksNum(epoch idx.Epoch) (idx.Pack, bool) {
	b, err := s.table.PacksNum.Get(epoch.Bytes())
	if err != nil {
		s.Log.Crit("Failed to get key-value", "err", err)
	}
	if b == nil {
		return 0, false
	}
	return idx.BytesToPack(b), true
}

func (s *Store) GetPacksNumOrDefault(epoch idx.Epoch) idx.Pack {
	num, ok := s.GetPacksNum(epoch)
	if !ok {
		return 1
	}
	return num
}

func (s *Store) SetPacksNum(epoch idx.Epoch, num idx.Pack) {
	err := s.table.PacksNum.Put(epoch.Bytes(), num.Bytes())
	if err != nil {
		s.Log.Crit("Failed to put key-value", "err", err)
	}
}
