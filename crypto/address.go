package crypto

import (
	"crypto/ecdsa"

	"github.com/Fantom-foundation/go-ethereum/crypto"
	
	"github.com/Fantom-foundation/go-lachesis/ethapi/common"
)

// PubkeyToAddress is a double of go-ethereum/crypto.PubkeyToAddress
// to don't import both packages.
func PubkeyToAddress(p ecdsa.PublicKey) common.Address {
	return crypto.PubkeyToAddress(p)
}
