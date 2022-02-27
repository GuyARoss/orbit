package orbitgen


var reactv16_13_1 = []string{
`<script src="https://unpkg.com/react/umd/react.production.min.js" crossorigin></script><script src="https://unpkg.com/react-dom/umd/react-dom.production.min.js" crossorigin></script><script src="https://unpkg.com/react-bootstrap@next/dist/react-bootstrap.min.js" crossorigin></script>`,
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
	ExamplePage: reactv16_13_1,
	ExampleTwoPage: reactv16_13_1,
}

type BundleMode int32

const (
	DevBundleMode  BundleMode = 0
	ProdBundleMode BundleMode = 1
)

var CurrentDevMode BundleMode = DevBundleMode