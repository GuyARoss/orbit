package webwrapper

import (
	"errors"
	"testing"

	"github.com/GuyARoss/orbit/pkg/jsparse"
)

type mockJSDoc struct {
	extension string
	name      string
}

func (s *mockJSDoc) WriteFile(string) error                                          { return nil }
func (s *mockJSDoc) Key() string                                                     { return "" }
func (s *mockJSDoc) Name() string                                                    { return s.name }
func (s *mockJSDoc) Imports() []*jsparse.ImportDependency                            { return nil }
func (s *mockJSDoc) AddImport(*jsparse.ImportDependency) []*jsparse.ImportDependency { return nil }
func (s *mockJSDoc) Other() []string                                                 { return nil }
func (s *mockJSDoc) AddOther(string) []string                                        { return nil }
func (s *mockJSDoc) Extension() string                                               { return s.extension }

func TestApplyReact_Error(t *testing.T) {
	tt := []struct {
		doc *mockJSDoc
		err error
	}{
		{&mockJSDoc{extension: "blah"}, ErrInvalidComponent},
		{&mockJSDoc{extension: "jsx", name: "lowercaseThing"}, ErrComponentExport},
	}

	r := &ReactWebWrapper{}

	for i, d := range tt {
		_, err := r.Apply(d.doc)

		if !errors.Is(err, d.err) {
			t.Errorf("(%d) got incorrect error", i)
		}
	}
}

func TestApplyReact(t *testing.T) {
	r := &ReactWebWrapper{}

	p, err := r.Apply(&mockJSDoc{extension: "jsx", name: "Thing"})
	if err != nil {
		t.Error("should not expect error during valid jsx parsing")
	}

	if p.Name() != "Thing" {
		t.Errorf("expected name 'Thing' got '%s'", p.Name())
	}
}
