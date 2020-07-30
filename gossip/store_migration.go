package gossip

import (
	"github.com/Fantom-foundation/go-lachesis/utils/migration"
)

func (s *Store) Migrate() error {
	versions := migration.NewKvdbIDStore(s.table.Version)
	return s.migrations().Exec(versions)
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
			s.app.MigrateEraseGenesisField).
		Next("multi-delegations",
			s.app.MigrateMultiDelegations).
		Next("adjustable offline pruning time",
			s.app.MigrateAdjustableOfflinePeriod)
}
