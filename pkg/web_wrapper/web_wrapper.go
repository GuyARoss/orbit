package webwrapper

import "github.com/GuyARoss/orbit/pkg/jsparse"

type BaseWebWrapper struct {
	WebDir string
}

type JSWebWrapper interface {
	Apply(page jsparse.JSDocument, toFilePath string) jsparse.JSDocument
	NodeDependencies() map[string]string
	DoesSatisfyConstraints(fileExtension string) bool
}

type JSWebWrapperMap []JSWebWrapper

func NewMap() JSWebWrapperMap {
	return []JSWebWrapper{
		&ReactWebWrapper{},
	}
}
