// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package mock

import (
	"context"

	"github.com/GuyARoss/orbit/internal/libout"
	"github.com/GuyARoss/orbit/internal/srcpack"
	"github.com/GuyARoss/orbit/pkg/webwrap"
)

type MockBundleWriter struct{}

func (m *MockBundleWriter) WriteLibout(files libout.Libout, fOpts *libout.FilePathOpts) error {
	return nil
}
func (m *MockBundleWriter) AcceptComponent(ctx context.Context, c srcpack.PackComponent, cacheOpts *webwrap.CacheDOMOpts) error {
	return nil
}
func (m *MockBundleWriter) AcceptComponents(ctx context.Context, comps []srcpack.PackComponent, cacheOpts *webwrap.CacheDOMOpts) error {
	return nil
}
