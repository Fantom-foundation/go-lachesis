package main

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestConfigMigrations(t *testing.T) {
	t.Run("Test parse fixture config", func(t *testing.T){
		source := filepath.Join("testdata", "test_config.toml")

		cfg := config{}
		err := loadAllConfigs(source, &cfg)
		assert.NoError(t, err, "Parse fixture config without error")
	})
}
