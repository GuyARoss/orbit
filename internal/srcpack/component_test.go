// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package srcpack

import (
	"context"
	"errors"
	"testing"

	"github.com/GuyARoss/orbit/pkg/bundler"
	"github.com/GuyARoss/orbit/pkg/jsparse"
	webwrapper "github.com/GuyARoss/orbit/pkg/web_wrapper"
)

type mockJsDocument struct{}

func (m *mockJsDocument) WriteFile(string) error { return nil }
func (m *mockJsDocument) Key() string            { return "" }
func (m *mockJsDocument) Name() string           { return "" }
func (m *mockJsDocument) Imports() []*jsparse.ImportDependency {
	return make([]*jsparse.ImportDependency, 0)
}
func (m *mockJsDocument) AddImport(*jsparse.ImportDependency) []*jsparse.ImportDependency {
	return make([]*jsparse.ImportDependency, 0)
}
func (m *mockJsDocument) Other() []string          { return []string{} }
func (m *mockJsDocument) AddOther(string) []string { return []string{} }
func (m *mockJsDocument) Extension() string        { return "" }

type mockWrapper struct {
	satisfy bool
}

func (m *mockWrapper) Apply(doc jsparse.JSDocument, t string) jsparse.JSDocument {
	return &mockJsDocument{}
}

func (m *mockWrapper) NodeDependencies() map[string]string { return make(map[string]string) }

func (m *mockWrapper) DoesSatisfyConstraints(p string) bool { return m.satisfy }
func (m *mockWrapper) Version() string                      { return "" }
func (m *mockWrapper) RequiredBodyDOMElements(ctx context.Context, opts *webwrapper.CacheDOMOpts) []string {
	return nil
}

type EmptyBundler struct {
	FailBundle bool
}

func (b *EmptyBundler) Setup(context.Context, *bundler.BundleOpts) (*bundler.BundledResource, error) {
	return &bundler.BundledResource{
		ConfiguratorPage: &mockJsDocument{},
	}, nil
}

func (b *EmptyBundler) Bundle(string) error {
	if b.FailBundle {
		return errors.New("fail")
	}

	return nil
}

func (b *EmptyBundler) NodeDependencies() map[string]string {
	return make(map[string]string)
}

func TestNewComponent_BundleKey(t *testing.T) {
	tt := []struct {
		s *NewComponentOpts
		k string
	}{
		{&NewComponentOpts{
			FilePath:      "something.test",
			WebDir:        "./webDir",
			DefaultKey:    "thing",
			JSParser:      &jsparse.EmptyParser{},
			Bundler:       &EmptyBundler{false},
			JSWebWrappers: []webwrapper.JSWebWrapper{&mockWrapper{satisfy: true}},
		}, "thing"},
	}

	for _, d := range tt {
		c, err := NewComponent(context.TODO(), d.s)

		if err != nil {
			t.Error("error should not be thrown")
		}

		if c.BundleKey != d.k {
			t.Errorf("expected %s got %s", d.k, c.BundleKey)
		}
	}
}

func TestNewComponent_Failures(t *testing.T) {
	tt := []struct {
		s *NewComponentOpts
	}{
		// cannot parse inital doc failure
		{&NewComponentOpts{JSParser: &jsparse.JSFileParser{}}},

		// bad web wrap
		{&NewComponentOpts{JSParser: &jsparse.EmptyParser{}}},

		// bundler failure
		{&NewComponentOpts{JSParser: &jsparse.EmptyParser{},
			JSWebWrappers: []webwrapper.JSWebWrapper{&mockWrapper{satisfy: true}},
			Bundler:       &EmptyBundler{true},
		}},
	}

	for _, d := range tt {
		_, err := NewComponent(context.TODO(), d.s)
		if err == nil {
			t.Errorf("expected failure upon component creation")
		}
	}
}
