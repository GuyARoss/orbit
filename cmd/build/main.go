package build

import (
	"github.com/GuyARoss/orbit/pkg/fs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CMD = &cobra.Command{
	Use: "build",
	Run: func(cmd *cobra.Command, args []string) {
		as := &GenPagesSettings{}

		err := as.CleanPathing()
		if err != nil {
			panic(err)
		}

		fs.SetupDirs()
		pages := fs.Pack("example", ".orbit/base/pages")
		preparedPages := as.SetupAutoGenPages(pages)

		preparedPages.CreateAndOverwrite()
	},
}

func init() {
	var outDir string
	var pacname string

	CMD.PersistentFlags().StringVar(&outDir, "out", "", "specifies the out directory of the generated code files")
	CMD.PersistentFlags().StringVar(&pacname, "pacname", "orbit", "specifies the package name of the generated code files")

	CMD.MarkFlagRequired("out")

	viper.BindPFlag("out", CMD.PersistentFlags().Lookup("out"))
	viper.BindPFlag("pacname", CMD.PersistentFlags().Lookup("pacname"))

}
