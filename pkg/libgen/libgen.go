package libgen

import (
	"fmt"
	"log"
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
	fmt.Println("package name", l.PackageName)
	out.WriteString(fmt.Sprintf("package %s\n\n", l.PackageName))

	if len(l.BaseBundleOut) > 0 {
		out.WriteString(fmt.Sprintf(`var bundleDir string = "%s"`, l.BaseBundleOut))
		out.WriteString("\n\n")
	}

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

	if err != nil {
		fmt.Println("cannot open file correctly", dir)
		log.Fatal(err)
	}
	defer f.Close()

	err = f.Truncate(0)
	if err != nil {
		fmt.Println("cannot truncate file correctly", dir)
		log.Fatal(err)
	}
	_, err = fmt.Fprintf(f, "%s", out.String())
	if err != nil {
		fmt.Println("uhh something stupid.", dir)
		log.Fatal(err)
	}
}
