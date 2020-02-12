package fetcher

import (
	"errors"
	"math/big"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/Fantom-foundation/go-lachesis/eventcheck/heavycheck"
	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/pos"
	"github.com/Fantom-foundation/go-lachesis/lachesis"
	"github.com/Fantom-foundation/go-lachesis/lachesis/genesis"
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

type notifiationTestCase struct {
	peer          string
	events        hash.Events
	time          time.Time
	requester     error
	announcesNum  int
	allInterested bool
}

func newCallback() Callback {
	net := lachesis.FakeNetConfig(genesis.FakeAccounts(0, 5, big.NewInt(0), pos.StakeToBalance(1)))
	heavyCheckReader := &TestHeavyCheckReader{}
	ledgerID := net.EvmChainConfig().ChainID
	hc := heavycheck.NewDefault(&net.Dag, heavyCheckReader, types.NewEIP155Signer(ledgerID))

	return Callback{
		OnlyInterested: func(ids hash.Events) hash.Events {
			return ids
		},
		FirstCheck: func(*inter.Event) error {
			return nil
		},
		HeavyCheck: &SyncChecker{hc},
		PushEvent:  func(e *inter.Event, peer string) {},
		DropPeer:   func(peer string) {},
	}
}

func TestFetcher_Notification(t *testing.T) {
	testCases := generateNotificationTestCases()

	for _, v := range testCases {
		c := newCallback()
		if !v.allInterested {
			c.OnlyInterested = func(ids hash.Events) hash.Events {
				return nil
			}
		}
		f := New(c)
		runTestNotification(f, v, t)
	}
}

func TestFetcher_NotArrived(t *testing.T) {
	require := require.New(t)

	peer := "peer"
	c := newCallback()
	f := New(c)

	event1 := hash.FakeEvent()
	event2 := hash.FakeEvent()
	batchOnTime := announcesBatch{
		hashes:      []hash.Event{event1},
		time:        time.Now(),
		peer:        peer,
		fetchEvents: func(hash.Events) error { return nil },
	}
	batchLate := announcesBatch{
		hashes:      []hash.Event{event2},
		time:        time.Unix(0, 0),
		peer:        peer,
		fetchEvents: func(hash.Events) error { return nil },
	}

	announceOnTime := oneAnnounce{batch: &batchOnTime, i: 0}
	announceLate := oneAnnounce{batch: &batchLate, i: 0}

	f.announced[event1] = []*oneAnnounce{&announceOnTime}
	f.announced[event2] = []*oneAnnounce{&announceLate}

	f.processNotifications()

	_, ok := f.fetching[event1]
	require.True(ok)

	_, ok = f.fetching[event2]
	require.False(ok)

	require.False(f.Overloaded())
}

func generateNotificationTestCases() []notifiationTestCase {
	var testCases []notifiationTestCase

	var eventNumbers = []int{0, 12, maxAnnounceBatch + 1}
	var times = []time.Time{time.Now(), time.Unix(0, 0)}
	var announcesNums = []int{0, 12, hashLimit}

	for _, name := range []string{"", "peer"} {
		for _, eventNum := range eventNumbers {
			for _, t := range times {
				for _, requester := range []error{nil, errors.New("some error text")} {
					for _, announcesNum := range announcesNums {
						for _, allInterested := range []bool{true, false} {
							n := notifiationTestCase{
								name,
								hash.FakeEvents(eventNum),
								t,
								requester,
								announcesNum,
								allInterested,
							}
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
	require := require.New(t)

	f.callbackHandler.HeavyCheck.Start()
	defer f.callbackHandler.HeavyCheck.Stop()

	f.announces[testData.peer] = testData.announcesNum
	err := f.Notify(testData.peer, testData.events, testData.time,
		func(hash.Events) error {
			return testData.requester
		})
	require.Nil(err)

	for len(f.notify) > 0 || len(f.inject) > 0 {
		f.processNotifications()
	}
	f.cleanupExpiredEvents()

	checkState(f, testData, t)
}

func checkState(f *Fetcher, testData notifiationTestCase, t *testing.T) {
	require := require.New(t)

	require.True(f.announces[testData.peer] >= 0)
	require.True(f.announces[testData.peer] <= hashLimit)

	if testData.peer == "" {
		require.Equal(0, len(f.fetching))
		return
	}

	if !testData.allInterested {
		require.Equal(0, len(f.fetching))
		return
	}

	if time.Since(testData.time) > fetchTimeout {
		require.Equal(0, len(f.fetching))
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
			require.Equal(0, len(f.fetching))
			return
		}
	}

	if len(testData.events) != len(f.fetching) {
		l1 := len(testData.events)
		l2 := len(f.fetching)
		require.Equal(l1, l2)
	}

	for _, v := range testData.events {
		_, ok := f.fetching[v]
		require.True(ok)
	}

	require.Equal(eventsMaxNum+testData.announcesNum, f.announces[testData.peer])
}
