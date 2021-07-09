package fs

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/GuyARoss/orbit/jsparse"
	"github.com/google/uuid"
)

func SetupDirs() {
	if !doesDirExist("./.orbit") {
		os.Mkdir(".orbit", 0755)
		os.Mkdir(".orbit/base", 0755)
		os.Mkdir(".orbit/dist", 0755)
		os.Mkdir(".orbit/assets", 0755)
	}
}

func applyLibTooling(dir string) *jsparse.Page {
	page, err := jsparse.ParsePage(dir)
	if err != nil {
		panic(err)
	}

	page.Imports = append(page.Imports, "import ReactDOM from 'react-dom'")
	// @@todo, could generate this element id, and pass it around.
	page.Other = append(page.Other, fmt.Sprintf("ReactDOM.render(<%s {...JSON.parse(document.getElementById('orbit_manifest').textContent)}/>, document.getElementById('root'))", page.Name))

	page.WriteFile(dir)

	return page
}

type bundlerOut struct {
	BundlerConfigPath string
	BundleName        string
}

func setupPageBundler(dir string, fileName string, name string) *bundlerOut {
	page := jsparse.Page{}
	page.Imports = append(page.Imports, "const {merge} = require('webpack-merge')")
	page.Imports = append(page.Imports, "const baseConfig = require('../../assets/base.config.js')")

	outputFileName := fmt.Sprintf("%s.js", name)

	page.Other = append(page.Other, fmt.Sprintf(`module.exports = merge(baseConfig, {
		entry: ['./%s'],
		output: {
			filename: '%s'
		},
	})`, fileName, outputFileName))
	configPath := fmt.Sprintf("%s/%s.config.js", dir, name)

	page.WriteFile(configPath)

	return &bundlerOut{
		BundlerConfigPath: configPath,
		BundleName:        outputFileName,
	}
}

func bundle(bundleFile string) error {
	cmd := exec.Command("bash", "node_modules/.bin/webpack", "--config", bundleFile)
	_, err := cmd.Output()

	return err
}

type PackedPage struct {
	PageName  string
	BundleKey string
}

func hashKey(idx int, name string) string {
	id := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(fmt.Sprintf("%d-%s", idx, name)))

	return strings.ReplaceAll(id.String(), "-", "")
}

func Pack(baseDir string, bundleOut string) []*PackedPage {
	dirs := copyDir(baseDir, baseDir, ".orbit/base")
	copyDir("assets", "assets", ".orbit/assets")

	pages := make([]*PackedPage, 0)
	for idx, dir := range dirs {
		if strings.Contains(dir, "pages") {
			page := applyLibTooling(dir)

			bundleKey := hashKey(idx, page.Name)
			err := os.Rename(dir, fmt.Sprintf("%s/%s.js", bundleOut, bundleKey))

			if err != nil {
				panic(err)
			}
			buildOut := setupPageBundler(bundleOut, fmt.Sprintf("%s/%s.js", bundleOut, bundleKey), bundleKey)
			bundleErr := bundle(buildOut.BundlerConfigPath)
			if bundleErr != nil {
				panic(bundleErr)
			}

			// @@todo(debug)
			fmt.Printf("successfully packed %s \n", page.Name)

			pages = append(pages, &PackedPage{
				PageName:  page.Name,
				BundleKey: bundleKey,
			})
		}
	}

	return pages
}
