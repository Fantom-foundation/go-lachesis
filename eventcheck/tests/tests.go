package tests

import (
	"crypto/ecdsa"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/lachesis"
	"github.com/Fantom-foundation/go-lachesis/lachesis/params"
	"github.com/Fantom-foundation/go-lachesis/vector"
)

const validSigLength = 65

// Events creates a set of testing events
func Events() []*inter.Event {
	var events []*inter.Event
	var eventHeaders []inter.EventHeader
	transactionSets := [][]*types.Transaction{
		ValidTransactions(),
		Transactions(),
	}
	eventHeaderDatas := makeEventHeaderDatas()
	eventSigs := [][]byte{[]byte{}, nil, makeDataWithLen(validSigLength)}

	for _, eventHeaderData := range eventHeaderDatas {
		for _, eventSig := range eventSigs {
			header := inter.EventHeader{
				EventHeaderData: eventHeaderData,
				Sig:             eventSig,
			}
			eventHeaders = append(eventHeaders, header)
		}
	}

	for _, transactionSet := range transactionSets {
		for _, eventHeader := range eventHeaders {
			event := &inter.Event{
				EventHeader:  eventHeader,
				Transactions: transactionSet,
			}
			events = append(events, event)
		}
	}

	events = append(events, nil)
	return events
}

// makeDataWithLen creates array of bytes
func makeDataWithLen(len int) []byte {
	var data []byte
	for i := len; i > 0; i-- {
		data = append(data, 0x00)
	}
	return data
}

// makeEventHeaderDatas create test data with event headers
func makeEventHeaderDatas() []inter.EventHeaderData {

	var int32cases = []int32{
		0, 1, math.MaxInt32 / 2, (math.MaxInt32 / 2) + 1,
	}

	var eventHeadersDatas []inter.EventHeaderData
	versions := []uint32{
		0, 1,
	}
	isRoots := []bool{
		true, false,
	}
	extras := [][]byte{
		[]byte{}, nil, makeDataWithLen(10), makeDataWithLen(params.MaxExtraData), makeDataWithLen(params.MaxExtraData + 1),
	}
	claimedTimes := []inter.Timestamp{
		0, 1, inter.Timestamp(uint64(time.Now().Unix())),
	}
	parentss := []hash.Events{
		hash.FakeEvents(0),
		hash.FakeEvents(1),
		hash.FakeEvents(1e3),
	}
	creators := []idx.StakerID{1, 102}

	for _, version := range versions {
		for _, isRoot := range isRoots {
			for _, extra := range extras {
				for _, parents := range parentss {
					for _, x := range int32cases {
						for _, claimedTime := range claimedTimes {
							for _, creator := range creators {
								data := inter.EventHeaderData{
									Version:     version,
									Creator:     creator,
									IsRoot:      isRoot,
									Seq:         idx.Event(x),
									Extra:       extra,
									Parents:     parents,
									Epoch:       idx.Epoch(x),
									Frame:       idx.Frame(x),
									Lamport:     idx.Lamport(x),
									ClaimedTime: claimedTime,
								}
								eventHeadersDatas = append(eventHeadersDatas, data)
							}
						}
					}
				}
			}
		}
	}
	return eventHeadersDatas
}

// Transactions for tests
func Transactions() []*types.Transaction {
	gasLimits := []uint64{
		1e10, 0, 10, 1e6,
	}
	amounts := []*big.Int{
		big.NewInt(0),
		big.NewInt(1),
		big.NewInt(1e10),
		big.NewInt(-1),
		big.NewInt(1e6),
	}
	gasPrices := []*big.Int{
		big.NewInt(1e9), big.NewInt(1), big.NewInt(0), big.NewInt(-1),
	}
	datas := [][]byte{
		[]byte{0x01},
		[]byte("some transaction data"),
		nil,
		[]byte{},
	}
	var transactions []*types.Transaction
	for _, gasLimit := range gasLimits {
		for _, amount := range amounts {
			for _, gasPrice := range gasPrices {
				for _, data := range datas {
					tx := NewTransaction(NewWallet().Address, amount, gasLimit, gasPrice, data)
					transactions = append(transactions, tx)
				}
			}
		}
	}
	return transactions
}

var nonce uint64 = 0

// NewTransaction creates new tx for tests
func NewTransaction(address common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *types.Transaction {
	tx := types.NewTransaction(nonce, address, amount, gasLimit, gasPrice, data)
	nonce += 1
	return tx
}

// ValidTransactions creates set of transactions
func ValidTransactions() []*types.Transaction {
	return []*types.Transaction{
		NewTransaction(NewWallet().Address, big.NewInt(1e6), 1e6, big.NewInt(1e9), []byte{0x01}),
	}
}

// DagConfigs for tests.
func DagConfigs() []*lachesis.DagConfig {
	return []*lachesis.DagConfig{
		&lachesis.DagConfig{
			MaxParents:                0,
			MaxFreeParents:            0,
			MaxEpochBlocks:            0,
			MaxEpochDuration:          0,
			VectorClockConfig:         vector.IndexConfig{},
			MaxValidatorEventsInBlock: 0,
		},
		&lachesis.DagConfig{
			MaxParents:                1e10,
			MaxFreeParents:            1,
			MaxEpochBlocks:            20,
			MaxEpochDuration:          2000,
			VectorClockConfig:         vector.IndexConfig{},
			MaxValidatorEventsInBlock: 10,
		},
	}
}

// Wallet is just a wallet for tests
type Wallet struct {
	Address    common.Address
	PubKey     []byte
	PrivateKey ecdsa.PrivateKey
}

// NewWallet creates test wallet
func NewWallet() *Wallet {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("cannot create publicKeyECDSA")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return &Wallet{
		PrivateKey: *privateKey,
		PubKey:     publicKeyBytes,
		Address:    address,
	}
}
