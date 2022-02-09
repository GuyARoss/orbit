package internal

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/GuyARoss/orbit/internal/assets"
	"github.com/GuyARoss/orbit/internal/srcpack"
	"github.com/GuyARoss/orbit/pkg/bundler"
	"github.com/GuyARoss/orbit/pkg/fs"
	"github.com/GuyARoss/orbit/pkg/jsparse"
	"github.com/GuyARoss/orbit/pkg/libgen"
)

type AutoGenPages struct {
	BundleData *libgen.LibOut
	Master     *libgen.LibOut
	Pages      []*srcpack.Component
	OutDir     string
}

type GenPagesSettings struct {
	PackageName    string
	OutDir         string
	WebDir         string
	BundlerMode    string
	NodeModulePath string
	PublicDir      string
}

func (s *GenPagesSettings) SetupPack(ctx context.Context) (context.Context, *srcpack.Packer) {
	ctx = context.WithValue(ctx, bundler.BundlerModeKey, bundler.BundlerMode(s.BundlerMode))

	return ctx, &srcpack.Packer{
		// Bundler: &bundler.WebPackBundler{
		// BundleSettings: &bundler.BundleSettings{
		// 	Mode:          bundler.BundlerMode(s.BundlerMode),
		// 	WebDir:        s.WebDir,
		// 	PageOutputDir: ".orbit/base/pages",
		// },
		// NodeModulesDir: s.NodeModulePath,
		// },
		WebDir:   s.WebDir,
		JsParser: &jsparse.JSFileParser{},
	}
}

func (s *GenPagesSettings) PackWebDir(ctx context.Context, hook srcpack.Hooks) (*AutoGenPages, error) {
	ctx, settings := s.SetupPack(ctx)

	err := assets.WriteAssetsDir(".orbit/assets")
	if err != nil {
		return nil, err
	}

	pageFiles := fs.DirFiles(fmt.Sprintf("%s/pages", s.WebDir))
	pages, err := settings.PackMany(pageFiles, hook)
	if err != nil {
		return nil, err
	}

	lg := &libgen.BundleGroup{
		PackageName:   s.PackageName,
		BaseBundleOut: ".orbit/dist",
		BundleMode:    string(s.BundlerMode),
		PublicDir:     s.PublicDir,
	}

	for _, p := range pages {
		lg.ApplyBundle(p.Name, p.BundleKey)
	}

	libStaticContent, parseErr := libgen.ParseStaticFile(".orbit/assets/orbit.go")
	if parseErr != nil {
		return nil, parseErr
	}

	return &AutoGenPages{
		OutDir:     s.OutDir,
		BundleData: lg.CreateBundleLib(),
		Master: &libgen.LibOut{
			Body:        libStaticContent,
			PackageName: s.PackageName,
		},
		Pages: pages,
	}, nil
}

func (s *GenPagesSettings) Repack(p *srcpack.Component) error {
	h := &srcpack.SyncHook{}
	h.Pre(p.OriginalFilePath())

	r := p.Repack()

	return r
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

	dirs := []string{".orbit", ".orbit/base", ".orbit/base/pages", ".orbit/dist", ".orbit/assets"}
	for _, dir := range dirs {
		_, err := os.Stat(dir)
		if errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(dir, 0755)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
