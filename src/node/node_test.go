package node

import (
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/crypto"
	"github.com/Fantom-foundation/go-lachesis/src/dummy"
	"github.com/Fantom-foundation/go-lachesis/src/net"
	"github.com/Fantom-foundation/go-lachesis/src/peers"
	"github.com/Fantom-foundation/go-lachesis/src/poset"
	"github.com/Fantom-foundation/go-lachesis/src/utils"
	"github.com/sirupsen/logrus"
)

func initPeers(n int) ([]*ecdsa.PrivateKey, *peers.Peers) {
	var keys []*ecdsa.PrivateKey
	ps := peers.NewPeers()

	for i := 0; i < n; i++ {
		key, _ := crypto.GenerateECDSAKey()
		keys = append(keys, key)

		ps.AddPeer(peers.NewPeer(
			fmt.Sprintf("0x%X", crypto.FromECDSAPub(&keys[i].PublicKey)),
			fmt.Sprintf("127.0.0.1:%d", i),
		))
	}

	return keys, ps
}

func TestProcessSync(t *testing.T) {
	keys, p := initPeers(2)
	testLogger := common.NewTestLogger(t)
	config := TestConfig(t)

	// Start two nodes

	ps := p.ToPeerSlice()

	peer0Trans, err := net.NewTCPTransport(utils.GetUnusedNetAddr(t), nil, 2,
		time.Second, testLogger)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer peer0Trans.Close()

	node0 := NewNode(config, ps[0].ID, keys[0], p,
		poset.NewInmemStore(p, config.CacheSize),
		peer0Trans,
		dummy.NewInmemDummyApp(testLogger))
	node0.Init()

	node0.RunAsync(false)
	defer node0.Shutdown()

	peer1Trans, err := net.NewTCPTransport(utils.GetUnusedNetAddr(t), nil, 2,
		time.Second, testLogger)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer peer1Trans.Close()

	node1 := NewNode(config, ps[1].ID, keys[1], p,
		poset.NewInmemStore(p, config.CacheSize),
		peer1Trans,
		dummy.NewInmemDummyApp(testLogger))
	node1.Init()

	node1.RunAsync(false)
	defer node1.Shutdown()

	// Manually prepare SyncRequest and expected SyncResponse

	node0KnownEvents := node0.core.KnownEvents()
	node1KnownEvents := node1.core.KnownEvents()

	unknownEvents, err := node1.core.EventDiff(node0KnownEvents)
	if err != nil {
		t.Fatal(err)
	}

	unknownWireEvents, err := node1.core.ToWire(unknownEvents)
	if err != nil {
		t.Fatal(err)
	}

	args := net.SyncRequest{
		FromID: node0.id,
		Known:  node0KnownEvents,
	}
	expectedResp := net.SyncResponse{
		FromID: node1.id,
		Events: unknownWireEvents,
		Known:  node1KnownEvents,
	}

	// Make actual SyncRequest and check SyncResponse

	testLogger.Println("SYNCING...")
	time.Sleep(2000 * time.Millisecond)
	var out net.SyncResponse
	if err := peer0Trans.Sync(peer1Trans.LocalAddr(), &args, &out); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Verify the response
	if expectedResp.FromID != out.FromID {
		t.Fatalf("SyncResponse.FromID should be %d, not %d",
			expectedResp.FromID, out.FromID)
	}

	if l := len(out.Events); l != len(expectedResp.Events) {
		t.Fatalf("SyncResponse.Events should contain %d items, not %d",
			len(expectedResp.Events), l)
	}

	for i, e := range expectedResp.Events {
		ex := out.Events[i]
		if !reflect.DeepEqual(e.Body, ex.Body) {
			t.Fatalf("SyncResponse.Events[%d] should be %v, not %v",
				i, e.Body, ex.Body)
		}
	}

	if !reflect.DeepEqual(expectedResp.Known, out.Known) {
		t.Fatalf("SyncResponse.KnownEvents should be %#v, not %#v",
			expectedResp.Known, out.Known)
	}

}

