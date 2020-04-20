package gossip

import (
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/tendermint/tendermint/abci/types"
)

func (s *Service) initApp() {
	req := s.config.Net.ChainInfo()
	_ = s.abciApp.InitChain(req)
	// TODO: check the resp
}

func (s *Service) deliverTxRequest(tx *eth.Transaction) types.RequestDeliverTx {
	buf, err := rlp.EncodeToBytes(tx)
	if err != nil {
		s.Log.Crit("Failed to encode rlp", "err", err)
	}

	return types.RequestDeliverTx{
		Tx: buf,
	}
}
