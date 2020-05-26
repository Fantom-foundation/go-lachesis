package main

import (
	"errors"

	"gopkg.in/urfave/cli.v1"

	"github.com/Fantom-foundation/go-lachesis/utils/toml"
)

var ConfigFileFlag = cli.StringFlag{
	Name:  "config",
	Usage: "TOML configuration file",
	Value: "config.toml",
}

type Config struct {
	ChainId uint     // chain id for sign transactions
	URLs    []string // WS nodes API URL
	Accs    struct {
		Count  uint // count of predefined fake accounts
		Offset uint // offset of predefined fake accounts
	}
}

func openConfig(ctx *cli.Context) *Config {
	cfg := &Config{}
	f := ctx.GlobalString(ConfigFileFlag.Name)
	err := cfg.Load(f)
	if err != nil {
		panic(err)
	}
	return cfg
}

func (cfg *Config) Load(file string) error {
	data, err := toml.ParseFile(file)
	if err != nil {
		return err
	}

	err = toml.Settings.UnmarshalTable(data, cfg)
	if err != nil {
		err = errors.New(file + ", " + err.Error())
		return err
	}

	return nil
}
