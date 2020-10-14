package gossip

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"

	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/sfctype"
	"github.com/Fantom-foundation/go-lachesis/kvdb/table"
	"github.com/Fantom-foundation/go-lachesis/utils/migration"
)

func isEmptyDB(db ethdb.Iteratee) bool {
	it := db.NewIterator(nil, nil)
	defer it.Release()
	return !it.Next()
}

func (s *Store) migrate() {
	versions := migration.NewKvdbIDStore(s.table.Version)
	if isEmptyDB(s.mainDb) && isEmptyDB(s.async.mainDb) {
		// short circuit if empty DB
		versions.SetID(s.migrations().ID())
		return
	}
	err := s.migrations().Exec(versions)
	if err != nil {
		s.Log.Crit("app store migrations", "err", err)
	}
}

func (s *Store) migrations() *migration.Migration {
	return migration.
		Begin("lachesis-gossip-store").
		Next("remove async data from sync DBs",
			func() error {
				s.rmPrefix(s.table.PackInfos, "serverPool")
				s.rmPrefix(s.mainDb, "Z")
				return nil
			}).
		Next("remove legacy genesis field",
			s.migrateEraseGenesisField).
		Next("multi-delegations",
			s.migrateMultiDelegations).
		Next("adjustable offline pruning time",
			s.migrateAdjustableOfflinePeriod).
		Next("adjustable minimum gas price",
			s.migrateAdjustableMinGasPrice)
}

func (s *Store) migrateEraseGenesisField() error {
	it := s.mainDb.NewIterator([]byte("G"), nil)
	defer it.Release()
	s.dropTable(it, s.mainDb)
	return nil
}

type (
	legacySfcDelegation struct {
		CreatedEpoch idx.Epoch
		CreatedTime  inter.Timestamp

		DeactivatedEpoch idx.Epoch
		DeactivatedTime  inter.Timestamp

		Amount *big.Int

		ToStakerID idx.StakerID
	}

	legacySfcConstants_06 struct {
		ShortGasPowerAllocPerSec uint64
		LongGasPowerAllocPerSec  uint64
		BaseRewardPerSec         *big.Int
	}

	legacyBlocksMissed struct {
		Num    idx.Block
		Period inter.Timestamp
	}

	legacySfcConstants_07rc1 struct {
		ShortGasPowerAllocPerSec uint64
		LongGasPowerAllocPerSec  uint64
		BaseRewardPerSec         *big.Int
		OfflinePenaltyThreshold  legacyBlocksMissed
	}

	legacySfcConstants_07rc3 struct {
		ShortGasPowerAllocPerSec uint64
		LongGasPowerAllocPerSec  uint64
		BaseRewardPerSec         *big.Int
		OfflinePenaltyThreshold  legacyBlocksMissed
		MinGasPrice              *big.Int
	}
)

func (s *Store) getLegacySfcDelegationsByAddr(addr common.Address, limit int) []sfctype.SfcDelegationAndID {
	forEachSfcDelegation := func(it ethdb.Iterator, do func(sfctype.SfcDelegationAndID) bool) {
		_continue := true
		for _continue && it.Next() {
			delegation := &sfctype.SfcDelegation{}
			err := rlp.DecodeBytes(it.Value(), delegation)
			if err != nil {
				s.Log.Crit("Failed to decode rlp while iteration", "err", err)
			}

			addr := it.Key()[len(it.Key())-sfctype.DelegationIDSize:]
			_continue = do(sfctype.SfcDelegationAndID{
				ID:         sfctype.BytesToDelegationID(addr),
				Delegation: delegation,
			})
		}
	}

	tableDelegations := table.New(s.mainDb, []byte("3"))
	it := tableDelegations.NewIterator(addr.Bytes(), nil)
	defer it.Release()
	res := make([]sfctype.SfcDelegationAndID, 0, limit)

	forEachSfcDelegation(it, func(id sfctype.SfcDelegationAndID) bool {
		if limit == 0 {
			return false
		}
		limit--
		res = append(res, id)
		return true
	})
	return res
}

