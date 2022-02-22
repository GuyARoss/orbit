package orbitgen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

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
	body := append(s.Body, wrapBody[page]...)

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
	`, strings.Join(s.Head, ""), string(data), strings.Join(body, ""), page)
}

func initHtmlDoc() (*htmlDoc, error) {
	base := &htmlDoc{
		Head: []string{`<meta charset="utf-8" />`},
		Body: []string{},
	}

	if CurrentDevMode == DevBundleMode {
		base.Body = append(base.Body, `<script> const getProps = () => JSON.parse(document.getElementById("orbit_manifest").textContent) </script>`)
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
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			html := doc.build(d, page)

			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(html))
		}

		ctx := &RuntimeCtx{
			RenderPage: renderPage,
			Request:    r,
			Response:   rw,
			Slugs:      parseSlug(slugKeys, r.URL.Path),
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

