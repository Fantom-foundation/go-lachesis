package fetcher

import (
	"github.com/Fantom-foundation/go-lachesis/eventcheck/heavycheck"
	"github.com/Fantom-foundation/go-lachesis/inter"
)

// Checker is an interface that represents abstract logic for a checker object
type Checker interface {
	Start()
	Stop()
	Overloaded() bool
	Enqueue(events inter.Events, onValidated heavycheck.OnValidatedFn) error
}

// MockTestData is a mock Checker for a testing purposes
type MockTestData struct {
	events      inter.Events
	resultError error
	onValidated heavycheck.OnValidatedFn
}

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

// MockChecker is a mock checker made for tests only
type MockChecker struct {
	resultErrors error
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
