// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package dependgraph

import (
	"encoding/json"
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

type CryptoScapeAVSDFGraphBuilderData[T any] struct {
	Data T `json:"data"`
}

type CryptoScapeAVSDFGraphBuilderNode struct {
	ID string `json:"id"`
}

type CryptoScrapeAVSDFGraphBuilderEdge struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type CryptoScapeAVSDFGraphBuilderElements struct {
	Edges []CryptoScapeAVSDFGraphBuilderData[CryptoScrapeAVSDFGraphBuilderEdge] `json:"edges"`
	Nodes []CryptoScapeAVSDFGraphBuilderData[CryptoScapeAVSDFGraphBuilderNode]  `json:"nodes"`
}

type CryptoScapeAVSDFGraphBuilder struct {
	dependencies []string
	renderer     strings.Builder
	elements     CryptoScapeAVSDFGraphBuilderElements
}

func (s *CryptoScapeAVSDFGraphBuilder) Graph(page *GraphPage) error {
	nodeMap := make(map[string]bool)
	edges := make([]CryptoScapeAVSDFGraphBuilderData[CryptoScrapeAVSDFGraphBuilderEdge], 0)

	for _, e := range page.Edges {
		nodeMap[e.Key] = true
		nodeMap[e.Value] = true

		edges = append(edges, CryptoScapeAVSDFGraphBuilderData[CryptoScrapeAVSDFGraphBuilderEdge]{
			Data: CryptoScrapeAVSDFGraphBuilderEdge{e.Key, e.Value},
		})
	}

	nodes := make([]CryptoScapeAVSDFGraphBuilderData[CryptoScapeAVSDFGraphBuilderNode], 0)
	for k := range nodeMap {
		nodes = append(nodes, CryptoScapeAVSDFGraphBuilderData[CryptoScapeAVSDFGraphBuilderNode]{
			Data: CryptoScapeAVSDFGraphBuilderNode{k},
		})
	}

	s.elements = CryptoScapeAVSDFGraphBuilderElements{
		Nodes: nodes,
		Edges: edges,
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
	<body><h1>%s</h1><div id="cy"></div><script>%s</script></body>
	</html>
`, strings.Join(s.dependencies, " "), path, s.renderer.String()))

	return nil
}
func (s *CryptoScapeAVSDFGraphBuilder) Renderer() error {
	elements, err := json.Marshal(s.elements)
	if err != nil {
		return err
	}

	s.renderer.WriteString("document.addEventListener('DOMContentLoaded', function() {")
	s.renderer.WriteString("var cy = window.cy = cytoscape({")
	s.renderer.WriteString("container: document.getElementById('cy'),")
	s.renderer.WriteString("layout: {name: 'avsdf',nodeSeparation: 120},")
	s.renderer.WriteString("style: [{selector: 'node',style: {	'label': 'data(id)',	'text-valign': 'center',	'color': '#000000',	'background-color': '#3a7ecf'}},{selector: 'edge',style: {'width': 2,'line-color': '#3a7ecf','opacity': 0.5	}}],")
	s.renderer.WriteString(fmt.Sprintf("elements: %s", string(elements)))

	s.renderer.WriteString("})")
	s.renderer.WriteString("})")
	return nil
}

func (s *CryptoScapeAVSDFGraphBuilder) Dependencies() error {
	s.dependencies = append(s.dependencies, `<meta name="viewport" content="width=device-width, user-scalable=no, initial-scale=1, maximum-scale=1">`)
	s.dependencies = append(s.dependencies, `<script src="https://unpkg.com/cytoscape/dist/cytoscape.min.js"></script>`)
	s.dependencies = append(s.dependencies, `<script src="https://unpkg.com/layout-base/layout-base.js"></script>`)
	s.dependencies = append(s.dependencies, `<script src="https://unpkg.com/avsdf-base/avsdf-base.js"></script>`)
	s.dependencies = append(s.dependencies, `<script src="https://unpkg.com/cytoscape-avsdf@1.0.0/cytoscape-avsdf.js"></script>`)
	s.dependencies = append(s.dependencies, `<style>body { font-family: helvetica; font-size: 15px; } h1 {opacity: 0.5;font-size: 1em;font-weight: bold;} #cy {width: 100%;height: 90%;z-index: 999;}</style>`)

	return nil
}

func NewCryptoScapeAVSDFGraphBuilder() *CryptoScapeAVSDFGraphBuilder {
	return &CryptoScapeAVSDFGraphBuilder{
		dependencies: make([]string, 0),
		renderer:     strings.Builder{},
		elements:     CryptoScapeAVSDFGraphBuilderElements{},
	}
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
