// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package mock

import "github.com/GuyARoss/orbit/pkg/jsparse"

type MockJsDocument struct {
	name      string
	extension string
}

func (m *MockJsDocument) WriteFile(string) error { return nil }
func (m *MockJsDocument) Key() string            { return "" }
func (m *MockJsDocument) Name() string           { return m.name }
func (m *MockJsDocument) Imports() []*jsparse.ImportDependency {
	return make([]*jsparse.ImportDependency, 0)
}
func (m *MockJsDocument) AddImport(*jsparse.ImportDependency) []*jsparse.ImportDependency {
	return make([]*jsparse.ImportDependency, 0)
}
func (m *MockJsDocument) Other() []string                         { return []string{} }
func (m *MockJsDocument) AddOther(string) []string                { return []string{} }
func (m *MockJsDocument) Extension() string                       { return m.extension }
func (m *MockJsDocument) AddSerializable(s jsparse.JSSerialize)   {}
func (m *MockJsDocument) DefaultExport() *jsparse.JsDocumentScope { return nil }

func NewMockJSDocument(name string, extension string) *MockJsDocument {
	return &MockJsDocument{
		name:      name,
		extension: extension,
	}
}
