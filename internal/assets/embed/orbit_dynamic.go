package orbit

import "fmt"

var bundleDir string = ".orbit/dist"

func deleteMeThing(bundleKey string, data []byte, doc htmlDoc) htmlDoc {
	doc.Head = append(doc.Head, fmt.Sprintf(`<script id="orbit_manifest" type="application/json">%s</script>`, data))
	doc.Body = append(doc.Body, fmt.Sprintf(`<script id="orbit_bk" src="/p/%s.js">`, bundleKey))

	return doc
}

var wrapDocRender = map[PageRender][]func(string, []byte, htmlDoc) htmlDoc{
	"test": {deleteMeThing},
}

var wrapBody = map[PageRender][]string{}

type PageRender string

var publicDir string = "./public/index.html"

type BundleMode int32

const (
	DevBundleMode  BundleMode = 0
	ProdBundleMode BundleMode = 1
)

var CurrentDevMode BundleMode

var hotReloadPort = 1000
