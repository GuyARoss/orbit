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

	lg := &libgen.LibOut{
		PackageName:   s.PackageName,
		BaseBundleOut: ".orbit/base/pages",
	}

	for _, p := range pages {
		lg.ApplyPage(p.PageName, p.BundleKey)
	}

	return &AutoGenPages{
		BundleData: lg,
	}
}

func (s *GenPagesSettings) ApplyPages() {
	pages := s.SetupAutoGenPages()

	pages.BundleData.WriteFile(fmt.Sprintf("%s/autogen_bundle.go", s.OutDir))
	// pages.Master.WriteFile(fmt.Sprintf("%s/autogen_master", s.OutDir))
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
