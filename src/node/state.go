package node

import (
	"context"
	"sync"
)

// LimitGoFunc the number of simultaneously working functions that change the
// state of a graph.
const LimitGoFunc = 1000

const (
	// Gossiping is the initial state of a Lachesis node.
	Gossiping state = iota
	// CatchingUp is the fast forward state
	CatchingUp
	// Shutdown is the shut down state
)

type state int

type nodeState struct {
	cond         *sync.Cond
	ctx          context.Context
	getStateChan chan state
	limit        int
	setStateChan chan state
	state        state
	tickets      chan func()
	units        chan struct{}
}

func newNodeState(ctx context.Context, limit int) *nodeState {
	ns := &nodeState{
		cond:         sync.NewCond(&sync.Mutex{}),
		ctx:          ctx,
		getStateChan: make(chan state),
		limit:        limit,
		setStateChan: make(chan state),
		units:        make(chan struct{}, limit),
	}

	for i := 0; i < limit; i++ {
		ns.units <- struct{}{}
	}

	go ns.mtx()
	return ns
}

func (s state) String() string {
	switch s {
	case Gossiping:
		return "Gossiping"
	case CatchingUp:
		return "CatchingUp"
	default:
		return "Unknown"
	}
}

func (s *nodeState) mtx() {
	for {
		select {
		case s.state = <-s.setStateChan:
		case s.getStateChan <- s.state:
		}
	}
}

func (s *nodeState) goFunc(fu func()) {
	select {
	case <-s.ctx.Done():
		return
	default:
		// Exchanges units for functions.
		<-s.units
		go func() {
			fu()
			// Returns a unit to channel.
			s.units <- struct{}{}

			s.cond.L.Lock()
			s.cond.Broadcast()
			s.cond.L.Unlock()
		}()
	}
}

func (s *nodeState) waitRoutines() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	for {
		if len(s.units) < s.limit {
			s.cond.Wait()
			continue
		}
		break
	}
}

func (s *nodeState) getState() state {
	return <-s.getStateChan
}

func (s *nodeState) setState(state state) {
	s.setStateChan <- state
}
