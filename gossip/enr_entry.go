package gossip

import (
	"github.com/Fantom-foundation/go-ethereum/core/forkid"
	"github.com/Fantom-foundation/go-ethereum/rlp"
)

// Enr is ENR entry which advertises eth protocol
// on the discovery network.
type Enr struct {
	ForkID forkid.ID
	// Ignore additional fields (for forward compatibility).
	Rest []rlp.RawValue `rlp:"tail"`
}

// ENRKey implements enr.Entry.
func (e Enr) ENRKey() string {
	return "lachesis"
}

func (s *Service) currentEnr() *Enr {
	return &Enr{}
}
