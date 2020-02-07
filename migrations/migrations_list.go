package migrations

import "github.com/Fantom-foundation/go-lachesis/utils/migration"

func List() *migration.Migration {
	/*
		Example:

		  return migration.New("init", nil, func()error{
				return nil
			}).New("20200207120000 <migration description>", func()error{
				... // Some actions for migrations
				return err
			}).New("20200208120000 <migration description>", func()error{
				... // Some actions for migrations
				return err
			}).New("20200209120000 <migration description>", func()error{
				... // Some actions for migrations
				return err
			})
			...
	*/

	return migration.New("init", nil, func()error{
		return nil
	})
}
