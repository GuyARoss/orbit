package dependgraph

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
)

type GraphEdge struct {
	Key   string
	Value string
}

type GraphPage struct {
	Edges []GraphEdge
}

type GraphBuilder interface {
	Graph(edges *GraphPage) error
	Write(string) error
	Renderer() error
	Dependencies() error
}

func RenderGraph(g GraphBuilder, page *GraphPage) error {
	err := g.Dependencies()
	if err != nil {
		return err
	}

	err = g.Graph(page)
	if err != nil {
		return err
	}

	return g.Renderer()
}

func varname(size int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")

	s := make([]rune, size)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}

	return string(s)
}

type CryptoScapeAVSDFGraphBuilderElementsNode struct {
}
type CryptoScapeAVSDFGraphBuilderElements struct {
	Elements []struct {
		Nodes []struct {
			Data struct {
				Id string `json:"id"`
			} `json:"data"`
		}
	} `json:"elements"`
}

type CryptoScapeAVSDFGraphBuilder struct {
	dependencies []string
	renderer     strings.Builder
	elements     strings.Builder
}

func (s *CryptoScapeAVSDFGraphBuilder) Graph(edges *GraphPage) error {
	nodes := make(map[string]bool)

	for _, e := range edges.Edges {
		nodes[e.Key] = true
	}

	return nil
}

func (s *CryptoScapeAVSDFGraphBuilder) Write(path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}

	defer out.Close()

	out.WriteString(fmt.Sprintf(`<html>
	<head>%s</head>
	<body><h1>filename</h1><div id="cy"></div><script></script></body>
	</html>
`, strings.Join(s.dependencies, " ")))

	return nil
}
func (s *CryptoScapeAVSDFGraphBuilder) Renderer() error {
	s.renderer.WriteString("document.addEventListener('DOMContentLoaded', function() {")
	s.renderer.WriteString("var cy = window.cy = cytoscape({")
	s.renderer.WriteString("container: document.getElementById('cy'),")
	s.renderer.WriteString("layout: {name: 'avsdf',nodeSeparation: 120},")
	s.renderer.WriteString("style: [{selector: 'node',style: {	'label': 'data(id)',	'text-valign': 'center',	'color': '#000000',	'background-color': '#3a7ecf'}},{selector: 'edge',style: {'width': 2,'line-color': '#3a7ecf','opacity': 0.5	}}],")

	s.renderer.WriteString("})")
	return nil
}

func (s *CryptoScapeAVSDFGraphBuilder) Dependencies() error {
	s.dependencies = append(s.dependencies, `<meta name="viewport" content="width=device-width, user-scalable=no, initial-scale=1, maximum-scale=1">`)
	s.dependencies = append(s.dependencies, `<script src="https://unpkg.com/cytoscape/dist/cytoscape.min.js"></script>`)
	s.dependencies = append(s.dependencies, `<script src="https://unpkg.com/layout-base/layout-base.js"></script>`)
	s.dependencies = append(s.dependencies, `<script src="https://unpkg.com/cytoscape-avsdf@1.0.0/cytoscape-avsdf.js"></script>`)
	s.dependencies = append(s.dependencies, `<style>body { font-family: helvetica; font-size: 15px; } h1 {opacity: 0.5;font-size: 1em;font-weight: bold;} #cy {width: 100%;height: 90%;z-index: 999;}</style>`)

	return nil
}

type DraculaGraphBuilder struct {
	edges        *strings.Builder
	dependencies *strings.Builder
	renderer     *strings.Builder

	varmap map[string]bool
}

func (g *DraculaGraphBuilder) Graph(page *GraphPage) error {
	key := varname(1)

	for s := 0; g.varmap[key]; s++ {
		key = varname(s)
	}

	g.varmap[key] = true

	g.edges.WriteString(fmt.Sprintf(`const %s = new Dracula.Graph();`, key) + "\n")
	for _, fp := range page.Edges {
		g.edges.WriteString(fmt.Sprintf(`%s.addEdge("%s", "%s");`, key, fp.Key, fp.Value) + "\n")
	}

	return nil
}

func (g *DraculaGraphBuilder) Renderer() error {
	for k := range g.varmap {
		g.renderer.WriteString(fmt.Sprintf("const layouter_%s = new Dracula.Layout.Spring(%s);\n", k, k))
		g.renderer.WriteString(fmt.Sprintf("layouter_%s.layout();\n", k))

		g.renderer.WriteString(fmt.Sprintf("const renderer_%s = new Dracula.Renderer.Raphael(document.getElementById('canvas'), %s, 1000, 1000);\n", k, k))
		g.renderer.WriteString(fmt.Sprintf("renderer_%s.draw();\n", k))
	}

	return nil
}

func (g *DraculaGraphBuilder) Dependencies() error {
	g.dependencies.WriteString(`<script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/raphael/2.3.0/raphael.min.js"></script>` + "\n")
	g.dependencies.WriteString(`<script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/graphdracula/1.2.1/dracula.min.js"></script>` + "\n")

	return nil
}

func (g *DraculaGraphBuilder) Write(path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}

	defer out.Close()

	out.WriteString(fmt.Sprintf(`<html>
	<head>%s</head>
	<body><div id="canvas"></div><script>%s %s</script></body>
	</html>
	`, g.dependencies.String(), g.edges.String(), g.renderer.String()))

	return nil
}

func NewDraculaGraph() GraphBuilder {
	return &DraculaGraphBuilder{
		edges:        &strings.Builder{},
		dependencies: &strings.Builder{},
		renderer:     &strings.Builder{},
		varmap:       make(map[string]bool),
	}
}
