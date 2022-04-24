package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCMD = &cobra.Command{
	Use:   "version",
	Short: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version 0.3.6")
	},
}
