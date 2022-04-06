package webwrap

import (
	context "context"
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
