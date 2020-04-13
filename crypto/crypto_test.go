package crypto

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"testing"
)

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
func newTestWallet(t *testing.T) TestWallet {
	privateKey, err := crypto.GenerateKey()
	require.Nil(t, err)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)

	//publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return TestWallet{
		PrivateKey: *privateKey,
		PubKey:     *publicKeyECDSA,
		Address:    address,
	}
}

// TestWallet is just a wallet for tests
type TestWallet struct {
	Address    common.Address
	PubKey     ecdsa.PublicKey
	PrivateKey ecdsa.PrivateKey
}
