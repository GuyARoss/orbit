package cmd

import "github.com/spf13/cobra"

var toolsCMD = &cobra.Command{
	Use:   "tool",
	Long:  "orbit suportted tooling",
	Short: "orbit supported tooling",
}

func init() {
	toolsCMD.AddCommand(dependgraphToolsCMD)
}
