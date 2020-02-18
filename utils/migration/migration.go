package migration

import (
	"crypto/sha256"
	"fmt"
)

// Migration is a migration step.
type Migration struct {
	name string
	exec func() error
	prev *Migration
}

// Begin with empty unique migration step.
// Use it in optional.
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

func (m *Migration) Prev() *Migration {
	return m.prev
}

func (m *Migration) Run() error {
	if m.exec == nil {
		return nil
	}

	return m.exec()
}
