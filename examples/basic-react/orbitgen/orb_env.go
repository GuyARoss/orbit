package orbitgen

import (
	"context"
	"fmt"
)
func reactCSR(ctx context.Context, bundleKey string, data []byte, doc *htmlDoc) (*htmlDoc, context.Context) {
	if v := ctx.Value(OrbitManifest); v == nil {
		doc.Head = append(doc.Head, fmt.Sprintf(`<script id="orbit_manifest" type="application/json">%s</script>`, data))
		ctx = context.WithValue(ctx, OrbitManifest, true)
	}
	doc.Body = append(doc.Body, fmt.Sprintf(`<script class="orbit_bk" src="/p/%s.js"></script>`, bundleKey))
	copy := doc.Body
	// the doc body is adjusted +1 indices to insert the react frame at the front of the list
	// this is due to react requiring the div id to exist before the necessary javascript is loaded in
	doc.Body = make([]string, len(doc.Body)+1)
	doc.Body[0] = fmt.Sprintf(`<div id="%s_react_frame"></div>`, bundleKey)
	for i, c := range copy {
		doc.Body[i+1] = c
	}
	return doc, ctx
}
var staticResourceMap = map[PageRender]bool{
	ExampleTwoPage: true,
	ExamplePage: false,
}
var serverStartupTasks = []func(){}
type RenderFunction func(context.Context, string, []byte, *htmlDoc) (*htmlDoc, context.Context)
var wrapDocRender = map[PageRender]*DocumentRenderer{
	ExampleTwoPage: {fn: reactCSR, version: "reactCSR"},
	ExamplePage: {fn: reactCSR, version: "reactCSR"},
}

type DocumentRenderer struct {
	fn RenderFunction
	version string
}
var bundleDir string = ".orbit/dist"

var publicDir string = "./public/index.html"
var hotReloadPort int = 0
type PageRender string

const ( 
	// orbit:page .//pages/example2.jsx
	ExampleTwoPage PageRender = "fe9faa2750e8559c8c213c2c25c4ce73"
	// orbit:page .//pages/example.jsx
	ExamplePage PageRender = "496a05464c3f5aa89e1d8bed7afe59d4"
)

var pageDependencies = map[PageRender][]string{
	ExampleTwoPage: {`<script src="/p/fc38086145547d465be97fec2e412a16.js"></script>`,
`<script src="/p/a63649d90703a7b09f22aed8d310be5b.js"></script>`,
},
	ExamplePage: {`<script src="/p/fc38086145547d465be97fec2e412a16.js"></script>`,
`<script src="/p/a63649d90703a7b09f22aed8d310be5b.js"></script>`,
},
}

	
type HydrationCtxKey string

const (
	OrbitManifest HydrationCtxKey = "orbitManifest"
)

type BundleMode int32

const (
	DevBundleMode  BundleMode = 0
	ProdBundleMode BundleMode = 1
)

var CurrentDevMode BundleMode = ProdBundleMode
var routeTable = map[PageRender]string{
ExampleTwoPage: "/second",}