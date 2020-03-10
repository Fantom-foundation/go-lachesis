package app

import (
	"bytes"
	"math/big"

	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/kvdb"
)

// IncBlocksMissed add count of missed blocks for validator
func (s *Store) IncBlocksMissed(stakerID idx.StakerID, periodDiff inter.Timestamp) {
	s.mutex.Inc.Lock()
	defer s.mutex.Inc.Unlock()

	missed := s.GetBlocksMissed(stakerID)
	missed.Num++
	missed.Period += periodDiff
	s.set(s.table.BlockDowntime, stakerID.Bytes(), &missed)

	s.cache.BlockDowntime.Add(stakerID, missed)
}

// NewDowntimeSnapshotEpoch add count of missed blocks for validator by epoch
func (s *Store) NewDowntimeSnapshotEpoch(epoch idx.Epoch) {
	s.newEpochSnapshot(s.table.BlockDowntime, s.table.BlockDowntimeEpoch, epoch)
}

// NewScoreSnapshotEpoch add scores for validator by epoch
func (s *Store) NewScoreSnapshotEpoch(epoch idx.Epoch) {
	s.newEpochSnapshot(s.table.ActiveValidationScore, s.table.ActiveValidationScoreEpoch, epoch)
}

// ResetBlocksMissed set to 0 missed blocks for validator
func (s *Store) ResetBlocksMissed(stakerID idx.StakerID) {
	s.mutex.Inc.Lock()
	defer s.mutex.Inc.Unlock()

	err := s.table.BlockDowntime.Delete(stakerID.Bytes())
	if err != nil {
		s.Log.Crit("Failed to set key-value", "err", err)
	}

	s.cache.BlockDowntime.Add(stakerID, BlocksMissed{})
}

// GetBlocksMissed return blocks missed num for validator
func (s *Store) GetBlocksMissed(stakerID idx.StakerID) BlocksMissed {
	missedVal, ok := s.cache.BlockDowntime.Get(stakerID)
	if ok {
		if missed, ok := missedVal.(BlocksMissed); ok {
			return missed
		}
	}

	pMissed, _ := s.get(s.table.BlockDowntime, stakerID.Bytes(), &BlocksMissed{}).(*BlocksMissed)
	if pMissed == nil {
		return BlocksMissed{}
	}
	missed := *pMissed

	s.cache.BlockDowntime.Add(stakerID, missed)

	return missed
}

func (s *Store) GetBlocksMissedEpoch(stakerID idx.StakerID, epoch idx.Epoch) BlocksMissed {
	key := bytes.Buffer{}
	if epoch != 0 {
		key.Write(epoch.Bytes())
	}
	key.Write(stakerID.Bytes())

	pMissed, _ := s.get(s.table.BlockDowntimeEpoch, key.Bytes(), &BlocksMissed{}).(*BlocksMissed)
	if pMissed == nil {
		return BlocksMissed{}
	}
	missed := *pMissed

	return missed
}

// GetActiveValidationScore return gas value for active validator score
func (s *Store) GetActiveValidationScore(stakerID idx.StakerID) *big.Int {
	return s.getValidationScore(s.table.ActiveValidationScore, stakerID)
}

// GetActiveValidationScore return gas value for active validator score
func (s *Store) GetActiveValidationScoreEpoch(stakerID idx.StakerID, epoch idx.Epoch) *big.Int {
	return s.getValidationScoreEpoch(s.table.ActiveValidationScoreEpoch, stakerID, epoch)
}

// AddDirtyValidationScore add gas value for active validation score
func (s *Store) AddDirtyValidationScore(stakerID idx.StakerID, v *big.Int) {
	s.addValidationScore(s.table.DirtyValidationScore, stakerID, v)
}

// DelActiveValidationScore deletes record about active validation score of a staker
func (s *Store) DelActiveValidationScore(stakerID idx.StakerID) {
	err := s.table.ActiveValidationScore.Delete(stakerID.Bytes())
	if err != nil {
		s.Log.Crit("Failed to erase key-value", "err", err)
	}
}

// DelDirtyValidationScore deletes record about dirty validation score of a staker
func (s *Store) DelDirtyValidationScore(stakerID idx.StakerID) {
	err := s.table.DirtyValidationScore.Delete(stakerID.Bytes())
	if err != nil {
		s.Log.Crit("Failed to erase key-value", "err", err)
	}
}

// GetDirtyValidationScore return gas value for active validator score
func (s *Store) GetDirtyValidationScore(stakerID idx.StakerID) *big.Int {
	return s.getValidationScore(s.table.DirtyValidationScore, stakerID)
}

func (s *Store) addValidationScore(t kvdb.KeyValueStore, stakerID idx.StakerID, diff *big.Int) {
	score := s.getValidationScore(t, stakerID)
	score.Add(score, diff)
	err := t.Put(stakerID.Bytes(), score.Bytes())
	if err != nil {
		s.Log.Crit("Failed to set key-value", "err", err)
	}
}

func (s *Store) getValidationScore(t kvdb.KeyValueStore, stakerID idx.StakerID) *big.Int {
	scoreBytes, err := t.Get(stakerID.Bytes())
	if err != nil {
		s.Log.Crit("Failed to get key-value", "err", err)
	}
	if scoreBytes == nil {
		return big.NewInt(0)
	}
	return new(big.Int).SetBytes(scoreBytes)
}

