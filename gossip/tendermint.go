package gossip

import (
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/tendermint/tendermint/abci/types"

	"github.com/Fantom-foundation/go-lachesis/app"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
)

func (s *Service) initApp() {
	req := s.config.Net.ChainInfo()
	_ = s.abciApp.InitChain(req)
	// TODO: check the resp
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
