// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package jsparse

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/google/uuid"
)

// JSDocument is an interface that describes the behavior of a JSDocument
type JSDocument interface {
	WriteFile(string) error
	Key() string
	Imports() []*ImportDependency
	AddImport(*ImportDependency) []*ImportDependency
	Other() []string
	// AddOther(string) []string
	AddOther(...string)
	Extension() string
	AddSerializable(s JSSerialize)
	Name() string
	DefaultExport() *JsDocumentScope
	Clone() JSDocument
}

// DefaultJSDocument is a struct that implements the JSDocument interface
// this struct can be used as an output for JSDocument parsing.
type DefaultJSDocument struct {
	imports      []*ImportDependency
	other        []string
	serializable []JSSerialize

	webDir    string
	pageDir   string
	extension string
	scope     map[string]*JsDocumentScope

	defaultExport *JsDocumentScope
	name          string
	inDeadBlock   bool
}

// JSToken is some tokens found in javascript used to tokenize js statements.
type JSToken string

const (
	ImportToken        JSToken = "import"
	ExportDefaultToken JSToken = "export default"
	ExportConstToken   JSToken = "export const"
	ConstToken         JSToken = "const"
	FuncToken          JSToken = "function"
	VarToken           JSToken = "var"
	LetToken           JSToken = "let"
	CommentToken       JSToken = "//"
	DoubleQuoteToken   JSToken = `"`
	SingleQuoteToken   JSToken = "'"
	MultiStringToken   JSToken = "`"
)

var declarationTokens = []JSToken{VarToken, ConstToken, FuncToken, LetToken, ExportDefaultToken, ImportToken}
var exportTokens = []JSToken{ExportDefaultToken, ExportConstToken}
var stringTokens = []JSToken{DoubleQuoteToken, SingleQuoteToken, MultiStringToken}

type JsDocumentScope struct {
	TokenType JSToken
	Name      string
	Export    JSExport
	Args      JSDocArgList
}

func removeCenterOfToken(line string, token string) (string, int) {
	parsedLine := line

	startIdx := 0
	opened := false
	foundCount := 0
	for idx, c := range line {
		if string(c) == string(token) {
			foundCount += 1
			if !opened {
				startIdx = idx
				opened = true
			} else {
				opened = false
				subset := line[startIdx+1 : idx]
				parsedLine = strings.ReplaceAll(parsedLine, subset, "")
			}
		}
	}

	return parsedLine, foundCount
}

func (p *DefaultJSDocument) Clone() JSDocument {
	return &DefaultJSDocument{
		imports:       p.imports,
		other:         p.other,
		serializable:  p.serializable,
		webDir:        p.webDir,
		pageDir:       p.pageDir,
		extension:     p.extension,
		scope:         p.scope,
		defaultExport: p.defaultExport,
		name:          p.name,
		inDeadBlock:   p.inDeadBlock,
	}
}

// tokenizeLine tokenizes each line and serializes it to the provided JSDocument
func (p *DefaultJSDocument) tokenizeLine(ctx context.Context, pageDir string, line string) (context.Context, error) {
	if len(line) == 0 {
		return ctx, nil
	}

	parsedLine := line
	for _, t := range stringTokens {
		out, v := removeCenterOfToken(parsedLine, string(t))
		if t == MultiStringToken && v%2 == 1 {
			p.inDeadBlock = !p.inDeadBlock
		}

		parsedLine = out
	}

	if p.inDeadBlock {
		p.AddOther(line)

		return ctx, nil
	}

	if strings.Contains(parsedLine, string(CommentToken)) {
		// the only part of the comment line that is valid would be everything before the comment
		commentDelimited := strings.Split(line, string(CommentToken))
		return p.tokenizeLine(ctx, pageDir, commentDelimited[0])
	}

	for _, decToken := range declarationTokens {
		if strings.Contains(parsedLine, string(decToken)) {
			switch decToken {
			case ImportToken:
				p.imports = append(p.imports, p.formatImportLine(line))
				return ctx, nil
			case ExportDefaultToken, ExportConstToken:
				name, err := extractJSTokenName(line, decToken)

				if err != nil {
					return ctx, err
				}

				if p.scope[name] != nil {
					if decToken == ExportDefaultToken {
						p.defaultExport = p.scope[name]
						p.name = name
					}
					return ctx, nil
				} else {
					v, err := p.parseInformalExportDefault(pageDir, line)
					if err != nil || v {
						return ctx, err
					}
				}
			}

			name, err := extractJSTokenName(line, decToken)
			if err != nil {
				return ctx, err
			}
			exportMethod := ExportNone

			isDefault := false
			for _, e := range exportTokens {
				if strings.Contains(line, string(e)) {
					switch e {
					case ExportDefaultToken:
						exportMethod = ExportDefault
						isDefault = true
					case ExportConstToken:
						exportMethod = ExportConst
					}
				}
			}

			args, err := parseArgs(line)
			if err != nil {
				return ctx, err
			}

			scope := &JsDocumentScope{
				Name:      name,
				Export:    exportMethod,
				TokenType: decToken,
				Args:      args,
			}

			p.scope[name] = scope

			if isDefault {
				if decToken == ExportDefaultToken {
					p.defaultExport = p.scope[name]
					p.name = name
				}
				return ctx, nil
			}
		}
	}

	p.AddOther(line)

	return ctx, nil
}

