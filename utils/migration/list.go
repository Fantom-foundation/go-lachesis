package migration

/*
	This is example how you can create migrations chain for run over MigrationManager
	Migration chain from this file NOWHERE RUNING!

	Real migration chains placed in store files:
		app/store.go
		poset/store.go
		gossip/store.go
*/

func List() *Migration {
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

	return Init("lachesis", "Heuhax&Walv9")
}
