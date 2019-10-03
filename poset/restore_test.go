package poset

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/log"
	"github.com/stretchr/testify/assert"

	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/kvdb/fallible"
	"github.com/Fantom-foundation/go-lachesis/logger"
)

func TestRestore(t *testing.T) {
	logger.SetTestMode(t)
	assertar := assert.New(t)

	const (
		COUNT     = 3 // two poset instances
		GENERATOR = 0 // event generator
		EXPECTED  = 1 // first as etalon
		RESTORED  = 2 // second with restoring
	)

	nodes := inter.GenNodes(5)
	posets := make([]*ExtendedPoset, 0, COUNT)
	inputs := make([]*EventStore, 0, COUNT)
	namespaces := make([]string, 0, COUNT)
	for i := 0; i < COUNT; i++ {
		namespace := uniqNamespace()
		poset, _, input := FakePoset(namespace, nodes)
		posets = append(posets, poset)
		inputs = append(inputs, input)
		namespaces = append(namespaces, namespace)
	}

	posets[GENERATOR].
		SetName("generator")
	posets[GENERATOR].store.
		SetName("generator")

	const epochs = 2
	var epochLen = int(posets[GENERATOR].dag.EpochLen)

	// create events
	var ordered []*inter.Event
	for epoch := idx.Epoch(1); epoch <= idx.Epoch(epochs); epoch++ {
		r := rand.New(rand.NewSource(int64((epoch))))
		_ = inter.ForEachRandEvent(nodes, epochLen*2, COUNT, r, inter.ForEachEvent{
			Process: func(e *inter.Event, name string) {
				inputs[GENERATOR].SetEvent(e)
				assertar.NoError(posets[GENERATOR].ProcessEvent(e))

				ordered = append(ordered, e)
			},
			Build: func(e *inter.Event, name string) *inter.Event {
				e.Epoch = epoch
				return posets[GENERATOR].Prepare(e)
			},
		})
	}

	posets[EXPECTED].
		SetName("expected")
	posets[EXPECTED].store.
		SetName("expected")

	posets[RESTORED].
		SetName("restored-0")
	posets[RESTORED].store.
		SetName("restored-0")

	// use pre-ordered events, call consensus(e) directly, to avoid issues with restoring state of EventBuffer
	x := 0
	for n, e := range ordered {
		if n%20 == 0 {
			log.Info("Restart poset")
			prev := posets[RESTORED]
			x++
			fs := newFakeFS(namespaces[RESTORED])
			store := NewStore(fs.OpenFakeDB(""), fs.OpenFakeDB)
			store.SetName(fmt.Sprintf("restored-%d", x))

			restored := New(prev.dag, store, prev.input)
			restored.SetName(fmt.Sprintf("restored-%d", x))
			restored.Bootstrap(prev.applyBlock)

			posets[RESTORED].Poset = restored
		}

		inputs[EXPECTED].SetEvent(e)
		assertar.NoError(posets[EXPECTED].ProcessEvent(e))

		inputs[RESTORED].SetEvent(e)
		assertar.NoError(posets[RESTORED].ProcessEvent(e))

		compareStates(assertar, posets[EXPECTED], posets[RESTORED])
		if t.Failed() {
			return
		}
	}

	if !assertar.Equal(epochLen*epochs, len(posets[EXPECTED].blocks)) {
		return
	}
	compareBlocks(assertar, posets[EXPECTED], posets[RESTORED])
}

