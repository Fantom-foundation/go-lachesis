package poset2

import (
	"github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/poset2/model"
	"github.com/pkg/errors"
)

func (p *Poset2) InsertEvent(event *model.Event) error {
	if !p.VerifyEvent(event) {
		return errors.New("invalid event")
	}

	if err := p.SetEvent(event); err != nil {
		return err
	}

	p.store.AddPendingEvent(event.Hash(p.codec))

	return nil
}

// READ
func (p *Poset2) isRoot(event *model.Event) (bool, error) {
	if model.IsLeafEvent(p.codec, event) {
		return true, nil
	}

	round, err := p.Round(event)
	if err != nil {
		return false, err
	}

	spEvent, err := p.GetEvent(event.SelfParent())
	if err != nil {
		return false, err
	}

	spRound, err := p.Round(spEvent)
	if err != nil {
		return false, err
	}

	return round > spRound, nil
}

func (p *Poset2) VerifyEvent(event *model.Event) bool {
	valid, err := event.VerifySignature(p.codec)
	if err != nil {
		return false
	}

	if !valid {
		return false
	}

	if withDuplicates(event.RootProof()) || withDuplicates(event.Flags()) ||
		withDuplicates(event.Parents()) {
		return false
	}

	// TODO
	if err := p.checkParents(event); err != nil {
		return false
	}

	return true
}

// TODO
func (p *Poset2) checkParents(event *model.Event) error {
	if model.IsLeafEvent(p.codec, event) {
		return nil
	}

	selfParentHash := event.SelfParent()
	lastEvent := p.LastEvent(event.Creator())
	if selfParentHash != lastEvent {
		return errors.New("self-parent not last known event by creator")
	}

	selfParent, err := p.GetEvent(selfParentHash)
	if err != nil {
		return err
	}

	if event.Index() != selfParent.Index()+1 {
		return errors.New("wrong index")
	}

	otherParent := event.SelfParent()
	if (otherParent != common.Hash{}) {

	}

	return nil
}
