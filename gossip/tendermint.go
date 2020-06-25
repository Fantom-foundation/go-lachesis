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

func beginBlockRequest(
	cheaters inter.Cheaters,
	stateHash common.Hash,
	block *inter.Block,
	participated map[idx.StakerID]bool,
) types.RequestBeginBlock {

	byzantines := make([]types.Evidence, len(cheaters))
	for i, stakerID := range cheaters {
		byzantines[i] = types.Evidence{
			Height: int64(stakerID),
		}
	}

	votes := make([]types.VoteInfo, 0, len(participated))
	for staker, ok := range participated {
		if !ok {
			continue
		}
		votes = append(votes, types.VoteInfo{
			Validator: types.Validator{
				Address: staker.Bytes(),
			},
		})
	}

	req := types.RequestBeginBlock{
		Hash: block.Atropos.Bytes(),
		Header: types.Header{
			Height:        int64(block.Index),
			Time:          block.Time.Time(),
			ConsensusHash: block.Atropos.Bytes(),
			LastBlockId: types.BlockID{
				Hash: block.PrevHash.Bytes(),
			},
			LastCommitHash: stateHash.Bytes(),
		},
		LastCommitInfo: types.LastCommitInfo{
			Votes: votes,
		},
		ByzantineValidators: byzantines,
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