func formatPathToPageName(path string) string {
	splitPath := strings.Split(path, string(os.PathSeparator))
	pageName := strings.Split(splitPath[len(splitPath)-1], ".")[0]

	caseTypes := []string{"_", "-", " "} // snake case, kebab case, monster case
	for _, c := range caseTypes {
		if strings.Contains(pageName, c) {
			final := ""
			for _, splitPageName := range strings.Split(pageName, c) {
				final += strings.ToUpper(splitPageName[:1]) + strings.ReplaceAll(splitPageName[1:], c, "")
			}
			return final
		}
	}

	pageName = strings.ToUpper(pageName[:1]) + pageName[1:]

	return pageName
}

func (p *DefaultJSDocument) parseInformalExportDefault(pageDir string, line string) (bool, error) {
	exportData := strings.TrimSpace(strings.Split(line, string(ExportDefaultToken))[1])
	if len(exportData) == 0 {
		return false, nil
	}
	if match := regexp.MustCompile("\\{|\\}|\\[|\\]|\\(|\\)"); !match.Match([]byte(exportData)) {
		return false, nil
	}

	pageName := formatPathToPageName(pageDir)

	p.AddOther(fmt.Sprintf("const %s = %s", pageName, exportData))
	p.scope[pageName] = &JsDocumentScope{
		Name:      pageName,
		Export:    ExportDefault,
		TokenType: ConstToken,
		Args:      make(JSDocArgList, 0),
	}

	p.defaultExport = p.scope[pageName]
	p.name = pageName

	return true, nil
}

func isIndexPath(finalPath string) bool {
	// first we check if it can be found locally
	stat, err := os.Stat("./" + finalPath)
	if err == nil && stat.IsDir() {
		return true
	}

	if err != nil {
		// next, if try to find it absolutely
		stat, err := os.Stat(finalPath)
		if err == nil && stat.IsDir() {
			return true
		}
	}

	return false
}

// formatImportLine parses an import line to create an import dependency
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
	path := subsetRune(line, '"', '"')
	if len(path) == 0 {
		path = subsetRune(line, '\'', '\'')
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
	for _, tk := range tokenPathPaths {
		// tk = "." should provide support for localized paths such as "./"
		if tk == "." {
			validPageDir := p.pageDir

			if p.pageDir[:2] == "./" {
				validPageDir = p.pageDir[2:]
			}

			pageDirs := strings.Split(validPageDir, "/")
			cleanWebDirPaths = pageDirs[0 : len(pageDirs)-1]

			continue
		}

		if strings.Contains(tk, "..") {
			continue
		}

		cleanWebDirPaths = append(cleanWebDirPaths, tk)
	}

	finalPath := strings.Join(cleanWebDirPaths, "/")

	// possible that the path is referencing an index
	// we can validate this by checking if the import path is a dir
	if isIndexPath(finalPath) {
		finalPath += "/index"
		cleanWebDirPaths = append(cleanWebDirPaths, "/index")
	}

	extension := pageExtension(verifyPath(finalPath))

	finalPath = strings.ReplaceAll(finalPath, fmt.Sprintf(".%s", extension), "")
	newPath := fmt.Sprintf("'../../../%s.%s'", finalPath, extension)

	statementWithoutPath := strings.Replace(line, fmt.Sprintf("%c%s%c", pathChar, path, pathChar), newPath, 1)

	initialPath := strings.ReplaceAll(strings.Join(cleanWebDirPaths, "/"), fmt.Sprintf(".%s", extension), "")

	return &ImportDependency{
		FinalStatement: statementWithoutPath,
		InitialPath:    fmt.Sprintf("%s.%s", initialPath, extension),
		Type:           importType,
	}
}

