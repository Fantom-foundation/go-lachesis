package election

import (
	"math/big"

	"github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/logger"
)

const (
	// TODO implement&test coinRound
	coinRound = 10 // every 10th round is a round with pseudorandom votes
)

type Election struct {
	// election params
	frameToDecide FrameHeight

	nodes         []ElectionNode
	totalStake    *big.Int // the sum of stakes (n)
	superMajority *big.Int // the quorum (should be 2/3n + 1)

	// election state
	decidedRoots map[RootSlot]voteValue
	votes        map[voteId]voteValue

	// external world
	stronglySee RootStronglySeeRootFn

	logger.Instance
}

type FrameHeight uint32

type ElectionNode struct {
	Nodeid      NodeId
	StakeAmount *big.Int
}

type RootHash struct {
	common.Hash
}

type NodeId struct {
	common.Hash
}

// specifies a slot {nodeid, frame}. Normal nodes can have only one root with this pair.
// Due to a fork, different roots may occupy the same slot
type RootSlot struct {
	Frame     FrameHeight
	Nodeid    NodeId
	nodeindex int // local nodeid, used for optimization
}

// @return hash of root B, if root A strongly sees root B.
// Due to a fork, there may be many roots B with the same slot,
// but strongly seen may be only one of them (if no more than 1/3n are Byzantine), with a specific hash.
type RootStronglySeeRootFn func(a RootHash, b RootSlot) *RootHash

type voteId struct {
	fromRoot  RootHash
	forNodeid NodeId
}
type voteValue struct {
	yes      bool
	seenRoot RootHash
}

type ElectionRes struct {
	DecidedFrame     FrameHeight
	DecidedSfWitness RootHash
}

func NewElection(
	nodes []ElectionNode,
	totalStake *big.Int,
	superMajority *big.Int,
	frameToDecide FrameHeight,
	stronglySeeFn RootStronglySeeRootFn,
) *Election {
	return &Election{
		nodes:         nodes,
		totalStake:    totalStake,
		superMajority: superMajority,
		frameToDecide: frameToDecide,
		decidedRoots:  make(map[RootSlot]voteValue),
		votes:         make(map[voteId]voteValue),
		stronglySee:   stronglySeeFn,
		Instance:      logger.MakeInstance(),
	}
}

// erase the current election state, prepare for new election frame
func (el *Election) ResetElection(frameToDecide FrameHeight) {
	el.frameToDecide = frameToDecide
	el.votes = make(map[voteId]voteValue)
	el.decidedRoots = make(map[RootSlot]voteValue)
}

// return root slots which are not within el.decidedRoots
func (el *Election) notDecidedRoots() []RootSlot {
	notDecidedRoots := make([]RootSlot, 0, len(el.nodes))

	for nodeindex, node := range el.nodes {
		slot := RootSlot{
			Frame:     el.frameToDecide,
			Nodeid:    node.Nodeid,
			nodeindex: nodeindex,
		}
		if _, ok := el.decidedRoots[slot]; !ok {
			notDecidedRoots = append(notDecidedRoots, slot)
		}
	}
	if len(notDecidedRoots)+len(el.decidedRoots) != len(el.nodes) { // sanity check
		el.Fatal("Mismatch of roots in notDecidedRoots()")
	}
	return notDecidedRoots
}

type weightedRoot struct {
	root        RootHash
	stakeAmount *big.Int
}

// @return all the roots which are strongly seen by the specified root at the specified frame
func (el *Election) stronglySeenRoots(root RootHash, frame FrameHeight) []weightedRoot {
	seenRoots := make([]weightedRoot, 0, len(el.nodes))
	for nodeindex, node := range el.nodes {
		slot := RootSlot{
			Frame:     frame,
			Nodeid:    node.Nodeid,
			nodeindex: nodeindex,
		}
		seenRoot := el.stronglySee(root, slot)
		if seenRoot != nil {
			seenRoots = append(seenRoots, weightedRoot{
				root:        *seenRoot,
				stakeAmount: el.nodes[nodeindex].StakeAmount,
			})
		}
	}
	return seenRoots
}

type GetRootsFn func(slot RootSlot) []RootHash

// The function is similar to ProcessRoot, but it fully re-processes the current voting.
// This routine should be called after node startup, and after each decided frame.
func (el *Election) ProcessKnownRoots(maxKnownFrame FrameHeight, getRootsFn GetRootsFn) (*ElectionRes, error) {
	// iterate all the roots from lowest frame to highest, call ProcessRootVotes for each
	for frame := el.frameToDecide + 1; frame <= maxKnownFrame; frame++ {
		for nodeindex, node := range el.nodes {
			slot := RootSlot{
				Frame:     frame,
				Nodeid:    node.Nodeid,
				nodeindex: nodeindex,
			}
			roots := getRootsFn(slot)
			// if there's more than 1 root, then all of them are forks. it's fine
			for _, root := range roots {
				decided, err := el.ProcessRoot(root, slot)
				if decided != nil || err != nil {
					return decided, err
				}
			}
		}
	}

	return nil, nil
}