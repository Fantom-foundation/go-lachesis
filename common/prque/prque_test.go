package prque

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// emptySetIndexCallback is just an empty function made for testing purposes
var emptySetIndexCallback = func(a interface{}, i int) {

}

// Test_Prque runs a set of tests for a prque package
func Test_Prque(t *testing.T) {
	p := New(emptySetIndexCallback)

	testPush(t, p)
	testPop(t, p)
	testRemove(t, p)
	testResetCont(t, p)
}

func testResetCont(t *testing.T, p *Prque) {
	// we assume that tests will run with an empty Prque
	require.True(t, p.Empty())

	p.Push(makeDataWithSeize(0), 0)
	require.False(t, p.Empty())

	p.cont.Reset()
	require.True(t, p.Empty())
}

// testPop tests that pop operations work correct
func testPop(t *testing.T, p *Prque) {
	// we assume that tests will run with an empty Prque
	require.True(t, p.Empty())
	popPanicFunc := func() { p.Pop() }
	popItemPanicFunc := func() { p.PopItem() }

	require.Panics(t, popPanicFunc)
	require.Panics(t, popItemPanicFunc)

	p.Push(makeDataWithSeize(0), 0)
	p.Push(makeDataWithSeize(0), 0)
	require.True(t, p.Size() == 2)

	p.Pop()
	require.True(t, p.Size() == 1)

	p.PopItem()
	require.True(t, p.Empty())
}

func testRemove(t *testing.T, p *Prque) {
	// we assume that tests will run with an empty Prque
	require.True(t, p.Empty())

	popRemoveFunc := func() { p.Remove(0) }
	require.Panics(t, popRemoveFunc)

	itemValueSeize1 := 2
	itemValueSeize2 := 3
	item1 := makeDataWithSeize(itemValueSeize1)
	item2 := makeDataWithSeize(itemValueSeize2)
	p.Push(item1, 0)
	p.Push(item2, 0)
	require.True(t, p.Size() == 2)

	removedItem := p.Remove(-1)
	require.Nil(t, removedItem)

	removedItem2 := p.Remove(1)
	removedData2 := removedItem2.(*item)
	require.True(t, len(removedData2.value.([]byte)) == itemValueSeize2)
	require.True(t, p.Size() == 1)
	removedItem1 := p.Remove(0)
	removedData1 := removedItem1.(*item)
	require.True(t, len(removedData1.value.([]byte)) == itemValueSeize1)
	require.True(t, p.Empty())
}

// testPush tests that push operations work correct
func testPush(t *testing.T, p *Prque) {
	// we assume that tests will run with an empty Prque
	require.True(t, p.Empty())

	testCommonPushLogic(t, p)
	testCapacityIncrease(t, p)
	p.Reset()
}

// testCapacityIncrease tests that a capacity of a stack expands properly
func testCapacityIncrease(t *testing.T, p *Prque) {
	var insertions []interface{}
	for i := 0; i < blockSize+1; i++ {
		insertions = append(insertions, makeDataWithSeize(0))
	}

	var order int64 = 0
	for _, v := range insertions {
		p.Push(v, order)
		order++
	}

	require.True(t, p.Size() > blockSize)
	require.True(t, p.cont.capacity == blockSize*2)
	require.Equal(t, p.cont.offset, p.cont.size-blockSize)
}

// testCommonPushLogic tests that an execution of a push changes stack state properly
func testCommonPushLogic(t *testing.T, p *Prque) {
	dataCases := []interface{}{
		nil,
		makeDataWithSeize(0),
		makeDataWithSeize(1),
		makeDataWithSeize(1000),
		makeDataWithSeize(1e8),
	}

	// the reason of chosen values is the following:
	// we want to check a result of a negative, positive and neutral value input
	// we also want to ensure that enormously big integers would be handled well (since the input type is int64)
	// and that arbitrary normal integers would be handled correct too
	orderCases := []int64{
		-1e18,
		-1,
		0,
		1,
		1e18,
	}

	var pushed int
	for _, dataCase := range dataCases {
		for _, orderCase := range orderCases {
			p.Push(dataCase, orderCase)
			require.NotNil(t, p.cont.active[pushed])
			pushed++
			require.Equal(t, p.Size(), pushed)
		}
	}
}

// makeDataWithSeize returns an array of bytes coerced to an interface type
func makeDataWithSeize(sieze int) interface{} {
	return make([]byte, sieze)
}
