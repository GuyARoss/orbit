package mock

import "github.com/GuyARoss/orbit/pkg/jsparse"

type MockJSParser struct {
	ParseDocument jsparse.JSDocument
	Err           error
}

func (m *MockJSParser) Parse(string, string) (jsparse.JSDocument, error) {
	return m.ParseDocument, m.Err
}
