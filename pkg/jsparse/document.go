// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package jsparse

import (
	"errors"
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
	Name() string
	Imports() []*ImportDependency
	AddImport(*ImportDependency) []*ImportDependency
	Other() []string
	AddOther(string) []string
	Extension() string
	AddSerializable(s JSSerialize)
}

type JSExport struct {
	IsDefault bool
	Name      string
	HasArgs   bool
}

// DefaultJSDocument is a struct that implements the JSDocument interface
// this struct can be used as an output for JSDocument parsing.
type DefaultJSDocument struct {
	imports      []*ImportDependency
	name         string
	other        []string
	serializable []JSSerialize

	webDir    string
	pageDir   string
	extension string
	exports   []JSExport
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
	for _, decToken := range declarationTokens {
		if strings.Contains(line, string(decToken)) {
			switch decToken {
			case ImportToken:
				p.imports = append(p.imports, p.formatImportLine(line))
				return nil
			case ExportToken:
				possibleName, err := extractDefaultExportName(line)
				if err != nil && !errors.Is(ErrFunctionExport, err) {
					return err
				}

				p.name = possibleName
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

func (p *DefaultJSDocument) Name() string { return p.name }

func (p *DefaultJSDocument) Other() []string              { return p.other }
func (p *DefaultJSDocument) Imports() []*ImportDependency { return p.imports }
func (p *DefaultJSDocument) Extension() string            { return p.extension }

func (p *DefaultJSDocument) Key() string {
	id := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(p.name))

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
	return &DefaultJSDocument{}
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
		webDir:       webDir,
		pageDir:      pageDir,
		extension:    pageExtension(pageDir),
		serializable: make([]JSSerialize, 0),
		exports:      make([]JSExport, 0),
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
