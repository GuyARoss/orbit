package dependgraph

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CMD = &cobra.Command{
	Use:   "dependgraph",
	Long:  "visualization dependendency graph output",
	Short: "visualization dependendency graph output",
	Run: func(cmd *cobra.Command, args []string) {
		g := make(dependgraph)

		for _, f := range args {
			err := g.ReadFile(f)
			if err != nil {
				panic(err)
			}
		}

		tmpFile, err := ioutil.TempDir("", "dependgraph-*.html")
		if err != nil {
			panic(err)
		}

		if err != nil {
			panic(err)
		}

		var graph GraphBuilder
		switch viper.GetString("graph") {
		case "avsd":
			graph = NewCryptoScapeAVSDFGraphBuilder()
		case "dracula":
			graph = NewDraculaGraph()
		default:
			panic(fmt.Sprintf("invalid error mode '%s'", viper.GetString("graph")))
		}

		err = RenderGraph(graph, &GraphPage{
			Edges: g.Edges(),
		})

		if err != nil {
			panic(err)
		}

		path := filepath.Join(tmpFile, "dependgraph.html")
		err = graph.Write(path)
		if err != nil {
			panic(err)
		}

		didStart := startBrowser(fmt.Sprintf("file://%s", path))
		if !didStart {
			fmt.Printf("visit this link to view the graph file://%s", path)
		}
	},
}

func init() {
	var graphmode string

	CMD.PersistentFlags().StringVar(&graphmode, "graph", "avsd", "specifies the graph mode used for building the graph output")
	viper.BindPFlag("graph", CMD.PersistentFlags().Lookup("graph"))
}
