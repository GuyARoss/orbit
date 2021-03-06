// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package webwrap

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/GuyARoss/orbit/pkg/embedutils"
	"github.com/GuyARoss/orbit/pkg/jsparse"
	parseerror "github.com/GuyARoss/orbit/pkg/parse_error"
)

type ReactWebWrapper struct {
	*BaseWebWrapper
	*BaseBundler
}

var ErrComponentExport = errors.New("prefer capitalization for jsx components")
var ErrInvalidComponent = errors.New("invalid jsx component")

const reactExtension string = "jsx"

func (s *ReactWebWrapper) Apply(page jsparse.JSDocument) (jsparse.JSDocument, error) {
	if page.Extension() != reactExtension { // @@todo bad pattern fix this
		return nil, ErrInvalidComponent
	}

	// react components should always be capitalized.
	if string(page.Name()[0]) != strings.ToUpper(string(page.Name()[0])) {
		return nil, ErrComponentExport
	}

	page.AddImport(&jsparse.ImportDependency{
		FinalStatement: "import ReactDOM from 'react-dom'",
		Type:           jsparse.ModuleImportType,
	})

	page.AddOther(fmt.Sprintf(
		"ReactDOM.render(<%s {...JSON.parse(document.getElementById('orbit_manifest').textContent)}/>, document.getElementById('%s_react_frame'))",
		page.Name(), page.Key()),
	)

	return page, nil
}

func (r *ReactWebWrapper) VerifyRequirements() error {
	webpackPath := fmt.Sprintf("%s%c%s%c%s", r.NodeModulesDir, os.PathSeparator, ".bin", os.PathSeparator, "webpack")

	// due to a "bug" with windows, it has an issue with shebang cmds, so we prefer the webpack.js file instead.
	if runtime.GOOS == "windows" {
		webpackPath = r.NodeModulesDir + "/webpack/bin/webpack.js"
	}

	_, err := os.Stat(webpackPath)
	if err != nil {
		return fmt.Errorf("node module not found: webpack. It is possible that you need to run `npm i` in your workspace directory to remedy this issue")
	}

	return nil
}

func (s *ReactWebWrapper) DoesSatisfyConstraints(fileExtension string) bool {
	return fileExtension == reactExtension
}

func (s *ReactWebWrapper) Version() string {
	return "reactManifestFallback"
}

func (s *ReactWebWrapper) Stats() *WrapStats {
	return &WrapStats{
		WebVersion: "react",
		Bundler:    "webpack",
	}
}

func (s *ReactWebWrapper) RequiredBodyDOMElements(ctx context.Context, cache *CacheDOMOpts) []string {
	mode := ctx.Value(BundlerID).(string)

	uris := make([]string, 0)
	switch BundlerMode(mode) {
	case DevelopmentBundle:
		uris = append(uris, "https://unpkg.com/react/umd/react.development.js")
		uris = append(uris, "https://unpkg.com/react-dom/umd/react-dom.development.js")
	case ProductionBundle:
		uris = append(uris, "https://unpkg.com/react/umd/react.production.min.js")
		uris = append(uris, "https://unpkg.com/react-dom/umd/react-dom.production.min.js")
	}

	files, _ := cache.CacheWebRequest(uris)

	// currently these files are just paths to a directory to refer
	// to them on the dom, we need to convert them to <script> tags.
	for i, f := range files {
		files[i] = fmt.Sprintf(`<script src="%s"></script>`, f)
	}

	return files
}

func (b *ReactWebWrapper) Setup(ctx context.Context, settings *BundleOpts) (*BundledResource, error) {
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

func (b *ReactWebWrapper) Bundle(configuratorFilePath string, filePath string) error {
	webpackPath := fmt.Sprintf("%s%c%s%c%s", b.NodeModulesDir, os.PathSeparator, ".bin", os.PathSeparator, "webpack")

	// due to a "bug" with windows, it has an issue with shebang cmds, so we prefer the webpack.js file instead.
	if runtime.GOOS == "windows" {
		webpackPath = b.NodeModulesDir + "/webpack/bin/webpack.js"
	}

	cmd := exec.Command("node", webpackPath, "--config", configuratorFilePath)
	output, err := cmd.Output()

	if err != nil {
		b.Logger.Warn(fmt.Sprintf(`invalid pack: "node %s --config %s"\n "%s"`, webpackPath, configuratorFilePath, string(output)))
		return parseerror.New("failed to bundle, this could denote a syntax error", filePath)
	}

	return nil
}

func (b *ReactWebWrapper) HydrationFile() []embedutils.FileReader {
	files, err := embedFiles.ReadDir("embed")
	if err != nil {
		return nil
	}

	for _, file := range files {
		if strings.Contains(file.Name(), "react_hydrate.go") {
			return []embedutils.FileReader{&embedFileReader{fileName: file.Name()}}
		}
	}
	return nil
}

func NewReactWebWrap(bundler *BaseBundler) *ReactWebWrapper {
	return &ReactWebWrapper{
		BaseBundler: bundler,
	}
}
