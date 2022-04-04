package orbitgen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)


// Request is the standard request payload for the orbit page handler
// this is just a fancy wrapper around the http request & response that will also assist
// the rendering of bundled pages & incoming path slugs
type Request struct {
	RenderPage  func(page PageRender, data interface{})
	RenderPages func(data interface{}, pages ...PageRender)
	Request     *http.Request
	Response    http.ResponseWriter
	Slugs       map[string]string
}

// DefaultPage defines the standard behavior for a orbit page handler
type DefaultPage interface {
	Handle(*Request)
}

// htmlDoc represents a basic document model that will be rendered upon build request
type htmlDoc struct {
	Head []string
	Body []string
}

// build builds the htmldocument given data for orbits manifest and the page's
// javascript bundle key to render the document out to a single string
func (s *htmlDoc) build(data []byte, pages ...PageRender) string {
	body := make([]string, 0)
	head := make([]string, 0)
	isWrapped := make(map[PageRender]bool)

	for _, p := range pages {
		if !isWrapped[p] {
			for _, b := range wrapBody[p] {
				head = append(head, b)
			}
			isWrapped[p] = true
		}

		// if the page is of static origin, we first check to see if it exists on the file system
		// if it does, it will be applied to the current html document, rather than returned directly
		// this is to support the usage of static html within micro-frontends
		if staticResourceMap[p] {
			_, err := os.Stat(publicDir)
			if !os.IsNotExist(err) {
				f, err := ioutil.ReadFile(fmt.Sprintf("%s%c%s", http.Dir(bundleDir), os.PathSeparator, p))
				if err != nil {
					continue
				}

				if len(pages) == 1 {
					return string(f)
				}

				body = append(body, innerHTML(string(f), "<body>", "</body>"))
				head = append(head, innerHTML(string(f), "<head>", "</head>"))
			}
		}
	}

	for _, p := range pages {
		op := wrapDocRender[p]

		if op == nil {
			continue
		}

		doc := op.fn(string(p), data, *s)
		for _, b := range doc.Body {
			body = append(body, b)
		}

		for _, h := range doc.Head {
			head = append(head, h)
		}
	}

	return fmt.Sprintf(`
	<!doctype html>
	<html lang="en">
	<head>%s</head>
	<body>%s</body>
	</html>`, strings.Join(head, ""), strings.Join(body, ""))
}

// innerHTML is a utility function that assists with the parsing the content of html tags
// it does this by returning the subset of the two provided strings "subStart" & "subEnd"
func innerHTML(str string, subStart string, subEnd string) string {
	return strings.Split(strings.Join(strings.Split(str, subStart)[1:], ""), subEnd)[0]
}

// defaultHTMLDoc builds a standard html doc for orbit that also verifies the public directory
// if override data exits, then it will use that as a base for the HTML document
func defaultHTMLDoc(override string) *htmlDoc {
	base := &htmlDoc{Head: []string{`<meta charset="utf-8" />`}, Body: []string{}}

	// we allow some special operations on the dom for debugging, currently supporting:
	// - geting the contents of orbit manifest with the function "getManifest"
	if CurrentDevMode == DevBundleMode {
		base.Body = append(base.Body, `<script class="debug"> const getManifest = () => JSON.parse(document.getElementById("orbit_manifest").textContent) </script>`)
		base.Body = append(base.Body, `<script class="debug" src="/p/hotreload.js"> </script>`)
		base.Body = append(base.Body, fmt.Sprintf(`<script class="debug" id="debug_data" type="application/json">{ "hotReloadPort": %d }</script>`, hotReloadPort))
	}

	// the html override that will provide a basis for the default html doc
	if override != "" {
		base.Body = append(base.Body, innerHTML(string(override), "<body>", "</body>"))
		base.Head = append(base.Head, innerHTML(string(override), "<head>", "</head>"))
	}

	return base
}

