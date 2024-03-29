package orbitgen

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)
func TestInnerHTML(t *testing.T) {
	var tt = []struct {
		f string
		s string
		t string
		e string
	}{
		{"<thing>DATA</blah>", "<thing>", "</blah>", "DATA"},
	}
	for _, d := range tt {
		c := innerHTML(d.f, d.s, d.t)
		if !strings.Contains(c, d.e) {
			t.Errorf("expected %s got %s", d.e, c)
		}
	}
}
func TestParseSlug(t *testing.T) {
	var tt = []struct {
		path   string
		k      map[int]string
		s      string
		expect string
	}{
		{"/thing/test", map[int]string{1: "slug"}, "slug", "thing"},
	}
	for _, d := range tt {
		s := parseSlug(d.k, d.path)
		got := s[d.s]
		if got != d.expect {
			t.Errorf("expected %s got %s", d.expect, got)
		}
	}
}
func TestParsePathSlugs(t *testing.T) {
	tt := []struct {
		path  string
		epath string
		i     int
		e     string
	}{
		{"/thing/{toast}", "/thing/", 2, "toast"},
		{"/thing", "/thing", 0, ""},
		{"/thing/{toast}/{cat}", "/thing/", 2, "toast"},
	}
	for i, d := range tt {
		slugs, path := parsePathSlugs(d.path)
		if slugs[d.i] != d.e {
			t.Errorf("(%d) expected %s got %s", i, d.e, slugs[d.i])
		}
		if path != d.epath {
			t.Errorf("(%d) expected %s got %s", i, d.epath, path)
		}
	}
}
func TestDefaultHTMLDoc_NoPath(t *testing.T) {
	currentMode := CurrentDevMode
	CurrentDevMode = DevBundleMode
	bodyContent := "<test>This is a thing here<test>"
	headContent := `<script id="thing"> {} </script>`
	doc := defaultHTMLDoc(fmt.Sprintf("<head>%s</head><body>%s</body>", headContent, bodyContent))
	hasMeta := false
	hasHeadContent := false
	for _, h := range doc.Head {
		if strings.Contains(h, "meta") {
			hasMeta = true
		}
		if strings.Contains(h, headContent) {
			hasHeadContent = true
		}
	}
	if !hasMeta {
		t.Error("default doc does not contain meta tag")
	}
	if !hasHeadContent {
		t.Error("expected override head to exist in final body")
	}
	hasNewBody := false
	hasDebugClass := false
	for _, h := range doc.Body {
		if strings.Contains(h, bodyContent) {
			hasNewBody = true
		}
		if strings.Contains(h, `class="debug"`) {
			hasDebugClass = true
		}
	}
	if !hasNewBody {
		t.Error("expected override body to exist in final body")
	}
	if !hasDebugClass {
		t.Error("expected debug class to be present during debug mode")
	}
	CurrentDevMode = currentMode
}
type mockResponseWriter struct {
	mockWriteHeader func(statusCode int)
	mockWrite       func([]byte) (int, error)
	mockHeader      func() http.Header
}
func (s *mockResponseWriter) Header() http.Header          { return s.mockHeader() }
func (s *mockResponseWriter) Write(in []byte) (int, error) { return s.mockWrite(in) }
func (s *mockResponseWriter) WriteHeader(statusCode int)   { s.mockWriteHeader(statusCode) }
type mockHandle struct {
	isErr         bool
	writer        http.ResponseWriter
	checkPath     func(string)
	requestPath   string
	requestMethod string
	headers       map[string][]string
}
func (m *mockHandle) HandleFunc(path string, handler func(rw http.ResponseWriter, r *http.Request)) {
	m.checkPath(path)
	if m.isErr {
		return
	}
	r, _ := http.NewRequest(m.requestMethod, m.requestPath, nil)
	r.Header = m.headers
	handler(m.writer, r)
}
func (m *mockHandle) Handle(path string, handler http.Handler)     { m.checkPath(path) }
func (m *mockHandle) ServeHTTP(http.ResponseWriter, *http.Request) {}
func TestSetupMuxRequirements(t *testing.T) {
	reg := false
	s := &Serve{
		mux: &mockHandle{
			headers: map[string][]string{
				"Accept-Encoding": {"gzip"},
			},
			requestPath:   "/test",
			requestMethod: "get",
			writer: &mockResponseWriter{
				mockWriteHeader: func(statusCode int) {},
				mockWrite:       func(b []byte) (int, error) { return 0, nil },
				mockHeader:      func() http.Header { return make(http.Header) },
			},
			checkPath: func(s string) {
				reg = true
				// bundler fileserver requires the /p/ path
				if s != "/p/" {
					t.Error("invalid path")
				}
			},
		},
		doc: &htmlDoc{[]string{}, []string{}},
	}
	newServe := s.setupMuxRequirements()
	if !reg {
		t.Error("test routine did not resolve")
	}
	if newServe == nil {
		t.Error("mux setup should not return nil")
	}
}
type mockPage struct {
	fn func(c *Request)
}
func (p *mockPage) Handle(c *Request) { p.fn(c) }
func TestHandlePage(t *testing.T) {
	p := &Serve{
		mux: &mockHandle{
			requestPath:   "/test",
			requestMethod: "get",
			writer: &mockResponseWriter{
				mockWriteHeader: func(statusCode int) {},
				mockWrite:       func(b []byte) (int, error) { return 0, nil },
				mockHeader:      func() http.Header { return make(http.Header) },
			},
			checkPath: func(s string) {},
		},
		doc: &htmlDoc{[]string{}, []string{}},
	}
	reg := false
	p.HandlePage("/test", &mockPage{
		fn: func(c *Request) { reg = true },
	})
	if !reg {
		t.Errorf("test routine did not resolve")
	}
}
func TestHandleFunc(t *testing.T) {
	writer := &mockResponseWriter{
		mockWriteHeader: func(statusCode int) {},
		mockWrite:       func(b []byte) (int, error) { return 0, nil },
		mockHeader:      func() http.Header { return make(http.Header) },
	}
	reg := false
	var tt = []struct {
		headFn      func(int)
		handler     func(*Request)
		path        string
		requestPath string
	}{
		// ensure the correct error code is transmitted when bad manifest
		// data is passed to render page.
		{
			func(code int) {
				reg = true
				if code != http.StatusInternalServerError {
					t.Errorf("expected status 500 upon bad manifest data got %d", code)
				}
			},
			func(c *Request) {
				c.RenderPage("", c)
			},
			"/test", "/test",
		},
		// status ok when manifest data is valid
		{
			func(code int) {
				reg = true
				if code != http.StatusOK {
					t.Errorf("expected status 200 upon valid manifest data got %d", code)
				}
			},
			func(c *Request) {
				props := make(map[string]interface{})
				props["test"] = "test_data"
				c.RenderPage("", props)
			},
			"/test", "/test",
		},
		// status ok when manifest data is nil
		{
			func(code int) {
				reg = true
				if code != http.StatusOK {
					t.Errorf("expected status 200 upon nil manifest data got %d", code)
				}
			},
			func(c *Request) {
				c.RenderPage("", nil)
			},
			"/test", "/test",
		},
		// status bad request slug count is invalid for this request.
		{
			func(code int) {
				reg = true
				if code != http.StatusBadRequest {
					t.Errorf("expected status 400 upon incorrectly formatted slugs got %d", code)
				}
			},
			func(c *Request) {
				c.RenderPage("", nil)
			},
			"/test/{cat}/{dog}", "/test/",
		},
	}
	for _, d := range tt {
		reg = false
		writer.mockWriteHeader = d.headFn
		s := &Serve{
			mux: &mockHandle{
				requestPath:   d.requestPath,
				requestMethod: "get",
				writer:        writer,
				checkPath:     func(s string) {},
			},
			doc: &htmlDoc{[]string{}, []string{}},
		}
		s.HandleFunc(d.path, func(c *Request) {
			d.handler(c)
		})
		// test case was not ran, fail anyways
		if !reg {
			t.Error("test routine did not resolve")
		}
	}
}
func TestSetupMuxRequirements_BundleModes(t *testing.T) {
	header := make(http.Header)
	w := &mockResponseWriter{
		mockWriteHeader: func(statusCode int) {},
		mockWrite:       func(b []byte) (int, error) { return 0, nil },
		mockHeader:      func() http.Header { return header },
	}
	s := &Serve{
		mux: &mockHandle{
			headers:       map[string][]string{},
			requestPath:   "/test",
			requestMethod: "get",
			writer:        w,
			checkPath:     func(s string) {},
		},
		doc: &htmlDoc{[]string{}, []string{}},
	}
	bundleModes := []struct {
		mode   BundleMode
		policy string
	}{
		{DevBundleMode, "no-cache, no-store, max-age=0, must-revalidate"},
		{ProdBundleMode, "public, max-age=31536000, immutable"},
	}
	startingDevMode := CurrentDevMode
	for _, m := range bundleModes {
		CurrentDevMode = m.mode
		s.setupMuxRequirements()
		if w.Header().Get("Cache-Control") != m.policy {
			t.Errorf("policy mismatch '%d' got '%s'", m.mode, w.Header().Get("Cache-Control"))
			return
		}
	}
	CurrentDevMode = startingDevMode
}
func TestHandleFuncs(t *testing.T) {
	writer := &mockResponseWriter{
		mockWriteHeader: func(statusCode int) {},
		mockWrite:       func(b []byte) (int, error) { return 0, nil },
		mockHeader:      func() http.Header { return make(http.Header) },
	}
	reg := false
	var tt = []struct {
		headFn      func(int)
		handler     func(*Request)
		path        string
		requestPath string
	}{
		// ensure the correct error code is transmitted when bad manifest
		// data is passed to render page.
		{
			func(code int) {
				reg = true
				if code != http.StatusInternalServerError {
					t.Errorf("expected status 500 upon bad manifest data got %d", code)
				}
			},
			func(c *Request) {
				c.RenderPages(c, "")
			},
			"/test", "/test",
		},
		// status ok when manifest data is valid
		{
			func(code int) {
				reg = true
				if code != http.StatusOK {
					t.Errorf("expected status 200 upon valid manifest data got %d", code)
				}
			},
			func(c *Request) {
				props := make(map[string]interface{})
				props["test"] = "test_data"
				c.RenderPages(props, "")
			},
			"/test", "/test",
		},
		// status ok when manifest data is nil
		{
			func(code int) {
				reg = true
				if code != http.StatusOK {
					t.Errorf("expected status 200 upon nil manifest data got %d", code)
				}
			},
			func(c *Request) {
				c.RenderPages(nil, "")
			},
			"/test", "/test",
		},
		// status bad request slug count is invalid for this request.
		{
			func(code int) {
				reg = true
				if code != http.StatusBadRequest {
					t.Errorf("expected status 400 upon incorrectly formatted slugs got %d", code)
				}
			},
			func(c *Request) {
				c.RenderPages(nil, "")
			},
			"/test/{cat}/{dog}", "/test/",
		},
		// status ok, standard multi-page render
		{
			func(code int) {
				reg = true
				if code != http.StatusOK {
					t.Errorf("expected status 200 upon correct multi-page request%d", code)
				}
			},
			func(c *Request) {
				c.RenderPages(nil, "", "something_2")
			},
			"/test", "/test",
		},
	}
	for _, d := range tt {
		reg = false
		writer.mockWriteHeader = d.headFn
		s := &Serve{
			mux: &mockHandle{
				requestPath:   d.requestPath,
				requestMethod: "get",
				writer:        writer,
				checkPath:     func(s string) {},
			},
			doc: &htmlDoc{[]string{}, []string{}},
		}
		s.HandleFunc(d.path, func(c *Request) {
			d.handler(c)
		})
		// test case was not ran, fail anyways
		if !reg {
			t.Error("test routine did not resolve")
		}
	}
}
func TestServe(t *testing.T) {
	s := &Serve{
		mux: &mockHandle{
			requestPath:   "/test",
			requestMethod: "get",
			writer:        nil,
			checkPath:     func(s string) {},
		},
		doc: &htmlDoc{[]string{}, []string{}},
	}
	f := s.Serve()
	if f == nil {
		t.Error("invalid serve")
	}
}
func TestBuildHTMLPages(t *testing.T) {
	t.Run("use static content", func(t *testing.T) {
		p := PageRender("thing")
		staticResourceMap[p] = true
		tmpDir := t.TempDir()
		tempBdir := bundleDir
		bundleDir = tmpDir
		t.Cleanup(func() {
			bundleDir = tempBdir
		})
		dir := fmt.Sprintf("%s%c%s", http.Dir(tmpDir), os.PathSeparator, p)
		ioutil.WriteFile(dir, []byte("<body> thing </body> <head> thint2 </head>"), 0666)
		o := buildHTMLPages([]byte(""), p)
		if len(o.Head) != 1 && len(o.Body) != 1 {
			t.Errorf("body and head len do not match")
		}
	})
	t.Run("wrap content", func(t *testing.T) {
		p := PageRender("wrapme")
		wrapDocRender[p] = &DocumentRenderer{
			fn: func(ctx context.Context, s string, b []byte, hd *htmlDoc) (*htmlDoc, context.Context) {
				hd.Body = append(hd.Body, "thing thing")
				return hd, ctx
			},
			version: "some_version",
		}
		t.Cleanup(func() {
			wrapDocRender[p] = nil
		})
		o := buildHTMLPages([]byte(""), p)
		if len(o.Body) != 1 {
			t.Errorf("did not apply wrap doc correctly")
		}
	})
}
func TestParseStaticDocument(t *testing.T) {
	t.Run("valid content", func(t *testing.T) {
		tempDir := t.TempDir()
		dir := fmt.Sprintf("%s/test.html", tempDir)
		ioutil.WriteFile(dir, []byte("stuff"), 0666)
		f, err := parseStaticDocument(dir)
		if err != nil {
			t.Errorf("should not throw error")
		}
		if f != "stuff" {
			t.Errorf("did not get expected static document content")
		}
	})
	t.Run("no content", func(t *testing.T) {
		_, err := parseStaticDocument("not real")
		if err == nil {
			t.Errorf("error should have been thrown")
		}
	})
}
func TestNew_ValidPublicDir(t *testing.T) {
	cpublicDir := publicDir
	tempDir := t.TempDir()
	publicDir = fmt.Sprintf("%s/public.html", tempDir)
	file, err := os.Create(publicDir)
	if err != nil {
		t.Error("error in test - publicDir cannot be created")
	}
	body := `<div id="test_id"> Thing </div>`
	file.Write([]byte(fmt.Sprintf(`<body>%s</body>`, body)))
	file.Close()
	s, err := New()
	if err != nil {
		t.Errorf("cannot create new orbit handler %s", err.Error())
	}
	if s == nil {
		t.Error("orbit handler should not be nil")
	}
	containsBody := false
	for _, b := range s.doc.Body {
		if strings.Contains(b, body) {
			containsBody = true
		}
	}
	if !containsBody {
		t.Error("body conent not applied correctly")
	}
	publicDir = cpublicDir
}
func TestNew_InvalidPublicDir(t *testing.T) {
	cpublicDir := publicDir
	publicDir = "not_a_valid_directory.fpea.toast"
	s, err := New()
	if err != nil {
		t.Error("should not through error upon invalid directory")
	}
	if s == nil {
		t.Error("orbit handler should not be nil")
	}
	publicDir = cpublicDir
}
