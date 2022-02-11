package orbitgen


var bundleDir string = ".orbit/dist"


var publicDir string = "./public/index.html"

type PageRender string

const ( 
	ExamplePage PageRender = "496a05464c3f5aa89e1d8bed7afe59d4"
)


type BundleMode int32

const (
	DevBundleMode  BundleMode = 0
	ProdBundleMode BundleMode = 1
)

var CurrentDevMode BundleMode = DevBundleMode