// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package webwrap

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/GuyARoss/orbit/pkg/bundler"
	"github.com/GuyARoss/orbit/pkg/jsparse"
)

type ReactWebWrapper struct {
	*BaseWebWrapper
}

var ErrComponentExport = errors.New("prefer capitalization for jsx components")
var ErrInvalidComponent = errors.New("invalid jsx component")

const reactExtension string = "jsx"

func (s *ReactWebWrapper) Apply(page jsparse.JSDocument) (jsparse.JSDocument, error) {
	if page.Extension() != reactExtension {
		return nil, ErrInvalidComponent
	}

	// react components should always be capitalized.
	if string(page.Name()[0]) != strings.ToUpper(string(page.Name()[0])) {
		return nil, ErrComponentExport
	}

	page.AddImport(&jsparse.ImportDependency{
		FinalStatement: "import ReactDOM from 'react-dom'",
		Type:           jsparse.ModuleImportType,
	})

	page.AddOther(fmt.Sprintf(
		"ReactDOM.render(<%s {...JSON.parse(document.getElementById('orbit_manifest').textContent)}/>, document.getElementById('root'))",
		page.Name()),
	)

	return page, nil
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
	return strings.Contains(fileExtension, reactExtension)
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
		uris = append(uris, "https://unpkg.com/react-dom/umd/react-dom.development.js")
	case bundler.ProductionBundle:
		uris = append(uris, "https://unpkg.com/react/umd/react.production.min.js")
		uris = append(uris, "https://unpkg.com/react-dom/umd/react-dom.production.min.js")
	}

	files, _ := cache.CacheWebRequest(uris)

	// currently these files are just paths to a directory to refer
	// to them on the dom, we need to convert them to <script> tags.
	for i, f := range files {
		files[i] = fmt.Sprintf(`<script src="%s"></script>`, f)
	}

	files = append(files, `<div id="root"></div>`)

	return files
}
