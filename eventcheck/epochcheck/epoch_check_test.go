package epochcheck

import (
	"github.com/Fantom-foundation/go-lachesis/eventcheck/testCommon"
	"github.com/Fantom-foundation/go-lachesis/inter"
	"github.com/Fantom-foundation/go-lachesis/inter/idx"
	"github.com/Fantom-foundation/go-lachesis/inter/pos"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestDagReader is a dagReader implementation made for testing purposes
type TestDagReader struct {
}

var testEpoch idx.Epoch = 1

// GetEpochValidators is an implementation used for testing
func (t TestDagReader) GetEpochValidators() (*pos.Validators, idx.Epoch) {
	vb := pos.NewBuilder()
	vb[1] = 1
	return vb.Build(), testEpoch
}

// makeTestEvents generates test events for a checker
func makeTestEvents() []inter.Event {
	return []inter.Event{
		{
			EventHeader: inter.EventHeader{
				EventHeaderData: inter.EventHeaderData{
					Version: 1,
					Seq:     0,
					Extra:   []byte{},
					Creator: 1,
					Epoch:   1,
				},
				Sig: []byte{},
			},
		},
		{
			EventHeader: inter.EventHeader{
				EventHeaderData: inter.EventHeaderData{
					Version: 1,
					Seq:     0,
					Extra:   []byte{},
					Creator: 2,
					Epoch:   1,
				},
				Sig: []byte{},
			},
		},
	}
}

// TestEpochCheck is a main testing func
func TestEpochCheck(t *testing.T) {
	lachesisConfigs := testCommon.MakeTestConfigs()
	for _, lachesisConfig := range lachesisConfigs {
		dagReader := TestDagReader{}
		checker := New(lachesisConfig, dagReader)
		for _, event := range makeTestEvents() {
			err := checker.Validate(&event)
			validators, epoch := checker.reader.GetEpochValidators()
			if event.Epoch != epoch {
				require.Equal(t, err, ErrNotRelevant)
				continue
			}
			if !validators.Exists(event.Creator) {
				require.Equal(t, err, ErrAuth)
				continue
			}
			require.Nil(t, err)
		}
		testEpoch = 0
	}
}
