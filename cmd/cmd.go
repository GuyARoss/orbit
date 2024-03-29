// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package cmd

import (
	"github.com/GuyARoss/orbit/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCMD = &cobra.Command{
	Use:   "orbit-cli",
	Short: "Orbit Golang SSR CLI",
}

func init() {
	var webdir string
	var outDir string
	var pacname string
	var nodeModDir string
	var publicDir string
	var dependout string
	var experimentalFeatures []string

	buildCmds := [4]*cobra.Command{
		buildCMD, devCMD, initCMD, deployCMD,
	}

	for _, cmd := range buildCmds {
		cmd.PersistentFlags().StringVar(&webdir, "webdir", "./", "specifies the directory of the web pages, leave blank for use of the root dir")
		viper.BindPFlag("webdir", cmd.PersistentFlags().Lookup("webdir"))

		cmd.PersistentFlags().StringVar(&outDir, "out", "./", "specifies the out directory of the generated code files")
		viper.BindPFlag("out", cmd.PersistentFlags().Lookup("out"))

		cmd.PersistentFlags().StringVar(&publicDir, "publicdir", "./public/index.html", "specifies the public directory for the base html webpage")
		viper.BindPFlag("publicdir", cmd.PersistentFlags().Lookup("publicdir"))

		cmd.PersistentFlags().StringVar(&pacname, "pacname", "orbit", "specifies the package name of the generated code files")
		viper.BindPFlag("pacname", cmd.PersistentFlags().Lookup("pacname"))

		cmd.PersistentFlags().StringVar(&nodeModDir, "nodemod", "./node_modules", "specifies the directory to find node modules")
		viper.BindPFlag("nodemod", cmd.PersistentFlags().Lookup("nodemod"))

		cmd.PersistentFlags().StringVar(&dependout, "depout", "", "specifies the directory to output a dependency map")
		viper.BindPFlag("depout", cmd.PersistentFlags().Lookup("depout"))

		cmd.PersistentFlags().StringSliceVar(&experimentalFeatures, "experimental", []string{}, "comma delimated list of experimental features to turn on, to view experiemental features use the command 'experimental'")
		viper.BindPFlag("experimental", cmd.PersistentFlags().Lookup("experimental"))
	}
}

func Execute() {
	logger := log.NewDefaultLogger()
	logger.Clear()
	logger.Title("orbit-ssr")

	RootCMD.AddCommand(versionCMD)
	RootCMD.AddCommand(devCMD)
	RootCMD.AddCommand(buildCMD)
	RootCMD.AddCommand(toolCMD)
	RootCMD.AddCommand(experimentalCMD)
	RootCMD.AddCommand(deployCMD)
	RootCMD.AddCommand(cleanCMD)
	RootCMD.AddCommand(initCMD)

	if err := RootCMD.Execute(); err != nil {
		panic(err)
	}
}
