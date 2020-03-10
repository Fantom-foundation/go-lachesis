package app

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
)

// TxPosition includes block and event.
type TxPosition struct {
	Block       idx.Block
	Event       hash.Event
	Creator     idx.StakerID
	EventOffset uint32
	BlockOffset uint32
}

// BlocksMissed is information about missed blocks from a staker.
type BlocksMissed struct {
	Num    idx.Block
	Period inter.Timestamp
}

// updateOriginationScores calculates the origination scores
func (a *App) updateOriginationScores(epoch idx.Epoch, evmBlock *evmcore.EvmBlock, receipts types.Receipts, txPositions map[common.Hash]TxPosition) {
	for i, tx := range evmBlock.Transactions {
		txEventPos := txPositions[receipts[i].TxHash]
		txFee := new(big.Int).Mul(new(big.Int).SetUint64(receipts[i].GasUsed), tx.GasPrice())
		a.store.AddDirtyOriginationScore(txEventPos.Creator, txFee)
	}

	if a.blockContext.sealEpoch {
		a.store.DelAllActiveOriginationScores()
		a.store.MoveDirtyOriginationScoresToActive()
	}
}
