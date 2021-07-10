package build

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CMD = &cobra.Command{
	Use: "build",
	Run: func(cmd *cobra.Command, args []string) {
		as := &GenPagesSettings{
			PackageName: viper.GetString("pacname"),
			OutDir:      viper.GetString("out"),
			WebDir:      viper.GetString("webdir"),
		}

		err := as.CleanPathing()
		if err != nil {
			log.Fatal(err)
		}

		as.ApplyPages()
	},
}

func init() {
	var outDir string
	var pacname string

	CMD.PersistentFlags().StringVar(&outDir, "out", "", "specifies the out directory of the generated code files")
	CMD.PersistentFlags().StringVar(&pacname, "pacname", "orbit", "specifies the package name of the generated code files")
	CMD.PersistentFlags().StringVar(&pacname, "webdir", "/", "specifies the directory of the web pages, leave blank for use of the root dir")

	CMD.MarkFlagRequired("out")

	viper.BindPFlag("out", CMD.PersistentFlags().Lookup("out"))
	viper.BindPFlag("pacname", CMD.PersistentFlags().Lookup("pacname"))
	viper.BindPFlag("webdir", CMD.PersistentFlags().Lookup("webdir"))
}
