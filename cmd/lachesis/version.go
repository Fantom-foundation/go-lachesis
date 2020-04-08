package main

import (
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/log"
	"gopkg.in/urfave/cli.v1"

	vers "github.com/Fantom-foundation/go-lachesis/version"
)

var noCheckVersionFlag = cli.BoolFlag{
	Name:  "nocheckversion",
	Usage: "Disable node version check on the latest",
}

func checkNodeVersion(ctx *cli.Context) {
	noCheckVersion := ctx.GlobalBool(noCheckVersionFlag.Name)
	if noCheckVersion {
		return
	}
	if err := vers.CheckNodeVersion(nil, App.Version); err != nil {
		if err.Error() == vers.FailedGetNodeVersionMsg {
			utils.Fatalf("Failed node version: %v", err)
		}
		if validator := ctx.GlobalString(validatorFlag.Name); validator != "no" && validator != "0" {
			utils.Fatalf("Failed node version: %v", err)
		}
		// TODO disable api for send transactions
		log.Warn(err.Error())
	}
}
