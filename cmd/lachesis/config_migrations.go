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
	return migration.Begin("lachesis-config").
		Next("From v0.5.0-rc.1 & v0.6.0-rc.1 to current", func()error{
			// v0.5.0-rc.1 -> HEAD

			_ = data.DeleteParam("omitempty", "Lachesis")
			_ = data.DeleteParam("omitempty", "Node")

			// If ForcedBroadcast exists - this is version v0.5.0-rc.1
			err := data.DeleteParam("ForcedBroadcast", "Lachesis")
			if err != toml.ErrorParamNotExists {
				// Lachesis
				err = data.AddSection("Lachesis.Protocol", "")
				if err != nil && err != toml.ErrorSectionAlreadyExists {
					return err
				}

				// EmitterConfig:
				err = data.AddSection("Lachesis.Emitter.EmitIntervals", "")
				if err != nil && err != toml.ErrorSectionAlreadyExists {
					return err
				}

				oldMin, err := data.GetParamInt("MinEmitInterval", "Lachesis.Emitter")
				if err == nil && oldMin >= 0 {
					err = data.AddParam("Min", "Lachesis.Emitter.EmitIntervals", oldMin)
					if err != nil && err != toml.ErrorParamAlreadyExists {
						return err
					}
				}

				oldMax, err := data.GetParamInt("MaxEmitInterval", "Lachesis.Emitter")
				if err == nil && oldMax >= 0 {
					err = data.AddParam("Max", "Lachesis.Emitter.EmitIntervals", oldMax)
					if err != nil && err != toml.ErrorParamAlreadyExists {
						return err
					}
				}

				oldSelfForkProtection, err := data.GetParamInt("MaxEmitInterval", "Lachesis.Emitter")
				if err == nil && oldSelfForkProtection >= 0 {
					err = data.AddParam("SelfForkProtection", "Lachesis.Emitter.EmitIntervals", oldSelfForkProtection)
					if err != nil && err != toml.ErrorParamAlreadyExists {
						return err
					}
				}

				_ = data.DeleteParam("MinEmitInterval", "Lachesis.Emitter")
				_ = data.DeleteParam("MaxEmitInterval", "Lachesis.Emitter")
				_ = data.DeleteParam("SelfForkProtectionInterval", "Lachesis.Emitter")
			}

			// v0.6.0-rc.1 -> HEAD
			// no changes

			return nil
		},
	)
}
