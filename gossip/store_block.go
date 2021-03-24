package gossip

import (
	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/ethereum/go-ethereum/rlp"
)

// SetBlock stores chain block.
func (s *Store) SetBlock(b *inter.Block) {
	s.set(s.table.Blocks, b.Index.Bytes(), b)

	// Add to LRU cache.
	if b != nil && s.cache.Blocks != nil {
		s.cache.Blocks.Add(b.Index, b)
	}
}

// GetBlock returns stored block.
func (s *Store) GetBlock(n idx.Block) *inter.Block {
	// Get block from LRU cache first.
	if s.cache.Blocks != nil {
		if c, ok := s.cache.Blocks.Get(n); ok {
			if b, ok := c.(*inter.Block); ok {
				return b
			}
		}
	}

	block, _ := s.get(s.table.Blocks, n.Bytes(), &inter.Block{}).(*inter.Block)

	// Add to LRU cache.
	if block != nil && s.cache.Blocks != nil {
		s.cache.Blocks.Add(n, block)
	}

	return block
}

func (s *Store) ForEachBlock(fn func(index idx.Block, block *inter.Block)) {
	it := s.table.Blocks.NewIterator(nil, nil)
	for it.Next() {
		var block inter.Block
		err := rlp.DecodeBytes(it.Value(), &block)
		if err != nil {
			s.Log.Crit("Failed to decode block", "err", err)
		}
		fn(idx.BytesToBlock(it.Key()), &block)
	}
}

// SetBlockIndex stores chain block index.
func (s *Store) SetBlockIndex(id hash.Event, n idx.Block) {
	if err := s.table.BlockHashes.Put(id.Bytes(), n.Bytes()); err != nil {
		s.Log.Crit("Failed to put key-value", "err", err)
	}

	s.cache.BlockHashes.Add(id, n)
}

// GetBlockIndex returns stored block index.
func (s *Store) GetBlockIndex(id hash.Event) *idx.Block {
	nVal, ok := s.cache.BlockHashes.Get(id)
	if ok {
		n, ok := nVal.(idx.Block)
		if ok {
			return &n
		}
	}

	buf, err := s.table.BlockHashes.Get(id.Bytes())
	if err != nil {
		s.Log.Crit("Failed to get key-value", "err", err)
	}
	if buf == nil {
		return nil
	}
	n := idx.BytesToBlock(buf)

	s.cache.BlockHashes.Add(id, n)

	return &n
}

// GetBlockByHash get block by block hash
func (s *Store) GetBlockByHash(id hash.Event) *inter.Block {
	return s.GetBlock(*s.GetBlockIndex(id))
}
