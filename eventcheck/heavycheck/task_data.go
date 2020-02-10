package heavycheck

import "github.com/Fantom-foundation/go-lachesis/inter"

// a struct that represents an input data for a check. Used also to validate events
type TaskData struct {
	Events inter.Events // events to validate
	Result []error      // resulting errors of events, nil if ok

	onValidated OnValidatedFn
}

// getter function for events
func (t *TaskData) GetEvents() inter.Events {
	return t.Events
}

// getter function for result
func (t *TaskData) GetResult() []error {
	return t.Result
}

// getter function for onValidated callback
func (t *TaskData) GetOnValidatedFn() OnValidatedFn {
	return t.onValidated
}

// an interface that represents task data getters
type ArbitraryTaskData interface {
	GetEvents() inter.Events
	GetResult() []error
	GetOnValidatedFn() OnValidatedFn
}
