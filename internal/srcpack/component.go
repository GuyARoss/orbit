package srcpack

import (
	"context"
	"sync"

	"github.com/GuyARoss/orbit/pkg/bundler"
	"github.com/GuyARoss/orbit/pkg/jsparse"
	webwrapper "github.com/GuyARoss/orbit/pkg/web_wrapper"
)

// PackedComponent
// Web-Component that has been successfully ran, and output from a packing method.
type Component struct {
	WebDir string

	Name      string
	BundleKey string

	WebWrapper webwrapper.JSWebWrapper
	Bundler    bundler.Bundler
	JsParser   jsparse.JSParser

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
	// parse original page
	page, err := s.JsParser.Parse(s.originalFilePath, s.WebDir)
	if err != nil {
		return err
	}

	// apply the nessasacary requirements for the web framework to the original page
	page = s.WebWrapper.Apply(page, s.originalFilePath)
	resource, err := s.Bundler.Setup(context.TODO(), &bundler.BundleOpts{
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
