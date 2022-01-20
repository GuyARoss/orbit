package orbit

// **__START_STATIC__**
import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// **__END_STATIC__**

var bundleDir string = ".orbit/dist"

type PageRender string

var hotReloadPipePath string = ""

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

func HandlePage(path string, dp DefaultPage) {
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

	http.HandleFunc(path, func(rw http.ResponseWriter, r *http.Request) {
		renderPage := func(page PageRender, data interface{}) {
			d, err := json.Marshal(data)
			if err != nil {
				// @@todo(debug): do something
				return
			}

			html := fmt.Sprintf(`
			<!doctype html><html lang="en"><head><meta charset="utf-8"><script id="orbit_manifest" type="application/json">%s</script></head>
			<body><script src="https://unpkg.com/react/umd/react.production.min.js" crossorigin></script><script src="https://unpkg.com/react-dom/umd/react-dom.production.min.js" crossorigin></script><script src="https://unpkg.com/react-bootstrap@next/dist/react-bootstrap.min.js" crossorigin></script>			
			<div id="root"></div><script src="/p/%s.js"></script>				
			</body></html>
			`, string(d), page)

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

		dp.Handle(ctx)
	})
}

func Start(port int) {
	http.Handle("/p/", http.StripPrefix("/p/", http.FileServer(http.Dir(bundleDir))))
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

// **__END_STATIC__**
