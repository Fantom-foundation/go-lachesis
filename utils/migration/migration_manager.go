package migration

import (
	"github.com/pkg/errors"
)

// Implementation of run all required migrations
type Manager struct {
	lastMigration *Migration
	idProd        IdProducer
}

func NewManager(last *Migration, idProducer IdProducer) *Manager {
	return &Manager{
		lastMigration: last,
		idProd:        idProducer,
	}
}

// Run all required migrations
func (mm *Manager) Run() error {
	// Read migration id from dbKey
	curId := mm.idProd.GetId()

	// Search last executed transaction

	list := make([]*Migration, 0)
	for startMigration := mm.lastMigration; startMigration != nil && string(curId) != startMigration.Id(); startMigration = startMigration.Prev() {
		list = append(list, startMigration)
	}

	// Execute migrations from list in reverse order (first runed - last in list)
	for i := len(list) - 1; i >= 0; i-- {
		err := list[i].Run()
		if err != nil {
			return errors.Wrap(err, "migration: "+list[i].Id())
		}
		mm.idProd.SetId(list[i].Id())
	}

	return nil
}
