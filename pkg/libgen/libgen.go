package libgen

import (
	"fmt"
	"os"
	"strings"
)

type page struct {
	name      string
	bundleKey string
}

type LibOut struct {
	pages         []*page
	PackageName   string
	BaseBundleOut string
}

func (l *LibOut) ApplyPage(name string, bundleKey string) {
	l.pages = append(l.pages, &page{name, bundleKey})
}

func (l *LibOut) WriteFile(dir string) {
	out := strings.Builder{}
	out.WriteString(fmt.Sprintf("package %s\n\n", l.PackageName))

	out.WriteString(fmt.Sprintf(`var bundleDir string = "%s"`, l.BaseBundleOut))
	out.WriteString("\n\n")

	for idx, p := range l.pages {
		if idx == 0 {
			out.WriteString("type PageRender string\n\n")
			out.WriteString("const ( \n")
		}

		out.WriteString(fmt.Sprintf(`	%sPage PageRender = "%s"`, p.name, p.bundleKey))
		out.WriteString("\n")

		if idx == len(l.pages)-1 {
			out.WriteString(")\n")
		}
	}

	f, err := os.OpenFile(dir, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	defer f.Close()

	if err != nil {
		panic(err)
	}

	err = f.Truncate(0)
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fprintf(f, "%s", out.String())
	if err != nil {
		panic(err)
	}
}
