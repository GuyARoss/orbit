package internal

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/GuyARoss/orbit/pkg/bundler"
	"github.com/GuyARoss/orbit/pkg/fs"
	"github.com/GuyARoss/orbit/pkg/jsparse"
	"github.com/GuyARoss/orbit/pkg/log"
	webwrapper "github.com/GuyARoss/orbit/pkg/web_wrapper"
	"github.com/google/uuid"
)

type PackSettings struct {
	Bundler    bundler.Bundler
	WebWrapper webwrapper.WebWrapper

	AssetDir string
	WebDir   string
}

func (s *PackSettings) CopyAssets() ([]*fs.CopyResults, error) {
	results := fs.CopyDir(s.AssetDir, s.AssetDir, ".orbit/assets", false)

	return results, nil
}

func hashKey(name string) string {
	id := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(name))

	return strings.ReplaceAll(id.String(), "-", "")
}

// PackedComponent
// Web-Component that has been successfully ran, and output from a packing method.
type PackedComponent struct {
	PageName            string
	BundleKey           string
	OriginalFilePath    string
	PackDurationSeconds float64
	Dependencies        []*jsparse.ImportDependency

	settings *PackSettings
}

// PackSingle
// Packs the a single file paths into the orbit root directory
// Process includes:
// - Wrapping the component with the specified front-end web framework.
// - Bundling the component with the specified javascript bundler.
func (s *PackSettings) PackSingle(pageFilePath string) (*PackedComponent, error) {
	startTime := time.Now()
	page, err := jsparse.ParsePage(pageFilePath, s.WebDir)
	if err != nil {
		return nil, err
	}

	page = s.WebWrapper.Apply(page, pageFilePath)

	bundleKey := hashKey(page.Name)
	resource, err := s.Bundler.Setup(&bundler.BundleSetupSettings{
		FileName:  pageFilePath,
		BundleKey: bundleKey,
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
		PageName:            page.Name,
		Dependencies:        page.Imports,
		BundleKey:           bundleKey,
		OriginalFilePath:    pageFilePath,
		PackDurationSeconds: time.Since(startTime).Seconds(),
		settings:            s,
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

func (s *PackedComponent) Repack(hooks PackHooks) error {
	if hooks != nil {
		hooks.Pre(s.OriginalFilePath)
	}
	startTime := time.Now()

	page, err := jsparse.ParsePage(s.OriginalFilePath, s.settings.WebDir)
	if err != nil {
		return err
	}

	page = s.settings.WebWrapper.Apply(page, s.OriginalFilePath)
	resource, err := s.settings.Bundler.Setup(&bundler.BundleSetupSettings{
		FileName:  s.OriginalFilePath,
		BundleKey: s.BundleKey,
	})

	if err != nil {
		return err
	}

	bundlePageErr := page.WriteFile(resource.BundleFilePath)
	if bundlePageErr != nil {
		return bundlePageErr
	}

	configErr := resource.ConfiguratorPage.WriteFile(resource.ConfiguratorFilePath)
	if configErr != nil {
		return configErr
	}

	bundleErr := s.settings.Bundler.Bundle(resource.ConfiguratorFilePath)
	if bundleErr != nil {
		return bundleErr
	}
	s.PackDurationSeconds = time.Since(startTime).Seconds()

	if hooks != nil {
		hooks.Post(s.PackDurationSeconds)
	}

	return nil
}

type concPack struct {
	m           *sync.Mutex
	settings    *PackSettings
	packedPages []*PackedComponent
	packMap     map[string]bool
}

func (p *concPack) PackSingle(wg *sync.WaitGroup, dir string, hooks PackHooks) {
	p.m.Lock()
	if hooks != nil {
		hooks.Pre(dir)
	}
	p.m.Unlock()

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
		cp.PackSingle(wg, dir, hooks)
	}

	return cp.packedPages, nil
}