// parseSlug will parse slugs from the incoming path provided initial slugKeys and return a map of the slugs
// in map[string]string form where the key is the slug name & value is the value as it appears in the path
func parseSlug(slugKeys map[int]string, path string) map[string]string {
	slugs := make(map[string]string)
	if len(slugKeys) > 0 {
		paths := strings.Split(path, "/")

		for idx, p := range paths {
			key := slugKeys[idx]
			if key != "" {
				slugs[key] = p
			}
		}
	}

	return slugs
}

// parsePathSlugs will parse initial slugs found in a path, these slugs can be identified with
// the "{" char prepended & "}" appended to the path/subpath e.g "/path/{thing}" will represent "thing" as a slug key.
func parsePathSlugs(path string) (map[int]string, string) {
	slugKeys := make(map[int]string)

	validInitial := make([]string, 0)
	slide := true
	if strings.Contains(path, "{") {
		paths := strings.Split(path, "/")
		for idx, p := range paths {
			if strings.Contains(p, "{") {
				slide = false
				slugKeys[idx] = p[1 : len(p)-1]
			}

			if slide {
				validInitial = append(validInitial, p)
			}
		}
	}

	finalPath := path
	if len(slugKeys) > 0 {
		finalPath = fmt.Sprintf("%s/", strings.Join(validInitial, "/"))
	}

	return slugKeys, finalPath
}

// muxHandle is used to inject the base mux handler behavior
type MuxHandler interface {
	HandleFunc(string, func(http.ResponseWriter, *http.Request))
	Handle(string, http.Handler)
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type Serve struct {
	mux MuxHandler
	doc *htmlDoc
}

// HandleFunc attaches a handler to the specified pattern, this handler will be
// ran upon a match of the request path during an incoming http request.
func (s *Serve) HandleFunc(path string, handler func(c *Request)) {
	slugs, path := parsePathSlugs(path)

	s.mux.HandleFunc(path, func(rw http.ResponseWriter, r *http.Request) {
		requestSlugs := make(map[string]string)

		if len(slugs) > 0 {
			requestSlugs = parseSlug(slugs, r.URL.Path)

			if len(requestSlugs) != len(slugs) {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		renderPage := func(page PageRender, data interface{}) {
			d, err := json.Marshal(data)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			html := s.doc.build(d, page)

			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(html))
		}

		renderPages := func(data interface{}, pages ...PageRender) {
			d, err := json.Marshal(data)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			html := s.doc.build(d, pages...)

			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(html))
		}

		ctx := &Request{
			RenderPage:  renderPage,
			RenderPages: renderPages,
			Request:     r,
			Response:    rw,
			Slugs:       requestSlugs,
		}

		handler(ctx)
	})
}

// HandlePage attaches an orbit page to the specified pattern, this handler will be
// ran upon a match of the request path during an incoming http request
func (s *Serve) HandlePage(path string, dp DefaultPage) {
	s.HandleFunc(path, dp.Handle)
}

// setupMuxRequirements creates the required mux handlers for orbit, these include
// - fileserver for the bundle directory bound to the "/p/" directory
func (s *Serve) setupMuxRequirements() *Serve {
	s.mux.Handle("/p/", http.StripPrefix("/p/", http.FileServer(http.Dir(bundleDir))))

	return s
}

// Serve returns the mux server
func (s *Serve) Serve() MuxHandler {
	return s.mux
}

func setupDoc() *htmlDoc {
	html := ""

	_, err := os.Stat(publicDir)
	if !os.IsNotExist(err) {
		// im not entirely sure that an error here would warrant a change in behavior
		// invalid files should already be skipped, besides that, an empty []byte should suffice.
		data, _ := ioutil.ReadFile(publicDir)
		html = string(data)
	}

	return defaultHTMLDoc(html)
}

func New() (*Serve, error) {
	for _, sfn := range serverStartupTasks {
		sfn()
	}

	return (&Serve{
		mux: http.NewServeMux(),
		doc: setupDoc(),
	}).setupMuxRequirements(), nil
}
