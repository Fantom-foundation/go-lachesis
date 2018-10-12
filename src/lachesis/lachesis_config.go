package lachesis

import (
	"crypto/ecdsa"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"time"

	"github.com/andrecronje/lachesis/src/node"
	"github.com/andrecronje/lachesis/src/proxy"
	sproxy "github.com/andrecronje/lachesis/src/proxy/socket/app"
	"github.com/sirupsen/logrus"
)

type LachesisConfig struct {
	NodeConfig node.Config `mapstructure:",squash"`
	DataDir     string `mapstructure:"datadir"`
	BindAddr    string `mapstructure:"listen"`
	ServiceAddr string `mapstructure:"service-listen"`
	MaxPool     int    `mapstructure:"max-pool"`
	Store       bool   `mapstructure:"store"`
	LogLevel    string `mapstructure:"log"`
	LoadPeers bool
	Proxy     proxy.AppProxy
	Key       *ecdsa.PrivateKey
	Logger    *logrus.Logger

	Test  bool
	TestN uint64
}

func NewDefaultConfig() *LachesisConfig {
	config := &LachesisConfig{
		DataDir:    DefaultDataDir(),
		BindAddr:   ":1337",
		Proxy:      nil,
		Logger:     logrus.New(),
		MaxPool:    2,
		NodeConfig: *node.DefaultConfig(),
		Store:      false,
		LoadPeers:  true,
		Key:        nil,
		Test:        false,
		TestN:       ^uint64(0),
	}

	config.NodeConfig.Logger = config.Logger

	//XXX
	config.Proxy, _ = sproxy.NewSocketAppProxy("127.0.0.1:1338", "127.0.0.1:1339", 1*time.Second, config.Logger)

	return config
}

func (c *LachesisConfig) BadgerDir() string {
	return filepath.Join(c.DataDir, "badger_db")
}

func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	home := HomeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, ".lachesis")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "LACHESIS")
		} else {
			return filepath.Join(home, ".lachesis")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

func HomeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

func LogLevel(l string) logrus.Level {
	switch l {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.DebugLevel
	}
}
