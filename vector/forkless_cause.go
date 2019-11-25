package vector

import (
	"github.com/Fantom-foundation/go-lachesis/hash"
)

type kv struct {
	a, b hash.Event
}

// ForklessCause calculates "sufficient coherence" between the events.
// The A.HighestBefore array remembers the sequence number of the last
// event by each validator that is an ancestor of A. The array for
// B.LowestAfter remembers the sequence number of the earliest
// event by each validator that is a descendant of B. Compare the two arrays,
// and find how many elements in the A.HighestBefore array are greater
// than or equal to the corresponding element of the B.LowestAfter
// array. If there are more than 2n/3 such matches, then the A and B
// have achieved sufficient coherency.
//
// If B1 and B2 are forks, then they cannot BOTH forkless-cause any specific event A,
// unless more than 1/3W are Byzantine.
// This great property is the reason why this function exists,
// providing the base for the BFT algorithm.
func (vi *Index) ForklessCause(aID, bID hash.Event) bool {
	// TODO	
//	if res, ok := vi.forklessCauseCache.Get(kv{aID, bID}); ok {
//		return res.(bool)
//	}

	res := vi.forklessCause(aID, bID)

//	vi.forklessCauseCache.Add(kv{aID, bID}, res)
	return res
}

func (vi *Index) forklessCause(aID, bID hash.Event) bool {
	// TODO add logic to filter out cheaters
	vi.Log.Crit("forklessCause NOT_IMPLEMENTED yet", "aID", aID.String())

	// check if Event A not found
//	if a == nil {
//		vi.Log.Crit("Event A not found", "event", aID.String())
//		return false
//	}


	// check if Event B not found
//	if b == nil {
//		vi.Log.Crit("Event B not found", "event", bID.String())
//		return false
//	}

	// check A doesn't observe any forks from B

	// check A observes that {QUORUM} non-cheater-validators observe B
	yes := vi.validators.NewCounter()

	return yes.HasQuorum()
}

// NoCheaters excludes events which are observed by selfParents as cheaters.
// Called by emitter to exclude cheater's events from potential parents list.
func (vi *Index) NoCheaters(selfParent *hash.Event, options hash.Events) hash.Events {
	if selfParent == nil {
		return options
	}

	filtered := make(hash.Events, 0, len(options))
	// TODO add logic to filter out cheaters
	vi.Log.Crit("NoCheaters NOT_IMPLEMENTED yet", "selfParent", selfParent.String())


	return filtered
}
