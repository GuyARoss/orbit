// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package mock

type MockSocket struct {
	DidWrite bool
}

func (m *MockSocket) WriteJSON(interface{}) error {
	m.DidWrite = true
	return nil
}
func (m *MockSocket) Close() error {
	return nil
}
func (m *MockSocket) ReadJSON(interface{}) error {
	return nil
}
