package gossip

import (
	"github.com/Fantom-foundation/go-lachesis/kvdb"
	"github.com/Fantom-foundation/go-lachesis/utils/migration"
)

func (s *Store) migrate() {
	version := kvdb.NewIdProducer(s.table.Version)

	migrationManager := migration.NewManager(s.migrations(), version)
	err := migrationManager.Run()
	if err != nil {
		s.Log.Crit("gossip store migrations", "err", err)
	}
	err = s.Commit(nil, true)
	if err != nil {
		s.Log.Crit("gossip store commit", "err", err)
	}

}

func (s *Store) migrations() *migration.Migration {
	return migration.Begin("lachesis-gossip-store")
}
