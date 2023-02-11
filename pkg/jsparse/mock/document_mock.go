// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package mock

import "github.com/GuyARoss/orbit/pkg/jsparse"

type MockJsDocument struct {
	name          string
	extension     string
	defaultExport string
}

func (m *MockJsDocument) Clone() jsparse.JSDocument {
	return nil
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
func (m *MockJsDocument) Other() []string                       { return []string{} }
func (m *MockJsDocument) AddOther(...string)                    {}
func (m *MockJsDocument) Extension() string                     { return m.extension }
func (m *MockJsDocument) AddSerializable(s jsparse.JSSerialize) {}
func (m *MockJsDocument) DefaultExport() *jsparse.JsDocumentScope {
	return &jsparse.JsDocumentScope{
		TokenType: jsparse.ExportDefaultToken,
		Name:      m.name,
		Export:    jsparse.ExportDefault,
		Args:      make(jsparse.JSDocArgList, 0),
	}
}

func NewMockJSDocument(name string, extension string, defaultExport string) *MockJsDocument {
	return &MockJsDocument{
		name:          name,
		extension:     extension,
		defaultExport: defaultExport,
	}
}
