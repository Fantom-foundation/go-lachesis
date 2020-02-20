package gossip

import (
	"github.com/Fantom-foundation/go-lachesis/kvdb"
	"github.com/Fantom-foundation/go-lachesis/kvdb/table"
	"github.com/Fantom-foundation/go-lachesis/utils/migration"
)

func (s *Store) migrate() {
	versions := kvdb.NewIdStore(s.table.Version)
	err := s.migrations().Exec(versions)
	if err != nil {
		s.Log.Crit("gossip store migrations", "err", err)
	}
	err = s.Commit(nil, true)
	if err != nil {
		s.Log.Crit("gossip store commit", "err", err)
	}
}

func (s *Store) migrations() *migration.Migration {
	return migration.
		Begin("lachesis-gossip-store").
		Next("service db",
			func() error {
				dst := s.service.Peers

				old1 := s.table.PackInfos
				err := s.move(old1, dst, []byte("serverPool"))
				if err != nil {
					return err
				}

				old2 := table.New(s.mainDb, []byte("Z"))
				err = s.move(old2, dst, nil)

				return err
			})
}
