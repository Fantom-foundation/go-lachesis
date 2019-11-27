package common

import (
	"github.com/Fantom-foundation/go-ethereum/common"
)

// PrettyDuration is a pretty printed version of a time.Duration value that cuts
// the unnecessary precision off from the formatted textual representation.
type PrettyDuration = common.PrettyDuration

// PrettyAge is a pretty printed version of a time.Duration value that rounds
// the values up to a single most significant unit, days/weeks/years included.
type PrettyAge = common.PrettyAge
