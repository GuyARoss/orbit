package cmd

import (
	"fmt"
	"log"

	"github.com/GuyARoss/orbit/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var buildCMD = &cobra.Command{
	Use: "build",
	Run: func(cmd *cobra.Command, args []string) {
		settings := &internal.GenPagesSettings{
			PackageName:    viper.GetString("pacname"),
			OutDir:         viper.GetString("out"),
			WebDir:         viper.GetString("webdir"),
			BundlerMode:    viper.GetString("mode"),
			AssetDir:       viper.GetString("assetdir"),
			NodeModulePath: viper.GetString("nodemod"),
		}

		fmt.Println(settings.AssetDir)

		err := settings.CleanPathing()
		if err != nil {
			log.Fatal(err)
		}

		settings.ApplyPages()
	},
}

func init() {
	RootCMD.AddCommand(buildCMD)
}
