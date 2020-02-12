package gossip

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Fantom-foundation/go-lachesis/eventcheck/heavycheck"
	"github.com/Fantom-foundation/go-lachesis/gossip/fetcher"
	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
)

// fetcherWorld controls fetcher's environment for tests
type fetcherWorld struct {
	hChecker *fetcher.MockChecker
	events   map[string][]inter.Event
	waiting  chan struct{}

	fetcher.Callback
}

// newFetcherWorld constructor
func newFetcherWorld(
	interestFilter fetcher.FilterInterestedFn,
	firstChecker func(*inter.Event) error,
	heavyChecker *fetcher.MockChecker,
) *fetcherWorld {
	w := &fetcherWorld{
		hChecker: heavyChecker,
		events:   make(map[string][]inter.Event),
		waiting:  make(chan struct{}, 1),
	}

	w.Callback = fetcher.Callback{
		OnlyInterested: func(ids hash.Events) hash.Events {
			intrested := interestFilter(ids)
			for i := len(ids) - len(intrested); i > 0; i-- {
				w.waited()
			}
			return intrested
		},
		FirstCheck: func(e *inter.Event) error {
			err := firstChecker(e)
			if err != nil {
				w.waited()
			}
			return err
		},
		HeavyCheck: fetcher.Checker(w),
		PushEvent: func(e *inter.Event, peer string) {
			w.events[peer] = append(w.events[peer], *e)
			w.waited()
		},
		DropPeer: func(peer string) {
		},
	}

	return w
}

func (w *fetcherWorld) waited() {
	w.waiting <- struct{}{}
}

func (w *fetcherWorld) WaitFor(count int) bool {
	for i := 0; i < count; i++ {
		select {
		case <-w.waiting:
		case <-time.After(4 * time.Second):
			return false
		}
	}
	return true
}

func (w *fetcherWorld) Fetch(hash.Events) error {
	w.waited()
	return nil
}

// Start implements fetcher.Checker interface
func (w *fetcherWorld) Start() {
	w.hChecker.Start()
}

// Stop implements fetcher.Checker interface
func (w *fetcherWorld) Stop() {
	w.hChecker.Stop()
}

// Overloaded implements fetcher.Checker interface
func (w *fetcherWorld) Overloaded() bool {
	return w.hChecker.Overloaded()
}

// Enqueue implements fetcher.Checker interface
func (w *fetcherWorld) Enqueue(events inter.Events, onValidated heavycheck.OnValidatedFn) error {
	return w.hChecker.Enqueue(events, func(data heavycheck.ArbitraryTaskData) {
		res := data.GetResult()
		onValidated(data)
		for _, err := range res {
			if err != nil {
				w.waited()
			}
		}
	})
}

type enqueueTestCase struct {
	peer             string
	mtdError         error
	inEvents         inter.Events
	t                time.Time
	allInterested    bool
	firstCheckResult error
}

func TestFetcher_Enqueue(t *testing.T) {
	var testCases = getEnqueueTestCases()

	for _, v := range testCases {
		mc := fetcher.NewMockChecker(v.mtdError)

		fworld := newFetcherWorld(
			interested(v.allInterested),
			firstCheck(v.firstCheckResult),
			mc)
		f := fetcher.New(fworld.Callback)
		runTestEnqueue(t, f, v, fworld)
	}
}

func interested(all bool) fetcher.FilterInterestedFn {
	return func(ids hash.Events) hash.Events {
		if all {
			return ids
		}
		return nil
	}
}

func firstCheck(err error) func(*inter.Event) error {
	return func(*inter.Event) error {
		return err
	}
}

func getEnqueueTestCases() []enqueueTestCase {
	var testCases []enqueueTestCase
	var peers = []string{"peer"}
	var mtdErrors = []error{nil, errors.New("ban err")}
	var interEventsNum = []int{0, 1}
	var times = []time.Time{time.Now(), time.Unix(0, 0)}

	for _, peer := range peers {
		for _, errs := range mtdErrors {
			for _, iEventsNum := range interEventsNum {
				for _, t := range times {
					for _, firstCheckRes := range []error{nil, errors.New("some err")} {
						for _, allInterested := range []bool{true, false} {
							etc := enqueueTestCase{
								peer,
								errs,
								makeInterEvents(iEventsNum),
								t,
								allInterested,
								firstCheckRes,
							}
							testCases = append(testCases, etc)
						}
					}
				}
			}
		}
	}
	return testCases
}

func makeInterEvents(n int) inter.Events {
	var e inter.Events
	for n > 0 {
		event := inter.NewEvent()
		e = append(e, event)
		n--
	}
	return e
}

func runTestEnqueue(t *testing.T, f *fetcher.Fetcher, testData enqueueTestCase, fworld *fetcherWorld) {
	f.Start()
	defer f.Stop()

	err := f.Enqueue(testData.peer, testData.inEvents, testData.t, fworld.Fetch)

	require.True(t,
		fworld.WaitFor(len(testData.inEvents)),
		"not all the events were processed")

	if testData.firstCheckResult != nil {
		require.Equal(t, 0, len(fworld.events))
		if len(testData.inEvents) > 0 && testData.allInterested {
			require.NotNil(t, err)
		}
		return
	}

	if !testData.allInterested {
		require.Equal(t, 0, len(fworld.events))
		return
	}

	if testData.mtdError != nil {
		require.Equal(t, 0, len(fworld.events))
		return
	}

	require.Equal(t, len(testData.inEvents), len(fworld.events[testData.peer]))
	require.Nil(t, err)
}
