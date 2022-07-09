package orbitgen

import (
	"context"
	"fmt"
)
func javascriptWebpack(ctx context.Context, bundleKey string, data []byte, doc *htmlDoc) (*htmlDoc, context.Context) {
	if v := ctx.Value(OrbitManifest); v == nil {
		doc.Head = append(doc.Head, fmt.Sprintf(`<script id="orbit_manifest" type="application/json">%s</script>`, data))
		ctx = context.WithValue(ctx, OrbitManifest, true)
	}
	doc.Body = append(doc.Body, fmt.Sprintf(`<script class="orbit_bk" src="/p/%s.js"></script>`, bundleKey))
	return doc, ctx
}
func reactManifestFallback(ctx context.Context, bundleKey string, data []byte, doc *htmlDoc) (*htmlDoc, context.Context) {
	if v := ctx.Value(OrbitManifest); v == nil {
		doc.Head = append(doc.Head, fmt.Sprintf(`<script id="orbit_manifest" type="application/json">%s</script>`, data))
		ctx = context.WithValue(ctx, OrbitManifest, true)
	}
	doc.Body = append(doc.Body, fmt.Sprintf(`<script class="orbit_bk" src="/p/%s.js"></script>`, bundleKey))
	copy := doc.Body
	doc.Body = make([]string, len(doc.Body)+1)
	doc.Body[0] = fmt.Sprintf(`<div id="%s_react_frame"></div>`, bundleKey)
	for i, c := range copy {
		doc.Body[i+1] = c
	}
	return doc, ctx
}
var staticResourceMap = map[PageRender]bool{
	ThingPage: false,
	StaticPage: true,
	NamePage: false,
	AgePage: false,
}
var serverStartupTasks = []func(){}
var wrapDocRender = map[PageRender]*DocumentRenderer{
	ThingPage: {fn: reactManifestFallback, version: "reactManifestFallback"},
	StaticPage: {fn: javascriptWebpack, version: "javascriptWebpack"},
	NamePage: {fn: reactManifestFallback, version: "reactManifestFallback"},
	AgePage: {fn: javascriptWebpack, version: "javascriptWebpack"},
}

type DocumentRenderer struct {
	fn func(context.Context, string, []byte, *htmlDoc) (*htmlDoc, context.Context)
	version string
}
var javascriptWebpack_bodywrap = []string{
`<script> const onLoadTasks = []; window.onload = (e) => { onLoadTasks.forEach(t => t(e))} </script>`,
}

var reactManifestFallback_bodywrap = []string{
`<script src="/p/02bab3977c197c77b270370f110270b1.js"></script>`,
`<script src="/p/8cfc2b31824016492ec09fc306264efd.js"></script>`,
}

var bundleDir string = ".orbit/dist"

var publicDir string = "./public/index.html"
var hotReloadPort int = 0
type PageRender string

const ( 
	// orbit:page .//pages/thing.jsx
	ThingPage PageRender = "79e586fbcad45ddab385257f9f3d3eaf"
	// orbit:page .//pages/static.js
	StaticPage PageRender = "a33d65b63e235f0788c046da83f123c2"
	// orbit:page .//pages/name.jsx
	NamePage PageRender = "d3204a628de15bc7929ef30743f5ff2a"
	// orbit:page .//pages/age.js
	AgePage PageRender = "752668a8ac895cdea34ec499148eaa8b"
)

var pageDependencies = map[PageRender][]string{
	ThingPage: reactManifestFallback_bodywrap,
	StaticPage: javascriptWebpack_bodywrap,
	NamePage: reactManifestFallback_bodywrap,
	AgePage: javascriptWebpack_bodywrap,
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

var CurrentDevMode BundleMode = DevBundleMode