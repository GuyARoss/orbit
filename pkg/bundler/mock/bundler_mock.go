// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package mock

import (
	"context"
	"errors"

	"github.com/GuyARoss/orbit/pkg/bundler"
	jsparsemock "github.com/GuyARoss/orbit/pkg/jsparse/mock"
)

type EmptyBundler struct {
	FailBundle bool
}

func (b *EmptyBundler) Setup(context.Context, *bundler.BundleOpts) (*bundler.BundledResource, error) {
	return &bundler.BundledResource{
		ConfiguratorPage: &jsparsemock.MockJsDocument{},
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
