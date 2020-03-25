package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Fantom-foundation/go-lachesis/utils/migration"
	"github.com/Fantom-foundation/go-lachesis/utils/toml"
)

func TestConfigParse(t *testing.T) {

	t.Run("Test parse fixture config with version", func(t *testing.T) {
		require := require.New(t)

		source := filepath.Join("testdata", "test_config.toml")
		modified := filepath.Join("testdata", "test_config_modified.toml")
		err := copyFile(source, modified)
		require.NoError(err, "Copy error")

		cfg := config{}
		err = cfg.Load(modified)
		require.Error(err, "Error when load version without migrations")

		os.Remove(modified)
		os.Remove(modified + ".init")
	})

	t.Run("Test parse fixture config without version", func(t *testing.T) {
		require := require.New(t)

		source := filepath.Join("testdata", "test_config_wo_version.toml")
		modified := filepath.Join("testdata", "test_config_wo_version_modified.toml")
		err := copyFile(source, modified)
		require.NoError(err, "Copy error")

		cfg := config{}
		err = cfg.Load(modified)
		require.NoError(err, "Parse fixture config without error")

		os.Remove(modified)
		os.Remove(modified + ".init")
	})
}

func TestConfigMigrations(t *testing.T) {

	t.Run("Rename section", func(t *testing.T) {
		require := require.New(t)

		data := execMigrations(t, "test_config_wo_version.toml", func(data *toml.Helper) *migration.Migration {
			return migration.Begin("lachesis-config-test").
				Next("v.2.0", func() error {
					err := data.RenameSection("Node", "NodeRenamed")
					return err
				})
		})

		newVersion, err := data.GetParamString("Version", "")
		require.NoError(err, "Get version from config data after config migration")
		require.Equal("v.2.0", newVersion, "Version right after config migrations")

		p, err := data.GetParamString("IPCPath", "NodeRenamed")
		require.NoError(err, "Get param from renamed section after config migration")
		require.Equal("lachesis.ipc", p, "Param right from renamed section after config migrations")

		p, err = data.GetParamString("ListenAddr", "NodeRenamed.P2P")
		require.NoError(err, "Get param from sub renamed section after config migration")
		require.Equal(":7946", p, "Param right from sub renamed section after config migrations")
	})

	t.Run("Add section", func(t *testing.T) {

		t.Run("in root", func(t *testing.T) {
			require := require.New(t)
			data := execMigrations(t, "test_config_wo_version.toml", func(data *toml.Helper) *migration.Migration {
				return migration.Begin("lachesis-config-test").
					Next("v.2.0", func() error {
						err := data.AddSection("NewSection", "")
						return err
					})
			})

			newVersion, err := data.GetParamString("Version", "")
			require.NoError(err, "Get version from config data after config migration")
			require.Equal("v.2.0", newVersion, "Version right after config migrations")

			sec, err := data.FindSection("NewSection")
			require.NoError(err, "New section search no error after add section in config migrations")
			require.NotNil(sec, "New section exists after add section in config migrations")
		})

		t.Run("sub section", func(t *testing.T) {

			t.Run("parent not exists", func(t *testing.T) {
				require := require.New(t)
				data := execMigrations(t, "test_config_wo_version.toml", func(data *toml.Helper) *migration.Migration {
					return migration.Begin("lachesis-config-test").
						Next("v.2.0", func() error {
							err := data.AddSection("NewSection.NewSubsection.NewSubsection2", "")
							return err
						})
				})

				newVersion, err := data.GetParamString("Version", "")
				require.NoError(err, "Get version from config data after config migration")
				require.Equal("v.2.0", newVersion, "Version right after config migrations")

				sec, err := data.FindSection("NewSection.NewSubsection.NewSubsection2")
				require.NoError(err, "New section search no error after add section in config migrations")
				require.NotNil(sec, "New section exists after add section in config migrations")
			})

			t.Run("parent exists", func(t *testing.T) {
				require := require.New(t)
				data := execMigrations(t, "test_config_wo_version.toml", func(data *toml.Helper) *migration.Migration {
					return migration.Begin("lachesis-config-test").
						Next("v.2.0", func() error {
							err := data.AddSection("Node.NewSubSection.NewSubsection2", "")
							return err
						})
				})

				newVersion, err := data.GetParamString("Version", "")
				require.NoError(err, "Get version from config data after config migration")
				require.Equal("v.2.0", newVersion, "Version right after config migrations")

				sec, err := data.FindSection("Node.NewSubSection.NewSubsection2")
				require.NoError(err, "New section search no error after add section in config migrations")
				require.NotNil(sec, "New section exists after add section in config migrations")
			})
		})
	})

	t.Run("Delete section", func(t *testing.T) {

		t.Run("in root", func(t *testing.T) {
			require := require.New(t)
			data := execMigrations(t, "test_config_wo_version.toml", func(data *toml.Helper) *migration.Migration {
				return migration.Begin("lachesis-config-test").
					Next("v.2.0", func() error {
						err := data.DeleteSection("Node")
						return err
					})
			})

			newVersion, err := data.GetParamString("Version", "")
			require.NoError(err, "Get version from config data after config migration")
			require.Equal("v.2.0", newVersion, "Version right after config migrations")

			sec, err := data.FindSection("Node")
			require.Error(err, "Deleted section search should return error after delete section in config migrations")
			require.Nil(sec, "Deleted section not exists after delete section in config migrations")
		})

		t.Run("in sub section", func(t *testing.T) {
			require := require.New(t)
			data := execMigrations(t, "test_config_wo_version.toml", func(data *toml.Helper) *migration.Migration {
				return migration.Begin("lachesis-config-test").
					Next("v.2.0", func() error {
						err := data.DeleteSection("Lachesis.Emitter.EmitIntervals")
						return err
					})
			})

			newVersion, err := data.GetParamString("Version", "")
			require.NoError(err, "Get version from config data after config migration")
			require.Equal("v.2.0", newVersion, "Version right after config migrations")

			sec, err := data.FindSection("Lachesis.Emitter.EmitIntervals")
			require.Error(err, "Deleted section search should return error after delete section in config migrations")
			require.Nil(sec, "Deleted section not exists after delete section in config migrations")

			// Parent section should be exists after delete child
			sec, err = data.FindSection("Lachesis.Emitter")
			require.NoError(err, "Parent section search no error after delete section in config migrations")
			require.NotNil(sec, "Parent section exists after delete section in config migrations")
		})
	})

	t.Run("Rename param", func(t *testing.T) {
		require := require.New(t)
		data := execMigrations(t, "test_config_wo_version.toml", func(data *toml.Helper) *migration.Migration {
			return migration.Begin("lachesis-config-test").
				Next("v.2.0", func() error {
					err := data.RenameParam("Validator", "Lachesis.Emitter", "ValidatorAddr")
					return err
				})
		})

		newVersion, err := data.GetParamString("Version", "")
		require.NoError(err, "Get version from config data after config migration")
		require.Equal("v.2.0", newVersion, "Version right after config migrations")

		paramStr, err := data.GetParamString("ValidatorAddr", "Lachesis.Emitter")
		require.NoError(err, "Get param data no error after rename param in config migrations")
		require.NotEmpty(paramStr, "Param data not empty after rename param in config migrations")
	})

	t.Run("Add param", func(t *testing.T) {

		t.Run("to root", func(t *testing.T) {
			require := require.New(t)
			data := execMigrations(t, "test_config_wo_version.toml", func(data *toml.Helper) *migration.Migration {
				return migration.Begin("lachesis-config-test").
					Next("v.2.0", func() error {
						err := data.AddParam("NewParamInRoot", "", "string value")
						return err
					})
			})

			newVersion, err := data.GetParamString("Version", "")
			require.NoError(err, "Get version from config data after config migration")
			require.Equal("v.2.0", newVersion, "Version right after config migrations")

			paramStr, err := data.GetParamString("NewParamInRoot", "")
			require.NoError(err, "Get param data no error after add param to root in config migrations")
			require.NotEmpty(paramStr, "Param data not empty after add param to root in config migrations")
		})

		t.Run("to section", func(t *testing.T) {
			require := require.New(t)
			data := execMigrations(t, "test_config_wo_version.toml", func(data *toml.Helper) *migration.Migration {
				return migration.Begin("lachesis-config-test").
					Next("v.2.0", func() error {
						err := data.AddParam("NewParamInSection", "Lachesis.Emitter", "string value")
						return err
					})
			})

			newVersion, err := data.GetParamString("Version", "")
			require.NoError(err, "Get version from config data after config migration")
			require.Equal("v.2.0", newVersion, "Version right after config migrations")

			paramStr, err := data.GetParamString("NewParamInSection", "Lachesis.Emitter")
			require.NoError(err, "Get param data no error after add param in section in config migrations")
			require.NotEmpty(paramStr, "Param data not empty after add param in section in config migrations")
		})
	})

	t.Run("Delete param", func(t *testing.T) {
		require := require.New(t)
		data := execMigrations(t, "test_config_wo_version.toml", func(data *toml.Helper) *migration.Migration {
			return migration.Begin("lachesis-config-test").
				Next("v.2.0", func() error {
					err := data.DeleteParam("Validator", "Lachesis.Emitter")
					return err
				})
		})

		newVersion, err := data.GetParamString("Version", "")
		require.NoError(err, "Get version from config data after config migration")
		require.Equal("v.2.0", newVersion, "Version right after config migrations")

		paramStr, err := data.GetParamString("Validator", "Lachesis.Emitter")
		require.Error(err, "Get param data return error after delete param in config migrations")
		require.Empty(paramStr, "Param data is empty after delete param in config migrations")
	})

	t.Run("Set param value", func(t *testing.T) {

		t.Run("String", func(t *testing.T) {
			require := require.New(t)
			data := execMigrations(t, "test_config_wo_version.toml", func(data *toml.Helper) *migration.Migration {
				return migration.Begin("lachesis-config-test").
					Next("v.2.0", func() error {
						err := data.SetParam("Validator", "Lachesis.Emitter", "test value")
						return err
					})
			})

			newVersion, err := data.GetParamString("Version", "")
			require.NoError(err, "Get version from config data after config migration")
			require.Equal("v.2.0", newVersion, "Version right after config migrations")

			paramStr, err := data.GetParamString("Validator", "Lachesis.Emitter")
			require.NoError(err, "Get param data return no error after set param in config migrations")
			require.Equal("test value", paramStr, "Param data correct after set param in config migrations")
		})

		t.Run("Int", func(t *testing.T) {
			require := require.New(t)
			data := execMigrations(t, "test_config_wo_version.toml", func(data *toml.Helper) *migration.Migration {
				return migration.Begin("lachesis-config-test").
					Next("v.2.0", func() error {
						err := data.SetParam("MaxTxsFromSender", "Lachesis.Emitter", 1112233)
						return err
					})
			})

			newVersion, err := data.GetParamString("Version", "")
			require.NoError(err, "Get version from config data after config migration")
			require.Equal("v.2.0", newVersion, "Version right after config migrations")

			paramInt, err := data.GetParamInt("MaxTxsFromSender", "Lachesis.Emitter")
			require.NoError(err, "Get param data return no error after set int param in config migrations")
			require.Equal(int64(1112233), paramInt, "Param data correct after set int param in config migrations")
		})

		t.Run("Float", func(t *testing.T) {
			require := require.New(t)
			data := execMigrations(t, "test_config_wo_version.toml", func(data *toml.Helper) *migration.Migration {
				return migration.Begin("lachesis-config-test").
					Next("v.2.0", func() error {
						err := data.SetParam("MaxGasRateGrowthFactor", "Lachesis.Emitter", 1.234)
						return err
					})
			})

			newVersion, err := data.GetParamString("Version", "")
			require.NoError(err, "Get version from config data after config migration")
			require.Equal("v.2.0", newVersion, "Version right after config migrations")

			paramFloat, err := data.GetParamFloat("MaxGasRateGrowthFactor", "Lachesis.Emitter")
			require.NoError(err, "Get param data return no error after set float param in config migrations")
			require.Equal(1.234, paramFloat, "Param data correct after set float param in config migrations")
		})

		t.Run("Bool", func(t *testing.T) {
			require := require.New(t)
			data := execMigrations(t, "test_config_wo_version.toml", func(data *toml.Helper) *migration.Migration {
				return migration.Begin("lachesis-config-test").
					Next("v.2.0", func() error {
						err := data.SetParam("TxIndex", "Lachesis", false)
						return err
					})
			})

			newVersion, err := data.GetParamString("Version", "")
			require.NoError(err, "Get version from config data after config migration")
			require.Equal("v.2.0", newVersion, "Version right after config migrations")

			paramBool, err := data.GetParamBool("TxIndex", "Lachesis")
			require.NoError(err, "Get param data return no error after set bool param in config migrations")
			require.False(paramBool, "Param data correct after set bool param in config migrations")
		})

		t.Run("DateTime", func(t *testing.T) {
			require := require.New(t)

			testTime0 := time.Now().UTC().Add(time.Hour)
			testTime := time.Now().UTC()

			data := execMigrations(t, "test_config_wo_version.toml", func(data *toml.Helper) *migration.Migration {
				return migration.Begin("lachesis-config-test").
					Next("v.2.0", func() error {
						err := data.AddParam("TestTime", "Node", testTime0)
						if err != nil {
							return err
						}
						err = data.SetParam("TestTime", "Node", testTime)
						return err
					})
			})

			newVersion, err := data.GetParamString("Version", "")
			require.NoError(err, "Get version from config data after config migration")
			require.Equal("v.2.0", newVersion, "Version right after config migrations")

			paramTime, err := data.GetParamTime("TestTime", "Node")
			require.NoError(err, "Get param data return no error after set datetime param in config migrations")
			require.Equal(testTime, paramTime, "Param data correct after set datetime param in config migrations")
		})
	})
}

func execMigrations(t *testing.T, file string, migrationsFunc func(helper *toml.Helper) *migration.Migration) *toml.Helper {
	source := filepath.Join("testdata", file)
	table, err := toml.ParseFile(source)
	require.NoError(t, err, "Parse config to Table")
	helper := toml.NewTomlHelper(table)

	migrations := migrationsFunc(helper)
	idProd := toml.NewIDStore(helper, migrations.IdChain())
	err = migrations.Exec(idProd)
	require.NoError(t, err, "Exec migrations")

	return helper
}
