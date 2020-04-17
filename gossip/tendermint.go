package gossip

import (
	"github.com/tendermint/tendermint/abci/types"
)

func (s *Service) chainInfo() types.RequestInitChain {
	return types.RequestInitChain{}
}
