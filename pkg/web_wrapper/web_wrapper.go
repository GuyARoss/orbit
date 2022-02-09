package webwrapper

import "github.com/GuyARoss/orbit/pkg/jsparse"

type WebWrapSettings struct {
	WebDir string
}

type WebWrapper interface {
	Apply(page jsparse.JSDocument, toFilePath string) jsparse.JSDocument
	NodeDependencies() map[string]string
	DoesSatisfyConstraints(fileExtension string) bool
}
