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

	"github.com/GuyARoss/orbit/pkg/fsutils"
)

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
		webDir:    webDir,
		pageDir:   pageDir,
		extension: pageExtension(pageDir),
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

// JSToken is some keyword(s) found in javascript used to tokenize js documents.
type JSToken string

const (
	ImportToken JSToken = "import"
	ExportToken JSToken = "export default"
)

var declarationTokens = []JSToken{ImportToken, ExportToken}

var ErrFunctionExport = errors.New("function export cannot be the name of the default export")

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
