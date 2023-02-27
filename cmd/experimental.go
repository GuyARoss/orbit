// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var experimentalCMD = &cobra.Command{
	Use:   "experimental",
	Long:  "shows the list of available experimental options",
	Short: "list of experimental options",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`
Experimental Options:

ssr 	enables the usage of ssr functionality for available web wrappers 
swc 	enables the usage of the swc compiler in place of babel
	
		`)
	},
}
