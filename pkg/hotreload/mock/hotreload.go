// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package mock

import (
	"net/http"
)

type MockHotReload struct {
	DidReload         bool
	currentBundleKeys []string
	reloadErr         error

	Active bool
}

func (m *MockHotReload) IsActive() bool {
	return m.Active
}
func (m *MockHotReload) IsActiveBundle(string) bool {
	return m.Active
}

func (m *MockHotReload) ReloadSignal() error {
	m.DidReload = true
	return m.reloadErr
}

func (m *MockHotReload) HandleWebSocket(w http.ResponseWriter, r *http.Request) {}
func (m *MockHotReload) CurrentBundleKeys() []string {
	return m.currentBundleKeys
}
