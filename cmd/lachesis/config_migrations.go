package main

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/naoina/toml/ast"

	"github.com/Fantom-foundation/go-lachesis/utils/migration"
	"github.com/Fantom-foundation/go-lachesis/utils/toml"
)

func (c *config) migrate(t *ast.Table) (changed bool, err error) {
	data := toml.NewTomlHelper(t)
	migrations := c.migrations(data)
	versions := toml.NewIDStore(data, migrations.IdChain())

	before := versions.GetID()

	err = migrations.Exec(versions)
	if err != nil && err != toml.ErrorParamNotExists {
		return
	}

	after := versions.GetID()
	changed = before != after

	return
}

func (c *config) migrations(data *toml.Helper) *migration.Migration {

	isCritical := func(err error) bool {
		return err != nil &&
			err != toml.ErrorSectionAlreadyExists &&
			err != toml.ErrorParamAlreadyExists
	}

	return migration.Begin("lachesis-config").
		Next("from v0.5.0-rc.1 & v0.6.0-rc.1", func() (err error) {
			defer func() {
				if err == nil {
					log.Warn("config migration from v0.5.0-rc.1 & v0.6.0-rc.1 has been applied")
				}
			}()

			_ = data.DeleteParam("omitempty", "Lachesis")
			_ = data.DeleteParam("omitempty", "Node")

			// detect v0.5.0-rc.1 version
			err = data.DeleteParam("ForcedBroadcast", "Lachesis")
			if err == toml.ErrorParamNotExists {
				return nil
			}

			// Lachesis
			err = data.AddSection("Lachesis.Protocol", "")
			if isCritical(err) {
				return err
			}

			// EmitterConfig:
			err = data.AddSection("Lachesis.Emitter.EmitIntervals", "")
			if isCritical(err) {
				return err
			}

			oldMin, err := data.GetParamInt("MinEmitInterval", "Lachesis.Emitter")
			if err == nil && oldMin >= 0 {
				err = data.AddParam("Min", "Lachesis.Emitter.EmitIntervals", oldMin)
				if isCritical(err) {
					return err
				}
			}

			oldMax, err := data.GetParamInt("MaxEmitInterval", "Lachesis.Emitter")
			if err == nil && oldMax >= 0 {
				err = data.AddParam("Max", "Lachesis.Emitter.EmitIntervals", oldMax)
				if isCritical(err) {
					return err
				}
			}

			oldSelfForkProtection, err := data.GetParamInt("SelfForkProtectionInterval", "Lachesis.Emitter")
			if err == nil && oldSelfForkProtection >= 0 {
				err = data.AddParam("SelfForkProtection", "Lachesis.Emitter.EmitIntervals", oldSelfForkProtection)
				if isCritical(err) {
					return err
				}
			}

			_ = data.DeleteParam("MinEmitInterval", "Lachesis.Emitter")
			_ = data.DeleteParam("MaxEmitInterval", "Lachesis.Emitter")
			_ = data.DeleteParam("SelfForkProtectionInterval", "Lachesis.Emitter")

			return nil

		}).
		Next("app's store", func() (err error) {
			err = data.AddSection("App.StoreConfig", "")
			if isCritical(err) {
				return
			}

			var val int64
			for _, param := range []string{
				"BlockCacheSize",
				"ReceiptsCacheSize",
				"StakersCacheSize",
				"DelegatorsCacheSize",
			} {
				val, err = data.GetParamInt(param, "Lachesis.StoreConfig")
				if err != nil {
					continue
				}
				err = data.AddParam(param, "App.StoreConfig", val)
				if isCritical(err) {
					return
				}
				_ = data.DeleteParam(param, "Lachesis.StoreConfig")
			}

			txi, err := data.GetParamBool("TxIndex", "Lachesis")
			if err == nil {
				err = data.AddParam("TxIndex", "App", txi)
				if isCritical(err) {
					return err
				}
			}

			log.Warn("app's store config migration has been applied")
			return

		})
}
