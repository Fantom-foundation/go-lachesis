package main

import (
	"net/url"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/log"
	"gopkg.in/urfave/cli.v1"

	vers "github.com/Fantom-foundation/go-lachesis/version"
)

var noCheckVersionFlag = cli.BoolFlag{
	Name:  "nocheckversion",
	Usage: "Disable node version check on the latest",
}

func checkNodeVersion(uri *url.URL, ctx *cli.Context, cfg *config, version string) {
	noCheckVersion := ctx.GlobalBool(noCheckVersionFlag.Name)
	if noCheckVersion {
		return
	}
	resp, err := vers.CheckNodeVersion(uri, version)
	if err != nil {
		utils.Fatalf("Failed node version: %v", err)
	}
	if resp.Message == "" { // version is latest
		return
	}
	if resp.IsNightlyBuild { // version is nightly build
		log.Warn(resp.Message)
		return
	}
	// exit if node is validator and version is not nightly build
	validator := ctx.GlobalString(validatorFlag.Name)
	if validator != "no" && validator != "0" && validator != "" {
		utils.Fatalf("Failed node version for validator: %v", resp.Message)
	}
	cfg.Lachesis.DisablePrivateAccountAPI = true // this node is not latest, disable private account api
	log.Warn(resp.Message)
}
