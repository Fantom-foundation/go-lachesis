package migrations

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/Fantom-foundation/go-lachesis/utils/migration"
)

type MemIdProducer struct {
	lastId string
}

func (p *MemIdProducer) GetId() (string, error) {
	return string(p.lastId), nil
}

func (p *MemIdProducer) SetId(id string) error {
	p.lastId = id
	return nil
}


func TestList(t *testing.T) {
	testData := map[string]int{}

	t.Run("Success run migrations", func(t *testing.T) {
		list := migration.Init("lachesis-test", "123456")

		num := 1
		list = list.New(func()error{
			testData["migration1"] = num
			num++
			return nil
		}).New(func()error{
			testData["migration2"] = num
			num++
			return nil
		}).New(func()error{
			testData["migration3"] = num
			num++
			return nil
		}).New(func()error{
			testData["migration4"] = num
			num++
			return nil
		})

		idProducer := &MemIdProducer{}
		mgr := migration.NewManager(list, idProducer)
		err := mgr.Run()

		assert.NoError(t, err, "Error when run migration manager")

		assert.Equal(t, 1, testData["migration1"], "Bad value after run migration1")
		assert.Equal(t, 2, testData["migration2"], "Bad value after run migration2")
		assert.Equal(t, 3, testData["migration3"], "Bad value after run migration3")
		assert.Equal(t, 4, testData["migration4"], "Bad value after run migration4")

		/*
			Additional migrations after already executed
		*/
		list = list.New(func()error{
			testData["migration5"] = num
			num++
			return nil
		}).New(func()error{
			testData["migration6"] = num
			num++
			return nil
		})

		mgr = migration.NewManager(list, idProducer)
		err = mgr.Run()

		assert.NoError(t, err, "Error when run migration manager")

		// No dup run executed migrations
		assert.Equal(t, 1, testData["migration1"], "Duplicate run migration1")
		assert.Equal(t, 2, testData["migration2"], "Duplicate run migration2")
		assert.Equal(t, 3, testData["migration3"], "Duplicate after run migration3")
		assert.Equal(t, 4, testData["migration4"], "Duplicate after run migration4")

		// Run new migrations
		assert.Equal(t, 5, testData["migration5"], "Bad value after run migration5")
		assert.Equal(t, 6, testData["migration6"], "Bad value after run migration6")

		/*
			Additional migrations with modify data
		*/
		list = list.New(func()error{
			testData["migration1"] = 100
			testData["migration3"] = 100
			testData["migration5"] = 100
			return nil
		})

		mgr = migration.NewManager(list, idProducer)
		err = mgr.Run()

		assert.NoError(t, err, "Error when run migration manager")

		assert.Equal(t, 100, testData["migration1"], "Bad value after run modify migration1")
		assert.Equal(t, 2, testData["migration2"], "Bad value after run modify migration2")
		assert.Equal(t, 100, testData["migration3"], "Bad value after run modify migration3")
		assert.Equal(t, 4, testData["migration4"], "Bad value after run modify migration4")
		assert.Equal(t, 100, testData["migration5"], "Bad value after run modify migration5")
		assert.Equal(t, 6, testData["migration6"], "Bad value after run modify migration6")

		/*
			Additional migrations with delete data
		*/
		list = list.New(func()error{
			delete(testData, "migration1")
			delete(testData, "migration2")
			delete(testData, "migration3")
			delete(testData, "migration4")
			delete(testData, "migration5")
			delete(testData, "migration6")
			return nil
		})

		mgr = migration.NewManager(list, idProducer)
		err = mgr.Run()

		assert.NoError(t, err, "Error when run migration manager")

		assert.Empty(t, testData, "Detect data after clear migration")
	})

	t.Run("Failed run migrations", func(t *testing.T) {
		list := migration.Init("lachesis-test", "123456")

		num := 1
		lastGood := list.New(func()error{
			testData["migration1"] = num
			num++
			return nil
		}).New(func()error{
			testData["migration2"] = num
			num++
			return nil
		})

		afterBad := lastGood.New(func()error{
			testData["migration3"] = num
			num++
			return errors.New("test migration error")
		}).New(func()error{
			testData["migration4"] = num
			num++
			return nil
		})

		idProducer := &MemIdProducer{}
		mgr := migration.NewManager(afterBad, idProducer)
		err := mgr.Run()

		assert.Error(t, err, "Success run migration manager with error migrations")

		lastId, _ := idProducer.GetId()
		assert.Equal(t, lastGood.Id(), lastId, "Bad last id in idProducer after migration error")

		assert.Equal(t, 1, testData["migration1"], "Bad value after run migration1")
		assert.Equal(t, 2, testData["migration2"], "Bad value after run migration2")
		assert.Equal(t, 3, testData["migration3"], "Bad value after run migration3")

		assert.Empty(t, testData["migration4"], "Bad data for migration4 - should by empty")

		/*
			Continue with fixed transactions
		*/
		num = 3
		fixed := lastGood.New(func()error{
			testData["migration3"] = num
			num++
			return nil
		}).New(func()error{
			testData["migration4"] = num
			num++
			return nil
		})

		mgr = migration.NewManager(fixed, idProducer)
		err = mgr.Run()

		assert.NoError(t, err, "Error when run migration manager")

		assert.Equal(t, 1, testData["migration1"], "Bad value after run migration1")
		assert.Equal(t, 2, testData["migration2"], "Bad value after run migration2")
		assert.Equal(t, 3, testData["migration3"], "Bad value after run migration3")
		assert.Equal(t, 4, testData["migration4"], "Bad value after run migration4")
	})
}
