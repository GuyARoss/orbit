package webwrap

import (
	"io/ioutil"
	"strings"
)

// htmlDoc represents a basic document model that will be rendered upon build request
type htmlDoc struct {
	Head []string
	Body []string
}

func innerHTML(str string, start string, end string) string {
	return strings.Split(strings.Join(strings.Split(str, start)[1:], ""), end)[0]
}

func DocFromFile(path string) *htmlDoc {
	data, _ := ioutil.ReadFile(path)

	if len(data) == 0 {
		return &htmlDoc{}
	}

	return &htmlDoc{
		Head: []string{innerHTML(string(data), "<head>", "</head>")},
		Body: []string{innerHTML(string(data), "<body>", "</body>")},
	}
}

func (s *htmlDoc) build(data []byte, page string) string {
	return ""
}

func setupDoc() *htmlDoc { return &htmlDoc{} }

var bundleDir string = ".orbit/dist"

var staticResourceMap map[PageRender]bool

var pageDependencies map[PageRender][]string

var serverStartupTasks = []func(){}

type PageRender string

type DocumentRenderer struct {
	fn      func(string, []byte, htmlDoc) htmlDoc
	version string
}

func NewEmptyDocumentRenderer(version string) *DocumentRenderer {
	return &DocumentRenderer{
		version: version,
		fn: func(s string, b []byte, hd htmlDoc) htmlDoc {
			return hd
		},
	}
}

var wrapDocRender = map[PageRender]*DocumentRenderer{}

type HydrationCtxKey string

const (
	OrbitManifest HydrationCtxKey = "orbitManifest"
)
