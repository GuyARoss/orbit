// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package jsparse

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/GuyARoss/orbit/pkg/fsutils"
	"github.com/google/uuid"
)

// FunctionDefinition is a light-weight struct to define JS function definition
type FunctionDefinition struct {
	Content    string
	Name       string
	IsExported bool
}

// ImportType represents a javascript import type
type ImportType int32

const (
	// LocalImportType represents an import that appears to be a local module
	// e.g import Thing from '../stuff/help.js'
	LocalImportType ImportType = 0

	// ModuleImportType represents an import that appears to be located in
	// as a node module. e.g import Thing from '@someorg/help.js'
	ModuleImportType ImportType = 1
)

// ImportDependency represents an entire import path within a javascript file.
type ImportDependency struct {
	FinalStatement string
	InitialPath    string
	Type           ImportType
}

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
}

// DefaultJSDocument is a struct that implements the JSDocument interface
// this struct can be used as an output for JSDocument parsing.
type DefaultJSDocument struct {
	imports []*ImportDependency
	name    string
	other   []string

	webDir    string
	pageDir   string
	extension string
}

// JSToken is some keyword(s) found in javascript used to tokenize js documents.
type JSToken string

const (
	ImportToken JSToken = "import"
	ExportToken JSToken = "export default"
)

var declarationTokens = []JSToken{ImportToken, ExportToken}

var ErrFunctionExport = errors.New("function export cannot be the name of the default export")

// extractDefaultExportName finds and returns an export name
// (if applicable) found within the provided line.
func extractDefaultExportName(line string) (string, error) {
	exportData := strings.Split(line, string(ExportToken))
	if len(exportData[1]) == 0 {
		return "", ErrFunctionExport
	}

	possibleName := strings.Trim(exportData[1][1:], " ")

	if !unicode.IsLetter(rune(possibleName[0])) {
		return "", nil
	}

	return possibleName, nil
}

// subsetRune returns a string subset found within two runes (subStart & subEnd)
func subsetRune(str string, subStart rune, subEnd rune) string {
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

// pageExtension attempts to determine the provided strings (importPath) file extension
// this can be used to determine an extension when one is not present on the file path
// if one is present, it returns in empty string, rather than the one present on the line.
func pageExtension(importPath string) string {
	split := strings.Split(importPath, ".")

	if len(split) > 1 {
		// an extension is already present on the resource.
		return split[len(split)-1]
	}

	extension := "js"
	// todo(issue/#11): context should be used here to pass in a "defaultExtension" type
	// provided by the pages web wrapper method.
	_, err := os.Stat(fmt.Sprintf("%s.js", importPath))
	if err != nil {
		extension = "jsx"
	}
	return extension
}

// pathToken finds the first valid JS path token (" or ') within the line
// if no valid path token is found \u002 is used by default.
func pathToken(line string) rune {
	for i := len(line) - 1; i > 0; i-- {
		if string(line[i]) == `'` {
			return rune('\u0027')
		}

		if string(line[i]) == `"` {
			return rune('\u0022')
		}
	}

	// @@todo an exception could be raised here in the case that a suitable path token is not found.
	return rune('\u0022')
}

// lineImportType finds the valid ImportType provided a valid import line
func lineImportType(line string) ImportType {
	pathToken := pathToken(line)
	path := subsetRune(line, rune(pathToken), rune(pathToken))

	if path[1] == '.' || path[1] == '/' {
		return LocalImportType
	}

	return ModuleImportType
}

// NewEmptyDocument creates a new empty JSDocument
func NewEmptyDocument() *DefaultJSDocument {
	return &DefaultJSDocument{}
}

// NewDocument creates a new JS document
func NewDocument(webDir string, pageDir string) *DefaultJSDocument {
	return &DefaultJSDocument{
		webDir:    webDir,
		pageDir:   pageDir,
		extension: pageExtension(pageDir),
	}
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
	extension := pageExtension(finalPath)

	newPath := fsutils.NormalizePath(fmt.Sprintf("'../../../%s.%s'", strings.Join(cleanWebDirPaths, "/"), extension))
	statementWithoutPath := strings.Replace(line, fmt.Sprintf("%c%s%c", pathChar, path, pathChar), newPath, 1)

	return &ImportDependency{
		FinalStatement: statementWithoutPath,
		InitialPath:    fsutils.NormalizePath(fmt.Sprintf("%s.%s", strings.Join(cleanWebDirPaths, "/"), extension)),
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

type EmptyParser struct {
	BadParse bool
}

func (p *EmptyParser) Parse(string, string) (JSDocument, error) {
	return &DefaultJSDocument{}, nil
}

type JSFileParser struct{}

func (p *JSFileParser) Parse(pageDir string, webDir string) (JSDocument, error) {
	if len(pageDir) >= 2 && pageDir[0:2] != fsutils.NormalizePath("./") {
		pageDir = fsutils.NormalizePath(fmt.Sprintf("./%s", pageDir))
	}

	file, err := os.Open(fsutils.NormalizePath(pageDir))
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
