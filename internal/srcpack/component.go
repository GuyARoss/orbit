// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package srcpack

import (
	"context"
	"errors"
	"os"
	"sync"

	"github.com/GuyARoss/orbit/pkg/jsparse"
	"github.com/GuyARoss/orbit/pkg/webwrap"
)

type PackComponent interface {
	Repack() error
	RepackForWaitGroup(wg *sync.WaitGroup) error
	OriginalFilePath() string
	Dependencies() []*jsparse.ImportDependency
	BundleKey() string
	Name() string
	WebWrapper() webwrap.JSWebWrapper
	IsStaticResource() bool
	JsDocument() jsparse.JSDocument
}

// component that has been successfully ran, and output from a packing method.
// a component represents a source file that has a valid parser
type Component struct {
	WebDir           string
	JsParser         jsparse.JSParser
	webWrapper       webwrap.JSWebWrapper
	bundleKey        string
	dependencies     []*jsparse.ImportDependency
	originalFilePath string
	name             string
	m                *sync.Mutex
	isStaticResource bool
	document         jsparse.JSDocument
}

// NewComponentOpts options for creating a new component
type NewComponentOpts struct {
	FilePath            string
	WebDir              string
	DefaultKey          string
	JSParser            jsparse.JSParser
	JSWebWrappers       webwrap.JSWebWrapperList
	SkipFirstPassBundle bool
}

var ErrInvalidComponentType = errors.New("invalid component type")
var ErrInvalidPageName = errors.New("invalid page name")
var ErrComponentNotExported = errors.New("component not exported")
var ErrParserError = errors.New("page can not be parsed")

// NewComponent creates a new component that represents a packaged & bundled web component
func NewComponent(ctx context.Context, opts *NewComponentOpts) (PackComponent, error) {
	initPage, err := opts.JSParser.Parse(opts.FilePath, opts.WebDir)
	if initPage == nil {
		return nil, ErrParserError
	}

	if err != nil {
		return nil, err
	}

	if initPage == nil || initPage.DefaultExport() == nil || initPage.DefaultExport().Name == "" {
		return nil, ErrComponentNotExported
	}

	// we attempt to find the first web wrapper that satisfies the extension requirements
	// this same js wrapper will be used when we go to repack.
	wrapMethod := opts.JSWebWrappers.FindFirst(initPage)

	if wrapMethod == nil {
		return nil, ErrInvalidComponentType
	}

	if err := wrapMethod.VerifyRequirements(); err != nil {
		return nil, err
	}

	bundleKey := opts.DefaultKey
	if bundleKey == "" {
		bundleKey = initPage.Key()
	}

	resource, err := wrapMethod.Setup(ctx, &webwrap.BundleOpts{
		FileName:  opts.FilePath,
		BundleKey: bundleKey,
		Name:      initPage.Name(),
	})

	if err != nil {
		return nil, err
	}

	wrapAppliedPages, err := wrapMethod.Apply(initPage.Clone())
	if err != nil {
		return nil, err
	}

	for bundleOp, filePath := range resource.BundleOpFileDescriptor {
		_, err = os.Stat(filePath)
		// this addresses a performance issue that resulted in slow startup times for bundles that already existed.
		// this bundle process can be skipped during the initial startup of the application as long as stale bundle data exists.
		if !opts.SkipFirstPassBundle || errors.Is(err, os.ErrNotExist) {
			bundlePageErr := wrapAppliedPages[bundleOp].WriteFile(filePath)
			if bundlePageErr != nil {
				return nil, bundlePageErr
			}
		}
	}

	for _, r := range resource.Configurators {
		configErr := r.Page.WriteFile(r.FilePath)
		if configErr != nil {
			return nil, configErr
		}

		bundleErr := wrapMethod.Bundle(r.FilePath, opts.FilePath)
		if bundleErr != nil {
			return nil, bundleErr
		}
	}

	return &Component{
		name:             initPage.Name(),
		bundleKey:        bundleKey,
		dependencies:     initPage.Imports(),
		originalFilePath: opts.FilePath,
		m:                &sync.Mutex{},
		webWrapper:       wrapMethod,
		JsParser:         opts.JSParser,
		WebDir:           opts.WebDir,
		isStaticResource: len(initPage.DefaultExport().Args) == 0,
		document:         initPage,
	}, nil
}

func (s *Component) JsDocument() jsparse.JSDocument { return s.document }

func (s *Component) IsStaticResource() bool { return s.isStaticResource }

// Repack repacks a component following the following processes
//   - parses the provided filepath with the the components jsparser
//   - reapplies the component web wrapper
//   - bundles the component
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

	if err := s.webWrapper.VerifyRequirements(); err != nil {
		return err
	}

	// apply the necessary requirements for the web framework to the original page
	pages, err := s.webWrapper.Apply(page)
	if err != nil {
		return err
	}

	resource, err := s.webWrapper.Setup(context.TODO(), &webwrap.BundleOpts{
		FileName:  s.originalFilePath,
		BundleKey: s.BundleKey(),
		Name:      page.Name(),
	})

	if err != nil {
		return err
	}

	s.m.Lock()
	for opt, filePath := range resource.BundleOpFileDescriptor {
		err = pages[opt].WriteFile(filePath)
		if err != nil {
			return err
		}
	}

	for _, b := range resource.Configurators {
		err = b.Page.WriteFile(b.FilePath)
		if err != nil {
			return err
		}

		err = s.webWrapper.Bundle(b.FilePath, s.originalFilePath)
		if err != nil {
			return err
		}
	}

	s.m.Unlock()

	return nil
}

// RepackForWaitGroup given a wait group, repacks the component using the underlying "Repack" method.
func (s *Component) RepackForWaitGroup(wg *sync.WaitGroup) error {
	err := s.Repack()

	if err != nil {
		return err
	}

	wg.Done()
	return nil
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

// parsePath is a utility that verifies that the provided path is of a valid structure
func parsePath(p string) string {
	skip := 0
	for _, c := range p {
		if c == '.' || c == '/' {
			skip += 1
			continue
		}

		break
	}
	return p[skip:]
}

type PackComponentFileMap map[string]PackComponent

// Find finds the provided key if one exists
func (m PackComponentFileMap) Find(key string) PackComponent {
	return m[parsePath(key)]
}

// Find finds and returns the the first component with provided bundle key
func (m PackComponentFileMap) FindBundleKey(key string) PackComponent {
	for _, c := range m {
		if c.BundleKey() == key {
			return c
		}
	}

	return nil
}

func (m PackComponentFileMap) Set(component PackComponent) {
	m[parsePath(component.OriginalFilePath())] = component
}
