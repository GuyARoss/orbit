// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.
package bundler

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/GuyARoss/orbit/pkg/fsutils"
	"github.com/GuyARoss/orbit/pkg/jsparse"
)

type WebPackBundler struct {
	*BaseBundler
}

func (b *WebPackBundler) Setup(ctx context.Context, settings *BundleOpts) (*BundledResource, error) {
	page := jsparse.NewEmptyDocument()

	page.AddImport(&jsparse.ImportDependency{
		FinalStatement: "const {merge} = require('webpack-merge')",
		Type:           jsparse.ModuleImportType,
	})

	// @@todo(guy): this webpack config is currently based off of react, if we want to add support in the future
	// we will need to update this to apply a type context depending on which of the frontend frameworks are selected.
	// * we could also parse the file to determine which of the front-end frameworks are attached. then use the correct config *
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
		BundleFilePath:       bundleFilePath,
		ConfiguratorFilePath: fmt.Sprintf("%s/%s.config.js", b.PageOutputDir, settings.BundleKey),
		ConfiguratorPage:     page,
	}, nil
}

func (b *WebPackBundler) Bundle(configuratorFilePath string) error {
	cmd := exec.Command("node", fsutils.NormalizePath(fmt.Sprintf("%s/.bin/webpack", b.NodeModulesDir)), "--config", configuratorFilePath)
	_, err := cmd.Output()

	if err != nil {
		b.Logger.Warn(fmt.Sprintf(`invalid pack: "node %s --config %s"`, fsutils.NormalizePath(fmt.Sprintf("%s/.bin/webpack", b.NodeModulesDir)), configuratorFilePath))
	}

	return err
}

func (b *WebPackBundler) NodeDependencies() map[string]string {
	return map[string]string{
		"@babel/core": "^7.11.1",
		"@babel/plugin-proposal-export-default-from": "^7.12.13",
		"@babel/polyfill":     "^7.12.1",
		"@babel/preset-env":   "^7.11.0",
		"@babel/preset-react": "^7.10.4",
		"babel-loader":        "^8.1.0",
		"css-loader":          "^4.2.2",
		"html-loader":         "^1.1.0",
		"html-webpack-plugin": "^4.3.0",
		"style-loader":        "^1.2.1",
		"webpack":             "^4.44.1",
		"webpack-cli":         "^3.3.12",
		"webpack-merge":       "^5.8.0",
	}
}
