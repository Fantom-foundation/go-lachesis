package main

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/sfctype"
)

func TestSize(t *testing.T) {

	t.Run("cache.Delegators", func(t *testing.T) {
		var key common.Address
		var val sfctype.SfcDelegator
		val.Amount = big.NewInt(10000000000)
		size := memSizeOf(key) + memSizeOf(val)
		t.Logf("Item size: ~ %d bytes", size)
		t.Logf("Cache size: ~ %d kbytes", size*2*4000/1024)
	})

	t.Run("cache.Stakers", func(t *testing.T) {
		var key idx.StakerID
		var val sfctype.SfcStaker
		val.StakeAmount = big.NewInt(10000000000)
		val.DelegatedMe = big.NewInt(10000000000)
		size := memSizeOf(key) + memSizeOf(val)
		t.Logf("Item size: ~ %d bytes", size)
		t.Logf("Cache size: ~ %d kbytes", size*2*4000/1024)
	})

	t.Run("cache.Receipts", func(t *testing.T) {
		var key idx.Block
		var val types.Receipts
		size := memSizeOf(key) + memSizeOf(val)
		t.Logf("Item size: ~ %d bytes", size)
		t.Logf("Cache size: ~ %d kbytes", size*2*100/1024)
	})
}
