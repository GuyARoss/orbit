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
	page.Imports = append(page.Imports, "const {merge} = require('webpack-merge')")
	page.Imports = append(page.Imports, "const baseConfig = require('../../assets/base.config.js')")

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

	return err
}
