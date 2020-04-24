package gossip

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	"github.com/Fantom-foundation/go-lachesis/app"
	"github.com/Fantom-foundation/go-lachesis/eventcheck"
	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/pos"
	"github.com/Fantom-foundation/go-lachesis/tracing"
)

// processEvent extends the engine.ProcessEvent with gossip-specific actions on each event processing
func (s *Service) processEvent(realEngine Consensus, e *inter.Event) error {
	// s.engineMu is locked here

	if s.store.HasEventHeader(e.Hash()) { // sanity check
		return eventcheck.ErrAlreadyConnectedEvent
	}

	// Trace arrival time of events
	if s.config.EventLocalTimeIndex {
		s.store.SetEventReceivingTime(e.Hash(), inter.Timestamp(time.Now().UnixNano()))
	}
	if s.config.DecisiveEventsIndex {
		s.currentEvent = e.Hash()
	}

	oldEpoch := e.Epoch

	s.store.SetEvent(e)
	if realEngine != nil {
		err := realEngine.ProcessEvent(e)
		if err != nil { // TODO make it possible to write only on success
			s.store.DeleteEvent(e.Epoch, e.Hash())
			return err
		}
	}
	_ = s.occurredTxs.CollectNotConfirmedTxs(e.Transactions)

	// set validator's last event. we don't care about forks, because this index is used only for emitter
	s.store.SetLastEvent(e.Epoch, e.Creator, e.Hash())

	// track events with no descendants, i.e. heads
	for _, parent := range e.Parents {
		s.store.DelHead(e.Epoch, parent)
	}
	s.store.AddHead(e.Epoch, e.Hash())

	s.packsOnNewEvent(e, e.Epoch)
	s.emitter.OnNewEvent(e)

	newEpoch := oldEpoch
	if realEngine != nil {
		newEpoch = realEngine.GetEpoch()
	}

	if newEpoch != oldEpoch {
		// notify event checkers about new validation data
		s.heavyCheckReader.Addrs.Store(ReadEpochPubKeys(s.abciApp, newEpoch))
		s.gasPowerCheckReader.Ctx.Store(ReadGasPowerContext(s.store, s.abciApp, s.engine.GetValidators(), newEpoch, &s.config.Net.Economy))

		// sealings/prunings
		s.packsOnNewEpoch(oldEpoch, newEpoch)
		s.store.delEpochStore(oldEpoch)
		s.store.getEpochStore(newEpoch)
		s.occurredTxs.Clear()

		// notify about new epoch after event connection
		s.emitter.OnNewEpoch(s.engine.GetValidators(), newEpoch)
		s.feed.newEpoch.Send(newEpoch)
	}

	immediately := (newEpoch != oldEpoch)
	return s.store.Commit(e.Hash().Bytes(), immediately)
}

