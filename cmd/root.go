package cmd

import (
	"github.com/GuyARoss/orbit/cmd/build"
	"github.com/GuyARoss/orbit/cmd/dev"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCMD = &cobra.Command{
	Use:   "orbit-cli",
	Short: "Orbit Golang SSR CLI",
}

func init() {
	var mode string
	var webdir string
	var outDir string
	var pacname string

	buildCmds := [2]*cobra.Command{
		build.CMD, dev.CMD,
	}

	for _, cmd := range buildCmds {
		if cmd.Use == "build" {
			cmd.PersistentFlags().StringVar(&mode, "mode", "production", "specifies the underlying bundler mode to run in")
		} else {
			cmd.PersistentFlags().StringVar(&mode, "mode", "development", "specifies the underlying bundler mode to run in")
		}
		viper.BindPFlag("mode", cmd.PersistentFlags().Lookup("mode"))

		cmd.PersistentFlags().StringVar(&webdir, "webdir", "./", "specifies the directory of the web pages, leave blank for use of the root dir")
		viper.BindPFlag("webdir", cmd.PersistentFlags().Lookup("webdir"))

		cmd.PersistentFlags().StringVar(&outDir, "out", "./orbit", "specifies the out directory of the generated code files")
		viper.BindPFlag("out", cmd.PersistentFlags().Lookup("out"))

		cmd.PersistentFlags().StringVar(&pacname, "pacname", "orbit", "specifies the package name of the generated code files")
		viper.BindPFlag("pacname", cmd.PersistentFlags().Lookup("pacname"))
	}

	RootCMD.AddCommand(build.CMD)
	RootCMD.AddCommand(dev.CMD)
}

func Execute() {
	if err := RootCMD.Execute(); err != nil {
		panic(err)
	}
}
