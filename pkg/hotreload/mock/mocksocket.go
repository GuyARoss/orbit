// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package mock

import "encoding/json"

type MockSocket struct {
	DidWrite bool
	ReadData interface{}
}

func (m *MockSocket) WriteJSON(interface{}) error {
	m.DidWrite = true
	return nil
}
func (m *MockSocket) Close() error {
	return nil
}
func (m *MockSocket) ReadJSON(data interface{}) error {
	b, _ := json.Marshal(m.ReadData)
	json.Unmarshal(b, data)

	return nil
}
