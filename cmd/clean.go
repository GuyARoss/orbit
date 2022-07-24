package cmd

import (
	"github.com/GuyARoss/orbit/internal"
	"github.com/GuyARoss/orbit/pkg/log"
	"github.com/spf13/cobra"
)

var cleanCMD = &cobra.Command{
	Use:   "clean",
	Short: "clean",
	Run: func(cmd *cobra.Command, args []string) {
		err := internal.OrbitRemoveFileStructure()
		if err != nil {
			panic(err)
		}

		log.NewDefaultLogger().Info("paths cleaned")
	},
}
