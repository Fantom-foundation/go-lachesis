package app

// EthAPIBackend provides methods for ethapi.Backend
type EthAPIBackend struct {
	app *App
}

// EthAPIBackend getter
func (a *App) EthAPIBackend() *EthAPIBackend {
	return &EthAPIBackend{a}
}
