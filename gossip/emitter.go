package gossip

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/hashicorp/golang-lru"

	"github.com/Fantom-foundation/go-lachesis/common/bigendian"
	"github.com/Fantom-foundation/go-lachesis/eventcheck"
	"github.com/Fantom-foundation/go-lachesis/eventcheck/basiccheck"
	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/gossip/occuredtxs"
	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/ancestor"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/pos"
	"github.com/Fantom-foundation/go-lachesis/lachesis"
	"github.com/Fantom-foundation/go-lachesis/lachesis/params"
	"github.com/Fantom-foundation/go-lachesis/logger"
	"github.com/Fantom-foundation/go-lachesis/tracing"
	"github.com/Fantom-foundation/go-lachesis/utils"
)

const (
	MimetypeEvent    = "application/event"
	TxTimeBufferSize = 20000
	TxTurnPeriod     = 4 * time.Second
	TxTurnNonces     = 8
)

// EmitterWorld is emitter's external world
type EmitterWorld struct {
	Store       *Store
	Engine      Consensus
	EngineMu    *sync.RWMutex
	Txpool      txPool
	Am          *accounts.Manager
	OccurredTxs *occuredtxs.Buffer

	Checkers *eventcheck.Checkers

	OnEmitted func(e *inter.Event)
	IsSynced  func() bool
	PeersNum  func() int
}

type Emitter struct {
	txTime *lru.Cache // tx hash -> tx time

	net    *lachesis.Config
	config *EmitterConfig

	world EmitterWorld

	creator   common.Address
	creatorMu sync.RWMutex

	syncStatus selfForkProtection

	gasRate         metrics.Meter
	prevEmittedTime time.Time

	done chan struct{}
	wg   sync.WaitGroup

	logger.Periodic
}

type selfForkProtection struct {
	connectedTime           time.Time
	prevLocalEmittedID      hash.Event
	prevExternalEmittedTime time.Time
	becameValidatorTime     time.Time
	wasValidator            bool
}

// NewEmitter creation.
func NewEmitter(
	config *Config,
	world EmitterWorld,
) *Emitter {

	txTime, _ := lru.New(TxTimeBufferSize)
	loggerInstance := logger.MakeInstance()
	return &Emitter{
		net:      &config.Net,
		config:   &config.Emitter,
		world:    world,
		creator:  config.Emitter.Validator,
		gasRate:  metrics.NewMeterForced(),
		txTime:   txTime,
		Periodic: logger.Periodic{Instance: loggerInstance},
	}
}

// StartEventEmission starts event emission.
func (em *Emitter) StartEventEmission() {
	if em.done != nil {
		return
	}
	em.done = make(chan struct{})

	newTxsCh := make(chan evmcore.NewTxsNotify)
	em.world.Txpool.SubscribeNewTxsNotify(newTxsCh)

	em.prevEmittedTime = em.loadPrevEmitTime()
	em.syncStatus.connectedTime = time.Now()
	em.syncStatus.wasValidator = true

	done := em.done
	em.wg.Add(1)
	go func() {
		defer em.wg.Done()
		ticker := time.NewTicker(em.config.MinEmitInterval / 5)
		for {
			select {
			case txNotify := <-newTxsCh:
				em.memorizeTxTimes(txNotify.Txs)
			case <-ticker.C:
				// must pass at least MinEmitInterval since last event
				if time.Since(em.prevEmittedTime) >= em.config.MinEmitInterval {
					em.EmitEvent()
				}
			case <-done:
				return
			}
		}
	}()
}

// StopEventEmission stops event emission.
func (em *Emitter) StopEventEmission() {
	if em.done == nil {
		return
	}

	close(em.done)
	em.done = nil
	em.wg.Wait()
}

// SetValidator sets event creator.
func (em *Emitter) SetValidator(addr common.Address) {
	em.creatorMu.Lock()
	defer em.creatorMu.Unlock()
	em.creator = addr
}

// GetValidator gets event creator.
func (em *Emitter) GetValidator() common.Address {
	em.creatorMu.RLock()
	defer em.creatorMu.RUnlock()
	return em.creator
}

