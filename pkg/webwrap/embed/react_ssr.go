package webwrap

import (
	context "context"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func reactSSR(bundleKey string, data []byte, doc htmlDoc) htmlDoc {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// @@todo(guy): prefer unix domian socket here
	conn, err := grpc.Dial("0.0.0.0:50051", opts...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := NewReactRendererClient(conn)

	response, err := client.Render(context.TODO(), &RenderRequest{
		BundleID: bundleKey,
		JSONData: string(data),
	})

	if err != nil {
		// @@todo(guy): this should not be handled here
		return htmlDoc{
			Body: []string{"ssr failed to resolve"},
		}
	}

	doc.Body = append(doc.Body, response.StaticContent)

	return doc
}
