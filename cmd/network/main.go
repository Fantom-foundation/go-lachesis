package main

import (
	_ "net/http/pprof"
	"os"

	cmd "github.com/Fantom-foundation/go-lachesis/cmd/network/commands"
)

func main() {
	rootCmd := cmd.RootCmd

	rootCmd.AddCommand(
		cmd.VersionCmd,
		cmd.NewProxyCmd(),
		cmd.NewRunCmd())

	//Do not print usage when error occurs
	rootCmd.SilenceUsage = true

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
