// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package cmd

import (
	"github.com/GuyARoss/orbit/internal"
	"github.com/GuyARoss/orbit/pkg/experiments"
	"github.com/GuyARoss/orbit/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deployCMD = &cobra.Command{
	Use:   "deploy",
	Short: "deploy",
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.NewDefaultLogger()

		// this entire "deploy" cmd is based off of an experimental feature
		err := experiments.Load(logger, []string{"ssr"})

		if err != nil {
			logger.Warn(err.Error())
		}

		buildOpts := internal.NewBuildOptsFromViper()
		buildOpts.RequiredDirs = []string{
			viper.GetString("static_out_dir"),
		}

		components, err := internal.Build(buildOpts)
		if err != nil {
			panic(err)
		}

		staticBuild := internal.NewStaticBuild(buildOpts, viper.GetString("static_out_dir"))
		err = staticBuild.Build(components)

		if err != nil {
			logger.Error(err.Error())
		}
	},
}

func init() {
	var staticOut string

	deployCMD.PersistentFlags().StringVar(&staticOut, "static_out_dir", "./static", "path for the static file directory")
	viper.BindPFlag("staticout", deployCMD.PersistentFlags().Lookup("static_out_dir"))
}
