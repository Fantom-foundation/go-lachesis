package migration

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMigrations(t *testing.T) {
	require := require.New(t)

	testData := map[string]int{}
	current := &inmemIdProducer{}
	list := Begin("TestMigrations")

	num := 1
	lastGood := list.Next("01",
		func() error {
			testData["migration1"] = num
			num++
			return nil
		},
	).Next("02",
		func() error {
			testData["migration2"] = num
			num++
			return nil
		},
	)

	afterBad := lastGood.Next("03",
		func() error {
			testData["migration3"] = num
			num++
			return errors.New("test migration error")
		},
	).Next("04",
		func() error {
			testData["migration4"] = num
			num++
			return nil
		},
	)

	err := afterBad.Exec(current)
	require.Error(err, "Success run migration manager with error migrations")

	lastId := current.GetId()
	require.Equal(lastGood.name, lastId, "Bad last id in idProducer after migration error")

	require.Equal(1, testData["migration1"], "Bad value after run migration1")
	require.Equal(2, testData["migration2"], "Bad value after run migration2")
	require.Equal(3, testData["migration3"], "Bad value after run migration3")
	require.Empty(testData["migration4"], "Bad data for migration4 - should by empty")

	// Continue with fixed transactions

	num = 3
	fixed := lastGood.Next("03",
		func() error {
			testData["migration3"] = num
			num++
			return nil
		},
	).Next("04",
		func() error {
			testData["migration4"] = num
			num++
			return nil
		},
	)

	err = fixed.Exec(current)
	require.NoError(err, "Error when run migration manager")

	require.Equal(1, testData["migration1"], "Bad value after run migration1")
	require.Equal(2, testData["migration2"], "Bad value after run migration2")
	require.Equal(3, testData["migration3"], "Bad value after run migration3")
	require.Equal(4, testData["migration4"], "Bad value after run migration4")

	require.Equal("03", fixed.PrevByName("04").Name(), "Bad search for PrevByName")
	require.Equal("02", fixed.PrevByName("03").Name(), "Bad search for PrevByName")
	require.Equal("01", fixed.PrevByName("02").Name(), "Bad search for PrevByName")
	require.Equal("TestMigrations", fixed.PrevByName("01").Name(), "Bad search for PrevByName")
	require.Nil(fixed.PrevByName("NotExisted"), "Bad search for PrevByName")
}
