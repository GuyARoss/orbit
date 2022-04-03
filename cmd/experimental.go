package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var experiementalCMD = &cobra.Command{
	Use:   "experimental",
	Long:  "shows the list of available expiremental options",
	Short: "list of experimental options",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`
Experimental Options:

ssr 	enables the usage of ssr functionality for available web wrappers 
	
		`)
	},
}
