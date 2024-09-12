// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package cmd

import (
	"fmt"

	"github.com/GuyARoss/orbit/internal"
	"github.com/GuyARoss/orbit/internal/srcpack"
	"github.com/GuyARoss/orbit/pkg/experiments"
	"github.com/GuyARoss/orbit/pkg/jsparse"
	"github.com/GuyARoss/orbit/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var buildCMD = &cobra.Command{
	Use:   "build",
	Long:  "bundle data given the specified pages in prod mode",
	Short: "bundle data given the specified pages in prod mode",
	Run: func(cmd *cobra.Command, args []string) {
		logger := log.NewDefaultLogger()

		err := experiments.Load(logger, viper.GetStringSlice("experimental"))
		if err != nil {
			logger.Warn(err.Error())
		}

		buildOpts := internal.NewBuildOptsFromViper()
		if buildOpts.Mode != "production" {
			logger.Warn(fmt.Sprintf("bundling mode '%s'\n", viper.GetString("build_bundle_mode")))
		}

		components, err := internal.Build(buildOpts)
		if len(components) < 1 {
			logger.Warn("no components were found, exiting")
			return
		}

		if viper.GetString("spa_entry_path") != "" {
			// if we are using the spa settings, there will only be a single bundle component
			if err := internal.BuildSPA(components[0], &internal.SPABuildOpts{
				PublicHTMLPath: viper.GetString("public_path"),
				SpaOutDir:      viper.GetString("spa_out_dir"),
			}); err != nil {
				logger.Error(err.Error())
				return
			}
		}

		if err != nil {
			logger.Error(err.Error())
			return
		}

		if viper.GetString("audit_path") != "" {
			components.Write(viper.GetString("audit_path"))
		}

		if viper.GetString("dep_map_out_dir") != "" {
			sourceMap, err := srcpack.New(viper.GetString("app_dir"), components, &srcpack.NewSourceMapOpts{
				WebDirPath: buildOpts.ApplicationDir,
				Parser:     &jsparse.JSFileParser{},
			})
			if err != nil {
				panic(err)
			}

			err = sourceMap.Write(viper.GetString("dep_map_out_dir"))
			if err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	var pageaudit string
	var mode string

	buildCMD.PersistentFlags().StringVar(&pageaudit, "audit_path", "", "file path used to output an audit file for the pages")
	viper.BindPFlag("audit_path", buildCMD.PersistentFlags().Lookup("audit_path"))

	buildCMD.PersistentFlags().StringVar(&mode, "mode", "production", "specifies the underlying bundler mode to run in")
	viper.BindPFlag("mode", buildCMD.PersistentFlags().Lookup("mode"))
}
