// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package cmd

import (
	"github.com/GuyARoss/orbit/internal"
	"github.com/GuyARoss/orbit/pkg/fsutils"
	webwrap "github.com/GuyARoss/orbit/pkg/webwrap/embed"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deployCMD = &cobra.Command{
	Use:   "deploy",
	Short: "deploy",
	Run: func(cmd *cobra.Command, args []string) {
		components, err := internal.Build(&internal.BuildOpts{
			Packname:       viper.GetString("pacname"),
			OutDir:         viper.GetString("out"),
			WebDir:         viper.GetString("webdir"),
			Mode:           viper.GetString("mode"),
			NodeModulePath: viper.GetString("nodemod"),
			PublicDir:      viper.GetString("publicdir"),
			NoWrite:        true,
			Dirs: []string{
				viper.GetString("staticout"),
			},
		})
		if err != nil {
			panic(err)
		}

		staticMap := make(map[webwrap.PageRender]bool)
		pages := make(map[webwrap.PageRender]*webwrap.DocumentRenderer)
		bundleToPath := make(map[webwrap.PageRender]string)

		for _, c := range components {
			pages[webwrap.PageRender(c.BundleKey())] = webwrap.NewEmptyDocumentRenderer(c.WebWrapper().Version())
			bundleToPath[webwrap.PageRender(c.BundleKey())] = fsutils.LastPathIndex(c.OriginalFilePath()) + ".html"

			if c.IsStaticResource() {
				switch c.WebWrapper().Version() {
				case "reactSSR":
					staticMap[webwrap.PageRender(c.BundleKey())] = true
				default:
					// TODO: error not supported
					continue
				}
			}
		}

		if len(staticMap) == 0 {
			return
		}
		doc := webwrap.DocFromFile("./public/index.html")

		defer webwrap.Close()
		webwrap.StartupTaskReactSSR(viper.GetString("staticout"), pages, staticMap, bundleToPath, *doc)()
	},
}

func init() {
	var staticOut string

	deployCMD.PersistentFlags().StringVar(&staticOut, "staticout", "./static", "path for the static file directory")
	viper.BindPFlag("staticout", deployCMD.PersistentFlags().Lookup("staticout"))
}
