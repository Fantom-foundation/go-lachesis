package node

import (
	"crypto/ecdsa"
	"fmt"
	"reflect"
	"sort"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/golang/protobuf/proto"

	"github.com/Fantom-foundation/go-lachesis/src/crypto"
	"github.com/Fantom-foundation/go-lachesis/src/log"
	"github.com/Fantom-foundation/go-lachesis/src/peers"
	"github.com/Fantom-foundation/go-lachesis/src/poset"
	"github.com/Fantom-foundation/go-lachesis/src/utils"
)

type Core struct {
	id     int64
	key    *ecdsa.PrivateKey
	pubKey []byte
	hexID  string
	poset  *poset.Poset

	inDegrees map[string]uint64

	participants *peers.Peers // [PubKey] => id
	head         string
	Seq          int64

	transactionPool         [][]byte
	internalTransactionPool []poset.InternalTransaction
	blockSignaturePool      []poset.BlockSignature

	logger *logrus.Entry

	maxTransactionsInEvent int

	addSelfEventBlockLocker       sync.Mutex
	transactionPoolLocker         sync.RWMutex
	internalTransactionPoolLocker sync.RWMutex
	blockSignaturePoolLocker      sync.RWMutex
}

func NewCore(id int64, key *ecdsa.PrivateKey, participants *peers.Peers,
	store poset.Store, commitCh chan poset.Block, logger *logrus.Logger) *Core {

	if logger == nil {
		logger = logrus.New()
		logger.Level = logrus.DebugLevel
		lachesis_log.NewLocal(logger, logger.Level.String())
	}
	logEntry := logger.WithField("id", id)

	inDegrees := make(map[string]uint64)
	for pubKey := range participants.GetByPubKeys() {
		inDegrees[pubKey] = 0
	}

	p2 := poset.NewPoset(participants, store, commitCh, logEntry)
	core := &Core{
		id:                      id,
		key:                     key,
		poset:                   p2,
		inDegrees:               inDegrees,
		participants:            participants,
		transactionPool:         [][]byte{},
		internalTransactionPool: []poset.InternalTransaction{},
		blockSignaturePool:      []poset.BlockSignature{},
		logger:                  logEntry,
		head:                    "",
		Seq:                     -1,
		// MaxReceiveMessageSize limitation in grpc: https://github.com/grpc/grpc-go/blob/master/clientconn.go#L96
		// default value is 4 * 1024 * 1024 bytes
		// we use transactions of 120 bytes in tester, thus rounding it down to 16384
		maxTransactionsInEvent: 16384,
	}

	p2.SetCore(core)
	core.setRootEvents()

	return core
}

func (c *Core) ID() int64 {
	return c.id
}

func (c *Core) PubKey() []byte {
	if c.pubKey == nil {
		c.pubKey = crypto.FromECDSAPub(&c.key.PublicKey)
	}
	return c.pubKey
}

func (c *Core) HexID() string {
	if c.hexID == "" {
		pubKey := c.PubKey()
		c.hexID = fmt.Sprintf("0x%X", pubKey)
	}
	return c.hexID
}

func (c *Core) Head() string {
	return c.head
}

// Heights returns map with heights for each participants
func (c *Core) Heights() map[string]uint64 {
	heights := make(map[string]uint64)
	for pubKey := range c.participants.GetByPubKeys() {
		participantEvents, err := c.poset.Store.ParticipantEvents(pubKey, -1)
		if err == nil {
			heights[pubKey] = uint64(len(participantEvents))
		} else {
			heights[pubKey] = 0
		}
	}
	return heights
}

func (c *Core) InDegrees() map[string]uint64 {
	return c.inDegrees
}

func (c *Core) SetHeadAndSeq() error {

	var head string
	var seq int64

	last, isRoot, err := c.poset.Store.LastEventFrom(c.HexID())
	if err != nil {
		return err
	}

	if isRoot {
		root, err := c.poset.Store.GetRoot(c.HexID())
		if err != nil {
			return err
		}
		head = root.SelfParent.Hash
		seq = root.SelfParent.Index
	} else {
		lastEvent, err := c.GetEventBlock(last)
		if err != nil {
			return err
		}
		head = last
		seq = lastEvent.Index()
	}

	c.head = head
	c.Seq = seq

	c.logger.WithFields(logrus.Fields{
		"core.head": c.head,
		"core.Seq":  c.Seq,
		"is_root":   isRoot,
	}).Debugf("SetHeadAndSeq()")

	return nil
}

