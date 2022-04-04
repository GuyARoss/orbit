package orbitgen

import (
	"fmt"
)


func javascriptWebpack(bundleKey string, data []byte, doc htmlDoc) htmlDoc {
	doc.Head = append(doc.Head, fmt.Sprintf(`<script id="orbit_manifest" type="application/json">%s</script>`, data))
	doc.Head = append(doc.Head, `<script> const onLoadTasks = []; window.onload = (e) => { onLoadTasks.forEach(t => t(e))} </script>`)

	doc.Body = append(doc.Body, fmt.Sprintf(`<script id="orbit_bk" src="/p/%s.js"></script>`, bundleKey))

	return doc
}


func reactManifestFallback(bundleKey string, data []byte, doc htmlDoc) htmlDoc {
	// the "orbit_manifest" refers to the object content that the specified
	// web javascript bundle can make use of
	doc.Head = append(doc.Head, fmt.Sprintf(`<script id="orbit_manifest" type="application/json">%s</script>`, data))
	doc.Body = append(doc.Body, fmt.Sprintf(`<script id="orbit_bk" src="/p/%s.js"></script>`, bundleKey))

	return doc
}
var staticResourceMap = map[PageRender]bool{
	StaticPage: true,
	AgePage: false,
	NamePage: false,
}
var serverStartupTasks = []func(){}
var wrapDocRender = map[PageRender]DocumentRenderer{
	StaticPage: {fn: javascriptWebpack, version: "javascriptWebpack"},
	AgePage: {fn: javascriptWebpack, version: "javascriptWebpack"},
	NamePage: {fn: reactManifestFallback, version: "reactManifestFallback"},
}

type DocumentRenderer struct {
	fn func(string, []byte, htmlDoc) htmlDoc
	version string
}
var javascriptWebpack_bodywrap = []string{
}

var reactManifestFallback_bodywrap = []string{
`<script src="/p/02bab3977c197c77b270370f110270b1.js"></script>`,
`<script src="/p/8cfc2b31824016492ec09fc306264efd.js"></script>`,
`<div id="2dbd580d-e418-4bc1-9b8c-363f9608117d"></div>`,
}

var bundleDir string = ".orbit/dist"

var publicDir string = "./public/index.html"
var hotReloadPort int = 0
type PageRender string

const ( 
	// orbit:page .//pages/static.js
	StaticPage PageRender = "a33d65b63e235f0788c046da83f123c2"
	// orbit:page .//pages/age.js
	AgePage PageRender = "752668a8ac895cdea34ec499148eaa8b"
	// orbit:page .//pages/name.jsx
	NamePage PageRender = "d3204a628de15bc7929ef30743f5ff2a"
)

var wrapBody = map[PageRender][]string{
	StaticPage: javascriptWebpack_bodywrap,
	AgePage: javascriptWebpack_bodywrap,
	NamePage: reactManifestFallback_bodywrap,
}

type BundleMode int32

const (
	DevBundleMode  BundleMode = 0
	ProdBundleMode BundleMode = 1
)

var CurrentDevMode BundleMode = DevBundleMode