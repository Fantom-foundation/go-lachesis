package gaspowercheck

import (
	"github.com/Fantom-foundation/go-lachesis/eventcheck/epochcheck"
	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/pos"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestDagReader is a dagReader implementation made for testing purposes
type TestDagReader struct {
	Ctx ValidationContext
}

// GetValidationContext is just an implementation
func (t TestDagReader) GetValidationContext() *ValidationContext {
	return &t.Ctx
}

// makeTestDagReaders creates dag readers for testing purposes
func makeTestDagReaders() []TestDagReader {
	var dagReaders []TestDagReader
	//type ValidationContext struct {
	//	Epoch                idx.Epoch
	//	Configs              [2]Config
	//	Validators           *pos.Validators
	//	PrevEpochLastHeaders inter.HeadersByCreator
	//	PrevEpochEndTime     inter.Timestamp
	//	PrevEpochRefunds     map[idx.StakerID]uint64
	//}

	epochs := []idx.Epoch{0, 1, 10}
	configIdxs := []int{-1, 0, 1, 10, 120}
	configAllocPerSpecs := []uint64{0, 1, 10}
	configMaxAllocPeriods := []inter.Timestamp{0, 10, 100}
	configStartupAllocPeriods := []inter.Timestamp{0, 10, 100}
	configMinStartupGasSet := []uint64{0, 1, 10}
	var configsSet [][2]Config

	for range configIdxs {
		for _, configAllocPerSpec := range configAllocPerSpecs {
			for _, configMaxAllocPeriod := range configMaxAllocPeriods {
				for _, configStartupAllocPeriod := range configStartupAllocPeriods {
					for _, configMinStartupGas := range configMinStartupGasSet {
						cfgs := [2]Config{
							{
								0,
								configAllocPerSpec,
								configMaxAllocPeriod,
								configStartupAllocPeriod,
								configMinStartupGas,
							},
							{
								1,
								configAllocPerSpec,
								configMaxAllocPeriod,
								configStartupAllocPeriod,
								configMinStartupGas,
							},
						}
						configsSet = append(configsSet, cfgs)
					}
				}
			}
		}
	}

	valBuilder := pos.NewBuilder()
	valBuilder.Set(1, 100)

	validators := valBuilder.Build()

	hbc1 := inter.HeadersByCreator{}
	hbc1[1] = &inter.EventHeaderData{
		Version: 1,
		Creator: 1,
		Seq:     0,
		Epoch:   1,
		Frame:   1,
		Lamport: 2,
	}
	per := map[idx.StakerID]uint64{}
	per[1] = 10
	prevEpochLastHeaders := []inter.HeadersByCreator{nil, {}, hbc1}
	prevEpochRefunds := []map[idx.StakerID]uint64{nil, per}
	for _, epoch := range epochs {
		for _, configs := range configsSet {
			for _, prevEpochLastHeader := range prevEpochLastHeaders {
				for _, prevEpochRefund := range prevEpochRefunds {
					dr := TestDagReader{}
					dr.Ctx = ValidationContext{
						Epoch:                epoch,
						Configs:              configs,
						Validators:           validators,
						PrevEpochLastHeaders: prevEpochLastHeader,
						PrevEpochEndTime:     0,
						PrevEpochRefunds:     prevEpochRefund,
					}
					dagReaders = append(dagReaders, dr)
				}
			}
		}
	}
	return dagReaders
}

// TestGasPowerCheck is a main testing func
func TestGasPowerCheck(t *testing.T) {
	dagReaders := makeTestDagReaders()
	for _, dagReader := range dagReaders {
		checker := New(dagReader)
		events, parentSet := makeTestEventsAndParents(true)
		for _, parents := range parentSet {
			for _, event := range events {
				for _, parent := range parents {

					err := checker.Validate(event, parent)
					checkErrorResponse(t, err, event, parent, checker)
				}
			}
		}
	}
}

// TestCalcGasPower runs special tests for a CalcGasPower func
func TestCalcGasPower(t *testing.T) {
	events := makeEventsForTestCalcGasPower()
	dagReaders := makeTestDagReaders()
	var gasPower inter.GasPowerLeft
	var err error

	for _, dagReader := range dagReaders {
		for _, event := range events {
			checker := New(dagReader)
			if event.SelfParent() != nil {
				gasPower, err = checker.CalcGasPower(&event.EventHeaderData, &event.EventHeaderData)
				continue
			} else {
				gasPower, err = checker.CalcGasPower(&event.EventHeaderData, nil)
			}
			validateCalcGpResponse(t, &event, gasPower, err, checker)
		}
	}
}

// validateCalcGpResponse checks if returned data during a test is valid
func validateCalcGpResponse(t *testing.T, e *inter.Event ,gp inter.GasPowerLeft, err error, checker *Checker) {
	ctx := checker.reader.GetValidationContext()
	// check that all the data is for the same epoch
	if ctx.Epoch != e.Epoch {
		require.Equal(t, gp, inter.GasPowerLeft{})
		require.Equal(t, err, epochcheck.ErrNotRelevant)
		return
	}

	var expectedGPLeft inter.GasPowerLeft
	for i := range ctx.Configs {
		if e.SelfParent() != nil {
			expectedGPLeft.Gas[i] = calcGasPower(&e.EventHeaderData, &e.EventHeaderData, ctx, &ctx.Configs[i])
			continue
		} else {
			expectedGPLeft.Gas[i] = calcGasPower(&e.EventHeaderData, nil, ctx, &ctx.Configs[i])
		}
	}
	require.Equal(t, expectedGPLeft, gp)
	require.Equal(t, err, nil)
}

// makeEventsForTestCalcGasPower creates special events for calcGasPower func
func makeEventsForTestCalcGasPower() []inter.Event {
	parentData := makeParentsForTests(2)
	commonParentHash := (&inter.Event{
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
	}).Hash()

	return []inter.Event{
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
					Seq:     2,
					Parents: hash.Events{commonParentHash},
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
					Seq:     2,
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
					Creator: 0,
					Seq:     2,
					Epoch:   1,
					Frame:   1,
					Lamport: 2,
				}, Sig: nil,
			},
		},
		{},
	}
}

