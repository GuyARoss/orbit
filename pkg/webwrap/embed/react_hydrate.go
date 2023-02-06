package webwrap

import (
	context "context"
	"fmt"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func serverRenderInnerHTML(bundleKey string, data []byte) string {
	if nodeProcess == nil {
		fmt.Println("react ssr process has not yet boot")
		return ""
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial("0.0.0.0:3024", opts...)
	if err != nil {
		return ""
	}

	defer conn.Close()
	client := NewReactRendererClient(conn)

	response, err := client.Render(context.TODO(), &RenderRequest{
		BundleID: bundleKey,
		JSONData: string(data),
	})

	if err != nil {
		return ""
	}

	return response.StaticContent
}

func reactHydrate(ctx context.Context, bundleKey string, data []byte, doc *htmlDoc) (*htmlDoc, context.Context) {
	innerServerHTML := serverRenderInnerHTML(bundleKey, data)

	if v := ctx.Value(OrbitManifest); v == nil {
		doc.Head = append(doc.Head, fmt.Sprintf(`<script id="orbit_manifest" type="application/json">%s</script>`, data))
		ctx = context.WithValue(ctx, OrbitManifest, true)
	}

	doc.Body = append(doc.Body, fmt.Sprintf(`<script class="orbit_bk" src="/p/%s.js"></script>`, bundleKey))
	copy := doc.Body

	// the doc body is adjusted +1 indices to insert the react frame at the front of the list
	// this is due to react requiring the div id to exist before the necessary javascript is loaded in
	doc.Body = make([]string, len(doc.Body)+1)
	doc.Body[0] = fmt.Sprintf(`<div id="%s_react_frame">%s</div>`, bundleKey, innerServerHTML)

	for i, c := range copy {
		doc.Body[i+1] = c
	}

	return doc, ctx
}
