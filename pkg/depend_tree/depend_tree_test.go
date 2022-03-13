// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package dependtree

import (
	"testing"
)

type MockDependencyTree struct {
	Dirs         []string
	Dependencies map[string][]string
}

func (s *MockDependencyTree) DirList(path string) ([]string, error) {
	return s.Dirs, nil
}

func (s *MockDependencyTree) PathDependencies(path string) ([]string, error) {
	return s.Dependencies[path], nil
}

func TestCreateDependTree(t *testing.T) {
	dep := make(map[string][]string)
	dep["/pages"] = []string{"../components/modal.jsx", "../components/layout.jsx"}
	dep["../components/modal.jsx"] = []string{"../components/form.jsx"}
	dep["../components/layout.jsx"] = []string{"../components/header.jsx"}
	dep["/files"] = []string{"../thing.jsx"}

	f := &MockDependencyTree{
		Dirs:         []string{"/pages"},
		Dependencies: dep,
	}

	s := ManagedDependencyTree{
		Settings: f,
	}
	resp, err := s.Create("/pages")
	if err != nil {
		t.FailNow()
	}

	shallowMap := resp.SourceMap()
	for f, k := range shallowMap.sourceMap {
		if f == "/pages" || f == "/files" {
			continue
		}
		switch f {
		case "/pages":
			{
				c := []string{"../components/modal.jsx", "../components/layout.jsx", "../components/form.jsx", "../components/header.jsx"}
				found := false
				for _, k := range c {
					if k == f {
						found = true
					}
				}
				if !found {
					t.Errorf("%s does not include %s", k, f)
				}
			}
		case "/files":
			{
				c := []string{"../thing.jsx"}
				found := false
				for _, k := range c {
					if k == f {
						found = true
					}
				}
				if !found {
					t.Errorf("%s does not include %s", k, f)
				}
			}
		}
	}
}
