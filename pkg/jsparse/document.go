// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package jsparse

import (
	"fmt"
	"os"
	"strings"

	"github.com/GuyARoss/orbit/pkg/fsutils"
	"github.com/google/uuid"
)

// JSDocument is an interface that describes the behavior of a JSDocument
type JSDocument interface {
	WriteFile(string) error
	Key() string
	Imports() []*ImportDependency
	AddImport(*ImportDependency) []*ImportDependency
	Other() []string
	AddOther(string) []string
	Extension() string
	AddSerializable(s JSSerialize)
	Name() string
	DefaultExport() *JsDocumentScope
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
}

// JSToken is some keyword(s) found in javascript used to tokenize js documents.
type JSToken string

const (
	ImportToken        JSToken = "import"
	ExportDefaultToken JSToken = "export default"
	ExportConstToken   JSToken = "export const"
	ConstToken         JSToken = "const"
	FuncToken          JSToken = "function"
	VarToken           JSToken = "var"
	LetToken           JSToken = "let"
)

var declarationTokens = []JSToken{VarToken, ConstToken, FuncToken, LetToken, ExportDefaultToken, ImportToken}
var exportTokens = []JSToken{ExportDefaultToken, ExportConstToken}

type JsDocumentScope struct {
	TokenType JSToken
	Name      string
	Export    JSExport
	Args      JSDocArgList
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
	extension := pageExtension(verifyPath(finalPath))

	finalPath = strings.ReplaceAll(finalPath, fmt.Sprintf(".%s", extension), "")

	newPath := fsutils.NormalizePath(fmt.Sprintf("'../../../%s.%s'", finalPath, extension))
	statementWithoutPath := strings.Replace(line, fmt.Sprintf("%c%s%c", pathChar, path, pathChar), newPath, 1)

	initialPath := strings.ReplaceAll(strings.Join(cleanWebDirPaths, "/"), fmt.Sprintf(".%s", extension), "")

	return &ImportDependency{
		FinalStatement: statementWithoutPath,
		InitialPath:    fsutils.NormalizePath(fmt.Sprintf("%s.%s", initialPath, extension)),
		Type:           importType,
	}
}

// tokenizeLine tokenizes each line and serializes it to the provided JSDocument
func (p *DefaultJSDocument) tokenizeLine(line string) error {
	// @@todo(guy): replace this process with a ll instead of reading from the lines directly
	for _, decToken := range declarationTokens {
		if strings.Contains(line, string(decToken)) {
			switch decToken {
			case ImportToken:
				p.imports = append(p.imports, p.formatImportLine(line))
				return nil
			case ExportDefaultToken, ExportConstToken:
				// does the name already exist?
				name, err := extractJSTokenName(line, decToken)

				if err != nil {
					return err
				}

				if p.scope[name] != nil {
					if decToken == ExportDefaultToken {
						p.defaultExport = p.scope[name]
						p.name = name
					}
					return nil
				}
			}

			name, err := extractJSTokenName(line, decToken)
			if err != nil {
				return err
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
				return err
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
				return nil
			}
		}
	}

	p.AddOther(line)

	return nil
}

func (p *DefaultJSDocument) WriteFile(dir string) error {
	out := strings.Builder{}
	for _, imp := range p.imports {
		out.WriteString(fsutils.NormalizePath(fmt.Sprintf("%s\n", imp.FinalStatement)))
	}

	for _, other := range p.Other() {
		out.WriteString(fsutils.NormalizePath(fmt.Sprintf("%s\n", other)))
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
		// @@ return error that a key is not present
		return ""
	}

	id := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(p.defaultExport.Name))

	return strings.ReplaceAll(id.String(), "-", "")
}

func (p *DefaultJSDocument) AddImport(dependency *ImportDependency) []*ImportDependency {
	p.imports = append(p.imports, dependency)

	return p.imports
}

func (p *DefaultJSDocument) AddOther(new string) []string {
	p.other = append(p.other, new)

	return p.other
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
	varname      string
	valueBodyMap map[string]jsSwitchValue
}

func (s *JsDocSwitch) Serialize() string {
	w := strings.Builder{}

	w.WriteString(fmt.Sprintf("switch (%s) {", s.varname))

	for _, v := range s.valueBodyMap {
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
		varname:      varname,
		valueBodyMap: make(map[string]jsSwitchValue),
	}
}

type JSType string

const (
	JSString JSType = "string"
	JSNumber JSType = "number"
)

func (s *JsDocSwitch) Add(t JSType, value string, body string) {
	s.valueBodyMap[value] = jsSwitchValue{
		Value:  value,
		JSType: t,
		Body:   body,
	}
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
