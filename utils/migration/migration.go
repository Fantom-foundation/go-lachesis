package migration

import "github.com/Fantom-foundation/go-lachesis/hash"

var (
	hashAppName string
	hashSalt string
)

// Implementation for DbMigration object
type Migration struct {
	id string
	prev *Migration
	run func() error
}

func Init(appName, salt string) *Migration {
	hashAppName = appName
	hashSalt = salt

	return newNamed("init", nil, func()error{
		return nil
	})
}

func newAuto(prev *Migration, runFunc func() error) *Migration {
	return newNamed("", prev, runFunc)
}

func newNamed(id string, prev *Migration, runFunc func() error) *Migration {
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
	return newAuto(m.prev, runFunc)
}

func (m *Migration) NewNamed(id string, runFunc func() error) *Migration {
	return newNamed(id, m.prev, runFunc)
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
