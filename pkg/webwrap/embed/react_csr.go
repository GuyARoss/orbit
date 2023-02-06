package webwrap

import (
	context "context"
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
