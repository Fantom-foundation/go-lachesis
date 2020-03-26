package parentscheck

import (
	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/lachesis"
	"github.com/Fantom-foundation/go-lachesis/vector"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestParentsCheck is a main testing func
func TestParentsCheck(t *testing.T) {
	lachesisConfigs := []*lachesis.DagConfig{
		nil,
		&lachesis.DagConfig{
			MaxParents:                0,
			MaxFreeParents:            0,
			MaxEpochBlocks:            0,
			MaxEpochDuration:          0,
			VectorClockConfig:         vector.IndexConfig{},
			MaxValidatorEventsInBlock: 0,
		},
		&lachesis.DagConfig{
			MaxParents:                1e10,
			MaxFreeParents:            1,
			MaxEpochBlocks:            20,
			MaxEpochDuration:          2000,
			VectorClockConfig:         vector.IndexConfig{},
			MaxValidatorEventsInBlock: 10,
		},
	}

	for _, lachesisConfig := range lachesisConfigs {
		checker := New(lachesisConfig)

		events, parentSet := makeTestEventsAndParents(true)
		for _, parents := range parentSet {
			for _, event := range events {
				err := checker.Validate(event, parents)
				checkErrorResponse(t, err, event, parents)
			}
		}

		events, parentSet = makeTestEventsAndParents(false)
		for _, parents := range parentSet {
			for _, event := range events {
				err := checker.Validate(event, parents)
				checkErrorResponse(t, err, event, parents)
			}
		}
	}
}

// checkErrorResponse resolves if returned error is expected
func checkErrorResponse(t *testing.T, err error, event *inter.Event, parents []*inter.EventHeaderData) {

	if len(event.Parents) != len(parents) {
		require.Equal(t, err, ErrIncorrectParents)
		return
	}

	maxLamport := idx.Lamport(0)
	for _, p := range parents {
		maxLamport = idx.MaxLamport(maxLamport, p.Lamport)
	}

	if event.Lamport != maxLamport+1 {
		require.Equal(t, err, ErrWrongLamport)
		return
	}

	if len(event.Parents.Set()) != len(event.Parents) {
		require.Equal(t, err, ErrDoubleParents)
		return
	}

	for i, p := range parents {
		if (p.Creator == event.Creator) != event.IsSelfParent(event.Parents[i]) {
			require.Equal(t, err, ErrWrongSelfParent)
			return
		}
	}

	if (event.Seq <= 1) != (event.SelfParent() == nil) {
		require.Equal(t, err, ErrWrongSeq)
		return
	}

	if event.SelfParent() != nil {
		selfParent := parents[0]
		if !event.IsSelfParent(selfParent.Hash()) {
			require.Equal(t, err, ErrWrongSelfParent)
			return
		}
		if event.Seq != selfParent.Seq+1 {
			require.Equal(t, err, ErrWrongSeq)
			return
		}

		if event.ClaimedTime <= selfParent.ClaimedTime {
			require.Equal(t, err, ErrPastTime)
			return
		}
	}

	require.Nil(t, err)
}

// makeParentsForTests creates parents for an event
func makeParentsForTests(num int) hash.Events {
	var hashEvents hash.Events
	var h common.Hash
	for i := num; i > 0; i-- {
		arrId := i % 32
		h[arrId] = h[arrId] + 1
		hashEvents = append(hashEvents, hash.Event(h))
	}
	return hashEvents
}

// makeTestEventsAndParents creates test data of events and parents
func makeTestEventsAndParents(valid bool) ([]*inter.Event, [][]*inter.EventHeaderData) {
	var events []*inter.Event
	var eventHeaderDataSets [][]*inter.EventHeaderData
	if valid {
		parentData := makeParentsForTests(2)
		events = []*inter.Event{
			{
				EventHeader: inter.EventHeader{
					EventHeaderData: inter.EventHeaderData{
						Version: 1,
						Creator: 1,
						Seq:     0,
						Parents: parentData,
						Epoch:   1,
						Frame:   1,
						Lamport: 2,
					}, Sig: nil,
				},
			},
			{
				EventHeader: inter.EventHeader{
					EventHeaderData: inter.EventHeaderData{
						Version: 1,
						Creator: 1,
						Seq:     0,
						Parents: parentData,
						Epoch:   1,
						Frame:   1,
						Lamport: 3,
					}, Sig: nil,
				},
			},
		}
		eventHeaderDataSets = [][]*inter.EventHeaderData{{
			{
				Version: 1,
				Seq:     0,
				Epoch:   1,
				Frame:   1,
				Lamport: 2,
			},
			{
				Version: 1,
				Seq:     0,
				Epoch:   1,
				Frame:   1,
				Lamport: 2,
			},
		}}
	} else {
		mainParent := (&inter.EventHeaderData{}).Hash()

		parentDatas := []hash.Events{makeParentsForTests(0), makeParentsForTests(2), {hash.Event{}, hash.Event{}}, {mainParent}}
		creators := []idx.StakerID{0, 1, 2}
		lamports := []idx.Lamport{0, 1, 3}
		claimedTimes := []inter.Timestamp{0, 10}
		seqs := []idx.Event{0, 1, 2}
		for _, parentData := range parentDatas {
			for _, creator := range creators {
				for _, lamport := range lamports {
					for _, seq := range seqs {
						for _, claimedTime := range claimedTimes {

							event := &inter.Event{
								EventHeader: inter.EventHeader{
									EventHeaderData: inter.EventHeaderData{
										Version:     1,
										Creator:     creator,
										Seq:         seq,
										ClaimedTime: claimedTime,
										Parents:     parentData,
										Epoch:       1,
										Frame:       1,
										Lamport:     lamport,
									}, Sig: nil,
								},
							}
							events = append(events, event)
						}
					}
				}
			}
		}

		parentNums := []int{
			0, 1, 2,
		}
		claimedTimes = append(claimedTimes, 1)

		for _, parentData := range parentDatas {
			for _, creator := range creators {
				for _, lamport := range lamports {
					for _, seq := range seqs {
						for _, claimedTime := range claimedTimes {
							for _, parentNum := range parentNums {
								var _eventHeaderDataSets []*inter.EventHeaderData
								for i := 0; i < parentNum; i++ {
									header := &inter.EventHeaderData{
										Version:     1,
										Parents:     parentData,
										Creator:     creator,
										Seq:         seq,
										ClaimedTime: claimedTime,
										Epoch:       1,
										Frame:       1,
										Lamport:     lamport,
									}
									_eventHeaderDataSets = append(_eventHeaderDataSets, header)
								}
								eventHeaderDataSets = append(eventHeaderDataSets, _eventHeaderDataSets)
							}
						}
					}
				}
			}
		}
		eventHeaderDataSets = append(eventHeaderDataSets, []*inter.EventHeaderData{&inter.EventHeaderData{}})
	}

	return events, eventHeaderDataSets
}
