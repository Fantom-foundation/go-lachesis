package lachesis

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"

	"github.com/andrecronje/lachesis/src/poset"
	"github.com/andrecronje/lachesis/src/proxy"
	"github.com/sirupsen/logrus"
)

type SocketLachesisProxyServer struct {
	netListener *net.Listener
	rpcServer   *rpc.Server
	handler     proxy.ProxyHandler
	timeout     time.Duration
	logger      *logrus.Logger
}

func NewSocketLachesisProxyServer(
	bindAddress string,
	handler proxy.ProxyHandler,
	timeout time.Duration,
	logger *logrus.Logger,
) (*SocketLachesisProxyServer, error) {

	server := &SocketLachesisProxyServer{
		handler: handler,
		timeout: timeout,
		logger:  logger,
	}

	if err := server.register(bindAddress); err != nil {
		return nil, err
	}

	return server, nil
}

func (p *SocketLachesisProxyServer) register(bindAddress string) error {
	rpcServer := rpc.NewServer()
	rpcServer.RegisterName("State", p)

	p.rpcServer = rpcServer

	l, err := net.Listen("tcp", bindAddress)

	if err != nil {
		return err
	}

	p.netListener = &l

	return nil
}

func (p *SocketLachesisProxyServer) listen() error {
	for {
		conn, err := (*p.netListener).Accept()

		if err != nil {
			return err
		}

		go (*p.rpcServer).ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

func (p *SocketLachesisProxyServer) CommitBlock(block poset.Block, stateH *[]byte) (err error) {
	*stateH, err = p.handler.CommitHandler(block)

	p.logger.WithFields(logrus.Fields{
		"block":      block.Index(),
		"state_hash": stateH,
		"err":        err,
	}).Debug("LachesisProxyServer.CommitBlock")

	if err != nil {
		return err
	}

	return
}

func (p *SocketLachesisProxyServer) GetSnapshot(blockIndex int, snapshot *[]byte) (err error) {
	*snapshot, err = p.handler.SnapshotHandler(blockIndex)

	if err != nil {
		return err
	}

	p.logger.WithFields(logrus.Fields{
		"block":    blockIndex,
		"snapshot": snapshot,
		"err":      err,
	}).Debug("LachesisProxyServer.GetSnapshot")

	return
}

func (p *SocketLachesisProxyServer) Restore(snapshot []byte, stateHash *[]byte) (err error) {
	*stateHash, err = p.handler.RestoreHandler(snapshot)

	if err != nil {
		return err
	}

	p.logger.WithFields(logrus.Fields{
		"state_hash": stateHash,
		"err":        err,
	}).Debug("LachesisProxyServer.Restore")

	return
}
