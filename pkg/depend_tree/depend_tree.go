package dependtree

type DependencyTreeNode struct {
	Value string
	Right *DependencyTreeNode
	Child *DependencyTreeNode
}

func (d *DependencyTreeNode) values(current []string) []string {
	if d == nil || d.Value == "" {
		return current
	}

	if current == nil {
		current = make([]string, 0)
	}

	current = append(current, d.Value)
	current = d.Right.values(current)
	current = d.Child.values(current)

	return current
}

type DependencySourceMap struct {
	sourceMap map[string]string
}

func (d *DependencySourceMap) PathParent(path string) string {
	return d.sourceMap[path]
}

func (d *DependencyTreeNode) SourceMap() *DependencySourceMap {
	m := make(map[string]string)

	current := d
	for current != nil {
		values := d.values(nil)
		for _, v := range values {
			m[v] = d.Value
		}

		current = current.Right
	}

	return &DependencySourceMap{m}
}

type DependencySettings interface {
	// DirList fetches the current paths within a directory
	DirList(path string) ([]string, error)

	// PathDependencies finds the given dependencies for a given path
	PathDependencies(path string) ([]string, error)
}

type ManagedDependencyTree struct {
	rootNode *DependencyTreeNode
	Settings DependencySettings
}

func mapNode(s DependencySettings, path string) (*DependencyTreeNode, error) {
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

		current = current.Right
	}

	s.rootNode = root

	return root, nil
}
