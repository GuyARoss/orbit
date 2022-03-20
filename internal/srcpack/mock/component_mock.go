// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package mock

import (
	"sync"

	"github.com/GuyARoss/orbit/pkg/jsparse"
	"github.com/GuyARoss/orbit/pkg/webwrap"
)

type MockPackedComponent struct {
	WasRepacked bool
}

func (m *MockPackedComponent) Repack() error {
	m.WasRepacked = true
	return nil
}

func (m *MockPackedComponent) RepackForWaitGroup(wg *sync.WaitGroup, c chan error) {}
func (m *MockPackedComponent) OriginalFilePath() string                            { return "" }
func (m *MockPackedComponent) Dependencies() []*jsparse.ImportDependency           { return nil }
func (m *MockPackedComponent) BundleKey() string                                   { return "" }
func (m *MockPackedComponent) Name() string                                        { return "" }
func (m *MockPackedComponent) WebWrapper() webwrap.JSWebWrapper                    { return nil }
