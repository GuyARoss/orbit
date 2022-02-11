package cmd

import (
	"context"
	"fmt"

	"github.com/GuyARoss/orbit/internal"
	"github.com/GuyARoss/orbit/internal/srcpack"
	"github.com/GuyARoss/orbit/pkg/log"
	"github.com/GuyARoss/orbit/pkg/runtimeanalytics"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var buildCMD = &cobra.Command{
	Use:   "build",
	Long:  "bundle data given the specified pages in prod mode",
	Short: "bundle data given the specified pages in prod mode",
	Run: func(cmd *cobra.Command, args []string) {
		settings := &internal.GenPagesSettings{
			PackageName:    viper.GetString("pacname"),
			OutDir:         viper.GetString("out"),
			WebDir:         viper.GetString("webdir"),
			BundlerMode:    viper.GetString("mode"),
			NodeModulePath: viper.GetString("nodemod"),
			PublicDir:      viper.GetString("publicdir"),
		}

		analytics := &runtimeanalytics.RuntimeAnalytics{}

		if viper.GetBool("debugduration") {
			analytics.StartCapture()
		}

		err := settings.CleanPathing()
		if err != nil {
			panic(err)
		}

		pages, err := settings.PackWebDir(context.Background(), srcpack.NewSyncHook(log.NewDefaultLogger()))
		if err != nil {
			panic(err)
		}

		err = pages.WriteOut()
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
