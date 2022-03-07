package srcpack

import (
	"context"
	"errors"
	"sync"

	"github.com/GuyARoss/orbit/pkg/bundler"
	"github.com/GuyARoss/orbit/pkg/jsparse"
	webwrapper "github.com/GuyARoss/orbit/pkg/web_wrapper"
)

// component that has been successfully ran, and output from a packing method.
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

type NewComponentOpts struct {
	FilePath   string
	WebDir     string
	DefaultKey string

	JSParser      jsparse.JSParser
	Bundler       bundler.Bundler
	JSWebWrappers webwrapper.JSWebWrapperMap
}

var ErrInvalidComponentType = errors.New("invalid component type")

func NewComponent(ctx context.Context, opts *NewComponentOpts) (*Component, error) {
	page, err := opts.JSParser.Parse(opts.FilePath, opts.WebDir)
	if err != nil {
		return nil, err
	}

	// we attempt to find the first web wrapper that satisfies the extension requirements
	// this same js wrapper will be used when we go to repack.
	webwrap := opts.JSWebWrappers.FirstMatch(page.Extension())

	// no webwrapper is available
	if webwrap == nil {
		return nil, ErrInvalidComponentType
	}

	page = webwrap.Apply(page, opts.FilePath)

	bundleKey := opts.DefaultKey
	if bundleKey == "" {
		bundleKey = page.Key()
	}

	resource, err := opts.Bundler.Setup(ctx, &bundler.BundleOpts{
		FileName:  opts.FilePath,
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

	bundleErr := opts.Bundler.Bundle(resource.ConfiguratorFilePath)
	if bundleErr != nil {
		return nil, bundleErr
	}

	return &Component{
		Name:             page.Name(),
		BundleKey:        bundleKey,
		dependencies:     page.Imports(),
		originalFilePath: opts.FilePath,
		m:                &sync.Mutex{},
		WebWrapper:       webwrap,
		Bundler:          opts.Bundler,
		JsParser:         opts.JSParser,
		WebDir:           opts.WebDir,
	}, nil
}

func (s *Component) OriginalFilePath() string {
	return s.originalFilePath
}

func (s *Component) Dependencies() []*jsparse.ImportDependency {
	return s.dependencies
}

func (s *Component) Repack() error {
	// parse the original javascript page, provided our javascript parser.
	// we later mutate this page to apply the rest of the required web wrapper
	page, err := s.JsParser.Parse(s.originalFilePath, s.WebDir)
	if err != nil {
		return err
	}

	// apply the necessary requirements for the web framework to the original page
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

func (s *Component) RepackForWaitGroup(wg *sync.WaitGroup, c chan error) {
	err := s.Repack()

	if err != nil {
		c <- err
	}

	wg.Done()
}
