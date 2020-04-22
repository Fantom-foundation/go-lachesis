package app

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
)

// BlockInfo includes only necessary information about inter.Block
type BlockInfo struct {
	Index      idx.Block
	Hash       hash.Event
	ParentHash hash.Event
	Root       common.Hash
	Time       inter.Timestamp
}

func blockInfo(b *evmcore.EvmHeader) *BlockInfo {
	return &BlockInfo{
		Index:      idx.Block(b.Number.Uint64()),
		Hash:       hash.Event(b.Hash),
		ParentHash: hash.Event(b.ParentHash),
		Root:       b.Root,
		Time:       b.Time,
	}
}

// BlockChain dummy reader.
func (a *App) BlockChain() evmcore.DummyChain {
	return a.store
}

// LastBlock returns last block info.
func (a *App) LastBlock() *BlockInfo {
	n := a.lastBlock()
	return a.store.GetBlock(idx.Block(n))
}
