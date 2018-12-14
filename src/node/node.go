package node

import (
	"crypto/ecdsa"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"strconv"

	"github.com/Fantom-foundation/go-lachesis/src/net"
	"github.com/Fantom-foundation/go-lachesis/src/peers"
	"github.com/Fantom-foundation/go-lachesis/src/poset"
	"github.com/Fantom-foundation/go-lachesis/src/proxy"
)

type Node struct {
	nodeState

	conf   *Config
	logger *logrus.Entry

	id       uint32
	core     *Core
	coreLock sync.Mutex

	localAddr string

	peerSelector PeerSelector
	selectorLock sync.Mutex

	trans net.Transport
	netCh <-chan net.RPC

	proxy            proxy.AppProxy
	submitCh         chan []byte
	submitInternalCh chan poset.InternalTransaction

	commitCh chan poset.Block

	shutdownCh chan struct{}

	controlTimer *ControlTimer

	start        time.Time
	syncRequests int
	syncErrors   int

	needBoostrap bool
	gossipJobs   count64
	rpcJobs      count64
}

func NewNode(conf *Config,
	id uint32,
	key *ecdsa.PrivateKey,
	peerSet *peers.PeerSet,
	store poset.Store,
	trans net.Transport,
	proxy proxy.AppProxy) *Node {

	localAddr := trans.LocalAddr()

	pmap, _ := store.GetLastPeerSet()

	commitCh := make(chan poset.Block, 400)
	core := NewCore(id, key, pmap, store, commitCh, conf.Logger)

	pubKey := core.HexID()

//	peerSelector := NewRandomPeerSelector(peerSet, localAddr)
	peerSelector := NewSmartPeerSelector(peerSet, pubKey,
		core.poset.GetFlagTableOfRandomUndeterminedEvent)

	node := Node{
		id:               id,
		conf:             conf,
		core:             core,
		localAddr:        localAddr,
		logger:           conf.Logger.WithField("this_id", id),
		peerSelector:     peerSelector,
		trans:            trans,
		netCh:            trans.Consumer(),
		proxy:            proxy,
		submitCh:         proxy.SubmitCh(),
		submitInternalCh: proxy.SubmitInternalCh(),
		commitCh:         commitCh,
		shutdownCh:       make(chan struct{}),
		controlTimer:     NewRandomControlTimer(),
		start:            time.Now(),
		gossipJobs:       0,
		rpcJobs:          0,
	}

	node.logger.WithField("peers", pmap).Debug("pmap")
	node.logger.WithField("pubKey", pubKey).Debug("pubKey")

	node.needBoostrap = store.NeedBoostrap()

	// Initialize
	node.setState(Gossiping)

	return &node
}

func (n *Node) Init() error {
	var peerAddresses []string
	for _, p := range n.peerSelector.PeerSet().Peers {
		peerAddresses = append(peerAddresses, p.NetAddr)
	}
	n.logger.WithField("peers", peerAddresses).Debug("Initialize Node")

	if n.needBoostrap {
		n.logger.Debug("Bootstrap")
		if err := n.core.Bootstrap(); err != nil {
			return err
		}
	}
	n.Register()

	return n.core.SetHeadAndSeq()
}

func (n *Node) RunAsync(gossip bool) {
	n.logger.Debug("RunAsync(gossip bool)")
	go n.Run(gossip)
}

func (n *Node) Run(gossip bool) {
	// The ControlTimer allows the background routines to control the
	// heartbeat timer when the node is in the Gossiping state. The timer should
	// only be running when there are uncommitted transactions in the system.
	go n.controlTimer.Run(n.conf.HeartbeatTimeout)

	// Execute some background work regardless of the state of the node.
	// Process SubmitTx and CommitBlock requests
	go n.doBackgroundWork()

	// pause before gossiping test transactions to allow all nodes come up
	time.Sleep(time.Duration(n.conf.TestDelay) * time.Second)

	// Execute Node State Machine
	for {
		// Run different routines depending on node state
		state := n.getState()
		n.logger.WithField("state", state.String()).Debug("RunAsync(gossip bool)")

		switch state {
		case Gossiping:
			n.lachesis(gossip)
		case CatchingUp:
			n.fastForward()
		case Stop:
			// do nothing in Stop state
		case Shutdown:
			return
		}
	}
}

