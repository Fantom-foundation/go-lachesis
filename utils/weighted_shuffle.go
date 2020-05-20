package utils

import (
	"crypto/sha256"

	"github.com/ethereum/go-ethereum/common"

	"github.com/Fantom-foundation/go-lachesis/common/littleendian"
	"github.com/Fantom-foundation/go-lachesis/inter/pos"
)

type weightedShuffleNode struct {
	thisWeight  pos.Stake
	leftWeight  pos.Stake
	rightWeight pos.Stake
}

type weightedShuffleTree struct {
	seed      common.Hash
	seedIndex int

	weights []pos.Stake
	nodes   []weightedShuffleNode
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

func (t *weightedShuffleTree) rand64() uint64 {
	if t.seedIndex == 32 {
		hasher := sha256.New() // use sha2 instead of sha3 for speed
		hasher.Write(t.seed.Bytes())
		t.seed = common.BytesToHash(hasher.Sum(nil))
		t.seedIndex = 0
	}
	// use not used parts of old seed, instead of calculating new one
	res := littleendian.BytesToInt64(t.seed[t.seedIndex : t.seedIndex+8])
	t.seedIndex += 8
	return res
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
