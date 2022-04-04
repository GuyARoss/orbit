// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package cmd

import (
	"github.com/GuyARoss/orbit/cmd/dependgraph"
	"github.com/spf13/cobra"
)

var toolCMD = &cobra.Command{
	Use:   "tool",
	Long:  "orbit suportted tooling",
	Short: "orbit supported tooling",
}

func init() {
	toolCMD.AddCommand(dependgraph.CMD)
}
