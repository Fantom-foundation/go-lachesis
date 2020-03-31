package ethapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// PublicTxPoolAPI

func TestPublicTxPoolAPI_Content(t *testing.T) {
	b := NewTestBackend()

	api := NewPublicTxPoolAPI(b)
	assert.NotPanics(t, func() {
		res := api.Content()
		assert.NotEmpty(t, res)
	})
}
func TestPublicTxPoolAPI_Status(t *testing.T) {
	b := NewTestBackend()

	api := NewPublicTxPoolAPI(b)
	assert.NotPanics(t, func() {
		res := api.Status()
		assert.NotEmpty(t, res)
	})
}
func TestPublicTxPoolAPI_Inspect(t *testing.T) {
	b := NewTestBackend()

	api := NewPublicTxPoolAPI(b)
	assert.NotPanics(t, func() {
		res := api.Inspect()
		assert.NotEmpty(t, res)
	})
}
