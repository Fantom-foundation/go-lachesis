// +build debug

// Package node these functions are used only in debugging
package node

import (
//	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/Fantom-foundation/go-lachesis/src/poset"
	"github.com/sirupsen/logrus"
	"github.com/tebeka/atexit"
)

// InfosLite small subset of debug info for node
type InfosLite struct {
	ParticipantEvents map[string]map[string]EventLite
	Rounds            []poset.RoundCreated
	Blocks            []poset.Block
}

// EventBodyLite small subset of event body for debugging
type EventBodyLite struct {
	Parents      [][]byte // hashes of the event's parents, self-parent first
	Creator      string   // creator's NetAddr //public key
	Index        int64    // index in the sequence of events created by Creator
	Transactions [][]byte
}

// EventMessageLite small subset of event body for debugging
type EventMessageLite struct {
	Body             EventBodyLite
	Signature        string // creator's digital signature of body
	TopologicalIndex int64
	Hash             string
	Frame            int64
	RoundReceived    int64

	ClothoProof []string
	// FlagTable []byte // FlagTable stores connection information
}

// EventLite small subset of event for debugging
type EventLite struct {
	CreatorID            uint64
	OtherParentCreatorID uint64
	Root                 bool
	Clotho               bool
	Atropos              bool
	LamportTimestamp     int64
	AtroposTimestamp     int64
	Message              EventMessageLite
	FlagTableBytes       []byte // FlagTable stores connection information
	RootTableBytes       []byte // FlagTable stores connection information
	AtTimes              []int64
	AtVisited            int64
	RecFrames            []int64
	FrameReceived        int64
}

// GetParticipantEventsLite returns all participants
func (g *Graph) GetParticipantEventsLite() map[string]map[string]EventLite {
	res := make(map[string]map[string]EventLite)

	store := g.Node.core.poset.Store
	peers := g.Node.core.poset.Participants

	// evs, err := store.ParticipantEvents(p.PubKeyHex, root.SelfParent.Index)
	evs, err := store.TopologicalEvents()

	if err != nil {
		panic(err)
	}

	res[g.Node.localAddr /*p.PubKeyHex*/] = make(map[string]EventLite)

	for _, event := range evs {

		if err != nil {
			panic(err)
		}

		peer, ok := peers.ReadByPubKey(event.GetCreator())
		if !ok {
			panic(fmt.Sprintf("Creator %v not found", event.GetCreator()))
		}
		creatorParts := strings.Split(peer.Message.NetAddr, ":")

		hash := event.Hash()

		liteEvent := EventLite{
			CreatorID:            event.CreatorID(),
			OtherParentCreatorID: event.OtherParentCreatorID(),
			Root: event.Root,
			Clotho: event.Clotho,
			Atropos: event.Atropos,
			LamportTimestamp: event.LamportTimestamp,
			AtroposTimestamp: event.AtroposTimestamp,
			FlagTableBytes: event.FlagTableBytes,
			RootTableBytes: event.RootTableBytes,
			AtTimes: event.AtTimes,
			AtVisited: event.AtVisited,
			RecFrames: event.RecFrames,
			FrameReceived: event.FrameReceived,
			Message: EventMessageLite{
				Body: EventBodyLite{
					Parents:      event.Message.Body.Parents,
					Creator:      creatorParts[1], //peer.NetAddr,
					Index:        event.Message.Body.Index,
					Transactions: event.Message.Body.Transactions,
				},
				Hash:             hash.String(),
				Signature:        event.Message.Signature,
//				ClothoProof:      event.Message.ClothoProof,
				Frame:            event.Frame,
//				RoundReceived:    event.Message.RoundReceived,
				TopologicalIndex: event.Message.TopologicalIndex,
				// 				FlagTable: event.FlagTable,
			},
		}

		res[g.Node.localAddr /*p.PubKeyHex*/][hash.String()] = liteEvent
	}

	return res
}

// GetInfosLite returns debug stats
func (g *Graph) GetInfosLite() InfosLite {
	return InfosLite{
		ParticipantEvents: g.GetParticipantEventsLite(),
		Rounds:            g.GetRounds(),
		Blocks:            g.GetBlocks(),
	}
}

