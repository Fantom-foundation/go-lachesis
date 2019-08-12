package posposet

import (
	"github.com/Fantom-foundation/go-lachesis/src/inter"
	"github.com/Fantom-foundation/go-lachesis/src/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/src/logger"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestRestore(t *testing.T) {
	logger.SetTestMode(t)

	const posetCount = 3 // 2 last will be restored
	const epochs = idx.SuperFrame(2)

	nodes := inter.GenNodes(5)

	posets := make([]*Poset, 0, posetCount)
	inputs := make([]*EventStore, 0, posetCount)

	makePoset := func(i int) *Store {
		poset, store, input := FakePoset(nodes)
		n := i % len(nodes)
		poset.SetName(nodes[n].String())
		store.SetName(nodes[n].String())
		poset.Start()
		posets = append(posets, poset)
		inputs = append(inputs, input)
		return store
	}

	for i := 0; i < posetCount-1; i++ {
		_ = makePoset(i)
	}

	// create events on poset0
	var ordered []*inter.Event
	for epoch := idx.SuperFrame(1); epoch <= epochs; epoch++ {
		buildEvent := func(e *inter.Event) *inter.Event {
			e.Epoch = epoch
			return posets[0].Prepare(e)
		}
		onNewEvent := func(e *inter.Event) {
			inputs[0].SetEvent(e)
			posets[0].PushEventSync(e.Hash())

			ordered = append(ordered, e)
		}
		r := rand.New(rand.NewSource(int64((epoch))))
		_ = inter.GenEventsByNode(nodes, int(SuperFrameLen)*3, 3, buildEvent, onNewEvent, r)
	}

	t.Run("Restore", func(t *testing.T) {
		assertar := assert.New(t)

		i := posetCount - 1
		j := posetCount - 2
		store := makePoset(i)

		// use pre-ordered events, call consensus(e) directly, to avoid issues with restoring state of EventBuffer
		for x, e := range ordered {
			if (x < len(ordered)/4) || x%20 == 0 {
				// restore
				posets[i].Stop()
				restored := New(store, inputs[i])
				n := i % len(nodes)
				restored.SetName("restored_" + nodes[n].String())
				store.SetName("restored_" + nodes[n].String())
				restored.Bootstrap()
				posets[i] = restored
				posets[i].Start()
			}
			// push on restore i, and non-restored j
			inputs[i].SetEvent(e)
			posets[i].consensus(e)

			inputs[j].SetEvent(e)
			posets[j].consensus(e)
			// compare state on i/j
			assertar.Equal(*posets[j].checkpoint, *posets[i].checkpoint)
			assertar.Equal(posets[j].superFrame, posets[i].superFrame)
		}
	})

	for i := 0; i < len(posets); i++ {
		posets[i].Stop()
	}
}