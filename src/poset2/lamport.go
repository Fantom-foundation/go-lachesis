package poset2

import (
	"github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/poset2/model"
)

// READ
func (p *Poset2) calcLamport(event *model.Event) (uint64, error) {
	// TODO: Add cache.

	return p.calcLamport2(event)
}

// READ
func (p *Poset2) calcLamport2(event *model.Event) (uint64, error) {
	if model.IsLeafEvent(p.codec, event) {
		return 0, nil
	}

	selfParent, err := p.GetEvent(event.SelfParent())
	if err != nil {
		return 0, err
	}

	parentTimestamp, err := p.calcLamport(selfParent)
	if err != nil {
		return 0, err
	}

	if (event.OtherParent() != common.Hash{}) {
		otherParent, err := p.GetEvent(event.OtherParent())
		if err != nil {
			return 0, err
		}

		otherParentTimestamp, err := p.calcLamport(otherParent)
		if err != nil {
			return 0, err
		}

		if otherParentTimestamp > parentTimestamp {
			parentTimestamp = otherParentTimestamp
		}
	}

	return parentTimestamp + 1, nil
}
