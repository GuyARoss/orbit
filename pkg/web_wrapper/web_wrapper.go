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

func NewActiveMap() JSWebWrapperMap {
	return []JSWebWrapper{
		&ReactWebWrapper{},
	}
}

func (j *JSWebWrapperMap) FirstMatch(fileExtension string) JSWebWrapper {
	for _, f := range *j {
		if f.DoesSatisfyConstraints(fileExtension) {
			return f
		}
	}

	return nil
}
