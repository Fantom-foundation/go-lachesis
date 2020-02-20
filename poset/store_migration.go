package poset

import (
	"github.com/Fantom-foundation/go-lachesis/kvdb"
	"github.com/Fantom-foundation/go-lachesis/utils/migration"
)

func (s *Store) migrate() {
	versions := kvdb.NewIdStore(s.table.Version)
	err := s.migrations().Exec(versions)
	if err != nil {
		s.Log.Crit("poset store migrations", "err", err)
	}
}

func (s *Store) migrations() *migration.Migration {
	return migration.Begin("lachesis-poset-store")
}
