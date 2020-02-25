package main

import (
	"github.com/naoina/toml/ast"

	"github.com/Fantom-foundation/go-lachesis/utils/migration"
	"github.com/Fantom-foundation/go-lachesis/utils/toml"
)

/*
	Methods for migrations (*toml_helper.Helper):
	- AddSection
	- DeleteSection
	- RenameSection

	- AddParam
	- DeleteParam
	- RenameParam
	- SetParam
	- GetParam[String|Int|Float|Bool|Time]
*/

func (c *config) migrate(t *ast.Table) (oldVersion string, newVersion string) {
	cfgData := toml.NewTomlHelper(t)
	var err error
	oldVersion, _ = cfgData.GetParamString("Version", "")

	migrations := c.migrations(cfgData)
	idProd := toml.NewIDStore(cfgData, migrations.IdChain())
	err = migrations.Exec(idProd)
	if err != nil && err != toml.ErrorParamNotExists {
		panic("error when run config migration: " + err.Error())
	}
	newVersion, err = cfgData.GetParamString("Version", "")
	if err != nil && err != toml.ErrorParamNotExists {
		panic("error when read new version after config migration: " + err.Error())
	}

	return
}

func (c *config) migrations(data *toml.Helper) *migration.Migration {
	return migration.Begin("lachesis-config")

	/*
		Use here only named migrations. Migration name - version of config.
		Example ():

		  return migration.Begin("lachesis-config").
			Next("v1", func()error{
				... // Some actions for migrations
				return err
			}).
			Next("v2", func()error{
				... // Some actions for migrations
				return err
			})
			...
	*/

}
