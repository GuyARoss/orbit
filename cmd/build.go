// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package cmd

import (
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

		err := experiments.LoadSingleton(logger, viper.GetStringSlice("experimental"))
		if err != nil {
			logger.Warn(err.Error())
		}

		components, err := internal.Build(&internal.BuildOpts{
			Packname:       viper.GetString("pacname"),
			OutDir:         viper.GetString("out"),
			WebDir:         viper.GetString("webdir"),
			Mode:           viper.GetString("mode"),
			NodeModulePath: viper.GetString("nodemod"),
			PublicDir:      viper.GetString("publicdir"),
		})

		if err != nil {
			logger.Error(err.Error())
			return
		}

		if viper.GetString("auditpage") != "" {
			components.Write(viper.GetString("auditpage"))
		}

		if viper.GetString("depout") != "" {
			sourceMap, err := srcpack.New(viper.GetString("webdir"), components, &srcpack.NewSourceMapOpts{
				WebDirPath: viper.GetString("webdir"),
				Parser:     &jsparse.JSFileParser{},
			})
			if err != nil {
				panic(err)
			}

			err = sourceMap.Write(viper.GetString("depout"))
			if err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	var pageaudit string

	buildCMD.PersistentFlags().StringVar(&pageaudit, "auditpage", "", "file path used to output an audit file for the pages")
	viper.BindPFlag("auditpage", buildCMD.PersistentFlags().Lookup("auditpage"))
}
