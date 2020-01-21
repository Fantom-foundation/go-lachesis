package gossip

/*
	In LRU cache data stored like pointer
*/

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/Fantom-foundation/go-lachesis/app"
)

// SetTxPosition stores transaction block and position.
func (s *Store) SetTxPosition(txid common.Hash, position *app.TxPosition) {
	s.set(s.table.TxPositions, txid.Bytes(), position)

	// Add to LRU cache.
	if position != nil && s.cache.TxPositions != nil {
		s.cache.TxPositions.Add(txid.String(), position)
	}
}

// GetTxPosition returns stored transaction block and position.
func (s *Store) GetTxPosition(txid common.Hash) *app.TxPosition {
	// Get data from LRU cache first.
	if s.cache.TxPositions != nil {
		if c, ok := s.cache.TxPositions.Get(txid.String()); ok {
			if b, ok := c.(*app.TxPosition); ok {
				return b
			}
		}
	}

	txPosition, _ := s.get(s.table.TxPositions, txid.Bytes(), &app.TxPosition{}).(*app.TxPosition)

	// Add to LRU cache.
	if txPosition != nil && s.cache.TxPositions != nil {
		s.cache.TxPositions.Add(txid.String(), txPosition)
	}

	return txPosition
}
