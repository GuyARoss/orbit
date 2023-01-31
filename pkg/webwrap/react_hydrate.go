package webwrap

import (
	"context"
	"fmt"
	"strings"

	"github.com/GuyARoss/orbit/pkg/embedutils"
	"github.com/GuyARoss/orbit/pkg/experiments"
	"github.com/GuyARoss/orbit/pkg/jsparse"
)

type ReactHydrate struct {
	csr *ReactCSR
}

func (s *ReactHydrate) Apply(page jsparse.JSDocument) (jsparse.JSDocument, error) {
	// react components should always be capitalized.
	if string(page.Name()[0]) != strings.ToUpper(string(page.Name()[0])) {
		return nil, ErrComponentExport
	}

	page.AddImport(&jsparse.ImportDependency{
		FinalStatement: "import ReactDOM from 'react-dom'",
		Type:           jsparse.ModuleImportType,
	})

	page.AddOther(fmt.Sprintf(
		"ReactDOM.hydrate(React.createElement(%s, JSON.parse(document.getElementById('orbit_manifest').textContent)), document.getElementById('%s_react_frame'))",
		page.Name(), page.Key()),
	)

	return page, nil
}

func (r *ReactHydrate) VerifyRequirements() error {
	return r.csr.VerifyRequirements()
}

func (s *ReactHydrate) Version() string {
	return "reactHydrate"
}

func (s *ReactHydrate) Stats() *WrapStats {
	if experiments.GlobalExperimentalFeatures.PreferSWCCompiler {
		return &WrapStats{
			WebVersion: "React Hydrate",
			Bundler:    "swc",
		}
	}

	return &WrapStats{
		WebVersion: "React Hydrate",
		Bundler:    "webpack",
	}
}

func (s *ReactHydrate) RequiredBodyDOMElements(ctx context.Context, cache *CacheDOMOpts) []string {
	return s.csr.RequiredBodyDOMElements(ctx, cache)
}

func (b *ReactHydrate) Setup(ctx context.Context, settings *BundleOpts) (*BundledResource, error) {
	return b.csr.Setup(ctx, settings)
}

func (b *ReactHydrate) Bundle(configuratorFilePath string, filePath string) error {
	return b.csr.Bundle(configuratorFilePath, filePath)
}

func (b *ReactHydrate) HydrationFile() []embedutils.FileReader {
	return b.csr.HydrationFile()
}

func NewReactHydrate(bundler *BaseBundler) *ReactHydrate {
	return &ReactHydrate{
		csr: NewReactCSR(bundler),
	}
}
