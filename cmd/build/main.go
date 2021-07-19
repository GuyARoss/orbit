package build

import (
	"log"

	"github.com/GuyARoss/orbit/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CMD = &cobra.Command{
	Use: "build",
	Run: func(cmd *cobra.Command, args []string) {
		as := &internal.GenPagesSettings{
			PackageName: viper.GetString("pacname"),
			OutDir:      viper.GetString("out"),
			WebDir:      viper.GetString("webdir"),
			BundlerMode: viper.GetString("mode"),
		}

		execute(as)
	},
}

func execute(settings *internal.GenPagesSettings) {
	err := settings.CleanPathing()
	if err != nil {
		log.Fatal(err)
	}

	settings.ApplyPages()
}
