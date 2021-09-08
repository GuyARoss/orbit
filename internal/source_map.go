package internal

import (
	"sync"

	dependtree "github.com/GuyARoss/orbit/pkg/depend_tree"
	"github.com/GuyARoss/orbit/pkg/jsparse"
)

type DependencySettings struct {
	WebDir            string
	dirList           *[]string
	pathDependencyMap *map[string][]string
}

func (s *DependencySettings) cacheRootDirList(c []*PackedComponent, wg *sync.WaitGroup) {
	defer wg.Done()

	lst := make([]string, len(c))

	for i, c := range c {
		lst[i] = c.OriginalFilePath
	}
	s.dirList = &lst
}

func flatPackedImports(dependencies []*jsparse.ImportDependency) []string {
	finalDependendices := make([]string, 0)
	for _, d := range dependencies {
		if d.Type == jsparse.LocalImportType {
			finalDependendices = append(finalDependendices, d.InitialPath)
		}
	}
	return finalDependendices
}

func (s *DependencySettings) cacheRootPathDependencyMap(c []*PackedComponent, wg *sync.WaitGroup) {
	defer wg.Done()

	m := make(map[string][]string)

	for _, component := range c {
		m[component.OriginalFilePath] = flatPackedImports(component.Dependencies)
	}

	s.pathDependencyMap = &m
}

func (s *DependencySettings) DirList(path string) ([]string, error) {
	return *s.dirList, nil
}

func (s *DependencySettings) PathDependencies(path string) ([]string, error) {
	derefMap := *s.pathDependencyMap
	c := derefMap[path]

	if c != nil {
		return c, nil
	}

	page, err := jsparse.ParsePage(path, s.WebDir)
	if err != nil {
		return nil, err
	}

	return flatPackedImports(page.Imports), nil
}

func CreateSourceMap(path string, c []*PackedComponent, webDirPath string) (*dependtree.DependencySourceMap, error) {
	var wg sync.WaitGroup

	dependSettings := &DependencySettings{
		WebDir: webDirPath,
	}

	m := &dependtree.ManagedDependencyTree{
		Settings: dependSettings,
	}

	wg.Add(2)
	go dependSettings.cacheRootDirList(c, &wg)
	go dependSettings.cacheRootPathDependencyMap(c, &wg)
	wg.Wait()

	treeNode, err := m.Create(path)
	if err != nil {
		return nil, err
	}

	return treeNode.SourceMap(), nil
}
