package evm_core

import (
	"github.com/quan8/go-ethereum/common"
)

// BadHashes represent a set of manually tracked bad hashes (usually hard forks)
var BadHashes = map[common.Hash]bool{}
