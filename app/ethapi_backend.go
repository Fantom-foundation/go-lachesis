package app

import (
	"github.com/ethereum/go-ethereum/ethdb"

	"github.com/Fantom-foundation/go-lachesis/topicsdb"
)

// EthAPIBackend provides methods for ethapi.Backend
type EthAPIBackend struct {
	app *App
}

// EthAPIBackend getter
func (a *App) EthAPIBackend() *EthAPIBackend {
	return &EthAPIBackend{a}
}

func (b *EthAPIBackend) ChainDb() ethdb.Database {
	return b.app.store.EvmTable()
}

func (b *EthAPIBackend) EvmLogIndex() *topicsdb.Index {
	return b.app.store.EvmLogs()
}
