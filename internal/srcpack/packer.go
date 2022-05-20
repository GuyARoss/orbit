// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package srcpack

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/GuyARoss/orbit/pkg/jsparse"
	"github.com/GuyARoss/orbit/pkg/log"
	"github.com/GuyARoss/orbit/pkg/webwrap"
)

type Packer interface {
	PackMany(pages []string) (PackedComponentList, error)
	PackSingle(logger log.Logger, file string) (PackComponent, error)
	ReattachLogger(logger log.Logger) Packer
}

// CachedEnvKeys represents a map where the key is the filepath
// for the env setting and where the value is a bundler key
type CachedEnvKeys map[string]string

// packer is the primary struct used for packing a directory of javascript files into
// valid web components.
type JSPacker struct {
	JsParser            jsparse.JSParser
	ValidWebWrappers    webwrap.JSWebWrapperList
	Logger              log.Logger
	SkipFirstPassBundle bool
	AssetDir            string
	WebDir              string
	cachedBundleKeys    CachedEnvKeys
}

// concpack is a private packing mechanism embedding the packer to pack a set of files concurrently.
type concPack struct {
	*JSPacker
	m sync.Mutex

	packedPages      []PackComponent
	packMap          map[string]bool
	cachedBundleKeys CachedEnvKeys
}

// packs the provided file paths into the orbit root directory
func (s *JSPacker) PackMany(pages []string) (PackedComponentList, error) {
	cp := &concPack{
		JSPacker:         s,
		packedPages:      make([]PackComponent, 0),
		packMap:          make(map[string]bool),
		cachedBundleKeys: s.cachedBundleKeys,
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(pages))

	var errOnce sync.Once
	var packErr error

	sh := NewSyncHook(s.Logger)

	defer sh.Close()

	for _, dir := range pages {
		// we copy dir here to avoid the pointer of dir being passed to our wrap func.
		t := dir
		// go routine to pack every page found in the pages directory
		// we wrap this routine with the sync hook to measure & log time deltas.
		go sh.WrapFunc(dir, func() {
			err := cp.PackSingle(wg, t)
			if err != nil {
				errOnce.Do(func() {
					packErr = err
				})
			}
		})
	}

	wg.Wait()

	return cp.packedPages, packErr
}

func (p *JSPacker) PackSingle(logger log.Logger, file string) (PackComponent, error) {
	return NewComponent(context.TODO(), &NewComponentOpts{
		DefaultKey:    p.cachedBundleKeys[file],
		FilePath:      file,
		WebDir:        p.WebDir,
		JSWebWrappers: p.ValidWebWrappers,
		JSParser:      p.JsParser,
	})
}

func (p *JSPacker) ReattachLogger(logger log.Logger) Packer {
	p.Logger = logger
	return p
}

// DefaultPackerOpts options for creating a new default packer
type DefaultPackerOpts struct {
	WebDir              string
	BundlerMode         string
	NodeModuleDir       string
	CachedBundleKeys    CachedEnvKeys
	SkipFirstPassBundle bool
}

// pack single packs a single file path into a usable web component
// this process includes the following:
// 1. wrapping the component with the specified front-end web framework.
// 2. bundling the component with the specified javascript bundler.
func (p *concPack) PackSingle(wg *sync.WaitGroup, path string) error {
	// this page has already been packed before and does not need to be repacked.
	if p.packMap[path] {
		wg.Done()
		return nil
	}

	// @@todo: we should validate if these components exist on our source map yet, if so we should
	// inherit the metadata, rather than generate new metadata.
	page, err := NewComponent(context.TODO(), &NewComponentOpts{
		DefaultKey:          p.cachedBundleKeys[path],
		FilePath:            path,
		WebDir:              p.WebDir,
		JSWebWrappers:       p.ValidWebWrappers,
		JSParser:            p.JsParser,
		SkipFirstPassBundle: p.SkipFirstPassBundle,
	})

	if err != nil {
		wg.Done()
		return err
	}

	p.m.Lock()
	p.packedPages = append(p.packedPages, page)
	p.packMap[path] = true
	p.m.Unlock()

	wg.Done()
	return nil
}

type PackedComponentList []PackComponent

func (l *PackedComponentList) RepackMany(logger log.Logger) error {
	wg := &sync.WaitGroup{}
	wg.Add(len(*l))

	var errOnce sync.Once
	var packErr error

	sh := NewSyncHook(logger)

	defer sh.Close()

	for _, comp := range *l {
		// we copy dir here to avoid the pointer of dir being passed to our wrap func.
		t := comp
		// go routine to pack every page found in the pages directory
		// we wrap this routine with the sync hook to measure & log time deltas.
		go sh.WrapFunc(t.OriginalFilePath(), func() {
			err := comp.RepackForWaitGroup(wg)
			if err != nil {
				errOnce.Do(func() {
					packErr = err
				})
			}
		})
	}

	wg.Wait()

	return packErr
}

// Write creates an audit file of all the current components to the specified file
func (l *PackedComponentList) Write(path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	f.WriteString("audit: components\n")

	for _, c := range *l {
		f.WriteString(fmt.Sprintf("%s %s", c.Name(), c.BundleKey()) + "\n")
	}

	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func NewDefaultPacker(logger log.Logger, opts *DefaultPackerOpts) Packer {
	packer := &JSPacker{
		WebDir:   opts.WebDir,
		JsParser: &jsparse.JSFileParser{},
		ValidWebWrappers: webwrap.NewActiveMap(&webwrap.BaseBundler{
			Mode:           webwrap.BundlerMode(opts.BundlerMode),
			WebDir:         opts.WebDir,
			PageOutputDir:  ".orbit/base/pages",
			NodeModulesDir: opts.NodeModuleDir,
			Logger:         logger,
		}),
		Logger:              logger,
		cachedBundleKeys:    opts.CachedBundleKeys,
		SkipFirstPassBundle: opts.SkipFirstPassBundle,
	}

	return packer
}
