package srcpack

import (
	"sync"
	"time"

	"github.com/GuyARoss/orbit/pkg/bundler"
	"github.com/GuyARoss/orbit/pkg/jsparse"
)

// PackedComponent
// Web-Component that has been successfully ran, and output from a packing method.
type Component struct {
	*Packer

	PageName            string
	BundleKey           string
	PackDurationSeconds float64

	dependencies     []*jsparse.ImportDependency
	originalFilePath string
	m                *sync.Mutex
}

func (s *Component) OriginalFilePath() string {
	return s.originalFilePath
}

func (s *Component) Dependencies() []*jsparse.ImportDependency {
	return s.dependencies
}

func (s *Component) Repack() error {
	startTime := time.Now()

	// parse original page
	page, err := s.JsParser.Parse(s.originalFilePath, s.WebDir)
	if err != nil {
		return err
	}

	// apply the nessasacary requirements for the web framework to the original page
	page = s.WebWrapper.Apply(page, s.originalFilePath)
	resource, err := s.Bundler.Setup(&bundler.BundleSetupSettings{
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

	bundleErr := s.Bundler.Bundle(resource.ConfiguratorFilePath)
	if bundleErr != nil {
		return bundleErr
	}
	s.PackDurationSeconds = time.Since(startTime).Seconds()

	return nil
}

func (s *Component) RepackForWaitGroup(wg *sync.WaitGroup) error {
	err := s.Repack()
	wg.Done()

	return err
}

type PackedComponentList []*Component

func (l *PackedComponentList) RepackMany(hooks Hooks) error {
	wg := &sync.WaitGroup{}
	for _, comp := range *l {
		wg.Add(1)
		go comp.RepackForWaitGroup(wg)
	}

	wg.Wait()

	return nil
}