func (n *Node) resetTimer() {
	if !n.controlTimer.GetSet() {
		ts := n.conf.HeartbeatTimeout
		//Slow gossip if nothing interesting to say
		if n.core.poset.GetPendingLoadedEvents() == 0 &&
			n.core.GetTransactionPoolCount() == 0 &&
			n.core.GetBlockSignaturePoolCount() == 0 {
			ts = time.Duration(time.Second)
		}
		n.controlTimer.resetCh <- ts
	}
}

func (n *Node) doBackgroundWork() {
	for {
		select {
		case t := <-n.submitCh:
			n.logger.Debug("Adding Transactions to Transaction Pool")
			n.addTransaction(t)
			n.resetTimer()
		case t := <-n.submitInternalCh:
			n.logger.Debug("Adding Internal Transaction")
			n.addInternalTransaction(&t)
			n.resetTimer()
		case block := <-n.commitCh:
			n.logger.WithFields(logrus.Fields{
				"index":          block.Index(),
				"round_received": block.RoundReceived(),
				"transactions":   len(block.Transactions()),
			}).Debug("Adding EventBlock")
			if err := n.commit(block); err != nil {
				n.logger.WithField("error", err).Error("Adding EventBlock")
			}
		case <-n.shutdownCh:
			return
		}
	}
}

// lachesis is interrupted when a gossip function, launched asynchronously, changes
// the state from Gossiping to CatchingUp, or when the node is shutdown.
// Otherwise, it processes RPC requests, periodicaly initiates gossip while there
// is something to gossip about, or waits.
func (n *Node) lachesis(gossip bool) {
	returnCh := make(chan struct{}, 100)
	for {
		select {
		case rpc := <-n.netCh:
			n.goFunc(func() {
				n.rpcJobs.increment()
				n.logger.Debug("Processing RPC")
				n.processRPC(rpc)
				n.resetTimer()
				n.rpcJobs.decrement()
			})
		case <-n.controlTimer.tickCh:
			if gossip && n.gossipJobs.get() < 1 {
				n.selectorLock.Lock()
				peerAddr := n.peerSelector.Next().NetAddr
				n.selectorLock.Unlock()
				n.goFunc(func() {
					n.gossipJobs.increment()
					n.gossip(peerAddr, returnCh)
					n.gossipJobs.decrement()
				})
				n.logger.Debug("Gossip")
			}
			n.logStats()
			n.resetTimer()
		case <-returnCh:
			return
		case <-n.shutdownCh:
			return
		}
	}
}

func (n *Node) processRPC(rpc net.RPC) {
	switch cmd := rpc.Command.(type) {
	case *net.SyncRequest:
		n.processSyncRequest(rpc, cmd)
	case *net.EagerSyncRequest:
		n.processEagerSyncRequest(rpc, cmd)
	case *net.FastForwardRequest:
		n.processFastForwardRequest(rpc, cmd)
	default:
		n.logger.WithField("cmd", rpc.Command).Error("Unexpected RPC command")
		rpc.Respond(nil, fmt.Errorf("unexpected command"))
	}
}

