// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/GuyARoss/orbit/internal"
	"github.com/GuyARoss/orbit/internal/srcpack"
	"github.com/GuyARoss/orbit/pkg/experiments"
	"github.com/GuyARoss/orbit/pkg/fsutils"
	"github.com/GuyARoss/orbit/pkg/htmlparse"
	"github.com/GuyARoss/orbit/pkg/jsparse"
	"github.com/GuyARoss/orbit/pkg/log"
	"github.com/GuyARoss/orbit/pkg/webwrap"
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
			spaEntryComponent := components[0]

			// this bundle component should get copied from the .orbit/dist to the spa_output
			bundlePath := fmt.Sprintf(".orbit/dist/%s.js", spaEntryComponent.BundleKey())
			outDir := viper.GetString("spa_out_dir")

			if _, err := os.Stat(outDir); os.IsNotExist(err) {
				if err := os.Mkdir(outDir, os.ModePerm); err != nil {
					logger.Error(fmt.Sprintf("cannot setup spa out directory %s, exiting", err))
					return
				}
			}

			if err = fsutils.CopyFile(bundlePath, fmt.Sprintf("%s/%s.js", outDir, spaEntryComponent.BundleKey())); err != nil {
				logger.Error(fmt.Sprintf("cannot copy file %s, exiting", err))
				return
			}

			if viper.GetString("public_path") != "" {
				// parse the html page (if exists) and add the javascript to it.
				htmlDoc := htmlparse.DocFromFile(viper.GetString("public_path"))
				wr := spaEntryComponent.WebWrapper()
				if wr == nil {
					return
				}

				body := wr.RequiredBodyDOMElements(context.TODO(), &webwrap.CacheDOMOpts{
					WebPrefix: outDir,
					CacheDir:  "",
				})
				// note: altering the order of the appends will break functionality
				htmlDoc.Body = append(htmlDoc.Body, wr.DocumentTag(spaEntryComponent.BundleKey()))

				htmlDoc.Body = append(htmlDoc.Body, body...)
				htmlDoc.Body = append(htmlDoc.Body, fmt.Sprintf(`<script src="./%s.js"></script>`, spaEntryComponent.BundleKey()))

				htmlDoc.SaveToFile(fmt.Sprintf("%s/index.html", outDir))
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
