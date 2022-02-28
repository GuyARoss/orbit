package cmd

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/GuyARoss/orbit/internal"
	"github.com/GuyARoss/orbit/internal/assets"
	"github.com/GuyARoss/orbit/internal/libout"
	"github.com/GuyARoss/orbit/internal/srcpack"
	"github.com/GuyARoss/orbit/pkg/bundler"
	"github.com/GuyARoss/orbit/pkg/fsutils"
	"github.com/GuyARoss/orbit/pkg/log"
	"github.com/GuyARoss/orbit/pkg/runtimeanalytics"
	webwrapper "github.com/GuyARoss/orbit/pkg/web_wrapper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var buildCMD = &cobra.Command{
	Use:   "build",
	Long:  "bundle data given the specified pages in prod mode",
	Short: "bundle data given the specified pages in prod mode",
	Run: func(cmd *cobra.Command, args []string) {
		analytics := &runtimeanalytics.RuntimeAnalytics{}

		if viper.GetBool("debugduration") {
			analytics.StartCapture()
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
			panic(err)
		}

		pageFiles := fsutils.DirFiles(fsutils.NormalizePath(fmt.Sprintf("%s/pages", viper.GetString("webdir"))))

		packer := srcpack.NewDefaultPacker(log.NewDefaultLogger(), &srcpack.DefaultPackerOpts{
			WebDir:        viper.GetString("webdir"),
			BundlerMode:   viper.GetString("mode"),
			NodeModuleDir: viper.GetString("nodemod"),
		})
		components, err := packer.PackMany(pageFiles)
		if err != nil {
			panic(err)
		}

		bg := libout.New(&libout.BundleGroupOpts{
			PackageName:   viper.GetString("pacname"),
			BaseBundleOut: fsutils.NormalizePath(".orbit/dist"),
			BundleMode:    string(viper.GetString("mode")),
			PublicDir:     viper.GetString("publicdir"),
		})

		ctx := context.Background()
		ctx = context.WithValue(ctx, bundler.BundlerID, viper.GetString("mode"))

		bg.AcceptComponents(ctx, components, &webwrapper.CacheDOMOpts{
			CacheDir:  fsutils.NormalizePath(".orbit/dist"),
			WebPrefix: fsutils.NormalizePath("/p/"),
		})

		err = bg.WriteLibout(libout.NewGOLibout(
			ats.AssetKey(assets.Tests),
			ats.AssetKey(assets.PrimaryPackage),
		), &libout.FilePathOpts{
			TestFile: fsutils.NormalizePath(fmt.Sprintf("%s/%s/orb_test.go", viper.GetString("webdir"), viper.GetString("pacname"))),
			EnvFile:  fsutils.NormalizePath(fmt.Sprintf("%s/%s/orb_env.go", viper.GetString("webdir"), viper.GetString("pacname"))),
			HTTPFile: fsutils.NormalizePath(fmt.Sprintf("%s/%s/orb_http.go", viper.GetString("webdir"), viper.GetString("pacname"))),
		})

		if err != nil {
			panic(err)
		}

		if viper.GetBool("debugduration") {
			end := analytics.StopCapture()
			fmt.Printf("total build duration: %fms\n", end)
		}
	},
}

func init() {
	RootCMD.AddCommand(buildCMD)
}
