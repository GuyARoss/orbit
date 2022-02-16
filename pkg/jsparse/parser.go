package jsparse

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
)

// FunctionDefinition
// light-weight struct to define JS function definition
type FunctionDefinition struct {
	Content    string
	Name       string
	IsExported bool
}

// ImportType
// Defines type of import with binary determination
type ImportType int32

const (
	LocalImportType  ImportType = 0
	ModuleImportType ImportType = 1
)

type ImportDependency struct {
	FinalStatement string
	InitialPath    string
	Type           ImportType
}

type JSDocument interface {
	WriteFile(string) error
	Key() string
	Name() string
	Imports() []*ImportDependency
	AddImport(*ImportDependency) []*ImportDependency
	Other() []string
	AddOther(string) []string
	Extension() string
}

type DefaultJSDocument struct {
	imports []*ImportDependency
	name    string
	other   []string

	webDir  string
	pageDir string
}

type JSToken string

const (
	ImportToken JSToken = "import"
	ExportToken JSToken = "export default"
)

var declarationTokens = []JSToken{ImportToken, ExportToken}

var ErrFunctionExport = errors.New("function export cannot be the name of the default export")

var ErrExportNotCapitalized = errors.New("default export of component should be capitalized")

func extractDefaultExportName(line string) (string, error) {
	exportData := strings.Split(line, string(ExportToken))
	possibleName := strings.Trim(exportData[1][1:], " ")

	if string(possibleName[0]) != strings.ToUpper(string(possibleName[0])) {
		return "", ErrExportNotCapitalized
	}

	if len(strings.Split(possibleName, " ")) > 1 {
		return "", ErrFunctionExport
	}

	return possibleName, nil
}

func filterCenter(str string, subStart rune, subEnd rune) string {
	final := make([]rune, 0)

	started := false
	for _, c := range str {
		if started && c == subEnd {
			return string(final)
		}

		if started {
			final = append(final, c)
		}

		if !started && c == subStart {
			started = true
		}
	}

	return string(final)
}

func pageExtension(importPath string) string {
	split := strings.Split(importPath, ".")
	if len(split) > 1 {
		return ""
	}

	extension := ".js"
	_, err := os.Stat(fmt.Sprintf("%s.js", importPath))
	if err != nil {
		extension = ".jsx"
	}
	return extension
}

func pathToken(line string) rune {
	for i := len(line) - 1; i > 0; i-- {
		if string(line[i]) == `'` {
			return rune('\u0027')
		}

		if string(line[i]) == `"` {
			return rune('\u0022')
		}
	}

	// @@todo raise an exception here
	return rune('\u0022')
}

func lineImportType(line string) ImportType {
	pathToken := pathToken(line)

	path := filterCenter(line, rune(pathToken), rune(pathToken))

	if path[1] == '.' || path[1] == '/' {
		return LocalImportType
	}

	return ModuleImportType
}

func NewEmptyDocument() *DefaultJSDocument {
	return &DefaultJSDocument{}
}

func (p *DefaultJSDocument) formatImportLine(line string) *ImportDependency {
	importType := lineImportType(line)
	if importType == ModuleImportType {
		return &ImportDependency{
			InitialPath:    line,
			FinalStatement: line,
			Type:           ModuleImportType,
		}
	}

	pathChar := '"'
	path := filterCenter(line, '"', '"')
	if len(path) == 0 {
		path = filterCenter(line, '\'', '\'')
		pathChar = '\''
	}

	webDirPaths := strings.Split(p.webDir, "/")
	cleanWebDirPaths := make([]string, 0)

	for _, dp := range webDirPaths {
		if len(dp) <= 1 && strings.Contains(dp, ".") {
			continue
		}

		cleanWebDirPaths = append(cleanWebDirPaths, dp)
	}

	tokenPathPaths := strings.Split(path, "/")
	hasProceedingDirectory := false
	for _, tk := range tokenPathPaths {
		if strings.Contains(tk, "..") {
			if hasProceedingDirectory {
				// @@todo(debug) throw error cuz this is out of range
			}
			hasProceedingDirectory = true
			continue
		}

		cleanWebDirPaths = append(cleanWebDirPaths, tk)
	}

	finalPath := strings.Join(cleanWebDirPaths, "/")
	extension := pageExtension(finalPath)

	newPath := fmt.Sprintf("'../../../%s%s'", strings.Join(cleanWebDirPaths, "/"), extension)
	statementWithoutPath := strings.Replace(line, fmt.Sprintf("%c%s%c", pathChar, path, pathChar), newPath, 1)

	return &ImportDependency{
		FinalStatement: statementWithoutPath,
		InitialPath:    fmt.Sprintf("%s%s", strings.Join(cleanWebDirPaths, "/"), extension),
		Type:           importType,
	}
}

