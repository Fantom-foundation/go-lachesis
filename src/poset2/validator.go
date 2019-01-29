package poset2

import (
	"github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/poset2/model"
	"github.com/Fantom-foundation/go-lachesis/src/poset2/store"
)

type DefaultValidator struct {
	codec         model.Codec
	superMajority int
	store         *store.Store
}

func NewDefaultValidator(store *store.Store, codec model.Codec, partyLen int) *DefaultValidator {
	superMajority := 2*partyLen/3 + 1

	return &DefaultValidator{
		codec:         codec,
		superMajority: superMajority,
		store:         store,
	}
}

// TODO duplicate
func (v *DefaultValidator) GetEvent(hash common.Hash) (*model.Event, error) {
	data := v.store.GetEvent(hash)
	return model.DecodeEvent(v.codec, data)
}

func (v *DefaultValidator) CheckClothoCandidate(candidate *model.Event,
	proofEvents, roots []common.Hash) (decided bool, somebody bool) {
	if len(proofEvents) >= v.superMajority {
		hash := candidate.Hash(v.codec)
		if (hash == common.Hash{}) {
			return false, false
		}

		count := 0

		for k := range proofEvents {
			for _, root := range roots {
				if proofEvents[k] != root || root == hash ||
					proofEvents[k] == hash {
					continue
				}

				event, err := v.GetEvent(proofEvents[k])
				if err != nil {
					continue
				}

				dominate, err := v.Dominated(candidate, event)
				if err != nil {
					continue
				}

				if dominate {
					if !somebody {
						somebody = true
					}
					count++
				}
			}
		}

		if count >= v.superMajority {
			return true, somebody
		}
	}
	return false, somebody
}

// READ
func (v *DefaultValidator) Dominated(
	newestEvent, oldestEvent *model.Event) (bool, error) {
	// TODO: Add cache.

	return v.dominated2(newestEvent, oldestEvent)
}

// READ
func (v *DefaultValidator) dominated2(
	newestEvent, oldestEvent *model.Event) (bool, error) {
	if newestEvent.Hash(v.codec) == oldestEvent.Hash(v.codec) {
		return false, nil
	}

	if model.IsLeafEvent(v.codec, newestEvent) {
		return false, nil
	}

	selfParentHash := newestEvent.SelfParent()

	if (selfParentHash == common.Hash{}) {
		return false, nil
	}

	if oldestEvent.Hash(v.codec) == selfParentHash {
		return true, nil
	}

	selfParent, err := v.GetEvent(selfParentHash)
	if err != nil {
		// TODO process error
		return false, nil
	}

	found, err := v.Dominated(selfParent, oldestEvent)
	if err != nil {
		return false, err
	}

	if found {
		return true, err
	}

	otherParentHash := newestEvent.OtherParent()
	if (selfParentHash == common.Hash{}) {
		return false, err
	}

	if oldestEvent.Hash(v.codec) == otherParentHash {
		return true, nil
	}

	otherParent, err := v.GetEvent(otherParentHash)
	if err != nil {
		// TODO process error
		return false, nil
	}

	return v.Dominated(otherParent, oldestEvent)
}
