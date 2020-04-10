package inter

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Fantom-foundation/go-lachesis/hash"
)

func TestEventsByParents(t *testing.T) {
	require := require.New(t)
	nodes := GenNodes(5)
	events := GenRandEvents(nodes, 10, 3, nil)
	var unordered Events
	for _, ee := range events {
		unordered = append(unordered, ee...)
	}

	ordered := unordered.ByParents()
	position := make(map[hash.Event]int)
	for i, e := range ordered {
		position[e.Hash()] = i
	}

	for i, e := range ordered {
		for _, p := range e.Parents {
			pos, ok := position[p]
			if !ok {
				continue
			}
			require.LessOrEqualf(pos, i, "parent %s is not before %s", p.String(), e.Hash().String())
		}
	}
}
