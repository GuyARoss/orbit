// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/GuyARoss/orbit/internal"
	"github.com/GuyARoss/orbit/internal/assets"
	"github.com/GuyARoss/orbit/internal/libout"
	"github.com/GuyARoss/orbit/internal/srcpack"
	"github.com/GuyARoss/orbit/pkg/fsutils"
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
		ats, err := assets.AssetKeys()
		if err != nil {
			panic(err)
		}

		err = internal.OrbitFileStructure(&internal.FileStructureOpts{
			PackageName: viper.GetString("pacname"),
			OutDir:      viper.GetString("out"),
			Assets: []fs.DirEntry{
				ats.AssetEntry(assets.WebPackConfig),
				ats.AssetEntry(assets.SSRProtoFile),
				ats.AssetEntry(assets.JsWebPackConfig),
			},
		})

		if err != nil {
			panic(err)
		}

		pageFiles := fsutils.DirFiles(fmt.Sprintf("%s/pages", viper.GetString("webdir")))

		c, err := internal.CachedEnvFromFile(fmt.Sprintf("%s/%s/orb_env.go", viper.GetString("out"), viper.GetString("pacname")))
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}

		packer := srcpack.NewDefaultPacker(log.NewDefaultLogger(), &srcpack.DefaultPackerOpts{
			WebDir:           viper.GetString("webdir"),
			BundlerMode:      viper.GetString("mode"),
			NodeModuleDir:    viper.GetString("nodemod"),
			CachedBundleKeys: c,
		})

		components, err := packer.PackMany(pageFiles)
		if err != nil {
			panic(err)
		}

		bg := libout.New(&libout.BundleGroupOpts{
			PackageName:   viper.GetString("pacname"),
			BaseBundleOut: ".orbit/dist",
			BundleMode:    string(viper.GetString("mode")),
			PublicDir:     viper.GetString("publicdir"),
		})

		ctx := context.Background()
		ctx = context.WithValue(ctx, webwrap.BundlerID, viper.GetString("mode"))

		bg.AcceptComponents(ctx, components, &webwrap.CacheDOMOpts{
			CacheDir:  ".orbit/dist",
			WebPrefix: "/p/",
		})

		err = bg.WriteLibout(libout.NewGOLibout(
			ats.AssetKey(assets.Tests),
			ats.AssetKey(assets.PrimaryPackage),
		), &libout.FilePathOpts{
			TestFile: fmt.Sprintf("%s/%s/orb_test.go", viper.GetString("out"), viper.GetString("pacname")),
			EnvFile:  fmt.Sprintf("%s/%s/orb_env.go", viper.GetString("out"), viper.GetString("pacname")),
			HTTPFile: fmt.Sprintf("%s/%s/orb_http.go", viper.GetString("out"), viper.GetString("pacname")),
		})

		if err != nil {
			panic(err)
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
