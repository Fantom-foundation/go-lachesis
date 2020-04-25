package gossip

import (
	"math"

	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/tendermint/tendermint/abci/types"

	"github.com/Fantom-foundation/go-lachesis/app"
	"github.com/Fantom-foundation/go-lachesis/evmcore"
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

func beginBlockRequest(
	cheaters inter.Cheaters,
	stateHash common.Hash,
	block *evmcore.EvmHeader,
) types.RequestBeginBlock {
	byzantines := make([]types.Evidence, len(cheaters))
	for i, stakerID := range cheaters {
		byzantines[i] = types.Evidence{
			Height: int64(stakerID),
		}
	}

	req := types.RequestBeginBlock{
		Hash: block.Hash.Bytes(),
		Header: types.Header{
			Height:        block.Number.Int64(),
			Time:          block.Time.Time(),
			ConsensusHash: block.Hash.Bytes(),
			LastBlockId: types.BlockID{
				Hash: block.ParentHash.Bytes(),
			},
			LastCommitHash: stateHash.Bytes(),
		},
		ByzantineValidators: byzantines,
	}

	if block.GasLimit != math.MaxUint64 {
		panic("we need to pass GasLimit throuth request")
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