func (p *DefaultJSDocument) WriteFile(dir string) error {
	out := strings.Builder{}
	for _, imp := range p.imports {
		out.WriteString(fmt.Sprintf("%s\n", imp.FinalStatement))
	}

	for _, other := range p.other {
		out.WriteString(fmt.Sprintf("%s\n", other))
	}

	for _, s := range p.serializable {
		out.WriteString(s.Serialize())
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

func (p *DefaultJSDocument) DefaultExport() *JsDocumentScope {
	return p.defaultExport
}

func (p *DefaultJSDocument) Name() string {
	return p.name
}

func (p *DefaultJSDocument) Other() []string              { return p.other }
func (p *DefaultJSDocument) Imports() []*ImportDependency { return p.imports }
func (p *DefaultJSDocument) Extension() string            { return p.extension }

func (p *DefaultJSDocument) Key() string {
	if p.defaultExport == nil {
		return ""
	}

	id := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(p.defaultExport.Name))

	return strings.ReplaceAll(id.String(), "-", "")
}

func (p *DefaultJSDocument) AddImport(dependency *ImportDependency) []*ImportDependency {
	p.imports = append(p.imports, dependency)

	return p.imports
}

func (p *DefaultJSDocument) AddOther(new ...string) {
	for _, n := range new {
		p.other = append(p.other, n)
	}
}

func (p *DefaultJSDocument) AddSerializable(s JSSerialize) {
	p.serializable = append(p.serializable, s)
}

// NewEmptyDocument creates a new empty JSDocument
func NewEmptyDocument() *DefaultJSDocument {
	return &DefaultJSDocument{
		scope:         make(map[string]*JsDocumentScope),
		imports:       make([]*ImportDependency, 0),
		defaultExport: &JsDocumentScope{},
		other:         make([]string, 0),
	}
}

func NewImportDocument(imports ...*ImportDependency) *DefaultJSDocument {
	doc := NewEmptyDocument()

	for _, i := range imports {
		doc.AddImport(i)
	}
	return doc
}

// NewDocument creates a new JS document
func NewDocument(webDir string, pageDir string) *DefaultJSDocument {
	return &DefaultJSDocument{
		webDir:        webDir,
		pageDir:       pageDir,
		extension:     pageExtension(pageDir),
		serializable:  make([]JSSerialize, 0),
		scope:         make(map[string]*JsDocumentScope),
		imports:       make([]*ImportDependency, 0),
		defaultExport: &JsDocumentScope{},
	}
}

type JSSerialize interface {
	Serialize() string
}

type jsSwitchValue struct {
	Value  string
	JSType JSType
	Body   string
}

type JsDocSwitch struct {
	varname string

	varChecker map[string]bool
	valueBody  []jsSwitchValue

	m sync.Mutex
}

func (s *JsDocSwitch) Serialize() string {
	w := strings.Builder{}

	w.WriteString(fmt.Sprintf("switch (%s) {", s.varname))

	for _, v := range s.valueBody {
		var e string

		switch v.JSType {
		case JSNumber:
			e = fmt.Sprintf("%s", v.Value)
		default:
			e = fmt.Sprintf("'%s'", v.Value)
		}

		w.WriteString(fmt.Sprintf(`case %s: { %s }`, e, v.Body))
	}

	w.WriteString("}")

	return w.String()
}

func NewSwitch(varname string) *JsDocSwitch {
	return &JsDocSwitch{
		varname:    varname,
		valueBody:  make([]jsSwitchValue, 0),
		varChecker: make(map[string]bool),
		m:          sync.Mutex{},
	}
}

type JSType string

const (
	JSString JSType = "string"
	JSNumber JSType = "number"
)

func (s *JsDocSwitch) Add(t JSType, value string, body string) {
	s.m.Lock()

	if !s.varChecker[value] {
		s.varChecker[value] = true
		s.valueBody = append(s.valueBody, jsSwitchValue{
			Value:  value,
			JSType: t,
			Body:   body,
		})
	}

	s.m.Unlock()
}

type JSExport int

const (
	ExportNone    JSExport = 0
	ExportDefault JSExport = 1
	ExportConst   JSExport = 2
)

type JSDocArgList []string

func (s JSDocArgList) ToString() string {
	return strings.Join(s, ",")
}

type JsDocFunc struct {
	Declaration string
	body        JSSerialize
}

func (s *JsDocFunc) Serialize() string {
	w := strings.Builder{}

	w.WriteString(s.Declaration)
	w.WriteString("{")
	w.WriteString(s.body.Serialize())
	w.WriteString("}")

	return w.String()
}

func NewFunc(declaration string, body JSSerialize) *JsDocFunc {
	return &JsDocFunc{
		Declaration: declaration,
		body:        body,
	}
}

func verifyPath(path string) string {
	extra := path[:2]

	if extra == ".." {
		return path
	}

	extra = strings.Replace(extra, ".", "", 1)
	extra = strings.Replace(extra, "/", "", 1)

	return fmt.Sprintf("%s%s%s", "./", extra, path[2:])
}
