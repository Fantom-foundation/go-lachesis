package election

import (
	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/common"

	"github.com/Fantom-foundation/go-lachesis/src/hash"
)

// may be used in tests to match election state
func (el *Election) DebugStateHash() common.Hash {
	hasher := sha3.New256()
	write := func(bb []byte) {
		if _, err := hasher.Write(bb); err != nil {
			panic(err)
		}
	}

	for vid, vote := range el.votes {
		write(vid.fromRoot.Bytes())
		write(vote.seenRoot.Bytes())
	}
	for member, vote := range el.decidedRoots {
		write(member.Bytes())
		write(vote.seenRoot.Bytes())
	}
	return hash.FromBytes(hasher.Sum(nil))
}

// @param (optional) voters is roots to print votes for. May be nil
// @return election summary in a human readable format
func (el *Election) String(voters []hash.Event) string {
	if voters == nil {
		votersM := make(map[hash.Event]bool)
		for vid := range el.votes {
			votersM[vid.fromRoot] = true
		}
		for voter := range votersM {
			voters = append(voters, voter)
		}
	}

	info := "Every line contains votes from a root, for each subject. y is yes, n is no. Upper case means 'decided'. '-' means that subject was already decided when root was processed.\n"
	for _, root := range voters { // voter
		info += root.String() + ": "
		for forM := range el.members { // subject
			vid := voteId{
				fromRoot:  root,
				forMember: forM,
			}
			vote, ok := el.votes[vid]
			if !ok { // i.e. subject was decided when root processed
				info += "-"
				continue
			}
			if vote.yes {
				if vote.decided {
					info += "Y"
				} else {
					info += "y"
				}
			} else {
				if vote.decided {
					info += "N"
				} else {
					info += "n"
				}
			}
		}
		info += "\n"
	}
	return info
}
