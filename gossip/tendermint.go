package gossip

import (
	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/tendermint/tendermint/abci/types"

	"github.com/Fantom-foundation/go-lachesis/app"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
)

const (
	txIsFullyValid = 0
)

func (s *Service) initApp() {
	req := s.config.Net.ChainInfo()
	_ = s.abciApp.InitChain(req)
	// TODO: check the resp
}

func beginBlockRequest(cheaters inter.Cheaters, stateHash common.Hash) types.RequestBeginBlock {
	req := types.RequestBeginBlock{
		Header: types.Header{
			LastCommitHash: stateHash.Bytes(),
		},
		ByzantineValidators: make([]types.Evidence, len(cheaters)),
	}

	for i, stakerID := range cheaters {
		req.ByzantineValidators[i] = types.Evidence{
			Height: int64(stakerID),
		}
	}

	return req
}

func deliverTxRequest(tx *eth.Transaction, originator idx.StakerID) types.RequestDeliverTx {
	t := app.DagTx{
		Transaction: tx,
		Originator:  originator,
	}

	return types.RequestDeliverTx{
		Tx: t.Bytes(),
	}
}

func endBlockRequest(block idx.Block) types.RequestEndBlock {
	return types.RequestEndBlock{
		Height: int64(block),
	}
}
