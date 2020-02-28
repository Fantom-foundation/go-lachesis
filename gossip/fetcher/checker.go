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
