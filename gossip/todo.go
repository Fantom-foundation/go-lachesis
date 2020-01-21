package gossip

import (
	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/pos"
	"github.com/Fantom-foundation/go-lachesis/inter/sfctype"
)

type appdeps struct {
	*Service
}

func (s *appdeps) GetBlock(n idx.Block) *inter.Block {
	return s.store.GetBlock(n)
}

func (s *appdeps) GetValidators() *pos.Validators {
	return s.engine.GetValidators()
}

func (s *appdeps) GetEventHeader(epoch idx.Epoch, h hash.Event) *inter.EventHeaderData {
	return s.store.GetEventHeader(epoch, h)
}

func (s *appdeps) BlockParticipated(staker idx.StakerID) bool {
	return s.blockParticipated[staker]
}

func (s *appdeps) GetDirtyEpochStats() *sfctype.EpochStats {
	return s.store.GetDirtyEpochStats()
}

func (s *appdeps) SetDirtyEpochStats(val *sfctype.EpochStats) {
	s.store.SetDirtyEpochStats(val)
}

func (s *appdeps) GetEpochStats(epoch idx.Epoch) *sfctype.EpochStats {
	return s.store.GetEpochStats(epoch)
}

func (s *appdeps) SetEpochStats(epoch idx.Epoch, val *sfctype.EpochStats) {
	s.store.SetEpochStats(epoch, val)
}

func (s *appdeps) GetDummyChainReader() evmcore.DummyChain {
	r := s.GetEvmStateReader()
	return evmcore.DummyChain(r)
}