func (p *DefaultJSDocument) tokenizeLine(line string) error {
	skip := false
	for _, decToken := range declarationTokens {
		if strings.Contains(line, string(decToken)) {
			switch decToken {
			case ImportToken:
				p.imports = append(p.imports, p.formatImportLine(line))

				skip = true
			case ExportToken:
				possibleName, err := extractDefaultExportName(line)
				if err != nil && !errors.Is(ErrFunctionExport, err) {
					return err
				}

				p.name = possibleName
				skip = true
			}
		}
	}

	if !skip {
		p.AddOther(line)
	}
	return nil
}

func (p *DefaultJSDocument) WriteFile(dir string) error {
	out := strings.Builder{}
	for _, imp := range p.imports {
		out.WriteString(fmt.Sprintf("%s\n", imp.FinalStatement))
	}

	for _, other := range p.Other() {
		out.WriteString(fmt.Sprintf("%s\n", other))
	}

	f, err := os.OpenFile(dir, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
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

func defaultPageName(pageDir string) string {
	// @@todo: validate that the DefaultJSDocument has a valid name,
	// if not, make one out of the pageDir
	basePageDir := strings.Split(pageDir, ".")[0]

	splitPath := strings.FieldsFunc(basePageDir, func(r rune) bool {
		return r == '_' || r == ' ' || r == '-'
	})

	for i, p := range splitPath {
		splitPath[i] = fmt.Sprintf("%s%s", strings.ToUpper(string(p[0])), p[1:])
	}

	return strings.Join(splitPath, "")
}

type JSParser interface {
	Parse(string, string) (JSDocument, error)
}

type JSFileParser struct{}

func (p *JSFileParser) Parse(pageDir string, webDir string) (JSDocument, error) {
	if pageDir[0:2] != "./" {
		pageDir = fmt.Sprintf("./%s", pageDir)
	}

	file, err := os.Open(pageDir)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	page := &DefaultJSDocument{
		webDir:  webDir,
		pageDir: pageDir,
	}

	for scanner.Scan() {
		err := page.tokenizeLine(scanner.Text())
		if err != nil {
			return nil, err
		}
	}

	if page.name == "" {
		page.name = defaultPageName(pageDir)
	}

	return page, nil
}

func (p *DefaultJSDocument) Key() string {
	id := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(p.name))

	return strings.ReplaceAll(id.String(), "-", "")
}

func (p *DefaultJSDocument) Name() string {
	return p.name
}

func (p *DefaultJSDocument) AddImport(dependency *ImportDependency) []*ImportDependency {
	p.imports = append(p.imports, dependency)

	return p.imports
}

func (p *DefaultJSDocument) Imports() []*ImportDependency {
	return p.imports
}

func (p *DefaultJSDocument) Other() []string {
	return p.other
}

func (p *DefaultJSDocument) AddOther(new string) []string {
	p.other = append(p.other, new)

	return p.other
}

func (p *DefaultJSDocument) Extension() string {
	f := strings.Split(p.pageDir, ".")

	if len(f) == 0 {
		return ""
	}

	return f[len(f)-1]
}
