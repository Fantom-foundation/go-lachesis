package app

import (
	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
)

// TxPosition includes block and event.
type TxPosition struct {
	Block       idx.Block
	Event       hash.Event
	EventOffset uint32
	BlockOffset uint32
}

// BlocksMissed is information about missed blocks from a staker.
type BlocksMissed struct {
	Num    idx.Block
	Period inter.Timestamp
}
