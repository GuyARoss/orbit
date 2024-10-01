// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package mock

import (
	"context"
	"errors"

	"github.com/GuyARoss/orbit/pkg/embedutils"
	"github.com/GuyARoss/orbit/pkg/jsparse"
	"github.com/GuyARoss/orbit/pkg/jsparse/mock"
	jsparsemock "github.com/GuyARoss/orbit/pkg/jsparse/mock"
	"github.com/GuyARoss/orbit/pkg/webwrap"
)

type MockWrapper struct {
	Satisfy    bool
	FailBundle bool
}

func (m *MockWrapper) DocumentTag(string) string { return "" }

func (m *MockWrapper) VerifyRequirements() error {
	return nil
}

func (m *MockWrapper) Apply(doc jsparse.JSDocument) (map[string]jsparse.JSDocument, error) {
	f := map[string]jsparse.JSDocument{"normal": &jsparsemock.MockJsDocument{}}
	return f, nil
}

func (m *MockWrapper) DoesSatisfyConstraints(doc jsparse.JSDocument) bool { return m.Satisfy }
func (m *MockWrapper) Version() string                                    { return "" }
func (m *MockWrapper) RequiredBodyDOMElements(ctx context.Context, opts *webwrap.CacheDOMOpts) []string {
	return nil
}

func (m *MockWrapper) Stats() *webwrap.WrapStats {
	return &webwrap.WrapStats{}
}

func (b *MockWrapper) Setup(context.Context, *webwrap.BundleOpts) (*webwrap.BundledResource, error) {
	return &webwrap.BundledResource{
		BundleOpFileDescriptor: map[string]string{
			"normal": "test",
		},
		Configurators: []webwrap.BundleConfigurator{
			{Page: mock.NewMockJSDocument("test", "jsx", "test"), FilePath: ""},
		},
	}, nil
}

func (b *MockWrapper) HydrationFile() []embedutils.FileReader {
	return nil
}

func (b *MockWrapper) Bundle(string, string) error {
	if b.FailBundle {
		return errors.New("fail")
	}

	return nil
}
