package gossip

import (
	"math/big"

	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/sfctype"
)

func (s *Service) updateEpochStats(
	epoch idx.Epoch,
	block *inter.Block,
	blockFee *big.Int,
	sealEpoch bool,
) *sfctype.EpochStats {
	stats := s.store.GetDirtyEpochStats()
	stats.TotalFee = new(big.Int).Add(stats.TotalFee, blockFee)
	if sealEpoch {
		// dirty EpochStats becomes active
		stats.End = block.Time
		s.store.SetEpochStats(epoch, stats)

		// new dirty EpochStats
		s.store.SetDirtyEpochStats(&sfctype.EpochStats{
			Start:    block.Time,
			TotalFee: new(big.Int),
		})
	} else {
		s.store.SetDirtyEpochStats(stats)
	}

	return stats
}
