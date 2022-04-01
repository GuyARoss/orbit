package webwrap

import "fmt"

func reactManifestFallback(bundleKey string, data []byte, doc htmlDoc) htmlDoc {
	// the "orbit_manifest" refers to the object content that the specified
	// web javascript bundle can make use of
	doc.Head = append(doc.Head, fmt.Sprintf(`<script id="orbit_manifest" type="application/json">%s</script>`, data))
	doc.Body = append(doc.Body, fmt.Sprintf(`<script id="orbit_bk" src="/p/%s.js"></script>`, bundleKey))

	return doc
}
