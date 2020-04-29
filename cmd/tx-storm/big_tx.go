package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"gopkg.in/urfave/cli.v1"

	//ethparams "github.com/ethereum/go-ethereum/params"
	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/lachesis"
)

var (
	bigTx = cli.Command{
		Action:    sendBigTx,
		Name:      "big-transaction",
		Usage:     "Send a big fake transaction",
		ArgsUsage: " ",
		Category:  "MISCELLANEOUS COMMANDS",
		Description: `
For testing purpose.
`,
	}
)

func sendBigTx(cmd *cli.Context) error {
	glogger := log.NewGlogHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(false)))
	glogger.Verbosity(log.Lvl(cmd.GlobalInt(VerbosityFlag.Name)))
	log.Root().SetHandler(glogger)

	args := cmd.Args()
	if len(args) != 1 {
		log.Error("url expected")
		return fmt.Errorf("url expected")
	}

	url := args[0]

	from := MakeAcc(1)
	to := MakeAcc(2)

	tx := from.BigTxTo(to, 0, 31*1024)

	client, err := ethclient.Dial(url)
	if err != nil {
		log.Error("ethclient.Dial", "url", url, "err", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.SendTransaction(ctx, tx)
	if err != nil {
		log.Error("client.SendTransaction", "err", err)
		return err
	}

	return nil
}

func (a *Acc) BigTxTo(b *Acc, nonce uint, size int) *types.Transaction {
	data := make([]byte, size)
	gas, err := evmcore.IntrinsicGas(data, false, false)
	if err != nil {
		panic(err)
	}

	tx := types.NewTransaction(
		uint64(nonce),
		*b.Addr,
		big.NewInt(0),
		gas*10,
		gasPrice,
		data,
	)

	signed, err := types.SignTx(
		tx,
		types.NewEIP155Signer(big.NewInt(int64(lachesis.FakeNetworkID))),
		a.Key,
	)
	if err != nil {
		panic(err)
	}

	return signed
}
