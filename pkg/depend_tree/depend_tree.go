// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.
package dependtree

type DependencyTreeNode struct {
	Value  string
	Right  *DependencyTreeNode
	Child  *DependencyTreeNode
	IsRoot bool
}

type DependencySourceMap struct {
	sourceMap map[string][]string
}

func (d *DependencySourceMap) FindRoot(path string) []string {
	return d.sourceMap[path]
}

func (d *DependencyTreeNode) values(current []string) []string {
	if d == nil || d.Value == "" {
		return current
	}

	if current == nil {
		current = make([]string, 0)
	}

	if !d.IsRoot {
		current = append(current, d.Value)
		current = d.Right.values(current)
	}

	current = d.Child.values(current)

	return current
}

func (d *DependencyTreeNode) SourceMap() *DependencySourceMap {
	m := make(map[string][]string)

	current := d
	for current != nil {
		values := current.values(nil)
		if len(values) == 0 || current.Value == "" {
			current = current.Right
			continue
		}

		for _, v := range values {
			if m[v] == nil {
				m[v] = make([]string, 0)
			}

			m[v] = append(m[v], current.Value)
		}

		current = current.Right
	}
	return &DependencySourceMap{m}
}

type DependencyTree interface {
	// DirList fetches the current paths within a directory
	DirList(path string) ([]string, error)

	// PathDependencies finds the given dependencies for a given path
	PathDependencies(path string) ([]string, error)
}

type ManagedDependencyTree struct {
	rootNode *DependencyTreeNode
	Settings DependencyTree
}

func mapNode(s DependencyTree, path string) (*DependencyTreeNode, error) {
	dependencies, err := s.PathDependencies(path)

	if err != nil {
		return nil, err
	}
	root := &DependencyTreeNode{}
	current := root
	for _, d := range dependencies {
		mapResp, err := mapNode(s, d)
		if err != nil {
			return nil, err
		}

		current.IsRoot = false
		current.Value = d
		current.Child = mapResp
		current.Right = &DependencyTreeNode{}

		current = current.Right
	}
	return root, nil
}

func (s *ManagedDependencyTree) Create(initialPath string) (*DependencyTreeNode, error) {
	dirs, err := s.Settings.DirList(initialPath)
	if err != nil {
		return nil, err
	}
	root := &DependencyTreeNode{}

	current := root

	for _, d := range dirs {
		mapResp, err := mapNode(s.Settings, d)
		if err != nil {
			return nil, err
		}

		current.Value = d
		current.Child = mapResp
		current.Right = &DependencyTreeNode{}
		current.IsRoot = true

		current = current.Right
	}

	s.rootNode = root

	return root, nil
}
