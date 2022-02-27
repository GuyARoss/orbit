package webwrapper

import (
	"context"
	"fmt"
	"strings"

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
	return "react-v16.13.1"
}

func (s *ReactWebWrapper) RequiredBodyDOMElements(ctx context.Context, cache *CacheDOMOpts) []string {
	files := cache.CacheWebRequest([]string{
		"https://unpkg.com/react/umd/react.production.min.js",
		"https://unpkg.com/react-dom/umd/react-dom.production.min.js",
	})

	files = append(files, `<div id="root"></div>`)

	return files
}
