package webwrap

import (
	"context"
	"fmt"
	"strings"

	"github.com/GuyARoss/orbit/pkg/embedutils"
	"github.com/GuyARoss/orbit/pkg/jsparse"
)

type ReactSSR struct {
	*BaseWebWrapper
	*BaseBundler
	bundlerProcessStarted bool
	sourceMapDoc          *jsparse.DefaultJSDocument
	initDoc               *jsparse.DefaultJSDocument
	jsSwitch              *jsparse.JsDocSwitch
}

type NewReactSSROpts struct {
	SourceMapDoc *jsparse.DefaultJSDocument
	InitDoc      *jsparse.DefaultJSDocument
	Bundler      *BaseBundler
}

func NewReactSSR(opts *NewReactSSROpts) *ReactSSR {
	opts.SourceMapDoc.AddImport(&jsparse.ImportDependency{
		FinalStatement: "import React from 'react'",
		Type:           jsparse.ModuleImportType,
	})

	opts.SourceMapDoc.AddImport(&jsparse.ImportDependency{
		FinalStatement: "import ReactDOMServer from 'react-dom/server'",
		Type:           jsparse.ModuleImportType,
	})

	opts.InitDoc.AddImport(&jsparse.ImportDependency{
		Type:           jsparse.ModuleImportType,
		FinalStatement: `import * as grpc from "@grpc/grpc-js"`,
	})

	opts.InitDoc.AddImport(&jsparse.ImportDependency{
		Type:           jsparse.ModuleImportType,
		FinalStatement: `import { loadSync } from "@grpc/proto-loader"`,
	})

	jsSwitch := jsparse.NewSwitch(`BundleID`)
	fn := jsparse.NewFunc(`const buildStaticContent = ({ BundleID, JSONData }) => `, jsSwitch)

	opts.InitDoc.AddSerializable(fn)

	// @@ should we put this in a embed file?
	opts.InitDoc.AddOther(`
	const options = {
		keepCase: true,
		longs: String,
		enums: String,
		defaults: true,
		oneofs: true,
	}

	const PROTO_PATH = "./.orbit/assets/com.proto"

	var packageDefinition = loadSync(PROTO_PATH, options)
	const proto = grpc.loadPackageDefinition(packageDefinition)

	const server = new grpc.Server()

	server.addService(proto.main.ReactRenderer.service, {
		Render: ({ request }, callback) => {        
			callback(null, {
				StaticContent: buildStaticContent(request),
			})
		},
	})

	server.bindAsync(
		"0.0.0.0:50051",
		grpc.ServerCredentials.createInsecure(),
		(error, port) => {
			console.log("Server running at http://0.0.0.0:50051")
			server.start()
		}
	)
`)

	return &ReactSSR{
		sourceMapDoc: opts.SourceMapDoc,
		initDoc:      opts.InitDoc,
		BaseBundler:  opts.Bundler,
		jsSwitch:     jsSwitch,
	}
}

func (r *ReactSSR) RequiredBodyDOMElements(context.Context, *CacheDOMOpts) []string {
	return []string{}
}

// @@ add support for multiple bundled resources?
func (r *ReactSSR) Setup(ctx context.Context, settings *BundleOpts) ([]*BundledResource, error) {
	bundleFilePath := fmt.Sprintf("%s/%s.js", r.PageOutputDir, settings.BundleKey)
	r.sourceMapDoc.AddImport(&jsparse.ImportDependency{
		FinalStatement: fmt.Sprintf("import {%s} from '%s'", settings.Name, fmt.Sprintf("./%s", settings.BundleKey)),
		Type:           jsparse.LocalImportType,
	})

	r.sourceMapDoc.AddOther(fmt.Sprintf(`export const %s = (d) => ReactDOMServer.renderToString(<%s {...d}/>)`, strings.ToLower(settings.Name), settings.Name))
	r.initDoc.AddImport(&jsparse.ImportDependency{
		FinalStatement: fmt.Sprintf("import %s from '%s'", strings.ToLower(settings.Name), fmt.Sprintf("./%s", "react_ssr.map.js")),
		Type:           jsparse.LocalImportType,
	})

	r.jsSwitch.Add(jsparse.JSString, settings.BundleKey, fmt.Sprintf(`return %s(JSONData)`, strings.ToLower(settings.Name)))

	return []*BundledResource{
		{BundleFilePath: bundleFilePath,
			ConfiguratorFilePath: fmt.Sprintf("%s/react_ssr.map.js", r.PageOutputDir),
			ConfiguratorPage:     r.sourceMapDoc},
		{BundleFilePath: bundleFilePath,
			ConfiguratorFilePath: fmt.Sprintf("%s/react_ssr.js", r.PageOutputDir),
			ConfiguratorPage:     r.initDoc},
	}, nil
}
func (r *ReactSSR) Apply(doc jsparse.JSDocument) (jsparse.JSDocument, error) {
	doc.AddOther(fmt.Sprintf(
		"export default %s",
		doc.Name()),
	)

	return doc, nil
}

func (r *ReactSSR) DoesSatisfyConstraints(fileExtension string) bool {
	return strings.Contains(fileExtension, reactExtension)
}

func (r *ReactSSR) Version() string {
	return "reactSSR"
}

var reactSSRParentDocument = jsparse.NewImportDocument(&jsparse.ImportDependency{
	FinalStatement: "import React from 'react'",
	Type:           jsparse.ModuleImportType,
}, &jsparse.ImportDependency{
	FinalStatement: "import ReactDOMServer from 'react-dom/server'",
	Type:           jsparse.ModuleImportType,
})

func (r *ReactSSR) Bundle(configuratorFilePath string) error {
	return nil
}

func (r *ReactSSR) HydrationFile() []embedutils.FileReader {
	files, err := embedFiles.ReadDir("embed")
	if err != nil {
		return nil
	}

	u := []embedutils.FileReader{}

	for _, file := range files {
		if strings.Contains(file.Name(), "react_ssr.go") {
			u = append(u, &embedFileReader{fileName: file.Name()})
			continue
		}
		if strings.Contains(file.Name(), "pb.go") {
			u = append(u, &embedFileReader{fileName: file.Name()})
		}
	}
	return u
}
