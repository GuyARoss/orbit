package bundler

import (
	"fmt"
	"os/exec"

	"github.com/GuyARoss/orbit/pkg/jsparse"
)

type WebPackBundler struct {
	*BundleSettings

	NodeModulesDir string
}

func (b *WebPackBundler) Setup(settings *BundleSetupSettings) (*BundledResource, error) {
	page := &jsparse.Page{}
	page.Imports = append(page.Imports, &jsparse.ImportDependency{
		FinalStatement: "const {merge} = require('webpack-merge')",
		Type:           jsparse.ModuleImportType,
	})
	// @@todo(guy): this webpack config is currently based off of react, if we want to add support in the future
	// we will need to update this to apply a type context depending on which of the frontend frameworks are selected.
	// * we could also parse the file to determine which of the front-end frameworks are attached. then use the correct config *
	page.Imports = append(page.Imports, &jsparse.ImportDependency{
		FinalStatement: "const baseConfig = require('../../assets/base.config.js')",
		Type:           jsparse.ModuleImportType,
	})

	outputFileName := fmt.Sprintf("%s.js", settings.BundleKey)
	bundleFilePath := fmt.Sprintf("%s/%s.js", b.PageOutputDir, settings.BundleKey)

	page.Other = append(page.Other, fmt.Sprintf(`module.exports = merge(baseConfig, {
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
	cmd := exec.Command("node", fmt.Sprintf("%s/.bin/webpack", b.NodeModulesDir), "--config", configuratorFilePath)
	_, err := cmd.Output()

	if err != nil {
		fmt.Println("node", fmt.Sprintf("%s/.bin/webpack", b.NodeModulesDir), "--config", configuratorFilePath)
	}

	return err
}
