// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package cmd

import (
	"github.com/GuyARoss/orbit/internal"
	"github.com/GuyARoss/orbit/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deployCMD = &cobra.Command{
	Use:   "deploy",
	Short: "deploy",
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.NewDefaultLogger()

		buildOpts := &internal.BuildOpts{
			Packname:       viper.GetString("pacname"),
			OutDir:         viper.GetString("out"),
			WebDir:         viper.GetString("webdir"),
			Mode:           viper.GetString("deploy_bundle_mode"),
			NodeModulePath: viper.GetString("nodemod"),
			PublicDir:      viper.GetString("publicdir"),
			NoWrite:        true,
			Dirs: []string{
				viper.GetString("staticout"),
			},
		}

		components, err := internal.Build(buildOpts)
		if err != nil {
			panic(err)
		}

		staticBuild := internal.NewStaticBuild(buildOpts, viper.GetString("staticout"))
		err = staticBuild.Build(components)

		if err != nil {
			logger.Error(err.Error())
		}
	},
}

func init() {
	var staticOut string
	var mode string

	deployCMD.PersistentFlags().StringVar(&staticOut, "staticout", "./static", "path for the static file directory")
	viper.BindPFlag("staticout", deployCMD.PersistentFlags().Lookup("staticout"))

	deployCMD.PersistentFlags().StringVar(&mode, "mode", "production", "specifies the underlying bundler mode to run in")
	viper.BindPFlag("deploy_bundle_mode", deployCMD.PersistentFlags().Lookup("mode"))
}
