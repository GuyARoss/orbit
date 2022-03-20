// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package mock

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