func TestDbFailure(t *testing.T) {
	logger.SetTestMode(t)
	assertar := assert.New(t)

	const (
		COUNT     = 3 // two poset instances
		GENERATOR = 0 // event generator
		EXPECTED  = 1 // first as etalon
		RESTORED  = 2 // second with db failures
	)
	nodes := inter.GenNodes(5)

	posets := make([]*ExtendedPoset, 0, COUNT)
	inputs := make([]*EventStore, 0, COUNT)
	namespaces := make([]string, 0, COUNT)
	for i := 0; i < COUNT; i++ {
		namespace := uniqNamespace()
		poset, _, input := FakePoset(namespace, nodes)
		posets = append(posets, poset)
		inputs = append(inputs, input)
		namespaces = append(namespaces, namespace)
	}

	posets[GENERATOR].
		SetName("generator")
	posets[GENERATOR].store.
		SetName("generator")

	epochLen := int(posets[GENERATOR].dag.EpochLen)

	// create events on etalon poset
	var ordered inter.Events
	inter.ForEachRandEvent(nodes, epochLen-1, COUNT, nil, inter.ForEachEvent{
		Process: func(e *inter.Event, name string) {
			ordered = append(ordered, e)

			inputs[GENERATOR].SetEvent(e)
			assertar.NoError(
				posets[GENERATOR].ProcessEvent(e))
		},
		Build: func(e *inter.Event, name string) *inter.Event {
			e.Epoch = 1
			return posets[GENERATOR].Prepare(e)
		},
	})

	posets[EXPECTED].
		SetName("expected")
	posets[EXPECTED].store.
		SetName("expected")

	posets[RESTORED].
		SetName("restored-0")
	posets[RESTORED].store.
		SetName("restored-0")

	// db writes limit
	db := posets[RESTORED].store.UnderlyingDB().(*fallible.Fallible)
	db.SetWriteCount(100)

	x := 0
	process := func(e *inter.Event) (ok bool) {
		ok = true
		defer func() {
			// catch a panic
			if r := recover(); r == nil {
				return
			}
			ok = false

			db.SetWriteCount(100)

			log.Info("Restart poset after db failure")
			prev := posets[RESTORED]
			x++
			fs := newFakeFS(namespaces[RESTORED])
			store := NewStore(fs.OpenFakeDB(""), fs.OpenFakeDB)
			store.SetName(fmt.Sprintf("restored-%d", x))

			restored := New(prev.dag, store, prev.input)
			restored.SetName(fmt.Sprintf("restored-%d", x))
			restored.Bootstrap(prev.applyBlock)

			posets[RESTORED].Poset = restored
		}()

		inputs[RESTORED].SetEvent(e)
		assertar.NoError(
			posets[RESTORED].ProcessEvent(e))

		inputs[EXPECTED].SetEvent(e)
		assertar.NoError(
			posets[EXPECTED].ProcessEvent(e))

		return
	}

	for len(ordered) > 0 {
		e := ordered[0]
		if e.Epoch != 1 {
			panic("sanity check")
		}

		if !process(e) {
			continue
		}

		ordered = ordered[1:]

		compareStates(assertar, posets[EXPECTED], posets[RESTORED])
		if t.Failed() {
			return
		}
	}

	compareBlocks(assertar, posets[EXPECTED], posets[RESTORED])
}

func compareStates(assertar *assert.Assertions, expected, restored *ExtendedPoset) {
	assertar.Equal(
		*expected.checkpoint, *restored.checkpoint)
	assertar.Equal(
		expected.epochState.PrevEpoch.Hash(), restored.epochState.PrevEpoch.Hash())
	assertar.Equal(
		expected.epochState.Validators, restored.epochState.Validators)
	assertar.Equal(
		expected.epochState.EpochN, restored.epochState.EpochN)
	// check LastAtropos and Head() method
	if expected.checkpoint.LastBlockN != 0 {
		assertar.Equal(
			expected.blocks[idx.Block(len(expected.blocks))].Hash(),
			restored.checkpoint.LastAtropos,
			"atropos must be last event in block")
	}
}

func compareBlocks(assertar *assert.Assertions, expected, restored *ExtendedPoset) {
	assertar.Equal(len(expected.blocks), len(restored.blocks))
	assertar.Equal(len(expected.blocks), int(restored.LastBlockN))
	for i := idx.Block(1); i <= idx.Block(len(restored.blocks)); i++ {
		if !assertar.NotNil(restored.blocks[i]) ||
			!assertar.Equal(expected.blocks[i], restored.blocks[i]) {
			return
		}

	}
}