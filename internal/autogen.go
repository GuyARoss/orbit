package internal

import (
	"fmt"
	"os"

	"github.com/GuyARoss/orbit/pkg/bundler"
	"github.com/GuyARoss/orbit/pkg/fs"
	"github.com/GuyARoss/orbit/pkg/libgen"
	"github.com/GuyARoss/orbit/pkg/log"
	webwrapper "github.com/GuyARoss/orbit/pkg/web_wrapper"
)

type AutoGenPages struct {
	BundleData *libgen.LibOut
	Master     *libgen.LibOut

	// Pages  []*fs.PackedPage
	Pages []*PackedComponent

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

// func (s *GenPagesSettings) PackWebDir() *AutoGenPages {
// 	settings := &fs.PackSettings{
// 		BundlerSettings: &fs.BundlerSettings{
// 			Mode:           fs.BundlerMode(s.BundlerMode),
// 			NodeModulePath: s.NodeModulePath,
// 			WebDir:         s.WebDir,
// 		},
// 		AssetDir: s.AssetDir,
// 	}

// 	pages := settings.Pack(s.WebDir, ".orbit/base/pages")

// 	lg := &libgen.BundleGroup{
// 		PackageName:   s.PackageName,
// 		BaseBundleOut: ".orbit/dist",
// 		BundleMode:    string(settings.BundlerSettings.Mode),
// 	}

// 	for _, p := range *pages {
// 		lg.ApplyBundle(p.PageName, p.BundleKey)
// 	}

// 	libStaticContent, parseErr := libgen.ParseStaticFile(".orbit/assets/orbit.go")
// 	if parseErr != nil {
// 		panic(parseErr)
// 	}

// 	return &AutoGenPages{
// 		OutDir:     s.OutDir,
// 		BundleData: lg.CreateBundleLib(),
// 		Master: &libgen.LibOut{
// 			Body:        libStaticContent,
// 			PackageName: s.PackageName,
// 		},
// 		Pages: *pages,
// 	}
// }

func (s *GenPagesSettings) PackWebDir() *AutoGenPages {
	settings := PackSettings{
		Bundler: &bundler.WebPackBundler{
			BundleSettings: &bundler.BundleSettings{
				Mode:          bundler.BundlerMode(s.BundlerMode),
				WebDir:        s.WebDir,
				PageOutputDir: ".orbit/base/pages",
			},
			NodeModulesDir: s.NodeModulePath,
		},
		AssetDir: s.AssetDir,
		WebDir:   s.WebDir,
		WebWrapper: &webwrapper.ReactWebWrap{
			WebWrapSettings: &webwrapper.WebWrapSettings{
				WebDir: s.WebDir,
			},
		},
	}

	pageFiles := fs.DirFiles(fmt.Sprintf("%s/pages", s.WebDir))
	pages, err := settings.PackMany(pageFiles, nil)
	if err != nil {
		log.Error(err.Error())
	}

	lg := &libgen.BundleGroup{
		PackageName:   s.PackageName,
		BaseBundleOut: ".orbit/dist",
		BundleMode:    string(s.BundlerMode),
	}

	for _, p := range pages {
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
		Pages: pages,
	}
}

func (s *GenPagesSettings) Repack(p *PackedComponent) {
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