func (em *Emitter) loadPrevEmitTime() time.Time {
	myID, ok := em.myStakerID()
	if !ok {
		return em.prevEmittedTime
	}

	prevEventID := em.world.Store.GetLastEvent(em.world.Engine.GetEpoch(), myID)
	if prevEventID == nil {
		return em.prevEmittedTime
	}
	prevEvent := em.world.Store.GetEventHeader(prevEventID.Epoch(), *prevEventID)
	if prevEvent == nil {
		return em.prevEmittedTime
	}
	return prevEvent.ClaimedTime.Time()
}

// safe for concurrent use
func (em *Emitter) memorizeTxTimes(txs types.Transactions) {
	now := time.Now()
	for _, tx := range txs {
		_, ok := em.txTime.Get(tx.Hash())
		if !ok {
			em.txTime.Add(tx.Hash(), now)
		}
	}
}

func (em *Emitter) myStakerID() (idx.StakerID, bool) {
	coinbase := em.GetValidator()

	validators := em.world.Store.app.GetEpochValidators(em.world.Engine.GetEpoch())
	for _, it := range validators {
		if it.Staker.Address == coinbase {
			return it.StakerID, true
		}
	}
	return 0, false
}

// safe for concurrent use
func (em *Emitter) isMyTxTurn(txHash common.Hash, sender common.Address, accountNonce uint64, now time.Time, validatorsArr []idx.StakerID, validatorsArrStakes []pos.Stake, me idx.StakerID) bool {
	turnHash := hash.Of(sender.Bytes(), bigendian.Int64ToBytes(accountNonce/TxTurnNonces), em.world.Engine.GetEpoch().Bytes())

	var txTime time.Time
	txTimeI, ok := em.txTime.Get(txHash)
	if !ok {
		txTime = now
		em.txTime.Add(txHash, txTime)
	} else {
		txTime = txTimeI.(time.Time)
	}

	roundIndex := int((now.Sub(txTime) / TxTurnPeriod) % time.Duration(len(validatorsArr)))
	turns := utils.WeightedPermutation(roundIndex+1, validatorsArrStakes, turnHash)

	return validatorsArr[turns[roundIndex]] == me
}

func (em *Emitter) addTxs(e *inter.Event, poolTxs map[common.Address]types.Transactions) *inter.Event {
	if poolTxs == nil || len(poolTxs) == 0 {
		return e
	}

	maxGasUsed := em.maxGasPowerToUse(e)

	now := time.Now()
	validators := em.world.Engine.GetValidators()
	validatorsArr := validators.SortedIDs() // validators must be sorted deterministically
	validatorsArrStakes := make([]pos.Stake, len(validatorsArr))
	for i, addr := range validatorsArr {
		validatorsArrStakes[i] = validators.Get(addr)
	}

	for sender, txs := range poolTxs {
		if txs.Len() > em.config.MaxTxsFromSender { // no more than MaxTxsFromSender txs from 1 sender
			txs = txs[:em.config.MaxTxsFromSender]
		}

		// txs is the chain of dependent txs
		for _, tx := range txs {
			// enough gas power
			if tx.Gas() >= e.GasPowerLeft.Min() || e.GasPowerUsed+tx.Gas() >= maxGasUsed {
				break // txs are dependent, so break the loop
			}
			// check not conflicted with already included txs (in any connected event)
			if em.world.OccurredTxs.MayBeConflicted(sender, tx.Hash()) {
				break // txs are dependent, so break the loop
			}
			// my turn, i.e. try to not include the same tx simultaneously by different validators
			if !em.isMyTxTurn(tx.Hash(), sender, tx.Nonce(), now, validatorsArr, validatorsArrStakes, e.Creator) {
				break // txs are dependent, so break the loop
			}

			// add
			e.GasPowerUsed += tx.Gas()
			e.GasPowerLeft.Sub(tx.Gas())
			e.Transactions = append(e.Transactions, tx)
		}
	}
	return e
}

