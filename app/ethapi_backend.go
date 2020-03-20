package app

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/ethdb"

	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
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

// GetValidationScore returns staker's ValidationScore.
func (b *EthAPIBackend) GetValidationScore(ctx context.Context, stakerID idx.StakerID) (*big.Int, error) {
	if !b.app.store.HasSfcStaker(stakerID) {
		return nil, nil
	}
	return b.app.store.GetActiveValidationScore(stakerID), nil
}

// GetOriginationScore returns staker's OriginationScore.
func (b *EthAPIBackend) GetOriginationScore(ctx context.Context, stakerID idx.StakerID) (*big.Int, error) {
	if !b.app.store.HasSfcStaker(stakerID) {
		return nil, nil
	}
	return b.app.store.GetActiveOriginationScore(stakerID), nil
}

// GetStakerPoI returns staker's PoI.
func (b *EthAPIBackend) GetStakerPoI(ctx context.Context, stakerID idx.StakerID) (*big.Int, error) {
	if !b.app.store.HasSfcStaker(stakerID) {
		return nil, nil
	}
	return b.app.store.GetStakerPOI(stakerID), nil
}

// GetDowntime returns staker's Downtime.
func (b *EthAPIBackend) GetDowntime(ctx context.Context, stakerID idx.StakerID) (idx.Block, inter.Timestamp, error) {
	missed := b.app.store.GetBlocksMissed(stakerID)
	return missed.Num, missed.Period, nil
}
