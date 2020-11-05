package main

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/Fantom-foundation/go-lachesis/crypto"
	"github.com/Fantom-foundation/go-lachesis/lachesis"
)

var (
	gasLimit = uint64(21000)
	gasPrice = lachesis.FakeEconomyConfig().InitialMinGasPrice // minimal
)

type Acc struct {
	Key  *ecdsa.PrivateKey
	Addr *common.Address
}

func MakeAcc(n uint) *Acc {
	key := crypto.FakeKey(int(n))
	addr := crypto.PubkeyToAddress(key.PublicKey)

	return &Acc{
		Key:  key,
		Addr: &addr,
	}
}

func (a *Acc) TransactionTo(b *Acc, nonce uint, amount *big.Int, chainId uint) *types.Transaction {
	tx := types.NewTransaction(
		uint64(nonce),
		*b.Addr,
		amount,
		gasLimit,
		gasPrice,
		[]byte{},
	)

	signed, err := types.SignTx(
		tx,
		types.NewEIP155Signer(big.NewInt(int64(chainId))),
		a.Key,
	)
	if err != nil {
		panic(err)
	}

	return signed
}
