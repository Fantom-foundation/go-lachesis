package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/sha3"
	"log"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"
)

type Config struct {
	RPCHost    string
	RPCPort    int
	OutPath    string
	LvlLimit   int
	OnlyEpoch  bool
	RenderFile bool
}

const defaultJsonKey = `{"address":"239fa7623354ec26520de878b52f13fe84b06971","crypto":{"cipher":"aes-128-ctr","ciphertext":"d25a3ce3381aef33d8c8e6345e3dc5547514cb4fe3daa78f1bb6a12b1b3a6400","cipherparams":{"iv":"0cdb328cb3b7d71b09efa90c90b23157"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"fb9296bf20661d0c4774514d914b1b0e33da86a616d179b604d06ca7196c6652"},"mac":"6396b71f2dabe1ae2d165692f7ff3f71a71a18ad5a71fbcd17a1f6d92cae15a6"},"id":"d2d05c35-66a9-4972-8216-c5a434cc72ff","version":3}`
var keyJson JsonKey
const contractAddr = "0xfc00face00000000000000000000000000000000"
const stakersNum = 3
const delegatorsNum = 1

var zeroInt = big.NewInt(0)

func TestSFC(t *testing.T) {
	json.Unmarshal([]byte(defaultJsonKey), &keyJson)
	createFakeNet()
	contract := connect(t, contractAddr)
	stakers := createStakers(t, contract, stakersNum)
	require.Equal(t, stakersNum, len(stakers))
}

func createFakeNet() {
	//
}

// this function is made to await an execution of a transaction
// since there is a certain period between adding tx to a pool and executing it
// should be replaced with a better approach
func waitForTx() {
	timeToAwait := time.Second * 3 // we assume that all transactions will be executed within the following block
	time.Sleep(timeToAwait)
}


func delegateWithDefault() {

}

// temprorary not used
func createDelegation(t *testing.T, contract *Main, staker createdStaker, delegatorsNum int64) {
	var callerOpts *bind.CallOpts = nil
	prevDelegNum, err := contract.DelegationsNum(callerOpts)
	require.Nil(t, err)

	delegation, err := contract.Delegations(nil, staker.Address)
	require.Nil(t, err)
	require.Equal(t, delegation.Amount.Cmp(zeroInt), 0)

	// replace a call below with a valid trOpts
	trOpts := defaultTransactionOpts(t)
	for i := delegatorsNum; i > 0; i-- {
		contract.CreateDelegation(trOpts, staker.StakerId)
	}
	log.Println(prevDelegNum)
}
func getMainStaker(t *testing.T) {

}

func createStakers(t *testing.T, contract *Main, stakersNum int64) createdStakers {
	var callerOpts *bind.CallOpts = nil
	var stakers = make(createdStakers, 0)
	prevStakersNum, err := contract.StakersNum(callerOpts)
	require.Nil(t, err)

	lastStakerId, err := contract.StakersLastID(callerOpts)
	require.Nil(t, err)

	for i := stakersNum; i > 0; i-- {
		trOpts := defaultTransactionOpts(t)
		addr := newFakeAddress()
		metadata := []byte{}
		tx, err := contract.CreateStakeWithAddresses(trOpts, addr, addr, metadata)
		require.Nil(t, err)

		log.Printf("CreateStake pending: 0x%x\n", tx.Hash())
		newStaker := createdStaker{
			Address:  addr,
			StakerId: big.NewInt(0).Add(lastStakerId, big.NewInt(1)),
		}
		stakers = append(stakers, newStaker)
		waitForTx()
	}

	newStakersNum, err := contract.StakersNum(callerOpts)
	require.Nil(t, err)
	stakersToAdd := big.NewInt(stakersNum)
	expectedStakersNum := big.NewInt(0).Add(prevStakersNum, stakersToAdd)
	require.Equal(t, newStakersNum.Cmp(expectedStakersNum), 0)
	return stakers
}

