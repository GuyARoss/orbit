package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// startBrowser tries to open the URL in a browser
// and reports whether it succeeds.
func startBrowser(url string) bool {
	// try to start the browser
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}

type dependgraphOut map[string][]string

func (o dependgraphOut) readFile(path string) error {
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

var dependgraphToolsCMD = &cobra.Command{
	Use:   "dependgraph",
	Long:  "visualization dependendency graph output",
	Short: "visualization dependendency graph output",
	Run: func(cmd *cobra.Command, args []string) {
		g := make(dependgraphOut)

		for _, f := range args {
			err := g.readFile(f)
			if err != nil {
				panic(err)
			}
		}

		tmpFile, err := ioutil.TempDir("", "dependgraph-*.html")
		if err != nil {
			panic(err)
		}
		out, err := os.Create(filepath.Join(tmpFile, "coverage.html"))
		if err != nil {
			panic(err)
		}

		defer out.Close()

		out.WriteString(`
		<!DOCTYPE html>
<html>

<head>
    <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/raphael/2.3.0/raphael.min.js"></script>
    <script type="text/javascript"
        src="https://cdnjs.cloudflare.com/ajax/libs/graphdracula/1.2.1/dracula.min.js"></script>
</head>

<body>
    <div id="canvas"></div>
    <script>
	var g = new Dracula.Graph();
		`)

		for ff, fp := range g {
			for _, p := range fp {
				out.WriteString(fmt.Sprintf(`g.addEdge("%s", "%s");`, ff, p) + "\n")
			}
		}

		out.WriteString(`
		var layouter = new Dracula.Layout.Spring(g);
        layouter.layout();

        var renderer = new Dracula.Renderer.Raphael(document.getElementById('canvas'), g, 1000, 1000);
        renderer.draw();
    </script>
</body>

</html>
		`)

		didStart := startBrowser("file://" + out.Name())

		fmt.Println(out.Name(), didStart)
	},
}
