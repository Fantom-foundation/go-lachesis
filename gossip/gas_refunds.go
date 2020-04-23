package gossip

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/Fantom-foundation/go-lachesis/app"
	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
)

const (
	minGasPowerRefund = 800
)

// incGasPowerRefund calculates the origination gas power refund
func (s *Service) incGasPowerRefund(epoch idx.Epoch, evmBlock *evmcore.EvmBlock, gasUsed []uint64, txPositions map[common.Hash]app.TxPosition, sealEpoch bool) {
	// Calc origination scores
	for i, tx := range evmBlock.Transactions {
		txHash := tx.Hash()
		txEventPos := txPositions[txHash]

		if tx.Gas() < gasUsed[i] {
			s.Log.Crit("Transaction gas used is higher than tx gas limit", "tx", txHash)
		}
		notUsedGas := tx.Gas() - gasUsed[i]
		if notUsedGas >= minGasPowerRefund { // do not refund if refunding is more costly than refunded value
			s.store.IncGasPowerRefund(epoch, txEventPos.Creator, notUsedGas)
		}
	}

	if sealEpoch {
		// prune not needed gas power records
		s.store.DelGasPowerRefunds(epoch - 1)
	}
}
