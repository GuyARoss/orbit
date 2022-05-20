// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package webwrap

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/GuyARoss/orbit/pkg/embedutils"
	"github.com/GuyARoss/orbit/pkg/jsparse"
)

type JavascriptWrapper struct {
	*BaseWebWrapper
	*BaseBundler
}

const javascriptExtension string = "js"

func (s *JavascriptWrapper) Apply(page jsparse.JSDocument) (jsparse.JSDocument, error) {
	if page.Extension() != javascriptExtension { // @@todo bad pattern fix this
		return nil, fmt.Errorf("invalid extension %s", page.Extension())
	}

	page.AddOther(fmt.Sprintf(
		`onLoadTasks.push(
			() => %s({...JSON.parse(document.getElementById('orbit_manifest').textContent)})
		)`,
		page.Name()),
	)

	return page, nil
}

func (b *JavascriptWrapper) VerifyRequirements() error {
	webpackPath := fmt.Sprintf("%s%c%s%c%s", b.NodeModulesDir, os.PathSeparator, ".bin", os.PathSeparator, "webpack")

	// due to a "bug" with windows, it has an issue with shebang cmds, so we prefer the webpack.js file instead.
	if runtime.GOOS == "windows" {
		webpackPath = b.NodeModulesDir + "/webpack/bin/webpack.js"
	}

	_, err := os.Stat(webpackPath)
	if err != nil {
		return fmt.Errorf("node module not found: webpack. It is possible that you need to run `npm i` in your workspace directory to remedy this issue.")
	}

	return nil
}

func (s *JavascriptWrapper) DoesSatisfyConstraints(fileExtension string) bool {
	return fileExtension == javascriptExtension
}

func (s *JavascriptWrapper) Version() string {
	return "javascriptWebpack"
}

func (s *JavascriptWrapper) RequiredBodyDOMElements(ctx context.Context, cache *CacheDOMOpts) []string {
	return []string{
		`<script> const onLoadTasks = []; window.onload = (e) => { onLoadTasks.forEach(t => t(e))} </script>`,
	}
}

func (b *JavascriptWrapper) Setup(ctx context.Context, settings *BundleOpts) (*BundledResource, error) {
	page := jsparse.NewEmptyDocument()

	page.AddImport(&jsparse.ImportDependency{
		FinalStatement: "const {merge} = require('webpack-merge')",
		Type:           jsparse.ModuleImportType,
	})

	page.AddImport(&jsparse.ImportDependency{
		FinalStatement: "const baseConfig = require('../../assets/base.config.js')",
		Type:           jsparse.ModuleImportType,
	})

	outputFileName := fmt.Sprintf("%s.js", settings.BundleKey)
	bundleFilePath := fmt.Sprintf("%s/%s.js", b.PageOutputDir, settings.BundleKey)

	page.AddOther(fmt.Sprintf(`module.exports = merge(baseConfig, {
		entry: ['./%s'],
		mode: '%s',
		output: {
			filename: '%s'
		},
	})`, bundleFilePath, string(b.Mode), outputFileName))

	return &BundledResource{
		BundleFilePath: bundleFilePath,
		Configurators: []BundleConfigurator{
			{
				FilePath: fmt.Sprintf("%s/%s.config.js", b.PageOutputDir, settings.BundleKey),
				Page:     page,
			},
		},
	}, nil
}

func (b *JavascriptWrapper) Bundle(configuratorFilePath string) error {
	webpackPath := fmt.Sprintf("%s%c%s%c%s", b.NodeModulesDir, os.PathSeparator, ".bin", os.PathSeparator, "webpack")

	// due to a "bug" with windows, it has an issue with shebang cmds, so we prefer the webpack.js file instead.
	if runtime.GOOS == "windows" {
		webpackPath = b.NodeModulesDir + "/webpack/bin/webpack.js"
	}

	cmd := exec.Command("node", webpackPath, "--config", configuratorFilePath)

	_, err := cmd.Output()
	if err != nil {
		b.Logger.Warn(fmt.Sprintf(`invalid pack: "node %s --config %s"`, fmt.Sprintf("%s/.bin/webpack", b.NodeModulesDir), configuratorFilePath))
	}

	return err
}

func (b *JavascriptWrapper) HydrationFile() []embedutils.FileReader {
	files, err := embedFiles.ReadDir("embed")
	if err != nil {
		return nil
	}

	for _, file := range files {
		if strings.Contains(file.Name(), "javascript.go") {
			return []embedutils.FileReader{&embedFileReader{fileName: file.Name()}}
		}
	}
	return nil
}
