// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package dependtree

import (
	"testing"

	"github.com/GuyARoss/orbit/pkg/depend_tree/mock"
)

func TestMergeOverKey(t *testing.T) {
	first := map[string][]string{
		"thing": {"fish"},
	}

	final := DependencySourceMap(first).MergeOverKey(map[string][]string{
		"thing2": {"stuff", "stuff2"},
		"thing":  {"apple"},
	})

	if len(final) != 2 {
		t.Errorf("final does not overwrite preexisting items on merge")
	}

	if len(final["thing"]) != 1 {
		t.Errorf("did not merge initial keys with merge keys")
	}
}

func TestMergeDependTree(t *testing.T) {
	first := map[string][]string{
		"thing": {"fish"},
	}

	final := DependencySourceMap(first).Merge(map[string][]string{
		"thing2": {"stuff", "stuff2"},
		"thing":  {"apple"},
	})

	if len(final) != 2 {
		t.Errorf("final does not include merged items")
	}

	if len(final["thing"]) != 2 {
		t.Errorf("did not merge initial keys with merge keys")
	}
}

func TestCreateDependTree(t *testing.T) {
	dep := make(map[string][]string)
	dep["/pages"] = []string{"../components/modal.jsx", "../components/layout.jsx"}
	dep["../components/modal.jsx"] = []string{"../components/form.jsx"}
	dep["../components/layout.jsx"] = []string{"../components/header.jsx"}
	dep["/files"] = []string{"../thing.jsx"}

	f := &mock.MockDependencyTree{
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
	for f, k := range shallowMap {
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
