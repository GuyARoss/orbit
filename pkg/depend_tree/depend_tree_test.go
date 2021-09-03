package dependtree

import (
	"fmt"
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

	f := &MockDependencyTree{
		Dirs:         []string{"/pages"},
		Dependencies: dep,
	}
	resp, err := Create(f, "/pages")
	if err != nil {
		t.FailNow()
	}

	fmt.Println(resp)
	t.Log(resp)
}
