package heavycheck

import "github.com/Fantom-foundation/go-lachesis/inter"

type TaskData struct {
	Events inter.Events // events to validate
	Result []error      // resulting errors of events, nil if ok

	onValidated OnValidatedFn
}

func (t *TaskData) GetEvents() inter.Events {
	return t.Events
}

func (t *TaskData) GetResult() []error {
	return t.Result
}

func (t *TaskData) GetOnValidatedFn() OnValidatedFn {
	return t.onValidated
}

type ArbitraryTaskData interface {
	GetEvents() inter.Events
	GetResult() []error
	GetOnValidatedFn() OnValidatedFn
}
