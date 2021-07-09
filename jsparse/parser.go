package jsparse

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type FunctionDefinition struct {
	Content    string
	Name       string
	IsExported bool
}

type Page struct {
	Imports []string
	Name    string
	Other   []string
}

type JSToken string

const (
	ImportToken JSToken = "import"
	ExportToken JSToken = "export default"
)

var declarationTokens = []JSToken{ImportToken, ExportToken}

func cleanExportDefaultName(line string) string {
	exportData := strings.Split(line, "export default")
	return exportData[1][1:]
}

func (p *Page) tokenizeLine(line string) {
	skip := false
	for _, decToken := range declarationTokens {
		if strings.Contains(line, string(decToken)) {
			switch decToken {
			case ImportToken:
				p.Imports = append(p.Imports, line)
				skip = true
			case ExportToken:
				p.Name = cleanExportDefaultName(line)
				skip = true
			}
		}
	}

	if !skip {
		p.Other = append(p.Other, line)
	}
}

func (p *Page) WriteFile(dir string) {
	out := strings.Builder{}
	for _, imp := range p.Imports {
		out.WriteString(fmt.Sprintf("%s\n", imp))
	}

	for _, other := range p.Other {
		out.WriteString(fmt.Sprintf("%s\n", other))
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

func ParsePage(pageDir string) (*Page, error) {
	file, err := os.Open(pageDir)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	p := &Page{}

	for scanner.Scan() {
		p.tokenizeLine(scanner.Text())
	}

	return p, nil
}
