package app

import (
	"sync/atomic"

	"github.com/Fantom-foundation/go-lachesis/inter/idx"
)

// GetEpoch returns current epoch num to 3rd party.
func (a *App) GetEpoch() idx.Epoch {
	e := atomic.LoadUint32((*uint32)(&a.epoch))
	return idx.Epoch(e)
}

func (a *App) setEpoch(e idx.Epoch) {
	atomic.StoreUint32((*uint32)(&a.epoch), uint32(e))
}

func (a *App) incEpoch() idx.Epoch {
	e := atomic.AddUint32((*uint32)(&a.epoch), 1)
	return idx.Epoch(e)
}
