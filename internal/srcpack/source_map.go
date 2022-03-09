// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// LICENSE file in the root directory of this source tree.
package srcpack

import (
	"sync"

	dependtree "github.com/GuyARoss/orbit/pkg/depend_tree"
	"github.com/GuyARoss/orbit/pkg/jsparse"
)

// javascript dependency tree, used to create a dependency tree from javascript files
// this struct should salsify the requirements for "DependencyTree" interface
type JSDependencyTree struct {
	WebDir            string
	dirList           *[]string
	pathDependencyMap *map[string][]string

	JsParser jsparse.JSParser
}

type SrcDependency interface {
	OriginalFilePath() string
	Dependencies() []*jsparse.ImportDependency
}

// given a slice of import dependencies, returns a string of local import paths
func localDependencies(dependencies []*jsparse.ImportDependency) []string {
	finalDependendices := make([]string, 0)
	for _, d := range dependencies {
		if d.Type == jsparse.LocalImportType {
			path := d.InitialPath

			// common issue with parsed paths is that they could be formatted differently
			// no resolve this, we normalize the process
			if path[0] == '/' {
				path = path[1:]
			}

			finalDependendices = append(finalDependendices, path)
		}
	}
	return finalDependendices
}

func (s *JSDependencyTree) cacheRootDirList(c []*Component, wg *sync.WaitGroup) {
	defer wg.Done()

	lst := make([]string, len(c))

	for i, c := range c {
		lst[i] = c.OriginalFilePath()
	}
	s.dirList = &lst
}

func (s *JSDependencyTree) cacheRootPathDependencyMap(c []*Component, wg *sync.WaitGroup) {
	defer wg.Done()

	m := make(map[string][]string)

	for _, component := range c {
		m[component.OriginalFilePath()] = localDependencies(component.Dependencies())
	}

	s.pathDependencyMap = &m
}

func (s *JSDependencyTree) DirList(path string) ([]string, error) {
	return *s.dirList, nil
}

// uses the js parser to get all of the dependencies for the specified file.
func (s *JSDependencyTree) PathDependencies(path string) ([]string, error) {
	pdm := *s.pathDependencyMap

	if c := pdm[path]; c != nil {
		return c, nil
	}

	page, err := s.JsParser.Parse(path, s.WebDir)
	if err != nil {
		return nil, err
	}

	return localDependencies(page.Imports()), nil
}

func New(path string, c []*Component, webDirPath string) (*dependtree.DependencySourceMap, error) {
	var wg sync.WaitGroup

	dependSettings := &JSDependencyTree{
		WebDir:   webDirPath,
		JsParser: &jsparse.JSFileParser{},
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
