package migration

import (
	"crypto/sha256"
	"fmt"

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

// ID is an uniq migration's id.
func (m *Migration) Id() string {
	digest := sha256.New()
	if m.prev != nil {
		digest.Write([]byte(m.prev.Id()))
	}
	digest.Write([]byte(m.name))

	bytes := digest.Sum(nil)
	return fmt.Sprintf("%x", bytes)
}

func (m *Migration) Exec(curr IdStore) error {
	if m.exec == nil {
		// only 1st empty migration
		return nil
	}

	myId := m.Id()

	if curr.GetId() == myId {
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
