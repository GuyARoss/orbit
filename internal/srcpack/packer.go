package srcpack

import (
	"context"
	"sync"

	"github.com/GuyARoss/orbit/pkg/bundler"
	"github.com/GuyARoss/orbit/pkg/fs"
	"github.com/GuyARoss/orbit/pkg/jsparse"
	webwrapper "github.com/GuyARoss/orbit/pkg/web_wrapper"
)

// packer is the primary struct used for packing a directory of javascript files into
// valid web components.
type Packer struct {
	Bundler          bundler.Bundler
	JsParser         jsparse.JSParser
	ValidWebWrappers webwrapper.JSWebWrapperMap

	AssetDir string
	WebDir   string
}

// copies the required assets to the asset directory
func (s *Packer) CopyAssets() ([]*fs.CopyResults, error) {
	results := fs.CopyDir(s.AssetDir, s.AssetDir, ".orbit/assets", false)

	return results, nil
}

// concpack is a private packing mechanism embedding the packer to pack a set of files concurrently.
type concPack struct {
	*Packer

	packedPages []*Component
	packMap     map[string]bool
}

// pack single packs a single file path into a usable web component
// this process includes the following:
// 1. wrapping the component with the specified front-end web framework.
// 2. bundling the component with the specified javascript bundler.
func (p *concPack) PackSingle(wg *sync.WaitGroup, path string) {
	// @@todo: we should validate if these components exist on our source map yet, if so we should
	// inherit the metadata, rather than generate new metadata.
	page, err := NewComponent(context.TODO(), &NewComponentOpts{
		FilePath:      path,
		WebDir:        p.WebDir,
		JSWebWrappers: p.ValidWebWrappers,
		Bundler:       p.Bundler,
		JSParser:      p.JsParser,
	})

	if p.packMap[page.Name] {
		return
	}

	if err != nil {
		// @@report error with packing via channel
		return
	}

	p.packedPages = append(p.packedPages, page)
	p.packMap[page.Name] = true

	wg.Done()
}

// packs the provoided file paths into the orbit root directory
func (s *Packer) PackMany(pages []string, hooks Hooks) ([]*Component, error) {
	cp := &concPack{
		Packer:      s,
		packedPages: make([]*Component, 0),
		packMap:     make(map[string]bool),
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(pages))

	for _, dir := range pages {
		// currently using a go routine to pack every page found in the pages directory
		// in the future, this should be wrapped with a routine to measure & log time deltas.
		go cp.PackSingle(wg, dir)
	}

	wg.Wait()

	return cp.packedPages, nil
}
