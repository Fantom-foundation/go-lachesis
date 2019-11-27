package types

import (
	"github.com/Fantom-foundation/go-ethereum/core/types"
)

type Log = types.Log

//type logMarshaling = types.logMarshaling
//
//type rlpLog= types.rlpLog
//
//// rlpStorageLog is the storage encoding of a log.
//type rlpStorageLog = types.rlpStorageLog
//
//// legacyRlpStorageLog is the previous storage encoding of a log including some redundant fields.
//type legacyRlpStorageLog = types.legacyRlpStorageLog


// LogForStorage is a wrapper around a Log that flattens and parses the entire content of
// a log including non-consensus fields.
type LogForStorage = types.LogForStorage