func (em *Emitter) findBestParents(epoch idx.Epoch, myStakerID idx.StakerID) (*hash.Event, hash.Events, bool) {
	selfParent := em.world.Store.GetLastEvent(epoch, myStakerID)
	heads := em.world.Store.GetHeads(epoch) // events with no descendants

	var strategy ancestor.SearchStrategy
	vecClock := em.world.Engine.GetVectorIndex()
	if vecClock != nil {
		strategy = ancestor.NewCasualityStrategy(vecClock, em.world.Engine.GetValidators())
		if rand.Intn(20) == 0 { // every 20th event uses random strategy is avoid repeating patterns in DAG
			strategy = ancestor.NewRandomStrategy(rand.New(rand.NewSource(time.Now().UnixNano())))
		}

		// don't link to known cheaters
		heads = vecClock.NoCheaters(selfParent, heads)
		if selfParent != nil && len(vecClock.NoCheaters(selfParent, hash.Events{*selfParent})) == 0 {
			em.Periodic.Error(5*time.Second, "I've created a fork, events emitting isn't allowed", "creator", myStakerID)
			return nil, nil, false
		}
	} else {
		// use dummy strategy in engine-less tests
		strategy = ancestor.NewRandomStrategy(nil)
	}

	maxParents := em.config.MaxParents
	if maxParents < em.net.Dag.MaxFreeParents {
		maxParents = em.net.Dag.MaxFreeParents
	}
	if maxParents > em.net.Dag.MaxParents {
		maxParents = em.net.Dag.MaxParents
	}
	_, parents := ancestor.FindBestParents(maxParents, heads, selfParent, strategy)
	return selfParent, parents, true
}

// createEvent is not safe for concurrent use.
func (em *Emitter) createEvent(poolTxs map[common.Address]types.Transactions) *inter.Event {
	myStakerID, ok := em.myStakerID()
	if !ok {
		// not a validator
		return nil
	}
	validators := em.world.Engine.GetValidators()

	if synced, _, _ := em.logging(em.isSynced()); !synced {
		// I'm reindexing my old events, so don't create events until connect all the existing self-events
		return nil
	}

	var (
		epoch          = em.world.Engine.GetEpoch()
		selfParentSeq  idx.Event
		selfParentTime inter.Timestamp
		parents        hash.Events
		maxLamport     idx.Lamport
	)

	// Find parents
	selfParent, parents, ok := em.findBestParents(epoch, myStakerID)
	if !ok {
		return nil
	}

	// Set parent-dependent fields
	parentHeaders := make([]*inter.EventHeaderData, len(parents))
	for i, p := range parents {
		parent := em.world.Store.GetEventHeader(epoch, p)
		if parent == nil {
			em.Log.Crit("Emitter: head not found", "event", p.String())
		}
		parentHeaders[i] = parent
		if parentHeaders[i].Creator == myStakerID && i != 0 {
			// there're 2 heads from me, i.e. due to a fork, findBestParents could have found multiple self-parents
			em.Periodic.Error(5*time.Second, "I've created a fork, events emitting isn't allowed", "creator", myStakerID)
			return nil
		}
		maxLamport = idx.MaxLamport(maxLamport, parent.Lamport)
	}

	selfParentSeq = 0
	selfParentTime = 0
	var selfParentHeader *inter.EventHeaderData
	if selfParent != nil {
		selfParentHeader = parentHeaders[0]
		selfParentSeq = selfParentHeader.Seq
		selfParentTime = selfParentHeader.ClaimedTime
	}

	event := inter.NewEvent()
	event.Epoch = epoch
	event.Seq = selfParentSeq + 1
	event.Creator = myStakerID

	event.Parents = parents
	event.Lamport = maxLamport + 1
	event.ClaimedTime = inter.MaxTimestamp(inter.Timestamp(time.Now().UnixNano()), selfParentTime+1)

	// set consensus fields
	event = em.world.Engine.Prepare(event)
	if event == nil {
		em.Log.Warn("Dropped event while emitting")
		return nil
	}

	// calc initial GasPower
	event.GasPowerUsed = basiccheck.CalcGasPowerUsed(event, &em.net.Dag)
	availableGasPower, err := em.world.Checkers.Gaspowercheck.CalcGasPower(&event.EventHeaderData, selfParentHeader)
	if err != nil {
		em.Log.Warn("Gas power calculation failed", "err", err)
		return nil
	}
	if event.GasPowerUsed > availableGasPower.Min() {
		em.Periodic.Warn(time.Second, "Not enough gas power to emit event. Too small stake?",
			"gasPower", availableGasPower,
			"stake%", 100*float64(validators.Get(myStakerID))/float64(validators.TotalStake()))
		return nil
	}
	event.GasPowerLeft = *availableGasPower.Sub(event.GasPowerUsed)

	// Add txs
	event = em.addTxs(event, poolTxs)

	if !em.isAllowedToEmit(event, selfParentHeader) {
		return nil
	}

	// calc Merkle root
	event.TxHash = types.DeriveSha(event.Transactions)

	// sign
	coinbase := em.GetValidator()
	signer := func(data []byte) (sig []byte, err error) {
		acc := accounts.Account{
			Address: coinbase,
		}
		w, err := em.world.Am.Find(acc)
		if err != nil {
			return
		}
		return w.SignData(acc, MimetypeEvent, data)
	}
	if err := event.Sign(signer); err != nil {
		em.Periodic.Error(time.Second, "Failed to sign event. Please unlock account.", "err", err)
		return nil
	}
	// calc hash after event is fully built
	event.RecacheHash()
	event.RecacheSize()
	{
		// sanity check
		if em.world.Checkers != nil {
			if err := em.world.Checkers.Validate(event, parentHeaders); err != nil {
				em.Periodic.Error(time.Second, "Signed event incorrectly", "err", err)
				return nil
			}
		}
	}

	// set event name for debug
	em.nameEventForDebug(event)

	return event
}

