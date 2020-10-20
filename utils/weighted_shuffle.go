package utils

import (
	"github.com/Fantom-foundation/go-lachesis/inter/pos"
)

type weightedShuffleNode struct {
	thisWeight  pos.Stake
	leftWeight  pos.Stake
	rightWeight pos.Stake
}

type weightedShuffleTree struct {
	weights []pos.Stake
	nodes   []weightedShuffleNode
	deterministicRand
}

func (t *weightedShuffleTree) leftIndex(i int) int {
	return i*2 + 1
}

func (t *weightedShuffleTree) rightIndex(i int) int {
	return i*2 + 2
}

func (t *weightedShuffleTree) build(i int) pos.Stake {
	if i >= len(t.weights) {
		return 0
	}
	thisW := t.weights[i]
	leftW := t.build(t.leftIndex(i))
	rightW := t.build(t.rightIndex(i))

	if thisW <= 0 {
		panic("all the weight must be positive")
	}

	t.nodes[i] = weightedShuffleNode{
		thisWeight:  thisW,
		leftWeight:  leftW,
		rightWeight: rightW,
	}
	return thisW + leftW + rightW
}

func (t *weightedShuffleTree) retrieve(i int) int {
	node := t.nodes[i]
	total := node.rightWeight + node.leftWeight + node.thisWeight

	r := pos.Stake(t.rand64()) % total

	if r < node.thisWeight {
		t.nodes[i].thisWeight = 0
		return i
	} else if r < node.thisWeight+node.leftWeight {
		chosen := t.retrieve(t.leftIndex(i))
		t.nodes[i].leftWeight -= t.weights[chosen]
		return chosen
	} else {
		chosen := t.retrieve(t.rightIndex(i))
		t.nodes[i].rightWeight -= t.weights[chosen]
		return chosen
	}
}
