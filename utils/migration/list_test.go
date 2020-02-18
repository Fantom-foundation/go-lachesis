package migration

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Fantom-foundation/go-lachesis/kvdb"
	"github.com/Fantom-foundation/go-lachesis/kvdb/flushable"
	"github.com/Fantom-foundation/go-lachesis/kvdb/leveldb"
	"github.com/Fantom-foundation/go-lachesis/kvdb/table"
)

type MemIdProducer struct {
	lastId string
}

func (p *MemIdProducer) GetId() string {
	return string(p.lastId)
}

func (p *MemIdProducer) SetId(id string) {
	p.lastId = id
}

func TestList(t *testing.T) {
	testData := map[string]int{}

	t.Run("Success run migrations", func(t *testing.T) {
		list := Init("lachesis-test", "123456")

		num := 1
		list = list.New(func() error {
			testData["migration1"] = num
			num++
			return nil
		}).New(func() error {
			testData["migration2"] = num
			num++
			return nil
		}).New(func() error {
			testData["migration3"] = num
			num++
			return nil
		}).New(func() error {
			testData["migration4"] = num
			num++
			return nil
		})

		idProducer := &MemIdProducer{}
		mgr := NewManager(list, idProducer)
		err := mgr.Run()

		assert.NoError(t, err, "Error when run migration manager")

		assert.Equal(t, 1, testData["migration1"], "Bad value after run migration1")
		assert.Equal(t, 2, testData["migration2"], "Bad value after run migration2")
		assert.Equal(t, 3, testData["migration3"], "Bad value after run migration3")
		assert.Equal(t, 4, testData["migration4"], "Bad value after run migration4")

		/*
			Additional migrations after already executed
		*/
		list = list.New(func() error {
			testData["migration5"] = num
			num++
			return nil
		}).New(func() error {
			testData["migration6"] = num
			num++
			return nil
		})

		mgr = NewManager(list, idProducer)
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
		list = list.New(func() error {
			testData["migration1"] = 100
			testData["migration3"] = 100
			testData["migration5"] = 100
			return nil
		})

		mgr = NewManager(list, idProducer)
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
		list = list.New(func() error {
			delete(testData, "migration1")
			delete(testData, "migration2")
			delete(testData, "migration3")
			delete(testData, "migration4")
			delete(testData, "migration5")
			delete(testData, "migration6")
			return nil
		})

		mgr = NewManager(list, idProducer)
		err = mgr.Run()

		assert.NoError(t, err, "Error when run migration manager")

		assert.Empty(t, testData, "Detect data after clear migration")
	})

	t.Run("Failed run migrations", func(t *testing.T) {
		list := Init("lachesis-test", "123456")

		num := 1
		lastGood := list.New(func() error {
			testData["migration1"] = num
			num++
			return nil
		}).New(func() error {
			testData["migration2"] = num
			num++
			return nil
		})

		afterBad := lastGood.New(func() error {
			testData["migration3"] = num
			num++
			return errors.New("test migration error")
		}).New(func() error {
			testData["migration4"] = num
			num++
			return nil
		})

		idProducer := &MemIdProducer{}
		mgr := NewManager(afterBad, idProducer)
		err := mgr.Run()

		assert.Error(t, err, "Success run migration manager with error migrations")

		lastId := idProducer.GetId()
		assert.Equal(t, lastGood.Id(), lastId, "Bad last id in idProducer after migration error")

		assert.Equal(t, 1, testData["migration1"], "Bad value after run migration1")
		assert.Equal(t, 2, testData["migration2"], "Bad value after run migration2")
		assert.Equal(t, 3, testData["migration3"], "Bad value after run migration3")

		assert.Empty(t, testData["migration4"], "Bad data for migration4 - should by empty")

		/*
			Continue with fixed transactions
		*/
		num = 3
		fixed := lastGood.New(func() error {
			testData["migration3"] = num
			num++
			return nil
		}).New(func() error {
			testData["migration4"] = num
			num++
			return nil
		})

		mgr = NewManager(fixed, idProducer)
		err = mgr.Run()

		assert.NoError(t, err, "Error when run migration manager")

		assert.Equal(t, 1, testData["migration1"], "Bad value after run migration1")
		assert.Equal(t, 2, testData["migration2"], "Bad value after run migration2")
		assert.Equal(t, 3, testData["migration3"], "Bad value after run migration3")
		assert.Equal(t, 4, testData["migration4"], "Bad value after run migration4")
	})

	t.Run("Success run migrations with DB", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "test-migrations")
		if err != nil {
			panic(fmt.Sprintf("can't create temporary directory %s: %v", dir, err))
		}
		disk := leveldb.NewProducer(dir)

		// open raw databases
		leveldb1 := disk.OpenDb("1")
		defer leveldb1.Drop()
		defer leveldb1.Close()

		flushableDB := flushable.Wrap(leveldb1)
		idTable := table.New(flushableDB, []byte("migration_id"))
		idProducer := kvdb.NewIdProducer(idTable)

		list := Init("lachesis-test-db", "654321")

		num := int64(1)
		list = list.New(func() error {
			err := flushableDB.Put([]byte("migration1"), []byte(strconv.FormatInt(num, 10)))
			if err != nil {
				return err
			}
			err = flushableDB.Flush()
			if err != nil {
				return err
			}
			num++
			return nil
		}).New(func() error {
			flushableDB.Put([]byte("migration2"), []byte(strconv.FormatInt(num, 10)))
			if err != nil {
				return err
			}
			err = flushableDB.Flush()
			if err != nil {
				return err
			}
			num++
			return nil
		}).New(func() error {
			flushableDB.Put([]byte("migration3"), []byte(strconv.FormatInt(num, 10)))
			if err != nil {
				return err
			}
			err = flushableDB.Flush()
			if err != nil {
				return err
			}
			num++
			return nil
		}).New(func() error {
			flushableDB.Put([]byte("migration4"), []byte(strconv.FormatInt(num, 10)))
			if err != nil {
				return err
			}
			err = flushableDB.Flush()
			if err != nil {
				return err
			}
			num++
			return nil
		})

		mgr := NewManager(list, idProducer)
		err = mgr.Run()

		assert.NoError(t, err, "Error when run migration manager")

		testData1, _ := flushableDB.Get([]byte("migration1"))
		testData2, _ := flushableDB.Get([]byte("migration2"))
		testData3, _ := flushableDB.Get([]byte("migration3"))
		testData4, _ := flushableDB.Get([]byte("migration4"))

		assert.Equal(t, []byte("1"), testData1, "Bad value after run migration1")
		assert.Equal(t, []byte("2"), testData2, "Bad value after run migration2")
		assert.Equal(t, []byte("3"), testData3, "Bad value after run migration3")
		assert.Equal(t, []byte("4"), testData4, "Bad value after run migration4")

		/*
			Additional migrations after already executed
		*/
		list = list.New(func() error {
			flushableDB.Put([]byte("migration5"), []byte(strconv.FormatInt(num, 10)))
			if err != nil {
				return err
			}
			err = flushableDB.Flush()
			if err != nil {
				return err
			}
			num++
			return nil
		}).New(func() error {
			flushableDB.Put([]byte("migration6"), []byte(strconv.FormatInt(num, 10)))
			if err != nil {
				return err
			}
			err = flushableDB.Flush()
			if err != nil {
				return err
			}
			num++
			return nil
		})

		mgr = NewManager(list, idProducer)
		err = mgr.Run()

		assert.NoError(t, err, "Error when run migration manager")

		testData1, _ = flushableDB.Get([]byte("migration1"))
		testData2, _ = flushableDB.Get([]byte("migration2"))
		testData3, _ = flushableDB.Get([]byte("migration3"))
		testData4, _ = flushableDB.Get([]byte("migration4"))
		testData5, _ := flushableDB.Get([]byte("migration5"))
		testData6, _ := flushableDB.Get([]byte("migration6"))

		// No dup run executed migrations
		assert.Equal(t, []byte("1"), testData1, "Duplicate run migration1")
		assert.Equal(t, []byte("2"), testData2, "Duplicate run migration2")
		assert.Equal(t, []byte("3"), testData3, "Duplicate after run migration3")
		assert.Equal(t, []byte("4"), testData4, "Duplicate after run migration4")

		// Run new migrations
		assert.Equal(t, []byte("5"), testData5, "Bad value after run migration5")
		assert.Equal(t, []byte("6"), testData6, "Bad value after run migration6")

		/*
			Additional migrations with modify data
		*/
		list = list.New(func() error {
			flushableDB.Put([]byte("migration1"), []byte("100"))
			if err != nil {
				return err
			}
			flushableDB.Put([]byte("migration3"), []byte("100"))
			if err != nil {
				return err
			}
			flushableDB.Put([]byte("migration5"), []byte("100"))
			if err != nil {
				return err
			}
			err = flushableDB.Flush()
			if err != nil {
				return err
			}
			num++
			return nil
		})

		mgr = NewManager(list, idProducer)
		err = mgr.Run()

		assert.NoError(t, err, "Error when run migration manager")

		testData1, _ = flushableDB.Get([]byte("migration1"))
		testData2, _ = flushableDB.Get([]byte("migration2"))
		testData3, _ = flushableDB.Get([]byte("migration3"))
		testData4, _ = flushableDB.Get([]byte("migration4"))
		testData5, _ = flushableDB.Get([]byte("migration5"))
		testData6, _ = flushableDB.Get([]byte("migration6"))

		assert.Equal(t, []byte("100"), testData1, "Bad value after run modify migration1")
		assert.Equal(t, []byte("2"), testData2, "Bad value after run modify migration2")
		assert.Equal(t, []byte("100"), testData3, "Bad value after run modify migration3")
		assert.Equal(t, []byte("4"), testData4, "Bad value after run modify migration4")
		assert.Equal(t, []byte("100"), testData5, "Bad value after run modify migration5")
		assert.Equal(t, []byte("6"), testData6, "Bad value after run modify migration6")

		/*
			Additional migrations with delete data
		*/
		list = list.New(func() error {
			err = flushableDB.Delete([]byte("migration1"))
			if err != nil {
				return err
			}
			err = flushableDB.Delete([]byte("migration2"))
			if err != nil {
				return err
			}
			err = flushableDB.Delete([]byte("migration3"))
			if err != nil {
				return err
			}
			err = flushableDB.Delete([]byte("migration4"))
			if err != nil {
				return err
			}
			err = flushableDB.Delete([]byte("migration5"))
			if err != nil {
				return err
			}
			err = flushableDB.Delete([]byte("migration6"))
			if err != nil {
				return err
			}
			err = flushableDB.Flush()
			if err != nil {
				return err
			}
			return nil
		})

		mgr = NewManager(list, idProducer)
		err = mgr.Run()

		assert.NoError(t, err, "Error when run migration manager")

		testData1, _ = flushableDB.Get([]byte("migration1"))
		testData2, _ = flushableDB.Get([]byte("migration2"))
		testData3, _ = flushableDB.Get([]byte("migration3"))
		testData4, _ = flushableDB.Get([]byte("migration4"))
		testData5, _ = flushableDB.Get([]byte("migration5"))
		testData6, _ = flushableDB.Get([]byte("migration6"))

		assert.Empty(t, testData1, "Detect data after clear migration1")
		assert.Empty(t, testData2, "Detect data after clear migration2")
		assert.Empty(t, testData3, "Detect data after clear migration3")
		assert.Empty(t, testData4, "Detect data after clear migration4")
		assert.Empty(t, testData5, "Detect data after clear migration5")
		assert.Empty(t, testData6, "Detect data after clear migration6")
	})

	t.Run("Failed run migrations with DB", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "test-migrations")
		if err != nil {
			panic(fmt.Sprintf("can't create temporary directory %s: %v", dir, err))
		}
		disk := leveldb.NewProducer(dir)

		// open raw databases
		db := disk.OpenDb("1")
		defer db.Drop()
		defer db.Close()

		dataTable := table.New(db, []byte("test_data"))
		idTable := table.New(db, []byte("id_last_migration"))
		idProducer := kvdb.NewIdProducer(idTable)

		list := Init("lachesis-test-db-fail", "654321")

		num := int64(1)
		lastGood := list.New(func() error {
			err := dataTable.Put([]byte("migration1"), []byte(strconv.FormatInt(num, 10)))
			if err != nil {
				return err
			}
			num++
			return nil
		}).New(func() error {
			err := dataTable.Put([]byte("migration2"), []byte(strconv.FormatInt(num, 10)))
			if err != nil {
				return err
			}
			num++
			return nil
		})

		afterBad := lastGood.New(func() error {
			return errors.New("test migration error")
		}).New(func() error {
			err := dataTable.Put([]byte("migration4"), []byte(strconv.FormatInt(num, 10)))
			if err != nil {
				return err
			}
			num++
			return nil
		})

		mgr := NewManager(afterBad, idProducer)
		err = mgr.Run()

		assert.Error(t, err, "Success run migration manager with error migrations")

		lastId := idProducer.GetId()
		assert.Equal(t, lastGood.Id(), lastId, "Bad last id in idProducer after migration error")

		testData1, _ := dataTable.Get([]byte("migration1"))
		testData2, _ := dataTable.Get([]byte("migration2"))
		testData3, _ := dataTable.Get([]byte("migration3"))
		testData4, _ := dataTable.Get([]byte("migration4"))

		assert.Equal(t, []byte("1"), testData1, "Bad value after run migration1")
		assert.Equal(t, []byte("2"), testData2, "Bad value after run migration2")

		assert.Empty(t, testData3, "Bad data for migration3 - should by empty")
		assert.Empty(t, testData4, "Bad data for migration4 - should by empty")

		/*
			Continue with fixed transactions
		*/
		num = 3
		fixed := lastGood.New(func() error {
			err := dataTable.Put([]byte("migration3"), []byte(strconv.FormatInt(num, 10)))
			if err != nil {
				return err
			}
			num++
			return nil
		}).New(func() error {
			err := dataTable.Put([]byte("migration4"), []byte(strconv.FormatInt(num, 10)))
			if err != nil {
				return err
			}
			num++
			return nil
		})

		mgr = NewManager(fixed, idProducer)
		err = mgr.Run()

		assert.NoError(t, err, "Error when run migration manager")

		testData1, _ = dataTable.Get([]byte("migration1"))
		testData2, _ = dataTable.Get([]byte("migration2"))
		testData3, _ = dataTable.Get([]byte("migration3"))
		testData4, _ = dataTable.Get([]byte("migration4"))

		assert.Equal(t, []byte("1"), testData1, "Bad value after run migration1")
		assert.Equal(t, []byte("2"), testData2, "Bad value after run migration2")
		assert.Equal(t, []byte("3"), testData3, "Bad value after run migration3")
		assert.Equal(t, []byte("4"), testData4, "Bad value after run migration4")
	})
}
