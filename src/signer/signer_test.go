package signer

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/signer/core"

	"github.com/Fantom-foundation/go-lachesis/src/utils"
)

const pass = "password_with_more_than_10_chars"

var tmpDir = utils.NewTempDir("lachesis-config")

func TestSignerAPI_New(t *testing.T) {
	// Init new signer api & ui handler
	manager := NewSignerTestManager(tmpDir)

	address, err := manager.NewAccount(pass)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(address.Hex())
}

func TestSignerAPI_List(t *testing.T) {
	// Init new signer api & ui handler
	manager := NewSignerTestManager(tmpDir)

	// create new account
	_, err := manager.NewAccount(pass)
	if err != nil {
		t.Fatal(err)
	}

	// Some time to allow changes to propagate
	time.Sleep(250 * time.Millisecond)

	// get list account
	addresses, err := manager.ListAccounts()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(addresses[0].Hex())
}

func TestSignerAPI_SignTransaction(t *testing.T) {
	// Init new signer api & ui handler
	manager := NewSignerTestManager(tmpDir)

	// create new account
	_, err := manager.NewAccount(pass)
	if err != nil {
		t.Fatal(err)
	}

	// Some time to allow changes to propagate
	time.Sleep(250 * time.Millisecond)

	// get list account
	addresses, err := manager.ListAccounts()
	if err != nil {
		t.Fatal(err)
	}

	// make transaction
	tx := testTx(common.NewMixedcaseAddress(addresses[0]))

	// sign transaction
	err = manager.SignTransaction(tx, pass)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSignerAPI_SignData(t *testing.T) {
	// Init new signer api & ui handler
	manager := NewSignerTestManager(tmpDir)

	// create new account
	_, err := manager.NewAccount(pass)
	if err != nil {
		t.Fatal(err)
	}

	// Some time to allow changes to propagate
	time.Sleep(250 * time.Millisecond)

	// get list account
	addresses, err := manager.ListAccounts()
	if err != nil {
		t.Fatal(err)
	}

	// sign hash
	_, err = manager.SignData(addresses[0], []byte("test hash"), pass)
	if err != nil {
		t.Fatal(err)
	}
}

// TODO: Do we really need it?
func TestSignerAPI_SignTypedData(t *testing.T) {
	// Init new signer api & ui handler
	manager := NewSignerTestManager(tmpDir)

	// create new account
	manager.ui.inputCh <- "password_with_more_than_10_chars"
	_, err := manager.signer.New(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// Some time to allow changes to propagate
	time.Sleep(250 * time.Millisecond)

	// get list account
	addresses, err := manager.signer.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// make mixed address creator
	mixedAddress := common.NewMixedcaseAddress(addresses[0])

	// sign typed data
	manager.ui.inputCh <- "password_with_more_than_10_chars"
	signature, err := manager.signer.SignTypedData(context.Background(), mixedAddress, typedData)
	if err != nil {
		t.Fatal(err)
	}
	if signature == nil || len(signature) != 65 {
		t.Errorf("Expected 65 byte signature (got %d bytes)", len(signature))
	}

	t.Log(signature.String())
}

func testTx(from common.MixedcaseAddress) core.SendTxArgs {
	to := common.NewMixedcaseAddress(common.HexToAddress("0x1337"))
	gas := hexutil.Uint64(21000)
	gasPrice := (hexutil.Big)(*big.NewInt(2000000000))
	value := (hexutil.Big)(*big.NewInt(1e18))
	nonce := (hexutil.Uint64)(0)
	data := hexutil.Bytes(common.Hex2Bytes("01020304050607080a"))
	tx := core.SendTxArgs{
		From:     from,
		To:       &to,
		Gas:      gas,
		GasPrice: gasPrice,
		Value:    value,
		Data:     &data,
		Nonce:    nonce}
	return tx
}