func TestProcessEagerSync(t *testing.T) {
	keys, p := initPeers(2)
	testLogger := common.NewTestLogger(t)
	config := TestConfig(t)

	// Start two nodes

	ps := p.ToPeerSlice()

	peer0Trans, err := net.NewTCPTransport(utils.GetUnusedNetAddr(t), nil, 2,
		time.Second, testLogger)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer peer0Trans.Close()

	node0 := NewNode(config, ps[0].ID, keys[0], p,
		poset.NewInmemStore(p, config.CacheSize),
		peer0Trans,
		dummy.NewInmemDummyApp(testLogger))
	node0.Init()

	node0.RunAsync(false)
	defer node0.Shutdown()

	peer1Trans, err := net.NewTCPTransport(utils.GetUnusedNetAddr(t), nil, 2,
		time.Second, testLogger)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer peer1Trans.Close()

	node1 := NewNode(config, ps[1].ID, keys[1], p,
		poset.NewInmemStore(p, config.CacheSize),
		peer1Trans,
		dummy.NewInmemDummyApp(testLogger))
	node1.Init()

	node1.RunAsync(false)
	defer node1.Shutdown()

	// Manually prepare EagerSyncRequest and expected EagerSyncResponse

	node1KnownEvents := node1.core.KnownEvents()

	unknownEvents, err := node0.core.EventDiff(node1KnownEvents)
	if err != nil {
		t.Fatal(err)
	}

	unknownWireEvents, err := node0.core.ToWire(unknownEvents)
	if err != nil {
		t.Fatal(err)
	}

	args := net.EagerSyncRequest{
		FromID: node0.id,
		Events: unknownWireEvents,
	}
	expectedResp := net.EagerSyncResponse{
		FromID:  node1.id,
		Success: true,
	}

	time.Sleep(2000 * time.Millisecond)
	// Make actual EagerSyncRequest and check EagerSyncResponse
	var out net.EagerSyncResponse
	if err := peer0Trans.EagerSync(
		peer1Trans.LocalAddr(), &args, &out); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Verify the response
	if expectedResp.Success != out.Success {
		t.Fatalf("EagerSyncResponse.Sucess should be %v, not %v",
			expectedResp.Success, out.Success)
	}
}

func TestAddTransaction(t *testing.T) {
	keys, p := initPeers(2)
	testLogger := common.NewTestLogger(t)
	config := TestConfig(t)

	// Start two nodes

	ps := p.ToPeerSlice()

	peer0Trans, err := net.NewTCPTransport(utils.GetUnusedNetAddr(t), nil, 2,
		time.Second, common.NewTestLogger(t))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	peer0Proxy := dummy.NewInmemDummyApp(testLogger)
	defer peer0Trans.Close()

	node0 := NewNode(TestConfig(t), ps[0].ID, keys[0], p,
		poset.NewInmemStore(p, config.CacheSize),
		peer0Trans,
		peer0Proxy)
	node0.Init()

	node0.RunAsync(false)
	defer node0.Shutdown()

	peer1Trans, err := net.NewTCPTransport(utils.GetUnusedNetAddr(t), nil, 2,
		time.Second, common.NewTestLogger(t))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	peer1Proxy := dummy.NewInmemDummyApp(testLogger)
	defer peer1Trans.Close()

	node1 := NewNode(TestConfig(t), ps[1].ID, keys[1], p,
		poset.NewInmemStore(p, config.CacheSize),
		peer1Trans,
		peer1Proxy)
	node1.Init()

	node1.RunAsync(false)
	defer node1.Shutdown()
	// Submit a Tx to node0

	time.Sleep(2000 * time.Millisecond)
	message := "Hello World!"
	peer0Proxy.SubmitCh() <- []byte(message)

	// simulate a SyncRequest from node0 to node1

	node0KnownEvents := node0.core.KnownEvents()
	args := net.SyncRequest{
		FromID: node0.id,
		Known:  node0KnownEvents,
	}

	peer1Trans.LocalAddr()
	var out net.SyncResponse
	if err := peer0Trans.Sync(peer1Trans.LocalAddr(), &args, &out); err != nil {
		t.Fatal(err)
	}

	if err := node0.sync(out.Events); err != nil {
		t.Fatal(err)
	}

	// check the Tx was removed from the transactionPool
	// and added to the new Head

	if l := len(node0.core.transactionPool); l > 0 {
		t.Fatalf("node0's transactionPool should have 0 elements, not %d\n", l)
	}

	node0Head, _ := node0.core.GetHead()
	if l := len(node0Head.Transactions()); l != 1 {
		t.Fatalf("node0's Head should have 1 element, not %d\n", l)
	}

	if m := string(node0Head.Transactions()[0]); m != message {
		t.Fatalf("Transaction message should be '%s' not, not %s\n",
			message, m)
	}
}

