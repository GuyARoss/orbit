package internal

import (
	"fmt"
	"os"

	"github.com/GuyARoss/orbit/pkg/fs"
	"github.com/GuyARoss/orbit/pkg/libgen"
)

type AutoGenPages struct {
	BundleData *libgen.LibOut
	Master     *libgen.LibOut

	Pages  []*fs.PackedPage
	OutDir string
}

type GenPagesSettings struct {
	PackageName    string
	OutDir         string
	WebDir         string
	BundlerMode    string
	AssetDir       string
	NodeModulePath string
}

func (s *GenPagesSettings) PackWebDir() *AutoGenPages {
	settings := &fs.PackSettings{
		BundlerSettings: &fs.BundlerSettings{
			Mode:           fs.BundlerMode(s.BundlerMode),
			NodeModulePath: s.NodeModulePath,
			WebDir:         s.WebDir,
		},
		AssetDir: s.AssetDir,
	}

	pages := settings.Pack(s.WebDir, ".orbit/base/pages")

	lg := &libgen.BundleGroup{
		PackageName:   s.PackageName,
		BaseBundleOut: ".orbit/dist",
		BundleMode:    string(settings.BundlerSettings.Mode),
	}

	for _, p := range *pages {
		lg.ApplyBundle(p.PageName, p.BundleKey)
	}

	libStaticContent, parseErr := libgen.ParseStaticFile(".orbit/assets/orbit.go")
	if parseErr != nil {
		panic(parseErr)
	}

	return &AutoGenPages{
		OutDir:     s.OutDir,
		BundleData: lg.CreateBundleLib(),
		Master: &libgen.LibOut{
			Body:        libStaticContent,
			PackageName: s.PackageName,
		},
		Pages: *pages,
	}
}

func (s *GenPagesSettings) Repack(p *fs.PackedPage) {
	lg := &libgen.BundleGroup{
		PackageName:   s.PackageName,
		BaseBundleOut: ".orbit/dist",
	}
	lg.ApplyBundle(p.PageName, p.BundleKey)
}

func (s *AutoGenPages) WriteOut() error {
	err := s.BundleData.WriteFile(fmt.Sprintf("%s/autogen_bundle.go", s.OutDir))
	if err != nil {
		return err
	}
	err = s.Master.WriteFile(fmt.Sprintf("%s/autogen_master.go", s.OutDir))
	if err != nil {
		return err
	}

	return nil
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
