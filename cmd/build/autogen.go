package build

import (
	"fmt"
	"os"

	"github.com/GuyARoss/orbit/pkg/fs"
	"github.com/GuyARoss/orbit/pkg/libgen"
)

type AutoGenPages struct {
	BundleData *libgen.LibOut
	Master     *libgen.LibOut
}

type GenPagesSettings struct {
	PackageName string
	OutDir      string
	WebDir      string
}

func (s *GenPagesSettings) SetupAutoGenPages() *AutoGenPages {
	pages := fs.Pack(s.WebDir, ".orbit/base/pages")

	lg := &libgen.BundleGroup{
		PackageName:   s.PackageName,
		BaseBundleOut: ".orbit/base/dist",
	}

	for _, p := range pages {
		lg.ApplyBundle(p.PageName, p.BundleKey)
	}

	libStaticContent, parseErr := libgen.ParseStaticFile(".orbit/base/assets/orbit.go")
	if parseErr != nil {
		panic(parseErr)
	}

	return &AutoGenPages{
		BundleData: lg.CreatePage(),
		Master: &libgen.LibOut{
			Body:        libStaticContent,
			PackageName: s.PackageName,
		},
	}
}

func (s *GenPagesSettings) ApplyPages() {
	pages := s.SetupAutoGenPages()

	pages.BundleData.WriteFile(fmt.Sprintf("%s/autogen_bundle.go", s.OutDir))
	pages.Master.WriteFile(fmt.Sprintf("%s/autogen_master.go", s.OutDir))
}

func (s *GenPagesSettings) CleanPathing() error {
	err := os.RemoveAll(".orbit/")
	if err != nil {
		return err
	}

	if !fs.DoesDirExist(s.OutDir) {
		err := os.Mkdir(s.OutDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// @@todo(debug) return err
	fs.SetupDirs()

	return nil
}
