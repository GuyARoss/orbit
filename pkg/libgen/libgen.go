package libgen

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/GuyARoss/orbit/internal/srcpack"
	webwrapper "github.com/GuyARoss/orbit/pkg/web_wrapper"
)

type LibOut struct {
	PackageName string
	Body        string
}

func (l *LibOut) WriteFile(dir string) error {
	out := strings.Builder{}
	out.WriteString(fmt.Sprintf("package %s\n\n", l.PackageName))
	out.WriteString(l.Body)

	if _, err := os.Stat(dir); err != nil {
		_, cerr := os.Create(dir)
		if cerr != nil {
			return cerr
		}
	}

	f, err := os.OpenFile(dir, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)

	if err != nil {
		return err
	}
	defer f.Close()

	err = f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(f, "%s", out.String())
	if err != nil {
		return err
	}

	return nil
}

type page struct {
	name        string
	bundleKey   string
	wrapVersion string
}

type BundleGroup struct {
	*BundleGroupOpts

	pages []*page
	// compRequiredBody map[string][]string

	// component web wrapper body data
	compww map[string][]string
}

type BundleGroupOpts struct {
	PackageName   string
	BaseBundleOut string
	BundleMode    string
	PublicDir     string
}

func New(opts *BundleGroupOpts) *BundleGroup {
	return &BundleGroup{
		BundleGroupOpts: opts,
		pages:           make([]*page, 0),
		compww:          make(map[string][]string),
	}
}

func parseVersionKey(k string) string {
	// @@todo: check if first char is int
	f := strings.ReplaceAll(k, ".", "_")
	return strings.ReplaceAll(f, "-", "")
}

func (l *BundleGroup) AcceptComponents(ctx context.Context, comps []*srcpack.Component, cacheOpts *webwrapper.CacheDOMOpts) {
	for _, c := range comps {
		v := parseVersionKey(c.WebWrapper.Version())

		l.pages = append(l.pages, &page{c.Name, c.BundleKey, v})
		l.compww[v] = c.WebWrapper.RequiredBodyDOMElements(ctx, cacheOpts)
	}
}

func (l *BundleGroup) CreateBundleLib() *LibOut {
	out := strings.Builder{}

	for rd, v := range l.compww {
		out.WriteString("\n")
		out.WriteString(fmt.Sprintf(`var %s = []string{`, rd))
		out.WriteString("\n")

		for _, b := range v {
			out.WriteString(fmt.Sprintf("`%s`,", b))
			out.WriteString("\n")
		}

		out.WriteString("}")
		out.WriteString("\n")
	}

	if len(l.BaseBundleOut) > 0 {
		out.WriteString("\n")
		out.WriteString(fmt.Sprintf(`var bundleDir string = "%s"`, l.BaseBundleOut))
		out.WriteString("\n")
	}

	if len(l.PublicDir) > 0 {
		out.WriteString("\n")
		out.WriteString(fmt.Sprintf(`var publicDir string = "%s"`, l.PublicDir))
		out.WriteString("\n")
	}

	out.WriteString("type PageRender string\n\n")

	for idx, p := range l.pages {
		if idx == 0 {
			out.WriteString("const ( \n")
		}

		if !strings.Contains(p.name, "Page") {
			p.name = fmt.Sprintf("%sPage", p.name)
		}

		out.WriteString(fmt.Sprintf(`	%s PageRender = "%s"`, p.name, p.bundleKey))
		out.WriteString("\n")

		if idx == len(l.pages)-1 {
			out.WriteString(")\n")
		}
	}

	out.WriteString("\n")
	out.WriteString(`var wrapBody = map[PageRender][]string{`)
	out.WriteString("\n")

	for _, p := range l.pages {
		out.WriteString(fmt.Sprintf(`	%s: %s,`, p.name, p.wrapVersion))
		out.WriteString("\n")
	}

	out.WriteString("}")

	out.WriteString("\n")

	out.WriteString(`
type BundleMode int32

const (
	DevBundleMode  BundleMode = 0
	ProdBundleMode BundleMode = 1
)

`)

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