// PrintStat prints debug stats
func (c *Core) PrintStat(logger *logrus.Entry) {
	logger.Warn("**core.HexID=", c.HexID())
	logger.Warn("****Known events:")
	for pidID, index := range c.KnownEvents() {
		peer, ok := c.participants.ReadByID(uint64(pidID))
		if ok {
			logger.Warn("    index=", index, " peer=", peer.Message.NetAddr,
				" pubKeyHex=", peer.Message.PubKeyHex)
		}
	}
//	c.poset.PrintStat(logger)
}

// PrintStat prints debug stats
func (n *Node) PrintStat() {
	n.logger.Warn("*Node=", n.localAddr)
	g := NewGraph(n)
//	encoder := json.NewEncoder(n.logger.Logger.Out)
//	encoder.SetIndent("", "  ")
	res := g.GetInfosLite()
//	encoder.Encode(res)
	n.core.PrintStat(n.logger)
	file, err := os.Create(fmt.Sprintf("Node_%v.gv", n.localAddr))
	if err != nil {
		panic(err)
	}
	graf, err := os.Create(fmt.Sprintf("Node_%v.graph", n.localAddr))
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(file, "digraph ANode { /* %v */\nrankdir=RL; ranksep=1.5; splines=false;\n", n.localAddr)
	fmt.Fprintf(graf, "// Node: %v\n#Vertexes\n", n.localAddr)
	fmt.Fprintf(graf, "// ID frame creator LamportTimetamp AtroposTimestamp root clotho atropos\n")

	fr := make(map[int64]map[string][]EventLite)
	cr := make(map[string][]EventLite)
	maxFrame := int64(0)
	for _, events := range res.ParticipantEvents {
		for _, le := range events {
			if le.Message.Frame > maxFrame {
				maxFrame = le.Message.Frame
			}
			fr[le.Message.Frame] = make(map[string][]EventLite)
			cr[le.Message.Body.Creator] = append(cr[le.Message.Body.Creator], le)
		}
	}

	for _, events := range res.ParticipantEvents {
		for _, le := range events {
			fr[le.Message.Frame][le.Message.Body.Creator] = append(fr[le.Message.Frame][le.Message.Body.Creator], le)
		}
	}


	fmt.Fprint(file, "layers = \"f0")
	for i := int64(1); i <= maxFrame; i++ {
		fmt.Fprintf(file, ":f%v", i)
	}
	fmt.Fprintln(file, "\";")


	var keys []string
	for k, _ := range cr {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, creator := range keys {
		lightEvents := cr[creator]
		fmt.Fprintf(file, "subgraph cluster_%v { rank = same; ranksep = 2.5; ", creator)
		for _, le := range lightEvents {
			fmt.Fprintf(file, "v%v [shape=none,layer=\"f%v\" label=<<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"4\"><TR><TD>f</TD><td>l</td><td>a</td><td>atr</td><td>cl</td><td>roo</td><td>cr</td></TR><tr><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td></tr>",
				le.Message.Hash, le.Message.Frame, le.Message.Frame, le.LamportTimestamp, le.AtroposTimestamp, le.Atropos, le.Clotho, le.Root, le.Message.Body.Creator )

			fmt.Fprintf(file, "<tr><td>%v</td><td colspan=\"6\">at:", le.AtVisited);
			for k, v := range le.AtTimes {
				if k > 0 {
					fmt.Fprintf(file, ", ")
				}
				fmt.Fprintf(file, "%v", v)
			}
			fmt.Fprintf(file, "</td></tr>")

			fmt.Fprintf(file, "<tr><td>%v</td><td colspan=\"6\">fr:", le.FrameReceived);
			for k, v := range le.RecFrames {
				if k > 0 {
					fmt.Fprintf(file, ": ")
				}
				fmt.Fprintf(file, "%v", v)
			}
			fmt.Fprintf(file, "</td></tr>")

			// Uncomment the following if FlagTable and RootTable are needed in the node visualisation
/*
			fmt.Frpintf(file, "<tr><td colspan=\"7\">ft:");
			ft := poset.NewFlagTable()
			ft.Unmarshal(le.FlagTableBytes)
			for k, v := range ft {
				ev, err := n.core.poset.Store.GetEventBlock(k)
				if err != nil {
					panic(err)
				}
				peer, ok := n.core.poset.Participants.ReadByPubKey(ev.GetCreator())
				if !ok {
					panic(fmt.Sprintf("Peer %v not found", k))
				}
				creatorParts := strings.Split(peer.NetAddr, ":")
				fmt.Fprintf(file, " %v:%v", creatorParts[1], v)
			}

			fmt.Fprintf(file, "</td></tr><tr><td colspan=\"7\">rt:")
			ft = poset.NewFlagTable()
			ft.Unmarshal(le.RootTableBytes)
			for k, v := range ft {
				ev, err := n.core.poset.Store.GetEventBlock(k)
				if err != nil {
					panic(err)
				}
				peer, ok := n.core.poset.Participants.ReadByPubKey(ev.GetCreator())
				if !ok {
					panic(fmt.Sprintf("Peer %v not found", k))
				}
				creatorParts := strings.Split(peer.NetAddr, ":")
				fmt.Fprintf(file, " %v:%v", creatorParts[1], v)
			}

			fmt.Fprintf(file, "</td></tr>")
*/
			fmt.Fprintf(file, "</TABLE>>];\n")
			n.logger.Warnf("v%v f:%v cl:%v roo:%v cr:%v l:%v", le.Message.Hash, le.Message.Frame, le.Clotho, le.Root, le.Message.Body.Creator, le.LamportTimestamp )
			fmt.Fprintf(graf, "%v %v %v %v %v %v %v %v\n", le.Message.Hash, le.Message.Frame, le.Message.Body.Creator, le.LamportTimestamp, le.AtroposTimestamp, le.Root, le.Clotho, le.Atropos)
		}
		fmt.Fprint(file, " }\n")
	}

/*
	for i := int64(0); i <= maxFrame; i++ {
		fmt.Fprintf(file, "subgraph cluster_f%v { rank = same; ", i)
		var keys []string
		for k, _ := range fr[i] {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, cr := range keys {
			lightEvents := fr[i][cr]
			fmt.Fprintf(file, "subgraph cluster%v_%v { rankdir=TB; rank = same; ", i, cr)
			for _, le := range lightEvents {
				fmt.Fprintf(file, "v%v [shape=none,layer=\"f%v\" label=<<TABLE BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\" CELLPADDING=\"4\"><TR><TD>f</TD><td>l</td><td>cl</td><td>roo</td><td>cr</td></TR><tr><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td></tr></TABLE>>];\n",
					le.Message.Hash, le.Message.Frame, le.Message.Frame, le.LamportTimestamp, le.Clotho, le.Root, le.Message.Body.Creator )
			}
			fmt.Fprint(file, " }\n")
		}
		fmt.Fprint(file, " }\n")
	}
*/

	fmt.Fprintf(graf, "#Edges\n//From To\n")

	for _/*localAddr*/, events := range res.ParticipantEvents {
		for _/*hash*/, le := range events {
			var parent, otherParent poset.EventHash
			parent.Set(le.Message.Body.Parents[0])
			otherParent.Set(le.Message.Body.Parents[1])
			if !parent.Zero() {
				fmt.Fprintf(file, "v%v -> v%v;\n", le.Message.Hash, parent.String())
			}
			if !otherParent.Zero() {
				fmt.Fprintf(file, "v%v -> v%v;\n", le.Message.Hash, otherParent.String())
			}
			fmt.Fprintf(graf, "%v %v\n", parent.String(), le.Message.Hash)
			fmt.Fprintf(graf, "%v %v\n", otherParent.String(), le.Message.Hash)
//			fmt.Fprintf(file, "v%v [shape=record,layer=\"f%v\" label=\"f:%v | l:%v | c:%v | r:%v | cr:%v\"];\n",
//				le.Message.Hash, le.Message.Frame, le.Message.Frame, le.LamportTimestamp, le.Clotho, le.Root, le.Message.Body.Creator )
		}
	}
	fmt.Fprintln(file, "}\n")
	file.Close()
	graf.Close()
}

// Register a print listener
func (n *Node) Register() {
	var once sync.Once
	onceBody := func() {
		file, err := os.OpenFile("lachesis.trace", os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("*** Open  err: %v", err)
			return
		}
		defer file.Close()
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, true)
		fmt.Fprintln(file, string(buf[0:stackSize]))
		// You must build with debug tag to have PrintStat() defined
		n.PrintStat()
	}
	atexit.Register(func() {
		once.Do(onceBody)
	})
	// use the following way of exit to execute registered atexit handlers:
	// import "github.com/tebeka/atexit"
	// atexit.Exit(0)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		atexit.Exit(13)
	}()
}
