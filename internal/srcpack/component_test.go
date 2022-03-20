// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package srcpack

import (
	"context"
	"testing"

	bundlermock "github.com/GuyARoss/orbit/pkg/bundler/mock"
	"github.com/GuyARoss/orbit/pkg/jsparse"
	"github.com/GuyARoss/orbit/pkg/webwrap"
	webwrapmock "github.com/GuyARoss/orbit/pkg/webwrap/mock"
)

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
			Bundler:       &bundlermock.EmptyBundler{false},
			JSWebWrappers: []webwrap.JSWebWrapper{&webwrapmock.MockWrapper{Satisfy: true}},
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
		// cannot parse inital doc failure
		{&NewComponentOpts{JSParser: &jsparse.JSFileParser{}}},

		// bad web wrap
		{&NewComponentOpts{JSParser: &jsparse.EmptyParser{}}},

		// bundler failure
		{&NewComponentOpts{JSParser: &jsparse.EmptyParser{},
			JSWebWrappers: []webwrap.JSWebWrapper{&webwrapmock.MockWrapper{Satisfy: true}},
			Bundler:       &bundlermock.EmptyBundler{true},
		}},
	}

	for _, d := range tt {
		_, err := NewComponent(context.TODO(), d.s)
		if err == nil {
			t.Errorf("expected failure upon component creation")
		}
	}
}
