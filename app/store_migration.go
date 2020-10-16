package app

import (
	"github.com/ethereum/go-ethereum/ethdb"

	"github.com/Fantom-foundation/go-lachesis/kvdb/flushable"
	"github.com/Fantom-foundation/go-lachesis/utils/migration"
)

func isEmptyDB(db ethdb.Iteratee) bool {
	it := db.NewIterator(nil, nil)
	defer it.Release()
	return !it.Next()
}

func (s *Store) migrate() {
	versions := migration.NewKvdbIDStore(s.table.Version)
	migrations := s.migrations(s.dbs)
	if isEmptyDB(s.mainDb) {
		// short circuit if empty DB
		versions.SetID(migrations.ID())
		return
	}
	err := migrations.Exec(versions)
	if err != nil {
		s.Log.Crit("app store migrations", "err", err)
	}
}

func (s *Store) migrations(dbs *flushable.SyncedPool) *migration.Migration {
	return migration.Begin("lachesis-app-store")
}