func initNodes(keys []*ecdsa.PrivateKey,
	peers *peers.Peers,
	cacheSize int,
	syncLimit int64,
	storeType string,
	logger *logrus.Logger,
	t testing.TB) []*Node {

	var nodes []*Node

	for _, k := range keys {
		key := fmt.Sprintf("0x%X", crypto.FromECDSAPub(&k.PublicKey))
		peer, ok := peers.GetByPubKey(key)
		if !ok {
			t.Fatalf("participant %v nor found", key)
		}
		id := peer.ID

		conf := NewConfig(
			5*time.Millisecond,
			time.Second,
			cacheSize,
			syncLimit,
			logger,
		)

		trans, err := net.NewTCPTransport(utils.GetUnusedNetAddr(t),
			nil, 2, time.Second, logger)
		if err != nil {
			t.Fatalf("failed to create transport for peer %d: %s", id, err)
		}

		peers.Lock()
		peer.NetAddr = trans.LocalAddr()
		peers.Unlock()

		var store poset.Store
		switch storeType {
		case "badger":
			path, _ := ioutil.TempDir("", "badger")
			store, err = poset.NewBadgerStore(peers, conf.CacheSize, path)
			if err != nil {
				t.Fatalf("failed to create BadgerStore for peer %d: %s",
					id, err)
			}
		case "inmem":
			store = poset.NewInmemStore(peers, conf.CacheSize)
		}
		prox := dummy.NewInmemDummyApp(logger)

		node := NewNode(conf,
			id,
			k,
			peers,
			store,
			trans,
			prox)

		if err := node.Init(); err != nil {
			t.Fatalf("failed to initialize node%d: %s", id, err)
		}
		nodes = append(nodes, node)
	}
	return nodes
}

func recycleNodes(
	oldNodes []*Node, logger *logrus.Logger, t *testing.T) []*Node {
	var newNodes []*Node
	for _, oldNode := range oldNodes {
		newNode := recycleNode(oldNode, logger, t)
		newNodes = append(newNodes, newNode)
	}
	return newNodes
}

func recycleNode(oldNode *Node, logger *logrus.Logger, t *testing.T) *Node {
	conf := oldNode.conf
	id := oldNode.id
	key := oldNode.core.key
	ps := oldNode.peerSelector.Peers()

	var store poset.Store
	var err error
	if _, ok := oldNode.core.poset.Store.(*poset.BadgerStore); ok {
		store, err = poset.LoadBadgerStore(
			conf.CacheSize, oldNode.core.poset.Store.StorePath())
		if err != nil {
			t.Fatal(err)
		}
	} else {
		store = poset.NewInmemStore(oldNode.core.participants, conf.CacheSize)
	}

	trans, err := net.NewTCPTransport(oldNode.localAddr,
		nil, 2, time.Second, logger)
	if err != nil {
		t.Fatal(err)
	}
	prox := dummy.NewInmemDummyApp(logger)

	newNode := NewNode(conf, id, key, ps, store, trans, prox)

	if err := newNode.Init(); err != nil {
		t.Fatal(err)
	}

	return newNode
}

func runNodes(nodes []*Node, gossip bool) {
	for _, n := range nodes {
		node := n
		go func() {
			node.Run(gossip)
		}()
	}
}

func shutdownNodes(nodes []*Node) {
	for _, n := range nodes {
		n.Shutdown()
	}
}

