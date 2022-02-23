package orbit

var bundleDir string = ".orbit/dist"

type PageRender string

var publicDir string = "./public/index.html"

type BundleMode int32

const (
	DevBundleMode  BundleMode = 0
	ProdBundleMode BundleMode = 1
)

var CurrentDevMode BundleMode

var wrapBody = map[PageRender][]string{}
