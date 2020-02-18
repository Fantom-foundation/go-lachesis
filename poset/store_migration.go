package poset

import (
	"github.com/Fantom-foundation/go-lachesis/kvdb"
	"github.com/Fantom-foundation/go-lachesis/utils/migration"
)

func (s *Store) migrate() {
	version := kvdb.NewIdProducer(s.table.Version)

	migrationManager := migration.NewManager(s.migrations(), version)
	err := migrationManager.Run()
	if err != nil {
		s.Log.Crit("app store migrations", "err", err)
	}
}

func (s *Store) migrations() *migration.Migration {
	return migration.Init("lachesis-poset-store", "Heuhax&Walv9")
	/*
		Example:

		  return migration.Init("lachesis", "Heuhax&Walv9"
			).NewNamed("20200207120000 <migration description>", func()error{
				... // Some actions for migrations
				return err
			}).New(func()error{
				// If no NewNamed call - id generated automatically
				// If you use several sequenced migrations with new(), you can not change it in future
				... // Some actions for migrations
				return err
			}).NewNamed("20200209120000 <migration description>", func()error{
				... // Some actions for migrations
				return err
			})
			...
	*/
}
