package lachesis

import (
	"strconv"

	"github.com/tendermint/tendermint/abci/types"
)

func (c *Config) ChainInfo() types.RequestInitChain {
	return types.RequestInitChain{
		ChainId: strconv.FormatUint(c.NetworkID, 10),
		Time:    c.Genesis.Time.Time(),
		// TODO: fill validators
	}
}
