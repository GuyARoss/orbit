package srcpack

import (
	"context"
	"sync"

	"github.com/GuyARoss/orbit/pkg/bundler"
	"github.com/GuyARoss/orbit/pkg/fs"
	"github.com/GuyARoss/orbit/pkg/jsparse"
	webwrapper "github.com/GuyARoss/orbit/pkg/web_wrapper"
)

type Packer struct {
	Bundler    bundler.Bundler
	WebWrapper webwrapper.JSWebWrapper

	AssetDir string
	WebDir   string
	JsParser jsparse.JSParser
}

func (s *Packer) CopyAssets() ([]*fs.CopyResults, error) {
	results := fs.CopyDir(s.AssetDir, s.AssetDir, ".orbit/assets", false)

	return results, nil
}

// newComponent creates a new component given a page file path, this
// component will inherit the necessary parameters from the packer.
// this process involves the following:
// 1. wrapping the component with the specified front-end web framework.
// 2. bundling the component with the specified javascript bundler.
func (s *Packer) NewComponent(pageFilePath string) (*Component, error) {
	page, err := s.JsParser.Parse(pageFilePath, s.WebDir)
	if err != nil {
		return nil, err
	}

	page = s.WebWrapper.Apply(page, pageFilePath)

	resource, err := s.Bundler.Setup(context.TODO(), &bundler.BundleOpts{
		FileName:  pageFilePath,
		BundleKey: page.Key(),
	})
	if err != nil {
		return nil, err
	}

	bundlePageErr := page.WriteFile(resource.BundleFilePath)
	if bundlePageErr != nil {
		return nil, bundlePageErr
	}

	configErr := resource.ConfiguratorPage.WriteFile(resource.ConfiguratorFilePath)
	if configErr != nil {
		return nil, configErr
	}

	bundleErr := s.Bundler.Bundle(resource.ConfiguratorFilePath)
	return &Component{
		Name:             page.Name(),
		BundleKey:        page.Key(),
		dependencies:     page.Imports(),
		originalFilePath: pageFilePath,
		m:                &sync.Mutex{},
		WebWrapper:       s.WebWrapper,
		Bundler:          s.Bundler,
		JsParser:         s.JsParser,
	}, bundleErr
}

type concPack struct {
	m        *sync.Mutex
	settings *Packer

	packedPages []*Component
	packMap     map[string]bool
}

func (p *concPack) PackSingle(wg *sync.WaitGroup, dir string) {
	page, err := p.settings.NewComponent(dir)
	if p.packMap[page.Name] {
		return
	}

	if err != nil {
		// @@report error with packing via channel
		return
	}

	p.packedPages = append(p.packedPages, page)
	p.packMap[page.Name] = true

	wg.Done()
}

// packs the provoided file paths into the orbit root directory
// this process includes the following:
// 1. wrapping the component with the specified front-end web framework.
// 2. bundling the component with the specified javascript bundler.
func (s *Packer) PackMany(pages []string, hooks Hooks) ([]*Component, error) {
	cp := &concPack{
		settings:    s,
		packedPages: make([]*Component, 0),
		packMap:     make(map[string]bool),
	}

	wg := &sync.WaitGroup{}
	for _, dir := range pages {
		wg.Add(1)
		go cp.PackSingle(wg, dir)
	}

	wg.Wait()

	return cp.packedPages, nil
}
