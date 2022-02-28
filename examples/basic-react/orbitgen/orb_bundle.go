package orbitgen


var react = []string{
`<script src="/p/fc38086145547d465be97fec2e412a16.js"></script>`,
`<script src="/p/fc38086145547d465be97fec2e412a16.js"></script>`,
`<div id="root"></div>`,
}

var bundleDir string = ".orbit/dist"

var publicDir string = "./public/index.html"
type PageRender string

const ( 
	ExamplePage PageRender = "496a05464c3f5aa89e1d8bed7afe59d4"
	ExampleTwoPage PageRender = "fe9faa2750e8559c8c213c2c25c4ce73"
)

var wrapBody = map[PageRender][]string{
	ExamplePage: react,
	ExampleTwoPage: react,
}

type BundleMode int32

const (
	DevBundleMode  BundleMode = 0
	ProdBundleMode BundleMode = 1
)

var CurrentDevMode BundleMode = ProdBundleMode