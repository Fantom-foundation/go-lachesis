package migration

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrations(t *testing.T) {
	testData := map[string]int{}
	idProducer := &inmemIdProducer{}
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

	mgr := NewManager(afterBad, idProducer)

	err := mgr.Run()
	assert.Error(t, err, "Success run migration manager with error migrations")

	lastId := idProducer.GetId()
	assert.Equal(t, lastGood.Id(), lastId, "Bad last id in idProducer after migration error")

	assert.Equal(t, 1, testData["migration1"], "Bad value after run migration1")
	assert.Equal(t, 2, testData["migration2"], "Bad value after run migration2")
	assert.Equal(t, 3, testData["migration3"], "Bad value after run migration3")
	assert.Empty(t, testData["migration4"], "Bad data for migration4 - should by empty")

	// Continue with fixed transactions

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
}
