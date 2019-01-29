package poset2

import (
	"errors"
	"github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/poset2/model"
)

func (p *Poset2) Round(event *model.Event) (uint64, error) {
	// TODO: Add cache.

	return p.round2(event)
}

func (p *Poset2) round2(event *model.Event) (uint64, error) {
	if model.IsLeafEvent(p.codec, event) {
		return 0, nil
	}

	if (event.SelfParent() == common.Hash{}) {
		return 0, errors.New("self parent hash is empty")
	}

	sParent, err := p.GetEvent(event.SelfParent())
	if err != nil {
		return 0, err
	}

	parentRound, err := p.Round(sParent)
	if err != nil {
		return 0, err
	}

	if (event.OtherParent() != common.Hash{}) {
		oParent, err := p.GetEvent(event.OtherParent())
		if err != nil {
			return 0, err
		}

		oParentRound, err := p.Round(oParent)
		if err != nil {
			return 0, err
		}

		if oParentRound > parentRound {
			round, err := p.GetRound(oParentRound)
			if err != nil {
				return 0, err
			}

			var (
				roots = round.Roots()
				flags = event.Flags()
			)

			cc, somebody := p.Validator.CheckClothoCandidate(
				event, flags, roots)

			if cc {
				return oParentRound + 1, nil
			}

			if somebody {
				return oParentRound, nil
			}

			parentRound = oParentRound
		}
	}

	round, err := p.GetRound(parentRound)
	if err != nil {
		return 0, err
	}

	roundRoots := round.Roots()

	// check wp
	if cc, _ := p.Validator.CheckClothoCandidate(
		event, event.RootProof(), roundRoots); cc {
		return parentRound + 1, nil
	}

	// check ft
	if cc, _ := p.Validator.CheckClothoCandidate(
		event, event.Flags(), roundRoots); cc {
		return parentRound + 1, nil
	}

	return parentRound, nil
}
