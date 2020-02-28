package heavycheck

import "github.com/Fantom-foundation/go-lachesis/inter"

// TaskData is a struct that represents an input data for a check. Used also to validate events
type TaskData struct {
	Events inter.Events // events to validate
	Result []error      // resulting errors of events, nil if ok

	onValidated OnValidatedFn
}

// GetEvents is a getter function for events
func (t *TaskData) GetEvents() inter.Events {
	return t.Events
}

// GetResult is a getter function for result
func (t *TaskData) GetResult() []error {
	return t.Result
}

// GetOnValidatedFn is a getter function for onValidated callback
func (t *TaskData) GetOnValidatedFn() OnValidatedFn {
	return t.onValidated
}

// ArbitraryTaskData is an interface that represents task data getters
type ArbitraryTaskData interface {
	GetEvents() inter.Events
	GetResult() []error
	GetOnValidatedFn() OnValidatedFn
}
