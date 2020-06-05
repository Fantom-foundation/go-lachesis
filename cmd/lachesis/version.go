package main

import (
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"gopkg.in/urfave/cli.v1"

	"github.com/Fantom-foundation/go-lachesis/version"
)

var noCheckVersionFlag = cli.BoolFlag{
	Name:  "nocheckversion",
	Usage: "Disable node version check on the latest",
}

func applyVersionCheck(
	cfg *config,
	status version.BuildStatus,
	msg string,
	err error,
) {
	hint := "Use --" + noCheckVersionFlag.Name + " flag to disable checking."

	if err != nil {
		utils.Fatalf("Failed version check: %v. %s", err, hint)
	}

	if status == version.Actual {
		return
	}

	if status == version.Nightly {
		log.Warn(msg, "hint", hint)
		return
	}

	if cfg.Lachesis.Emitter.Validator != (common.Address{}) {
		log.Error(msg, "hint", hint)
		utils.Fatalf("Emission on outdated node is not allowed!")
	}

	cfg.Lachesis.DisablePrivateAccountAPI = true
	log.Warn(msg, "hint", hint)
	log.Warn("PrivateAccountAPI is disabled on outdated node")
}
