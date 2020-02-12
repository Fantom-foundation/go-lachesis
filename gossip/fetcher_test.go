package gossip

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/Fantom-foundation/go-lachesis/eventcheck/heavycheck"
	"github.com/Fantom-foundation/go-lachesis/gossip/fetcher"
	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
)

type TestValidatorsPubKeys struct {
	Epoch     idx.Epoch
	Addresses map[idx.StakerID]common.Address
}

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

func filterInterestedFn(ids hash.Events) hash.Events {
	return ids
}

func filterInterestedFnEmpty(ids hash.Events) hash.Events {
	return nil
}

func firstCheck(*inter.Event) error {
	return nil
}

func firstCheckWithErr(*inter.Event) error {
	return errors.New("some err")
}

type enqueueTestCase struct {
	peer             string
	mtdError         error
	inEvents         inter.Events
	t                time.Time
	filterInterested fetcher.FilterInterestedFn
	checkFn          func(*inter.Event) error
}

func TestFetcher_Enqueue(t *testing.T) {
	var testCases = getEnqueueTestCases()

	for _, v := range testCases {
		mc := fetcher.NewMockChecker(v.mtdError)
		fworld := newFetcherWorld(v.filterInterested, v.checkFn, mc)
		f := fetcher.New(fworld.Callback)
		runTestEnqueue(t, f, v, fworld)
	}
}

func getEnqueueTestCases() []enqueueTestCase {
	var testCases []enqueueTestCase
	var peers = []string{"peer"}
	var mtdErrors = []error{nil, errors.New("ban err")}
	var interEventsNum = []int{0, 1}
	var times = []time.Time{time.Now(), time.Unix(0, 0)}
	var filterFuncs = []fetcher.FilterInterestedFn{filterInterestedFn, filterInterestedFnEmpty}
	var checkFuncs = []func(*inter.Event) error{firstCheck, firstCheckWithErr}

	for _, peer := range peers {
		for _, errs := range mtdErrors {
			for _, iEventsNum := range interEventsNum {
				for _, t := range times {
					for _, checFn := range checkFuncs {
						for _, filterFunc := range filterFuncs {
							etc := enqueueTestCase{
								peer,
								errs,
								makeInterEvents(iEventsNum),
								t,
								filterFunc,
								checFn,
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

	if funcIsFirstCheckWithErr(testData.checkFn) {
		require.Equal(t, 0, len(fworld.events))
		if len(testData.inEvents) > 0 && !funcIsInterestedFnEmpty(testData.filterInterested) {
			require.NotNil(t, err)
		}
		return
	}

	if funcIsInterestedFnEmpty(testData.filterInterested) {
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

func funcIsInterestedFnEmpty(f fetcher.FilterInterestedFn) bool {
	return reflect.ValueOf(f).Pointer() == reflect.ValueOf(filterInterestedFnEmpty).Pointer()
}

func funcIsFirstCheckWithErr(f func(*inter.Event) error) bool {
	return reflect.ValueOf(f).Pointer() == reflect.ValueOf(firstCheckWithErr).Pointer()
}
