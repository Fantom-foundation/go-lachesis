package node

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/Fantom-foundation/go-lachesis/src/peer"
	"github.com/Fantom-foundation/go-lachesis/src/peers"
	"github.com/Fantom-foundation/go-lachesis/src/poset"
	"github.com/Fantom-foundation/go-lachesis/src/proxy"
	"github.com/Fantom-foundation/go-lachesis/src/service"
)

// Stats used for logStats.
type Stats interface {
	Stats() map[string]string
	LastRound() int64
}

// Node struct that keeps all high level node functions
type Node struct {
	*nodeState2
	stats Stats

	conf   *Config
	logger *logrus.Entry

	id       uint64
	core     *Core
	coreLock sync.Mutex

	localAddr string

	peerSelector PeerSelector

	trans peer.SyncPeer
	proxy proxy.AppProxy

	submitCh         chan []byte
	submitInternalCh chan poset.InternalTransaction
	commitCh         chan poset.Block
	shutdownCh       chan struct{}
	signalTERMch     chan os.Signal

	controlTimer *ControlTimer

	start        time.Time
	syncRequests int
	syncErrors   int

	needBootstrap bool
	gossipJobs    count64
	rpcJobs       count64
}

// NewNode create a new node struct
func NewNode(conf *Config,
	id uint64,
	key *ecdsa.PrivateKey,
	participants *peers.Peers,
	pst *poset.Poset,
	commitCh chan poset.Block,
	needBootstrap bool,
	trans peer.SyncPeer,
	proxy proxy.AppProxy,
	selectorInitFunc SelectorCreationFn,
	selectorInitArgs SelectorCreationFnArgs,
	localAddr string) *Node {

	core := NewCore(id, key, participants, pst, conf.Logger)

	pubKey := core.HexID()

	if args, ok := selectorInitArgs.(SmartPeerSelectorCreationFnArgs); ok {
		args.GetFlagTable = core.poset.GetPeerFlagTableOfRandomUndeterminedEvent
		args.LocalAddr = localAddr
		selectorInitArgs = args
	}

	peerSelector := selectorInitFunc(participants, selectorInitArgs)

	node := Node{
		id:               id,
		conf:             conf,
		core:             core,
		logger:           conf.Logger.WithField("this_id", id),
		peerSelector:     peerSelector,
		trans:            trans,
		proxy:            proxy,
		submitCh:         proxy.SubmitCh(),
		submitInternalCh: proxy.SubmitInternalCh(),
		commitCh:         commitCh,
		shutdownCh:       make(chan struct{}),
		controlTimer:     NewRandomControlTimer(),
		start:            time.Now(),
		gossipJobs:       0,
		rpcJobs:          0,
		nodeState2:       newNodeState2(),
		signalTERMch:     make(chan os.Signal, 1),
	}

	stats := service.NewStats(pst.Store, pst, &node)
	node.stats = stats

	signal.Notify(node.signalTERMch, syscall.SIGTERM, os.Kill)

	node.logger.WithField("participants", participants).Debug("participants")
	node.logger.WithField("pubKey", pubKey).Debug("pubKey")

	node.needBootstrap = needBootstrap

	// Initialize
	node.setState(Gossiping)

	return &node
}

// Init initializes all the node processes
func (n *Node) Init() error {
	var peerAddresses []string
	for _, p := range n.peerSelector.Peers().ToPeerSlice() {
		peerAddresses = append(peerAddresses, p.NetAddr)
	}
	n.logger.WithField("peers", peerAddresses).Debug("Initialize Node")

	if n.needBootstrap {
		n.logger.Debug("Bootstrap")
		if err := n.core.Bootstrap(); err != nil {
			return err
		}
	}
	n.Register()

	return n.core.SetHeadAndHeight()
}

// RunAsync run the background processes asynchronously
func (n *Node) RunAsync(gossip bool) {
	n.logger.Debug("RunAsync(gossip bool)")
	go n.Run(gossip)
}

// Run core run loop, takes care of all processes
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
		n.logger.WithField("state", state.String()).Debug("Run(gossip bool)")

		switch state {
		case Gossiping:
			n.lachesis(gossip)
		case CatchingUp:
			if err := n.fastForward(); err != nil {
				n.logger.WithField("state", "fastForward").WithError(err).Debug("Run(gossip bool)")
			}
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
		// Slow gossip if nothing interesting to say
		if n.core.poset.GetPendingLoadedEvents() == 0 &&
			n.core.GetTransactionPoolCount() == 0 &&
			n.core.GetBlockSignaturePoolCount() == 0 {
			ts = time.Duration(time.Second)
		}
		n.controlTimer.resetCh <- ts
	}
}

