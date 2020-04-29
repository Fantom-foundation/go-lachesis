package app

import (
	"sync"

	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
)

// Checkpoint is for persistent storing.
type Checkpoint struct {
	// fields can change only after a block commit
	EpochN idx.Epoch
	BlockN idx.Block

	EpochBlock idx.Block
	EpochStart inter.Timestamp

	*sync.RWMutex
}

// Bootstrap restores poset's state from store.
func (a *App) Bootstrap() {
	cp := a.store.GetCheckpoint()
	a.checkpoint = *cp
}

// State saves Checkpoint.
func (a *App) saveCheckpoint() {
	a.store.SetCheckpoint(a.checkpoint)
}

// GetLastVoting says when last voting was (sealed epoch)
func (a *App) GetLastVoting() (block idx.Block, start inter.Timestamp) {
	a.checkpoint.RLock()
	defer a.checkpoint.RUnlock()

	return a.checkpoint.EpochBlock, a.checkpoint.EpochStart
}

// SetLastVoting saves when last voting was (sealed epoch)
func (a *App) SetLastVoting(block idx.Block, start inter.Timestamp) {
	a.checkpoint.Lock()
	defer a.checkpoint.Unlock()

	a.checkpoint.EpochBlock = block
	a.checkpoint.EpochStart = start
	a.saveCheckpoint()
}

// GetEpoch returns current epoch num to 3rd party.
func (a *App) GetEpoch() idx.Epoch {
	a.checkpoint.RLock()
	defer a.checkpoint.RUnlock()

	return a.checkpoint.EpochN
}

func (a *App) incEpoch() idx.Epoch {
	a.checkpoint.Lock()
	defer a.checkpoint.Unlock()

	a.checkpoint.EpochN += 1
	a.saveCheckpoint()
	return a.checkpoint.EpochN
}

func (a *App) lastBlock() idx.Block {
	a.checkpoint.RLock()
	defer a.checkpoint.RUnlock()

	return a.checkpoint.BlockN
}

func (a *App) incLastBlock() idx.Block {
	a.checkpoint.Lock()
	defer a.checkpoint.Unlock()

	a.checkpoint.BlockN += 1
	a.saveCheckpoint()
	return a.checkpoint.BlockN
}
