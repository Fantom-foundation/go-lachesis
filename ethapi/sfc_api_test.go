package ethapi

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestPublicSfcAPI_GetDelegator(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicSfcAPI(b)
	require.NotPanics(t, func() {
		res, err := api.GetDelegator(ctx, common.Address{1}, 2)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
	require.NotPanics(t, func() {
		res, err := api.GetDelegator(ctx, common.Address{1}, 0)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
}
func TestPublicSfcAPI_GetDelegatorClaimedRewards(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicSfcAPI(b)
	require.NotPanics(t, func() {
		res, err := api.GetDelegatorClaimedRewards(ctx, common.Address{1})
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
}
func TestPublicSfcAPI_GetDelegatorsOf(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicSfcAPI(b)
	require.NotPanics(t, func() {
		res, err := api.GetDelegatorsOf(ctx, 1, 2)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
	require.NotPanics(t, func() {
		res, err := api.GetDelegatorsOf(ctx, 1, 0)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
}
func TestPublicSfcAPI_GetDowntime(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicSfcAPI(b)
	require.NotPanics(t, func() {
		res, err := api.GetDowntime(ctx, 2)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
}
func TestPublicSfcAPI_GetOriginationScore(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicSfcAPI(b)
	require.NotPanics(t, func() {
		res, err := api.GetOriginationScore(ctx, 1)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
}
func TestPublicSfcAPI_GetRewardWeights(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicSfcAPI(b)
	require.NotPanics(t, func() {
		res, err := api.GetRewardWeights(ctx, 1)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
}
func TestPublicSfcAPI_GetStaker(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicSfcAPI(b)
	require.NotPanics(t, func() {
		res, err := api.GetStaker(ctx, 1, 4)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
}
func TestPublicSfcAPI_GetStakers(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicSfcAPI(b)
	require.NotPanics(t, func() {
		res, err := api.GetStakers(ctx, 4)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
	require.NotPanics(t, func() {
		res, err := api.GetStakers(ctx, 0)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
}
func TestPublicSfcAPI_GetStakerByAddress(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicSfcAPI(b)
	require.NotPanics(t, func() {
		res, err := api.GetStakerByAddress(ctx, common.Address{1}, 4)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
	require.NotPanics(t, func() {
		res, err := api.GetStakerByAddress(ctx, common.Address{1}, 0)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
}
func TestPublicSfcAPI_GetStakerClaimedRewards(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicSfcAPI(b)

	require.NotPanics(t, func() {
		res, err := api.GetStakerClaimedRewards(ctx, 1)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
}
func TestPublicSfcAPI_GetStakerDelegatorsClaimedRewards(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicSfcAPI(b)

	require.NotPanics(t, func() {
		res, err := api.GetStakerDelegatorsClaimedRewards(ctx, 1)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
}
func TestPublicSfcAPI_GetStakerPoI(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicSfcAPI(b)

	require.NotPanics(t, func() {
		res, err := api.GetStakerPoI(ctx, 1)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
}
func TestPublicSfcAPI_GetValidationScore(t *testing.T) {
	ctx := context.TODO()
	b := NewTestBackend()

	api := NewPublicSfcAPI(b)

	require.NotPanics(t, func() {
		res, err := api.GetValidationScore(ctx, 1)
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})
}
