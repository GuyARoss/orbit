package webwrap

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/GuyARoss/orbit/pkg/embedutils"
	"github.com/GuyARoss/orbit/pkg/fsutils"
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

func (s *JavascriptWrapper) DoesSatisfyConstraints(fileExtension string) bool {
	return fileExtension == javascriptExtension
}

func (s *JavascriptWrapper) Version() string {
	return "javascriptWebpack"
}

func (s *JavascriptWrapper) RequiredBodyDOMElements(ctx context.Context, cache *CacheDOMOpts) []string {
	return []string{}
}

func (b *JavascriptWrapper) Setup(ctx context.Context, settings *BundleOpts) ([]*BundledResource, error) {
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

	return []*BundledResource{{
		BundleFilePath:       bundleFilePath,
		ConfiguratorFilePath: fmt.Sprintf("%s/%s.config.js", b.PageOutputDir, settings.BundleKey),
		ConfiguratorPage:     page,
	}}, nil
}

func (b *JavascriptWrapper) Bundle(configuratorFilePath string) error {
	cmd := exec.Command("node", fsutils.NormalizePath(fmt.Sprintf("%s/.bin/webpack", b.NodeModulesDir)), "--config", configuratorFilePath)
	_, err := cmd.Output()

	if err != nil {
		b.Logger.Warn(fmt.Sprintf(`invalid pack: "node %s --config %s"`, fsutils.NormalizePath(fmt.Sprintf("%s/.bin/webpack", b.NodeModulesDir)), configuratorFilePath))
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
