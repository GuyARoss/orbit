package dependtree

type DependencyTreeNode struct {
	Value string
	Right *DependencyTreeNode
	Child *DependencyTreeNode
}

type DependencyTree interface {
	// DirList fetches the current paths within a directory
	DirList(path string) ([]string, error)

	// PathDependencies finds the given dependencies for a given path
	PathDependencies(path string) ([]string, error)
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

		current.Value = d
		current.Child = mapResp
		current.Right = &DependencyTreeNode{}

		current = current.Right
	}
	return root, nil
}

func Create(s DependencyTree, initialPath string) (*DependencyTreeNode, error) {
	dirs, err := s.DirList(initialPath)
	if err != nil {
		return nil, err
	}
	root := &DependencyTreeNode{}
	current := root
	for _, d := range dirs {
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
