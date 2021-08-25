package webwrapper

import (
	"fmt"

	"github.com/GuyARoss/orbit/pkg/jsparse"
)

type ReactWebWrap struct {
	*WebWrapSettings
}

func (s *ReactWebWrap) Apply(page *jsparse.Page, toFilePath string) *jsparse.Page {
	page.Imports = append(page.Imports, "import ReactDOM from 'react-dom'")
	page.Other = append(page.Other, fmt.Sprintf("ReactDOM.render(<%s {...JSON.parse(document.getElementById('orbit_manifest').textContent)}/>, document.getElementById('root'))", page.Name))

	return page
}
