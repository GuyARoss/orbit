// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package mock

import (
	"context"
	"errors"

	"github.com/GuyARoss/orbit/pkg/embedutils"
	"github.com/GuyARoss/orbit/pkg/jsparse"
	jsparsemock "github.com/GuyARoss/orbit/pkg/jsparse/mock"
	"github.com/GuyARoss/orbit/pkg/webwrap"
)

type MockWrapper struct {
	Satisfy    bool
	FailBundle bool
}

func (m *MockWrapper) VerifyRequirements() error {
	return nil
}

func (m *MockWrapper) Apply(doc jsparse.JSDocument) (jsparse.JSDocument, error) {
	return &jsparsemock.MockJsDocument{}, nil
}

func (m *MockWrapper) DoesSatisfyConstraints(p string) bool { return m.Satisfy }
func (m *MockWrapper) Version() string                      { return "" }
func (m *MockWrapper) RequiredBodyDOMElements(ctx context.Context, opts *webwrap.CacheDOMOpts) []string {
	return nil
}

func (b *MockWrapper) Setup(context.Context, *webwrap.BundleOpts) ([]*webwrap.BundledResource, error) {
	return []*webwrap.BundledResource{{
		ConfiguratorPage: &jsparsemock.MockJsDocument{},
	}}, nil
}

func (b *MockWrapper) HydrationFile() []embedutils.FileReader {
	return nil
}

func (b *MockWrapper) Bundle(string) error {
	if b.FailBundle {
		return errors.New("fail")
	}

	return nil
}
