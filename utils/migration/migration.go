package migration

import (
	"crypto/sha256"
	"fmt"
)

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
		digest := sha256.New()
		digest.Write([]byte(prev.id + hashSalt))
		bytes := digest.Sum(nil)
		id = fmt.Sprintf("%s?%x", hashAppName, bytes)
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
