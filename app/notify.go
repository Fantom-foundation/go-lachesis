package app

import (
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	notify "github.com/ethereum/go-ethereum/event"
)

// Feed with notifications.
type Feed struct {
	scope notify.SubscriptionScope

	newTxs  notify.Feed
	newLogs notify.Feed
}

func (f *Feed) Close() {
	f.scope.Close()
}

func (f *Feed) SubscribeNewTxs(ch chan<- core.NewTxsEvent) notify.Subscription {
	return f.scope.Track(f.newTxs.Subscribe(ch))
}

func (f *Feed) SubscribeNewLogs(ch chan<- []*types.Log) notify.Subscription {
	return f.scope.Track(f.newLogs.Subscribe(ch))
}

// Start apps service (non-tendermint).
func (a *App) Start() {

}

// Stop apps service (non-tendermint).
func (a *App) Stop() {
	a.Feed.Close()
}
