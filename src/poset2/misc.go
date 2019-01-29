package poset2

import "github.com/Fantom-foundation/go-lachesis/src/common"

func withDuplicates(items []common.Hash) bool {
	table := make(map[common.Hash]bool)
	for v := range items {
		if table[items[v]] {
			return true
		} else {
			table[items[v]] = true
		}
	}
	return false
}