func (s *Store) getValidationScoreEpoch(t kvdb.KeyValueStore, stakerID idx.StakerID, epoch idx.Epoch) *big.Int {
	key := bytes.Buffer{}
	key.Write(epoch.Bytes())
	key.Write(stakerID.Bytes())

	scoreBytes, err := t.Get(key.Bytes())
	if err != nil {
		s.Log.Crit("Failed to get key-value", "err", err)
	}
	if scoreBytes == nil {
		return big.NewInt(0)
	}
	return new(big.Int).SetBytes(scoreBytes)
}

// DelAllActiveValidationScores deletes all the record about dirty validation scores of stakers
func (s *Store) DelAllActiveValidationScores() {
	it := s.table.ActiveValidationScore.NewIterator()
	defer it.Release()
	s.dropTable(it, s.table.ActiveValidationScore)
}

// MoveDirtyValidationScoresToActive moves all the dirty records to active
func (s *Store) MoveDirtyValidationScoresToActive() {
	it := s.table.DirtyValidationScore.NewIterator()
	defer it.Release()

	keys := make([][]byte, 0, 500) // don't write during iteration
	vals := make([][]byte, 0, 500)

	for it.Next() {
		keys = append(keys, it.Key())
		vals = append(vals, it.Value())
	}

	for i := range keys {
		err := s.table.ActiveValidationScore.Put(keys[i], vals[i])
		if err != nil {
			s.Log.Crit("Failed to set key-value", "err", err)
		}
		err = s.table.DirtyValidationScore.Delete(keys[i])
		if err != nil {
			s.Log.Crit("Failed to erase key-value", "err", err)
		}
	}
}

// GetActiveOriginationScore return gas value for active validator score
func (s *Store) GetActiveOriginationScore(stakerID idx.StakerID) *big.Int {
	return s.getOriginationScore(s.table.ActiveOriginationScore, stakerID)
}

// AddDirtyOriginationScore add gas value for active validation score
func (s *Store) AddDirtyOriginationScore(stakerID idx.StakerID, v *big.Int) {
	s.addOriginationScore(s.table.DirtyOriginationScore, stakerID, v)
}

// DelActiveOriginationScore deletes record about active origination score of a staker
func (s *Store) DelActiveOriginationScore(stakerID idx.StakerID) {
	err := s.table.ActiveOriginationScore.Delete(stakerID.Bytes())
	if err != nil {
		s.Log.Crit("Failed to erase key-value", "err", err)
	}
}

// DelDirtyOriginationScore deletes record about dirty origination score of a staker
func (s *Store) DelDirtyOriginationScore(stakerID idx.StakerID) {
	err := s.table.DirtyOriginationScore.Delete(stakerID.Bytes())
	if err != nil {
		s.Log.Crit("Failed to erase key-value", "err", err)
	}
}

// GetDirtyOriginationScore return gas value for active validator score
func (s *Store) GetDirtyOriginationScore(stakerID idx.StakerID) *big.Int {
	return s.getOriginationScore(s.table.DirtyOriginationScore, stakerID)
}

func (s *Store) addOriginationScore(t kvdb.KeyValueStore, stakerID idx.StakerID, diff *big.Int) {
	score := s.getOriginationScore(t, stakerID)
	score.Add(score, diff)
	err := t.Put(stakerID.Bytes(), score.Bytes())
	if err != nil {
		s.Log.Crit("Failed to set key-value", "err", err)
	}
}

func (s *Store) getOriginationScore(t kvdb.KeyValueStore, stakerID idx.StakerID) *big.Int {
	scoreBytes, err := t.Get(stakerID.Bytes())
	if err != nil {
		s.Log.Crit("Failed to get key-value", "err", err)
	}
	if scoreBytes == nil {
		return big.NewInt(0)
	}
	return new(big.Int).SetBytes(scoreBytes)
}

// DelAllActiveOriginationScores deletes all the record about dirty origination scores of stakers
func (s *Store) DelAllActiveOriginationScores() {
	it := s.table.ActiveOriginationScore.NewIterator()
	defer it.Release()
	s.dropTable(it, s.table.ActiveOriginationScore)
}

// MoveDirtyOriginationScoresToActive moves all the dirty records to active
func (s *Store) MoveDirtyOriginationScoresToActive() {
	it := s.table.DirtyOriginationScore.NewIterator()
	defer it.Release()

	keys := make([][]byte, 0, 500) // don't write during iteration
	vals := make([][]byte, 0, 500)

	for it.Next() {
		keys = append(keys, it.Key())
		vals = append(vals, it.Value())
	}

	for i := range keys {
		err := s.table.ActiveOriginationScore.Put(keys[i], vals[i])
		if err != nil {
			s.Log.Crit("Failed to set key-value", "err", err)
		}
		err = s.table.DirtyOriginationScore.Delete(keys[i])
		if err != nil {
			s.Log.Crit("Failed to erase key-value", "err", err)
		}
	}
}

func (s *Store) newEpochSnapshot(from, to kvdb.KeyValueStore, epoch idx.Epoch) {
	it := from.NewIterator()
	for it.Next() {
		k := it.Key()
		v := it.Value()

		newKey := bytes.Buffer{}
		newKey.Write(epoch.Bytes())
		newKey.Write(k)

		err := to.Put(newKey.Bytes(), v)
		if err != nil {
			s.Log.Error("error when write to epoch snapshot table")
		}
	}
}