func connect(t *testing.T, contractAddressHex string) *Main {
	var contractAddress = common.HexToAddress(contractAddressHex)
	// Create an IPC based RPC connection to a remote node
	conn, err := ethclient.Dial("http://localhost:18545") //"localhost:18545 4001"
	if err != nil {
		require.Nil(t, err)
	}

	m, err := NewMain(contractAddress, conn)
	if err != nil {
		require.Nil(t, err)
	}
	return m
}

// legacy. needed for tests
// will be removed in an upcomming commit
func connect1() {
	var contractAddressHex = "0xfc00face00000000000000000000000000000000" // 0xfc00face00000000000000000000000000000000 0xfc00beef00000000000000000000000000000101
	var contractAddress = common.HexToAddress(contractAddressHex)
	// Create an IPC based RPC connection to a remote node
	conn, err := ethclient.Dial("http://localhost:18545") //"localhost:18545 4001"
	if err != nil {
		logAndExit(fmt.Sprintf("Failed to connect to the Ethereum client: %v", err))
	}

	tkn, err := NewMain(contractAddress, conn)
	if err != nil {
		logAndExit(fmt.Sprintf("Failed to instantiate a Token contract: %v", err))
	}
	stkrs(tkn)
	return

	log.Println("tkn:", tkn)
	epoch, err := tkn.CurrentEpoch(nil)
	if err != nil {
		logAndExit(fmt.Sprintf("Failed to instantiate a Token contract: %v", err))
	}
	log.Println("epoch1", epoch)

	br, err := tkn.BondedRatio(nil)
	if err != nil {
		logAndExit(fmt.Sprintf("Failed to get BondedRatio: %v", err))
	}
	log.Println("br", br)

	// Create an authorized transactor
	auth, err := NewMainTransactor(contractAddress, nil)
	if err != nil {
		logAndExit(fmt.Sprintf("Failed to create authorized transactor: %v", err))
	}
	log.Println("auth", auth)

	opts := makeOpts()
	//opts.GasLimit = 100000
	newEph := big.NewInt(100000)
	tx, err := tkn.MakeEpochSnapshots(opts, newEph)
	if err != nil {
		logAndExit(fmt.Sprintf("Failed to MakeEpochSnapshots: %v", err))
	}
	log.Printf("MakeEpochSnapshots pending: 0x%x\n", tx.Hash())
	log.Printf("tx.Data: %v", tx.Data())

	prog, err := conn.SyncProgress(context.TODO())
	if err != nil {
		logAndExit(fmt.Sprintf("Failed to SyncProgress: %v", err))
	}
	log.Printf("SyncProgress :%v", prog)

	pNum, err := conn.PendingTransactionCount(context.TODO())
	if err != nil {
		logAndExit(fmt.Sprintf("Failed to PendingTransactionCount: %v", err))
	}
	log.Println("PendingTransactionCount", pNum)

	//err = conn.SendTransaction(context.TODO(), tx)
	//if err != nil {
	//	log.Fatalf("Failed to SendTransaction: %v", err)
	//}

	log.Println("tx", tx)
	epoch, err = tkn.CurrentEpoch(nil)
	if err != nil {
		logAndExit(fmt.Sprintf("Failed to instantiate a Token contract: %v", err))
	}
	log.Println("epoch2", epoch)
}

