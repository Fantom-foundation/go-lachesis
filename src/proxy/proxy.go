package proxy

import (
	"github.com/andrecronje/lachesis/src/poset"
)

type AppProxy interface {
	SubmitCh() chan []byte
	CommitBlock(block poset.Block) ([]byte, error)
	GetSnapshot(blockIndex int) ([]byte, error)
	Restore(snapshot []byte) error
}