func (s *Store) migrateMultiDelegations() error {
	{ // migrate s.table.Delegations
		tableDelegations := table.New(s.mainDb, []byte("3"))

		newKeys := make([][]byte, 0, 10000)
		newValues := make([][]byte, 0, 10000)
		{

			it := tableDelegations.NewIterator(nil, nil)
			defer it.Release()
			for it.Next() {
				delegation := &legacySfcDelegation{}
				err := rlp.DecodeBytes(it.Value(), delegation)
				if err != nil {
					return errors.Wrap(err, "failed legacy delegation deserialization during migration")
				}

				addr := common.BytesToAddress(it.Key())
				id := sfctype.DelegationID{
					Delegator: addr,
					StakerID:  delegation.ToStakerID,
				}
				newValue, err := rlp.EncodeToBytes(sfctype.SfcDelegation{
					CreatedEpoch:     delegation.CreatedEpoch,
					CreatedTime:      delegation.CreatedTime,
					DeactivatedEpoch: delegation.DeactivatedEpoch,
					DeactivatedTime:  delegation.DeactivatedTime,
					Amount:           delegation.Amount,
				})
				if err != nil {
					return err
				}

				// don't write into DB during iteration
				newKeys = append(newKeys, id.Bytes())
				newValues = append(newValues, newValue)
			}
		}
		{
			it := tableDelegations.NewIterator(nil, nil)
			defer it.Release()
			s.dropTable(it, tableDelegations)
		}
		for i := range newKeys {
			err := tableDelegations.Put(newKeys[i], newValues[i])
			if err != nil {
				return err
			}
		}
	}
	{ // migrate s.table.DelegationOldRewards
		tableDelegationOldRewards := table.New(s.mainDb, []byte("6"))
		newKeys := make([][]byte, 0, 10000)
		newValues := make([][]byte, 0, 10000)
		{
			it := tableDelegationOldRewards.NewIterator(nil, nil)
			defer it.Release()
			for it.Next() {
				addr := common.BytesToAddress(it.Key())
				delegations := s.getLegacySfcDelegationsByAddr(addr, 2)
				if len(delegations) > 1 {
					return errors.New("more than one delegation during multi-delegation migration")
				}
				if len(delegations) == 0 {
					continue
				}
				toStakerID := delegations[0].ID.StakerID
				id := sfctype.DelegationID{
					Delegator: addr,
					StakerID:  toStakerID,
				}

				// don't write into DB during iteration
				newKeys = append(newKeys, id.Bytes())
				newValues = append(newValues, it.Value())
			}
		}
		{
			it := tableDelegationOldRewards.NewIterator(nil, nil)
			defer it.Release()
			s.dropTable(it, tableDelegationOldRewards)
		}
		for i := range newKeys {
			err := tableDelegationOldRewards.Put(newKeys[i], newValues[i])
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (s *Store) migrateAdjustableOfflinePeriod() error {
	{ // migrate s.table.SfcConstants
		tableSfcConstants := table.New(s.mainDb, []byte("4"))
		newKeys := make([][]byte, 0, 10000)
		newValues := make([][]byte, 0, 10000)
		{
			it := tableSfcConstants.NewIterator(nil, nil)
			defer it.Release()
			for it.Next() {
				constants := &legacySfcConstants_06{}
				err := rlp.DecodeBytes(it.Value(), constants)
				if err != nil {
					return errors.Wrap(err, "failed legacy constants deserialization during migration")
				}

				newConstants := legacySfcConstants_07rc1{
					ShortGasPowerAllocPerSec: constants.ShortGasPowerAllocPerSec,
					LongGasPowerAllocPerSec:  constants.LongGasPowerAllocPerSec,
					BaseRewardPerSec:         constants.BaseRewardPerSec,
				}
				newValue, err := rlp.EncodeToBytes(newConstants)
				if err != nil {
					return err
				}

				// don't write into DB during iteration
				newKeys = append(newKeys, it.Key())
				newValues = append(newValues, newValue)
			}
		}
		{
			it := tableSfcConstants.NewIterator(nil, nil)
			defer it.Release()
			s.dropTable(it, tableSfcConstants)
		}
		for i := range newKeys {
			err := tableSfcConstants.Put(newKeys[i], newValues[i])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Store) migrateAdjustableMinGasPrice() error {
	{ // migrate s.table.SfcConstants
		tableSfcConstants := table.New(s.mainDb, []byte("4"))
		newKeys := make([][]byte, 0, 10000)
		newValues := make([][]byte, 0, 10000)
		{
			it := tableSfcConstants.NewIterator(nil, nil)
			defer it.Release()
			for it.Next() {
				constants := &legacySfcConstants_07rc1{}
				err := rlp.DecodeBytes(it.Value(), constants)
				if err != nil {
					return errors.Wrap(err, "failed legacy constants deserialization during migration")
				}

				newConstants := legacySfcConstants_07rc3{
					ShortGasPowerAllocPerSec: constants.ShortGasPowerAllocPerSec,
					LongGasPowerAllocPerSec:  constants.LongGasPowerAllocPerSec,
					BaseRewardPerSec:         constants.BaseRewardPerSec,
					OfflinePenaltyThreshold:  constants.OfflinePenaltyThreshold,
					MinGasPrice:              big.NewInt(0),
				}
				newValue, err := rlp.EncodeToBytes(newConstants)
				if err != nil {
					return err
				}

				// don't write into DB during iteration
				newKeys = append(newKeys, it.Key())
				newValues = append(newValues, newValue)
			}
		}
		{
			it := tableSfcConstants.NewIterator(nil, nil)
			defer it.Release()
			s.dropTable(it, tableSfcConstants)
		}
		for i := range newKeys {
			err := tableSfcConstants.Put(newKeys[i], newValues[i])
			if err != nil {
				return err
			}
		}
	}
	return nil
}
