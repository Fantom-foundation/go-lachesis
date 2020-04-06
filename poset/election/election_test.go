package election

import (
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/pos"
	"github.com/Fantom-foundation/go-lachesis/utils"
)

type fakeEdge struct {
	from hash.Event
	to   hash.Event
}

type (
	stakes map[string]pos.Stake
)

type testExpected struct {
	DecidedFrame   idx.Frame
	DecidedAtropos string
	DecisiveRoots  map[string]bool
}

func TestProcessRoot(t *testing.T) {
	t.Run("4 equalStakes notDecided", func(t *testing.T) {
		testProcessRoot(t,
			1,
			nil,
			stakes{
				"nodeA": 1,
				"nodeB": 1,
				"nodeC": 1,
				"nodeD": 1,
			}, `
a0    b0    c0    d0
║     ║     ║     ║
a1════╬═════╣     ║
║     ║     ║     ║
║╚════b1════╣     ║
║     ║     ║     ║
║     ║╚════c1════╣
║     ║     ║     ║
║     ║╚═══─╫╩════d1
║     ║     ║     ║
a2════╬═════╬═════╣
║     ║     ║     ║
`)
	})

	t.Run("4 equalStakes", func(t *testing.T) {
		testProcessRoot(t,
			1,
			&testExpected{
				DecidedFrame:   1,
				DecidedAtropos: "b0",
				DecisiveRoots:  map[string]bool{"a2": true},
			},
			stakes{
				"nodeA": 1,
				"nodeB": 1,
				"nodeC": 1,
				"nodeD": 1,
			}, `
a0    b0    c0    d0
║     ║     ║     ║
a1════╬═════╣     ║
║     ║     ║     ║
║     b1════╬═════╣
║     ║     ║     ║
║     ║╚════c1════╣
║     ║     ║     ║
║     ║╚═══─╫╩════d1
║     ║     ║     ║
a2════╬═════╬═════╣
║     ║     ║     ║
`)
	})

	t.Run("4 equalStakes missingRoot", func(t *testing.T) {
		testProcessRoot(t,
			1,
			&testExpected{
				DecidedFrame:   1,
				DecidedAtropos: "a0",
				DecisiveRoots:  map[string]bool{"a2": true},
			},
			stakes{
				"nodeA": 1,
				"nodeB": 1,
				"nodeC": 1,
				"nodeD": 1,
			}, `
a0    b0    c0    d0
║     ║     ║     ║
a1════╬═════╣     ║
║     ║     ║     ║
║╚════b1════╣     ║
║     ║     ║     ║
║╚═══─╫╩════c1    ║
║     ║     ║     ║
a2════╬═════╣     ║
║     ║     ║     ║
`)
	})

	t.Run("4 differentStakes", func(t *testing.T) {
		testProcessRoot(t,
			1,
			&testExpected{
				DecidedFrame:   1,
				DecidedAtropos: "a0",
				DecisiveRoots:  map[string]bool{"b2": true},
			},
			stakes{
				"nodeA": 1000000000000000000,
				"nodeB": 1,
				"nodeC": 1,
				"nodeD": 1,
			}, `
a0    b0    c0    d0  
║     ║     ║     ║
a1════╬═════╣     ║
║     ║     ║     ║
║╚════b1    ║     ║
║     ║     ║     ║
║╚═══─╫─════c1    ║
║     ║     ║     ║
║╚═══─╫╩═══─╫╩════d1
║     ║     ║     ║
╠═════b2════╬═════╣
║     ║     ║     ║
`)
	})

	t.Run("5 differentStakes 3 rounds", func(t *testing.T) {
		testProcessRoot(t,
			1,
			&testExpected{
				DecidedFrame:   1,
				DecidedAtropos: "nodeB001",
				DecisiveRoots:  map[string]bool{"nodeA004": true},
			},
			stakes{
				"nodeA": 5,
				"nodeB": 4,
				"nodeC": 3,
				"nodeD": 3,
				"nodeE": 3,
			}, `
 ║         
 nodeA1    
 ║          ║         
 ║          nodeB1    
 ║          ║          ║         
 ║          ║          nodeC1    
 ║          ║          ║          ║         
 ║          ║          ║          nodeD1    
 ║          ║          ║          ║          ║         
 ║          ║          ║          ║          nodeE1    
 ║          ║          ║          ║          ║         
 nodeA2═════╬══════════╬══════════╣          ║         
 ║          ║          ║          ║          ║         
 ║          nodeB2═════╬══════════╬══════════╣         
 ║║         ║║         ║          ║          ║         
 ║╚════════─╫╩════════ nodeC2═════╣          ║         
 ║          ║║         ║║         ║          ║         
 ║          ║╚════════─╫╩════════ nodeD2═════╣         
 ║          ║║         ║║         ║║         ║         
 ║          ║╚════════─╫╩════════─╫╩════════ nodeE2    
 ║║        ║║         ║║         ║║         ║║         
 ║╚════════╚ nodeB_2══╩╫─════════╩╫─════════╝║        
 ║          ║          ║          ║          ║         
 nodeA3════─╫─═════════╬══════════╬══════════╣         
 ║          ║║         ║          ║          ║         
 ║          ║╚════════ nodeC3═════╬══════════╣         
 ║          ║║         ║║         ║          ║         
 ║          ║╚════════─╫╩════════ nodeD3═════╣         
 ║          ║║         ║║         ║║         ║         
 ║          ║╚════════─╫╩════════─╫╩════════ nodeE3    
 ║          ║         ║║         ║║         ║║         
 ║          nodeB3════╩╫─════════╩╫─════════╝║         
 ║          ║          ║          ║          ║         
 nodeA4═════╬══════════╬══════════╬══════════╣         
`)
	})

}

