package node

import (
	"fmt"

	"github.com/Fantom-foundation/go-lachesis/src/poset"
)

// Infos struct for graph data (visualizer)
type Infos struct {
	ParticipantEvents map[string]map[string]*poset.Event
	Rounds            []*poset.RoundCreated
	Blocks            []*poset.Block
}

// Graph stuct to represent the DAG
type Graph struct {
	*Node
}

// GetBlocks returns all blocks in the DAG
func (g *Graph) GetBlocks() []*poset.Block {
	var res []*poset.Block
	store := g.Node.core.poset.Store
	blockIdx := store.LastBlockIndex() - 10

	if blockIdx < 0 {
		blockIdx = 0
	}

	for blockIdx <= store.LastBlockIndex() {
		r, err := store.GetBlock(blockIdx)
		if err != nil {
			break
		}
		res = append(res, r)
		blockIdx++
	}
	return res
}

// GetParticipantEvents returns all known events per participant
func (g *Graph) GetParticipantEvents() map[string]map[string]*poset.Event {
	res := make(map[string]map[string]*poset.Event)

	store := g.Node.core.poset.Store
	repertoire := g.Node.core.poset.Store.RepertoireByPubKey()
	known := store.KnownEvents()
	for _, p := range repertoire {
		root, err := store.GetRoot(p.PubKeyHex)

		if err != nil {
			panic(err)
		}

		skip := known[p.ID] - 30
		if skip < 0 {
			skip = -1
		}

		evs, err := store.ParticipantEvents(p.PubKeyHex, skip)

		if err != nil {
			panic(err)
		}

		res[p.PubKeyHex] = make(map[string]*poset.Event)

		selfParent := fmt.Sprintf("Root%d", p.ID)

		flagTable := make(map[string]int64)
		flagTable[selfParent] = 1

		// Create and save the first Event
		initialEvent := poset.NewEvent([][]byte{},
			[]*poset.InternalTransaction{},
			[]poset.BlockSignature{},
			[]string{}, []byte{}, 0, flagTable)

		res[p.PubKeyHex][root.SelfParent.Hash] = initialEvent

		for _, e := range evs {
			event, err := store.GetEventBlock(e)

			if err != nil {
				panic(err)
			}

			hash := event.Hex()

			res[p.PubKeyHex][hash] = event
		}
	}

	return res
}

// GetRounds returns the created rounds for the DAG
func (g *Graph) GetRounds() []*poset.RoundCreated {
	var res []*poset.RoundCreated

	store := g.Node.core.poset.Store

	round := store.LastRound() - 20

	if round < 0 {
		round = 0
	}

	for round <= store.LastRound() {
		r, err := store.GetRoundCreated(round)

		if err != nil {
			break
		}

		res = append(res, r)

		round++
	}

	return res
}

// GetInfos returns the info subset for the DAG
func (g *Graph) GetInfos() Infos {
	return Infos{
		ParticipantEvents: g.GetParticipantEvents(),
		Rounds:            g.GetRounds(),
		Blocks:            g.GetBlocks(),
	}
}

// NewGraph creates a new DAG
func NewGraph(n *Node) *Graph {
	return &Graph{
		Node: n,
	}
}
