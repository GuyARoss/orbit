package orbit

import (
	"fmt"
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

func TestHTMLDocBuild(t *testing.T) {
	data := "thingy"
	page := "this_is_a_page_key"

	head := "<div>something</div>"
	body := "<h1>header</h1>"

	var tt = []struct {
		expect string
		l      string
		r      string
	}{
		{data, `<script id="orbit_manifest" type="application/json">`, "</script>"},
		{page, `script src="/p/`, `.js">`},
		{head, "<head>", "</head>"},
		{body, "<body>", "</body>"},
	}

	doc := &htmlDoc{
		Head: []string{head},
		Body: []string{body},
	}

	flatdoc := doc.build([]byte(data), PageRender(page))

	for _, d := range tt {
		content := innerHTML(flatdoc, d.l, d.r)
		if !strings.Contains(content, d.expect) {
			t.Errorf("%s expected on parent tag but not found", d.expect)
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
	path := "/thing/{toast}"

	p := parsePathSlugs(&path)

	if p[2] != "toast" {
		t.Errorf("expected %s got %s", "toast", p[1])
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
	isErr     bool
	writer    http.ResponseWriter
	checkPath func(string)
}

func (m *mockHandle) HandleFunc(path string, handler func(rw http.ResponseWriter, r *http.Request)) {
	m.checkPath(path)
	if m.isErr {
		return
	}

	r, _ := http.NewRequest("get", "/thing", nil)

	handler(m.writer, r)
}

func (m *mockHandle) Handle(path string, handler http.Handler)     { m.checkPath(path) }
func (m *mockHandle) ServeHTTP(http.ResponseWriter, *http.Request) {}

func TestSetupMuxRequirements(t *testing.T) {
	reg := false

	s := &serve{
		mux: &mockHandle{
			writer: nil,
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
	p := &serve{
		mux: &mockHandle{
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
		headFn  func(int)
		handler func(*Request)
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
		},
	}

	for _, d := range tt {
		reg = false

		writer.mockWriteHeader = d.headFn

		s := &serve{
			mux: &mockHandle{
				writer:    writer,
				checkPath: func(s string) {},
			},
			doc: &htmlDoc{[]string{}, []string{}},
		}

		s.HandleFunc("/test", func(c *Request) {
			d.handler(c)
		})

		// test case was not ran, fail anyways
		if !reg {
			t.Error("test routine did not resolve")
		}
	}
}

func TestServe(t *testing.T) {
	s := &serve{
		mux: &mockHandle{
			writer:    nil,
			checkPath: func(s string) {},
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