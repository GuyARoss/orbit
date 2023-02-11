// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package srcpack

import (
	"context"
	"testing"

	"github.com/GuyARoss/orbit/pkg/jsparse"
	"github.com/GuyARoss/orbit/pkg/jsparse/mock"
	"github.com/GuyARoss/orbit/pkg/webwrap"
	webwrapmock "github.com/GuyARoss/orbit/pkg/webwrap/mock"
)

func TestNewComponent_BundleKey(t *testing.T) {
	tt := []struct {
		s *NewComponentOpts
		k string
	}{
		{&NewComponentOpts{
			FilePath:   "something.test",
			WebDir:     "./webDir",
			DefaultKey: "thing",
			JSParser: &mock.MockJSParser{
				Err:           nil,
				ParseDocument: mock.NewMockJSDocument("test", "jsx", "test"),
			},
			JSWebWrappers: []webwrap.JSWebWrapper{&webwrapmock.MockWrapper{Satisfy: true, FailBundle: false}},
		}, "thing"},
	}

	for _, d := range tt {
		c, err := NewComponent(context.TODO(), d.s)

		if err != nil {
			t.Error("error should not be thrown")
		}

		if c.BundleKey() != d.k {
			t.Errorf("expected %s got %s", d.k, c.BundleKey())
		}
	}
}

func TestNewComponent_Failures(t *testing.T) {
	tt := []struct {
		s *NewComponentOpts
	}{
		// cannot parse initial doc failure
		{&NewComponentOpts{JSParser: &jsparse.JSFileParser{}}},

		// bad web wrap
		{&NewComponentOpts{JSParser: &mock.MockJSParser{
			Err:           nil,
			ParseDocument: mock.NewMockJSDocument("test", "jsx", "test"),
		}}},

		// bundler failure
		{&NewComponentOpts{
			JSParser: &mock.MockJSParser{
				Err:           nil,
				ParseDocument: mock.NewMockJSDocument("test", "jsx", "test"),
			},
			JSWebWrappers:       []webwrap.JSWebWrapper{&webwrapmock.MockWrapper{Satisfy: true, FailBundle: true}},
			SkipFirstPassBundle: false,
		}},

		// default export is not present
		{&NewComponentOpts{JSParser: &mock.MockJSParser{
			Err:           nil,
			ParseDocument: mock.NewMockJSDocument("test", "jsx", ""),
		},
			JSWebWrappers: []webwrap.JSWebWrapper{&webwrapmock.MockWrapper{Satisfy: true, FailBundle: true}},
		}},
	}

	for i, d := range tt {
		_, err := NewComponent(context.TODO(), d.s)
		if err == nil {
			t.Errorf("(%d) expected failure upon component creation '%s'", i, err)
		}
	}
}

func TestParsePath(t *testing.T) {
	expected := "thing/something"
	if got := parsePath("./thing/something"); got != expected {
		t.Errorf("expected '%s' got '%s'", expected, got)
	}
}

func TestFindBundleKey(t *testing.T) {
	m := PackComponentFileMap{
		"test":    &Component{bundleKey: "something", name: "Working"},
		"test123": &Component{bundleKey: "notsomething", name: "Broken"},
	}

	if m.FindBundleKey("something").Name() != "Working" {
		t.Errorf("did not return correct component got '%s'", m.FindBundleKey("something").Name())
		return
	}

	if m.FindBundleKey("do_not_exist") != nil {
		t.Errorf("did not return correct component")
		return
	}
}

func TestRepack(t *testing.T) {
	comp, err := NewComponent(context.TODO(), &NewComponentOpts{
		FilePath:   "something.test",
		WebDir:     "./webDir",
		DefaultKey: "thing",
		JSParser: &mock.MockJSParser{
			Err:           nil,
			ParseDocument: mock.NewMockJSDocument("test", "jsx", "test"),
		},
		JSWebWrappers: []webwrap.JSWebWrapper{&webwrapmock.MockWrapper{Satisfy: true, FailBundle: false}},
	})
	if err != nil {
		t.Errorf("error should not be thrown '%s'", err)
		return
	}

	err = comp.Repack()
	if err != nil {
		t.Errorf("error should not be thrown during repack '%s'", err)
		return
	}
}
