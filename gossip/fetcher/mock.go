package fetcher

import (
	"sync"

	"github.com/Fantom-foundation/go-lachesis/eventcheck/heavycheck"
	"github.com/Fantom-foundation/go-lachesis/inter"
)

type (
	// MockTestData is a mock Checker for a testing purposes
	MockTestData struct {
		events      inter.Events
		resultError error
		onValidated heavycheck.OnValidatedFn
	}

	// MockChecker is a mock checker made for tests only
	MockChecker struct {
		resultErrors error
	}

	// SyncChecker calls checker.Enqueue synchronously for a testing purposes
	SyncChecker struct {
		checker Checker
	}
)

// NewMockTestData is a constructor for a mock test data
func NewMockTestData(resultError error) *MockTestData {
	return &MockTestData{
		resultError: resultError,
	}
}

// GetEvents is a getter for events
func (m *MockTestData) GetEvents() inter.Events {
	return m.events
}

// GetResult is a getter for result
func (m *MockTestData) GetResult() []error {
	var errors []error
	for range m.events {
		errors = append(errors, m.resultError)
	}

	return errors
}

// GetOnValidatedFn is a getter for onValidated callback
func (m *MockTestData) GetOnValidatedFn() heavycheck.OnValidatedFn {
	return m.onValidated
}

// NewMockChecker is a constructor for a mock checker
func NewMockChecker(resultError error) *MockChecker {
	return &MockChecker{resultErrors: resultError}
}

// Start is an empty implementation
func (m *MockChecker) Start() {

}

// Stop is an empty implementation
func (m *MockChecker) Stop() {

}

// Overloaded is an empty implementation
func (m *MockChecker) Overloaded() bool {
	return false
}

// Enqueue made for tests only
func (m *MockChecker) Enqueue(events inter.Events, onValidated heavycheck.OnValidatedFn) error {
	const maxBatch = 4

	for start := 0; start < len(events); start += maxBatch {
		end := len(events)
		if end > start+maxBatch {
			end = start + maxBatch
		}
		op := &MockTestData{
			events:      events[start:end],
			onValidated: onValidated,
			resultError: m.resultErrors,
		}
		onValidated(op)
	}
	return nil
}

func (s *SyncChecker) Start() {
	s.checker.Start()
}

func (s *SyncChecker) Stop() {
	s.checker.Stop()
}

func (s *SyncChecker) Overloaded() bool {
	return s.checker.Overloaded()
}

func (s *SyncChecker) Enqueue(events inter.Events, onValidated heavycheck.OnValidatedFn) error {
	var wg sync.WaitGroup
	wg.Add(len(events))
	err := s.checker.Enqueue(events, func(data heavycheck.ArbitraryTaskData) {
		for range data.GetResult() {
			wg.Done()
		}
		onValidated(data)
	})

	wg.Wait()
	return err
}
