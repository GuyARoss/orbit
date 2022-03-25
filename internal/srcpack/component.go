// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package srcpack

import (
	"context"
	"errors"
	"sync"

	"github.com/GuyARoss/orbit/pkg/bundler"
	"github.com/GuyARoss/orbit/pkg/jsparse"
	"github.com/GuyARoss/orbit/pkg/webwrap"
)

type PackComponent interface {
	Repack() error
	RepackForWaitGroup(wg *sync.WaitGroup, c chan error)
	OriginalFilePath() string
	Dependencies() []*jsparse.ImportDependency
	BundleKey() string
	Name() string
	WebWrapper() webwrap.JSWebWrapper
}

// component that has been successfully ran, and output from a packing method.
// a component represents a source file that has a valid parser
type Component struct {
	WebDir string

	Bundler  bundler.Bundler
	JsParser jsparse.JSParser

	webWrapper       webwrap.JSWebWrapper
	bundleKey        string
	dependencies     []*jsparse.ImportDependency
	originalFilePath string
	name             string
	m                *sync.Mutex
}

// NewComponentOpts options for creating a new component
type NewComponentOpts struct {
	FilePath   string
	WebDir     string
	DefaultKey string

	JSParser      jsparse.JSParser
	Bundler       bundler.Bundler
	JSWebWrappers webwrap.JSWebWrapperList
}

var ErrInvalidComponentType = errors.New("invalid component type")

var ErrInvalidPageName = errors.New("invalid page name")

// NewComponent creates a new component that represents a packaged & bundled web component
func NewComponent(ctx context.Context, opts *NewComponentOpts) (PackComponent, error) {
	page, err := opts.JSParser.Parse(opts.FilePath, opts.WebDir)
	if err != nil {
		return nil, err
	}

	// we attempt to find the first web wrapper that satisfies the extension requirements
	// this same js wrapper will be used when we go to repack.
	webwrap := opts.JSWebWrappers.FirstMatch(page.Extension())

	if webwrap == nil {
		return nil, ErrInvalidComponentType
	}

	page, err = webwrap.Apply(page)
	if err != nil {
		return nil, err
	}

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
		name:             page.Name(),
		bundleKey:        bundleKey,
		dependencies:     page.Imports(),
		originalFilePath: opts.FilePath,
		m:                &sync.Mutex{},
		webWrapper:       webwrap,
		Bundler:          opts.Bundler,
		JsParser:         opts.JSParser,
		WebDir:           opts.WebDir,
	}, nil
}

// Repack repacks a component following the following processes
// 	- parses the provided filepath with the the components jsparser
// 	- reapplies the component web wrapper
// 	- bundles the component
func (s *Component) Repack() error {
	// parse the original javascript page, provided our javascript parser.
	// we later mutate this page to apply the rest of the required web wrapper
	page, err := s.JsParser.Parse(s.originalFilePath, s.WebDir)
	if err != nil {
		return err
	}

	if page.Name() == "" {
		return ErrInvalidPageName
	}

	// apply the necessary requirements for the web framework to the original page
	page, err = s.WebWrapper().Apply(page)
	if err != nil {
		return err
	}

	resource, err := s.Bundler.Setup(context.TODO(), &bundler.BundleOpts{
		FileName:  s.originalFilePath,
		BundleKey: s.BundleKey(),
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

// RepackForWaitGroup given a wait group, repacks the component using the underlying "Repack" method.
func (s *Component) RepackForWaitGroup(wg *sync.WaitGroup, c chan error) {
	err := s.Repack()

	if err != nil {
		c <- err
	}

	wg.Done()
}

// OriginalFilePath returns the original file path on the component
func (s *Component) OriginalFilePath() string { return s.originalFilePath }

// Dependencies returns the dependencies on the component
func (s *Component) Dependencies() []*jsparse.ImportDependency { return s.dependencies }

// BundleKey returns the bundle key for the packed component
func (s *Component) BundleKey() string { return s.bundleKey }

// Name returns the name of the component
func (s *Component) Name() string { return s.name }

// WebWrapper returns the instance of the webwrapper applied to the component
func (s *Component) WebWrapper() webwrap.JSWebWrapper { return s.webWrapper }
