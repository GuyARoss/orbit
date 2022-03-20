// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package mock

import "net/http"

type MockHotReload struct {
	DidReload        bool
	currentBundleKey string
	reloadErr        error
}

func (m *MockHotReload) ReloadSignal() error {
	m.DidReload = true

	return m.reloadErr
}

func (m *MockHotReload) HandleWebSocket(w http.ResponseWriter, r *http.Request) {}
func (m *MockHotReload) CurrentBundleKey() string {
	return m.currentBundleKey
}
