package orbitgen

import (
	"context"
)
var staticResourceMap = map[PageRender]bool{
}
var serverStartupTasks = []func(){}
var wrapDocRender = map[PageRender]*DocumentRenderer{
}

type DocumentRenderer struct {
	fn func(context.Context, string, []byte, *htmlDoc) (*htmlDoc, context.Context)
	version string
}
var bundleDir string = ".orbit/dist"

var publicDir string = "./public/index.html"
var hotReloadPort int = 0
type PageRender string


var pageDependencies = map[PageRender][]string{
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