// OnNewEvent tracks new events to find out am I properly synced or not
func (em *Emitter) OnNewEvent(e *inter.Event) {
	now := time.Now()
	myStakerID, isValidator := em.myStakerID()
	if isValidator && em.syncStatus.prevLocalEmittedID != e.Hash() {
		if e.Creator == myStakerID {
			// it was emitted by another instance with the same address
			em.syncStatus.prevExternalEmittedTime = now
		}
	}

	// track when I've became validator
	if isValidator && !em.syncStatus.wasValidator {
		em.syncStatus.becameValidatorTime = now
	}
	em.syncStatus.wasValidator = isValidator
}

func (em *Emitter) isSynced() (bool, string, time.Duration) {
	if em.config.SelfForkProtectionInterval == 0 {
		return true, "", 0 // protection disabled
	}
	if em.world.PeersNum() == 0 {
		em.syncStatus.connectedTime = time.Now() // move time of the first connection
		return false, "no connections", 0
	}
	if !em.world.IsSynced() {
		return false, "synchronizing (all the peers have higher/lower epoch)", 0
	}
	sinceLastExternalEvent := time.Since(em.syncStatus.prevExternalEmittedTime)
	if sinceLastExternalEvent < em.config.SelfForkProtectionInterval {
		return false, "synchronizing (not downloaded all the self-events)", em.config.SelfForkProtectionInterval - sinceLastExternalEvent
	}
	connectedTime := time.Since(em.syncStatus.connectedTime)
	if connectedTime < em.config.SelfForkProtectionInterval {
		return false, "synchronizing (just connected)", em.config.SelfForkProtectionInterval - connectedTime
	}
	sinceBecameValidator := time.Since(em.syncStatus.becameValidatorTime)
	if sinceBecameValidator < em.config.SelfForkProtectionInterval {
		return false, "synchronizing (just joined the validators group)", em.config.SelfForkProtectionInterval - sinceBecameValidator
	}

	return true, "", 0
}

func (em *Emitter) logging(synced bool, reason string, wait time.Duration) (bool, string, time.Duration) {
	if !synced {
		if wait == 0 {
			em.Periodic.Info(25*time.Second, "Emitting is paused", "reason", reason)
		} else {
			em.Periodic.Info(25*time.Second, "Emitting is paused", "reason", reason, "wait", wait)
		}
	}
	return synced, reason, wait
}

// return true if event is in epoch tail (unlikely to confirm)
func (em *Emitter) isEpochTail(e *inter.Event) bool {
	return e.Frame >= em.net.Dag.MaxEpochBlocks-em.config.EpochTailLength
}