// will be removed in an upcomming commit
func stkrs(m *Main) {
	stkNum, err := m.StakersNum(nil)
	if err != nil {
		logAndExit(fmt.Sprintf("Failed to get StakersNum: %v", err))
	}
	log.Println("stkNum", stkNum)

	stkId, err := m.StakersLastID(nil)
	if err != nil {
		logAndExit(fmt.Sprintf("Failed to get stkId: %v", err))
	}
	log.Println("stkId", stkId)

	stkr0, err := m.Stakers(nil, big.NewInt(0))
	if err != nil {
		logAndExit(fmt.Sprintf("Failed to get Staker: %v", err))
	}
	stkr1, err := m.Stakers(nil, big.NewInt(1))
	if err != nil {
		logAndExit(fmt.Sprintf("Failed to get Staker: %v", err))
	}
	stkr2, err := m.Stakers(nil, big.NewInt(2))
	if err != nil {
		logAndExit(fmt.Sprintf("Failed to get Staker: %v", err))
	}

	log.Println("stkr0", stkr0)
	log.Println("stkr1", stkr1)
	log.Println("stkr2", stkr2)
	log.Println("stkr2 SfcAddress", stkr2.SfcAddress)
	log.Println("stkr2 DagAddress", stkr2.DagAddress.Hex())

	opts := makeOpts()
	log.Println("opts", opts.Value.String())
	addr, amt := newStakerParams()
	opts.Value = amt
	tx, err := m.CreateStakeWithAddresses(opts, addr, addr, []byte{})
	if err != nil {
		logAndExit(fmt.Sprintf("Failed to CreateStakeWithAddresses: %v", err))
	}
	log.Printf("CreateStake pending: 0x%x\n", tx.Hash())

	stkNum, err = m.StakersNum(nil)
	if err != nil {
		logAndExit(fmt.Sprintf("Failed to get StakersNum: %v", err))
	}
	log.Println("stkNum", stkNum)
	return
	os.Exit(0)
}

func newStake() {

}

// will be removed in an upcomming commit
func makeOpts() *bind.TransactOpts {
	jsonKey := `{"address":"239fa7623354ec26520de878b52f13fe84b06971","crypto":{"cipher":"aes-128-ctr","ciphertext":"d25a3ce3381aef33d8c8e6345e3dc5547514cb4fe3daa78f1bb6a12b1b3a6400","cipherparams":{"iv":"0cdb328cb3b7d71b09efa90c90b23157"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"fb9296bf20661d0c4774514d914b1b0e33da86a616d179b604d06ca7196c6652"},"mac":"6396b71f2dabe1ae2d165692f7ff3f71a71a18ad5a71fbcd17a1f6d92cae15a6"},"id":"d2d05c35-66a9-4972-8216-c5a434cc72ff","version":3}`
	opts, err := bind.NewTransactor(strings.NewReader(jsonKey), "fakepassword")
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	opts.Value = bigintScaled(3175100)
	return opts
}

// will be removed in an upcomming commit
func defaultTransactionOpts(t *testing.T) *bind.TransactOpts {

	pass := "fakepassword"
	amt := int64(3175100)
	tr, err := newTransactionOpts(defaultJsonKey, pass, amt)
	require.Nil(t, err)
	return tr
}

func newTransactionOpts(jsonKey, password string, valueToScale int64) (*bind.TransactOpts, error) {
	opts, err := bind.NewTransactor(strings.NewReader(jsonKey), password)
	if err != nil {
		return nil, err
	}
	opts.Value = bigintScaled(valueToScale)
	return opts, nil
}

func bigintWithBase(base *big.Int, value *big.Int) *big.Int {
	return base.Mul(base, value)
}

func bigintScaled(value int64) *big.Int {
	var e18 = big.NewInt(1e18)
	val := big.NewInt(value)
	return bigintWithBase(e18, val)
}

func newStakerParams() (common.Address, *big.Int) {
	// Create an account
	key, err := crypto.GenerateKey()
	if err != nil {
		logAndExit(err.Error())
	}

	return crypto.PubkeyToAddress(key.PublicKey), bigintScaled(64984764)
}

func newFakeAddress() common.Address {
	key, err := crypto.GenerateKey()
	if err != nil {
		logAndExit(err.Error())
	}

	return crypto.PubkeyToAddress(key.PublicKey)
}

func opsRandAddr() *bind.TransactOpts {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println(hexutil.Encode(privateKeyBytes)[2:]) // 0xfad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println(hexutil.Encode(publicKeyBytes)[4:]) // 0x049a7df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b97934fdd7719d0810951e03418205868a5c1b40b192451367f28e0088dd75e15de40c05

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address) // 0x96216849c49358B10257cb55b28eA603c874b05E

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])

	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:]))
	return nil
}

// just not to panic
// will be removed in an upcomming commit
func logAndExit(txt string) {
	log.Println(txt)
	os.Exit(1)
}
