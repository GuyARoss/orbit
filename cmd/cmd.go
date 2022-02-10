package cmd

import (
	"github.com/GuyARoss/orbit/pkg/log"
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
	var nodeModDir string
	var publicDir string
	var recordBuildDuation bool

	buildCmds := [2]*cobra.Command{
		buildCMD, devCMD,
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

		cmd.PersistentFlags().StringVar(&publicDir, "publicdir", "./public/index.html", "specifies the public directory for the base html webpage")
		viper.BindPFlag("publicdir", cmd.PersistentFlags().Lookup("publicdir"))

		cmd.PersistentFlags().StringVar(&pacname, "pacname", "orbit", "specifies the package name of the generated code files")
		viper.BindPFlag("pacname", cmd.PersistentFlags().Lookup("pacname"))

		cmd.PersistentFlags().StringVar(&nodeModDir, "nodemod", "./node_modules", "specifies the directory to find node modules")
		viper.BindPFlag("nodemod", cmd.PersistentFlags().Lookup("nodemod"))

		cmd.PersistentFlags().BoolVar(&recordBuildDuation, "debugduration", false, "capture the duration for the entire bundle")
		viper.BindPFlag("debugduration", cmd.PersistentFlags().Lookup("debugduration"))
	}
}

func Execute() {
	logger := log.NewDefaultLogger()
	logger.Title("orbit-ssr")

	if err := RootCMD.Execute(); err != nil {
		panic(err)
	}
}
