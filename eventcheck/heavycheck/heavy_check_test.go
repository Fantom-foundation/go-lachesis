package heavycheck

import (
	"crypto/ecdsa"
	"github.com/Fantom-foundation/go-lachesis/eventcheck/epochcheck"
	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/pos"
	"github.com/Fantom-foundation/go-lachesis/lachesis"
	"github.com/Fantom-foundation/go-lachesis/lachesis/genesis"
	"github.com/Fantom-foundation/go-lachesis/lachesis/params"
	"github.com/Fantom-foundation/go-lachesis/vector"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"log"
	"math"
	"math/big"
	"runtime"
	"testing"
	"time"
)

type TestDagReader struct {
	epochPubKeys map[idx.StakerID]common.Address
	epoch        idx.Epoch
}

// GetEpochPubKeys returns pubkeys of a special testing dagReader
func (t *TestDagReader) GetEpochPubKeys() (map[idx.StakerID]common.Address, idx.Epoch) {
	return t.epochPubKeys, t.epoch
}

// getTestDagReaders creates dag readers for testing purposes
func getTestDagReaders(t *testing.T) []TestDagReader {
	var tdr []TestDagReader
	epochPubKeys1 := map[idx.StakerID]common.Address{}
	epochPubKeys1[1] = newTestWallet(t).Address
	epochPubKeys1[2] = newTestWallet(t).Address

	epochPubKeys2 := map[idx.StakerID]common.Address{}
	epochPubKeys2[11] = newTestWallet(t).Address
	epochPubKeys2[22] = newTestWallet(t).Address
	epochs := []idx.Epoch{0, 1, 10}
	for _, epoch := range epochs {
		tdr = append(tdr, TestDagReader{epochPubKeys1, epoch})
		tdr = append(tdr, TestDagReader{epochPubKeys2, epoch})
	}
	return tdr
}

// TestHeavyCheck main testing func
func TestHeavyCheck(t *testing.T) {
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
	dagReaders := getTestDagReaders(t)
	for _, cfg := range lachesisConfigs {
		for _, dagReader := range dagReaders {
			net := lachesis.FakeNetConfig(genesis.FakeAccounts(0, 5, big.NewInt(0), pos.StakeToBalance(1)))
			ledgerID := net.EvmChainConfig().ChainID

			testEvents := makeTestEvents(t)
			testEvent := testEvents[0]
			tw := newTestWallet(t)
			sig, err := crypto.Sign(testEvent.Hash().Bytes(), &tw.PrivateKey)
			require.Nil(t, err)
			testEvent.Sig = sig
			signedHash := crypto.Keccak256(testEvent.DataToSign())
			pk, err := crypto.SigToPub(signedHash, testEvent.Sig)
			require.Nil(t, err)
			dagReader.epochPubKeys[1] = crypto.PubkeyToAddress(*pk)

			checker := NewDefault(cfg, &dagReader, types.NewEIP155Signer(ledgerID))

			testChecker(t, checker, testEvents)
		}
	}

	net := lachesis.FakeNetConfig(genesis.FakeAccounts(0, 5, big.NewInt(0), pos.StakeToBalance(1)))
	ledgerID := net.EvmChainConfig().ChainID
	checker := NewDefault(lachesisConfigs[1], &dagReaders[0], types.NewEIP155Signer(ledgerID))
	testOverloaded(t, checker)
}

// testChecker runs set of main tests for a checker
func testChecker(t *testing.T, checker *Checker, testEvents []*inter.Event) {
	gNum := runtime.NumGoroutine()
	expectedGNum := gNum + checker.numOfThreads
	checker.Start()
	runtime.Gosched()
	require.Equal(t, expectedGNum, runtime.NumGoroutine())

	for _, event := range testEvents {
		testValidate(t, checker, event)
		testEnqueue(t, checker, inter.Events{event})
	}

	gNum = runtime.NumGoroutine()
	expectedGNum = gNum - checker.numOfThreads
	checker.Stop()
	runtime.Gosched()
	require.Equal(t, expectedGNum, runtime.NumGoroutine())
}

// testOverloaded is a small testing func for a Overloaded func
func testOverloaded(t *testing.T, checker *Checker) {
	runtime.GOMAXPROCS(1)
	checker.Start()
	log.Println("checker.tasksQ1", len(checker.tasksQ))
	require.False(t, checker.Overloaded())
	taskDatas := makeTaskData(maxQueuedTasks)
	for _, taskData := range taskDatas {
		checker.tasksQ <- taskData
	}
	log.Println("checker.tasksQ2", len(checker.tasksQ))
	log.Println("checker.tasksQ3", len(checker.tasksQ))
	require.True(t, checker.Overloaded())
}

// makeTaskData creates array of taskData objects
func makeTaskData(num int) []*TaskData {
	var taskDatas []*TaskData
	for i :=0 ; i< num; i++ {
		td := TaskData{}
		td.onValidated = func(ArbitraryTaskData) {}
		taskDatas = append(taskDatas, &td)
	}
	return taskDatas
}

type TestArbitraryTaskData struct {
}

// GetEvents is just an implementation
func (t *TestArbitraryTaskData) GetEvents() inter.Events {
	return nil
}

// GetResult is just an implementation
func (t *TestArbitraryTaskData) GetResult() []error {
	return nil
}

