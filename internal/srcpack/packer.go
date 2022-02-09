package srcpack

import (
	"sync"
	"time"

	"github.com/GuyARoss/orbit/pkg/bundler"
	"github.com/GuyARoss/orbit/pkg/fs"
	"github.com/GuyARoss/orbit/pkg/jsparse"
	webwrapper "github.com/GuyARoss/orbit/pkg/web_wrapper"
)

type Packer struct {
	Bundler    bundler.Bundler
	WebWrapper webwrapper.WebWrapper

	AssetDir string
	WebDir   string
	JsParser jsparse.JSParser
}

func (s *Packer) CopyAssets() ([]*fs.CopyResults, error) {
	results := fs.CopyDir(s.AssetDir, s.AssetDir, ".orbit/assets", false)

	return results, nil
}

// PackSingle Packs the a single file paths into the orbit root directory
// Process includes:
// - Wrapping the component with the specified front-end web framework.
// - Bundling the component with the specified javascript bundler.
func (s *Packer) PackSingle(pageFilePath string) (*Component, error) {
	startTime := time.Now()

	page, err := s.JsParser.Parse(pageFilePath, s.WebDir)
	if err != nil {
		return nil, err
	}

	page = s.WebWrapper.Apply(page, pageFilePath)

	resource, err := s.Bundler.Setup(&bundler.BundleSetupSettings{
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
		Packer:              s,
		PageName:            page.Name(),
		BundleKey:           page.Key(),
		PackDurationSeconds: time.Since(startTime).Seconds(),
		dependencies:        page.Imports(),
		originalFilePath:    pageFilePath,
		m:                   &sync.Mutex{},
	}, bundleErr
}

type concPack struct {
	m        *sync.Mutex
	settings *Packer

	packedPages []*Component
	packMap     map[string]bool
}

func (p *concPack) PackSingle(wg *sync.WaitGroup, dir string, hooks Hooks) {
	page, err := p.settings.PackSingle(dir)
	if p.packMap[page.PageName] {
		return
	}

	if err != nil {
		// @@report error with packing via channel
		return
	}

	p.m.Lock()
	p.packedPages = append(p.packedPages, page)
	p.packMap[page.PageName] = true

	if hooks != nil {
		hooks.Pre(dir)
		hooks.Post(page.PackDurationSeconds)
	}
	p.m.Unlock()

	wg.Done()
}

// PackPages
// Packs the provoided file paths into the orbit root directory
// Process includes:
// - Wrapping the component with the specified front-end web framework.
// - Bundling the component with the specified javascript bundler.
func (s *Packer) PackMany(pages []string, hooks Hooks) ([]*Component, error) {
	cp := &concPack{
		m:           &sync.Mutex{},
		settings:    s,
		packedPages: make([]*Component, 0),
		packMap:     make(map[string]bool),
	}

	wg := &sync.WaitGroup{}
	for _, dir := range pages {
		wg.Add(1)
		go cp.PackSingle(wg, dir, hooks)
	}

	wg.Wait()

	return cp.packedPages, nil
}