// checkErrorResponse checks that returned error is expected
func checkErrorResponse(t *testing.T, err error, event *inter.Event, parent *inter.EventHeaderData, checker *Checker) {
	gasPowers, err2 := checker.CalcGasPower(&event.EventHeaderData, parent)
	if err2 != nil {
		require.Equal(t, err, err2)
		return
	}

	for i := range gasPowers.Gas {
		if event.GasPowerLeft.Gas[i]+event.GasPowerUsed != gasPowers.Gas[i] { // GasPowerUsed is checked in basic_check
			require.Equal(t, err, ErrWrongGasPowerLeft)
		}
	}
	require.Equal(t, err, nil)
	return
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

// makeTestEventsAndParents creates events and parents for tests
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
		eventHeaderDataSets = [][]*inter.EventHeaderData{
			{
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
			},
		}
	} else {
		mainParent := (&inter.EventHeaderData{}).Hash()

		parentDatas := []hash.Events{makeParentsForTests(0), makeParentsForTests(2), {mainParent}}
		creators := []idx.StakerID{0, 1}
		lamports := []idx.Lamport{0, 1}
		gaspowerLeftLongs := []uint64{0, 10, 1e8}
		gaspowerLeftShorts := []uint64{0, 10, 1e8}
		claimedTimes := []inter.Timestamp{0, 10}
		seqs := []idx.Event{0, 1, 2}
		for _, parentData := range parentDatas {
			for _, creator := range creators {
				for _, lamport := range lamports {
					for _, seq := range seqs {
						for _, claimedTime := range claimedTimes {
							for _, gaspowerLeftLong := range gaspowerLeftLongs {
								for _, gaspowerLeftShort := range gaspowerLeftShorts {

									event := &inter.Event{
										EventHeader: inter.EventHeader{
											EventHeaderData: inter.EventHeaderData{
												Version:      1,
												Creator:      creator,
												Seq:          seq,
												ClaimedTime:  claimedTime,
												Parents:      parentData,
												Epoch:        1,
												GasPowerLeft: inter.GasPowerLeft{Gas: [2]uint64{gaspowerLeftLong, gaspowerLeftShort}},
												Frame:        1,
												Lamport:      lamport,
											}, Sig: nil,
										},
									}
									events = append(events, event)
								}
							}
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

