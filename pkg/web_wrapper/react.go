// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package webwrapper

import (
	"context"
	"fmt"
	"strings"

	"github.com/GuyARoss/orbit/pkg/bundler"
	"github.com/GuyARoss/orbit/pkg/jsparse"
)

type ReactWebWrapper struct {
	*BaseWebWrapper
}

func (s *ReactWebWrapper) Apply(page jsparse.JSDocument, toFilePath string) jsparse.JSDocument {
	page.AddImport(&jsparse.ImportDependency{
		FinalStatement: "import ReactDOM from 'react-dom'",
		Type:           jsparse.ModuleImportType,
	})

	page.AddOther(fmt.Sprintf(
		"ReactDOM.render(<%s {...JSON.parse(document.getElementById('orbit_manifest').textContent)}/>, document.getElementById('root'))",
		page.Name()),
	)

	return page
}

func (s *ReactWebWrapper) NodeDependencies() map[string]string {
	return map[string]string{
		"react":            "latest",
		"react-dom":        "latest",
		"react-hot-loader": "latest",
		"react-router-dom": "latest",
	}
}

func (s *ReactWebWrapper) DoesSatisfyConstraints(fileExtension string) bool {
	return strings.Contains(fileExtension, "jsx")
}

func (s *ReactWebWrapper) Version() string {
	return "react"
}

func (s *ReactWebWrapper) RequiredBodyDOMElements(ctx context.Context, cache *CacheDOMOpts) []string {
	mode := ctx.Value(bundler.BundlerID).(string)

	uris := make([]string, 0)
	switch bundler.BundlerMode(mode) {
	case bundler.DevelopmentBundle:
		uris = append(uris, "https://unpkg.com/react/umd/react.development.js")
		uris = append(uris, "https://unpkg.com/react/umd/react.development.js")
	case bundler.ProductionBundle:
		uris = append(uris, "https://unpkg.com/react/umd/react.production.min.js")
		uris = append(uris, "https://unpkg.com/react/umd/react.production.min.js")
	}

	files, err := cache.CacheWebRequest(uris)

	if err != nil {
		fmt.Println(err)
	}

	// currently these files are just paths to a directory to refer
	// to them on the dom, we need to convert them to <script> tags.
	for i, f := range files {
		files[i] = fmt.Sprintf(`<script src="%s"></script>`, f)
	}

	files = append(files, `<div id="root"></div>`)

	return files
}