// applyNewState moves the state according to new block (txs execution, SFC logic, epoch sealing)
func (s *Service) applyNewState(
	block *inter.Block,
	cheaters inter.Cheaters,
) (
	*inter.Block,
	map[common.Hash]app.TxPosition,
	common.Hash,
	bool,
) {
	// s.engineMu is locked here

	start := time.Now()

	events, allTxs := s.usedEvents(block)
	unusedCount := len(block.Events) - len(events)
	block.Events = block.Events[unusedCount:]

	// memorize position of each tx, for indexing and origination scores
	allTxsPosition := make(map[common.Hash]app.TxPosition)
	for _, e := range events {
		for i, tx := range e.Transactions {
			// if tx was met in multiple events, then assign to first ordered event
			if _, ok := allTxsPosition[tx.Hash()]; ok {
				continue
			}
			allTxsPosition[tx.Hash()] = app.TxPosition{
				Event:       e.Hash(),
				Creator:     e.Creator,
				EventOffset: uint32(i),
			}
		}
	}

	epoch := s.engine.GetEpoch()
	stateHash := s.store.GetBlock(block.Index - 1).Root
	evmHeader := evmcore.ToEvmHeader(block)

	s.abciApp.BeginBlock(beginBlockRequest(cheaters, stateHash), *evmHeader, s.blockParticipated)

	okTxs := make(types.Transactions, 0, len(allTxs))
	block.SkippedTxs = make([]uint, 0, len(allTxs))
	for i, tx := range allTxs {
		originator := allTxsPosition[tx.Hash()].Creator
		req := deliverTxRequest(tx, originator)
		resp := s.abciApp.DeliverTx(req)
		evmHeader.GasUsed += uint64(resp.GasUsed)

		if resp.Code != txIsFullyValid {
			block.SkippedTxs = append(block.SkippedTxs, uint(i))
			continue
		}

		okTxs = append(okTxs, tx)
		notUsedGas := resp.GasWanted - resp.GasUsed
		s.store.IncGasPowerRefund(epoch, originator, notUsedGas)

		if resp.Log != "" {
			s.Log.Info("tx processed", "log", resp.Log)
		}
	}

	_, sealEpoch := s.abciApp.EndBlock(endBlockRequest(block.Index))
	commit := s.abciApp.Commit()

	block.Root = common.BytesToHash(commit.Data)
	block.TxHash = types.DeriveSha(okTxs)
	block.GasUsed = evmHeader.GasUsed

	// memorize block position of each valid tx, for indexing and origination scores
	okTxsPosition := make(map[common.Hash]app.TxPosition, len(okTxs))
	for i, tx := range okTxs {
		// not skipped txs only
		position := allTxsPosition[tx.Hash()]
		position.Block = block.Index
		position.BlockOffset = uint32(i)
		okTxsPosition[tx.Hash()] = position
	}

	// process new epoch
	if sealEpoch {
		// prune not needed gas power records
		s.store.DelGasPowerRefunds(epoch - 1)
		s.onEpochSealed(block, cheaters)
	}

	// calc appHash
	appHash := block.TxHash

	log.Info("New block", "index", block.Index,
		"atropos", block.Atropos,
		"gasUsed", block.GasUsed,
		"skipped_txs", len(block.SkippedTxs),
		"txs", len(okTxs),
		"t", time.Since(start))

	return block, okTxsPosition, appHash, sealEpoch
}

// spillBlockEvents excludes first events which exceed BlockGasHardLimit
func (s *Service) spillBlockEvents(block *inter.Block) inter.Events {
	events := make(inter.Events, len(block.Events))
	if len(block.Events) == 0 {
		return events
	}
	gasPowerUsedSum := uint64(0)
	// iterate in reversed order
	for i := len(block.Events) - 1; i >= 0; i-- {
		id := block.Events[i]
		e := s.store.GetEvent(id)
		if e == nil {
			s.Log.Crit("Event not found", "event", id.String())
		}
		events[i] = e
		gasPowerUsedSum += e.GasPowerUsed
		// stop if limit is exceeded, erase [:i] events
		if gasPowerUsedSum > s.config.Net.Blocks.BlockGasHardLimit {
			// spill
			events = events[i+1:]
			break
		}
	}
	return events
}

// usedEvents and transactions of block for EVM
func (s *Service) usedEvents(block *inter.Block) (
	inter.Events, types.Transactions,
) {
	// s.engineMu is locked here
	if len(block.SkippedTxs) != 0 {
		log.Crit("Building with SkippedTxs isn't supported")
	}

	events := s.spillBlockEvents(block)
	transactions := make(types.Transactions, 0, len(events)*10)
	for _, e := range events {
		transactions = append(transactions, e.Transactions...)
	}

	return events, transactions
}

// onEpochSealed applies the new epoch sealing state
func (s *Service) onEpochSealed(block *inter.Block, cheaters inter.Cheaters) {
	// s.engineMu is locked here

	epoch := s.engine.GetEpoch()

	// delete last headers of cheaters
	for _, cheater := range cheaters {
		s.store.DelLastHeader(epoch, cheater) // for cheaters, it's uncertain which event is "last confirmed"
	}
	// prune not needed last headers
	s.store.DelLastHeaders(epoch - 1)
}

