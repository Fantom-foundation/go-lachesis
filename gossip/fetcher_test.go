package gossip

import (
	"errors"
	"github.com/Fantom-foundation/go-lachesis/gossip/fetcher"
	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"reflect"
	"runtime"
	"testing"
	"time"
)

type TestValidatorsPubKeys struct {
	Epoch     idx.Epoch
	Addresses map[idx.StakerID]common.Address
}

var tBuff = new(tBuffer)

func newFetcher(fn fetcher.FilterInterestedFn, hc *fetcher.MockChecker, fc func(*inter.Event) error, t *testing.T) *fetcher.Fetcher {
	return fetcher.New(fetcher.Callback{
		PushEvent:      pushEventFn,
		OnlyInterested: fn,
		DropPeer:       dropPeerFn,
		FirstCheck:     fc,
		HeavyCheck: hc,
	})
}

type tBuffer struct {
	events map[string][]inter.Event
}

func (t *tBuffer) Flush() {
	t.events = make(map[string][]inter.Event, 1)
}

func pushEventFn(e *inter.Event, peer string) {
	tBuff.events[peer] = append(tBuff.events[peer], *e)
}

func filterInterestedFn(ids hash.Events) hash.Events {
	return ids
}

func filterInterestedFnEmpty(ids hash.Events) hash.Events {
	return nil
}

func eventsRequesterFn(hash.Events) error {
	return nil
}

func dropPeerFn(peer string) {

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
	fetchEvents      fetcher.EventsRequesterFn
	filterInterested fetcher.FilterInterestedFn
	checkFn          func(*inter.Event) error
}

func TestFetcher_Enqueue(t *testing.T) {
	var testCases = getEnqueueTestCases()
	runtime.GOMAXPROCS(1)

	for _, v := range testCases {
		mc := fetcher.NewMockChecker(v.mtdError)
		f := newFetcher(v.filterInterested, mc, v.checkFn, t)
		runTestEnqueue(f, v, t)
		tBuff.Flush()
	}
}

func getEnqueueTestCases() []enqueueTestCase {
	var testCases []enqueueTestCase
	var peers = []string{"peer"}
	var mtdErrors = []error{nil, errors.New("ban err")}
	var interEventsNum = []int{0, 1}
	var times = []time.Time{time.Now(), time.Unix(0, 0)}
	var fetchEventsFuncs = []fetcher.EventsRequesterFn{eventsRequesterFn}
	var filterFuncs = []fetcher.FilterInterestedFn{filterInterestedFn, filterInterestedFnEmpty}
	var checkFuncs = []func(*inter.Event) error{firstCheck, firstCheckWithErr}

	for _, peer := range peers {
		for _, errs := range mtdErrors {
			for _, iEventsNum := range interEventsNum {
				for _, t := range times {
					for _, fetchEvent := range fetchEventsFuncs {
						for _, checFn := range checkFuncs {
							for _, filterFunc := range filterFuncs {
								etc := enqueueTestCase{
									peer,
									errs,
									makeInterEvents(iEventsNum),
									t,
									fetchEvent,
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

func runTestEnqueue(f *fetcher.Fetcher, testData enqueueTestCase, t *testing.T) {
	f.Start()
	runtime.Gosched() // we schedule every time after a call to a function with goroutine

	err := f.Enqueue(testData.peer, testData.inEvents, testData.t, testData.fetchEvents)
	runtime.Gosched()

	f.Stop()
	runtime.Gosched()

	time.Sleep(time.Millisecond * time.Duration(len(testData.inEvents))) // we cannot watch channel queue thus it is a private var. so we just wait for goroutine to handle the queue

	if funcIsFirstCheckWithErr(testData.checkFn) {
		require.Equal(t, 0, len(tBuff.events))
		if len(testData.inEvents) > 0 && !funcIsInterestedFnEmpty(testData.filterInterested) {
			require.NotNil(t, err)
		}
		return
	}

	if funcIsInterestedFnEmpty(testData.filterInterested) {
		require.Equal(t, 0, len(tBuff.events))
		return
	}

	if testData.mtdError != nil {
		require.Equal(t, 0, len(tBuff.events))
		return
	}

	require.Equal(t, len(testData.inEvents), len(tBuff.events[testData.peer]))
	require.Nil(t, err)
}

func funcIsInterestedFnEmpty(f fetcher.FilterInterestedFn) bool {
	return reflect.ValueOf(f).Pointer() == reflect.ValueOf(filterInterestedFnEmpty).Pointer()
}

func funcIsFirstCheckWithErr(f func(*inter.Event) error) bool {
	return reflect.ValueOf(f).Pointer() == reflect.ValueOf(firstCheckWithErr).Pointer()
}
