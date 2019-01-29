package poset2

import (
	"github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/poset2/model"
)

func (p *Poset2) GetEvent(hash common.Hash) (*model.Event, error) {
	data := p.store.GetEvent(hash)
	return model.DecodeEvent(p.codec, data)
}

func (p *Poset2) SetEvent(event *model.Event) error {
	data, err := event.Encode(p.codec)
	if err != nil {
		return err
	}

	hash := event.Hash(p.codec)
	p.store.SetEvent(hash, data)

	// TODO remove it
	p.store.SetLastEvent(event.Creator(), hash)
	return nil
}

func (p *Poset2) GetRound(index uint64) (*model.Round, error) {
	data, err := p.store.GetRound(index)
	if err != nil {
		return nil, err
	}

	return model.DecodeRound(p.codec, data)
}

func (p *Poset2) SetRound(index uint64, round *model.Round) error {
	data, err := round.Encode(p.codec)
	if err != nil {
		return err
	}

	p.store.SetRound(index, data)
	return nil
}

func (p *Poset2) LastEvent(publicKey []byte) common.Hash {
	return p.store.LastEvent(publicKey)

}
