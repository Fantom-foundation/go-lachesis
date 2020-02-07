package fetcher

import (
	"github.com/Fantom-foundation/go-lachesis/eventcheck/heavycheck"
	"github.com/Fantom-foundation/go-lachesis/inter"
)

type Checker interface {
	Start()
	Stop()
	Overloaded() bool
	Enqueue(events inter.Events, onValidated heavycheck.OnValidatedFn) error
}

type MockTestData struct {
	events      inter.Events
	resultError error
	onValidated heavycheck.OnValidatedFn
}

func NewMockTestData(resultError error) *MockTestData {
	return &MockTestData{
		resultError: resultError,
	}
}

func (m *MockTestData) GetEvents() inter.Events {
	return m.events
}

func (m *MockTestData) GetResult() []error {
	var errors []error
	for range m.events {
		errors = append(errors, m.resultError)
	}

	return errors
}

func (m *MockTestData) GetOnValidatedFn() heavycheck.OnValidatedFn {
	return m.onValidated
}

type MockChecker struct {
	resultErrors error
}

func NewMockChecker(resultError error) *MockChecker {
	return &MockChecker{resultErrors: resultError}
}

func (m *MockChecker) Start() {

}

func (m *MockChecker) Stop() {

}

func (m *MockChecker) Overloaded() bool {
	return false
}

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
