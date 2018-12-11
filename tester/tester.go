package tester

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	_ "os"
	_ "sync"

	"github.com/Fantom-foundation/go-lachesis/src/peers"
	"github.com/Fantom-foundation/go-lachesis/src/proxy"
	"github.com/sirupsen/logrus"
)

func PingNodesN(participants []*peers.Peer, p peers.PubKeyPeers, n uint64, delay uint64, logger *logrus.Logger, ProxyAddr string) {
	// pause before shooting test transactions
	time.Sleep(time.Duration(delay) * time.Second)

	proxies := make(map[int64]*proxy.GrpcLachesisProxy)
	for _, participant := range participants {
		node := p[participant.PubKeyHex]
		if node.NetAddr == "" {
			fmt.Printf("node missing NetAddr [%v]", node)
			continue
		}
		host_port := strings.Split(node.NetAddr, ":")
		port, err := strconv.Atoi(host_port[1])
		addr := fmt.Sprintf("%s:%d", host_port[0], port-3000 /*9000*/)
		lachesisProxy, err := proxy.NewGrpcLachesisProxy(addr, logger)
		if err != nil {
			fmt.Printf("error:\t\t\t%s\n", err.Error())
			fmt.Printf("Failed to create WebsocketLachesisProxy:\t\t\t%s (id=%d)\n", participant.NetAddr, node.ID)
		}
		proxies[node.ID] = lachesisProxy
	}
	for iteration := uint64(0); iteration < n; iteration++ {
		participant := participants[rand.Intn(len(participants))]
		node := p[participant.PubKeyHex]

		_, err := transact(proxies[node.ID], ProxyAddr, iteration)

		if err != nil {
			fmt.Printf("error:\t\t\t%s\n", err.Error())
			fmt.Printf("Failed to ping:\t\t\t%s (id=%d)\n", participant.NetAddr, node.ID)
			fmt.Printf("Failed to send transaction:\t%d\n\n", iteration)
		} /*else {
			fmt.Printf("Pinged:\t\t\t%s (id=%d)\n", participant.NetAddr, node)
			fmt.Printf("Last transaction sent:\t%d\n\n", iteration)
		}*/
	}

	for _, lachesisProxy := range proxies {
		lachesisProxy.Close()
	}
	fmt.Println("Pinging stopped after ", n, " iterations")
}

func transact(proxy *proxy.GrpcLachesisProxy, proxyAddr string, iteration uint64) (string, error) {

	// Ethereum txns are ~108 bytes. Bitcoin txns are ~250 bytes.
	// A good assumption is to make txns 120 bytes in size.
	// However, for speed, we're using 1 byte here. Modify accordingly.
	//var msg [1]byte
	for i := 0; i < 10; i++ {
		// Send 10 txns to the server.
		msg := fmt.Sprintf("%s.%d.%d", proxyAddr, iteration, i)
		err := proxy.SubmitTx([]byte(msg))
		if err != nil {
			return "", err
		}
	}
	// fmt.Println("Submitted tx, ack=", ack)  # `ack` is now `_`

	return "", nil
}
