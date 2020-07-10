package app

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/sfctype"
)

type delegatorAndStaker struct {
	Delegator common.Address
	Staker    idx.StakerID
}

// SetSfcDelegator stores SfcDelegator
func (s *Store) SetSfcDelegator(address common.Address, staker idx.StakerID, v *sfctype.SfcDelegator) {
	key := append(address.Bytes(), staker.Bytes()...)
	v.ToStakerID = staker
	s.set(s.table.Delegators, key, v)

	// Add to LRU cache.
	if s.cache.Delegators != nil {
		s.cache.Delegators.Add(delegatorAndStaker{address, staker}, v)
	}
}

// DelSfcDelegator deletes SfcDelegator
func (s *Store) DelSfcDelegator(address common.Address, staker idx.StakerID) {
	key := append(address.Bytes(), staker.Bytes()...)
	err := s.table.Delegators.Delete(key)
	if err != nil {
		s.Log.Crit("Failed to erase delegator")
	}

	// Add to LRU cache.
	if s.cache.Delegators != nil {
		s.cache.Delegators.Remove(delegatorAndStaker{address, staker})
	}
}

// ForEachSfcDelegator iterates all stored SfcDelegators
func (s *Store) ForEachSfcDelegator(do func(sfctype.SfcDelegatorAndAddr)) {
	it := s.table.Delegators.NewIterator()
	defer it.Release()
	s.forEachSfcDelegator(it, do)
}

func (s *Store) forEachSfcDelegator(it ethdb.Iterator, do func(sfctype.SfcDelegatorAndAddr)) {
	for it.Next() {
		delegator := new(sfctype.SfcDelegator)
		err := rlp.DecodeBytes(it.Value(), delegator)
		if err != nil {
			s.Log.Crit("Failed to decode rlp while iteration", "err", err)
		}

		addr := common.BytesToAddress(it.Key()[:20])
		do(sfctype.SfcDelegatorAndAddr{
			Addr:      addr,
			Delegator: delegator,
		})
	}
}

// GetSfcDelegator returns stored SfcDelegator
func (s *Store) GetSfcDelegator(address common.Address, staker idx.StakerID) *sfctype.SfcDelegator {
	key := append(address.Bytes(), staker.Bytes()...)

	// Get data from LRU cache first.
	if s.cache.Delegators != nil {
		if c, ok := s.cache.Delegators.Get(delegatorAndStaker{address, staker}); ok {
			if b, ok := c.(*sfctype.SfcDelegator); ok {
				return b
			}
		}
	}

	w, _ := s.get(s.table.Delegators, key, &sfctype.SfcDelegator{}).(*sfctype.SfcDelegator)

	// Add to LRU cache.
	if w != nil && s.cache.Delegators != nil {
		s.cache.Delegators.Add(delegatorAndStaker{address, staker}, w)
	}

	return w
}

// GetSfcDelegations returns all stored SfcDelegators
func (s *Store) GetSfcDelegations(address common.Address) []*sfctype.SfcDelegator {
	var all []*sfctype.SfcDelegator

	prefix := address.Bytes()
	it := s.table.Delegators.NewIteratorWithPrefix(prefix)
	defer it.Release()

	s.forEachSfcDelegator(it, func(d sfctype.SfcDelegatorAndAddr) {
		all = append(all, d.Delegator)
	})

	return all
}
