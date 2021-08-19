package fs

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/GuyARoss/orbit/pkg/jsparse"
	"github.com/google/uuid"
)

func SetupDirs() {
	if !DoesDirExist("./.orbit") {
		os.Mkdir(".orbit", 0755)
		os.Mkdir(".orbit/base", 0755)
		os.Mkdir(".orbit/dist", 0755)
		os.Mkdir(".orbit/assets", 0755)
	}
}

type bundlerOut struct {
	BundlerConfigPath string
	BundleName        string
}

type BundlerMode string

const (
	ProductionBundle  BundlerMode = "production"
	DevelopmentBundle BundlerMode = "development"
)

type BundlerSettings struct {
	Mode           BundlerMode
	NodeModulePath string

	WebDir string
}

func (s *BundlerSettings) setupPageBundler(dir string, fileName string, name string) *bundlerOut {
	page := jsparse.Page{}
	page.Imports = append(page.Imports, "const {merge} = require('webpack-merge')")

	// will this cause an issue? hmm..
	page.Imports = append(page.Imports, "const baseConfig = require('../../assets/base.config.js')")

	outputFileName := fmt.Sprintf("%s.js", name)

	page.Other = append(page.Other, fmt.Sprintf(`module.exports = merge(baseConfig, {
		entry: ['./%s'],
		mode: '%s',
		output: {
			filename: '%s'
		},
	})`, fileName, string(s.Mode), outputFileName))
	configPath := fmt.Sprintf("%s/%s.config.js", dir, name)

	page.WriteFile(configPath)

	return &bundlerOut{
		BundlerConfigPath: configPath,
		BundleName:        outputFileName,
	}
}

func bundle(bundleFile string, nodeModuleDir string) error {
	cmd := exec.Command("node", fmt.Sprintf("%s/.bin/webpack", nodeModuleDir), "--config", bundleFile)
	_, err := cmd.Output()

	fmt.Printf("\n\nnode %s --config %s \n\n\n\n", fmt.Sprintf("%s/.bin/webpack", nodeModuleDir), bundleFile)

	return err
}

type PackedPage struct {
	PageName  string
	BundleKey string
	BaseDir   string
}

func hashKey(idx int, name string) string {
	id := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(fmt.Sprintf("%d-%s", idx, name)))

	return strings.ReplaceAll(id.String(), "-", "")
}

type PackSettings struct {
	*BundlerSettings
	AssetDir string
}

func (s *PackSettings) applyLibTooling(dir string) *jsparse.Page {
	page, err := jsparse.ParsePage(dir, s.WebDir)
	if err != nil {
		panic(err)
	}

	page.Imports = append(page.Imports, "import ReactDOM from 'react-dom'")
	// @@todo, could generate this element id, and pass it around.
	page.Other = append(page.Other, fmt.Sprintf("ReactDOM.render(<%s {...JSON.parse(document.getElementById('orbit_manifest').textContent)}/>, document.getElementById('root'))", page.Name))

	page.WriteFile(dir)

	return page
}

func (s *PackSettings) Pack(baseDir string, bundleOut string) *[]*PackedPage {
	dirs := copyDir(baseDir, baseDir, ".orbit/base", true)
	copyDir(s.AssetDir, s.AssetDir, ".orbit/assets", false)

	pages := make([]*PackedPage, 0)
	for idx, dir := range dirs {
		if strings.Contains(dir.CopyDir, "pages") {
			page := s.applyLibTooling(dir.CopyDir)

			bundleKey := hashKey(idx, page.Name)
			err := os.Rename(dir.CopyDir, fmt.Sprintf("%s/%s.js", bundleOut, bundleKey))

			if err != nil {
				panic(err)
			}

			buildOut := s.BundlerSettings.setupPageBundler(bundleOut, fmt.Sprintf("%s/%s.js", bundleOut, bundleKey), bundleKey)
			bundleErr := bundle(buildOut.BundlerConfigPath, s.NodeModulePath)

			if bundleErr != nil {
				panic(bundleErr)
			}

			// @@todo(debug)
			fmt.Printf("successfully packed %s %s \n", page.Name, dir.BaseDir)

			pages = append(pages, &PackedPage{
				PageName:  page.Name,
				BundleKey: bundleKey,
				BaseDir:   dir.BaseDir,
			})
		}
	}

	return &pages
}
