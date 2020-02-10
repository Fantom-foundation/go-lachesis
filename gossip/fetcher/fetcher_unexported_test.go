package fetcher

import (
	"errors"
	"github.com/Fantom-foundation/go-lachesis/eventcheck/heavycheck"
	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/pos"
	"github.com/Fantom-foundation/go-lachesis/lachesis"
	"github.com/Fantom-foundation/go-lachesis/lachesis/genesis"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"math/big"
	"reflect"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

type TestValidatorsPubKeys struct {
	Epoch     idx.Epoch
	Addresses map[idx.StakerID]common.Address
}

type TestHeavyCheckReader struct {
	Addrs atomic.Value
}

func (r *TestHeavyCheckReader) GetEpochPubKeys() (map[idx.StakerID]common.Address, idx.Epoch) {
	auth := r.Addrs.Load().(*TestValidatorsPubKeys)

	return auth.Addresses, auth.Epoch
}

func filterInterestedFn(ids hash.Events) hash.Events {
	return ids
}

func emptyInterestedFn(ids hash.Events) hash.Events {
	return nil
}

type notifiationTestCase struct {
	peer               string
	events             hash.Events
	time               time.Time
	requesterFn        func(hash.Events) error
	announcesNum       int
	interestedFilterFn FilterInterestedFn
}

func testEventsRequesterFn(hash.Events) error {
	return nil
}

func testEventsRequesterFnWithErr(hash.Events) error {
	return errors.New("some error text")
}

func pushEventFn(e *inter.Event, peer string) {

}

func dropPeerFn(peer string) {}

func firstCheck(*inter.Event) error {
	return nil
}

func newFetcher(fn FilterInterestedFn) *Fetcher {
	net := lachesis.FakeNetConfig(genesis.FakeAccounts(0, 5, big.NewInt(0), pos.StakeToBalance(1)))
	heavyCheckReader := &TestHeavyCheckReader{}
	ledgerID := net.EvmChainConfig().ChainID

	hc := heavycheck.NewDefault(&net.Dag, heavyCheckReader, types.NewEIP155Signer(ledgerID))
	return New(Callback{
		PushEvent:      pushEventFn,
		OnlyInterested: fn,
		DropPeer:       dropPeerFn,
		FirstCheck:     firstCheck,
		HeavyCheck:     hc,
	})
}

func TestFetcher_Notification(t *testing.T) {
	runtime.GOMAXPROCS(1)
	var testCases = generateNotificationTestCases()

	for _, v := range testCases {
		f := newFetcher(v.interestedFilterFn)
		runTestNotification(f, v, t)
	}
}

func TestFetcher_NotArrived(t *testing.T) {
	runtime.GOMAXPROCS(1)
	//mc := NewMockChecker(nil)
	peer := "peer"
	f := newFetcher(filterInterestedFn)
	event1 := hash.FakeEvent()
	event2 := hash.FakeEvent()
	batchOnTime := announcesBatch{
		hashes:      []hash.Event{event1},
		time:        time.Now(),
		peer:        peer,
		fetchEvents: testEventsRequesterFn,
	}
	batchLate := announcesBatch{
		hashes:      []hash.Event{event2},
		time:        time.Unix(0, 0),
		peer:        peer,
		fetchEvents: testEventsRequesterFn,
	}

	announceOnTime := oneAnnounce{batch: &batchOnTime, i: 0}
	announceLate := oneAnnounce{batch: &batchLate, i: 0}
	f.Start()
	f.announced[event1] = []*oneAnnounce{&announceOnTime}
	f.announced[event2] = []*oneAnnounce{&announceLate}
	time.Sleep(time.Millisecond)
	_, ok := f.fetching[event1]
	require.True(t, ok)

	_, ok = f.fetching[event2]
	require.False(t, ok)

	require.False(t, f.Overloaded())
}

func generateNotificationTestCases() []notifiationTestCase {
	var testCases []notifiationTestCase
	var peerNames = []string{"", "peer"}
	var eventNumbers = []int{0, 12, maxAnnounceBatch + 1}
	var times = []time.Time{time.Now(), time.Unix(0, 0)}
	var requeserFns = []EventsRequesterFn{testEventsRequesterFn, testEventsRequesterFnWithErr}
	var announcesNums = []int{0, 12, hashLimit}
	var filterFns = []FilterInterestedFn{filterInterestedFn, emptyInterestedFn}

	for _, name := range peerNames {
		for _, eventNum := range eventNumbers {
			for _, t := range times {
				for _, requesterFn := range requeserFns {
					for _, announcesNum := range announcesNums {
						for _, filterFn := range filterFns {
							n := notifiationTestCase{name, hash.FakeEvents(eventNum), t, requesterFn, announcesNum, filterFn}
							testCases = append(testCases, n)
						}
					}
				}
			}
		}
	}
	return testCases
}

func runTestNotification(f *Fetcher, testData notifiationTestCase, t *testing.T) {
	f.Start()
	f.announces[testData.peer] = testData.announcesNum
	err := f.Notify(testData.peer, testData.events, testData.time, testData.requesterFn)
	require.Nil(t, err)
	runtime.Gosched()
	f.Stop()
	runtime.Gosched()

	checkState(f, testData, t)

}

func checkState(f *Fetcher, testData notifiationTestCase, t *testing.T) {
	require.True(t, f.announces[testData.peer] >= 0)
	require.True(t, f.announces[testData.peer] <= hashLimit)

	if testData.peer == "" {
		require.Equal(t, 0, len(f.fetching))
		return
	}

	if reflect.ValueOf(testData.interestedFilterFn).Pointer() == reflect.ValueOf(emptyInterestedFn).Pointer() {
		require.Equal(t, 0, len(f.fetching))
		return
	}

	if time.Since(testData.time) > fetchTimeout {
		require.Equal(t, 0, len(f.fetching))
		return
	}

	eventsMaxNum := len(testData.events)
	eventsActualNum := eventsMaxNum
	if eventsMaxNum > maxAnnounceBatch {
		eventsActualNum = maxAnnounceBatch
	}
	annTotalNum := eventsActualNum + testData.announcesNum
	if len(testData.events)+testData.announcesNum > hashLimit {
		if annTotalNum > hashLimit {
			require.Equal(t, 0, len(f.fetching))
			return
		}
	}

	if len(testData.events) != len(f.fetching) {
		l1 := len(testData.events)
		l2 := len(f.fetching)
		require.Equal(t, l1, l2)
	}

	for _, v := range testData.events {
		_, ok := f.fetching[v]
		require.True(t, ok)
	}

	require.Equal(t, eventsMaxNum+testData.announcesNum, f.announces[testData.peer])
}