func (s *Service) legacyShouldSealEpoch(block *inter.Block, decidedFrame idx.Frame, cheaters inter.Cheaters) (sealEpoch bool) {
	// if cheater is confirmed, seal epoch right away to prune them from of BFT validators list
	epochStart := s.abciApp.GetDirtyEpochStats().Start
	sealEpoch = decidedFrame >= s.config.Net.Dag.MaxEpochBlocks
	sealEpoch = sealEpoch || block.Time-epochStart >= inter.Timestamp(s.config.Net.Dag.MaxEpochDuration)
	sealEpoch = sealEpoch || cheaters.Len() > 0
	return sealEpoch
}

// applyBlock execs ordered txns of new block on state, and fills the block DB indexes.
func (s *Service) applyBlock(block *inter.Block, decidedFrame idx.Frame, cheaters inter.Cheaters) (newAppHash common.Hash, sealEpoch bool) {
	// s.engineMu is locked here

	confirmBlocksMeter.Inc(1)

	block, txPositions, newAppHash, sealEpoch := s.applyNewState(block, cheaters)

	s.store.SetBlock(block)
	s.store.SetBlockIndex(block.Atropos, block.Index)

	// Build index for txs
	if s.config.TxIndex {
		for tx, pos := range txPositions {
			s.store.SetTxPosition(tx, &pos)
		}
	}

	// Trace by which event this block was confirmed (only for API)
	if s.config.DecisiveEventsIndex {
		s.store.SetBlockDecidedBy(block.Index, s.currentEvent)
	}

	// trace confirmed transactions
	confirmTxnsMeter.Inc(int64(len(txPositions)))
	for tx := range txPositions {
		tracing.FinishTx(tx, "Service.onNewBlock()")
		if latency, err := txLatency.Finish(tx); err == nil {
			txTtfMeter.Update(latency.Milliseconds())
		}
	}

	s.blockParticipated = make(map[idx.StakerID]bool) // reset map of participated validators

	return newAppHash, sealEpoch
}

// selectValidatorsGroup is a callback type to select new validators group
func (s *Service) selectValidatorsGroup(oldEpoch, newEpoch idx.Epoch) (newValidators *pos.Validators) {
	// s.engineMu is locked here

	builder := pos.NewBuilder()
	for _, it := range s.abciApp.GetEpochValidators(newEpoch) {
		builder.Set(it.StakerID, pos.BalanceToStake(it.Staker.CalcTotalStake()))
	}

	return builder.Build()
}

// onEventConfirmed is callback type to notify about event confirmation
func (s *Service) onEventConfirmed(header *inter.EventHeaderData, seqDepth idx.Event) {
	// s.engineMu is locked here

	if !header.NoTransactions() {
		// erase confirmed txs from originated-but-non-confirmed
		// to allow to re-originate this transaction if it will get skipped or spilled
		event := s.store.GetEvent(header.Hash())
		s.occurredTxs.CollectConfirmedTxs(event.Transactions)
	}

	// track last confirmed events from each validator
	if seqDepth == 0 {
		s.store.AddLastHeader(header.Epoch, header)
	}

	// track validators who participated in the block
	s.blockParticipated[header.Creator] = true
}

// isEventAllowedIntoBlock is callback type to check is event may be within block or not
func (s *Service) isEventAllowedIntoBlock(header *inter.EventHeaderData, seqDepth idx.Event) bool {
	// s.engineMu is locked here

	if header.NoTransactions() {
		return false // block contains only non-empty events to speed up block retrieving and processing
	}
	if seqDepth > s.config.Net.Dag.MaxValidatorEventsInBlock {
		return false // block contains only MaxValidatorEventsInBlock highest events from a creator to prevent huge blocks
	}
	return true
}
