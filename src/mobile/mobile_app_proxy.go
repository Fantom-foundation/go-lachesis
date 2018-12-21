package mobile

import (
	"github.com/Fantom-foundation/go-lachesis/src/poset"
	"github.com/Fantom-foundation/go-lachesis/src/proxy"
	"github.com/sirupsen/logrus"
)

/*
This type is not exported
*/

// mobileAppProxy object
type mobileAppProxy struct {
	*proxy.InmemAppProxy

	commitHandler    CommitHandler
	exceptionHandler ExceptionHandler
	logger           *logrus.Logger
}

// newMobileAppProxy create proxy
func newMobileAppProxy(
	commitHandler CommitHandler,
	exceptionHandler ExceptionHandler,
	logger *logrus.Logger,
) *mobileAppProxy {

	mobileApp := &mobileAppProxy{
		commitHandler:    commitHandler,
		exceptionHandler: exceptionHandler,
		logger:           logger,
	}

	mobileApp.InmemAppProxy = proxy.NewInmemAppProxy(mobileApp, logger)

	return mobileApp
}

func (m *mobileAppProxy) CommitHandler(block poset.Block) ([]byte, error) {
	blockBytes, err := block.ProtoMarshal()
	if err != nil {
		m.logger.Debug("mobileAppProxy error marhsalling Block")
		return nil, err
	}
	stateHash := m.commitHandler.OnCommit(blockBytes)
	return stateHash, nil
}
func (m *mobileAppProxy) SnapshotHandler(blockIndex int64) ([]byte, error) {
	return []byte{}, nil
}
func (m *mobileAppProxy) RestoreHandler(snapshot []byte) ([]byte, error) {
	return []byte{}, nil
}

// CommitBlock commits a Block to the App and expects the resulting state hash
// gomobile cannot export a Block object because it doesn't support arrays of
// arrays of bytes; so we have to serialize the block.
// Overrides  InappProxy::CommitBlock
func (m *mobileAppProxy) CommitBlock(block poset.Block) ([]byte, error) {
	blockBytes, err := block.ProtoMarshal()
	if err != nil {
		m.logger.Debug("mobileAppProxy error marhsalling Block")
		return nil, err
	}
	stateHash := m.commitHandler.OnCommit(blockBytes)
	return stateHash, nil
}

//TODO - Implement these two functions
func (m *mobileAppProxy) GetSnapshot(blockIndex int64) ([]byte, error) {
	return []byte{}, nil
}

func (m *mobileAppProxy) Restore(snapshot []byte) error {
	return nil
}
