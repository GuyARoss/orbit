// Copyright (c) 2021 Guy A. Ross
// This source code is licensed under the GNU GPLv3 found in the
// license file in the root directory of this source tree.

package jsparse

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// pageExtension attempts to determine the provided strings (importPath) file extension
// this can be used to determine an extension when one is not present on the file path
// if one is present, it returns in empty string, rather than the one present on the line.
func pageExtension(importPath string) string {
	split := strings.Split(importPath, ".")

	if len(split) > 2 {
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

var ErrFunctionExport = errors.New("function export cannot be the name of the default export")
var ErrInvalidName = errors.New("variable name is invalid")

// extractDefaultExportName finds and returns an export name
// (if applicable) found within the provided line.
func extractDefaultExportName(line string) (string, error) {
	exportData := strings.Split(line, string(ExportDefaultToken))
	if len(exportData[1]) == 0 {
		return "", ErrFunctionExport
	}

	possibleName := strings.Trim(exportData[1][1:], " ")

	if !unicode.IsLetter(rune(possibleName[0])) {
		return "", nil
	}

	return possibleName, nil
}

func parseArgs(line string) (JSDocArgList, error) {
	isInCtx := false
	args := make(JSDocArgList, 0)

	vname := []rune{}
	anonLen := 0
	inAnonCtx := false
	for _, c := range line {
		if !isInCtx && c == '(' {
			isInCtx = true
		}

		if !isInCtx && c != '(' {
			continue
		}

		if c == ')' {
			isInCtx = false
		}

		if c == '{' {
			inAnonCtx = true
		}
		if inAnonCtx && c == '}' {
			inAnonCtx = false
			args = append(args, fmt.Sprintf("anon_%d", anonLen))
			anonLen += 1
		}

		if unicode.IsLetter(c) {
			vname = append(vname, c)
			continue
		}

		if len(vname) > 0 && unicode.IsNumber(c) {
			vname = append(vname, c)
			continue
		}

		if len(vname) > 0 {
			args = append(args, string(vname))
			vname = []rune{}
		}

		if inAnonCtx && c != '{' {
			inAnonCtx = false
			args = append(args, fmt.Sprintf("anon_%d", anonLen))
			anonLen += 1
		}
	}

	return args, nil
}

func extractJSTokenName(line string, token JSToken) (string, error) {
	exportData := strings.Split(line, string(token))
	if len(exportData) == 0 {
		return "", nil
	}

	if len(exportData[1]) == 0 {
		return "", ErrFunctionExport
	}

	valid := []rune{}
	for _, k := range exportData[1][1:] {
		// @@todo validate i of 0 to ensure it is letter
		if unicode.IsLetter(k) || unicode.IsNumber(k) {
			valid = append(valid, k)
			continue
		}
		break
	}

	return string(valid), nil
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

// lineImportType finds the valid ImportType provided a valid import line
func lineImportType(line string) ImportType {
	pathToken := pathToken(line)
	path := subsetRune(line, rune(pathToken), rune(pathToken))

	if path[0] == '.' || path[0] == '/' {
		return LocalImportType
	}

	if path[1] == '.' || path[1] == '/' {
		return LocalImportType
	}

	return ModuleImportType
}
