package ethapi

import (
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/stretchr/testify/assert"
	"testing"
)

// PublicAccountAPI

func TestPublicAccountAPI_Accounts(t *testing.T) {
	b := NewTestBackend()
	am := accounts.NewManager(&accounts.Config{InsecureUnlockAllowed:true}, b)

	api := NewPublicAccountAPI(am)
	assert.NotPanics(t, func() {
		api.Accounts()
	})
}

// PrivateAccountAPI

func TestPrivateAccountAPI_DeriveAccount(t *testing.T) {
	b := NewTestBackend()

	nonceLock := new(AddrLocker)
	api := NewPrivateAccountAPI(b, nonceLock)
	assert.NotPanics(t, func() {
		api.DeriveAccount("https://test.ru", "/test", nil)
	})
}
func TestPrivateAccountAPI_ImportRawKey(t *testing.T) {
	b := NewTestBackend()

	nonceLock := new(AddrLocker)
	api := NewPrivateAccountAPI(b, nonceLock)

	// Prepare keystore
	keyStore := keystore.NewKeyStore("/tmp", 2, 2)
	for _, ac := range keyStore.Accounts() {
		keyStore.Delete(ac, "1234")
	}

	// Prepare account manager
	api.am = accounts.NewManager(&accounts.Config{InsecureUnlockAllowed:true}, b, keyStore)

	assert.NotPanics(t, func() {
		res, err := api.ImportRawKey("11223344556677889900aabbccddff0011223344556677889900aabbccddff00", "1234")
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
func TestPrivateAccountAPI_ListAccounts(t *testing.T) {
	b := NewTestBackend()

	nonceLock := new(AddrLocker)
	api := NewPrivateAccountAPI(b, nonceLock)

	// Prepare keystore
	keyStore := keystore.NewKeyStore("/tmp", 2, 2)

	// Prepare account manager
	api.am = accounts.NewManager(&accounts.Config{InsecureUnlockAllowed:true}, b, keyStore)

	assert.NotPanics(t, func() {
		api.ListAccounts()
	})
}
func TestPrivateAccountAPI_ListWallets(t *testing.T) {
	b := NewTestBackend()

	nonceLock := new(AddrLocker)
	api := NewPrivateAccountAPI(b, nonceLock)

	// Prepare keystore
	keyStore := keystore.NewKeyStore("/tmp", 2, 2)

	// Prepare account manager
	api.am = accounts.NewManager(&accounts.Config{InsecureUnlockAllowed:true}, b, keyStore)

	assert.NotPanics(t, func() {
		api.ListWallets()
	})
}
func TestPrivateAccountAPI_NewAccount(t *testing.T) {
	b := NewTestBackend()

	nonceLock := new(AddrLocker)
	api := NewPrivateAccountAPI(b, nonceLock)

	// Prepare keystore
	keyStore := keystore.NewKeyStore("/tmp", 2, 2)
	for _, ac := range keyStore.Accounts() {
		keyStore.Delete(ac, "1234")
	}

	// Prepare account manager
	api.am = accounts.NewManager(&accounts.Config{InsecureUnlockAllowed:true}, b, keyStore)

	assert.NotPanics(t, func() {
		res, err := api.NewAccount("1234")
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}
