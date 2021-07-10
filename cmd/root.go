package cmd

import (
	"github.com/GuyARoss/orbit/cmd/build"
	"github.com/spf13/cobra"
)

var RootCMD = &cobra.Command{
	Use:   "orbit-cli",
	Short: "Orbit Golang SSR CLI",
}

func init() {
	RootCMD.AddCommand(build.CMD)
}

func Execute() {
	if err := RootCMD.Execute(); err != nil {
		panic(err)
	}
}
