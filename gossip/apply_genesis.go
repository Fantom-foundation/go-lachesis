package gossip

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/lachesis"
)

// GenesisMismatchError is raised when trying to overwrite an existing
// genesis block with an incompatible one.
type GenesisMismatchError struct {
	Stored, New hash.Event
}

// Error implements error interface.
func (e *GenesisMismatchError) Error() string {
	return fmt.Sprintf("database contains incompatible gossip genesis (have %s, new %s)", e.Stored.FullID(), e.New.FullID())
}

// ApplyGenesis writes initial state.
func (s *Store) ApplyGenesis(
	net *lachesis.Config, stateRoot common.Hash,
) (
	atropos hash.Event, appState common.Hash, isNew bool, err error,
) {
	s.migrate()

	stored := s.GetBlock(0)
	if stored != nil {
		isNew = false
		atropos = calcGenesisHash(net, stateRoot)
		if stored.Atropos != atropos {
			err = &GenesisMismatchError{stored.Atropos, atropos}
			return
		}
		appState = common.Hash(atropos) // is not state.Root because legacy
		return
	}
	// if we'here, then it's first time genesis is applied
	isNew = true
	atropos = s.applyGenesis(net, stateRoot)
	appState = common.Hash(atropos) // is not state.Root because legacy
	return
}

// calcGenesisHash calcs hash of genesis atropos.
func calcGenesisHash(net *lachesis.Config, stateRoot common.Hash) hash.Event {
	s := NewMemStore()
	defer s.Close()

	s.Log.SetHandler(log.DiscardHandler())

	atropos := s.applyGenesis(net, stateRoot)
	return atropos
}

func (s *Store) applyGenesis(net *lachesis.Config, stateRoot common.Hash) (atropos hash.Event) {
	prettyHash := func(net *lachesis.Config) hash.Event {
		e := inter.NewEvent()
		// for nice-looking ID
		e.Epoch = 0
		e.Lamport = idx.Lamport(net.Dag.MaxEpochBlocks)
		// actual data hashed
		e.Extra = net.Genesis.ExtraData
		e.ClaimedTime = net.Genesis.Time
		e.TxHash = net.Genesis.Alloc.Accounts.Hash()

		return e.CalcHash()
	}
	atropos = prettyHash(net)

	block := inter.NewBlock(0,
		net.Genesis.Time,
		atropos,
		hash.Event{},
		hash.Events{atropos},
	)
	block.Root = stateRoot

	s.SetBlock(block)
	s.SetBlockIndex(atropos, block.Index)

	return atropos
}
