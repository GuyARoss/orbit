package internal

import dependtree "github.com/GuyARoss/orbit/pkg/depend_tree"

type DependencySettings struct {
}

func (s *DependencySettings) DirList(path string) ([]string, error) {
	return nil, nil
}

func (s *DependencySettings) PathDependencies(path string) ([]string, error) {
	return nil, nil
}

func CreateSourceMap(path string) (*dependtree.DependencySourceMap, error) {
	m := &dependtree.ManagedDependencyTree{
		Settings: &DependencySettings{},
	}

	treeNode, err := m.Create(path)
	if err != nil {
		return nil, err
	}

	return treeNode.SourceMap(), nil
}
