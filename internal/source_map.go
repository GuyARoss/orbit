package internal

import (
	dependtree "github.com/GuyARoss/orbit/pkg/depend_tree"
	"github.com/GuyARoss/orbit/pkg/jsparse"
)

type DependencySettings struct {
	dirList          []string
	pathDependencies map[string][]string
}

func createDirListFromRoot(c []*PackedComponent) []string {
	lst := make([]string, len(c))

	for i, c := range c {
		lst[i] = c.OriginalFilePath
	}
	return lst
}

func createPathDependenciesMapFromRoot(c []*PackedComponent) map[string][]string {
	m := make(map[string][]string)

	for _, component := range c {
		finalDependendices := make([]string, 0)
		for _, d := range component.Dependencies {
			if d.Type == jsparse.LocalImportType {
				finalDependendices = append(finalDependendices, d.InitialPath)
			}
		}

		m[component.OriginalFilePath] = finalDependendices
	}

	return m
}

func (s *DependencySettings) DirList(path string) ([]string, error) {
	return s.dirList, nil
}

func (s *DependencySettings) PathDependencies(path string) ([]string, error) {
	return s.pathDependencies[path], nil
}

func CreateSourceMap(path string, c []*PackedComponent) (*dependtree.DependencySourceMap, error) {
	m := &dependtree.ManagedDependencyTree{
		Settings: &DependencySettings{
			// @@todo: parallelize these processes
			dirList:          createDirListFromRoot(c),
			pathDependencies: createPathDependenciesMapFromRoot(c),
		},
	}

	treeNode, err := m.Create(path)
	if err != nil {
		return nil, err
	}

	return treeNode.SourceMap(), nil
}
