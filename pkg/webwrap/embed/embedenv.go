package webwrap

// htmlDoc represents a basic document model that will be rendered upon build request
type htmlDoc struct {
	Head []string
	Body []string
}

func (s *htmlDoc) build(data []byte, page string) string {
	return ""
}

func setupDoc() *htmlDoc { return &htmlDoc{} }

var bundleDir string = ".orbit/dist"

var staticResourceMap map[PageRender]bool

var wrapBody map[PageRender][]string

var serverStartupTasks = []func(){}

type PageRender string

type DocumentRenderer struct {
	fn      func(string, []byte, htmlDoc) htmlDoc
	version string
}

var wrapDocRender = map[PageRender]DocumentRenderer{}