func testProcessRoot(
	t *testing.T,
	frameToDecide idx.Frame,
	expected *testExpected,
	stakes stakes,
	dag string,
) {
	assertar := assert.New(t)
	require := require.New(t)

	// events:
	ordered := make(inter.Events, 0)
	events := make(map[hash.Event]*inter.Event)
	frameRoots := make(map[idx.Frame][]RootAndSlot)
	vertices := make(map[hash.Event]Slot)
	edges := make(map[fakeEdge]bool)

	nodes, _, _ := inter.ASCIIschemeForEach(dag, inter.ForEachEvent{
		Process: func(root *inter.Event, name string) {
			// store all the events
			ordered = append(ordered, root)

			events[root.Hash()] = root

			slot := Slot{
				Frame:     root.Frame,
				Validator: root.Creator,
			}
			vertices[root.Hash()] = slot

			frameRoots[root.Frame] = append(frameRoots[root.Frame], RootAndSlot{
				ID:   root.Hash(),
				Slot: slot,
			})

			// build edges to be able to fake forkless cause fn
			from := root.Hash()
			for _, observed := range root.Parents {
				to := observed
				edge := fakeEdge{
					from: from,
					to:   to,
				}
				edges[edge] = true
			}
		},
		Build: func(e *inter.Event, name string) *inter.Event {
			e.Frame = idx.Frame(e.Seq)
			return e
		},
	})

	validatorsBuilder := pos.NewBuilder()
	for _, node := range nodes {
		validatorsBuilder.Set(node, stakes[utils.NameOf(node)])
	}
	validators := validatorsBuilder.Build()

	forklessCauseFn := func(a hash.Event, b hash.Event) bool {
		edge := fakeEdge{
			from: a,
			to:   b,
		}
		return edges[edge]
	}
	getFrameRootsFn := func(f idx.Frame) []RootAndSlot {
		return frameRoots[f]
	}

	// re-order events randomly, preserving parents order
	unordered := make(inter.Events, len(ordered))
	for i, j := range rand.Perm(len(ordered)) {
		unordered[i] = ordered[j]
	}
	ordered = unordered.ByParents()

	election := New(validators, frameToDecide, forklessCauseFn, getFrameRootsFn)

	// processing:
	var alreadyDecided bool
	for _, root := range ordered {
		rootHash := root.Hash()
		rootSlot, ok := vertices[rootHash]
		require.True(ok, "inconsistent vertices")

		got, err := election.ProcessRoot(RootAndSlot{
			ID:   rootHash,
			Slot: rootSlot,
		})
		require.NoError(err)

		// checking:
		decisive := expected != nil && expected.DecisiveRoots[root.Hash().String()]
		if decisive || alreadyDecided {
			require.NotNil(got)

			assertar.Equal(expected.DecidedFrame, got.Frame)
			assertar.Equal(expected.DecidedAtropos, got.Atropos.String())
			alreadyDecided = true
		} else {
			assertar.Nil(got)
		}
	}
}
