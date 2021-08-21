package libgen

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type LibOut struct {
	PackageName string
	Body        string
}

func (l *LibOut) WriteFile(dir string) {
	out := strings.Builder{}
	out.WriteString(fmt.Sprintf("package %s\n\n", l.PackageName))
	out.WriteString(l.Body)

	f, err := os.OpenFile(dir, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = f.Truncate(0)
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fprintf(f, "%s", out.String())
	if err != nil {
		log.Fatal(err)
	}
}

type page struct {
	name      string
	bundleKey string
}

type BundleGroup struct {
	PackageName   string
	BaseBundleOut string
	BundleMode    string

	pages []*page
}

func (l *BundleGroup) ApplyBundle(name string, bundleKey string) {
	l.pages = append(l.pages, &page{name, bundleKey})
}

func (l *BundleGroup) CreateBundleLib() *LibOut {
	out := strings.Builder{}

	out.WriteString(`var hotReloadPipePath = ".orbit/hotreload"`)

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

	out.WriteString(`\n
type BundleMode int32

const (
	DevBundleMode  BundleMode = 0
	ProdBundleMode BundleMode = 1
)`)

	if l.BundleMode == "production" {
		out.WriteString("var CurrentDevMode BundleMode = ProdBundleMode")
	} else {
		out.WriteString("var CurrentDevMode BundleMode = DevBundleMode")
	}

	return &LibOut{
		PackageName: l.PackageName,
		Body:        out.String(),
	}
}

type StaticToken string

const (
	StartToken StaticToken = "// **__START_STATIC__**"
	EndToken   StaticToken = "// **__END_STATIC__**"
)

var declarationTokens = []StaticToken{StartToken, EndToken}

func ParseStaticFile(dir string) (string, error) {
	file, err := os.Open(dir)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	out := strings.Builder{}
	isStatic := false

	for scanner.Scan() {
		line := scanner.Text()

		skip := false

		for _, decToken := range declarationTokens {
			if strings.Contains(line, string(decToken)) {
				switch decToken {
				case StartToken:
					{
						skip = true
						isStatic = true
					}
				case EndToken:
					isStatic = false
				}

				continue
			}
		}
		if isStatic && !skip {
			out.WriteString(fmt.Sprintf("%s\n", line))
		}
	}

	return out.String(), nil
}
