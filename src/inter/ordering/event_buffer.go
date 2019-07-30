package ordering

import (
	"fmt"
	"time"

	"github.com/Fantom-foundation/go-lachesis/src/hash"
	"github.com/Fantom-foundation/go-lachesis/src/inter"
)

const expiration = 1 * time.Hour

type (
	// Event is a inter.Event and data for ordering purpose.
	Event struct {
		*inter.Event

		parents map[hash.Event]*inter.Event
		expired time.Time
	}

	// Callback is a set of EventBuffer()'s args.
	Callback struct {
		Process func(*inter.Event)
		Drop    func(*inter.Event, error)
		Exists  func(hash.Event) *inter.Event
	}
)

// EventBuffer validates, bufferizes and drops() or processes() pushed() event
// if all their parents exists().
func EventBuffer(callback Callback) (push func(*inter.Event), get func() map[string]*Event) {

	var (
		incompletes = make(map[string]*Event)
		lastGC      = time.Now()
		onNewEvent  func(e *Event)
	)

	onNewEvent = func(e *Event) {
		reffs := newRefsValidator(e.Event)
		ltime := newLamportTimeValidator(e.Event)

		// fill event's parents index or hold it as incompleted
		for pHash := range e.Parents {
			if pHash.IsZero() {
				// first event of node
				if err := reffs.AddUniqueParent(e.Creator); err != nil {
					callback.Drop(e.Event, err)
					return
				}
				if e.SelfParent != pHash {
					callback.Drop(e.Event, fmt.Errorf("invalid SelfParent"))
					return
				}
				if err := ltime.AddParentTime(0); err != nil {
					callback.Drop(e.Event, err)
					return
				}
				continue
			}
			parent := e.parents[pHash]
			if parent == nil {
				parent = callback.Exists(pHash)
				if parent == nil {
					key := e.Creator.Hex() + string(e.Seq)
					incompletes[key] = e
					return
				}
				e.parents[pHash] = parent
			}
			if err := reffs.AddUniqueParent(parent.Creator); err != nil {
				callback.Drop(e.Event, err)
				return
			}
			if parent.Creator == e.Creator && e.SelfParent != pHash {
				callback.Drop(e.Event, fmt.Errorf("invalid SelfParent"))
				return
			}
			if err := ltime.AddParentTime(parent.LamportTime); err != nil {
				callback.Drop(e.Event, err)
				return
			}
		}
		if err := reffs.CheckSelfParent(); err != nil {
			callback.Drop(e.Event, err)
			return
		}
		if err := ltime.CheckSequential(); err != nil {
			callback.Drop(e.Event, err)
			return
		}

		// parents OK
		callback.Process(e.Event)

		// now child events may become complete, check it again
		for key, child := range incompletes {
			if parent, ok := child.parents[e.Hash()]; ok && parent == nil {
				child.parents[e.Hash()] = e.Event
				delete(incompletes, key)
				onNewEvent(child)
			}
		}

	}

	get = func() map[string]*Event {
		return incompletes
	}

	push = func(e *inter.Event) {
		if callback.Exists(e.Hash()) != nil {
			callback.Drop(e, fmt.Errorf("event %s had received already", e.Hash().String()))
			return
		}

		w := &Event{
			Event:   e,
			parents: make(map[hash.Event]*inter.Event, len(e.Parents)),
			expired: time.Now().Add(expiration),
		}
		for parentHash := range e.Parents {
			w.parents[parentHash] = nil
		}
		onNewEvent(w)

		// GC
		if time.Now().Add(-expiration / 4).Before(lastGC) {
			return
		}
		lastGC = time.Now()
		limit := time.Now().Add(-expiration)
		for k, e := range incompletes {
			if e.expired.Before(limit) {
				delete(incompletes, k)
			}
		}
	}

	return
}