func (n *Node) processSyncRequest(rpc net.RPC, cmd *net.SyncRequest) {
	n.logger.WithFields(logrus.Fields{
		"from_id": cmd.FromID,
		"known":   cmd.Known,
	}).Debug("processSyncRequest(rpc net.RPC, cmd *net.SyncRequest)")

	resp := &net.SyncResponse{
		FromID: n.id,
	}
	var respErr error

	// Check sync limit
	n.coreLock.Lock()
	overSyncLimit := n.core.OverSyncLimit(cmd.Known, n.conf.SyncLimit)
	n.coreLock.Unlock()
	if overSyncLimit {
		n.logger.Debug("n.core.OverSyncLimit(cmd.Known, n.conf.SyncLimit)")
		resp.SyncLimit = true
	} else {
		// Compute Diff
		start := time.Now()
		n.coreLock.Lock()
		eventDiff, err := n.core.EventDiff(cmd.Known)
		n.coreLock.Unlock()
		elapsed := time.Since(start)
		n.logger.WithField("Duration", elapsed.Nanoseconds()).Debug("n.core.EventBlockDiff(cmd.Known)")
		if err != nil {
			n.logger.WithField("Error", err).Error("n.core.EventBlockDiff(cmd.Known)")
			respErr = err
		}

		// Convert to WireEvents
		wireEvents, err := n.core.ToWire(eventDiff)
		if err != nil {
			n.logger.WithField("error", err).Debug("n.core.TransportEventBlock(eventDiff)")
			respErr = err
		} else {
			resp.Events = wireEvents
		}
	}

	// Get Self Known
	n.coreLock.Lock()
	knownEvents := n.core.KnownEvents()
	n.coreLock.Unlock()
	resp.Known = knownEvents

	n.logger.WithFields(logrus.Fields{
		"events":     len(resp.Events),
		"known":      resp.Known,
		"sync_limit": resp.SyncLimit,
		"error":      respErr,
	}).Debug("SyncRequest Received")

	rpc.Respond(resp, respErr)
}

func (n *Node) processEagerSyncRequest(rpc net.RPC, cmd *net.EagerSyncRequest) {
	n.logger.WithFields(logrus.Fields{
		"from_id": cmd.FromID,
		"events":  len(cmd.Events),
	}).Debug("processEagerSyncRequest(rpc net.RPC, cmd *net.EagerSyncRequest)")

	success := true
	n.coreLock.Lock()
	err := n.sync(cmd.Events)
	n.coreLock.Unlock()
	if err != nil {
		n.logger.WithField("error", err).Error("n.sync(cmd.Events)")
		success = false
	}

	resp := &net.EagerSyncResponse{
		FromID:  n.id,
		Success: success,
	}
	rpc.Respond(resp, err)
}

func (n *Node) processFastForwardRequest(rpc net.RPC, cmd *net.FastForwardRequest) {
	n.logger.WithFields(logrus.Fields{
		"from": cmd.FromID,
	}).Debug("processFastForwardRequest(rpc net.RPC, cmd *net.FastForwardRequest)")

	resp := &net.FastForwardResponse{
		FromID: n.id,
	}
	var respErr error

	// Get latest Frame
	n.coreLock.Lock()
	block, frame, err := n.core.GetAnchorBlockWithFrame()
	n.coreLock.Unlock()
	if err != nil {
		n.logger.WithField("error", err).Error("n.core.GetAnchorBlockWithFrame()")
		respErr = err
	} else {
		resp.Block = *block
		resp.Frame = *frame

		// Get snapshot
		snapshot, err := n.proxy.GetSnapshot(block.Index())
		if err != nil {
			n.logger.WithField("error", err).Error("n.proxy.GetSnapshot(block.Index())")
			respErr = err
		}
		resp.Snapshot = snapshot
	}

	n.logger.WithFields(logrus.Fields{
		"Events": len(resp.Frame.Events),
		"Error":  respErr,
	}).Debug("FastForwardRequest Received")
	rpc.Respond(resp, respErr)
}

// This function is usually called in a go-routine and needs to inform the
// calling routine (usually the lachesis routine) when it is time to exit the
// Gossiping state and return.
func (n *Node) gossip(peerAddr string, parentReturnCh chan struct{}) error {

	// pull
	syncLimit, otherKnownEvents, err := n.pull(peerAddr)
	if err != nil {
		return err
	}

	// check and handle syncLimit
	if syncLimit {
		n.logger.WithField("from", peerAddr).Debug("SyncLimit")
		n.setState(CatchingUp)
		parentReturnCh <- struct{}{}
		return nil
	}

	// push
	err = n.push(peerAddr, otherKnownEvents)
	if err != nil {
		return err
	}

	// update peer selector
	n.selectorLock.Lock()
	n.peerSelector.UpdateLast(peerAddr)
	n.selectorLock.Unlock()

	return nil
}

