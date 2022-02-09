package webwrapper

import (
	"fmt"

	"github.com/GuyARoss/orbit/pkg/jsparse"
)

type ReactWebWrap struct {
	*WebWrapSettings
}

func (s *ReactWebWrap) Apply(page jsparse.JSDocument, toFilePath string) jsparse.JSDocument {
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
