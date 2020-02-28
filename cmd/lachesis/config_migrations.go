package main

import (
	"github.com/naoina/toml/ast"

	"github.com/Fantom-foundation/go-lachesis/utils/migration"
	"github.com/Fantom-foundation/go-lachesis/utils/toml"
)

func (c *config) migrate(t *ast.Table) (before string, after string) {
	cfgData := toml.NewTomlHelper(t)
	migrations := c.migrations(cfgData)
	versions := toml.NewIDStore(cfgData, migrations.IdChain())

	before, _ = cfgData.GetParamString("Version", "")

	err := migrations.Exec(versions)
	if err != nil && err != toml.ErrorParamNotExists {
		panic(err)
	}

	after, _ = cfgData.GetParamString("Version", "")

	return
}

func (c *config) migrations(data *toml.Helper) *migration.Migration {
	return migration.Begin("lachesis-config")
}
