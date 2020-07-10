package app

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
)

// GetTotalSupply returns total supply
func (s *Store) GetTotalLocked() *big.Int {
	amountBytes, err := s.table.TotalLocked.Get([]byte("c"))
	if err != nil {
		s.Log.Crit("Failed to get key-value", "err", err)
	}
	if amountBytes == nil {
		return big.NewInt(0)
	}
	return new(big.Int).SetBytes(amountBytes)
}

// SetTotalSupply stores total supply
func (s *Store) SetTotalLocked(amount *big.Int) {
	err := s.table.TotalLocked.Put([]byte("c"), amount.Bytes())
	if err != nil {
		s.Log.Crit("Failed to set key-value", "err", err)
	}
}

type locked struct {
	FromEpoch idx.Epoch
	EndTime   inter.Timestamp
}

func (s *Store) SetStakeLock(stakerID idx.StakerID, fromEpoch idx.Epoch, endTime inter.Timestamp) {
	key := stakerID.Bytes()
	s.set(s.table.Locked, key, &locked{
		fromEpoch,
		endTime,
	})
}

func (s *Store) SetDelegateLock(addr common.Address, stakerID idx.StakerID, fromEpoch idx.Epoch, endTime inter.Timestamp) {
	key := append(stakerID.Bytes(), addr.Bytes()...)
	s.set(s.table.Locked, key, &locked{
		fromEpoch,
		endTime,
	})
}

func (s *Store) IsStakeLocked(stakerID idx.StakerID, epoch idx.Epoch) bool {
	key := stakerID.Bytes()
	locked, ok := s.get(s.table.Locked, key, &locked{}).(*locked)
	if !ok {
		return false
	}
	return locked != nil && locked.FromEpoch <= epoch
}

func (s *Store) IsDelegationLocked(addr common.Address, stakerID idx.StakerID, epoch idx.Epoch) bool {
	key := append(stakerID.Bytes(), addr.Bytes()...)
	locked, ok := s.get(s.table.Locked, key, &locked{}).(*locked)
	if !ok {
		return false
	}
	return locked != nil && locked.FromEpoch <= epoch
}

func (s *Store) DelOutdatedLocks(time inter.Timestamp) (amount *big.Int) {
	amount = big.NewInt(0)

	it := s.table.Locked.NewIterator()
	defer it.Release()

	keys := make([][]byte, 0, 500) // don't write during iteration

	for it.Next() {
		lock := &locked{}
		err := rlp.DecodeBytes(it.Value(), lock)
		if err != nil {
			s.Log.Crit("Failed to decode rlp while iteration", "err", err)
		}

		if lock.EndTime > time {
			continue
		}

		key := it.Key()
		keys = append(keys, key)

		stakerID := idx.BytesToStakerID(key[:4])

		if len(key) > 4 {
			address := common.BytesToAddress(key[4:])
			delegator := s.GetSfcDelegator(address, stakerID)
			amount.Add(amount, delegator.Amount)
		} else {
			staker := s.GetSfcStaker(stakerID)
			amount.Add(amount, staker.StakeAmount)
		}
	}

	for i := range keys {
		err := s.table.Locked.Delete(keys[i])
		if err != nil {
			s.Log.Crit("Failed to erase key-value", "err", err)
		}
	}

	return
}
