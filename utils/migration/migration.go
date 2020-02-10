package migration

import "github.com/Fantom-foundation/go-lachesis/hash"

const (
	hashAppName = "lachesis"
	hashSalt = "Heuhax&Walv9"
)

// Implementation for DbMigration object
type Migration struct {
	id string
	prev *Migration
	run func() error
}

func New(id string, prev *Migration, runFunc func() error) *Migration {
	if id == "" {
		id = hashAppName+"?"+hash.FromBytes([]byte(prev.id + hashSalt)).Hex()
	}

	return &Migration{
		id:   id,
		prev: prev,
		run:  runFunc,
	}
}

func (m *Migration) New(runFunc func() error) *Migration {
	return New("", m.prev, runFunc)
}

func (m *Migration) NewNamed(id string, runFunc func() error) *Migration {
	return New(id, m.prev, runFunc)
}

func (m *Migration) Id() string {
	return m.id
}

func (m *Migration) Prev() *Migration {
	return m.prev
}

func (m *Migration) Run() error {
	return m.run()
}
