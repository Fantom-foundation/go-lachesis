package basiccheck

import (
	"crypto/ecdsa"
	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/lachesis"
	"github.com/Fantom-foundation/go-lachesis/lachesis/params"
	"github.com/Fantom-foundation/go-lachesis/vector"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"log"
	"math"
	"math/big"
	"testing"
	"time"
)

// TestBasicCheck is a main testing func
func TestBasicCheck(t *testing.T) {
	lachesisConfigs := []*lachesis.DagConfig{
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

	for _, lachesisConfig := range lachesisConfigs {
		checker := New(lachesisConfig)
		testValidateTx(t, checker)
		testCheckTxs(t, checker)
		testCheckerValidate(t, checker)
	}
}

// testCheckerValidate runs tests for Validate func
func testCheckerValidate(t *testing.T, c *Checker) {
	testEvents := makeTestEvents()
	for _, event := range testEvents {
		if event != nil && event.GasPowerUsed == 0 && c != nil && c.config != nil {
			event.GasPowerUsed = CalcGasPowerUsed(event, c.config)
			//log.Println("event.GasPowerUsed", event.GasPowerUsed)
		}

		err := c.Validate(event)
		checkValidateTxError(t, event, err, c)
	}
}

// checkValidateTxError checks that returned error is expected
func checkValidateTxError(t *testing.T, event *inter.Event, err error, c *Checker) {
	if event == nil {
		require.Equal(t, err, ErrUninitialisedEvent)
		return
	}
	if c.config == nil {
		require.Equal(t, err, ErrUninitialisedConfig)
		return
	}
	if event.Version != 0 {
		require.Equal(t, err, ErrVersion)
		return
	}
	if len(event.Extra) > params.MaxExtraData {
		require.Equal(t, err, ErrExtraTooLarge)
		return
	}
	if len(event.Parents) > c.config.MaxParents {
		require.Equal(t, err, ErrTooManyParents)
		return
	}
	if event.Seq >= math.MaxInt32/2 || event.Epoch >= math.MaxInt32/2 || event.Frame >= math.MaxInt32/2 ||
		event.Lamport >= math.MaxInt32/2 || event.GasPowerUsed >= math.MaxInt64/2 || event.GasPowerLeft.Max() >= math.MaxInt64/2 {
		require.Equal(t, err, ErrHugeValue)
		return
	}
	if event.Seq == 0 || event.Epoch == 0 || event.Frame == 0 || event.Lamport == 0 {
		require.Equal(t, err, ErrNotInited)
		return
	}

	if event.ClaimedTime == 0 {
		require.Equal(t, err, ErrZeroTime)
		return
	}
	if event.Seq > 1 && len(event.Parents) == 0 {
		require.Equal(t, err, ErrNoParents)
		return
	}
	if len(event.Sig) != validSigLength {
		require.Equal(t, err, ErrSigMalformed)
		return
	}
	if event.GasPowerUsed > params.MaxGasPowerUsed {
		require.Equal(t, err, ErrTooBigGasUsed)
		return
	}
	if event.GasPowerUsed != CalcGasPowerUsed(event, c.config) {
		require.Equal(t, err, ErrWrongGasUsed)
		return
	}
}

// testCheckTxs runs checkTxs func of a checker and validates result
func testCheckTxs(t *testing.T, c *Checker) {
	transactions := makeTestTransactions()
	event := inter.Event{Transactions: transactions}
	err := c.checkTxs(&event)
	require.NotNil(t, err)

	validTransactions := makeValidTransactions()
	event = inter.Event{Transactions: validTransactions}
	err = c.checkTxs(&event)
	require.Nil(t, err)
}

// testValidateTx runs validateTx func of a checker and validates result
func testValidateTx(t *testing.T, c *Checker) {
	transactions := makeTestTransactions()
	for _, tx := range transactions {
		err := c.validateTx(tx)
		handleCheckTxError(t, tx, err)
	}
}

// handleCheckTxError checks that returned error is expected
func handleCheckTxError(t *testing.T, tx *types.Transaction, returnedError error) {
	if tx.Value().Sign() < 0 || tx.GasPrice().Sign() < 0 {
		require.Equal(t, returnedError, ErrNegativeValue)
		return
	}

	intrGas, err2 := evmcore.IntrinsicGas(tx.Data(), tx.To() == nil, true)
	if err2 != nil {
		require.Equal(t, returnedError, err2)
		return
	}

	if tx.Gas() < intrGas {
		require.Equal(t, returnedError, ErrIntrinsicGas)
		return
	}
	if tx.GasPrice().Cmp(params.MinGasPrice) < 0 {
		require.Equal(t, returnedError, ErrUnderpriced)
		return
	}
}

// makeValidTransactions creates set of transactions
func makeValidTransactions() []*types.Transaction {
	return []*types.Transaction{
		newTransaction(newWallet().Address, big.NewInt(1e6), 1e6, big.NewInt(1e9), []byte{0x01}),
	}
}

// makeParentsForTests creates parents for an event
func makeParentsForTests(num int) hash.Events {
	var hashEvents hash.Events
	var h common.Hash
	for i := num; i > 0; i-- {
		hashEvents = append(hashEvents, hash.Event(h))
	}
	return hashEvents
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
		makeParentsForTests(0),
		makeParentsForTests(1),
		makeParentsForTests(1e3),
	}

	for _, version := range versions {
		for _, isRoot := range isRoots {
			for _, seq := range seqs {
				for _, extra := range extras {
					for _, parents := range parentss {
						for _, epoch := range epochs {
							for _, frame := range frames {
								for _, lamport := range lamports {
									for _, claimedTime := range claimedTimes {
										data := inter.EventHeaderData{
											Version:     version,
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

	return eventHeadersDatas
}

// makeTestEvents creates a set of events for a testing purposes
func makeTestEvents() []*inter.Event {
	var events []*inter.Event
	var eventHeaders []inter.EventHeader
	transactionSets := [][]*types.Transaction{
		makeValidTransactions(),
		makeTestTransactions(),
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

// makeTestTransactions creates test transactions
func makeTestTransactions() []*types.Transaction {
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
					tx := newTransaction(newWallet().Address, amount, gasLimit, gasPrice, data)
					transactions = append(transactions, tx)
				}
			}
		}
	}
	return transactions
}

var nonce uint64 = 0

// newTransaction creates new tx for tests
func newTransaction(address common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *types.Transaction {
	defer func() { nonce += 1 }()
	tx := types.NewTransaction(nonce, address, amount, gasLimit, gasPrice, data)
	return tx
}

// newWallet creates new test wallet for testing purposes
func newWallet() TestWallet {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return TestWallet{
		PubKey:  publicKeyBytes,
		Address: address,
	}
}

// TestWallet is a special struct for tests
type TestWallet struct {
	PubKey  []byte
	Address common.Address
}
