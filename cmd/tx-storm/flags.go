package main

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/urfave/cli.v1"
)

var TxnsRateFlag = cli.IntFlag{
	Name:  "rate",
	Usage: "transactions per second (max sum of all instances)",
}

func getTxnsRate(ctx *cli.Context) uint {
	return uint(ctx.GlobalInt(TxnsRateFlag.Name))
}

var NumberFlag = cli.StringFlag{
	Name:  "num",
	Usage: "'N/X' - it is a N-th generator of X",
	Value: "1/1",
}

func getNumber(ctx *cli.Context) (num, total uint) {
	var err error
	num, total, err = parseNumber(ctx.GlobalString(NumberFlag.Name))
	if err != nil {
		panic(err)
	}
	return
}

func parseNumber(s string) (num, total uint, err error) {
	var i64 uint64

	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		err = fmt.Errorf("use %%d/%%d format")
		return
	}

	i64, err = strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return
	}
	num = uint(i64)

	i64, err = strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return
	}
	total = uint(i64)

	if 1 > num || num > total {
		err = fmt.Errorf("key-num should be in range from 1 to total : <key-num>/<total>")
	}

	return
}

var VerbosityFlag = cli.IntFlag{
	Name:  "verbosity",
	Usage: "sets the verbosity level",
	Value: 3,
}
