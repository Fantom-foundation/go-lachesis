package main

import (
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Fantom-foundation/go-lachesis/utils/migration"
)

func TestConfigParse(t *testing.T) {
	t.Run("Test parse fixture config with version", func(t *testing.T) {
		source := filepath.Join("testdata", "test_config.toml")
		modified := filepath.Join("testdata", "test_config_modified.toml")
		copy(source, modified)

		cfg := config{}
		err := loadAllConfigs(modified, &cfg)
		assert.NoError(t, err, "Parse fixture config without error")

		os.Remove(modified)
		os.Remove(modified + ".init")
	})
	t.Run("Test parse fixture config without version", func(t *testing.T) {
		source := filepath.Join("testdata", "test_config_wo_version.toml")
		modified := filepath.Join("testdata", "test_config_wo_version_modified.toml")
		copy(source, modified)

		cfg := config{}
		err := loadAllConfigs(modified, &cfg)
		assert.NoError(t, err, "Parse fixture config without error")

		os.Remove(modified)
		os.Remove(modified + ".init")
	})
}

func TestConfigMigrations(t *testing.T) {
	/*
		Test cases:
			- rename section
			- add section
			- delete section
			- rename param
			- add param
			- delete param
			- set param value
			- modify param value all types (String, Int, Float, Bool, Datetime)
	*/

	t.Run("Rename section", func(t *testing.T) {
		source := filepath.Join("testdata", "test_config_wo_version.toml")
		table, err := parseConfigToTable(source)
		assert.NoError(t, err, "Parse config to Table")
		data := NewConfigData(table)

		migrations := func(data *ConfigData) *migration.Migration {
			return migration.Begin("lachesis-config-test").
				Next("v1-test", func() error {
					err := data.RenameSection("Node", "NodeRenamed")
					return err
				})
		}(data)

		idProd := NewConfigIdProducer(data, migrations)
		err = migrations.Exec(idProd)

		assert.NoError(t, err, "Config migrations success run")

		newVersion, err := data.GetParamString("Version", "")
		assert.NoError(t, err, "Get version from config data after config migration")
		assert.Equal(t, "v1-test", newVersion, "Version right after config migrations")

		p, err := data.GetParamString("IPCPath", "NodeRenamed")
		assert.NoError(t, err, "Get param from renamed section after config migration")
		assert.Equal(t, "lachesis.ipc", p, "Param right from renamed section after config migrations")

		p, err = data.GetParamString("ListenAddr", "NodeRenamed.P2P")
		assert.NoError(t, err, "Get param from sub renamed section after config migration")
		assert.Equal(t, ":7946", p, "Param right from sub renamed section after config migrations")
	})

	t.Run("Add section", func(t *testing.T) {
		t.Run("in root", func(t *testing.T) {
			source := filepath.Join("testdata", "test_config_wo_version.toml")
			table, err := parseConfigToTable(source)
			assert.NoError(t, err, "Parse config to Table")
			data := NewConfigData(table)

			// Add root section
			migrations := func(data *ConfigData) *migration.Migration {
				return migration.Begin("lachesis-config-test").
					Next("v1-test", func() error {
						err := data.AddSection("NewSection", "")
						return err
					})
			}(data)

			idProd := NewConfigIdProducer(data, migrations)
			err = migrations.Exec(idProd)

			assert.NoError(t, err, "Config migrations success run")

			newVersion, err := data.GetParamString("Version", "")
			assert.NoError(t, err, "Get version from config data after config migration")
			assert.Equal(t, "v1-test", newVersion, "Version right after config migrations")

			sec, err := data.findSection("NewSection")
			assert.NoError(t, err, "New section search no error after add section in config migrations")
			assert.NotNil(t, sec, "New section exists after add section in config migrations")
		})

		t.Run("sub section", func(t *testing.T) {
			t.Run("parent not exists", func(t *testing.T) {
				source := filepath.Join("testdata", "test_config_wo_version.toml")
				table, err := parseConfigToTable(source)
				assert.NoError(t, err, "Parse config to Table")
				data := NewConfigData(table)

				// Add root section
				migrations := func(data *ConfigData) *migration.Migration {
					return migration.Begin("lachesis-config-test").
						Next("v1-test", func() error {
							err := data.AddSection("NewSection.NewSubsection.NewSubsection2", "")
							return err
						})
				}(data)

				idProd := NewConfigIdProducer(data, migrations)
				err = migrations.Exec(idProd)

				assert.NoError(t, err, "Config migrations success run")

				newVersion, err := data.GetParamString("Version", "")
				assert.NoError(t, err, "Get version from config data after config migration")
				assert.Equal(t, "v1-test", newVersion, "Version right after config migrations")

				sec, err := data.findSection("NewSection.NewSubsection.NewSubsection2")
				assert.NoError(t, err, "New section search no error after add section in config migrations")
				assert.NotNil(t, sec, "New section exists after add section in config migrations")
			})

			t.Run("parent exists", func(t *testing.T) {
				source := filepath.Join("testdata", "test_config_wo_version.toml")
				table, err := parseConfigToTable(source)
				assert.NoError(t, err, "Parse config to Table")
				data := NewConfigData(table)

				// Add root section
				migrations := func(data *ConfigData) *migration.Migration {
					return migration.Begin("lachesis-config-test").
						Next("v1-test", func() error {
							err := data.AddSection("Node.NewSubSection.NewSubsection2", "")
							return err
						})
				}(data)

				idProd := NewConfigIdProducer(data, migrations)
				err = migrations.Exec(idProd)

				assert.NoError(t, err, "Config migrations success run")

				newVersion, err := data.GetParamString("Version", "")
				assert.NoError(t, err, "Get version from config data after config migration")
				assert.Equal(t, "v1-test", newVersion, "Version right after config migrations")

				sec, err := data.findSection("Node.NewSubSection.NewSubsection2")
				assert.NoError(t, err, "New section search no error after add section in config migrations")
				assert.NotNil(t, sec, "New section exists after add section in config migrations")
			})
		})
	})

	t.Run("Delete section", func(t *testing.T) {
		t.Run("in root", func(t *testing.T) {
			source := filepath.Join("testdata", "test_config_wo_version.toml")
			table, err := parseConfigToTable(source)
			assert.NoError(t, err, "Parse config to Table")
			data := NewConfigData(table)

			// Add root section
			migrations := func(data *ConfigData) *migration.Migration {
				return migration.Begin("lachesis-config-test").
					Next("v1-test", func() error {
						err := data.DeleteSection("Node")
						return err
					})
			}(data)

			idProd := NewConfigIdProducer(data, migrations)
			err = migrations.Exec(idProd)

			assert.NoError(t, err, "Config migrations success run")

			newVersion, err := data.GetParamString("Version", "")
			assert.NoError(t, err, "Get version from config data after config migration")
			assert.Equal(t, "v1-test", newVersion, "Version right after config migrations")

			sec, err := data.findSection("Node")
			assert.Error(t, err, "Deleted section search should return error after delete section in config migrations")
			assert.Nil(t, sec, "Deleted section not exists after delete section in config migrations")
		})
		t.Run("in sub section", func(t *testing.T) {
			source := filepath.Join("testdata", "test_config_wo_version.toml")
			table, err := parseConfigToTable(source)
			assert.NoError(t, err, "Parse config to Table")
			data := NewConfigData(table)

			// Add root section
			migrations := func(data *ConfigData) *migration.Migration {
				return migration.Begin("lachesis-config-test").
					Next("v1-test", func() error {
						err := data.DeleteSection("Lachesis.Emitter.EmitIntervals")
						return err
					})
			}(data)

			idProd := NewConfigIdProducer(data, migrations)
			err = migrations.Exec(idProd)

			assert.NoError(t, err, "Config migrations success run")

			newVersion, err := data.GetParamString("Version", "")
			assert.NoError(t, err, "Get version from config data after config migration")
			assert.Equal(t, "v1-test", newVersion, "Version right after config migrations")

			sec, err := data.findSection("Lachesis.Emitter.EmitIntervals")
			assert.Error(t, err, "Deleted section search should return error after delete section in config migrations")
			assert.Nil(t, sec, "Deleted section not exists after delete section in config migrations")

			// Parent section should be exists after delete child
			sec, err = data.findSection("Lachesis.Emitter")
			assert.NoError(t, err, "Parent section search no error after delete section in config migrations")
			assert.NotNil(t, sec, "Parent section exists after delete section in config migrations")
		})
	})

	t.Run("Rename param", func(t *testing.T) {
		source := filepath.Join("testdata", "test_config_wo_version.toml")
		table, err := parseConfigToTable(source)
		assert.NoError(t, err, "Parse config to Table")
		data := NewConfigData(table)

		// Add root section
		migrations := func(data *ConfigData) *migration.Migration {
			return migration.Begin("lachesis-config-test").
				Next("v1-test", func() error {
					err := data.RenameParam("Validator", "Lachesis.Emitter", "ValidatorAddr")
					return err
				})
		}(data)

		idProd := NewConfigIdProducer(data, migrations)
		err = migrations.Exec(idProd)

		assert.NoError(t, err, "Config migrations success run")

		newVersion, err := data.GetParamString("Version", "")
		assert.NoError(t, err, "Get version from config data after config migration")
		assert.Equal(t, "v1-test", newVersion, "Version right after config migrations")

		paramStr, err := data.GetParamString("ValidatorAddr", "Lachesis.Emitter")
		assert.NoError(t, err, "Get param data no error after rename param in config migrations")
		assert.NotEmpty(t, paramStr, "Param data not empty after rename param in config migrations")
	})

	t.Run("Add param", func(t *testing.T) {
		t.Run("to root", func(t *testing.T) {
			source := filepath.Join("testdata", "test_config_wo_version.toml")
			table, err := parseConfigToTable(source)
			assert.NoError(t, err, "Parse config to Table")
			data := NewConfigData(table)

			// Add root section
			migrations := func(data *ConfigData) *migration.Migration {
				return migration.Begin("lachesis-config-test").
					Next("v1-test", func() error {
						err := data.AddParam("NewParamInRoot", "", "string value")
						return err
					})
			}(data)

			idProd := NewConfigIdProducer(data, migrations)
			err = migrations.Exec(idProd)

			assert.NoError(t, err, "Config migrations success run")

			newVersion, err := data.GetParamString("Version", "")
			assert.NoError(t, err, "Get version from config data after config migration")
			assert.Equal(t, "v1-test", newVersion, "Version right after config migrations")

			paramStr, err := data.GetParamString("NewParamInRoot", "")
			assert.NoError(t, err, "Get param data no error after add param to root in config migrations")
			assert.NotEmpty(t, paramStr, "Param data not empty after add param to root in config migrations")
		})
		t.Run("to section", func(t *testing.T) {
			source := filepath.Join("testdata", "test_config_wo_version.toml")
			table, err := parseConfigToTable(source)
			assert.NoError(t, err, "Parse config to Table")
			data := NewConfigData(table)

			// Add root section
			migrations := func(data *ConfigData) *migration.Migration {
				return migration.Begin("lachesis-config-test").
					Next("v1-test", func() error {
						err := data.AddParam("NewParamInSection", "Lachesis.Emitter", "string value")
						return err
					})
			}(data)

			idProd := NewConfigIdProducer(data, migrations)
			err = migrations.Exec(idProd)

			assert.NoError(t, err, "Config migrations success run")

			newVersion, err := data.GetParamString("Version", "")
			assert.NoError(t, err, "Get version from config data after config migration")
			assert.Equal(t, "v1-test", newVersion, "Version right after config migrations")

			paramStr, err := data.GetParamString("NewParamInSection", "Lachesis.Emitter")
			assert.NoError(t, err, "Get param data no error after add param in section in config migrations")
			assert.NotEmpty(t, paramStr, "Param data not empty after add param in section in config migrations")
		})
	})

	t.Run("Delete param", func(t *testing.T) {
		source := filepath.Join("testdata", "test_config_wo_version.toml")
		table, err := parseConfigToTable(source)
		assert.NoError(t, err, "Parse config to Table")
		data := NewConfigData(table)

		// Add root section
		migrations := func(data *ConfigData) *migration.Migration {
			return migration.Begin("lachesis-config-test").
				Next("v1-test", func() error {
					err := data.DeleteParam("Validator", "Lachesis.Emitter")
					return err
				})
		}(data)

		idProd := NewConfigIdProducer(data, migrations)
		err = migrations.Exec(idProd)

		assert.NoError(t, err, "Config migrations success run")

		newVersion, err := data.GetParamString("Version", "")
		assert.NoError(t, err, "Get version from config data after config migration")
		assert.Equal(t, "v1-test", newVersion, "Version right after config migrations")

		paramStr, err := data.GetParamString("Validator", "Lachesis.Emitter")
		assert.Error(t, err, "Get param data return error after delete param in config migrations")
		assert.Empty(t, paramStr, "Param data is empty after delete param in config migrations")
	})

	t.Run("Set param value", func(t *testing.T) {
		t.Run("String", func(t *testing.T) {
			source := filepath.Join("testdata", "test_config_wo_version.toml")
			table, err := parseConfigToTable(source)
			assert.NoError(t, err, "Parse config to Table")
			data := NewConfigData(table)

			// Add root section
			migrations := func(data *ConfigData) *migration.Migration {
				return migration.Begin("lachesis-config-test").
					Next("v1-test", func() error {
						err := data.SetParam("Validator", "Lachesis.Emitter", "test value")
						return err
					})
			}(data)

			idProd := NewConfigIdProducer(data, migrations)
			err = migrations.Exec(idProd)

			assert.NoError(t, err, "Config migrations success run")

			newVersion, err := data.GetParamString("Version", "")
			assert.NoError(t, err, "Get version from config data after config migration")
			assert.Equal(t, "v1-test", newVersion, "Version right after config migrations")

			paramStr, err := data.GetParamString("Validator", "Lachesis.Emitter")
			assert.NoError(t, err, "Get param data return no error after set param in config migrations")
			assert.Equal(t, "test value", paramStr, "Param data correct after set param in config migrations")
		})
		t.Run("Int", func(t *testing.T) {
			source := filepath.Join("testdata", "test_config_wo_version.toml")
			table, err := parseConfigToTable(source)
			assert.NoError(t, err, "Parse config to Table")
			data := NewConfigData(table)

			// Add root section
			migrations := func(data *ConfigData) *migration.Migration {
				return migration.Begin("lachesis-config-test").
					Next("v1-test", func() error {
						err := data.SetParam("MaxTxsFromSender", "Lachesis.Emitter", 1112233)
						return err
					})
			}(data)

			idProd := NewConfigIdProducer(data, migrations)
			err = migrations.Exec(idProd)

			assert.NoError(t, err, "Config migrations success run")

			newVersion, err := data.GetParamString("Version", "")
			assert.NoError(t, err, "Get version from config data after config migration")
			assert.Equal(t, "v1-test", newVersion, "Version right after config migrations")

			paramInt, err := data.GetParamInt("MaxTxsFromSender", "Lachesis.Emitter")
			assert.NoError(t, err, "Get param data return no error after set int param in config migrations")
			assert.Equal(t, int64(1112233), paramInt, "Param data correct after set int param in config migrations")
		})
		t.Run("Float", func(t *testing.T) {
			source := filepath.Join("testdata", "test_config_wo_version.toml")
			table, err := parseConfigToTable(source)
			assert.NoError(t, err, "Parse config to Table")
			data := NewConfigData(table)

			// Add root section
			migrations := func(data *ConfigData) *migration.Migration {
				return migration.Begin("lachesis-config-test").
					Next("v1-test", func() error {
						err := data.SetParam("MaxGasRateGrowthFactor", "Lachesis.Emitter", 1.234)
						return err
					})
			}(data)

			idProd := NewConfigIdProducer(data, migrations)
			err = migrations.Exec(idProd)

			assert.NoError(t, err, "Config migrations success run")

			newVersion, err := data.GetParamString("Version", "")
			assert.NoError(t, err, "Get version from config data after config migration")
			assert.Equal(t, "v1-test", newVersion, "Version right after config migrations")

			paramFloat, err := data.GetParamFloat("MaxGasRateGrowthFactor", "Lachesis.Emitter")
			assert.NoError(t, err, "Get param data return no error after set float param in config migrations")
			assert.Equal(t, 1.234, paramFloat, "Param data correct after set float param in config migrations")
		})
		t.Run("Bool", func(t *testing.T) {
			source := filepath.Join("testdata", "test_config_wo_version.toml")
			table, err := parseConfigToTable(source)
			assert.NoError(t, err, "Parse config to Table")
			data := NewConfigData(table)

			// Add root section
			migrations := func(data *ConfigData) *migration.Migration {
				return migration.Begin("lachesis-config-test").
					Next("v1-test", func() error {
						err := data.SetParam("TxIndex", "Lachesis", false)
						return err
					})
			}(data)

			idProd := NewConfigIdProducer(data, migrations)
			err = migrations.Exec(idProd)

			assert.NoError(t, err, "Config migrations success run")

			newVersion, err := data.GetParamString("Version", "")
			assert.NoError(t, err, "Get version from config data after config migration")
			assert.Equal(t, "v1-test", newVersion, "Version right after config migrations")

			paramBool, err := data.GetParamBool("TxIndex", "Lachesis")
			assert.NoError(t, err, "Get param data return no error after set bool param in config migrations")
			assert.False(t, paramBool, "Param data correct after set bool param in config migrations")
		})
		t.Run("DateTime", func(t *testing.T) {
			source := filepath.Join("testdata", "test_config_wo_version.toml")
			table, err := parseConfigToTable(source)
			assert.NoError(t, err, "Parse config to Table")
			data := NewConfigData(table)

			testTime0 := time.Now().UTC().Add(time.Hour)
			testTime := time.Now().UTC()

			// Add root section
			migrations := func(data *ConfigData) *migration.Migration {
				return migration.Begin("lachesis-config-test").
					Next("v1-test", func() error {
						err := data.AddParam("TestTime", "Node", testTime0)
						if err != nil {
							return err
						}
						err = data.SetParam("TestTime", "Node", testTime)
						return err
					})
			}(data)

			idProd := NewConfigIdProducer(data, migrations)
			err = migrations.Exec(idProd)

			assert.NoError(t, err, "Config migrations success run")

			newVersion, err := data.GetParamString("Version", "")
			assert.NoError(t, err, "Get version from config data after config migration")
			assert.Equal(t, "v1-test", newVersion, "Version right after config migrations")

			paramTime, err := data.GetParamTime("TestTime", "Node")
			assert.NoError(t, err, "Get param data return no error after set datetime param in config migrations")
			assert.Equal(t, testTime, paramTime, "Param data correct after set datetime param in config migrations")
		})
	})
}

func copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