func (n *Node) pull(peerAddr string) (syncLimit bool, otherKnownEvents map[uint32]int64, err error) {
	// Compute Known
	n.coreLock.Lock()
	knownEvents := n.core.KnownEvents()
	n.coreLock.Unlock()

	// Send SyncRequest
	start := time.Now()
	resp, err := n.requestSync(peerAddr, knownEvents)
	elapsed := time.Since(start)
	n.logger.WithField("Duration", elapsed.Nanoseconds()).Debug("n.requestSync(peerAddr, knownEvents)")
	// FIXIT: should we catch io.EOF error here and how we process it?
	//	if err == io.EOF {
	//		return false, nil, nil
	//	}
	if err != nil {
		n.logger.WithField("Error", err).Error("n.requestSync(peerAddr, knownEvents)")
		return false, nil, err
	}
	n.logger.WithFields(logrus.Fields{
		"from_id":     resp.FromID,
		"sync_limit":  resp.SyncLimit,
		"events":      len(resp.Events),
		"known":       resp.Known,
		"knownEvents": knownEvents,
	}).Debug("SyncResponse")

	if resp.SyncLimit {
		return true, nil, nil
	}

	// Add Events to poset and create new Head if necessary
	n.coreLock.Lock()
	err = n.sync(resp.Events)
	n.coreLock.Unlock()
	if err != nil {
		n.logger.WithField("error", err).Error("n.sync(resp.Events)")
		return false, nil, err
	}

	return false, resp.Known, nil
}

func (n *Node) push(peerAddr string, knownEvents map[uint32]int64) error {

	// Check SyncLimit
	n.coreLock.Lock()
	overSyncLimit := n.core.OverSyncLimit(knownEvents, n.conf.SyncLimit)
	n.coreLock.Unlock()
	if overSyncLimit {
		n.logger.Debug("n.core.OverSyncLimit(knownEvents, n.conf.SyncLimit)")
		return nil
	}

	// Compute Diff
	start := time.Now()
	n.coreLock.Lock()
	eventDiff, err := n.core.EventDiff(knownEvents)
	n.coreLock.Unlock()
	elapsed := time.Since(start)
	n.logger.WithField("Duration", elapsed.Nanoseconds()).Debug("n.core.EventDiff(knownEvents)")
	if err != nil {
		n.logger.WithField("Error", err).Error("n.core.EventDiff(knownEvents)")
		return err
	}

	if len(eventDiff) > 0 {
		// Convert to WireEvents
		wireEvents, err := n.core.ToWire(eventDiff)
		if err != nil {
			n.logger.WithField("Error", err).Debug("n.core.TransferEventBlock(eventDiff)")
			return err
		}

		// Create and Send EagerSyncRequest
		start = time.Now()
		n.logger.WithField("wireEvents", wireEvents).Debug("Sending n.requestEagerSync.wireEvents")
		resp2, err := n.requestEagerSync(peerAddr, wireEvents)
		elapsed = time.Since(start)
		n.logger.WithField("Duration", elapsed.Nanoseconds()).Debug("n.requestEagerSync(peerAddr, wireEvents)")
		if err != nil {
			n.logger.WithField("Error", err).Error("n.requestEagerSync(peerAddr, wireEvents)")
			return err
		}
		n.logger.WithFields(logrus.Fields{
			"from_id": resp2.FromID,
			"success": resp2.Success,
		}).Debug("EagerSyncResponse")
	}

	return nil
}

