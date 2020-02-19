package migration

import (
	"github.com/ethereum/go-ethereum/log"
)

// Migration is a migration step.
type Migration struct {
	name string
	exec func() error
	prev *Migration
}

// Begin with empty unique migration step.
func Begin(appName string) *Migration {
	return &Migration{
		name: appName,
	}
}

// Next creates next migration.
func (m *Migration) Next(name string, exec func() error) *Migration {
	if name == "" {
		panic("empty name")
	}

	if exec == nil {
		panic("empty exec")
	}

	return &Migration{
		name: name,
		exec: exec,
		prev: m,
	}
}

func (m *Migration) Exec(curr IdProducer) error {
	if m.exec == nil {
		// only 1st empty migration
		return nil
	}

	myId := m.name

	if curr.IsCurrent(myId) {
		return nil
	}

	err := m.prev.Exec(curr)
	if err != nil {
		return err
	}

	err = m.exec()
	if err != nil {
		log.Error("'"+m.name+"' migration failed", "err", err)
		return err
	}
	log.Warn("'" + m.name + "' migration has been applied")

	curr.SetId(myId)
	return nil
}
