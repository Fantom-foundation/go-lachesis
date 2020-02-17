package main

import (
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

}
