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

var staticResourceMap map[string]bool

var wrapBody map[string][]string

var serverStartupTasks = []func(){}
