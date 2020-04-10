package app

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/Fantom-foundation/go-lachesis/evmcore"
	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
)

// SetBlock stores chain block.
func (s *Store) SetBlock(b *BlockInfo) {
	s.set(s.table.Blocks, b.Index.Bytes(), b)

	// Add to LRU cache.
	if b != nil && s.cache.Blocks != nil {
		s.cache.Blocks.Add(b.Index, b)
	}
}

// GetBlock returns stored block.
func (s *Store) GetBlock(n idx.Block) *BlockInfo {
	// Get block from LRU cache first.
	if s.cache.Blocks != nil {
		if c, ok := s.cache.Blocks.Get(n); ok {
			if b, ok := c.(*BlockInfo); ok {
				return b
			}
		}
	}

	block, _ := s.get(s.table.Blocks, n.Bytes(), &BlockInfo{}).(*BlockInfo)

	// Add to LRU cache.
	if block != nil && s.cache.Blocks != nil {
		s.cache.Blocks.Add(n, block)
	}

	return block
}

// GetHeader implements evmcore.DummyChain interface to be used during transaction processing.
// Only Number and ParentHash of header are filled.
func (s *Store) GetHeader(h common.Hash, n uint64) *evmcore.EvmHeader {
	info := s.GetBlock(idx.Block(n))
	if info == nil {
		return nil
	}
	if (h != common.Hash{}) && (info.Hash != hash.Event(h)) {
		return nil
	}

	return &evmcore.EvmHeader{
		Number:     big.NewInt(int64(info.Index)),
		ParentHash: common.Hash(info.ParentHash),
	}
}
