package app

import (
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/Fantom-foundation/go-lachesis/inter/idx"
)

// DagTx is the DAG specific transaction.
// TODO: move it to the more common package.
type DagTx struct {
	Originator  idx.StakerID
	Transaction *eth.Transaction
}

// Bytes gets the byte representation of the tx.
func (tx *DagTx) Bytes() []byte {
	buf, err := rlp.EncodeToBytes(tx)
	if err != nil {
		log.Crit("Failed to encode DagTx", "err", err)
	}
	return buf
}

// BytesToDagTx converts bytes to tx.
func BytesToDagTx(buf []byte) *DagTx {
	tx := new(DagTx)
	err := rlp.DecodeBytes(buf, tx)
	if err != nil {
		log.Crit("Failed to decode DagTx", "err", err, "size", len(buf))
	}
	return tx
}
