package crypto

import (
	"crypto/ecdsa"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

// testWallet is just a wallet for tests
type testWallet struct {
	Address    common.Address
	PubKey     ecdsa.PublicKey
	PrivateKey ecdsa.PrivateKey
}

func TestFakeKey(t *testing.T) {
	k1 := FakeKey(1)
	k2 := FakeKey(2)
	require.NotEqual(t, k1, k2)
	require.Equal(t, k1, FakeKey(1))
}

func TestPubkeyToAddress(t *testing.T) {
	tw := newTestWallet(t)
	addr := PubkeyToAddress(tw.PubKey)
	require.Equal(t, addr, tw.Address)
}

// newTestWallet creates test wallet
func newTestWallet(t *testing.T) testWallet {
	privateKey, err := crypto.GenerateKey()
	require.Nil(t, err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return testWallet{
		PrivateKey: *privateKey,
		PubKey:     *publicKeyECDSA,
		Address:    address,
	}
}
