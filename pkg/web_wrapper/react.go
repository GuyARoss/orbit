package webwrapper

import (
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
		"react":            "^16.13.1",
		"react-dom":        "^16.13.1",
		"react-hot-loader": "^4.12.21",
		"react-router-dom": "^5.2.0",
	}
}

func (s *ReactWebWrapper) DoesSatisfyConstraints(fileExtension string) bool {
	return strings.Contains(fileExtension, ".jsx")
}
