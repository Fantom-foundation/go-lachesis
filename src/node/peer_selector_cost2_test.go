package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnfairSelectorEmpty(t *testing.T) {
	assert0 := assert.New(t)

	fp := fakePeers(0)

	fs := NewUnfairPeerSelector(
		fp,
		UnfairPeerSelectorCreationFnArgs{
			LocalAddr: "",
		},
	)

	assert0.Nil(fs.Next())
}

func TestUnfairSelectorLocalAddrOnly(t *testing.T) {
	assertO := assert.New(t)

	fp := fakePeers(1)
	fps := fp.ToPeerSlice()

	fs := NewUnfairPeerSelector(
		fp,
		UnfairPeerSelectorCreationFnArgs{
			LocalAddr: fps[0].NetAddr,
		},
	)

	assertO.Nil(fs.Next())
}

func TestUnfairSelectorGeneral(t *testing.T) {
	assertO := assert.New(t)

	fp := fakePeers(4)
	fps := fp.ToPeerSlice()

	ss := NewUnfairPeerSelector(
		fp,
		UnfairPeerSelectorCreationFnArgs{
			LocalAddr: fps[3].NetAddr,
		},
	)

	addresses := []string{
		fps[0].NetAddr,
		fps[1].NetAddr,
		fps[2].NetAddr,
		fps[3].NetAddr,
	}
	assertO.Contains(addresses, ss.Next().NetAddr)
	assertO.Contains(addresses, ss.Next().NetAddr)
	assertO.Contains(addresses, ss.Next().NetAddr)
	assertO.Contains(addresses, ss.Next().NetAddr)
}

/*
 * go test -bench "BenchmarkUnfairSelectorNext" -benchmem -run "^$" ./src/node
 */

func BenchmarkUnfairSelectorNext(b *testing.B) {
	const fakePeersCount = 50

	participants1 := fakePeers(fakePeersCount)
	participants2 := clonePeers(participants1)

	fs1 := NewUnfairPeerSelector(
		participants1,
		UnfairPeerSelectorCreationFnArgs{
			LocalAddr: fakeAddr(0),
		},
	)
	rnd := NewRandomPeerSelector(
		participants2,
		RandomPeerSelectorCreationFnArgs{
			LocalAddr: fakeAddr(0),
		},
	)

	b.ResetTimer()

	b.Run("fair Next()", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			p := fs1.Next()
			if p == nil {
				b.Fatal("No next peer")
				break
			}
			fs1.UpdateLast(p.PubKeyHex)
		}
	})

	b.Run("simple Next()", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			p := rnd.Next()
			if p == nil {
				b.Fatal("No next peer")
				break
			}
			rnd.UpdateLast(p.PubKeyHex)
		}
	})

}