func TestGossip(t *testing.T) {

	logger := common.NewTestLogger(t)

	keys, ps := initPeers(4)
	nodes := initNodes(keys, ps, 1000, 1000, "inmem", logger, t)

	target := int64(50)

	err := gossip(nodes, target, true, 13*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	checkGossip(nodes, 0, t)
}

func TestMissingNodeGossip(t *testing.T) {

	logger := common.NewTestLogger(t)

	keys, ps := initPeers(4)
	nodes := initNodes(keys, ps, 1000, 1000, "inmem", logger, t)
	defer shutdownNodes(nodes)

	err := gossip(nodes[1:], 10, true, 13*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	checkGossip(nodes[1:], 0, t)
}

func TestSyncLimit(t *testing.T) {

	logger := common.NewTestLogger(t)

	keys, ps := initPeers(4)
	nodes := initNodes(keys, ps, 1000, 1000, "inmem", logger, t)

	err := gossip(nodes, 10, false, 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer shutdownNodes(nodes)

	// create fake node[0] known to artificially reach SyncLimit
	node0KnownEvents := nodes[0].core.KnownEvents()
	for k := range node0KnownEvents {
		node0KnownEvents[k] = 0
	}

	args := net.SyncRequest{
		FromID: nodes[0].id,
		Known:  node0KnownEvents,
	}
	expectedResp := net.SyncResponse{
		FromID:    nodes[1].id,
		SyncLimit: true,
	}

	var out net.SyncResponse
	if err := nodes[0].trans.Sync(nodes[1].localAddr, &args, &out); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Verify the response
	if expectedResp.FromID != out.FromID {
		t.Fatalf("SyncResponse.FromID should be %d, not %d",
			expectedResp.FromID, out.FromID)
	}
	if !expectedResp.SyncLimit {
		t.Fatal("SyncResponse.SyncLimit should be true")
	}
}

func TestFastForward(t *testing.T) {

	logger := common.NewTestLogger(t)

	keys, ps := initPeers(4)
	nodes := initNodes(keys, ps, 1000, 1000,
		"inmem", logger, t)
	defer shutdownNodes(nodes)

	target := int64(20)
	err := gossip(nodes[1:], target, false, 15*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	err = nodes[0].fastForward()
	if err != nil {
		t.Fatalf("Error FastForwarding: %s", err)
	}

	lbi := nodes[0].core.GetLastBlockIndex()
	if lbi <= 0 {
		t.Fatalf("LastBlockIndex is too low: %d", lbi)
	}
	sBlock, err := nodes[0].GetBlock(lbi)
	if err != nil {
		t.Fatalf("Error retrieving latest Block"+
			" from reset hasposetraph: %v", err)
	}
	expectedBlock, err := nodes[1].GetBlock(lbi)
	if err != nil {
		t.Fatalf("Failed to retrieve block %d from node1: %v", lbi, err)
	}
	if !reflect.DeepEqual(sBlock.Body, expectedBlock.Body) {
		t.Fatalf("Blocks defer")
	}
}

func TestCatchUp(t *testing.T) {
	logger := common.NewTestLogger(t)

	// Create  config for 4 nodes
	keys, ps := initPeers(4)

	// Initialize the first 3 nodes only
	normalNodes := initNodes(keys[0:3], ps, 1000, 400, "inmem", logger, t)
	defer shutdownNodes(normalNodes)

	target := int64(50)

	err := gossip(normalNodes, target, false, 14*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	checkGossip(normalNodes, 0, t)

	node4 := initNodes(keys[3:], ps, 1000, 400, "inmem", logger, t)[0]

	// Run parallel routine to check node4 eventually reaches CatchingUp state.
	timeout := time.After(10 * time.Second)
	go func() {
		for {
			select {
			case <-timeout:
				t.Fatalf("Timeout waiting for node4 to enter CatchingUp state")
			default:
			}
			if node4.getState() == CatchingUp {
				break
			}
		}
	}()

	node4.RunAsync(true)
	defer node4.Shutdown()

	// Gossip some more
	nodes := append(normalNodes, node4)
	newTarget := target + 20
	err = bombardAndWait(nodes, newTarget, 10*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	start := node4.core.poset.FirstConsensusRound
	checkGossip(nodes, *start, t)
}

func TestFastSync(t *testing.T) {
	logger := common.NewTestLogger(t)

	// Create  config for 4 nodes
	keys, ps := initPeers(4)
	nodes := initNodes(keys, ps, 1000, 400, "inmem", logger, t)
	defer shutdownNodes(nodes)

	var target int64 = 50

	err := gossip(nodes, target, false, 13*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	checkGossip(nodes, 0, t)

	node4 := nodes[3]
	node4.Shutdown()

	secondTarget := target + 50
	err = bombardAndWait(nodes[0:3], secondTarget, 6*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	checkGossip(nodes[0:3], 0, t)

	// Can't re-run it; have to reinstantiate a new node.
	node4 = recycleNode(node4, logger, t)

	// Run parallel routine to check node4 eventually reaches CatchingUp state.
	timeout := time.After(6 * time.Second)
	go func() {
		for {
			select {
			case <-timeout:
				t.Fatalf("Timeout waiting for node4 to enter CatchingUp state")
			default:
			}
			if node4.getState() == CatchingUp {
				break
			}
		}
	}()

	node4.RunAsync(true)
	defer node4.Shutdown()

	nodes[3] = node4

	// Gossip some more
	thirdTarget := secondTarget + 20
	err = bombardAndWait(nodes, thirdTarget, 6*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	start := node4.core.poset.FirstConsensusRound
	checkGossip(nodes, *start, t)
}

func TestShutdown(t *testing.T) {
	logger := common.NewTestLogger(t)

	keys, ps := initPeers(4)
	nodes := initNodes(keys, ps, 1000, 1000, "inmem", logger, t)
	runNodes(nodes, false)

	nodes[0].Shutdown()

	err := nodes[1].gossip(nodes[0].localAddr, nil)
	if err == nil {
		t.Fatal("Expected Timeout Error")
	}

	nodes[1].Shutdown()
}

func TestBootstrapAllNodes(t *testing.T) {
	logger := common.NewTestLogger(t)

	os.RemoveAll("test_data")
	os.Mkdir("test_data", os.ModeDir|0777)

	// create a first network with BadgerStore
	// and wait till it reaches 10 consensus rounds before shutting it down
	keys, ps := initPeers(4)
	nodes := initNodes(keys, ps, 1000, 1000, "badger", logger, t)

	err := gossip(nodes, 10, false, 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	checkGossip(nodes, 0, t)
	shutdownNodes(nodes)

	// Now try to recreate a network from the databases created
	// in the first step and advance it to 20 consensus rounds
	newNodes := recycleNodes(nodes, logger, t)
	err = gossip(newNodes, 20, false, 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	checkGossip(newNodes, 0, t)
	shutdownNodes(newNodes)

	// Check that both networks did not have
	// completely different consensus events
	checkGossip([]*Node{nodes[0], newNodes[0]}, 0, t)
}

func gossip(
	nodes []*Node, target int64, shutdown bool, timeout time.Duration) error {
	runNodes(nodes, true)
	err := bombardAndWait(nodes, target, timeout)
	if err != nil {
		return err
	}
	if shutdown {
		shutdownNodes(nodes)
	}
	return nil
}

func bombardAndWait(nodes []*Node, target int64, timeout time.Duration) error {

	quit := make(chan struct{})
	makeRandomTransactions(nodes, quit)

	// wait until all nodes have at least 'target' blocks
	stopper := time.After(timeout)
	for {
		select {
		case <-stopper:
			return fmt.Errorf("timeout")
		default:
		}
		time.Sleep(10 * time.Millisecond)
		done := true
		for _, n := range nodes {
			ce := n.core.GetLastBlockIndex()
			if ce < target {
				done = false
				break
			} else {
				// wait until the target block has retrieved a state hash from
				// the app
				targetBlock, _ := n.core.poset.Store.GetBlock(target)
				if len(targetBlock.GetStateHash()) == 0 {
					done = false
					break
				}
			}
		}
		if done {
			break
		}
	}
	close(quit)
	return nil
}

func checkGossip(nodes []*Node, fromBlock int64, t *testing.T) {

	nodeBlocks := map[int64][]poset.Block{}
	for _, n := range nodes {
		var blocks []poset.Block
		for i := fromBlock; i < n.core.poset.Store.LastBlockIndex(); i++ {
			block, err := n.core.poset.Store.GetBlock(i)
			if err != nil {
				t.Fatalf("checkGossip: %v ", err)
			}
			blocks = append(blocks, block)
		}
		nodeBlocks[n.id] = blocks
	}

	minB := len(nodeBlocks[0])
	for k := int64(1); k < int64(len(nodes)); k++ {
		if len(nodeBlocks[k]) < minB {
			minB = len(nodeBlocks[k])
		}
	}

	for i, block := range nodeBlocks[0][:minB] {
		for k := int64(1); k < int64(len(nodes)); k++ {
			oBlock := nodeBlocks[k][i]
			if !reflect.DeepEqual(block.Body, oBlock.Body) {
				t.Fatalf("check gossip: difference in block %d."+
					" node 0: %v, node %d: %v",
					block.Index(), block.Body, k, oBlock.Body)
			}
		}
	}
}

func makeRandomTransactions(nodes []*Node, quit chan struct{}) {
	go func() {
		seq := make(map[int]int)
		for {
			select {
			case <-quit:
				return
			default:
				n := rand.Intn(len(nodes))
				node := nodes[n]
				submitTransaction(node, []byte(
					fmt.Sprintf("node%d transaction %d", n, seq[n])))
				seq[n] = seq[n] + 1
				time.Sleep(3 * time.Millisecond)
			}
		}
	}()
}

func submitTransaction(n *Node, tx []byte) error {

	n.proxy.SubmitCh() <- []byte(tx)
	return nil
}

func BenchmarkGossip(b *testing.B) {
	logger := common.NewTestLogger(b)
	for n := 0; n < b.N; n++ {
		keys, ps := initPeers(4)
		nodes := initNodes(keys, ps, 1000, 1000, "inmem", logger, b)
		gossip(nodes, 50, true, 3*time.Second)
	}
}
