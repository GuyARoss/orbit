package build

import (
	"os"

	"github.com/GuyARoss/orbit/pkg/fs"
	"github.com/GuyARoss/orbit/pkg/libgen"
)

type AutoGenPages struct {
	BundleData *libgen.LibOut
	Master     *libgen.LibOut
}

func (s *AutoGenPages) CreateAndOverwrite() {
	s.BundleData.WriteFile("example/orbit.go")
}

type GenPagesSettings struct {
	PackageName string
}

func (s *GenPagesSettings) SetupAutoGenPages(pages []*fs.PackedPage) *AutoGenPages {
	lg := &libgen.LibOut{
		PackageName:   "orbit",
		BaseBundleOut: ".orbit/base/pages",
	}

	for _, p := range pages {
		lg.ApplyPage(p.PageName, p.BundleKey)
	}

	return &AutoGenPages{
		BundleData: lg,
	}
}

func (s *GenPagesSettings) CleanPathing() error {
	err := os.RemoveAll(".orbit/")
	if err != nil {
		return err
	}

	return nil
}
