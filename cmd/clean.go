package cmd

import (
	"io/fs"

	"github.com/GuyARoss/orbit/internal"
	"github.com/GuyARoss/orbit/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cleanCMD = &cobra.Command{
	Use:   "clean",
	Short: "clean",
	Run: func(cmd *cobra.Command, args []string) {
		err := (&internal.FileStructure{
			PackageName: viper.GetString("pacname"),
			OutDir:      viper.GetString("out"),
			Assets:      []fs.DirEntry{},
		}).Cleanup()

		if err != nil {
			panic(err)
		}

		log.NewDefaultLogger().Info("paths cleaned")
	},
}