// GetOnValidatedFn is just an implementation
func (t *TestArbitraryTaskData) GetOnValidatedFn() OnValidatedFn {
	return nil
}

// testEnqueue tests Enqueue function
func testEnqueue(t *testing.T, checker *Checker, event inter.Events) {
	onValidatedFns := []func(ArbitraryTaskData) { func(ArbitraryTaskData) {}, }
	for _, fn := range onValidatedFns {
		err := checker.Enqueue(event, fn)
		require.Nil(t, err)
	}
}

// testValidate tests validate function
func testValidate(t *testing.T, checker *Checker, event *inter.Event) {
	err := checker.Validate(event)

	addrs, epoch := checker.reader.GetEpochPubKeys()
	if event.Epoch != epoch {
		require.Equal(t, epochcheck.ErrNotRelevant, err)
		return
	}

	addr, ok := addrs[event.Creator]
	if !ok {
		require.Equal(t, epochcheck.ErrAuth, err)
		return
	}

	if !event.VerifySignature(addr) {
		require.Equal(t, ErrWrongEventSig, err)
		return
	}

	for _, tx := range event.Transactions {
		_, err2 := types.Sender(checker.txSigner, tx)
		if err2 != nil {
			require.Equal(t, ErrMalformedTxSig, err)
			return
		}
	}

	if event.TxHash != types.DeriveSha(event.Transactions) {
		require.Equal(t, ErrWrongTxHash, err)
		return
	}

	require.Equal(t, nil, err)
}

// newTestWallet creates test wallet
func newTestWallet(t *testing.T) TestWallet {
	privateKey, err := crypto.GenerateKey()
	require.Nil(t, err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return TestWallet{
		PrivateKey: *privateKey,
		PubKey:     publicKeyBytes,
		Address:    address,
	}
}

// TestWallet is just a wallet for tests
type TestWallet struct {
	Address    common.Address
	PubKey     []byte
	PrivateKey ecdsa.PrivateKey
}

// makeParentsForTests creates parents for an event
func makeParentsForTests(num int) hash.Events {
	var hashEvents hash.Events
	var h common.Hash
	for i := num; i > 0; i-- {
		arrId := i % 32
		h[arrId] = h[arrId] + 1
		hashEvents = append(hashEvents, hash.Event(h))
	}
	return hashEvents
}

// makeDataWithLen creates array of bytes
func makeDataWithLen(len int) []byte {
	var data []byte
	for i := len; i > 0; i-- {
		data = append(data, 0x00)
	}
	return data
}

// makeTestEvents creates test events
func makeTestEvents(t *testing.T) []*inter.Event {
	var events []*inter.Event
	versions := []uint32{
		0, 1,
	}
	seqs := []idx.Event{
		0, 1, (math.MaxInt32 / 2) + 1,
	}
	epochs := []idx.Epoch{
		0, 1, (math.MaxInt32 / 2) + 1,
	}
	frames := []idx.Frame{
		0, 1, (math.MaxInt32 / 2) + 1,
	}
	lamports := []idx.Lamport{
		0, 1, (math.MaxInt32 / 2) + 1,
	}
	extras := [][]byte{
		[]byte{}, makeDataWithLen(1), makeDataWithLen(params.MaxExtraData + 1),
	}
	claimedTimes := []inter.Timestamp{
		0, 1, inter.Timestamp(uint64(time.Now().Unix())),
	}
	parentss := []hash.Events{
		makeParentsForTests(0),
		makeParentsForTests(1),
	}
	creators := []idx.StakerID{1, 102}
	txsSet := [][]*types.Transaction{makeTestTransactions(t), []*types.Transaction{}}

	for _, version := range versions {
		for _, seq := range seqs {
			for _, extra := range extras {
				for _, parents := range parentss {
					for _, epoch := range epochs {
						for _, frame := range frames {
							for _, lamport := range lamports {
								for _, claimedTime := range claimedTimes {
									for _, txs := range txsSet {
										for _, creator := range creators {
											event := inter.Event{
												EventHeader: inter.EventHeader{
													EventHeaderData: inter.EventHeaderData{
														Version:     version,
														Creator:     creator,
														Seq:         seq,
														Parents:     parents,
														Epoch:       epoch,
														Frame:       frame,
														Lamport:     lamport,
														ClaimedTime: claimedTime,
														Extra:       extra,
													}, Sig: nil,
												},
												Transactions: txs,
											}
											//pk := newTestWallet(t).PrivateKey
											//sig, err := crypto.Sign(event.DataToSign(), &pk)
											//require.Nil(t, err)
											//
											//event.Sig = sig
											events = append(events, &event)
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
	return events
}

// makeTestTransactions generates test transactions
func makeTestTransactions(t *testing.T) []*types.Transaction {
	var transactions []*types.Transaction
	transactions = append(transactions, newTransaction(newTestWallet(t).Address, big.NewInt(0), 1e10, big.NewInt(1), []byte{0x01}))
	return transactions
}

var nonce uint64 = 0

// newTransaction creates new testing transaction
func newTransaction(address common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *types.Transaction {
	defer func() { nonce += 1 }()
	tx := types.NewTransaction(nonce, address, amount, gasLimit, gasPrice, data)
	return tx
}
