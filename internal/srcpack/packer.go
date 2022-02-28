package srcpack

import (
	"context"
	"fmt"
	"sync"

	"github.com/GuyARoss/orbit/pkg/bundler"
	"github.com/GuyARoss/orbit/pkg/jsparse"
	"github.com/GuyARoss/orbit/pkg/log"
	webwrapper "github.com/GuyARoss/orbit/pkg/web_wrapper"
)

// packer is the primary struct used for packing a directory of javascript files into
// valid web components.
type Packer struct {
	Bundler          bundler.Bundler
	JsParser         jsparse.JSParser
	ValidWebWrappers webwrapper.JSWebWrapperMap
	Logger           log.Logger

	AssetDir string
	WebDir   string
}

// packs the provided file paths into the orbit root directory
func (s *Packer) PackMany(pages []string) ([]*Component, error) {
	cp := &concPack{
		Packer:      s,
		packedPages: make([]*Component, 0),
		packMap:     make(map[string]bool),
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(pages))

	errchan := make(chan error)

	go func() {
		err := <-errchan
		// @@todo: do something more with this error?
		fmt.Println("error occurred", err.Error())
	}()

	sh := NewSyncHook(s.Logger)

	defer sh.Close()

	for _, dir := range pages {
		// we copy dir here to avoid the pointer of dir being passed to our wrap func.
		t := dir
		// go routine to pack every page found in the pages directory
		// we wrap this routine with the sync hook to measure & log time deltas.
		go sh.WrapFunc(dir, func() { cp.PackSingle(errchan, wg, t) })
	}

	wg.Wait()

	return cp.packedPages, nil
}

func (p *Packer) ReattachLogger(logger log.Logger) *Packer {
	p.Logger = logger
	return p
}

type DefaultPackerOpts struct {
	WebDir        string
	BundlerMode   string
	NodeModuleDir string
}

func NewDefaultPacker(logger log.Logger, opts *DefaultPackerOpts) *Packer {
	return &Packer{
		Bundler: &bundler.WebPackBundler{
			BaseBundler: &bundler.BaseBundler{
				Mode:           bundler.BundlerMode(opts.BundlerMode),
				WebDir:         opts.WebDir,
				PageOutputDir:  ".orbit/base/pages",
				NodeModulesDir: opts.NodeModuleDir,
				Logger:         logger,
			},
		},
		WebDir:           opts.WebDir,
		JsParser:         &jsparse.JSFileParser{},
		ValidWebWrappers: webwrapper.NewActiveMap(),
		Logger:           logger,
	}
}

// concpack is a private packing mechanism embedding the packer to pack a set of files concurrently.
type concPack struct {
	*Packer
	m sync.Mutex

	packedPages []*Component
	packMap     map[string]bool
}

// pack single packs a single file path into a usable web component
// this process includes the following:
// 1. wrapping the component with the specified front-end web framework.
// 2. bundling the component with the specified javascript bundler.
func (p *concPack) PackSingle(errchan chan error, wg *sync.WaitGroup, path string) {
	// @@todo: we should validate if these components exist on our source map yet, if so we should
	// inherit the metadata, rather than generate new metadata.
	page, err := NewComponent(context.TODO(), &NewComponentOpts{
		FilePath:      path,
		WebDir:        p.WebDir,
		JSWebWrappers: p.ValidWebWrappers,
		Bundler:       p.Bundler,
		JSParser:      p.JsParser,
	})

	if err != nil {
		errchan <- err
		fmt.Println(err)

		wg.Done()
		return
	}

	if p.packMap[page.Name] {
		// this page has already been packed before
		// and does not need to be repacked.
		wg.Done()
		return
	}

	p.m.Lock()
	p.packedPages = append(p.packedPages, page)
	p.packMap[page.Name] = true
	p.m.Unlock()

	wg.Done()
}

type PackedComponentList []*Component

func (l *PackedComponentList) RepackMany(logger log.Logger) error {
	wg := &sync.WaitGroup{}
	wg.Add(len(*l))

	errchan := make(chan error)

	go func() {
		err := <-errchan
		// @@todo: do something more with this error?
		fmt.Println("error occurred", err.Error())
	}()

	sh := NewSyncHook(logger)

	defer sh.Close()

	for _, comp := range *l {
		// we copy dir here to avoid the pointer of dir being passed to our wrap func.
		t := comp
		// go routine to pack every page found in the pages directory
		// we wrap this routine with the sync hook to measure & log time deltas.
		go sh.WrapFunc(t.originalFilePath, func() { comp.RepackForWaitGroup(wg, errchan) })
	}

	wg.Wait()

	return nil
}
