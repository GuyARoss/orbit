package internal

import (
	"fmt"
	"sync"
	"time"

	"github.com/GuyARoss/orbit/pkg/bundler"
	"github.com/GuyARoss/orbit/pkg/fs"
	"github.com/GuyARoss/orbit/pkg/jsparse"
	"github.com/GuyARoss/orbit/pkg/log"
	webwrapper "github.com/GuyARoss/orbit/pkg/web_wrapper"
)

type PackSettings struct {
	Bundler    bundler.Bundler
	WebWrapper webwrapper.WebWrapper

	AssetDir string
	WebDir   string
	JsParser jsparse.JSParser
}

func (s *PackSettings) CopyAssets() ([]*fs.CopyResults, error) {
	results := fs.CopyDir(s.AssetDir, s.AssetDir, ".orbit/assets", false)

	return results, nil
}

// PackedComponent
// Web-Component that has been successfully ran, and output from a packing method.
type PackedComponent struct {
	PageName            string
	BundleKey           string
	PackDurationSeconds float64

	dependencies     []*jsparse.ImportDependency
	originalFilePath string
	settings         *PackSettings
	m                *sync.Mutex
}

func (s *PackedComponent) OriginalFilePath() string {
	return s.originalFilePath
}

func (s *PackedComponent) Dependencies() []*jsparse.ImportDependency {
	return s.dependencies
}

// PackSingle Packs the a single file paths into the orbit root directory
// Process includes:
// - Wrapping the component with the specified front-end web framework.
// - Bundling the component with the specified javascript bundler.
func (s *PackSettings) PackSingle(pageFilePath string) (*PackedComponent, error) {
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
	return &PackedComponent{
		PageName:            page.Name(),
		BundleKey:           page.Key(),
		PackDurationSeconds: time.Since(startTime).Seconds(),
		dependencies:        page.Imports(),
		originalFilePath:    pageFilePath,
		settings:            s,
		m:                   &sync.Mutex{},
	}, bundleErr
}

// PackHooks
// passing of "per" & "post" hooks for our iterative packing method "PackPages".
type PackHooks interface {
	Pre(filePath string)      // "pre" runs before each component packing iteration
	Post(elapsedTime float64) // "post" runs after each component packing iteration
}

type DefaultPackHook struct{}

func (s *DefaultPackHook) Pre(filePath string) {
	log.Info(fmt.Sprintf("bundling %s → ...", filePath))
}

func (s *DefaultPackHook) Post(elapsedTime float64) {
	log.Success(fmt.Sprintf("completed in %fs\n", elapsedTime))
}

func (s *PackedComponent) Repack() error {
	startTime := time.Now()

	// parse original page
	page, err := s.settings.JsParser.Parse(s.originalFilePath, s.settings.WebDir)
	if err != nil {
		return err
	}

	// apply the nessasacary requirements for the web framework to the original page
	page = s.settings.WebWrapper.Apply(page, s.originalFilePath)
	resource, err := s.settings.Bundler.Setup(&bundler.BundleSetupSettings{
		FileName:  s.originalFilePath,
		BundleKey: s.BundleKey,
	})

	if err != nil {
		return err
	}

	s.m.Lock()
	bundlePageErr := page.WriteFile(resource.BundleFilePath)
	s.m.Unlock()

	if bundlePageErr != nil {
		return bundlePageErr
	}

	s.m.Lock()
	configErr := resource.ConfiguratorPage.WriteFile(resource.ConfiguratorFilePath)
	s.m.Unlock()

	if configErr != nil {
		return configErr
	}

	bundleErr := s.settings.Bundler.Bundle(resource.ConfiguratorFilePath)
	if bundleErr != nil {
		return bundleErr
	}
	s.PackDurationSeconds = time.Since(startTime).Seconds()

	return nil
}

func (s *PackedComponent) RepackForWaitGroup(wg *sync.WaitGroup) error {
	err := s.Repack()
	wg.Done()

	return err
}

type PackedComponentList []*PackedComponent

func (l *PackedComponentList) RepackMany(hooks PackHooks) error {
	wg := &sync.WaitGroup{}
	for _, comp := range *l {
		wg.Add(1)
		go comp.RepackForWaitGroup(wg)
	}

	wg.Wait()

	return nil
}

type concPack struct {
	m           *sync.Mutex
	settings    *PackSettings
	packedPages []*PackedComponent
	packMap     map[string]bool
}

func (p *concPack) PackSingle(wg *sync.WaitGroup, dir string, hooks PackHooks) {
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
func (s *PackSettings) PackMany(pages []string, hooks PackHooks) ([]*PackedComponent, error) {
	cp := &concPack{
		m:           &sync.Mutex{},
		settings:    s,
		packedPages: make([]*PackedComponent, 0),
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