// CommitCh return channel to write commit messages.
// Maybe better way to call n.commit in poset
// through node interface ???
func (n *Node) CommitCh() chan<- poset.Block {
	return n.commitCh
}

func (n *Node) doBackgroundWork() {
	for {
		select {
		case t := <-n.submitCh:
			n.logger.Debug("Adding Transactions to Transaction Pool")
			err := n.addTransaction(t)
			if err != nil {
				n.logger.Errorf("Adding Transactions to Transaction Pool: %s", err)
			}
			n.resetTimer()
		case t := <-n.submitInternalCh:
			n.logger.Debug("Adding Internal Transaction")
			n.addInternalTransaction(t)
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
		case <-n.signalTERMch:
			n.Shutdown()
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
		case rpc, ok := <-n.trans.ReceiverChannel():
			if !ok {
				return
			}
			n.goFunc(func() {
				n.rpcJobs.increment()
				n.logger.Debug("Processing RPC")
				n.processRPC(rpc)
				n.resetTimer()
				n.rpcJobs.decrement()
			})
		case <-n.controlTimer.tickCh:
			n.logStats()
			if gossip && n.gossipJobs.get() < 1 {
				n.goFunc(func() {
					n.gossipJobs.increment()
					if err := n.gossip(returnCh); err != nil {
						n.logger.WithError(err).Debug("node::lachesis(bool)::n.controlTimer.tickCh")
					}
					n.gossipJobs.decrement()
				})
				n.logger.Debug("Gossip")
			}
			n.resetTimer()
		case <-returnCh:
			return
		case <-n.shutdownCh:
			return
		}
	}
}

func (n *Node) processRPC(rpc *peer.RPC) {
	logger := n.logger.WithFields(logrus.Fields{"method": "processRPC",
		"cmd": rpc.Command})
	switch cmd := rpc.Command.(type) {
	case *peer.SyncRequest:
		n.processSyncRequest(rpc, cmd)
	case *peer.ForceSyncRequest:
		n.processEagerSyncRequest(rpc, cmd)
	case *peer.FastForwardRequest:
		n.processFastForwardRequest(rpc, cmd)
	default:
		logger.Warn("unexpected RPC command")
		// TODO: context.Background
		rpc.SendResult(context.Background(), n.logger,
			nil, fmt.Errorf("unexpected command"))
	}
}