func (n *Node) fastForward() error {
	n.logger.Debug("fastForward()")

	// wait until sync routines finish
	n.waitRoutines()

	// fastForwardRequest
	n.selectorLock.Lock()
	peer := n.peerSelector.Next()
	n.selectorLock.Unlock()
	start := time.Now()
	resp, err := n.requestFastForward(peer.NetAddr)
	elapsed := time.Since(start)
	n.logger.WithField("Duration", elapsed.Nanoseconds()).Debug("n.requestFastForward(peer.NetAddr)")
	if err != nil {
		n.logger.WithField("Error", err).Error("n.requestFastForward(peer.NetAddr)")
		return err
	}
	n.logger.WithFields(logrus.Fields{
		"from_id":              resp.FromID,
		"block_index":          resp.Block.Index(),
		"block_round_received": resp.Block.RoundReceived(),
		"frame_events":         len(resp.Frame.Events),
		"frame_roots":          resp.Frame.Roots,
		"snapshot":             resp.Snapshot,
	}).Debug("FastForwardResponse")

	// prepare core. ie: fresh poset
	n.coreLock.Lock()
	err = n.core.FastForward(peer.PubKeyHex, &resp.Block, &resp.Frame)
	n.coreLock.Unlock()
	if err != nil {
		n.logger.WithField("Error", err).Error("n.core.FastForward(peer.PubKeyHex, resp.Block, resp.Frame)")
		return err
	}

	// update app from snapshot
	err = n.proxy.Restore(resp.Snapshot)
	if err != nil {
		n.logger.WithField("Error", err).Error("n.proxy.Restore(resp.Snapshot)")
		return err
	}

	n.setState(Gossiping)

	return nil
}

func (n *Node) requestSync(target string, known map[uint32]int64) (net.SyncResponse, error) {

	args := net.SyncRequest{
		FromID: n.id,
		Known:  known,
	}

	var out net.SyncResponse
	err := n.trans.Sync(target, &args, &out)
	//n.logger.WithField("out", out).Debug("requestSync(target string, known map[int]int)")
	return out, err
}

func (n *Node) requestEagerSync(target string, events []poset.WireEvent) (net.EagerSyncResponse, error) {
	args := net.EagerSyncRequest{
		FromID: n.id,
		Events: events,
	}

	var out net.EagerSyncResponse
	n.logger.WithFields(logrus.Fields{
		"target": target,
	}).Debug("requestEagerSync(target string, events []poset.WireEvent)")
	err := n.trans.EagerSync(target, &args, &out)

	return out, err
}

func (n *Node) requestFastForward(target string) (net.FastForwardResponse, error) {
	n.logger.WithFields(logrus.Fields{
		"target": target,
	}).Debug("requestFastForward(target string) (net.FastForwardResponse, error)")

	args := net.FastForwardRequest{
		FromID: n.id,
	}

	var out net.FastForwardResponse
	err := n.trans.FastForward(target, &args, &out)

	return out, err
}

