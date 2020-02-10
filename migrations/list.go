package migrations

import "github.com/Fantom-foundation/go-lachesis/utils/migration"

func List() *migration.Migration {
	/*
		Example:

		  return migration.New("init", nil, func()error{
				return nil
			}).NewNamed("20200207120000 <migration description>", func()error{
				... // Some actions for migrations
				return err
			}).New(func()error{
				// If no NewNamed call - id generated automatically
				// If you use several sequenced migrations with New(), you can not change it in future
				... // Some actions for migrations
				return err
			}).NewNamed("20200209120000 <migration description>", func()error{
				... // Some actions for migrations
				return err
			})
			...
	*/

	return migration.New("init", nil, func()error{
		return nil
	})
}
