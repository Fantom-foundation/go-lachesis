package main

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/urfave/cli.v1"

	"github.com/Fantom-foundation/go-lachesis/version"
)

func TestVersionCheck(t *testing.T) {
	t.Run("Actual node version", func(t *testing.T) {
		ctx := cli.NewContext(cli.NewApp(), flag.NewFlagSet("test", flag.ContinueOnError), nil)
		cfg := makeAllConfigs(ctx)
		applyVersionCheck(&cfg, version.Actual, "", nil)
		require.False(t, cfg.Lachesis.DisablePrivateAccountAPI)
	})

	t.Run("Nightly node version", func(t *testing.T) {
		ctx := cli.NewContext(cli.NewApp(), flag.NewFlagSet("test", flag.ContinueOnError), nil)
		cfg := makeAllConfigs(ctx)
		applyVersionCheck(&cfg, version.Nightly, "", nil)
		require.False(t, cfg.Lachesis.DisablePrivateAccountAPI)
	})

	t.Run("Outdated node version", func(t *testing.T) {
		ctx := cli.NewContext(cli.NewApp(), flag.NewFlagSet("test", flag.ContinueOnError), nil)
		cfg := makeAllConfigs(ctx)
		applyVersionCheck(&cfg, version.Outdated, "", nil)
		require.True(t, cfg.Lachesis.DisablePrivateAccountAPI)
	})
}
