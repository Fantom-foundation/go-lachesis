package proxy

import (
	"github.com/Fantom-foundation/go-lachesis/src/poset"
	"github.com/Fantom-foundation/go-lachesis/src/proxy/proto"
)

/*
These types are exported and need to be implemented and used by the calling
application.
*/

// ProxyHandler provides an interface for the application to set handlers for
// committing, retrieving and restoring state and transactions, to and from
// the DAG
type ProxyHandler interface {
	// CommitHandler is called when Lachesis commits a block to the DAG. It
	// returns the state hash resulting from applying the block's transactions to the
	// state
	CommitHandler(block poset.Block) (commitResponse proto.Response, err error)

	// SnapshotHandler is called by Lachesis to retrieve a snapshot corresponding to a
	// particular block
	SnapshotHandler(blockIndex int64) (snapshot []byte, err error)

	// RestoreHandler is called by Lachesis to restore the application to a specific
	// state
	RestoreHandler(snapshot []byte) (stateHash []byte, err error)
}
