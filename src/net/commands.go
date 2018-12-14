package net

import "github.com/Fantom-foundation/go-lachesis/src/poset"

type SyncRequest struct {
	FromID uint32
	Known  map[uint32]int64
}

type SyncResponse struct {
	FromID    uint32
	SyncLimit bool
	Events    []poset.WireEvent
	Known     map[uint32]int64
}

//++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

type EagerSyncRequest struct {
	FromID uint32
	Events []poset.WireEvent
}

type EagerSyncResponse struct {
	FromID  uint32
	Success bool
}

//++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

type FastForwardRequest struct {
	FromID uint32
}

type FastForwardResponse struct {
	FromID   uint32
	Block    poset.Block
	Frame    poset.Frame
	Snapshot []byte
}
