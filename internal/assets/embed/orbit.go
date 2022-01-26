package orbit

// **__START_STATIC__**
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// **__END_STATIC__**

var bundleDir string = ".orbit/dist"

type PageRender string

var hotReloadPipePath string = ""

var publicDir string = "./public/index.html"

// **__START_STATIC__**
type RuntimeCtx struct {
	RenderPage func(page PageRender, data interface{})
	Request    *http.Request
	Response   http.ResponseWriter
	Slugs      map[string]string
}

type DefaultPage interface {
	Handle(c *RuntimeCtx)
}

type Route struct {
	Path string
	Page DefaultPage
}

type htmlDoc struct {
	Head []string
	Body []string
}

func centerStr(str string, subStart string, subEnd string) string {
	init := strings.Split(str, subStart)[1:]

	return strings.Split(strings.Join(init, ""), subEnd)[0]
}

func (s *htmlDoc) build(data []byte, page PageRender) string {
	return fmt.Sprintf(`
	<!doctype html>
	<html lang="en">
			<head>
				%s	
				<script id="orbit_manifest" type="application/json">
				%s
				</script>
			</head>
		<body>
			%s
			<script src="/p/%s.js"></script>				
		</body>
	</html>
	`, strings.Join(s.Head, ""), string(data), strings.Join(s.Body, ""), page)
}

func initHtmlDoc() (*htmlDoc, error) {
	base := &htmlDoc{
		Head: []string{`<meta charset="utf-8" />`},
		Body: []string{
			`<script src="https://unpkg.com/react/umd/react.production.min.js" crossorigin></script><script src="https://unpkg.com/react-dom/umd/react-dom.production.min.js" crossorigin></script><script src="https://unpkg.com/react-bootstrap@next/dist/react-bootstrap.min.js" crossorigin></script>`,
			`<div id="root"></div>`,
		},
	}

	_, err := os.Stat(publicDir)
	if !os.IsNotExist(err) {
		data, err := ioutil.ReadFile(publicDir)
		if err != nil {
			return base, err
		}

		base.Body = append(base.Body, centerStr(string(data), "<body>", "</body>"))
		base.Head = append(base.Body, centerStr(string(data), "<head>", "</head>"))
	}

	return base, nil
}

func HandleFunc(path string, handler func(c *RuntimeCtx)) {
	slugKeys := make(map[int]string, 0)

	validInitial := make([]string, 0)
	stillValid := true
	if strings.Contains(path, "{") {
		paths := strings.Split(path, "/")
		for idx, p := range paths {
			if strings.Contains(p, "{") {
				stillValid = false
				slugKeys[idx] = p[1 : len(p)-1]
			}

			if stillValid {
				validInitial = append(validInitial, p)
			}
		}
		path = fmt.Sprintf("%s/", strings.Join(validInitial, "/"))
	}

	doc, err := initHtmlDoc()
	if err != nil {
		return
	}

	http.HandleFunc(path, func(rw http.ResponseWriter, r *http.Request) {
		renderPage := func(page PageRender, data interface{}) {
			d, err := json.Marshal(data)
			if err != nil {
				// @@todo(debug): do something
				return
			}

			html := doc.build(d, page)

			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(html))
		}

		slugs := make(map[string]string)
		if len(slugKeys) > 0 {
			paths := strings.Split(r.URL.Path, "/")
			for idx, p := range paths {
				key := slugKeys[idx]
				if key != "" {
					slugs[key] = p
				}
			}
		}

		ctx := &RuntimeCtx{
			RenderPage: renderPage,
			Request:    r,
			Response:   rw,
			Slugs:      slugs,
		}

		handler(ctx)
	})
}

func HandlePage(path string, dp DefaultPage) {
	HandleFunc(path, dp.Handle)
}

func Start(port int) {
	http.Handle("/p/", http.StripPrefix("/p/", http.FileServer(http.Dir(bundleDir))))
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

// **__END_STATIC__**