func (n *Node) processSyncRequest(rpc *peer.RPC, cmd *peer.SyncRequest) {
	n.logger.WithFields(logrus.Fields{
		"from_id": cmd.FromID,
		"known":   cmd.Known,
	}).Debug("processSyncRequest(rpc net.RPC, cmd *net.SyncRequest)")

	resp := &peer.SyncResponse{
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

	// TODO: context.Background
	rpc.SendResult(context.Background(), n.logger, resp, respErr)
}

func (n *Node) processEagerSyncRequest(rpc *peer.RPC, cmd *peer.ForceSyncRequest) {
	success := true
	participants, err := n.core.poset.Store.Participants()
	if err != nil {
		n.logger.WithField("error", err).Error("n.sync(cmd.Events)")
		success = false
	}
	p, ok := participants.ReadByID(cmd.FromID)
	if !ok {
		n.logger.WithField("error", err).Error("n.sync(cmd.Events)")
		success = false
	}
	n.logger.WithFields(logrus.Fields{
		"from":    p.NetAddr,
		"from_id": cmd.FromID,
		"events":  len(cmd.Events),
	}).Debug("processEagerSyncRequest(rpc net.RPC, cmd *net.ForceSyncRequest)")

	resp := &peer.ForceSyncResponse{
		FromID:  n.id,
		Success: success,
	}
	// TODO: context.Background
	rpc.SendResult(context.Background(), n.logger, resp, nil)

	n.coreLock.Lock()
	err = n.sync(&p, cmd.Events)
	n.coreLock.Unlock()

	if err != nil {
		n.logger.WithField("error", err).Error("n.sync(cmd.Events)")
		success = false
	}
}

func (n *Node) processFastForwardRequest(rpc *peer.RPC, cmd *peer.FastForwardRequest) {
	n.logger.WithFields(logrus.Fields{
		"from": cmd.FromID,
	}).Debug("processFastForwardRequest(rpc net.RPC, cmd *net.FastForwardRequest)")

	resp := &peer.FastForwardResponse{
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
		resp.Block = block
		resp.Frame = frame

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
	// TODO: context.Background
	rpc.SendResult(context.Background(), n.logger, resp, respErr)
}

// This function is usually called in a go-routine and needs to inform the
// calling routine (usually the lachesis routine) when it is time to exit the
// Gossiping state and return.
func (n *Node) gossip(parentReturnCh chan struct{}) error {

	peer := n.peerSelector.Next()
	if peer == nil {
		return fmt.Errorf("can't select next peer")
	}

	// pull
	syncLimit, otherKnownEvents, err := n.pull(peer)
	if err != nil {
		return err
	}

	// check and handle syncLimit
	if syncLimit {
		n.logger.WithField("from", peer.NetAddr).Debug("SyncLimit")
		n.setState(CatchingUp)
		parentReturnCh <- struct{}{}
		return nil
	}

	// push
	err = n.push(peer.NetAddr, otherKnownEvents)
	if err != nil {
		return err
	}

	// update peer selector
	n.peerSelector.UpdateLast(peer.NetAddr)

	return nil
}

func (n *Node) pull(peer *peers.Peer) (syncLimit bool, otherKnownEvents map[uint64]int64, err error) {
	// Compute Known
	n.coreLock.Lock()
	knownEvents := n.core.KnownEvents()
	n.coreLock.Unlock()

	// Send SyncRequest
	start := time.Now()
	resp, err := n.requestSync(peer.NetAddr, knownEvents)
	elapsed := time.Since(start)
	n.logger.WithField("Duration", elapsed.Nanoseconds()).Debug("n.requestSync(peer.NetAddr, knownEvents)")
	// FIXIT: should we catch io.EOF error here and how we process it?
	// 	if err == io.EOF {
	// 		return false, nil, nil
	// 	}
	if err != nil {
		n.logger.WithField("Error", err).Error("n.requestSync(peer.NetAddr, knownEvents)")
		return resp.SyncLimit, nil, err
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
	err = n.sync(peer, resp.Events)
	n.coreLock.Unlock()
	if err != nil {
		n.logger.WithField("error", err).Error("n.sync(peer, resp.Events)")
		return false, nil, err
	}

	return false, resp.Known, nil
}

func (n *Node) push(peerAddr string, knownEvents map[uint64]int64) error {

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

		// Create and Send ForceSyncRequest
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
		}).Debug("ForceSyncResponse")
	}

	return nil
}

