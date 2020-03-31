package prque

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// emptySetIndexCallback is just an empty function made for testing purposes
var emptySetIndexCallback = func(a interface{}, i int) {

}

// Test_Prque runs a set of tests for a prque package
func Test_Prque(t *testing.T) {
	p := New(emptySetIndexCallback)

	t.Run("push", func(t *testing.T) {
		p.Reset()
		testPush(t, p)
	})
	t.Run("pop", func(t *testing.T) {
		p.Reset()
		testPop(t, p)
	})
	t.Run("remove", func(t *testing.T) {
		p.Reset()
		testRemove(t, p)
	})
	t.Run("reset", func(t *testing.T) {
		p.Reset()
		testResetCont(t, p)
	})
}

func testResetCont(t *testing.T, p *Prque) {
	require := require.New(t)
	// we assume that tests will run with an empty Prque
	require.True(p.Empty())

	p.Push(makeData(0), 0)
	require.False(p.Empty())

	p.cont.Reset()
	require.True(p.Empty())
}

// testPop tests that pop operations work correct
func testPop(t *testing.T, p *Prque) {
	require := require.New(t)
	// we assume that tests will run with an empty Prque
	require.True(p.Empty())

	require.Panics(func() {
		_, _ = p.Pop()
	})
	require.Panics(func() {
		_ = p.PopItem()
	})

	p.Push(makeData(0), 0)
	p.Push(makeData(0), 0)
	require.Equal(2, p.Size())

	p.Pop()
	require.Equal(1, p.Size())

	p.PopItem()
	require.True(p.Empty())
}

func testRemove(t *testing.T, p *Prque) {
	require := require.New(t)
	// we assume that tests will run with an empty Prque
	require.True(p.Empty())

	require.Panics(func() {
		_ = p.Remove(0)
	})

	size1 := 2
	size2 := 3
	item1 := makeData(size1)
	item2 := makeData(size2)
	priority1 := int64(2)
	priority2 := int64(1)
	p.Push(item1, priority1)
	p.Push(item2, priority2)
	require.Equal(2, p.Size())

	require.Nil(
		p.Remove(-1))

	removed2 := p.Remove(1).(*item)
	require.Equal(size2, len(removed2.value.([]byte)))
	require.Equal(priority2, removed2.priority)
	require.Equal(1, p.Size())

	removed1 := p.Remove(0).(*item)
	require.Equal(size1, len(removed1.value.([]byte)))
	require.Equal(priority1, removed1.priority)
	require.True(p.Empty())
}

// testPush tests that push operations work correct
func testPush(t *testing.T, p *Prque) {
	// we assume that tests will run with an empty Prque
	require.True(t, p.Empty())

	testCommonPushLogic(t, p)
	testCapacityIncrease(t, p)
}

// testCapacityIncrease tests that a capacity of a stack expands properly
func testCapacityIncrease(t *testing.T, p *Prque) {
	require := require.New(t)

	var insertions []interface{}
	for i := 0; i < blockSize+1; i++ {
		insertions = append(insertions, makeData(0))
	}

	var order int64 = 0
	for _, v := range insertions {
		p.Push(v, order)
		order++
	}

	require.Greater(p.Size(), blockSize)
	require.Equal(blockSize*2, p.cont.capacity)
	require.Equal(p.cont.offset, p.cont.size-blockSize)
}

// testCommonPushLogic tests that an execution of a push changes stack state properly
func testCommonPushLogic(t *testing.T, p *Prque) {
	require := require.New(t)

	dataCases := []interface{}{
		nil,
		makeData(0),
		makeData(1),
		makeData(1000),
		makeData(1e8),
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
	for _, data := range dataCases {
		for _, order := range orderCases {
			p.Push(data, order)
			require.NotNil(p.cont.active[pushed])
			pushed++
			require.Equal(pushed, p.Size())
		}
	}
}

// makeData returns an array of bytes coerced to an interface type
func makeData(size int) interface{} {
	return make([]byte, size)
}
