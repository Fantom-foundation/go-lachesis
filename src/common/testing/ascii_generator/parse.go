// +build test

package ascii_generator

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/awalterschulze/gographviz"
)

func createParentGraph(parentId string) *gographviz.Graph {
	g := gographviz.NewGraph()

	g.SetName("\""+parentId+"\"")

	g.SetDir(true)

	return g
}

func writeDotFile(filename string, graph *gographviz.Graph) {
	fi, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fi.WriteString(graph.String())
}

// PlayGeneric struct to hold a test data
type PlayGeneric interface {
	Name() string
	SelfParent() string
	OtherParent() string
	Index() int64
}

// ConvertPlaysToAscii converts an array of play structs to ascii using graphviz and graph-easy
func ConvertPlaysToAscii(plays []interface{}, graphName string) {
	graph := createParentGraph(graphName)

	parents := make(map[string]bool)

	for _, playInterface := range plays {
		if play, ok := playInterface.(PlayGeneric); ok {
			id := play.Name()
			attrs := make(map[string]string)
			if err := graph.AddNode(graphName, id, attrs); err != nil {
				panic(err)
			}
			parents[play.Name()] = true

			// make sure parents exist before adding edges
			if play.SelfParent() != "" {
				if found := parents[play.SelfParent()]; !found {
					graph.AddNode(graphName, play.SelfParent(), nil)
				}
				if err := graph.AddEdge(play.SelfParent(), id, true, nil); err != nil {
					panic(err)
				}
			}

			if play.OtherParent() != "" {
				if found := parents[play.OtherParent()]; !found {
					graph.AddNode(graphName, play.OtherParent(), nil)
				}

				if err := graph.AddEdge(play.OtherParent(), id, true, nil); err != nil {
					panic(err)
				}
			}
		}
	}

	graph.AddAttr(graphName, string(gographviz.NewRank), "same")
	graph.AddAttr(graphName, string(gographviz.Rank), "same")

	convertFromGraphToAscii(graph, graphName)
}

func convertFromGraphToAscii(graph *gographviz.Graph, graphName string) {
	dotOutputPath := fmt.Sprintf("GRAPH_%s.dot", graphName)
	asciiOutputPath := fmt.Sprintf("GRAPH_%s_ascii.txt", graphName)
	writeDotFile(dotOutputPath, graph)
	processDotFileToAscii(dotOutputPath, asciiOutputPath)
}

func processDotFileToSvg(inputPath string, outputPath string) {
	// example convert to svg:
	//   dot -Tsvg src/poset/GRAPH_initConsensusPoset.dot -o output.svg
	out, err := exec.Command("dot", "Tsvg", inputPath, "-o", outputPath).Output()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Output:\n%s", out)
}

func processDotFileToAscii(inputPath string, outputPath string) {
	// convert to ascii:
	//   graph-easy src/poset/GRAPH_initConsensusPoset.dot GRAPH_initConsensusPoset_ascii.txt
	out, err := exec.Command("graph-easy", inputPath, outputPath).Output()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Output:\n%s", out)
}
