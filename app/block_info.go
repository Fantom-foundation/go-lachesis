package app

import (
	"sync/atomic"

	"github.com/ethereum/go-ethereum/common"

	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
)

// BlockInfo includes only necessary information about inter.Block
type BlockInfo struct {
	Index idx.Block
	Root  common.Hash
	Time  inter.Timestamp
}

func blockInfo(b *inter.Block) *BlockInfo {
	return &BlockInfo{
		Index: b.Index,
		Root:  b.Root,
		Time:  b.Time,
	}
}

// LastBlock returns last block info.
func (a *App) LastBlock() *BlockInfo {
	n := atomic.LoadUint64((*uint64)(&a.block))
	return a.store.GetBlock(idx.Block(n))
}

func (a *App) setLastBlock(n idx.Block) {
	atomic.StoreUint64((*uint64)(&a.block), uint64(n))
}

func (a *App) incLastBlock() idx.Block {
	n := atomic.AddUint64((*uint64)(&a.block), 1)
	return idx.Block(n)
}
