package main

import (
	"os"

	"github.com/GuyARoss/orbit/fs"
	"github.com/GuyARoss/orbit/libgen"
)

// @@todo: verify assets exist correctly
func main() {
	err := os.RemoveAll(".orbit/")
	if err != nil {
		panic(err)
	}

	fs.SetupDirs()
	pages := fs.Pack("example", ".orbit/base/pages")

	lg := &libgen.LibOut{
		PackageName:   "orbit",
		BaseBundleOut: ".orbit/base/pages",
	}

	for _, p := range pages {
		lg.ApplyPage(p.PageName, p.BundleKey)
	}

	lg.WriteFile("example/orbit.go")
}
