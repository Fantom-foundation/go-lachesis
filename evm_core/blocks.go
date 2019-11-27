package evm_core

import (
	"github.com/Fantom-foundation/go-lachesis/ethapi/common"
)

// BadHashes represent a set of manually tracked bad hashes (usually hard forks)
var BadHashes = map[common.Hash]bool{}
