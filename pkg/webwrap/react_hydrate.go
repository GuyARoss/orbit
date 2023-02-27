// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

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
	ssr *PartialWrapReactSSR
}

func (s *ReactHydrate) Apply(page jsparse.JSDocument) (map[string]jsparse.JSDocument, error) {
	// react components should always be capitalized.
	if string(page.Name()[0]) != strings.ToUpper(string(page.Name()[0])) {
		return nil, ErrComponentExport
	}

	csrHydratePage := page.Clone()

	csrHydratePage.AddImport(&jsparse.ImportDependency{
		FinalStatement: "import ReactDOM from 'react-dom'",
		Type:           jsparse.ModuleImportType,
	})

	csrHydratePage.AddOther("// testing", fmt.Sprintf(
		"ReactDOM.hydrate(React.createElement(%s, JSON.parse(document.getElementById('orbit_manifest').textContent)), document.getElementById('%s_react_frame'))",
		page.Name(), page.Key()),
	)

	ssrPage, err := s.ssr.Apply(page.Clone())
	if err != nil {
		return nil, err
	}

	response := map[string]jsparse.JSDocument{
		"csr": csrHydratePage,
		"ssr": ssrPage,
	}

	return response, nil
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
	page := jsparse.NewEmptyDocument()

	page.AddImport(&jsparse.ImportDependency{
		FinalStatement: "const {merge} = require('webpack-merge')",
		Type:           jsparse.ModuleImportType,
	})

	if experiments.GlobalExperimentalFeatures.PreferSWCCompiler {
		page.AddImport(&jsparse.ImportDependency{
			FinalStatement: "const baseConfig = require('../../assets/swc-base.config.js')",
			Type:           jsparse.ModuleImportType,
		})
	} else {
		page.AddImport(&jsparse.ImportDependency{
			FinalStatement: "const baseConfig = require('../../assets/base.config.js')",
			Type:           jsparse.ModuleImportType,
		})
	}

	outputFileName := fmt.Sprintf("%s.js", settings.BundleKey)
	clientBundleFilePath := fmt.Sprintf("%s/%s.js", b.csr.PageOutputDir, settings.BundleKey)

	page.AddOther(fmt.Sprintf(`module.exports = merge(baseConfig, {
		entry: ['./%s'],
		mode: '%s',
		output: {
			filename: '%s'
		},
	})`, clientBundleFilePath, string(b.csr.Mode), outputFileName))

	b.ssr.sourceMapDoc.AddImport(&jsparse.ImportDependency{
		FinalStatement: fmt.Sprintf("import %s from '%s'", settings.Name, fmt.Sprintf("./%s.ssr.js", settings.BundleKey)),
		Type:           jsparse.LocalImportType,
	})

	b.ssr.sourceMapDoc.AddOther(fmt.Sprintf(`export const %s = (d) => ReactDOMServer.renderToString(<%s {...d}/>)`, strings.ToLower(settings.Name), settings.Name))
	b.ssr.initDoc.AddImport(&jsparse.ImportDependency{
		FinalStatement: fmt.Sprintf("import { %s } from '%s'", strings.ToLower(settings.Name), fmt.Sprintf("./%s", "react_ssr.map.js")),
		Type:           jsparse.LocalImportType,
	})

	b.ssr.jsSwitch.Add(jsparse.JSString, settings.BundleKey, fmt.Sprintf(`return %s(JSON.parse(JSONData))`, strings.ToLower(settings.Name)))

	return &BundledResource{
		BundleOpFileDescriptor: map[string]string{
			"csr": clientBundleFilePath,
			"ssr": fmt.Sprintf("%s/%s.ssr.js", b.ssr.PageOutputDir, settings.BundleKey),
		},
		Configurators: []BundleConfigurator{
			{
				FilePath: fmt.Sprintf("%s/react_ssr.map.js", b.ssr.PageOutputDir),
				Page:     b.ssr.sourceMapDoc,
			}, {
				FilePath: fmt.Sprintf("%s/react_ssr.js", b.ssr.PageOutputDir),
				Page:     b.ssr.initDoc,
			},
			{
				FilePath: fmt.Sprintf("%s/%s.config.js", b.csr.PageOutputDir, settings.BundleKey),
				Page:     page,
			},
		},
	}, nil
}

func (b *ReactHydrate) Bundle(configuratorFilePath string, filePath string) error {
	if strings.Contains(configuratorFilePath, "ssr") { // todo: avoid doing this.
		return nil
	}

	return b.csr.Bundle(configuratorFilePath, filePath)
}

func (b *ReactHydrate) DoesSatisfyConstraints(page jsparse.JSDocument) bool {
	return page.Extension() == "jsx" && page.DefaultExport() != nil
}

func (b *ReactHydrate) HydrationFile() []embedutils.FileReader {
	files, err := embedFiles.ReadDir("embed")
	if err != nil {
		return nil
	}
	u := []embedutils.FileReader{}
	for _, file := range files {
		if strings.Contains(file.Name(), "react_hydrate.go") {
			u = append(u, &embedFileReader{fileName: file.Name()})
			continue
		}
		if strings.Contains(file.Name(), "react_ssr.go") {
			u = append(u, &embedFileReader{fileName: file.Name()})
			continue
		}
		if strings.Contains(file.Name(), "pb.go") {
			u = append(u, &embedFileReader{fileName: file.Name()})
		}
	}
	return u
}

func NewReactHydrate(bundler *BaseBundler) JSWebWrapper {
	return &ReactHydrate{
		csr: NewReactCSR(bundler),
		ssr: NewReactSSRPartial(&NewReactSSROpts{
			Bundler:      bundler,
			SourceMapDoc: jsparse.NewEmptyDocument(),
			InitDoc:      jsparse.NewEmptyDocument(),
		}),
	}
}
