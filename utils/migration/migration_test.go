package migration

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMigrations(t *testing.T) {
	require := require.New(t)

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
	require.Error(err, "Success run migration manager with error migrations")

	lastId := idProducer.GetId()
	require.Equal(lastGood.Id(), lastId, "Bad last id in idProducer after migration error")

	require.Equal(1, testData["migration1"], "Bad value after run migration1")
	require.Equal(2, testData["migration2"], "Bad value after run migration2")
	require.Equal(3, testData["migration3"], "Bad value after run migration3")
	require.Empty(testData["migration4"], "Bad data for migration4 - should by empty")

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
	require.NoError(err, "Error when run migration manager")

	require.Equal(1, testData["migration1"], "Bad value after run migration1")
	require.Equal(2, testData["migration2"], "Bad value after run migration2")
	require.Equal(3, testData["migration3"], "Bad value after run migration3")
	require.Equal(4, testData["migration4"], "Bad value after run migration4")
}