func (em *Emitter) maxGasPowerToUse(e *inter.Event) uint64 {
	// No txs in epoch tail, because tail events are unlikely to confirm
	{
		if em.isEpochTail(e) {
			return 0
		}
	}
	// No txs if power is low
	{
		threshold := em.config.NoTxsThreshold
		if e.GasPowerLeft.Min() <= threshold {
			return 0
		}
		if e.GasPowerLeft.Min() < threshold+params.MaxGasPowerUsed {
			return e.GasPowerLeft.Min() - threshold
		}
	}
	// Smooth TPS if power isn't big
	{
		threshold := em.config.SmoothTpsThreshold
		if e.GasPowerLeft.Min() <= threshold {
			// it's emitter, so no need in determinism => fine to use float
			passedTime := float64(e.ClaimedTime.Time().Sub(em.prevEmittedTime)) / (float64(time.Second))
			maxGasUsed := uint64(passedTime * em.gasRate.Rate1() * em.config.MaxGasRateGrowthFactor)
			if maxGasUsed > params.MaxGasPowerUsed {
				maxGasUsed = params.MaxGasPowerUsed
			}
			return maxGasUsed
		}
	}
	return params.MaxGasPowerUsed
}

func (em *Emitter) isAllowedToEmit(e *inter.Event, selfParent *inter.EventHeaderData) bool {
	passedTime := e.ClaimedTime.Time().Sub(em.prevEmittedTime)
	// Slow down emitting if power is low
	{
		threshold := em.config.NoTxsThreshold
		if e.GasPowerLeft.Min() <= threshold {
			// it's emitter, so no need in determinism => fine to use float
			minT := float64(em.config.MinEmitInterval)
			maxT := float64(em.config.MaxEmitInterval)
			factor := float64(e.GasPowerLeft.Min()) / float64(threshold)
			adjustedEmitInterval := time.Duration(maxT - (maxT-minT)*factor)
			if passedTime < adjustedEmitInterval {
				return false
			}
		}
	}
	// Forbid emitting if not enough power and power is decreasing
	{
		threshold := em.config.EmergencyThreshold
		if e.GasPowerLeft.Min() <= threshold {
			if selfParent != nil && e.GasPowerLeft.Min() < selfParent.GasPowerLeft.Min() {
				validators := em.world.Engine.GetValidators()
				em.Periodic.Warn(10*time.Second, "Not enough power to emit event, waiting",
					"power", e.GasPowerLeft.String(),
					"selfParentPower", selfParent.GasPowerLeft.String(),
					"stake%", 100*float64(validators.Get(e.Creator))/float64(validators.TotalStake()))
				return false
			}
		}
	}
	// Slow down emitting if no txs to confirm/post, and not at epoch tail
	{
		if passedTime < em.config.MaxEmitInterval &&
			em.world.OccurredTxs.Len() == 0 &&
			len(e.Transactions) == 0 &&
			!em.isEpochTail(e) {
			return false
		}
	}

	return true
}

func (em *Emitter) EmitEvent() *inter.Event {
	poolTxs, err := em.world.Txpool.Pending() // request txs before locking engineMu to prevent deadlock!
	if err != nil {
		em.Log.Error("Tx pool transactions fetching error", "err", err)
		return nil
	}

	for _, tt := range poolTxs {
		for _, t := range tt {
			span := tracing.CheckTx(t.Hash(), "Emitter.EmitEvent(candidate)")
			defer span.Finish()
		}
	}

	em.world.EngineMu.Lock()
	defer em.world.EngineMu.Unlock()

	e := em.createEvent(poolTxs)
	if e == nil {
		return nil
	}
	em.syncStatus.prevLocalEmittedID = e.Hash()

	if em.world.OnEmitted != nil {
		em.world.OnEmitted(e)
	}
	em.gasRate.Mark(int64(e.GasPowerUsed))
	em.prevEmittedTime = time.Now() // record time after connecting, to add the event processing time
	em.Log.Info("New event emitted", "event", e.String(), "txs", e.Transactions.Len(), "elapsed", time.Since(e.ClaimedTime.Time()))

	// metrics
	for _, t := range e.Transactions {
		span := tracing.CheckTx(t.Hash(), "Emitter.EmitEvent()")
		defer span.Finish()
	}

	return e
}

func (em *Emitter) nameEventForDebug(e *inter.Event) {
	name := []rune(hash.GetNodeName(e.Creator))
	if len(name) < 1 {
		return
	}

	name = name[len(name)-1:]
	hash.SetEventName(e.Hash(), fmt.Sprintf("%s%03d",
		strings.ToLower(string(name)),
		e.Seq))
}
