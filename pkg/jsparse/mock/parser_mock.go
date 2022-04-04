// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package mock

import "github.com/GuyARoss/orbit/pkg/jsparse"

type MockJSParser struct {
	ParseDocument jsparse.JSDocument
	Err           error
}

func (m *MockJSParser) Parse(string, string) (jsparse.JSDocument, error) {
	return m.ParseDocument, m.Err
}