func (c *Core) Bootstrap() error {
	if err := c.poset.Bootstrap(); err != nil {
		return err
	}
	c.bootstrapInDegrees()
	return nil
}

func (c *Core) bootstrapInDegrees() {
	for pubKey := range c.participants.GetByPubKeys() {
		c.inDegrees[pubKey] = 0
		eventHash, _, err := c.poset.Store.LastEventFrom(pubKey)
		if err != nil {
			continue
		}
		for otherPubKey := range c.participants.GetByPubKeys() {
			if otherPubKey == pubKey {
				continue
			}
			events, err := c.poset.Store.ParticipantEvents(otherPubKey, -1)
			if err != nil {
				continue
			}
			for _, eh := range events {
				event, err := c.poset.Store.GetEventBlock(eh)
				if err != nil {
					continue
				}
				if event.OtherParent() == eventHash {
					c.inDegrees[pubKey]++
				}
			}
		}
	}
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func (c *Core) SignAndInsertSelfEvent(event poset.Event) error {
	if err := c.poset.SetWireInfoAndSign(&event, c.key); err != nil {
		return err
	}

	return c.InsertEvent(event, true)
}

func (c *Core) InsertEvent(event poset.Event, setWireInfo bool) error {

	c.logger.WithFields(logrus.Fields{
		"event":      event,
		"creator":    event.GetCreator(),
		"selfParent": event.SelfParent(),
		"index":      event.Index(),
		"hex":        event.Hex(),
	}).Debug("InsertEvent(event poset.Event, setWireInfo bool)")

	if err := c.poset.InsertEvent(event, setWireInfo); err != nil {
		return err
	}

	if event.GetCreator() == c.HexID() {
		c.head = event.Hex()
		c.Seq = event.Index()
	}

	c.inDegrees[event.GetCreator()] = 0

	if otherEvent, err := c.poset.Store.GetEventBlock(event.OtherParent()); err == nil {
		c.inDegrees[otherEvent.GetCreator()]++
	}
	return nil
}

func (c *Core) KnownEvents() map[int64]int64 {
	return c.poset.Store.KnownEvents()
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func (c *Core) SignBlock(block poset.Block) (poset.BlockSignature, error) {
	sig, err := block.Sign(c.key)
	if err != nil {
		return poset.BlockSignature{}, err
	}
	if err := block.SetSignature(sig); err != nil {
		return poset.BlockSignature{}, err
	}
	return sig, c.poset.Store.SetBlock(block)
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func (c *Core) OverSyncLimit(knownEvents map[int64]int64, syncLimit int64) bool {
	totUnknown := int64(0)
	myKnownEvents := c.KnownEvents()
	for i, li := range myKnownEvents {
		if li > knownEvents[i] {
			totUnknown += li - knownEvents[i]
		}
	}
	if totUnknown > syncLimit {
		return true
	}
	return false
}

func (c *Core) GetAnchorBlockWithFrame() (poset.Block, poset.Frame, error) {
	return c.poset.GetAnchorBlockWithFrame()
}

// returns events that c knows about and are not in 'known'
func (c *Core) EventDiff(known map[int64]int64) (events []poset.Event, err error) {
	var unknown []poset.Event
	// known represents the index of the last event known for every participant
	// compare this to our view of events and fill unknown with events that we know of
	// and the other doesn't
	for id, ct := range known {
		peer, ok := c.participants.GetById(id)
		if !ok {
			// unknown peer detected.
			// TODO: we should handle this nicely
			continue
		}
		// get participant Events with index > ct
		participantEvents, err := c.poset.Store.ParticipantEvents(peer.PubKeyHex, ct)
		if err != nil {
			return []poset.Event{}, err
		}
		for _, e := range participantEvents {
			ev, err := c.poset.Store.GetEventBlock(e)
			if err != nil {
				return []poset.Event{}, err
			}
			c.logger.WithFields(logrus.Fields{
				"event":      ev,
				"creator":    ev.GetCreator(),
				"selfParent": ev.SelfParent(),
				"index":      ev.Index(),
				"hex":        ev.Hex(),
			}).Debugf("Sending Unknown Event")
			unknown = append(unknown, ev)
		}
	}
	sort.Stable(poset.ByTopologicalOrder(unknown))

	return unknown, nil
}

func (c *Core) Sync(unknownEvents []poset.WireEvent) error {

	c.logger.WithFields(logrus.Fields{
		"unknown_events":              len(unknownEvents),
		"transaction_pool":            c.GetTransactionPoolCount(),
		"internal_transaction_pool":   c.GetInternalTransactionPoolCount(),
		"block_signature_pool":        c.GetBlockSignaturePoolCount(),
		"c.poset.PendingLoadedEvents": c.poset.GetPendingLoadedEvents(),
	}).Debug("Sync(unknownEventBlocks []poset.EventBlock)")

	myKnownEvents := c.KnownEvents()
	otherHead := ""
	// add unknown events
	for k, we := range unknownEvents {
		c.logger.WithFields(logrus.Fields{
			"unknown_events": we,
		}).Debug("unknownEvents")
		ev, err := c.poset.ReadWireInfo(we)
		if err != nil {
			c.logger.WithField("EventBlock", we).Errorf("c.poset.ReadEventBlockInfo(we)")
			return err

		}
		if ev.Index() > myKnownEvents[ev.CreatorID()] {
			ev.SetLamportTimestamp(poset.LamportTimestampNIL)
			ev.SetRound(poset.RoundNIL)
			ev.SetRoundReceived(poset.RoundNIL)
			if err := c.InsertEvent(*ev, false); err != nil {
				return err
			}
		}

		// assume last event corresponds to other-head
		if k == len(unknownEvents)-1 {
			otherHead = ev.Hex()
		}
	}

	// create new event with self head and other head only if there are pending
	// loaded events or the pools are not empty
	if c.poset.GetPendingLoadedEvents() > 0 ||
		c.GetTransactionPoolCount() > 0 ||
		c.GetInternalTransactionPoolCount() > 0 ||
		c.GetBlockSignaturePoolCount() > 0 {
		return c.AddSelfEventBlock(otherHead)
	}
	return nil
}

func (c *Core) FastForward(peer string, block poset.Block, frame poset.Frame) error {

	// Check Block Signatures
	err := c.poset.CheckBlock(block)
	if err != nil {
		return err
	}

	// Check Frame Hash
	frameHash, err := frame.Hash()
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(block.GetFrameHash(), frameHash) {
		return fmt.Errorf("invalid Frame Hash")
	}

	err = c.poset.Reset(block, frame)
	if err != nil {
		return err
	}

	err = c.SetHeadAndSeq()
	if err != nil {
		return err
	}

	err = c.RunConsensus()
	if err != nil {
		return err
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (c *Core) AddSelfEventBlock(otherHead string) error {

	c.addSelfEventBlockLocker.Lock()
	defer c.addSelfEventBlockLocker.Unlock()

	// Get flag tables from parents
	parentEvent, errSelf := c.poset.Store.GetEventBlock(c.head)
	if errSelf != nil {
		c.logger.Warnf("failed to get parent: %s", errSelf)
	}
	otherParentEvent, errOther := c.poset.Store.GetEventBlock(otherHead)
	if errOther != nil {
		c.logger.Warnf("failed to get other parent: %s", errOther)
	}

	var (
		flagTable map[string]int64
		err       error
	)

	if errSelf != nil {
		flagTable = map[string]int64{c.head: 1}
	} else {
		flagTable, err = parentEvent.GetFlagTable()
		if err != nil {
			return fmt.Errorf("failed to get self flag table: %s", err)
		}
	}

	if errOther == nil {
		flagTable, err = otherParentEvent.MergeFlagTable(flagTable)
		if err != nil {
			return fmt.Errorf("failed to marge flag tables: %s", err)
		}
	}

	// create new event with self head and empty other parent
	// empty transaction pool in its payload
	var batch [][]byte
	nTxs := min(int(c.GetTransactionPoolCount()), c.maxTransactionsInEvent)
	c.transactionPoolLocker.RLock()
	batch = c.transactionPool[0:nTxs:nTxs]
	c.transactionPoolLocker.RUnlock()
	newHead := poset.NewEvent(batch,
		c.internalTransactionPool,
		c.blockSignaturePool,
		[]string{c.head, otherHead}, c.PubKey(), c.Seq+1, flagTable)

	if err := c.SignAndInsertSelfEvent(newHead); err != nil {
		return fmt.Errorf("newHead := poset.NewEventBlock: %s", err)
	}
	c.logger.WithFields(logrus.Fields{
		"transactions":          c.GetTransactionPoolCount(),
		"internal_transactions": c.GetInternalTransactionPoolCount(),
		"block_signatures":      c.GetBlockSignaturePoolCount(),
	}).Debug("newHead := poset.NewEventBlock")

	c.transactionPoolLocker.Lock()
	c.transactionPool = c.transactionPool[nTxs:] //[][]byte{}
	c.transactionPoolLocker.Unlock()
	c.internalTransactionPoolLocker.Lock()
	c.internalTransactionPool = []poset.InternalTransaction{}
	c.internalTransactionPoolLocker.Unlock()
	// retain c.blockSignaturePool until c.transactionPool is empty
	// FIXIT: is there any better strategy?
	if c.GetTransactionPoolCount() == 0 {
		c.blockSignaturePoolLocker.Lock()
		c.blockSignaturePool = []poset.BlockSignature{}
		c.blockSignaturePoolLocker.Unlock()
	}

	return nil
}

func (c *Core) FromWire(wireEvents []poset.WireEvent) ([]poset.Event, error) {
	events := make([]poset.Event, len(wireEvents), len(wireEvents))
	for i, w := range wireEvents {
		ev, err := c.poset.ReadWireInfo(w)
		if err != nil {
			return nil, err
		}
		events[i] = *ev
	}
	return events, nil
}

func (c *Core) ToWire(events []poset.Event) ([]poset.WireEvent, error) {
	wireEvents := make([]poset.WireEvent, len(events), len(events))
	for i, e := range events {
		wireEvents[i] = e.ToWire()
	}
	return wireEvents, nil
}

func (c *Core) RunConsensus() error {
	start := time.Now()
	err := c.poset.DivideRounds()
	c.logger.WithField("Duration", time.Since(start).Nanoseconds()).Debug("c.poset.DivideAtropos()")
	if err != nil {
		c.logger.WithField("Error", err).Error("c.poset.DivideAtropos()")
		return err
	}

	start = time.Now()
	err = c.poset.DecideAtropos()
	c.logger.WithField("Duration", time.Since(start).Nanoseconds()).Debug("c.poset.DecideClotho()")
	if err != nil {
		c.logger.WithField("Error", err).Error("c.poset.DecideClotho()")
		return err
	}

	start = time.Now()
	err = c.poset.DecideRoundReceived()
	c.logger.WithField("Duration", time.Since(start).Nanoseconds()).Debug("c.poset.DecideAtroposRoundReceived()")
	if err != nil {
		c.logger.WithField("Error", err).Error("c.poset.DecideAtroposRoundReceived()")
		return err
	}

	start = time.Now()
	err = c.poset.ProcessDecidedRounds()
	c.logger.WithField("Duration", time.Since(start).Nanoseconds()).Debug("c.poset.ProcessAtroposRounds()")
	if err != nil {
		c.logger.WithField("Error", err).Error("c.poset.ProcessAtroposRounds()")
		return err
	}

	start = time.Now()
	err = c.poset.ProcessSigPool()
	c.logger.WithField("Duration", time.Since(start).Nanoseconds()).Debug("c.poset.ProcessSigPool()")
	if err != nil {
		c.logger.WithField("Error", err).Error("c.poset.ProcessSigPool()")
		return err
	}

	c.logger.WithFields(logrus.Fields{
		"transaction_pool":            c.GetTransactionPoolCount(),
		"block_signature_pool":        c.GetBlockSignaturePoolCount(),
		"c.poset.pendingLoadedEvents": c.poset.GetPendingLoadedEvents(),
	}).Debug("c.RunConsensus()")

	return nil
}

func (c *Core) AddTransactions(txs [][]byte) {
	c.transactionPoolLocker.Lock()
	defer c.transactionPoolLocker.Unlock()
	c.transactionPool = append(c.transactionPool, txs...)
}

func (c *Core) AddInternalTransactions(txs []poset.InternalTransaction) {
	c.internalTransactionPoolLocker.Lock()
	defer c.internalTransactionPoolLocker.Unlock()
	c.internalTransactionPool = append(c.internalTransactionPool, txs...)
}

func (c *Core) AddBlockSignature(bs poset.BlockSignature) {
	c.blockSignaturePoolLocker.Lock()
	defer c.blockSignaturePoolLocker.Unlock()
	c.blockSignaturePool = append(c.blockSignaturePool, bs)
}

func (c *Core) GetHead() (poset.Event, error) {
	return c.poset.Store.GetEventBlock(c.head)
}

func (c *Core) GetEventBlock(hash string) (poset.Event, error) {
	return c.poset.Store.GetEventBlock(hash)
}

func (c *Core) GetEventBlockTransactions(hash string) ([][]byte, error) {
	var txs [][]byte
	ex, err := c.GetEventBlock(hash)
	if err != nil {
		return txs, err
	}
	txs = ex.Transactions()
	return txs, nil
}

func (c *Core) GetConsensusEvents() []string {
	return c.poset.Store.ConsensusEvents()
}

func (c *Core) GetConsensusEventsCount() int64 {
	return c.poset.Store.ConsensusEventsCount()
}

func (c *Core) GetUndeterminedEvents() []string {
	return c.poset.GetUndeterminedEvents()
}

func (c *Core) GetPendingLoadedEvents() int64 {
	return c.poset.GetPendingLoadedEvents()
}

func (c *Core) GetConsensusTransactions() ([][]byte, error) {
	var txs [][]byte
	for _, e := range c.GetConsensusEvents() {
		eTxs, err := c.GetEventBlockTransactions(e)
		if err != nil {
			return txs, fmt.Errorf("GetConsensusTransactions(): %s", e)
		}
		txs = append(txs, eTxs...)
	}
	return txs, nil
}

func (c *Core) GetLastConsensusRound() int64 {
	return c.poset.GetLastConsensusRound()
}

func (c *Core) GetConsensusTransactionsCount() uint64 {
	return c.poset.GetConsensusTransactionsCount()
}

func (c *Core) GetLastCommittedRoundEventsCount() int {
	return c.poset.LastCommitedRoundEvents
}

func (c *Core) GetLastBlockIndex() int64 {
	return c.poset.Store.LastBlockIndex()
}

func (c *Core) GetTransactionPoolCount() int64 {
	c.transactionPoolLocker.RLock()
	defer c.transactionPoolLocker.RUnlock()
	return int64(len(c.transactionPool))
}

func (c *Core) GetInternalTransactionPoolCount() int64 {
	c.internalTransactionPoolLocker.RLock()
	defer c.internalTransactionPoolLocker.RUnlock()
	return int64(len(c.internalTransactionPool))
}

func (c *Core) GetBlockSignaturePoolCount() int64 {
	c.blockSignaturePoolLocker.RLock()
	defer c.blockSignaturePoolLocker.RUnlock()
	return int64(len(c.blockSignaturePool))
}

func (c *Core) setRootEvents() error {
	roots, err := c.poset.Store.RootsBySelfParent()
	if err != nil {
		return err
	}
	for participant, root := range roots {
		var creator []byte
		fmt.Sscanf(participant, "0x%X", &creator)
		flagTable := map[string]int64{root.SelfParent.Hash: 1}
		ft, _ := proto.Marshal(&poset.FlagTableWrapper{Body: flagTable})
		body := poset.EventBody{
			Creator: creator, /*s.participants.ByPubKey[participant].PubKey,*/
			Index:   root.SelfParent.Index,
			Parents: []string{root.SelfParent.Hash, ""}, //root.SelfParent.Hash, root.SelfParent.Hash},
		}
		event := poset.Event{
			Message: &poset.EventMessage{
				Hash:             utils.HashFromHex(root.SelfParent.Hash),
				CreatorID:        root.SelfParent.CreatorID,
				TopologicalIndex: -1,
				Body:             &body,
				FlagTable:        ft,
				ClothoProof:      []string{root.SelfParent.Hash},
			},
//			lamportTimestamp: 0,
//			round:            0,
//			roundReceived:    0, /*RoundNIL*/
		}
		if _, err := c.poset.Store.GetEventBlock(event.Hex()); err != nil {
			// if we do not get root event from poset store we must be on a newly created database
			// so let create one
			event.Sign(c.key)
			if err := c.poset.Store.SetEvent(event); err != nil {
				return err
			}
		}
	}
	return nil
}
