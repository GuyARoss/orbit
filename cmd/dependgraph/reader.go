// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package dependgraph

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type dependgraph map[string][]string

func (o dependgraph) Edges() []GraphEdge {
	edges := make([]GraphEdge, 0)

	for ff, fp := range o {
		for _, p := range fp {
			edges = append(edges, GraphEdge{
				Key:   ff,
				Value: p,
			})
		}
	}

	return edges
}

func (o dependgraph) ReadFile(path string) error {
	file, err := os.Open(path)

	if err != nil {
		return err
	}

	defer file.Close()

	confirmMode := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		splitLine := strings.Split(scanner.Text(), " ")

		if !confirmMode {
			if splitLine[1] != "graph" {
				return fmt.Errorf("received invalid mode type '%s'", splitLine[2])
			}

			confirmMode = true
			continue
		}

		if o[splitLine[0]] == nil {
			o[splitLine[0]] = make([]string, 0)
		}

		o[splitLine[0]] = append(o[splitLine[0]], splitLine[1])
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
