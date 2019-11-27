// Package types contains data types related to Ethereum consensus.
package types

import (
	"github.com/Fantom-foundation/go-ethereum/core/types"
)

var (
	EmptyRootHash  = types.EmptyRootHash
	EmptyUncleHash = types.EmptyUncleHash
)

// A BlockNonce is a 64-bit hash which proves (combined with the
// mix-hash) that a sufficient amount of computation has been carried
// out on a block.
type BlockNonce = types.BlockNonce

// Header represents a block header in the Ethereum blockchain.
type Header = types.Header

// Body is a simple (mutable, non-safe) data container for storing and moving
// a block's data contents (transactions and uncles) together.
type Body = types.Body

// Block represents an entire block in the Ethereum blockchain.
type Block = types.Block

func NewBlock(header *Header, txs []*Transaction, uncles []*Header, receipts []*Receipt) *Block {
	return types.NewBlock(header, txs, uncles, receipts)
}


// [deprecated by eth/63]
// StorageBlock defines the RLP encoding of a Block stored in the
// state database. The StorageBlock encoding contains fields that
// would otherwise need to be recomputed.
type StorageBlock = types.StorageBlock

type Blocks = types.Blocks

type BlockBy = types.BlockBy
