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
	var appDir string
	var outDir string
	var pacname string
	var nodeModDir string
	var publicPath string
	var dependout string
	var spaEntry string
	var spaOutDir string
	var experimentalFeatures []string

	buildCmds := [4]*cobra.Command{
		buildCMD, devCMD, initCMD, deployCMD,
	}

	for _, cmd := range buildCmds {
		cmd.PersistentFlags().StringVar(&appDir, "app_dir", "./", "specifies the directory where the application lives, left blank will use the root directory")
		viper.BindPFlag("app_dir", cmd.PersistentFlags().Lookup("app_dir"))

		cmd.PersistentFlags().StringVar(&outDir, "out_dir", "./", "specifies the out directory of the generated code files")
		viper.BindPFlag("out_dir", cmd.PersistentFlags().Lookup("out_dir"))

		cmd.PersistentFlags().StringVar(&publicPath, "public_path", "./public/index.html", "specifies the path for the base html webpage default './public/index.html'")
		viper.BindPFlag("public_path", cmd.PersistentFlags().Lookup("public_path"))

		cmd.PersistentFlags().StringVar(&pacname, "package_name", "orbit", "specifies the package name of the generated code files")
		viper.BindPFlag("package_name", cmd.PersistentFlags().Lookup("package_name"))

		cmd.PersistentFlags().StringVar(&nodeModDir, "node_modules_dir", "./node_modules", "specifies the directory to find node modules")
		viper.BindPFlag("node_modules_dir", cmd.PersistentFlags().Lookup("node_modules_dir"))

		cmd.PersistentFlags().StringVar(&dependout, "dep_map_out_dir", "", "specifies the directory to output a dependency map")
		viper.BindPFlag("dep_map_out_dir", cmd.PersistentFlags().Lookup("dep_map_out_dir"))

		cmd.PersistentFlags().StringSliceVar(&experimentalFeatures, "experimental", []string{}, "comma delimited list of experimental features to turn on, to view experiemental features use the command 'experimental'")
		viper.BindPFlag("experimental", cmd.PersistentFlags().Lookup("experimental"))

		cmd.PersistentFlags().StringVar(&spaEntry, "spa_entry_path", "", "when specified this entry should be the file name of the entrypoint file. note this command will force the application to become an SPA")
		viper.BindPFlag("spa_entry_path", cmd.PersistentFlags().Lookup("spa_entry_path"))

		cmd.PersistentFlags().StringVar(&spaOutDir, "spa_out_dir", "./dist", "output directory to write an SPA, requires 'spa_entry_path' to be set")
		viper.BindPFlag("spa_out_dir", cmd.PersistentFlags().Lookup("spa_out_dir"))
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
