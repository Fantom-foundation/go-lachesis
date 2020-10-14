package app

import (
	"github.com/ethereum/go-ethereum/ethdb"

	"github.com/Fantom-foundation/go-lachesis/utils/migration"
)

func isEmptyDB(db ethdb.Iteratee) bool {
	it := db.NewIterator(nil, nil)
	defer it.Release()
	return !it.Next()
}

func (s *Store) migrate() {
	return // TODO: enable when dedicated db

	versions := migration.NewKvdbIDStore(s.table.Version)
	if isEmptyDB(s.mainDb) {
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
		Begin("lachesis-app-store")
}
