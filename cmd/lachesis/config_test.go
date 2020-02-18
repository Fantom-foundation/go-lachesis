package main

import (
	"github.com/Fantom-foundation/go-lachesis/utils/migration"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestConfigParse(t *testing.T) {
	t.Run("Test parse fixture config with version", func(t *testing.T){
		source := filepath.Join("testdata", "test_config.toml")

		cfg := config{}
		err := loadAllConfigs(source, &cfg)
		assert.NoError(t, err, "Parse fixture config without error")
	})
	t.Run("Test parse fixture config without version", func(t *testing.T){
		source := filepath.Join("testdata", "test_config_wo_version.toml")

		cfg := config{}
		err := loadAllConfigs(source, &cfg)
		assert.NoError(t, err, "Parse fixture config without error")
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

	t.Run("Rename section", func(t *testing.T){
		source := filepath.Join("testdata", "test_config_wo_version.toml")
		table, err := parseConfigToTable(source)
		assert.NoError(t, err, "Parse config to Table")
		data := NewConfigData(table)

		migrations := func(data *ConfigData) *migration.Migration {
			return migration.Init("lachesis-config-test", "ajIr@Quicuj9").
			NewNamed("v1-test", func()error{
				err := data.RenameSection("Node", "NodeRenamed")
				return err
			})
		}(data)

		idProd := NewConfigIdProducer(data)
		migrationManager := migration.NewManager(migrations, idProd)
		err = migrationManager.Run()

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
				return migration.Init("lachesis-config-test", "ajIr@Quicuj9").
					NewNamed("v1-test", func() error {
						err := data.AddSection("NewSection", "")
						return err
					})
			}(data)

			idProd := NewConfigIdProducer(data)
			migrationManager := migration.NewManager(migrations, idProd)
			err = migrationManager.Run()

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
					return migration.Init("lachesis-config-test", "ajIr@Quicuj9").
						NewNamed("v1-test", func() error {
							err := data.AddSection("NewSection.NewSubsection.NewSubsection2", "")
							return err
						})
				}(data)

				idProd := NewConfigIdProducer(data)
				migrationManager := migration.NewManager(migrations, idProd)
				err = migrationManager.Run()

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
					return migration.Init("lachesis-config-test", "ajIr@Quicuj9").
						NewNamed("v1-test", func() error {
							err := data.AddSection("Node.NewSubSection.NewSubsection2", "")
							return err
						})
				}(data)

				idProd := NewConfigIdProducer(data)
				migrationManager := migration.NewManager(migrations, idProd)
				err = migrationManager.Run()

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
				return migration.Init("lachesis-config-test", "ajIr@Quicuj9").
					NewNamed("v1-test", func() error {
						err := data.DeleteSection("Node")
						return err
					})
			}(data)

			idProd := NewConfigIdProducer(data)
			migrationManager := migration.NewManager(migrations, idProd)
			err = migrationManager.Run()

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
				return migration.Init("lachesis-config-test", "ajIr@Quicuj9").
					NewNamed("v1-test", func() error {
						err := data.DeleteSection("Lachesis.Emitter.EmitIntervals")
						return err
					})
			}(data)

			idProd := NewConfigIdProducer(data)
			migrationManager := migration.NewManager(migrations, idProd)
			err = migrationManager.Run()

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
}
