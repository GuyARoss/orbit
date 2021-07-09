package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var bundleDir string = "../.orbit/dist"

type PageRender string

const (
	HelloPage PageRender = "af9e40e28f955222b69dfe0f1b729754"
)

type RuntimeCtx struct {
	RenderPage func(page PageRender, data interface{})
}

type DefaultPage interface {
	Render(c *RuntimeCtx)
	// @@todo: GET & POST
}

type Route struct {
	Path string
	Page DefaultPage
}

func HandlePage(path string, dp DefaultPage) {
	http.HandleFunc(path, func(rw http.ResponseWriter, r *http.Request) {
		renderPage := func(page PageRender, data interface{}) {
			d, err := json.Marshal(data)
			if err != nil {
				// @ do something
				return
			}

			// @@todo: look into embeding this
			html := fmt.Sprintf(`
			<!doctype html><html lang="en"><head><meta charset="utf-8"><script id="orbit_manifest" type="application/json">%s</script></head>
			<body><script src="https://unpkg.com/react/umd/react.production.min.js" crossorigin></script><script src="https://unpkg.com/react-dom/umd/react-dom.production.min.js" crossorigin></script><script src="https://unpkg.com/react-bootstrap@next/dist/react-bootstrap.min.js" crossorigin></script>			
			<div id="root"></div><script src="/p/%s.js"></script>				
			</body></html>
			`, string(d), page)

			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(html))
		}

		dp.Render(&RuntimeCtx{
			RenderPage: renderPage,
		})
	})
}

func Start(port int) {
	http.Handle("/p/", http.StripPrefix("/p/", http.FileServer(http.Dir(bundleDir))))
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

type Test123 struct {
}

func (t *Test123) Render(c *RuntimeCtx) {
	d := make(map[string]interface{})
	d["name"] = "guy"
	d["age"] = 22

	c.RenderPage(HelloPage, d)
}

func main() {
	HandlePage("/test", &Test123{})

	Start(3001)
}