func (n *Node) sync(events []poset.WireEvent) error {
	// Insert Events in Poset and create new Head if necessary
	start := time.Now()
	err := n.core.Sync(events)
	elapsed := time.Since(start)
	n.logger.WithField("Duration", elapsed.Nanoseconds()).Debug("n.core.Sync(events)")
	if err != nil {
		return err
	}

	// Run consensus methods
	start = time.Now()
	err = n.core.RunConsensus()
	elapsed = time.Since(start)
	n.logger.WithField("Duration", elapsed.Nanoseconds()).Debug("n.core.RunConsensus()")
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) commit(block poset.Block) error {

	commitResponse, err := n.proxy.CommitBlock(block)
	if err != nil {
		n.logger.WithError(err).Debug("commit(block poset.Block)")
	}

	n.logger.WithFields(logrus.Fields{
		"block":      block.Index(),
		"state_hash": fmt.Sprintf("%X", commitResponse.StateHash),
		"accepted_internal_transactions": commitResponse.AcceptedInternalTransactions,
		// "err":        err,
	}).Debug("commit(eventBlock poset.EventBlock)")

	// XXX what do we do in case of error. Retry? This has to do with the
	// Lachesis <-> App interface. Think about it.

	// An error here could be that the endpoint is not configured, not all
	// nodes will be sending blocks to clients, in these cases -no_client can be
	// used, alternatively should check for the error here and handle it
	// appropriately

	// There is no point in using the stateHash if we know it is wrong
	// if err == nil {
	if true {
		// inmem statehash would be different than proxy statehash
		// inmem is simply the hash of transactions
		// this requires a 1:1 relationship with nodes and clients
		// multiple nodes can't read from the same client

		// Handle the response to set Block StateHash and process accepted
		// InternalTransactions which might update the PeerSet.
		block.StateHash = commitResponse.StateHash
		n.coreLock.Lock()
		defer n.coreLock.Unlock()
		sig, err := n.core.SignBlock(&block)
		if err != nil {
			return err
		}
		n.core.AddBlockSignature(sig)
		err = n.core.ProcessAcceptedInternalTransactions(block.RoundReceived(), commitResponse.AcceptedInternalTransactions)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *Node) addTransaction(tx []byte) {
	// we do not need coreLock here as n.core.AddTransactions has TransactionPoolLocker
	n.core.AddTransactions([][]byte{tx})
}

func (n *Node) addInternalTransaction(tx *poset.InternalTransaction) {
	n.coreLock.Lock()
	defer n.coreLock.Unlock()
	n.core.AddInternalTransactions([]*poset.InternalTransaction{tx})
}

func (n *Node) Shutdown() {
	if n.getState() != Shutdown {
		// n.mqtt.FireEvent("Shutdown()", "/mq/lachesis/node")
		n.logger.Debug("Shutdown()")

		// Exit any non-shutdown state immediately
		n.setState(Shutdown)

		// Stop and wait for concurrent operations
		close(n.shutdownCh)
		n.waitRoutines()

		// For some reason this needs to be called after closing the shutdownCh
		// Not entirely sure why...
		n.controlTimer.Shutdown()

		// transport and store should only be closed once all concurrent operations
		// are finished otherwise they will panic trying to use close objects
		n.trans.Close()
		n.core.poset.Store.Close()
	}
}

func (n *Node) GetStats() map[string]string {
	toString := func(i int64) string {
		if i <= 0 {
			return "nil"
		}
		return strconv.FormatInt(i, 10)
	}

	timeElapsed := time.Since(n.start)

	consensusEvents := n.core.GetConsensusEventsCount()
	consensusEventsPerSecond := float64(consensusEvents) / timeElapsed.Seconds()
	consensusTransactions := n.core.GetConsensusTransactionsCount()
	transactionsPerSecond := float64(consensusTransactions) / timeElapsed.Seconds()

	lastConsensusRound := n.core.GetLastConsensusRound()
	var consensusRoundsPerSecond float64
	if lastConsensusRound != -1 {
		consensusRoundsPerSecond = float64(lastConsensusRound) / timeElapsed.Seconds()
	}

	s := map[string]string{
		"last_consensus_round":    toString(lastConsensusRound),
		"time_elapsed":            strconv.FormatFloat(timeElapsed.Seconds(), 'f', 2, 64),
		"heartbeat":               strconv.FormatFloat(n.conf.HeartbeatTimeout.Seconds(), 'f', 2, 64),
		"node_current":            strconv.FormatInt(time.Now().Unix(), 10),
		"node_start":              strconv.FormatInt(n.start.Unix(), 10),
		"last_block_index":        strconv.FormatInt(n.core.GetLastBlockIndex(), 10),
		"consensus_events":        strconv.FormatInt(consensusEvents, 10),
		"sync_limit":              strconv.FormatInt(n.conf.SyncLimit, 10),
		"consensus_transactions":  strconv.FormatUint(consensusTransactions, 10),
		"undetermined_events":     strconv.Itoa(len(n.core.GetUndeterminedEvents())),
		"transaction_pool":        strconv.FormatInt(n.core.GetTransactionPoolCount(), 10),
		"num_peers":               strconv.Itoa(n.peerSelector.PeerSet().Len()),
		"sync_rate":               strconv.FormatFloat(n.SyncRate(), 'f', 2, 64),
		"transactions_per_second": strconv.FormatFloat(transactionsPerSecond, 'f', 2, 64),
		"events_per_second":       strconv.FormatFloat(consensusEventsPerSecond, 'f', 2, 64),
		"rounds_per_second":       strconv.FormatFloat(consensusRoundsPerSecond, 'f', 2, 64),
		"round_events":            strconv.Itoa(n.core.GetLastCommittedRoundEventsCount()),
		"id":                      fmt.Sprint(n.id),
		"state":                   n.getState().String(),
	}
	// n.mqtt.FireEvent(s, "/mq/lachesis/stats")
	return s
}

func (n *Node) logStats() {
	stats := n.GetStats()
	n.logger.WithFields(logrus.Fields{
		"last_consensus_round":   stats["last_consensus_round"],
		"last_block_index":       stats["last_block_index"],
		"consensus_events":       stats["consensus_events"],
		"consensus_transactions": stats["consensus_transactions"],
		"undetermined_events":    stats["undetermined_events"],
		"transaction_pool":       stats["transaction_pool"],
		"num_peers":              stats["num_peers"],
		"sync_rate":              stats["sync_rate"],
		"events/s":               stats["events_per_second"],
		"t/s":                    stats["transactions_per_second"],
		"rounds/s":               stats["rounds_per_second"],
		"round_events":           stats["round_events"],
		"id":                     stats["id"],
		"state":                  stats["state"],
		"z_gossipJobs":           n.gossipJobs.get(),
		"z_rpcJobs":              n.rpcJobs.get(),
		"addr":                   n.localAddr,
		"pending_loaded_events":  n.GetPendingLoadedEvents(),
		"last_round":             n.GetLastRound(),
	}).Warn("logStats()")
}

func (n *Node) SyncRate() float64 {
	var syncErrorRate float64
	if n.syncRequests != 0 {
		syncErrorRate = float64(n.syncErrors) / float64(n.syncRequests)
	}
	return 1 - syncErrorRate
}

func (n *Node) GetParticipants() (*peers.PeerSet, error) {
	return n.core.poset.Store.GetLastPeerSet()
}

func (n *Node) GetEvent(event string) (*poset.Event, error) {
	return n.core.poset.Store.GetEvent(event)
}

func (n *Node) GetLastEventFrom(participant string) (string, bool, error) {
	return n.core.poset.Store.LastEventFrom(participant)
}

func (n *Node) GetKnownEvents() map[uint32]int64 {
	return n.core.poset.Store.KnownEvents()
}

func (n *Node) GetEvents() (map[uint32]int64, error) {
	res := n.core.KnownEvents()
	return res, nil
}

func (n *Node) GetConsensusEvents() []string {
	return n.core.poset.Store.ConsensusEvents()
}

func (n *Node) GetConsensusTransactionsCount() uint64 {
	return n.core.GetConsensusTransactionsCount()
}

func (n *Node) GetPendingLoadedEvents() int64 {
	return n.core.GetPendingLoadedEvents()
}

func (n *Node) GetRound(roundIndex int64) (*poset.RoundCreated, error) {
	return n.core.poset.Store.GetRoundCreated(roundIndex)
}

func (n *Node) GetLastRound() int64 {
	return n.core.poset.Store.LastRound()
}

func (n *Node) GetRoundWitnesses(roundIndex int64) []string {
	return n.core.poset.Store.RoundWitnesses(roundIndex)
}

func (n *Node) GetRoundEvents(roundIndex int64) int {
	return n.core.poset.Store.RoundEvents(roundIndex)
}

func (n *Node) GetRoot(rootIndex string) (*poset.Root, error) {
	return n.core.poset.Store.GetRoot(rootIndex)
}

func (n *Node) GetBlock(blockIndex int64) (*poset.Block, error) {
	return n.core.poset.Store.GetBlock(blockIndex)
}

func (n *Node) ID() uint32 {
	return n.id
}

func (n *Node) Stop() {
	n.setState(Stop)
}
