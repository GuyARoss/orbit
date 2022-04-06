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
