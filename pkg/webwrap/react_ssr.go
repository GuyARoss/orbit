package webwrap

import (
	"context"
	"strings"

	"github.com/GuyARoss/orbit/pkg/embedutils"
	"github.com/GuyARoss/orbit/pkg/jsparse"
)

type ReactSSR struct{}

func NewReactSSR() *ReactSSR {
	return &ReactSSR{}
}

func (r *ReactSSR) RequiredBodyDOMElements(context.Context, *CacheDOMOpts) []string {
	// nothing goes here
	return []string{}
}

func (r *ReactSSR) Setup(context.Context, *BundleOpts) (*BundledResource, error) {
	// @@ create the manifest node application that includes a reference to the react component
	return nil, nil
}

func (r *ReactSSR) Apply(jsparse.JSDocument) (jsparse.JSDocument, error) {
	// @@ use original page
	return nil, nil
}
func (r *ReactSSR) DoesSatisfyConstraints(string) bool {
	// @@ ensure extension is jsx
	return false
}
func (r *ReactSSR) Version() string {
	return "reactssr"
}
func (r *ReactSSR) Bundle(configuratorFilePath string) error {
	// @@ add to the app file thingy
	return nil
}

func (r *ReactSSR) HydrationFile() embedutils.FileReader {
	files, err := embedFiles.ReadDir("embed")
	if err != nil {
		return nil
	}

	for _, file := range files {
		if strings.Contains(file.Name(), "react_ssr.go") {
			return &embedFileReader{fileName: file.Name()}
		}
	}
	return nil
}
