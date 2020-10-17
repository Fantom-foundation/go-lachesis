package app

import (
	"github.com/Fantom-foundation/go-lachesis/kvdb"
	"github.com/Fantom-foundation/go-lachesis/utils/migration"
)

func (s *Store) migrate() {
	versions := migration.NewKvdbIDStore(s.table.Version)
	migrations := s.migrations()
	if kvdb.IsEmptyDB(s.mainDb) {
		// short circuit if empty DB
		versions.SetID(migrations.ID())
		return
	}
	err := migrations.Exec(versions)
	if err != nil {
		s.Log.Crit("app store migrations", "err", err)
	}
}

func (s *Store) migrations() *migration.Migration {
	return migration.Begin("lachesis-app-store")
}
