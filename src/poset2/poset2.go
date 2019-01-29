package poset2

import (
	"bytes"
	"github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/poset2/model"
	"github.com/Fantom-foundation/go-lachesis/src/poset2/store"
	"github.com/sirupsen/logrus"
)

type Core interface {
	Head() common.Hash
	PubKey() []byte
}

type Poset2 struct {
	codec  model.Codec
	core   Core
	logger *logrus.Entry
	store  *store.Store
	// TODO hide it
	Validator *DefaultValidator
}

func NewPoset2(store *store.Store,
	validator *DefaultValidator, logger *logrus.Entry) *Poset2 {
	logger = logger.WithFields(logrus.Fields{"type": "poset"})

	return &Poset2{
		codec:     model.NewDefaultCodec(),
		logger:    logger,
		store:     store,
		Validator: validator,
	}
}

// Write
func (p *Poset2) DivideRounds() error {
	for _, pendingEvent := range p.store.PendingEvents() {
		event, err := p.GetEvent(pendingEvent)
		if err != nil {
			return err
		}

		if event.Round() != 0 {
			continue
		}

		roundIndex, err := p.Round(event)
		if err != nil {
			return err
		}

		event.SetRound(roundIndex)

		root, err := p.isRoot(event)
		if err != nil {
			return err
		}

		var round *model.Round
		round, err = p.GetRound(roundIndex)
		if err != nil {
			if err != store.ErrNotFound {
				return err
			}
			round = model.NewRound()
		}

		round.AddEvent(pendingEvent, root)

		if root {
			if p.core != nil && event.Hash(p.codec) == p.core.Head() &&
				bytes.Equal(event.Creator(), p.core.PubKey()) {

				event.SetFlags(round.Roots())

				parentRound, err := p.GetRound(roundIndex - 1)
				if err != nil {
					return err
				}
				event.SetRootProof(parentRound.Roots())
			}
		}

		if event.Lamport() == 0 {
			lamport, err := p.calcLamport(event)
			if err != nil {
				return err
			}
			event.SetLamport(lamport)
		}

		if err := p.SetEvent(event); err != nil {
			return err
		}

		if err := p.SetRound(roundIndex, round); err != nil {
			return err
		}

		if !round.IsQueued() && roundIndex >= p.store.LastRound() {
			p.store.AddPendingRound(&store.PendingRound{Index: roundIndex})
		}
	}
	return nil
}

// Write
func (p *Poset2) DecideClothos() error {
	return nil
}

////func (p *Poset2) round(x EventHash) (int64, error) {
////
////}
////
////func (p *Poset2) round2(x EventHash) (int64, error) {
////
////}
////
////func (p *Poset2) clotho(x EventHash) (bool, error) {
////
////}
////
////func (p *Poset2) roundDiff(x, y EventHash) (int64, error) {
////
////}
////
////func (p *Poset2) checkSelfParent(event Event) error {
////
////}
////
////func (p *Poset2) checkOtherParent(event Event) error {
////
////}
////
////func (p *Poset2) createSelfParentRootEvent(ev Event) (RootEvent, error) {
////
////}
////
////func (p *Poset2) createOtherParentRootEvent(ev Event) (RootEvent, error) {
////
////
////}
////
////func (p *Poset2) createRoot(ev Event) (Root, error) {
////
////}
////
////func (p *Poset2) SetWireInfo(event *Event) error {
////}
////
////func (p *Poset2) SetWireInfoAndSign(event *Event, privKey *ecdsa.PrivateKey) error {
////
////}
////
////func (p *Poset2) setWireInfo(event *Event) error {
////}
////
////func (p *Poset2) updatePendingRounds(decidedRounds map[int64]int64) {
////}
////
////func (p *Poset2) removeProcessedSignatures(processedSignatures map[int64]bool) {
////}
////
////func (p *Poset2) InsertEvent(event Event, setWireInfo bool) error {
////
////}
////
////func (p *Poset2) DivideRounds() error {
////}
////
////func (p *Poset2) DecideAtropos2() error {
////
////}
////
////func (p *Poset2) DecideRoundReceived2() error {
////}
////
////func (p *Poset2) ProcessDecidedRounds() error {
////}
////
////func (p *Poset2) GetFrame(roundNumber int64) (Frame, error) {
////}
////
////func (p *Poset2) ProcessSigPool() error {
////
////}
////
////func (p *Poset2) GetAnchorBlockWithFrame() (Block, Frame, error) {
////
////}
////
////func (p *Poset2) Reset(block Block, frame Frame) error {
////}
////
////func (p *Poset2) Bootstrap() error {
////
////}
////
////func (p *Poset2) ReadWireInfo(wevent WireEvent) (*Event, error) {
////}
////
////func (p *Poset2) CheckBlock(block Block) error {
////}
////
////func (p *Poset2) setLastConsensusRound(i int64) {
////}
////
////func (p *Poset2) setAnchorBlock(i int64) {
////}
////
////func (p *Poset2) GetPeerFlagTableOfRandomUndeterminedEvent() (map[string]int64, error) {
////}
////
////func (p *Poset2) GetUndeterminedEvents() EventHashes {
////
////}
////
////func (p *Poset2) GetPendingLoadedEvents() int64 {
////}
////
////func (p *Poset2) GetLastConsensusRound() int64 {
////}
////
////func (p *Poset2) GetConsensusTransactionsCount() uint64 {
////}
