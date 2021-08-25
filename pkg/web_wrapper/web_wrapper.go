package webwrapper

import "github.com/GuyARoss/orbit/pkg/jsparse"

type WebWrapSettings struct {
	WebDir string
}

type WebWrapper interface {
	Apply(page *jsparse.Page, toFilePath string) *jsparse.Page
}
