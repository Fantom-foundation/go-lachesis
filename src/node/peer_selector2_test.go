package node

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Fantom-foundation/go-lachesis/src/peers"
)

func TestSmartSelectorEmpty(t *testing.T) {
	assertO := assert.New(t)

	fp := fakePeers(0)

	ss := NewSmartPeerSelector(
		fp,
		SmartPeerSelectorCreationFnArgs{
			LocalAddr: "",
			GetFlagTable: func() (map[string]int64, error) {
				return nil, nil
			},
		},
	)

	assertO.Nil(ss.Next())
}

func TestSmartSelectorLocalAddrOnly(t *testing.T) {
	assertO := assert.New(t)

	fp := fakePeers(1)
	fps := fp.ToPeerSlice()

	ss := NewSmartPeerSelector(
		fp,
		SmartPeerSelectorCreationFnArgs{
			LocalAddr: fps[0].Message.NetAddr,
			GetFlagTable: func() (map[string]int64, error) {
				return nil, nil
			},
		},
	)

	assertO.Nil(ss.Next())
}

func TestSmartSelectorUsed(t *testing.T) {
	assertO := assert.New(t)

	fp := fakePeers(3)
	fps := fp.ToPeerSlice()

	ss := NewSmartPeerSelector(
		fp,
		SmartPeerSelectorCreationFnArgs{
			LocalAddr: fps[0].Message.NetAddr,
			GetFlagTable: func() (map[string]int64, error) {
				return nil, nil
			},
		},
	)

	choose1 := ss.Next().Message.NetAddr
	assertO.NotEqual(fps[0].Message.NetAddr, choose1)

	choose2 := ss.Next().Message.NetAddr
	assertO.NotEqual(fps[0].Message.NetAddr, choose2)
	assertO.NotEqual(choose1, choose2)

	choose3 := ss.Next().Message.NetAddr
	assertO.NotEqual(fps[0].Message.NetAddr, choose3)
}

func TestSmartSelectorFlagged(t *testing.T) {
	assertO := assert.New(t)

	fp := fakePeers(3)
	fps := fp.ToPeerSlice()

	ss := NewSmartPeerSelector(
		fp,
		SmartPeerSelectorCreationFnArgs{
			LocalAddr: fps[0].Message.NetAddr,
			GetFlagTable: func() (map[string]int64, error) {
				return map[string]int64{
					fps[2].Message.PubKeyHex: 1,
				}, nil
			},
		},
	)

	assertO.Equal(fps[1].Message.NetAddr, ss.Next().Message.NetAddr)
	assertO.Equal(fps[1].Message.NetAddr, ss.Next().Message.NetAddr)
	assertO.Equal(fps[1].Message.NetAddr, ss.Next().Message.NetAddr)
}

func TestSmartSelectorGeneral(t *testing.T) {
	assertO := assert.New(t)

	fp := fakePeers(4)
	fps := fp.ToPeerSlice()

	ss := NewSmartPeerSelector(
		fp,
		SmartPeerSelectorCreationFnArgs{
			LocalAddr: fps[3].Message.NetAddr,
			GetFlagTable: func() (map[string]int64, error) {
				return map[string]int64{
					fps[0].Message.PubKeyHex: 0,
					fps[1].Message.PubKeyHex: 0,
					fps[2].Message.PubKeyHex: 1,
					fps[3].Message.PubKeyHex: 0,
				}, nil
			},
		},
	)

	addresses := []string{fps[0].Message.NetAddr, fps[1].Message.NetAddr}
	assertO.Contains(addresses, ss.Next().Message.NetAddr)
	assertO.Contains(addresses, ss.Next().Message.NetAddr)
	assertO.Contains(addresses, ss.Next().Message.NetAddr)
	assertO.Contains(addresses, ss.Next().Message.NetAddr)
}

/*
 * go test -bench "BenchmarkSmartSelectorNext" -benchmem -run "^$" ./src/node
 */

func BenchmarkSmartSelectorNext(b *testing.B) {
	const fakePeersCount = 50

	participants1 := fakePeers(fakePeersCount)
	participants2 := clonePeers(participants1)

	flagTable1 := fakeFlagTable(participants1)

	ss1 := NewSmartPeerSelector(
		participants1,
		SmartPeerSelectorCreationFnArgs{
			LocalAddr: fakeAddr(0),
			GetFlagTable: func() (map[string]int64, error) {
				return flagTable1, nil
			},
		},
	)
	rnd := NewRandomPeerSelector(
		participants2,
		RandomPeerSelectorCreationFnArgs{
			LocalAddr: fakeAddr(0),
		},
	)

	b.ResetTimer()

	b.Run("smart Next()", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			p := ss1.Next()
			if p == nil {
				b.Fatal("No next peer")
			}
			ss1.UpdateLast(p.Message.PubKeyHex)
		}
	})

	b.Run("simple Next()", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			p := rnd.Next()
			if p == nil {
				b.Fatal("No next peer")
			}
			rnd.UpdateLast(p.Message.PubKeyHex)
		}
	})

}

/*
 * utility function for peer_selector2_test
 */

func fakeFlagTable(participants *peers.Peers) map[string]int64 {
	res := make(map[string]int64, participants.Len())
	for _, p := range participants.ToPeerSlice() {
		res[p.Message.PubKeyHex] = rand.Int63n(2)
	}
	return res
}
