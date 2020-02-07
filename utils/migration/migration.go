package migration

// Implementation for DbMigration object
type Migration struct {
	id string
	prev *Migration
	run func() error
}

func New(id string, prev *Migration, runFunc func() error) *Migration {
	return &Migration{
		id:   id,
		prev: prev,
		run:  runFunc,
	}
}

func (m *Migration) New(id string, runFunc func() error) *Migration {
	return &Migration{
		id:   id,
		prev: m,
		run:  runFunc,
	}
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
