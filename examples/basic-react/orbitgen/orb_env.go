package orbitgen

import (
	"context"
	"fmt"
)


func reactManifestFallback(ctx context.Context, bundleKey string, data []byte, doc *htmlDoc) (*htmlDoc, context.Context) {
	if v := ctx.Value(OrbitManifest); v == nil {
		doc.Head = append(doc.Head, fmt.Sprintf(`<script id="orbit_manifest" type="application/json">%s</script>`, data))
		ctx = context.WithValue(ctx, OrbitManifest, true)
	}

	doc.Body = append(doc.Body, fmt.Sprintf(`<script class="orbit_bk" src="/p/%s.js"></script>`, bundleKey))

	return doc, ctx
}
var staticResourceMap = map[PageRender]bool{
	ExampleTwoPage: true,
	ExamplePage: false,
}
var serverStartupTasks = []func(){}
var wrapDocRender = map[PageRender]*DocumentRenderer{
	ExampleTwoPage: {fn: reactManifestFallback, version: "reactManifestFallback"},
	ExamplePage: {fn: reactManifestFallback, version: "reactManifestFallback"},
}

type DocumentRenderer struct {
	fn func(context.Context, string, []byte, *htmlDoc) (*htmlDoc, context.Context)
	version string
}
var reactManifestFallback_bodywrap = []string{
`<script src="/p/02bab3977c197c77b270370f110270b1.js"></script>`,
`<script src="/p/8cfc2b31824016492ec09fc306264efd.js"></script>`,
`<div id="b647b6d4-da06-49a3-961a-4fe69b34eb1d"></div>`,
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
	ExampleTwoPage: reactManifestFallback_bodywrap,
	ExamplePage: reactManifestFallback_bodywrap,
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