func (n *Node) fastForward() error {
	n.logger.Debug("fastForward()")

	// wait until sync routines finish
	n.waitRoutines()

	// fastForwardRequest
	peer := n.peerSelector.Next()
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
	err = n.core.FastForward(peer.PubKeyHex, resp.Block, resp.Frame)
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

func (n *Node) requestSync(target string, known map[uint64]int64) (*peer.SyncResponse, error) {
	args := &peer.SyncRequest{FromID: n.id, Known: known}
	out := &peer.SyncResponse{}
	err := n.trans.Sync(context.Background(), target, args, out)

	return out, err
}

func (n *Node) requestEagerSync(target string, events []poset.WireEvent) (*peer.ForceSyncResponse, error) {
	args := &peer.ForceSyncRequest{FromID: n.id, Events: events}
	out := &peer.ForceSyncResponse{}
	err := n.trans.ForceSync(context.Background(), target, args, out)

	return out, err
}

func (n *Node) requestFastForward(target string) (*peer.FastForwardResponse, error) {
	args := &peer.FastForwardRequest{FromID: n.id}
	out := &peer.FastForwardResponse{}
	err := n.trans.FastForward(context.Background(), target, args, out)

	return out, err
}

func (n *Node) sync(peer *peers.Peer, events []poset.WireEvent) error {
	// Insert Events in Poset and create new Head if necessary
	start := time.Now()
	err := n.core.Sync(peer, events)
	elapsed := time.Since(start)
	n.logger.WithField("Duration", elapsed.Nanoseconds()).Debug("n.core.Sync(events)")
	if err != nil {
		return fmt.Errorf("n.core.Sync(peer, events): %v", err)
	}

	if err := n.core.RunConsensus(); err != nil {
		return err
	}

	return nil
}

func (n *Node) commit(block poset.Block) error {

	n.coreLock.Lock()
	defer n.coreLock.Unlock()

	stateHash := []byte{0, 1, 2}
	_, err := n.proxy.CommitBlock(block)
	if err != nil {
		n.logger.WithError(err).Debug("commit(block poset.Block)")
	}

	n.logger.WithFields(logrus.Fields{
		"block":      block.Index(),
		"state_hash": fmt.Sprintf("%X", stateHash),
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

		block.StateHash = stateHash
		sig, err := n.core.SignBlock(block)
		if err != nil {
			return err
		}
		n.core.AddBlockSignature(sig)
	}

	return nil
}

func (n *Node) addTransaction(tx []byte) error {
	// we do not need coreLock here as n.core.AddTransactions has TransactionPoolLocker
	return n.core.AddTransactions([][]byte{tx})
}

func (n *Node) addInternalTransaction(tx poset.InternalTransaction) {
	n.coreLock.Lock()
	defer n.coreLock.Unlock()
	n.core.AddInternalTransactions([]poset.InternalTransaction{tx})
}

// Shutdown the node
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
		if err := n.core.poset.Store.Close(); err != nil {
			n.logger.WithError(err).Debug("node::Shutdown::n.core.poset.Store.Close()")
		}
	}
}

// GetStats returns processing stats for the node
func (n *Node) logStats() {
	stats := n.stats.Stats()
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
		"state":                  stats["state"],
		"z_gossipJobs":           n.gossipJobs.get(),
		"z_rpcJobs":              n.rpcJobs.get(),
		"pending_loaded_events":  n.GetPendingLoadedEvents(),
		"last_round":             n.stats.LastRound(),
		// "addr" is already defined in Node.logger, see NewNode() function
		// uncomment when needed
		//		"addr":                   n.localAddr,
		// "id" is duplicate of "this_id" in Node.logger, see NewNode() function
		// uncomment when needed
		//		"id":                     stats["id"],
	}).Warn("logStats()")
}

// SyncRate returns the current synchronization (talking to over nodes) rate in ms
func (n *Node) SyncRate() float64 {
	var syncErrorRate float64
	if n.syncRequests != 0 {
		syncErrorRate = float64(n.syncErrors) / float64(n.syncRequests)
	}
	return 1 - syncErrorRate
}

// EventDiff returns events that n knows about and are not in 'known'
func (n *Node) EventDiff(
	known map[uint64]int64) (events []poset.Event, err error) {
	return n.core.EventDiff(known)
}

// GetConsensusTransactionsCount get the count of finalized transactions
func (n *Node) GetConsensusTransactionsCount() uint64 {
	// TODO: remove use poset.
	return n.core.GetConsensusTransactionsCount()
}

// GetPendingLoadedEvents returns all the pending events
func (n *Node) GetPendingLoadedEvents() int64 {
	// TODO: remove use poset.
	return n.core.GetPendingLoadedEvents()
}

// GetTransactionPoolCount returns the count of all pending transactions
func (n *Node) GetTransactionPoolCount() int64 {
	return n.core.GetTransactionPoolCount()
}

// HeartbeatTimeout returns heartbeat timeout.
func (n *Node) HeartbeatTimeout() time.Duration {
	return n.conf.HeartbeatTimeout
}

// KnownEvents return known events.
func (n *Node) KnownEvents() map[uint64]int64 {
	return n.core.KnownEvents()
}

// StartTime returns node start time.
func (n *Node) StartTime() time.Time {
	return n.start
}

// State returns current state.
func (n *Node) State() string {
	return n.state.String()
}

// SyncLimit returns sync limit from config.
func (n *Node) SyncLimit() int64 {
	return n.conf.SyncLimit
}

// ID shows the ID of the node
func (n *Node) ID() uint64 {
	return n.id
}

// SubmitCh func for test
func (n *Node) SubmitCh(tx []byte) error {
	n.proxy.SubmitCh() <- []byte(tx)
	return nil
}

// Stop stops the node from gossiping
func (n *Node) Stop() {
	n.setState(Stop)
}
