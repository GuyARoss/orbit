// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.
package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os/exec"

	"github.com/GuyARoss/orbit/internal"
	"github.com/GuyARoss/orbit/internal/assets"
	"github.com/GuyARoss/orbit/pkg/prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCMD = &cobra.Command{
	Use:  "init",
	Long: "initializes the project directory",
	Run: func(cmd *cobra.Command, args []string) {
		nodeDependencies := map[string]string{
			"@babel/core": "^7.11.1",
			"@babel/plugin-proposal-export-default-from": "^7.12.13",
			"@babel/polyfill":     "^7.12.1",
			"@babel/preset-env":   "^7.11.0",
			"@babel/preset-react": "^7.10.4",
			"babel-loader":        "^8.1.0",
			"css-loader":          "^4.2.2",
			"html-loader":         "^1.1.0",
			"html-webpack-plugin": "^4.3.0",
			"react":               "^16.13.1",
			"react-dom":           "^16.13.1",
			"react-hot-loader":    "^4.12.21",
			"react-router-dom":    "^5.2.0",
			"style-loader":        "^1.2.1",
			"webpack":             "^4.44.1",
			"webpack-cli":         "^3.3.12",
			"webpack-merge":       "^5.8.0",
		}

		pkgJson := &internal.PackageJSONTemplate{
			Name:         prompt.StringPrompt("Project Name: "),
			Version:      prompt.StringPrompt("Project Version: "),
			Description:  prompt.StringPrompt("Project Description: "),
			Author:       prompt.StringPrompt("Author: "),
			License:      prompt.StringPrompt("License: "),
			Dependencies: nodeDependencies,
		}

		err := pkgJson.Write(fmt.Sprintf("%s/package.json", viper.GetString("out")))
		if err != nil {
			log.Fatal(err)
		}

		execcmd := exec.Command("npm", "install")
		if err := execcmd.Run(); err != nil {
			log.Fatal(err)
		}

		ats, err := assets.AssetKeys()
		if err != nil {
			panic(err)
		}

		err = internal.OrbitFileStructure(&internal.FileStructureOpts{
			PackageName: viper.GetString("pacname"),
			OutDir:      viper.GetString("out"),
			Assets:      []fs.DirEntry{ats.AssetKey(assets.WebPackConfig)},
		})

		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCMD.AddCommand(initCMD)
}
