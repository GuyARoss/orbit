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
