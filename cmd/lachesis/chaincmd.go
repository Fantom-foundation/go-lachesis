package main

import (
	"github.com/ethereum/go-ethereum/cmd/utils"
	"gopkg.in/urfave/cli.v1"
)

var (
	exportCommand = cli.Command{
		Action:    utils.MigrateFlags(exportChain),
		Name:      "export",
		Usage:     "Export blockchain into file",
		ArgsUsage: "<filename> [<epochFrom> [<epochTo>]]",
		Category:  "MISCELLANEOUS COMMANDS",
		Description: `
Requires a first argument of the file to write to. The file will be appended
if already existing. If the file ends with .gz, the output will
be gzipped.
Optional second and third arguments control the first and last epoch to export.`,
	}

	importCommand = cli.Command{
		Action:    utils.MigrateFlags(importChain),
		Name:      "import",
		Usage:     "Import a blockchain file",
		ArgsUsage: "<filename> (<filename 2> ... <filename N>) ",
		Category:  "MISCELLANEOUS COMMANDS",
		Description: `
The import command imports DAG-events from an RLP-encoded form.
If only one file is used, import error will result in failure.
If several files are used, processing will proceed even
if an individual RLP-file import failure occurs.`,
	}
)
