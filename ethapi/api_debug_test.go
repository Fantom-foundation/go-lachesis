package ethapi

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

// PublicDebugAPI

func TestPublicDebugAPI_GetBlockRlp(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicDebugAPI(b)
	require.NotPanics(t, func() {
		res, err := api.GetBlockRlp(ctx, 1)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
}
func TestPublicDebugAPI_PrintBlock(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicDebugAPI(b)
	require.NotPanics(t, func() {
		res, err := api.PrintBlock(ctx, 1)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
}
func TestPublicDebugAPI_SeedHash(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicDebugAPI(b)
	require.NotPanics(t, func() {
		res, err := api.SeedHash(ctx, 1)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
}
func TestPublicDebugAPI_TestSignCliqueBlock(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicDebugAPI(b)
	require.NotPanics(t, func() {
		res, err := api.TestSignCliqueBlock(ctx, common.Address{}, 1)
		require.Error(t, err)
		require.NotEmpty(t, res)
	})
}

func TestPublicDebugAPI_TtfReport(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicDebugAPI(b)
	require.NotPanics(t, func() {
		res, err := api.TtfReport(ctx, 1, 1, "claimed_time", 3)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
}
func TestPublicDebugAPI_ValidatorTimeDrifts(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicDebugAPI(b)
	require.NotPanics(t, func() {
		res, err := api.ValidatorTimeDrifts(ctx, 1, 1, 3)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
}
func TestPublicDebugAPI_ValidatorVersions(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicDebugAPI(b)
	require.NotPanics(t, func() {
		_, err := api.ValidatorVersions(ctx, 1, 1)
		require.NoError(t, err)
	})
}

// PrivateDebugAPI

func TestPrivateDebugAPI_ChaindbProperty(t *testing.T) {
	b := NewTestBackend()

	api := NewPrivateDebugAPI(b)
	require.NotPanics(t, func() {
		_, _ = api.ChaindbProperty("p1")
	})
}

func TestPrivateDebugAPI_ChaindbCompact(t *testing.T) {
	b := NewTestBackend()

	api := NewPrivateDebugAPI(b)
	require.NotPanics(t, func() {
		err := api.ChaindbCompact()
		require.NoError(t, err)
	})
}

func TestPrivateDebugAPI_SetHead(t *testing.T) {
	b := NewTestBackend()

	api := NewPrivateDebugAPI(b)
	require.NotPanics(t, func() {
		err := api.SetHead(1)
		require.Error(t, err)
	})
}
