package inter

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"math"
	"math/rand"
	"testing"

	"github.com/Fantom-foundation/go-lachesis/hash"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/utils/fast"
	"github.com/ethereum/go-ethereum/rlp"
)

func TestEventHeaderDataSerialization(t *testing.T) {
	ee := map[string]EventHeaderData{
		"empty": EventHeaderData{
			Parents: hash.Events{},
			TxHash:  EmptyTxHash,
			Extra:   []uint8{},
		},
		"max": EventHeaderData{
			Epoch:        idx.Epoch(math.MaxUint32),
			GasPowerLeft: GasPowerLeft{Gas: [2]uint64{math.MaxUint64, math.MaxUint64}},
			Parents: hash.Events{
				hash.BytesToEvent(bytes.Repeat([]byte{math.MaxUint8}, 32)),
			},
			Extra: []uint8{},
		},
		"random": FakeEvent().EventHeaderData,
	}

	t.Run("ok", func(t *testing.T) {
		require := require.New(t)

		for name, header0 := range ee {
			buf, err := rlp.EncodeToBytes(&header0)
			require.NoError(err)

			var header1 EventHeaderData
			err = rlp.DecodeBytes(buf, &header1)
			require.NoError(err, name)

			require.EqualValues(header0, header1, name)
		}
	})

	t.Run("err", func(t *testing.T) {
		require := require.New(t)

		for name, header0 := range ee {
			bin, err := header0.MarshalBinary()
			require.NoError(err, name)

			n := rand.Intn(len(bin) - len(header0.Extra))
			bin = bin[0:n]

			buf, err := rlp.EncodeToBytes(bin)
			require.NoError(err, name)

			var header1 EventHeaderData
			err = rlp.DecodeBytes(buf, &header1)
			require.Error(err, name)
			//t.Log(err)
		}
	})
}

func BenchmarkEventHeaderData_EncodeRLP(b *testing.B) {
	header := FakeEvent().EventHeaderData

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf, err := rlp.EncodeToBytes(&header)
		require.NoError(b, err)
		b.ReportMetric(float64(len(buf)), "size")
	}
}

func BenchmarkEventHeaderData_DecodeRLP(b *testing.B) {
	header := FakeEvent().EventHeaderData

	buf, err := rlp.EncodeToBytes(&header)
	require.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err = rlp.DecodeBytes(buf, &header)
		require.NoError(b, err)
	}
}

func TestReadUintCompact(t *testing.T) {
	require := require.New(t)

	// canonical
	for exp, bb := range map[uint64][]byte{
		0x000000: []byte{0x00},
		0x0000FF: []byte{0xFF},
		0x010000: []byte{0x00, 0x00, 0x01},
	} {
		got, err := readUintCompact(fast.NewBuffer(bb), len(bb))
		require.NoError(err)
		require.Equal(exp, got, bb)
	}

	// non canonical
	for _, bb := range [][]byte{
		[]byte{0x00, 0x00},
		[]byte{0xFF, 0x00},
		[]byte{0x00, 0x00, 0x01, 0x00},
	} {
		_, err := readUintCompact(fast.NewBuffer(bb), len(bb))
		require.Error(err)

		require.Equal(ErrNonCanonicalEncoding, err, bb)
	}
}

// FakeEvent generates random event for testing purpose.
func FakeEvent() *Event {
	var epoch idx.Epoch = hash.FakeEpoch()

	e := NewEvent()
	e.Epoch = epoch
	e.Seq = idx.Event(9)
	e.GasPowerLeft.Gas[0] = 0x010000
	e.GasPowerLeft.Gas[1] = 100
	e.Creator = hash.FakePeer()
	e.Parents = hash.FakeEvents(8)
	e.PrevEpochHash = hash.FakeHash()
	e.Extra = make([]byte, 10, 10)
	e.Sig = []byte{}

	_, err := rand.Read(e.Extra)
	if err != nil {
		panic(err)
	}

	return e
}
