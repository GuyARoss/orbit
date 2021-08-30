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

	webDir string
}

type JSToken string

const (
	ImportToken JSToken = "import"
	ExportToken JSToken = "export default"
)

var declarationTokens = []JSToken{ImportToken, ExportToken}

func cleanExportDefaultName(line string) string {
	// @@todo(guy): detect other types of exporting.
	// @@todo(guy): validate that export type is capitalized

	exportData := strings.Split(line, "export default")
	return exportData[1][1:]
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

func (p *Page) formatImportLine(line string) string {
	// local imports "should" always include path information
	if !strings.Contains(line, "/") {
		return line
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

	return statementWithoutPath
}

func (p *Page) tokenizeLine(line string) {
	skip := false
	for _, decToken := range declarationTokens {
		if strings.Contains(line, string(decToken)) {
			switch decToken {
			case ImportToken:
				p.Imports = append(p.Imports, p.formatImportLine(line))

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

func (p *Page) WriteFile(dir string) error {
	out := strings.Builder{}
	for _, imp := range p.Imports {
		out.WriteString(fmt.Sprintf("%s\n", imp))
	}

	for _, other := range p.Other {
		out.WriteString(fmt.Sprintf("%s\n", other))
	}

	// if _, err := os.Stat(dir); os.IsNotExist(err) {
	// 	fmt.Println("uhh file doesn't exist?")
	// 	err := os.Mkdir(dir, 0755)
	// 	return err
	// }

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

func ParsePage(pageDir string, webDir string) (*Page, error) {
	file, err := os.Open(pageDir)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	p := &Page{
		webDir: webDir,
	}

	for scanner.Scan() {
		p.tokenizeLine(scanner.Text())
	}

	return p, nil
}
