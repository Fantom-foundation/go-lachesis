package testCommon

import (
	"crypto/ecdsa"
	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/lachesis"
	"github.com/Fantom-foundation/go-lachesis/lachesis/params"
	"github.com/Fantom-foundation/go-lachesis/vector"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math"
	"math/big"
	"time"
)

const validSigLength = 65

// MakeTestEvents creates a set of testing events
func MakeEventList() []*inter.Event {
	var events []*inter.Event
	var eventHeaders []inter.EventHeader
	transactionSets := [][]*types.Transaction{
		MakeValidTransactions(),
		MakeTestTransactions(),
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
	var eventHeadersDatas []inter.EventHeaderData
	versions := []uint32{
		0, 1,
	}
	isRoots := []bool{
		true, false,
	}
	seqs := []idx.Event{
		0, 1, math.MaxInt32 / 2, (math.MaxInt32 / 2) + 1,
	}
	epochs := []idx.Epoch{
		0, 1, math.MaxInt32 / 2, (math.MaxInt32 / 2) + 1,
	}
	frames := []idx.Frame{
		0, 1, math.MaxInt32 / 2, (math.MaxInt32 / 2) + 1,
	}
	lamports := []idx.Lamport{
		0, 1, math.MaxInt32 / 2, (math.MaxInt32 / 2) + 1,
	}
	extras := [][]byte{
		[]byte{}, nil, makeDataWithLen(10), makeDataWithLen(params.MaxExtraData), makeDataWithLen(params.MaxExtraData + 1),
	}
	claimedTimes := []inter.Timestamp{
		0, 1, inter.Timestamp(uint64(time.Now().Unix())),
	}
	parentss := []hash.Events{
		MakeParentsForTests(0),
		MakeParentsForTests(1),
		MakeParentsForTests(1e3),
	}
	creators := []idx.StakerID{1, 102}

	for _, version := range versions {
		for _, isRoot := range isRoots {
			for _, seq := range seqs {
				for _, extra := range extras {
					for _, parents := range parentss {
						for _, epoch := range epochs {
							for _, frame := range frames {
								for _, lamport := range lamports {
									for _, claimedTime := range claimedTimes {
										for _, creator := range creators {
											data := inter.EventHeaderData{
												Version:     version,
												Creator:     creator,
												IsRoot:      isRoot,
												Seq:         seq,
												Extra:       extra,
												Parents:     parents,
												Epoch:       epoch,
												Frame:       frame,
												Lamport:     lamport,
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
			}
		}
	}

	return eventHeadersDatas
}

// MakeParentsForTests creates parents for an event
func MakeParentsForTests(num int) hash.Events {
	var hashEvents hash.Events
	var h common.Hash
	for i := num; i > 0; i-- {
		hashEvents = append(hashEvents, hash.Event(h))
	}
	return hashEvents
}

// MakeTestTransactions creates test transactions
func MakeTestTransactions() []*types.Transaction {
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
					tx := NewTransaction(NewTestWallet().Address, amount, gasLimit, gasPrice, data)
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
	defer func() { nonce += 1 }()
	tx := types.NewTransaction(nonce, address, amount, gasLimit, gasPrice, data)
	return tx
}

// MakeValidTransactions creates set of transactions
func MakeValidTransactions() []*types.Transaction {
	return []*types.Transaction{
		NewTransaction(NewTestWallet().Address, big.NewInt(1e6), 1e6, big.NewInt(1e9), []byte{0x01}),
	}
}

func MakeTestConfigs() []*lachesis.DagConfig {
	return []*lachesis.DagConfig{
		nil,
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

// TestWallet is just a wallet for tests
type TestWallet struct {
	Address    common.Address
	PubKey     []byte
	PrivateKey ecdsa.PrivateKey
}

// newTestWallet creates test wallet
func NewTestWallet() TestWallet {
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
	return TestWallet{
		PrivateKey: *privateKey,
		PubKey:     publicKeyBytes,
		Address:    address,
	}
}
