package orbitgen

import (
	"strings"
	"testing"
	"fmt"
	"net/http"
	"os"
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
}

func (m *mockHandle) HandleFunc(path string, handler func(rw http.ResponseWriter, r *http.Request)) {
	m.checkPath(path)
	if m.isErr {
		return
	}

	r, _ := http.NewRequest(m.requestMethod, m.requestPath, nil)

	handler(m.writer, r)
}

func (m *mockHandle) Handle(path string, handler http.Handler)     { m.checkPath(path) }
func (m *mockHandle) ServeHTTP(http.ResponseWriter, *http.Request) {}

func TestSetupMuxRequirements(t *testing.T) {
	reg := false

	s := &Serve{
		mux: &mockHandle{
			requestPath:   "/test",
			requestMethod: "get",
			writer:        nil,
			checkPath: func(s string) {
				reg = true

				// our bundler fileserver requires the /p/ path
				if s != "/p/" {
					t.Error("invalid path")
				}
			},
		},
		doc: &htmlDoc{[]string{}, []string{}},
	}

	s.setupMuxRequirements()

	if !reg {
		t.Error("test routine did not resolve")
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
