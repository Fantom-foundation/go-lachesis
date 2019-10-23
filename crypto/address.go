package crypto

import (
	"crypto/ecdsa"

	"github.com/Fantom-foundation/go-ethereum/common"
	"github.com/Fantom-foundation/go-ethereum/crypto"
)

// PubkeyToAddress is a double of go-ethereum/crypto.PubkeyToAddress
// to don't import both packages.
func PubkeyToAddress(p ecdsa.PublicKey) common.Address {
	return crypto.PubkeyToAddress(p)
}
